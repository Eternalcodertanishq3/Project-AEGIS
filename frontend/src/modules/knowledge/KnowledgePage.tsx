import { useState, useEffect, useCallback } from 'react'
import {
  Search,
  BookOpen,
  FileText,
  HardDrive,
  AlertCircle,
  ArrowLeft,
  ExternalLink,
  Loader2,
} from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { apiFetch } from '@/hooks/useApi'

interface ZIMFile {
  id: string
  name: string
  path: string
  size_bytes: number
  size_human: string
}

interface SearchResult {
  title: string
  path: string
  snippet: string
  zim_id: string
}

interface KnowledgeStatus {
  kiwix_running: boolean
  kiwix_port: number
  zim_files_count: number
  zim_files: ZIMFile[]
  kiwix_sidecar: string
  kiwix_note?: string
}

type ViewState = 'home' | 'search' | 'article'

export function KnowledgePage() {
  const [status, setStatus] = useState<KnowledgeStatus | null>(null)
  const [query, setQuery] = useState('')
  const [results, setResults] = useState<SearchResult[]>([])
  const [searching, setSearching] = useState(false)
  const [viewState, setViewState] = useState<ViewState>('home')
  const [articleHtml, setArticleHtml] = useState('')
  const [articleTitle, setArticleTitle] = useState('')
  const [loading, setLoading] = useState(true)

  // Fetch knowledge module status
  useEffect(() => {
    async function fetchStatus() {
      try {
        const data = await apiFetch<KnowledgeStatus>('/knowledge/status')
        setStatus(data)
      } catch {
        // Backend may not have knowledge module yet
      } finally {
        setLoading(false)
      }
    }
    fetchStatus()
  }, [])

  // Search handler
  const handleSearch = useCallback(async (e: React.FormEvent) => {
    e.preventDefault()
    if (!query.trim()) return

    setSearching(true)
    setViewState('search')
    try {
      const data = await apiFetch<{ results: SearchResult[]; count: number }>(
        `/knowledge/search?q=${encodeURIComponent(query)}&limit=25`
      )
      setResults(data.results || [])
    } catch {
      setResults([])
    } finally {
      setSearching(false)
    }
  }, [query])

  // View article handler
  const handleViewArticle = useCallback(async (zimId: string, path: string) => {
    setViewState('article')
    setArticleTitle(path.split('/').pop() || 'Article')
    try {
      const response = await fetch(`/api/knowledge/article/${zimId}${path}`)
      const html = await response.text()
      setArticleHtml(html)
    } catch {
      setArticleHtml('<p>Failed to load article.</p>')
    }
  }, [])

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-6 w-6 animate-spin text-slate-500" />
      </div>
    )
  }

  return (
    <div className="p-6 max-w-6xl space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3">
        {viewState !== 'home' && (
          <button
            onClick={() => setViewState('home')}
            className="p-2 rounded-lg hover:bg-slate-800 transition-colors text-slate-400 hover:text-slate-200"
          >
            <ArrowLeft className="h-4 w-4" />
          </button>
        )}
        <div className="flex items-center gap-2">
          <div className="flex items-center justify-center w-8 h-8 rounded-md bg-blue-600/20 border border-blue-500/30">
            <BookOpen className="h-4 w-4 text-blue-400" />
          </div>
          <div>
            <h1 className="text-xl font-semibold text-slate-100">Knowledge library</h1>
            <p className="text-xs text-slate-500">
              Search offline encyclopedias and references
            </p>
          </div>
        </div>
      </div>

      {/* Search bar */}
      <form onSubmit={handleSearch} className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-500" />
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search articles, topics, references…"
          className="w-full pl-10 pr-4 py-3 rounded-lg bg-slate-800/60 border border-slate-700/50 text-slate-100 
                     placeholder:text-slate-500 focus:outline-none focus:border-emerald-500/50 focus:ring-1 
                     focus:ring-emerald-500/20 transition-all text-sm"
        />
        {searching && (
          <Loader2 className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 animate-spin text-emerald-500" />
        )}
      </form>

      {/* View states */}
      {viewState === 'home' && (
        <div className="space-y-4">
          {/* Status card */}
          <div className="glass-panel p-5">
            <div className="flex items-center gap-2 mb-4">
              <HardDrive className="h-4 w-4 text-slate-400" />
              <h2 className="text-sm font-semibold text-slate-200">Content status</h2>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              {/* Kiwix status */}
              <div className="rounded-lg border border-slate-700/30 bg-slate-800/20 p-4">
                <div className="flex items-center gap-2 mb-2">
                  <span className={`status-dot ${status?.kiwix_running ? 'status-dot-active' : 'status-dot-inactive'}`} />
                  <span className="text-xs font-medium text-slate-500 uppercase">Kiwix engine</span>
                </div>
                <p className="text-sm font-semibold text-slate-100">
                  {status?.kiwix_running ? 'Running' : 'Not running'}
                </p>
                {status?.kiwix_note && (
                  <p className="text-[11px] text-amber-500/80 mt-1">{status.kiwix_note}</p>
                )}
              </div>

              {/* ZIM files */}
              <div className="rounded-lg border border-slate-700/30 bg-slate-800/20 p-4">
                <div className="flex items-center gap-2 mb-2">
                  <FileText className="h-3.5 w-3.5 text-slate-500" />
                  <span className="text-xs font-medium text-slate-500 uppercase">ZIM files</span>
                </div>
                <p className="text-2xl font-semibold text-slate-100">{status?.zim_files_count || 0}</p>
                <p className="text-[11px] text-slate-500">content packs found</p>
              </div>

              {/* Articles available */}
              <div className="rounded-lg border border-slate-700/30 bg-slate-800/20 p-4">
                <div className="flex items-center gap-2 mb-2">
                  <BookOpen className="h-3.5 w-3.5 text-slate-500" />
                  <span className="text-xs font-medium text-slate-500 uppercase">Status</span>
                </div>
                <p className="text-sm font-semibold text-slate-100">
                  {status?.zim_files_count ? 'Ready to search' : 'No content'}
                </p>
                <p className="text-[11px] text-slate-500">
                  {status?.zim_files_count ? 'Enter a query above' : 'Add ZIM files to content-packs/'}
                </p>
              </div>
            </div>
          </div>

          {/* ZIM file list */}
          {status?.zim_files && status.zim_files.length > 0 && (
            <div className="glass-panel p-5">
              <h2 className="text-sm font-semibold text-slate-200 mb-3">Available content packs</h2>
              <div className="space-y-2">
                {status.zim_files.map((zf) => (
                  <div
                    key={zf.id}
                    className="flex items-center justify-between p-3 rounded-lg border border-slate-700/30 bg-slate-800/20"
                  >
                    <div className="flex items-center gap-3">
                      <FileText className="h-4 w-4 text-blue-400" />
                      <div>
                        <p className="text-sm font-medium text-slate-200">{zf.name}</p>
                        <p className="text-[11px] text-slate-500 font-mono">{zf.id}</p>
                      </div>
                    </div>
                    <Badge variant="outline" className="text-xs text-slate-400">{zf.size_human}</Badge>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Empty state */}
          {(!status?.zim_files || status.zim_files.length === 0) && (
            <div className="glass-panel p-8 text-center">
              <AlertCircle className="h-8 w-8 text-slate-600 mx-auto mb-3" />
              <h3 className="text-sm font-semibold text-slate-300 mb-1">No content packs found</h3>
              <p className="text-xs text-slate-500 max-w-md mx-auto">
                Download ZIM files (such as Wikipedia, Stack Overflow, or survival manuals) and place them in the{' '}
                <code className="text-emerald-400 bg-slate-800 px-1 py-0.5 rounded text-[11px]">content-packs/</code>{' '}
                directory, then restart AEGIS.
              </p>
            </div>
          )}
        </div>
      )}

      {viewState === 'search' && (
        <div className="space-y-2">
          <p className="text-xs text-slate-500 mb-3">
            {results.length} result{results.length !== 1 ? 's' : ''} for "{query}"
          </p>
          {results.length === 0 && !searching && (
            <div className="glass-panel p-8 text-center">
              <Search className="h-6 w-6 text-slate-600 mx-auto mb-2" />
              <p className="text-sm text-slate-400">No results found</p>
              <p className="text-xs text-slate-500 mt-1">Try a different query or check your content packs</p>
            </div>
          )}
          {results.map((result, idx) => (
            <button
              key={idx}
              onClick={() => handleViewArticle(result.zim_id, result.path)}
              className="w-full text-left glass-panel p-4 hover:border-slate-600/60 hover:bg-slate-800/60 transition-all group"
            >
              <div className="flex items-start justify-between gap-3">
                <div className="flex-1 min-w-0">
                  <h3 className="text-sm font-semibold text-slate-100 group-hover:text-emerald-400 transition-colors truncate">
                    {result.title}
                  </h3>
                  {result.snippet && (
                    <p className="text-xs text-slate-400 mt-1 line-clamp-2">{result.snippet}</p>
                  )}
                  <p className="text-[11px] text-slate-600 mt-1 font-mono">{result.zim_id}</p>
                </div>
                <ExternalLink className="h-3.5 w-3.5 text-slate-600 group-hover:text-emerald-400 flex-shrink-0 mt-1" />
              </div>
            </button>
          ))}
        </div>
      )}

      {viewState === 'article' && (
        <div className="glass-panel p-6">
          <h2 className="text-sm font-semibold text-slate-200 mb-4">{articleTitle}</h2>
          <div
            className="prose prose-invert prose-sm max-w-none 
                       prose-headings:text-slate-200 prose-p:text-slate-300 
                       prose-a:text-emerald-400 prose-a:no-underline hover:prose-a:underline
                       prose-code:text-emerald-400 prose-code:bg-slate-800/60 prose-code:px-1 prose-code:rounded
                       prose-img:rounded-lg prose-hr:border-slate-700"
            dangerouslySetInnerHTML={{ __html: articleHtml }}
          />
        </div>
      )}
    </div>
  )
}
