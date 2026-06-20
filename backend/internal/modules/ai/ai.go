package ai

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// ModelFile represents a discovered GGUF model file.
type ModelFile struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	SizeBytes int64  `json:"size_bytes"`
	SizeHuman string `json:"size_human"`
}

// AIManager manages the llama-server sidecar process lifecycle.
type AIManager struct {
	mu           sync.RWMutex
	dataDir      string
	sidecarPath  string
	port         int
	running      bool
	activeModel  string
	models       []ModelFile
	cancelFunc   context.CancelFunc
	cmd          *exec.Cmd
}

// NewAIManager creates a new AIManager.
func NewAIManager(dataDir string, llamaPort int) *AIManager {
	return &AIManager{
		dataDir: dataDir,
		port:    llamaPort,
		models:  make([]ModelFile, 0),
	}
}

// DiscoverModels scans content-packs and data directory for .gguf model files.
func (am *AIManager) DiscoverModels() ([]ModelFile, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	am.models = make([]ModelFile, 0)

	searchPaths := []string{
		am.dataDir,
		filepath.Join(filepath.Dir(am.dataDir), "content-packs"),
		filepath.Join(filepath.Dir(am.dataDir), "content-packs", "models-ai"),
		filepath.Join(am.dataDir, "models"),
	}

	seen := make(map[string]bool)

	for _, searchPath := range searchPaths {
		err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".gguf" {
				absPath, _ := filepath.Abs(path)
				if seen[absPath] {
					return nil
				}
				seen[absPath] = true

				id := strings.TrimSuffix(info.Name(), ext)
				am.models = append(am.models, ModelFile{
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
			continue
		}
	}

	log.Printf("  AI: discovered %d GGUF models", len(am.models))
	return am.models, nil
}

// GetModels returns currently discovered models.
func (am *AIManager) GetModels() []ModelFile {
	am.mu.RLock()
	defer am.mu.RUnlock()
	result := make([]ModelFile, len(am.models))
	copy(result, am.models)
	return result
}

// FindSidecar looks for the llama-server binary.
func (am *AIManager) FindSidecar() (string, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	candidates := []string{
		filepath.Join(filepath.Dir(am.dataDir), "sidecars", "llama", "windows", "llama-server.exe"),
		filepath.Join(filepath.Dir(am.dataDir), "sidecars", "llama", "llama-server.exe"),
		filepath.Join(filepath.Dir(am.dataDir), "sidecars", "llama", "windows", "server.exe"),
		filepath.Join(filepath.Dir(am.dataDir), "llama-server.exe"),
		"llama-server", // try PATH
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			am.sidecarPath = path
			return path, nil
		}
	}

	return "", fmt.Errorf("llama-server not found; place it in sidecars/llama/windows/")
}

// Start launches the llama-server with the specified model.
func (am *AIManager) Start(modelID string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if am.running {
		return fmt.Errorf("AI server is already running with model %q; stop it first", am.activeModel)
	}

	if am.sidecarPath == "" {
		return fmt.Errorf("llama-server not found")
	}

	// Find the model
	var modelPath string
	for _, m := range am.models {
		if m.ID == modelID {
			modelPath = m.Path
			break
		}
	}
	if modelPath == "" {
		return fmt.Errorf("model %q not found", modelID)
	}

	ctx, cancel := context.WithCancel(context.Background())
	am.cancelFunc = cancel

	// llama-server args: --model <path> --port <port> --host 127.0.0.1 --ctx-size 2048
	args := []string{
		"--model", modelPath,
		"--port", fmt.Sprintf("%d", am.port),
		"--host", "127.0.0.1",
		"--ctx-size", "2048",
	}

	am.cmd = exec.CommandContext(ctx, am.sidecarPath, args...)

	if err := am.cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("failed to start llama-server: %w", err)
	}

	am.running = true
	am.activeModel = modelID
	log.Printf("✓ AI module: started llama-server on port %d with model %q", am.port, modelID)

	// Monitor process in background
	go func() {
		err := am.cmd.Wait()
		am.mu.Lock()
		am.running = false
		am.activeModel = ""
		am.cmd = nil
		am.mu.Unlock()
		if err != nil && !strings.Contains(err.Error(), "killed") && !strings.Contains(err.Error(), "signal") {
			log.Printf("! AI module: llama-server exited with error: %v", err)
		} else {
			log.Printf("  AI module: llama-server stopped")
		}
	}()

	return nil
}

// Stop terminates the llama-server process.
func (am *AIManager) Stop() error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if !am.running {
		return fmt.Errorf("AI server is not running")
	}

	if am.cancelFunc != nil {
		am.cancelFunc()
		am.cancelFunc = nil
	}
	return nil
}

// IsRunning returns whether llama-server is currently running.
func (am *AIManager) IsRunning() bool {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.running
}

// ActiveModel returns the ID of the currently loaded model.
func (am *AIManager) ActiveModel() string {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.activeModel
}

// Port returns the port llama-server is running on.
func (am *AIManager) Port() int {
	return am.port
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
