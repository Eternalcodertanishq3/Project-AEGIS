import { useState, useEffect } from 'react'
import { Navigation, Loader2, MapPin, Crosshair, ArrowRight } from 'lucide-react'

interface Position { id: string; latitude: number; longitude: number; altitude: number; source: string; comment: string; timestamp: string }
interface APRSBeacon { callsign: string; latitude: number; longitude: number; comment: string; raw: string; timestamp: string }
interface DistanceResult { from_lat: number; from_lon: number; to_lat: number; to_lon: number; distance_km: number; distance_mi: number; distance_nm: number; bearing_deg: number; bearing_cardinal: string }
interface BeaconStatus { callsign: string; beacon_active: boolean; position_count: number; transport_mode: string }

export function BeaconPage() {
  const [status, setStatus] = useState<BeaconStatus | null>(null)
  const [positions, setPositions] = useState<Position[]>([])
  const [activeView, setActiveView] = useState<'log' | 'aprs' | 'distance'>('log')
  const [lat, setLat] = useState('28.6139')
  const [lon, setLon] = useState('77.2090')
  const [comment, setComment] = useState('')
  const [aprsResult, setAprsResult] = useState<APRSBeacon | null>(null)
  const [distResult, setDistResult] = useState<DistanceResult | null>(null)
  const [lat2, setLat2] = useState(''); const [lon2, setLon2] = useState('')
  const [loading, setLoading] = useState(true)

  const refresh = () => {
    Promise.all([
      fetch('/api/beacon/status').then(r => r.json()),
      fetch('/api/beacon/positions?limit=20').then(r => r.json()),
    ]).then(([s, p]) => { setStatus(s); setPositions(p.positions || []) }).finally(() => setLoading(false))
  }

  useEffect(() => { refresh() }, [])

  const logPosition = async () => {
    await fetch('/api/beacon/positions', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ latitude: parseFloat(lat), longitude: parseFloat(lon), source: 'manual', comment }),
    })
    setComment(''); refresh()
  }

  const generateAPRS = async () => {
    const res = await fetch('/api/beacon/aprs', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ latitude: parseFloat(lat), longitude: parseFloat(lon), comment: comment || 'AEGIS Node' }),
    })
    setAprsResult(await res.json())
  }

  const calcDistance = async () => {
    if (!lat2 || !lon2) return
    const res = await fetch(`/api/beacon/distance?lat1=${lat}&lon1=${lon}&lat2=${lat2}&lon2=${lon2}`)
    setDistResult(await res.json())
  }

  const useGPS = () => {
    if (navigator.geolocation) navigator.geolocation.getCurrentPosition(p => { setLat(p.coords.latitude.toFixed(4)); setLon(p.coords.longitude.toFixed(4)) })
  }

  if (loading) return <div className="flex items-center justify-center h-full"><Loader2 className="w-8 h-8 text-orange-500 animate-spin" /></div>

  return (
    <div className="flex flex-col h-full">
      <div className="p-6 pb-4 border-b border-slate-800/50">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-orange-500/20 border border-orange-500/30 flex items-center justify-center">
            <Navigation className="w-5 h-5 text-orange-400" />
          </div>
          <div>
            <h1 className="text-xl font-semibold text-slate-100">Position Beacon</h1>
            <p className="text-sm text-slate-400">
              {status?.callsign} · {status?.position_count} positions logged · {status?.transport_mode}
            </p>
          </div>
        </div>

        {/* Tabs */}
        <div className="flex gap-1 mt-4 p-1 rounded-lg bg-slate-800/50">
          {([['log', 'Position Log', MapPin], ['aprs', 'APRS Beacon', Crosshair], ['distance', 'Distance Calc', ArrowRight]] as const).map(([id, label, Icon]) => (
            <button key={id} onClick={() => setActiveView(id)}
              className={`flex-1 flex items-center justify-center gap-2 px-3 py-2 rounded-md text-sm font-medium transition-all ${
                activeView === id ? 'bg-orange-500/20 text-orange-400 border border-orange-500/30' : 'text-slate-500 hover:text-slate-300'
              }`}><Icon className="w-4 h-4" /> {label}</button>
          ))}
        </div>

        {/* Coordinate inputs */}
        <div className="flex gap-2 mt-4 items-end">
          <div className="flex-1">
            <label className="text-[10px] font-semibold text-slate-500 uppercase">Latitude</label>
            <input type="text" value={lat} onChange={e => setLat(e.target.value)}
              className="w-full mt-0.5 px-3 py-1.5 rounded-lg bg-slate-800/50 border border-slate-700/50 text-sm text-slate-200 font-mono focus:outline-none focus:border-orange-500/50" />
          </div>
          <div className="flex-1">
            <label className="text-[10px] font-semibold text-slate-500 uppercase">Longitude</label>
            <input type="text" value={lon} onChange={e => setLon(e.target.value)}
              className="w-full mt-0.5 px-3 py-1.5 rounded-lg bg-slate-800/50 border border-slate-700/50 text-sm text-slate-200 font-mono focus:outline-none focus:border-orange-500/50" />
          </div>
          <button onClick={useGPS} className="px-3 py-1.5 rounded-lg bg-slate-800 text-xs text-slate-400 hover:text-slate-200 transition-colors">📍 GPS</button>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-6 max-w-3xl">
        {activeView === 'log' && (
          <div className="space-y-4">
            <div className="flex gap-2">
              <input type="text" value={comment} onChange={e => setComment(e.target.value)} placeholder="Position comment (optional)"
                className="flex-1 px-3 py-2 rounded-lg bg-slate-800/50 border border-slate-700/50 text-sm text-slate-200 placeholder-slate-500 focus:outline-none focus:border-orange-500/50" />
              <button onClick={logPosition}
                className="px-4 py-2 rounded-lg bg-orange-500 hover:bg-orange-600 text-white text-sm font-semibold transition-colors">Log Position</button>
            </div>

            {positions.length === 0 ? (
              <div className="text-center py-8 rounded-xl bg-slate-800/20"><MapPin className="w-8 h-8 text-slate-700 mx-auto mb-2" /><p className="text-sm text-slate-600">No positions logged yet</p></div>
            ) : (
              <div className="space-y-2">
                {positions.map(pos => (
                  <div key={pos.id} className="p-3 rounded-xl bg-slate-800/30 border border-slate-700/30">
                    <div className="flex items-center justify-between">
                      <span className="text-sm font-mono text-orange-400">{pos.latitude.toFixed(4)}, {pos.longitude.toFixed(4)}</span>
                      <span className="text-[10px] text-slate-500">{new Date(pos.timestamp).toLocaleString()}</span>
                    </div>
                    {pos.comment && <p className="text-xs text-slate-400 mt-0.5">{pos.comment}</p>}
                    <span className="text-[10px] px-1.5 py-0.5 rounded bg-slate-800 text-slate-500">{pos.source}</span>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {activeView === 'aprs' && (
          <div className="space-y-4">
            <div className="flex gap-2">
              <input type="text" value={comment} onChange={e => setComment(e.target.value)} placeholder="Beacon comment"
                className="flex-1 px-3 py-2 rounded-lg bg-slate-800/50 border border-slate-700/50 text-sm text-slate-200 placeholder-slate-500 focus:outline-none" />
              <button onClick={generateAPRS}
                className="px-4 py-2 rounded-lg bg-orange-500 hover:bg-orange-600 text-white text-sm font-semibold transition-colors">Generate APRS</button>
            </div>
            {aprsResult && (
              <div className="p-4 rounded-xl bg-slate-800/30 border border-orange-500/30">
                <p className="text-[10px] font-semibold text-slate-500 uppercase mb-1">APRS Beacon String</p>
                <p className="text-sm font-mono text-orange-400 bg-slate-900/50 p-3 rounded-lg break-all">{aprsResult.raw}</p>
                <div className="grid grid-cols-2 gap-2 mt-3">
                  <div><p className="text-[10px] text-slate-500">Callsign</p><p className="text-sm text-slate-200">{aprsResult.callsign}</p></div>
                  <div><p className="text-[10px] text-slate-500">Position</p><p className="text-sm text-slate-200 font-mono">{aprsResult.latitude}, {aprsResult.longitude}</p></div>
                </div>
              </div>
            )}
          </div>
        )}

        {activeView === 'distance' && (
          <div className="space-y-4">
            <div className="flex gap-2 items-end">
              <div className="flex-1">
                <label className="text-[10px] font-semibold text-slate-500 uppercase">Target Lat</label>
                <input type="text" value={lat2} onChange={e => setLat2(e.target.value)}
                  className="w-full mt-0.5 px-3 py-1.5 rounded-lg bg-slate-800/50 border border-slate-700/50 text-sm text-slate-200 font-mono focus:outline-none" />
              </div>
              <div className="flex-1">
                <label className="text-[10px] font-semibold text-slate-500 uppercase">Target Lon</label>
                <input type="text" value={lon2} onChange={e => setLon2(e.target.value)}
                  className="w-full mt-0.5 px-3 py-1.5 rounded-lg bg-slate-800/50 border border-slate-700/50 text-sm text-slate-200 font-mono focus:outline-none" />
              </div>
              <button onClick={calcDistance}
                className="px-4 py-1.5 rounded-lg bg-orange-500 hover:bg-orange-600 text-white text-sm font-semibold transition-colors">Calculate</button>
            </div>
            {distResult && (
              <div className="grid grid-cols-2 gap-3">
                <InfoCard label="Distance" value={`${distResult.distance_km} km`} sub={`${distResult.distance_mi} mi · ${distResult.distance_nm} nm`} />
                <InfoCard label="Bearing" value={`${distResult.bearing_deg}° ${distResult.bearing_cardinal}`} sub="From true north" />
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}

function InfoCard({ label, value, sub }: { label: string; value: string; sub: string }) {
  return (
    <div className="p-4 rounded-xl bg-slate-800/30 border border-slate-700/30">
      <p className="text-[10px] font-semibold text-slate-500 uppercase tracking-wider">{label}</p>
      <p className="text-lg font-semibold text-slate-100 mt-1">{value}</p>
      <p className="text-xs text-slate-500">{sub}</p>
    </div>
  )
}
