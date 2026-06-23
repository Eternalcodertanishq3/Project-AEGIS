import { useState, useEffect } from 'react'
import {
  Leaf, Loader2, ArrowLeft, ChevronRight,
  Search, AlertTriangle, Skull, Check
} from 'lucide-react'

interface Plant {
  id: string
  common_name: string
  scientific_name: string
  edibility: string
  category: string
  habitat: string
  season: string
  description: string
  identification: string[]
  preparation?: string[]
  warnings?: string[]
  look_alikes?: string[]
}

interface PlantGroup {
  id: string
  name: string
  icon: string
  description: string
  plants: Plant[]
}

const edibilityStyles: Record<string, { bg: string; border: string; text: string; label: string; icon: React.ElementType }> = {
  'edible': { bg: 'bg-emerald-500/10', border: 'border-emerald-500/30', text: 'text-emerald-400', label: 'EDIBLE', icon: Check },
  'edible-caution': { bg: 'bg-amber-500/10', border: 'border-amber-500/30', text: 'text-amber-400', label: 'EDIBLE — CAUTION', icon: AlertTriangle },
  'poisonous': { bg: 'bg-red-500/10', border: 'border-red-500/30', text: 'text-red-400', label: 'POISONOUS', icon: Skull },
  'deadly': { bg: 'bg-red-500/20', border: 'border-red-500/50', text: 'text-red-400', label: '☠ DEADLY', icon: Skull },
}

