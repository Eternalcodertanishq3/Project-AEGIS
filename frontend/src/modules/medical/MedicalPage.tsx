import { useState, useEffect } from 'react'
import {
  Heart, Loader2, ArrowLeft, AlertTriangle,
  AlertCircle, Info, ChevronRight
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

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="w-8 h-8 text-red-500 animate-spin" />
      </div>
    )
  }

  // Entry detail view
  if (activeEntry) {
    const sev = severityStyles[activeEntry.severity] || severityStyles.info
    const SevIcon = sev.icon
    return (
      <div className="flex flex-col h-full">
        <div className="p-6 pb-4 border-b border-slate-800/50">
          <button
            onClick={() => setActiveEntry(null)}
            className="flex items-center gap-2 text-sm text-slate-400 hover:text-slate-200 mb-3 transition-colors"
          >
            <ArrowLeft className="w-4 h-4" /> Back to {activeCategory?.name}
          </button>
          <div className="flex items-center gap-3">
            <div className={`px-3 py-1 rounded-full ${sev.bg} ${sev.border} border flex items-center gap-1.5`}>
              <SevIcon className={`w-3.5 h-3.5 ${sev.text}`} />
              <span className={`text-xs font-semibold uppercase ${sev.text}`}>{activeEntry.severity}</span>
            </div>
            <h1 className="text-xl font-semibold text-slate-100">{activeEntry.title}</h1>
          </div>
          <p className="text-sm text-slate-400 mt-2">{activeEntry.summary}</p>
        </div>

        <div className="flex-1 overflow-y-auto p-6 space-y-6">
          {/* Steps */}
          <div>
            <h2 className="text-sm font-semibold text-slate-300 uppercase tracking-wider mb-3">
              Step-by-Step Instructions
            </h2>
            <div className="space-y-2">
              {activeEntry.steps.map((step, i) => (
                <div key={i} className="flex gap-3 p-3 rounded-lg bg-slate-800/30 border border-slate-700/30">
                  <div className="w-7 h-7 rounded-full bg-emerald-500/20 border border-emerald-500/30 flex items-center justify-center shrink-0">
                    <span className="text-xs font-bold text-emerald-400">{i + 1}</span>
                  </div>
                  <p className="text-sm text-slate-200 leading-relaxed pt-1">{step}</p>
                </div>
              ))}
            </div>
          </div>

          {/* Warnings */}
          {activeEntry.warnings && activeEntry.warnings.length > 0 && (
            <div>
              <h2 className="text-sm font-semibold text-red-400 uppercase tracking-wider mb-3">
                ⚠ Critical Warnings
              </h2>
              <div className="space-y-2">
                {activeEntry.warnings.map((warning, i) => (
                  <div key={i} className="flex gap-3 p-3 rounded-lg bg-red-500/5 border border-red-500/20">
                    <AlertTriangle className="w-5 h-5 text-red-400 shrink-0 mt-0.5" />
                    <p className="text-sm text-red-300 leading-relaxed">{warning}</p>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    )
  }

  // Category detail view
  if (activeCategory) {
    return (
      <div className="flex flex-col h-full">
        <div className="p-6 pb-4 border-b border-slate-800/50">
          <button
            onClick={() => setActiveCategory(null)}
            className="flex items-center gap-2 text-sm text-slate-400 hover:text-slate-200 mb-3 transition-colors"
          >
            <ArrowLeft className="w-4 h-4" /> All Categories
          </button>
          <h1 className="text-xl font-semibold text-slate-100">{activeCategory.name}</h1>
          <p className="text-sm text-slate-400 mt-1">{activeCategory.description}</p>
        </div>

        <div className="flex-1 overflow-y-auto p-6 space-y-3">
          {activeCategory.entries.map(entry => {
            const sev = severityStyles[entry.severity] || severityStyles.info
            const SevIcon = sev.icon
            return (
              <button
                key={entry.id}
                onClick={() => setActiveEntry(entry)}
                className="w-full text-left p-4 rounded-xl bg-slate-800/30 border border-slate-700/30 hover:border-slate-600/50 transition-all group"
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className={`px-2 py-0.5 rounded-full ${sev.bg} ${sev.border} border flex items-center gap-1`}>
                      <SevIcon className={`w-3 h-3 ${sev.text}`} />
                      <span className={`text-[10px] font-semibold uppercase ${sev.text}`}>{entry.severity}</span>
                    </div>
                    <span className="text-sm font-medium text-slate-200">{entry.title}</span>
                  </div>
                  <ChevronRight className="w-4 h-4 text-slate-600 group-hover:text-slate-400 transition-colors" />
                </div>
                <p className="text-xs text-slate-500 mt-2 ml-0">{entry.summary}</p>
                <p className="text-xs text-slate-600 mt-1">{entry.steps.length} steps</p>
              </button>
            )
          })}
        </div>
      </div>
    )
  }

  // Categories overview
  return (
    <div className="flex flex-col h-full">
      <div className="p-6 pb-4 border-b border-slate-800/50">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-red-500/20 border border-red-500/30 flex items-center justify-center">
            <Heart className="w-5 h-5 text-red-400" />
          </div>
          <div>
            <h1 className="text-xl font-semibold text-slate-100">Medical Triage</h1>
            <p className="text-sm text-slate-400">
              Offline emergency medical reference · {categories.length} categories
            </p>
          </div>
        </div>

        <div className="mt-4 p-3 rounded-lg bg-red-500/5 border border-red-500/20">
          <div className="flex items-start gap-2">
            <AlertTriangle className="w-4 h-4 text-red-400 shrink-0 mt-0.5" />
            <p className="text-xs text-red-300">
              <strong>Disclaimer:</strong> This is a basic first-aid reference for emergencies only. 
              It is NOT a substitute for professional medical care. Always seek qualified medical help when available.
            </p>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-6">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {categories.map(cat => (
            <button
              key={cat.id}
              onClick={() => setActiveCategory(cat)}
              className="text-left p-5 rounded-xl bg-slate-800/30 border border-slate-700/30 hover:border-red-500/30 hover:bg-red-500/5 transition-all group"
            >
              <div className="flex items-center justify-between mb-3">
                <h3 className="text-sm font-semibold text-slate-200 group-hover:text-red-300 transition-colors">
                  {cat.name}
                </h3>
                <ChevronRight className="w-4 h-4 text-slate-600 group-hover:text-red-400 transition-colors" />
              </div>
              <p className="text-xs text-slate-500 mb-3">{cat.description}</p>
              <div className="flex items-center gap-2">
                <span className="text-[10px] px-2 py-0.5 rounded-full bg-slate-700/50 text-slate-400">
                  {cat.entries.length} {cat.entries.length === 1 ? 'entry' : 'entries'}
                </span>
                {cat.entries.some(e => e.severity === 'critical') && (
                  <span className="text-[10px] px-2 py-0.5 rounded-full bg-red-500/10 text-red-400 border border-red-500/20">
                    CRITICAL
                  </span>
                )}
              </div>
            </button>
          ))}
        </div>
      </div>
    </div>
  )
}
