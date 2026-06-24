import { useState, useEffect, useMemo } from 'react'
import {
  Heart, Loader2, ArrowLeft, AlertTriangle,
  AlertCircle, Info, ChevronRight, Search
} from 'lucide-react'

interface Entry {
  id: string
  title: string
  severity: string
  summary: string
  steps: string[]
  warnings?: string[]
}

interface Category {
  id: string
  name: string
  icon: string
  description: string
  entries: Entry[]
}

const severityStyles: Record<string, { bg: string; border: string; text: string; icon: React.ElementType }> = {
  critical: { bg: 'bg-red-500/10', border: 'border-red-500/30', text: 'text-red-400', icon: AlertTriangle },
  warning: { bg: 'bg-amber-500/10', border: 'border-amber-500/30', text: 'text-amber-400', icon: AlertCircle },
  info: { bg: 'bg-sky-500/10', border: 'border-sky-500/30', text: 'text-sky-400', icon: Info },
}

export function MedicalPage() {
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [activeCategory, setActiveCategory] = useState<Category | null>(null)
  const [activeEntry, setActiveEntry] = useState<Entry | null>(null)
  const [searchQuery, setSearchQuery] = useState('')

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const res = await fetch('/api/medical/categories')
        const data = await res.json()
        setCategories(data.categories || [])
      } catch (e) {
        console.error('Failed to fetch medical categories:', e)
      } finally {
        setLoading(false)
      }
    }
    fetchCategories()
  }, [])

  // Flatten all entries for search
  const searchResults = useMemo(() => {
    if (!searchQuery.trim()) return []
    const query = searchQuery.toLowerCase()
    
    const results: { category: Category; entry: Entry }[] = []
    categories.forEach(cat => {
      cat.entries.forEach(entry => {
        if (
          entry.title.toLowerCase().includes(query) ||
          entry.summary.toLowerCase().includes(query) ||
          cat.name.toLowerCase().includes(query)
        ) {
          results.push({ category: cat, entry })
        }
      })
    })
    return results
  }, [searchQuery, categories])

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full bg-[#030712]">
        <Loader2 className="w-8 h-8 text-red-500 animate-spin" />
      </div>
    )
  }

  // Entry detail view
  if (activeEntry) {
    const sev = severityStyles[activeEntry.severity] || severityStyles.info
    const SevIcon = sev.icon
    return (
      <div className="flex flex-col h-full bg-[#030712] relative overflow-hidden">
        <div className="absolute top-0 right-0 w-[600px] h-[300px] bg-red-500/5 blur-[100px] rounded-full pointer-events-none" />
        
        <div className="p-6 md:p-8 pb-4 border-b border-slate-800/50 bg-slate-900/50 backdrop-blur-xl z-10 sticky top-0">
          <button
            onClick={() => { setActiveEntry(null); setSearchQuery(''); }}
            className="flex items-center gap-2 text-sm text-slate-400 hover:text-slate-200 mb-6 transition-colors group"
          >
            <ArrowLeft className="w-4 h-4 transition-transform group-hover:-translate-x-1" /> Back
          </button>
          <div className="flex items-center gap-3 mb-2">
            <div className={`px-3 py-1 rounded-full ${sev.bg} ${sev.border} border flex items-center gap-1.5 shadow-lg`}>
              <SevIcon className={`w-3.5 h-3.5 ${sev.text}`} />
              <span className={`text-xs font-bold tracking-widest uppercase ${sev.text}`}>{activeEntry.severity}</span>
            </div>
            <h1 className="text-2xl md:text-3xl font-bold text-slate-100 tracking-tight">{activeEntry.title}</h1>
          </div>
          <p className="text-base text-slate-400 max-w-3xl leading-relaxed">{activeEntry.summary}</p>
        </div>

        <div className="flex-1 overflow-y-auto p-6 md:p-8 max-w-4xl mx-auto w-full z-10 space-y-8">
          {/* Warnings */}
          {activeEntry.warnings && activeEntry.warnings.length > 0 && (
            <div>
              <h2 className="text-xs font-bold text-red-400/80 uppercase tracking-widest mb-3">
                ⚠ Critical Warnings
              </h2>
              <div className="space-y-3">
                {activeEntry.warnings.map((warning, i) => (
                  <div key={i} className="flex gap-4 p-4 rounded-xl bg-red-500/10 border border-red-500/20 shadow-lg shadow-red-500/5">
                    <AlertTriangle className="w-6 h-6 text-red-400 shrink-0" />
                    <p className="text-sm font-medium text-red-200/90 leading-relaxed">{warning}</p>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Steps */}
          <div>
            <h2 className="text-xs font-bold text-slate-500 uppercase tracking-widest mb-4">
              Step-by-Step Instructions
            </h2>
            <div className="space-y-3">
              {activeEntry.steps.map((step, i) => (
                <div key={i} className="flex gap-4 p-4 rounded-xl bg-slate-800/40 border border-slate-700/50 backdrop-blur-sm">
                  <div className="w-8 h-8 rounded-full bg-slate-700/50 border border-slate-600 flex items-center justify-center shrink-0 shadow-inner">
                    <span className="text-sm font-bold text-slate-300">{i + 1}</span>
                  </div>
                  <p className="text-base text-slate-200 leading-relaxed pt-0.5">{step}</p>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    )
  }

  // Category detail view
  if (activeCategory) {
    return (
      <div className="flex flex-col h-full bg-[#030712] relative overflow-hidden">
        <div className="absolute top-0 right-0 w-[500px] h-[300px] bg-red-500/5 blur-[100px] rounded-full pointer-events-none" />
        <div className="p-6 md:p-8 pb-6 border-b border-slate-800/50 bg-slate-900/50 backdrop-blur-xl z-10 sticky top-0">
          <button
            onClick={() => setActiveCategory(null)}
            className="flex items-center gap-2 text-sm text-slate-400 hover:text-slate-200 mb-6 transition-colors group"
          >
            <ArrowLeft className="w-4 h-4 transition-transform group-hover:-translate-x-1" /> All Categories
          </button>
          <h1 className="text-2xl md:text-3xl font-bold text-slate-100 tracking-tight">{activeCategory.name}</h1>
          <p className="text-sm text-slate-400 mt-2 max-w-2xl">{activeCategory.description}</p>
        </div>

        <div className="flex-1 overflow-y-auto p-6 md:p-8 max-w-4xl mx-auto w-full z-10 space-y-4">
          {activeCategory.entries.map(entry => {
            const sev = severityStyles[entry.severity] || severityStyles.info
            const SevIcon = sev.icon
            return (
              <button
                key={entry.id}
                onClick={() => setActiveEntry(entry)}
                className="w-full text-left p-5 rounded-xl bg-slate-900/60 backdrop-blur-md border border-slate-700/50 hover:border-slate-500/50 hover:bg-slate-800/80 transition-all group shadow-lg"
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <div className={`px-2.5 py-1 rounded-full ${sev.bg} ${sev.border} border flex items-center gap-1.5`}>
                      <SevIcon className={`w-3.5 h-3.5 ${sev.text}`} />
                      <span className={`text-[10px] font-bold tracking-wider uppercase ${sev.text}`}>{entry.severity}</span>
                    </div>
                    <span className="text-lg font-semibold text-slate-200 group-hover:text-slate-100 transition-colors">{entry.title}</span>
                  </div>
                  <ChevronRight className="w-5 h-5 text-slate-600 group-hover:text-slate-400 transition-colors" />
                </div>
                <p className="text-sm text-slate-400 mt-3">{entry.summary}</p>
                <p className="text-xs font-semibold text-slate-500 mt-3 flex items-center gap-1.5">
                  <span className="w-1.5 h-1.5 rounded-full bg-slate-600" />
                  {entry.steps.length} steps
                </p>
              </button>
            )
          })}
        </div>
      </div>
    )
  }

  // Search or Categories overview
  return (
    <div className="flex flex-col h-full bg-[#030712] relative overflow-hidden">
      <div className="absolute top-0 right-1/4 w-[600px] h-[300px] bg-red-500/5 blur-[120px] rounded-full pointer-events-none" />
      
      <div className="p-6 md:p-8 border-b border-slate-800/50 bg-slate-900/50 backdrop-blur-xl z-10 sticky top-0">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-6">
          <div className="flex items-center gap-4">
            <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-red-500/20 to-rose-500/10 border border-red-500/30 flex items-center justify-center shadow-[0_0_20px_rgba(239,68,68,0.15)]">
              <Heart className="w-6 h-6 text-red-400" />
            </div>
            <div>
              <h1 className="text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-red-400 to-rose-300">Medical Triage</h1>
              <p className="text-sm text-slate-400">Offline emergency medical reference</p>
            </div>
          </div>
          
          <div className="relative flex-1 max-w-xl group">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-slate-500 transition-colors group-focus-within:text-red-400" />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search all emergency procedures..."
              className="w-full pl-12 pr-4 py-3.5 rounded-xl bg-slate-900/80 border border-slate-700/50 text-slate-100 placeholder:text-slate-500/80 focus:outline-none focus:border-red-500/50 focus:ring-2 focus:ring-red-500/20 transition-all shadow-inner shadow-black/20"
            />
          </div>
        </div>

        {!searchQuery && (
          <div className="mt-6 p-4 rounded-xl bg-red-500/5 border border-red-500/20 flex items-start gap-3">
            <AlertTriangle className="w-5 h-5 text-red-400 shrink-0 mt-0.5" />
            <p className="text-sm text-red-200/90 leading-relaxed">
              <strong>Disclaimer:</strong> This is a basic first-aid reference for emergencies only. 
              It is NOT a substitute for professional medical care. Always seek qualified medical help when available.
            </p>
          </div>
        )}
      </div>

      <div className="flex-1 overflow-y-auto p-6 md:p-8 z-10">
        {searchQuery ? (
          // Search Results
          <div className="max-w-4xl mx-auto space-y-4">
            <h2 className="text-sm font-bold text-slate-500 uppercase tracking-widest mb-4">
              Search Results ({searchResults.length})
            </h2>
            
            {searchResults.length === 0 ? (
              <div className="py-20 text-center">
                <Search className="h-12 w-12 text-slate-700 mx-auto mb-4" />
                <h3 className="text-xl font-semibold text-slate-300 mb-2">No procedures found for "{searchQuery}"</h3>
                <p className="text-slate-500">Try using simpler terms or checking the categories directly.</p>
              </div>
            ) : (
              searchResults.map(({ category, entry }) => {
                const sev = severityStyles[entry.severity] || severityStyles.info
                const SevIcon = sev.icon
                return (
                  <button
                    key={`${category.id}-${entry.id}`}
                    onClick={() => setActiveEntry(entry)}
                    className="w-full text-left p-5 rounded-xl bg-slate-900/60 backdrop-blur-md border border-slate-700/50 hover:border-slate-500/50 hover:bg-slate-800/80 transition-all group shadow-lg"
                  >
                    <div className="flex items-center justify-between mb-2">
                      <div className="flex items-center gap-4">
                        <div className={`px-2.5 py-1 rounded-full ${sev.bg} ${sev.border} border flex items-center gap-1.5`}>
                          <SevIcon className={`w-3.5 h-3.5 ${sev.text}`} />
                          <span className={`text-[10px] font-bold tracking-wider uppercase ${sev.text}`}>{entry.severity}</span>
                        </div>
                        <span className="text-xs font-bold text-slate-500 uppercase tracking-wider">{category.name}</span>
                      </div>
                      <ChevronRight className="w-5 h-5 text-slate-600 group-hover:text-slate-400 transition-colors" />
                    </div>
                    <h3 className="text-lg font-semibold text-slate-200 group-hover:text-slate-100 transition-colors mb-2">{entry.title}</h3>
                    <p className="text-sm text-slate-400 line-clamp-2">{entry.summary}</p>
                  </button>
                )
              })
            )}
          </div>
        ) : (
          // Categories Grid
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 max-w-7xl mx-auto">
            {categories.map(cat => (
              <button
                key={cat.id}
                onClick={() => setActiveCategory(cat)}
                className="text-left p-6 rounded-2xl glass-panel group hover:-translate-y-1 transition-all duration-300"
              >
                <div className="flex items-start justify-between mb-4">
                  <h3 className="text-lg font-bold text-slate-200 group-hover:text-red-300 transition-colors">
                    {cat.name}
                  </h3>
                  <div className="w-8 h-8 rounded-full bg-slate-800/80 flex items-center justify-center group-hover:bg-red-500/20 transition-colors">
                    <ChevronRight className="w-4 h-4 text-slate-500 group-hover:text-red-400 transition-colors" />
                  </div>
                </div>
                <p className="text-sm text-slate-400 mb-6 line-clamp-2 leading-relaxed">{cat.description}</p>
                <div className="flex items-center gap-3">
                  <span className="text-xs font-semibold px-3 py-1 rounded-full bg-slate-800/80 text-slate-400 border border-slate-700/50">
                    {cat.entries.length} {cat.entries.length === 1 ? 'procedure' : 'procedures'}
                  </span>
                  {cat.entries.some(e => e.severity === 'critical') && (
                    <span className="text-xs font-bold px-3 py-1 rounded-full bg-red-500/10 text-red-400 border border-red-500/20 shadow-[0_0_10px_rgba(239,68,68,0.1)]">
                      CRITICAL
                    </span>
                  )}
                </div>
              </button>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
