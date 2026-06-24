import { useState, useEffect, useCallback } from 'react'
import {
  Search,
  BookOpen,
  HardDrive,
  AlertCircle,
  ArrowLeft,
  Loader2,
} from 'lucide-react'
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
    <div className="flex flex-col h-full bg-[#030712] relative overflow-hidden">
      {/* Background ambient glow */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[400px] bg-emerald-500/5 blur-[120px] rounded-full pointer-events-none" />

      {/* Header for non-home states */}
      {viewState !== 'home' && (
        <div className="flex-none p-4 md:p-6 border-b border-slate-800/50 bg-slate-900/50 backdrop-blur-xl z-10 sticky top-0 flex items-center gap-4">
          <button
            onClick={() => setViewState(viewState === 'article' ? 'search' : 'home')}
            className="p-2 rounded-xl hover:bg-slate-800/80 transition-all text-slate-400 hover:text-emerald-400 group"
          >
            <ArrowLeft className="h-5 w-5 transition-transform group-hover:-translate-x-1" />
          </button>
          
          <form onSubmit={handleSearch} className="flex-1 max-w-3xl relative">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-4 w-4 text-emerald-500/50" />
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Search offline encyclopedias..."
              className="w-full pl-11 pr-4 py-3 rounded-xl bg-slate-900/80 border border-slate-700/50 text-slate-100 placeholder:text-slate-500 focus:outline-none focus:border-emerald-500/50 focus:ring-1 focus:ring-emerald-500/20 transition-all shadow-inner shadow-black/20"
            />
            {searching && (
              <Loader2 className="absolute right-4 top-1/2 -translate-y-1/2 h-4 w-4 animate-spin text-emerald-500" />
            )}
          </form>
        </div>
      )}

      {/* Main Content Area */}
      <div className="flex-1 overflow-y-auto z-10">
        
        {/* Google-like Home State */}
        {viewState === 'home' && (
          <div className="flex flex-col items-center justify-center min-h-[80vh] px-4">
            <div className="w-20 h-20 bg-gradient-to-br from-emerald-400/20 to-cyan-400/20 rounded-3xl flex items-center justify-center mb-8 border border-emerald-500/20 shadow-[0_0_40px_rgba(16,185,129,0.1)]">
              <BookOpen className="h-10 w-10 text-emerald-400" />
            </div>
            
            <h1 className="text-4xl md:text-5xl font-bold tracking-tight text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-cyan-400 mb-8 text-center">
              Knowledge Library
            </h1>

            <form onSubmit={handleSearch} className="w-full max-w-2xl relative group">
              <Search className="absolute left-5 top-1/2 -translate-y-1/2 h-5 w-5 text-slate-400 transition-colors group-focus-within:text-emerald-400" />
              <input
                type="text"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                placeholder="Search offline encyclopedias, wikis, and manuals..."
                className="w-full pl-14 pr-6 py-4 rounded-2xl bg-slate-900/60 backdrop-blur-md border border-slate-700/50 text-lg text-slate-100 placeholder:text-slate-500/80 focus:outline-none focus:border-emerald-500/50 focus:ring-4 focus:ring-emerald-500/10 transition-all shadow-2xl shadow-black/40"
              />
              {searching && (
                <Loader2 className="absolute right-5 top-1/2 -translate-y-1/2 h-5 w-5 animate-spin text-emerald-500" />
              )}
            </form>

            <div className="mt-8 flex gap-4 text-sm text-slate-400">
              <div className="px-4 py-1.5 rounded-full bg-slate-800/50 border border-slate-700/50 backdrop-blur-sm flex items-center gap-2">
                <HardDrive className="h-4 w-4 text-emerald-400/70" />
                {status?.zim_files_count || 0} Content Packs
              </div>
              <div className="px-4 py-1.5 rounded-full bg-slate-800/50 border border-slate-700/50 backdrop-blur-sm flex items-center gap-2">
                <div className={`w-2 h-2 rounded-full ${status?.kiwix_running ? 'bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.8)]' : 'bg-red-500'}`} />
                {status?.kiwix_running ? 'Engine Online' : 'Engine Offline'}
              </div>
            </div>

            {status?.kiwix_note && (
              <div className="mt-6 p-4 rounded-xl bg-red-500/10 border border-red-500/20 flex items-start gap-3 max-w-2xl text-left">
                <AlertCircle className="h-5 w-5 text-red-400 shrink-0 mt-0.5" />
                <p className="text-sm text-red-200/90 leading-relaxed">{status.kiwix_note}</p>
              </div>
            )}
          </div>
        )}

        {/* Search Results */}
        {viewState === 'search' && (
          <div className="p-4 md:p-8 max-w-4xl mx-auto space-y-4">
            <p className="text-sm font-medium text-slate-400 mb-6">
              About {results.length} results
            </p>
            
            {results.length === 0 && !searching && (
              <div className="py-20 text-center">
                <Search className="h-12 w-12 text-slate-700 mx-auto mb-4" />
                <h3 className="text-xl font-semibold text-slate-300 mb-2">No results found for "{query}"</h3>
                <p className="text-slate-500">Try adjusting your keywords or adding more ZIM content packs.</p>
              </div>
            )}

            <div className="space-y-6">
              {results.map((result, idx) => (
                <div key={idx} className="group">
                  <p className="text-xs font-mono text-emerald-400/70 mb-1 flex items-center gap-2">
                    <BookOpen className="h-3 w-3" />
                    {result.zim_id}
                  </p>
                  <button
                    onClick={() => handleViewArticle(result.zim_id, result.path)}
                    className="block text-left"
                  >
                    <h3 className="text-xl font-semibold text-slate-200 group-hover:text-emerald-400 group-hover:underline decoration-emerald-400/30 underline-offset-4 transition-colors mb-2 leading-tight">
                      {result.title}
                    </h3>
                  </button>
                  {result.snippet && (
                    <p className="text-sm text-slate-400 leading-relaxed line-clamp-3 max-w-3xl">
                      {result.snippet}
                    </p>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Article Reader */}
        {viewState === 'article' && (
          <div className="p-4 md:p-8 max-w-4xl mx-auto">
            <div className="glass-panel p-6 md:p-10 shadow-2xl">
              <h1 className="text-3xl font-bold text-slate-100 mb-8 pb-6 border-b border-slate-700/50">{articleTitle}</h1>
              <div
                className="prose prose-invert prose-lg max-w-none 
                           prose-headings:text-slate-200 prose-headings:font-semibold
                           prose-p:text-slate-300 prose-p:leading-relaxed
                           prose-a:text-emerald-400 prose-a:no-underline hover:prose-a:underline
                           prose-code:text-emerald-300 prose-code:bg-slate-800/80 prose-code:px-1.5 prose-code:py-0.5 prose-code:rounded-md
                           prose-img:rounded-xl prose-img:shadow-xl
                           prose-hr:border-slate-700/50
                           prose-li:text-slate-300"
                dangerouslySetInnerHTML={{ __html: articleHtml }}
              />
            </div>
          </div>
        )}

      </div>
    </div>
  )
}
