package knowledge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ZIMFile represents a discovered ZIM content file.
type ZIMFile struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	SizeBytes int64 `json:"size_bytes"`
	SizeHuman string `json:"size_human"`
}

// KiwixManager manages the kiwix-serve sidecar process lifecycle.
type KiwixManager struct {
	mu            sync.RWMutex
	dataDir       string
	sidecarPath   string
	port          int
	running       bool
	zimFiles      []ZIMFile
	cancelFunc    context.CancelFunc
	cmd           *exec.Cmd
	lastError     string
}

// NewKiwixManager creates a new KiwixManager.
func NewKiwixManager(dataDir string, sidecarPort int) *KiwixManager {
	return &KiwixManager{
		dataDir: dataDir,
		port:    sidecarPort,
		zimFiles: make([]ZIMFile, 0),
		lastError: "",
	}
}

// DiscoverZIMFiles scans the data directory and content-packs for .zim files.
func (km *KiwixManager) DiscoverZIMFiles() ([]ZIMFile, error) {
	km.mu.Lock()
	defer km.mu.Unlock()

	km.zimFiles = make([]ZIMFile, 0)

	// Search paths: data dir, content-packs directory, and zim subdirectory
	searchPaths := []string{
		km.dataDir,
		filepath.Join(filepath.Dir(km.dataDir), "content-packs"),
		filepath.Join(filepath.Dir(km.dataDir), "content-packs", "zim-survival-pack"),
		filepath.Join(km.dataDir, "zim"),
	}

	seen := make(map[string]bool)

	for _, searchPath := range searchPaths {
		err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // skip inaccessible directories
			}
			if info.IsDir() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".zim" {
				absPath, _ := filepath.Abs(path)
				if seen[absPath] {
					return nil
				}
				seen[absPath] = true

				id := strings.TrimSuffix(info.Name(), ext)
				km.zimFiles = append(km.zimFiles, ZIMFile{
					ID:        id,
					Name:      info.Name(),
					Path:      absPath,
					SizeBytes: info.Size(),
					SizeHuman: humanizeBytes(info.Size()),
				})
			}
			return nil
		})
		if err != nil {
			// Non-fatal: path may not exist
			continue
		}
	}

	log.Printf("  Knowledge: discovered %d ZIM files", len(km.zimFiles))
	return km.zimFiles, nil
}

// GetZIMFiles returns currently discovered ZIM files.
func (km *KiwixManager) GetZIMFiles() []ZIMFile {
	km.mu.RLock()
	defer km.mu.RUnlock()
	result := make([]ZIMFile, len(km.zimFiles))
	copy(result, km.zimFiles)
	return result
}

// FindSidecar looks for the kiwix-serve binary in the expected sidecar locations.
func (km *KiwixManager) FindSidecar() (string, error) {
	km.mu.Lock()
	defer km.mu.Unlock()

	candidates := []string{
		filepath.Join(filepath.Dir(km.dataDir), "sidecars", "kiwix-serve", "windows", "kiwix-serve.exe"),
		filepath.Join(filepath.Dir(km.dataDir), "sidecars", "kiwix-serve", "kiwix-serve.exe"),
		filepath.Join(filepath.Dir(km.dataDir), "kiwix-serve.exe"),
		"kiwix-serve", // try PATH
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			km.sidecarPath = path
			return path, nil
		}
	}

	return "", fmt.Errorf("kiwix-serve binary not found; place it in sidecars/kiwix-serve/windows/")
}

// Start launches the kiwix-serve process if it exists and ZIM files are found.
func (km *KiwixManager) Start() error {
	km.mu.Lock()
	defer km.mu.Unlock()

	if km.running {
		return nil
	}

	if km.sidecarPath == "" {
		return fmt.Errorf("kiwix-serve not found")
	}

	if len(km.zimFiles) == 0 {
		return fmt.Errorf("no ZIM files found to serve")
	}

	ctx, cancel := context.WithCancel(context.Background())
	km.cancelFunc = cancel

	// Build args: --port=9080 file1.zim file2.zim
	args := []string{fmt.Sprintf("--port=%d", km.port)}
	for _, zf := range km.zimFiles {
		args = append(args, zf.Path)
	}

	km.cmd = exec.CommandContext(ctx, km.sidecarPath, args...)
	
	// Start the process
	if err := km.cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("failed to start kiwix-serve: %w", err)
	}

	km.running = true
	log.Printf("✓ Knowledge module: started kiwix-serve on port %d with %d ZIM files", km.port, len(km.zimFiles))

	// Monitor process in background
	go func() {
		err := km.cmd.Wait()
		km.mu.Lock()
		km.running = false
		km.cmd = nil
		km.mu.Unlock()
		if err != nil && err.Error() != "signal: killed" {
			log.Printf("! Knowledge module: kiwix-serve exited with error: %v", err)
			
			km.mu.Lock()
			km.running = false
			km.cmd = nil
			// Detect Windows missing DLL error (0xc0000135)
			if strings.Contains(err.Error(), "0xc0000135") || strings.Contains(err.Error(), "3221225781") {
				km.lastError = "Missing Windows C++ Redistributable (libgcc_s_seh-1.dll, etc). Run download-kiwix.ps1 or install VC++ Redist."
				log.Printf("! Knowledge module: CRITICAL - " + km.lastError)
			} else {
				km.lastError = err.Error()
			}
			km.mu.Unlock()
		} else {
			km.mu.Lock()
			km.lastError = ""
			km.mu.Unlock()
			log.Printf("  Knowledge module: kiwix-serve stopped")
		}
	}()

	return nil
}

