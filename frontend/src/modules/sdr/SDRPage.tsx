import { useState, useEffect } from 'react'
import { Activity, Loader2, Search, ChevronRight, ArrowLeft, AlertTriangle, Radio } from 'lucide-react'

interface FrequencyEntry { id: string; frequency: string; freq_mhz: number; name: string; mode: string; description: string; category: string; priority?: string }
interface FrequencyGroup { id: string; name: string; icon: string; description: string; frequencies: FrequencyEntry[] }
interface BandPlan { id: string; name: string; start_mhz: number; end_mhz: number; allocation: string; description: string }

const priorityStyles: Record<string, { bg: string; text: string }> = {
  critical: { bg: 'bg-red-500/10 border-red-500/30', text: 'text-red-400' },
  important: { bg: 'bg-amber-500/10 border-amber-500/30', text: 'text-amber-400' },
  useful: { bg: 'bg-blue-500/10 border-blue-500/30', text: 'text-blue-400' },
}

export function SDRPage() {
  const [groups, setGroups] = useState<FrequencyGroup[]>([])
  const [bandPlans, setBandPlans] = useState<BandPlan[]>([])
  const [activeGroup, setActiveGroup] = useState<FrequencyGroup | null>(null)
  const [search, setSearch] = useState('')
  const [searchResults, setSearchResults] = useState<FrequencyEntry[] | null>(null)
  const [activeView, setActiveView] = useState<'frequencies' | 'bandplans'>('frequencies')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([
      fetch('/api/sdr/frequencies').then(r => r.json()),
      fetch('/api/sdr/bandplans').then(r => r.json()),
    ]).then(([f, b]) => {
      setGroups(f.groups || []); setBandPlans(b.band_plans || [])
    }).finally(() => setLoading(false))
  }, [])

  const doSearch = async () => {
    if (!search.trim()) return
    const res = await fetch(`/api/sdr/search?q=${encodeURIComponent(search)}`)
    const data = await res.json()
    setSearchResults(data.results || [])
  }

  if (loading) return <div className="flex items-center justify-center h-full"><Loader2 className="w-8 h-8 text-cyan-500 animate-spin" /></div>

  if (activeGroup) {
    return (
      <div className="flex flex-col h-full">
        <div className="p-6 pb-4 border-b border-slate-800/50">
          <button onClick={() => setActiveGroup(null)} className="flex items-center gap-2 text-sm text-slate-400 hover:text-slate-200 mb-3 transition-colors">
            <ArrowLeft className="w-4 h-4" /> All Groups
          </button>
          <h1 className="text-xl font-semibold text-slate-100">{activeGroup.name}</h1>
          <p className="text-sm text-slate-400 mt-1">{activeGroup.description}</p>
        </div>
        <div className="flex-1 overflow-y-auto p-6 space-y-2">
          {activeGroup.frequencies.map(freq => {
            const pri = freq.priority ? priorityStyles[freq.priority] : null
            return (
              <div key={freq.id} className="p-4 rounded-xl bg-slate-800/30 border border-slate-700/30">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <span className="text-sm font-mono font-bold text-cyan-400">{freq.frequency}</span>
                    <span className="px-2 py-0.5 rounded text-[10px] bg-slate-700/50 text-slate-400 font-mono">{freq.mode}</span>
                    {pri && <span className={`px-2 py-0.5 rounded text-[10px] border font-bold uppercase ${pri.bg} ${pri.text}`}>{freq.priority}</span>}
                  </div>
                </div>
                <p className="text-sm font-semibold text-slate-200 mt-1">{freq.name}</p>
                <p className="text-xs text-slate-500 mt-0.5">{freq.description}</p>
              </div>
            )
          })}
        </div>
      </div>
    )
  }

  const totalFreqs = groups.reduce((s, g) => s + g.frequencies.length, 0)

  return (
    <div className="flex flex-col h-full">
      <div className="p-6 pb-4 border-b border-slate-800/50">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-cyan-500/20 border border-cyan-500/30 flex items-center justify-center">
            <Activity className="w-5 h-5 text-cyan-400" />
          </div>
          <div>
            <h1 className="text-xl font-semibold text-slate-100">SDR Monitor</h1>
            <p className="text-sm text-slate-400">{totalFreqs} frequencies · {bandPlans.length} band plans · Reference mode</p>
          </div>
        </div>

        <div className="flex gap-2 mt-4">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500" />
            <input type="text" value={search} onChange={e => setSearch(e.target.value)} onKeyDown={e => e.key === 'Enter' && doSearch()}
              placeholder="Search frequencies..." className="w-full pl-9 pr-3 py-2 rounded-lg bg-slate-800/50 border border-slate-700/50 text-sm text-slate-200 placeholder-slate-500 focus:outline-none focus:border-cyan-500/50" />
          </div>
          <button onClick={doSearch} className="px-4 py-2 rounded-lg bg-cyan-500 hover:bg-cyan-600 text-white text-sm font-semibold transition-colors">Search</button>
          {searchResults && <button onClick={() => { setSearch(''); setSearchResults(null) }} className="px-3 py-2 rounded-lg bg-slate-800 text-xs text-slate-400 hover:text-slate-200">Clear</button>}
        </div>

        <div className="flex gap-1 mt-3 p-1 rounded-lg bg-slate-800/50">
          {([['frequencies', 'Frequencies', Radio], ['bandplans', 'Band Plans', Activity]] as const).map(([id, label, Icon]) => (
            <button key={id} onClick={() => setActiveView(id)}
              className={`flex-1 flex items-center justify-center gap-2 px-3 py-2 rounded-md text-sm font-medium transition-all ${
                activeView === id ? 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/30' : 'text-slate-500 hover:text-slate-300'
              }`}><Icon className="w-4 h-4" /> {label}</button>
          ))}
        </div>

        <div className="mt-3 p-3 rounded-lg bg-amber-500/5 border border-amber-500/20">
          <div className="flex items-start gap-2">
            <AlertTriangle className="w-4 h-4 text-amber-400 shrink-0 mt-0.5" />
            <p className="text-xs text-amber-300">
              <strong>Reference Mode:</strong> No SDR hardware detected. Connect an RTL-SDR dongle to enable live monitoring.
              Frequency data is available offline for reference.
            </p>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-6">
        {searchResults ? (
          <div className="space-y-2">
            <p className="text-sm text-slate-400 mb-3">{searchResults.length} results for "{search}"</p>
            {searchResults.map(freq => (
              <div key={freq.id} className="p-3 rounded-xl bg-slate-800/30 border border-slate-700/30">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-mono font-bold text-cyan-400">{freq.frequency}</span>
                  <span className="px-1.5 py-0.5 rounded text-[10px] bg-slate-700/50 text-slate-400 font-mono">{freq.mode}</span>
                </div>
                <p className="text-xs text-slate-200 mt-0.5">{freq.name}</p>
              </div>
            ))}
          </div>
        ) : activeView === 'frequencies' ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {groups.map(g => (
              <button key={g.id} onClick={() => setActiveGroup(g)}
                className="text-left p-5 rounded-xl bg-slate-800/30 border border-slate-700/30 hover:border-cyan-500/30 hover:bg-cyan-500/5 transition-all group">
                <div className="flex items-center justify-between mb-2">
                  <h3 className="text-sm font-semibold text-slate-200 group-hover:text-cyan-300">{g.name}</h3>
                  <ChevronRight className="w-4 h-4 text-slate-600 group-hover:text-cyan-400" />
                </div>
                <p className="text-xs text-slate-500 mb-2">{g.description}</p>
                <span className="text-[10px] px-2 py-0.5 rounded-full bg-slate-700/50 text-slate-400">{g.frequencies.length} frequencies</span>
              </button>
            ))}
          </div>
        ) : (
          <div className="space-y-3">
            {bandPlans.map(bp => (
              <div key={bp.id} className="p-4 rounded-xl bg-slate-800/30 border border-slate-700/30">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-semibold text-slate-200">{bp.name}</span>
                  <span className="text-xs font-mono text-cyan-400">{bp.start_mhz} – {bp.end_mhz} MHz</span>
                </div>
                <p className="text-xs text-amber-400 mt-1">{bp.allocation}</p>
                <p className="text-xs text-slate-500 mt-0.5">{bp.description}</p>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