export function PlantIDPage() {
  const [groups, setGroups] = useState<PlantGroup[]>([])
  const [loading, setLoading] = useState(true)
  const [activeGroup, setActiveGroup] = useState<PlantGroup | null>(null)
  const [activePlant, setActivePlant] = useState<Plant | null>(null)
  const [search, setSearch] = useState('')
  const [searchResults, setSearchResults] = useState<Plant[] | null>(null)
  const [searching, setSearching] = useState(false)

  useEffect(() => {
    const fetchGroups = async () => {
      try {
        const res = await fetch('/api/plants/groups')
        const data = await res.json()
        setGroups(data.groups || [])
      } catch (e) {
        console.error('Failed to fetch plant groups:', e)
      } finally {
        setLoading(false)
      }
    }
    fetchGroups()
  }, [])

  const doSearch = async () => {
    if (!search.trim()) return
    setSearching(true)
    try {
      const res = await fetch(`/api/plants/search?q=${encodeURIComponent(search)}`)
      const data = await res.json()
      setSearchResults(data.results || [])
    } catch (e) {
      console.error('Search failed:', e)
    } finally {
      setSearching(false)
    }
  }

  const clearSearch = () => {
    setSearch('')
    setSearchResults(null)
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="w-8 h-8 text-green-500 animate-spin" />
      </div>
    )
  }

  // Plant detail view
  if (activePlant) {
    const edib = edibilityStyles[activePlant.edibility] || edibilityStyles.edible
    const EdibIcon = edib.icon
    return (
      <div className="flex flex-col h-full">
        <div className="p-6 pb-4 border-b border-slate-800/50">
          <button
            onClick={() => setActivePlant(null)}
            className="flex items-center gap-2 text-sm text-slate-400 hover:text-slate-200 mb-3 transition-colors"
          >
            <ArrowLeft className="w-4 h-4" /> Back
          </button>
          <h1 className="text-xl font-semibold text-slate-100">{activePlant.common_name}</h1>
          <p className="text-sm text-slate-500 italic">{activePlant.scientific_name}</p>
          <div className="flex items-center gap-3 mt-3">
            <span className={`px-2.5 py-0.5 rounded-full text-xs font-bold uppercase ${edib.bg} ${edib.border} border ${edib.text} flex items-center gap-1`}>
              <EdibIcon className="w-3 h-3" /> {edib.label}
            </span>
            <span className="text-xs text-slate-500">{activePlant.season}</span>
          </div>
          <p className="text-sm text-slate-400 mt-2">{activePlant.description}</p>
        </div>

        <div className="flex-1 overflow-y-auto p-6 space-y-6">
          {/* Habitat */}
          <div>
            <h2 className="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Habitat</h2>
            <p className="text-sm text-slate-300">{activePlant.habitat}</p>
          </div>

          {/* Identification */}
          <div>
            <h2 className="text-xs font-semibold text-emerald-400 uppercase tracking-wider mb-2">🔍 How to Identify</h2>
            <div className="space-y-1.5">
              {activePlant.identification.map((id, i) => (
                <div key={i} className="flex gap-2 p-2.5 rounded-lg bg-slate-800/20 border border-slate-700/20">
                  <span className="text-emerald-400 shrink-0">•</span>
                  <p className="text-sm text-slate-300">{id}</p>
                </div>
              ))}
            </div>
          </div>

          {/* Preparation */}
          {activePlant.preparation && activePlant.preparation.length > 0 && (
            <div>
              <h2 className="text-xs font-semibold text-amber-400 uppercase tracking-wider mb-2">🍽 Preparation</h2>
              <div className="space-y-1.5">
                {activePlant.preparation.map((p, i) => (
                  <div key={i} className="flex gap-2 p-2.5 rounded-lg bg-amber-500/5 border border-amber-500/10">
                    <span className="text-amber-400 shrink-0">•</span>
                    <p className="text-sm text-amber-300/80">{p}</p>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Warnings */}
          {activePlant.warnings && activePlant.warnings.length > 0 && (
            <div>
              <h2 className="text-xs font-semibold text-red-400 uppercase tracking-wider mb-2">⚠ Warnings</h2>
              <div className="space-y-1.5">
                {activePlant.warnings.map((w, i) => (
                  <div key={i} className="flex gap-2 p-2.5 rounded-lg bg-red-500/5 border border-red-500/20">
                    <Skull className="w-4 h-4 text-red-400 shrink-0 mt-0.5" />
                    <p className="text-sm text-red-300">{w}</p>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Look-alikes */}
          {activePlant.look_alikes && activePlant.look_alikes.length > 0 && (
            <div>
              <h2 className="text-xs font-semibold text-orange-400 uppercase tracking-wider mb-2">👀 Dangerous Look-Alikes</h2>
              <div className="space-y-1.5">
                {activePlant.look_alikes.map((la, i) => (
                  <div key={i} className="flex gap-2 p-2.5 rounded-lg bg-orange-500/5 border border-orange-500/20">
                    <AlertTriangle className="w-4 h-4 text-orange-400 shrink-0 mt-0.5" />
                    <p className="text-sm text-orange-300/80">{la}</p>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    )
  }

  // Group detail view
  if (activeGroup) {
    return (
      <div className="flex flex-col h-full">
        <div className="p-6 pb-4 border-b border-slate-800/50">
          <button
            onClick={() => setActiveGroup(null)}
            className="flex items-center gap-2 text-sm text-slate-400 hover:text-slate-200 mb-3 transition-colors"
          >
            <ArrowLeft className="w-4 h-4" /> All Groups
          </button>
          <h1 className="text-xl font-semibold text-slate-100">{activeGroup.name}</h1>
          <p className="text-sm text-slate-400 mt-1">{activeGroup.description}</p>
        </div>
        <div className="flex-1 overflow-y-auto p-6 space-y-3">
          {activeGroup.plants.map(plant => {
            const edib = edibilityStyles[plant.edibility] || edibilityStyles.edible
            const EdibIcon = edib.icon
            return (
              <button
                key={plant.id}
                onClick={() => setActivePlant(plant)}
                className="w-full text-left p-4 rounded-xl bg-slate-800/30 border border-slate-700/30 hover:border-green-500/30 transition-all group"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <span className="text-sm font-semibold text-slate-200">{plant.common_name}</span>
                    <span className="text-xs text-slate-600 italic ml-2">{plant.scientific_name}</span>
                  </div>
                  <ChevronRight className="w-4 h-4 text-slate-600 group-hover:text-green-400" />
                </div>
                <p className="text-xs text-slate-500 mt-1">{plant.description.substring(0, 100)}...</p>
                <div className="flex items-center gap-2 mt-2">
                  <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase ${edib.bg} ${edib.border} border ${edib.text} flex items-center gap-1`}>
                    <EdibIcon className="w-2.5 h-2.5" /> {edib.label}
                  </span>
                  <span className="text-[10px] text-slate-600">{plant.season}</span>
                </div>
              </button>
            )
          })}
        </div>
      </div>
    )
  }

  // Main view — Groups + Search
  const totalPlants = groups.reduce((sum, g) => sum + g.plants.length, 0)

  return (
    <div className="flex flex-col h-full">
      <div className="p-6 pb-4 border-b border-slate-800/50">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-green-500/20 border border-green-500/30 flex items-center justify-center">
            <Leaf className="w-5 h-5 text-green-400" />
          </div>
          <div>
            <h1 className="text-xl font-semibold text-slate-100">Plant & Fungi ID</h1>
            <p className="text-sm text-slate-400">
              {groups.length} groups · {totalPlants} species · Offline field guide
            </p>
          </div>
        </div>

        {/* Search */}
        <div className="flex gap-2 mt-4">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500" />
            <input
              type="text" value={search}
              onChange={e => setSearch(e.target.value)}
              onKeyDown={e => e.key === 'Enter' && doSearch()}
              placeholder="Search by name, habitat, or description..."
              className="w-full pl-9 pr-3 py-2 rounded-lg bg-slate-800/50 border border-slate-700/50 text-sm text-slate-200 placeholder-slate-500 focus:outline-none focus:border-green-500/50"
            />
          </div>
          <button onClick={doSearch} disabled={searching || !search.trim()}
            className="px-4 py-2 rounded-lg bg-green-500 hover:bg-green-600 text-white text-sm font-semibold transition-colors disabled:opacity-50">
            {searching ? <Loader2 className="w-4 h-4 animate-spin" /> : 'Search'}
          </button>
          {searchResults !== null && (
            <button onClick={clearSearch} className="px-3 py-2 rounded-lg bg-slate-800 text-xs text-slate-400 hover:text-slate-200 transition-colors">
              Clear
            </button>
          )}
        </div>

        <div className="mt-3 p-3 rounded-lg bg-red-500/5 border border-red-500/20">
          <div className="flex items-start gap-2">
            <Skull className="w-4 h-4 text-red-400 shrink-0 mt-0.5" />
            <p className="text-xs text-red-300">
              <strong>WARNING:</strong> Never eat anything you cannot identify with 100% certainty.
              When in doubt, DO NOT eat it. Many edible plants have deadly look-alikes.
            </p>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-6">
        {/* Search results */}
        {searchResults !== null ? (
          <div className="space-y-3">
            <p className="text-sm text-slate-400 mb-3">{searchResults.length} result{searchResults.length !== 1 ? 's' : ''} for "{search}"</p>
            {searchResults.map(plant => {
              const edib = edibilityStyles[plant.edibility] || edibilityStyles.edible
              return (
                <button key={plant.id} onClick={() => setActivePlant(plant)}
                  className="w-full text-left p-4 rounded-xl bg-slate-800/30 border border-slate-700/30 hover:border-green-500/30 transition-all group">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-semibold text-slate-200">{plant.common_name}</span>
                    <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase ${edib.bg} ${edib.text}`}>{edib.label}</span>
                  </div>
                  <p className="text-xs text-slate-500 italic">{plant.scientific_name}</p>
                </button>
              )
            })}
          </div>
        ) : (
          /* Groups grid */
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            {groups.map(group => (
              <button
                key={group.id}
                onClick={() => setActiveGroup(group)}
                className="text-left p-5 rounded-xl bg-slate-800/30 border border-slate-700/30 hover:border-green-500/30 hover:bg-green-500/5 transition-all group"
              >
                <div className="flex items-center justify-between mb-3">
                  <h3 className="text-sm font-semibold text-slate-200 group-hover:text-green-300">{group.name}</h3>
                  <ChevronRight className="w-4 h-4 text-slate-600 group-hover:text-green-400" />
                </div>
                <p className="text-xs text-slate-500 mb-3">{group.description}</p>
                <span className="text-[10px] px-2 py-0.5 rounded-full bg-slate-700/50 text-slate-400">
                  {group.plants.length} species
                </span>
              </button>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
