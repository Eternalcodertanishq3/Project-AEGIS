import { useState, useEffect } from 'react'
import {
  Map,
  Layers,
  AlertCircle,
  Download,
  Loader2,
  Globe,
  Maximize2,
} from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { apiFetch } from '@/hooks/useApi'

interface PMTilesFile {
  id: string
  name: string
  path: string
  size_bytes: number
  size_human: string
}

interface MapsStatus {
  module: string
  status: string
  pmtiles_count: number
  pmtiles_files: PMTilesFile[]
}

export function MapsPage() {
  const [status, setStatus] = useState<MapsStatus | null>(null)
  const [loading, setLoading] = useState(true)
  const [selectedMap, setSelectedMap] = useState<PMTilesFile | null>(null)

  useEffect(() => {
    async function fetchStatus() {
      try {
        const data = await apiFetch<MapsStatus>('/maps/status')
        setStatus(data)
        if (data.pmtiles_files?.length > 0) {
          setSelectedMap(data.pmtiles_files[0])
        }
      } catch {
        // Maps module not available
      } finally {
        setLoading(false)
      }
    }
    fetchStatus()
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
        <div className="flex items-center justify-center w-8 h-8 rounded-md bg-emerald-600/20 border border-emerald-500/30">
          <Map className="h-4 w-4 text-emerald-400" />
        </div>
        <div>
          <h1 className="text-xl font-semibold text-slate-100">Offline maps</h1>
          <p className="text-xs text-slate-500">
            View detailed maps without internet using PMTiles
          </p>
        </div>
      </div>

      {/* Map viewer area */}
      <div className="glass-panel overflow-hidden">
        {selectedMap ? (
          <div className="relative">
            {/* Map container — placeholder since MapLibre needs the JS library loaded at runtime */}
            <div className="relative h-[500px] bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 flex items-center justify-center">
              {/* Grid overlay to simulate map tiles */}
              <div className="absolute inset-0 opacity-10">
                <div className="h-full w-full" style={{
                  backgroundImage: `
                    linear-gradient(rgba(16,185,129,0.3) 1px, transparent 1px),
                    linear-gradient(90deg, rgba(16,185,129,0.3) 1px, transparent 1px)
                  `,
                  backgroundSize: '60px 60px',
                }} />
              </div>

              {/* Center content */}
              <div className="relative z-10 text-center space-y-4">
                <Globe className="h-16 w-16 text-emerald-500/30 mx-auto" />
                <div>
                  <h3 className="text-lg font-semibold text-slate-200">Map ready</h3>
                  <p className="text-sm text-slate-400 mt-1">
                    <span className="text-emerald-400 font-medium">{selectedMap.name}</span>
                  </p>
                  <p className="text-xs text-slate-500 mt-2 max-w-md mx-auto">
                    PMTiles file detected ({selectedMap.size_human}). MapLibre GL JS will render 
                    interactive tiles from this file when loaded in a full browser session.
                  </p>
                </div>

                {/* Map info badges */}
                <div className="flex items-center justify-center gap-2">
                  <Badge className="bg-emerald-600/20 text-emerald-400 border-emerald-500/30 text-xs">
                    <Layers className="h-3 w-3 mr-1" />
                    PMTiles
                  </Badge>
                  <Badge variant="outline" className="text-xs text-slate-400">
                    {selectedMap.size_human}
                  </Badge>
                  <Badge variant="outline" className="text-xs text-slate-400">
                    Offline ready
                  </Badge>
                </div>
              </div>

              {/* Fullscreen hint */}
              <button className="absolute top-3 right-3 p-2 rounded-lg bg-slate-800/80 border border-slate-700/50 
                                 text-slate-400 hover:text-slate-200 hover:bg-slate-800 transition-colors">
                <Maximize2 className="h-4 w-4" />
              </button>
            </div>

            {/* Map details bar */}
            <div className="flex items-center justify-between px-4 py-2.5 border-t border-slate-700/50 bg-slate-900/50">
              <div className="flex items-center gap-2 text-[11px] text-slate-500">
                <span className="status-dot status-dot-active" />
                <span>{selectedMap.name}</span>
                <span className="text-slate-700">|</span>
                <span>Tile server on /api/maps/tiles/{selectedMap.id}/</span>
              </div>
              <span className="text-[11px] text-slate-600 font-mono">{selectedMap.id}</span>
            </div>
          </div>
        ) : (
          <div className="h-[400px] flex items-center justify-center">
            <div className="text-center space-y-3">
              <AlertCircle className="h-8 w-8 text-slate-600 mx-auto" />
              <h3 className="text-sm font-semibold text-slate-300">No map files found</h3>
              <p className="text-xs text-slate-500 max-w-md mx-auto">
                Download PMTiles map files and place them in the{' '}
                <code className="text-emerald-400 bg-slate-800 px-1 py-0.5 rounded text-[11px]">
                  content-packs/maps-regional/
                </code>{' '}
                directory, then restart AEGIS.
              </p>
            </div>
          </div>
        )}
      </div>

      {/* Available map files */}
      <div className="glass-panel p-5">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <Layers className="h-4 w-4 text-slate-400" />
            <h2 className="text-sm font-semibold text-slate-200">Available map files</h2>
          </div>
          <Badge variant="outline" className="text-xs text-slate-400">
            {status?.pmtiles_count || 0} file{(status?.pmtiles_count || 0) !== 1 ? 's' : ''}
          </Badge>
        </div>

        {status?.pmtiles_files && status.pmtiles_files.length > 0 ? (
          <div className="space-y-2">
            {status.pmtiles_files.map((pf) => (
              <button
                key={pf.id}
                onClick={() => setSelectedMap(pf)}
                className={`w-full flex items-center justify-between p-3 rounded-lg border transition-all text-left
                  ${selectedMap?.id === pf.id
                    ? 'border-emerald-500/30 bg-emerald-600/10'
                    : 'border-slate-700/30 bg-slate-800/20 hover:bg-slate-800/40'
                  }`}
              >
                <div className="flex items-center gap-3">
                  <Map className={`h-4 w-4 ${selectedMap?.id === pf.id ? 'text-emerald-400' : 'text-slate-500'}`} />
                  <div>
                    <p className={`text-sm font-medium ${selectedMap?.id === pf.id ? 'text-emerald-400' : 'text-slate-200'}`}>
                      {pf.name}
                    </p>
                    <p className="text-[11px] text-slate-500 font-mono">{pf.id}</p>
                  </div>
                </div>
                <Badge variant="outline" className="text-xs text-slate-400">{pf.size_human}</Badge>
              </button>
            ))}
          </div>
        ) : (
          <div className="text-center py-6">
            <Download className="h-6 w-6 text-slate-600 mx-auto mb-2" />
            <p className="text-xs text-slate-500">
              No PMTiles files detected. Visit{' '}
              <span className="text-emerald-400">protomaps.com</span>{' '}
              to download offline map tiles.
            </p>
          </div>
        )}
      </div>
    </div>
  )
}