// Stop terminates the kiwix-serve process.
func (km *KiwixManager) Stop() {
	km.mu.Lock()
	defer km.mu.Unlock()
	if km.cancelFunc != nil {
		km.cancelFunc()
		km.cancelFunc = nil
	}
}

// IsRunning returns whether kiwix-serve is currently running.
func (km *KiwixManager) IsRunning() bool {
	km.mu.RLock()
	defer km.mu.RUnlock()
	return km.running
}

// GetLastError returns the last error encountered.
func (km *KiwixManager) GetLastError() string {
	km.mu.RLock()
	defer km.mu.RUnlock()
	return km.lastError
}

// Port returns the port kiwix-serve is running on.
func (km *KiwixManager) Port() int {
	return km.port
}

// SearchResult represents a single search result from kiwix-serve.
type SearchResult struct {
	Title   string `json:"title"`
	Path    string `json:"path"`
	Snippet string `json:"snippet"`
	ZimID   string `json:"zim_id"`
}

// ProxySearch forwards a search request to kiwix-serve and returns results.
// When kiwix-serve is not running, it falls back to a built-in stub search.
func (km *KiwixManager) ProxySearch(query string, limit int) ([]SearchResult, error) {
	if !km.IsRunning() {
		// Fallback: search ZIM file names for a match
		return km.fallbackSearch(query, limit), nil
	}

	kiwixURL := fmt.Sprintf("http://localhost:%d/search?pattern=%s&pageLength=%d",
		km.port, url.QueryEscape(query), limit)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(kiwixURL)
	if err != nil {
		return nil, fmt.Errorf("kiwix-serve search failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse kiwix-serve JSON search response
	var results []SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		// kiwix-serve may return HTML; try to handle gracefully
		return km.fallbackSearch(query, limit), nil
	}
	return results, nil
}

// ProxyArticle fetches an article from kiwix-serve.
func (km *KiwixManager) ProxyArticle(zimID, articlePath string) ([]byte, string, error) {
	if !km.IsRunning() {
		return []byte(km.generateOfflinePage(zimID, articlePath)), "text/html", nil
	}

	kiwixURL := fmt.Sprintf("http://localhost:%d/%s/%s", km.port, zimID, articlePath)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(kiwixURL)
	if err != nil {
		return nil, "", fmt.Errorf("kiwix-serve article fetch failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read article body: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "text/html"
	}
	return body, contentType, nil
}

// fallbackSearch provides basic search when kiwix-serve is unavailable.
func (km *KiwixManager) fallbackSearch(query string, limit int) []SearchResult {
	results := make([]SearchResult, 0)
	query = strings.ToLower(query)

	for _, zf := range km.GetZIMFiles() {
		if strings.Contains(strings.ToLower(zf.Name), query) ||
			strings.Contains(strings.ToLower(zf.ID), query) {
			results = append(results, SearchResult{
				Title:   zf.Name,
				Path:    "/" + zf.ID,
				Snippet: fmt.Sprintf("ZIM content pack: %s (%s)", zf.Name, zf.SizeHuman),
				ZimID:   zf.ID,
			})
		}
		if len(results) >= limit {
			break
		}
	}

	return results
}

// generateOfflinePage creates an informational HTML page when kiwix-serve isn't available.
func (km *KiwixManager) generateOfflinePage(zimID, articlePath string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><title>AEGIS Knowledge</title>
<style>
body { font-family: Inter, system-ui, sans-serif; background: #0f172a; color: #e2e8f0; 
       max-width: 600px; margin: 80px auto; padding: 20px; text-align: center; }
h1 { color: #10b981; }
.info { background: #1e293b; border: 1px solid #334155; border-radius: 8px; padding: 20px; margin: 20px 0; }
code { background: #1e293b; padding: 2px 6px; border-radius: 4px; color: #10b981; }
</style></head><body>
<h1>📚 Knowledge library</h1>
<div class="info">
<p><strong>Kiwix sidecar not running</strong></p>
<p>The article <code>%s/%s</code> is available but requires the kiwix-serve sidecar to render.</p>
<p>Place the <code>kiwix-serve</code> binary in:<br><code>sidecars/kiwix-serve/windows/</code></p>
<p>Then restart AEGIS to enable full article rendering.</p>
</div>
<p>%d ZIM files detected in your content packs.</p>
</body></html>`, zimID, articlePath, len(km.zimFiles))
}

// humanizeBytes converts bytes to a human-readable string.
func humanizeBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
