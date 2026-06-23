import { useState, useEffect } from 'react'
import {
  Star, Loader2, Sun, Moon, Compass,
  ChevronRight, ArrowLeft, Navigation
} from 'lucide-react'

interface CelestialResult {
  date_time: string
  latitude: number
  longitude: number
  sun_azimuth: number
  sun_altitude: number
  sun_rise: string
  sun_set: string
  moon_phase: string
  moon_phase_angle: number
  polaris_altitude: number
  true_north_hint: string
  day_length: string
}

interface NavStar {
  name: string
  constellation: string
  description: string
  magnitude: number
  usage: string
}

interface Technique {
  id: string
  name: string
  hemisphere: string
  description: string
  steps: string[]
}

export function CelestialPage() {
  const [lat, setLat] = useState('28.6139')
  const [lon, setLon] = useState('77.2090')
  const [result, setResult] = useState<CelestialResult | null>(null)
  const [stars, setStars] = useState<NavStar[]>([])
  const [techniques, setTechniques] = useState<Technique[]>([])
  const [loading, setLoading] = useState(false)
  const [activeView, setActiveView] = useState<'calculator' | 'stars' | 'techniques'>('calculator')
  const [activeTechnique, setActiveTechnique] = useState<Technique | null>(null)

  useEffect(() => {
    fetch('/api/celestial/stars')
      .then(r => r.json())
      .then(d => setStars(d.stars || []))
      .catch(console.error)
    fetch('/api/celestial/techniques')
      .then(r => r.json())
      .then(d => setTechniques(d.techniques || []))
      .catch(console.error)
  }, [])

  const calculate = async () => {
    setLoading(true)
    try {
      const res = await fetch(`/api/celestial/calculate?lat=${lat}&lon=${lon}`)
      const data = await res.json()
      setResult(data)
    } catch (e) {
      console.error('Calculation failed:', e)
    } finally {
      setLoading(false)
    }
  }

  const useGPS = () => {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (pos) => {
          setLat(pos.coords.latitude.toFixed(4))
          setLon(pos.coords.longitude.toFixed(4))
        },
        () => { /* GPS not available offline — that's expected */ }
      )
    }
  }

  const hemisphereLabel: Record<string, string> = {
    north: '🌍 Northern Hemisphere',
    south: '🌏 Southern Hemisphere',
    both: '🌐 Both Hemispheres',
  }

  // Technique detail
  if (activeTechnique) {
    return (
      <div className="flex flex-col h-full">
        <div className="p-6 pb-4 border-b border-slate-800/50">
          <button
            onClick={() => setActiveTechnique(null)}
            className="flex items-center gap-2 text-sm text-slate-400 hover:text-slate-200 mb-3 transition-colors"
          >
            <ArrowLeft className="w-4 h-4" /> All Techniques
          </button>
          <h1 className="text-xl font-semibold text-slate-100">{activeTechnique.name}</h1>
          <p className="text-sm text-slate-400 mt-1">{activeTechnique.description}</p>
          <span className="inline-block mt-2 px-2.5 py-0.5 rounded-full text-xs bg-indigo-500/10 border border-indigo-500/30 text-indigo-400">
            {hemisphereLabel[activeTechnique.hemisphere]}
          </span>
        </div>
        <div className="flex-1 overflow-y-auto p-6 space-y-2">
          {activeTechnique.steps.map((step, i) => (
            <div key={i} className="flex gap-3 p-3 rounded-lg bg-slate-800/30 border border-slate-700/30">
              <div className="w-7 h-7 rounded-full bg-indigo-500/20 border border-indigo-500/30 flex items-center justify-center shrink-0">
                <span className="text-xs font-bold text-indigo-400">{i + 1}</span>
              </div>
              <p className="text-sm text-slate-200 leading-relaxed pt-1">{step}</p>
            </div>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="p-6 pb-4 border-b border-slate-800/50">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-indigo-500/20 border border-indigo-500/30 flex items-center justify-center">
            <Navigation className="w-5 h-5 text-indigo-400" />
          </div>
          <div>
            <h1 className="text-xl font-semibold text-slate-100">Celestial Navigation</h1>
            <p className="text-sm text-slate-400">Sun position, moon phase, and star-based navigation</p>
          </div>
        </div>

        {/* Tab bar */}
        <div className="flex gap-1 mt-4 p-1 rounded-lg bg-slate-800/50">
          {([['calculator', 'Calculator', Sun], ['stars', 'Stars', Star], ['techniques', 'Techniques', Compass]] as const).map(([id, label, Icon]) => (
            <button
              key={id}
              onClick={() => setActiveView(id)}
              className={`flex-1 flex items-center justify-center gap-2 px-3 py-2 rounded-md text-sm font-medium transition-all ${
                activeView === id
                  ? 'bg-indigo-500/20 text-indigo-400 border border-indigo-500/30'
                  : 'text-slate-500 hover:text-slate-300'
              }`}
            >
              <Icon className="w-4 h-4" /> {label}
            </button>
          ))}
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-6">
        {/* Calculator Tab */}
        {activeView === 'calculator' && (
          <div className="space-y-6 max-w-2xl">
            <div className="flex gap-3 items-end">
              <div className="flex-1">
                <label className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Latitude</label>
                <input
                  type="text" value={lat} onChange={e => setLat(e.target.value)}
                  className="w-full mt-1 px-3 py-2 rounded-lg bg-slate-800/50 border border-slate-700/50 text-sm text-slate-200 focus:outline-none focus:border-indigo-500/50"
                />
              </div>
              <div className="flex-1">
                <label className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Longitude</label>
                <input
                  type="text" value={lon} onChange={e => setLon(e.target.value)}
                  className="w-full mt-1 px-3 py-2 rounded-lg bg-slate-800/50 border border-slate-700/50 text-sm text-slate-200 focus:outline-none focus:border-indigo-500/50"
                />
              </div>
              <button onClick={useGPS} className="px-3 py-2 rounded-lg bg-slate-800 text-xs text-slate-400 hover:text-slate-200 transition-colors whitespace-nowrap">
                📍 GPS
              </button>
              <button
                onClick={calculate} disabled={loading}
                className="px-4 py-2 rounded-lg bg-indigo-500 hover:bg-indigo-600 text-white text-sm font-semibold transition-colors disabled:opacity-50"
              >
                {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : 'Calculate'}
              </button>
            </div>

            {result && (
              <div className="grid grid-cols-2 gap-3">
                <InfoCard icon={<Sun className="w-5 h-5 text-amber-400" />} label="Sun Altitude" value={`${result.sun_altitude}°`} sub={result.sun_altitude > 0 ? 'Above horizon' : 'Below horizon'} />
                <InfoCard icon={<Compass className="w-5 h-5 text-blue-400" />} label="Sun Azimuth" value={`${result.sun_azimuth}°`} sub={`Bearing from true north`} />
                <InfoCard icon={<Sun className="w-5 h-5 text-orange-400" />} label="Sunrise" value={result.sun_rise} sub={`Sunset: ${result.sun_set}`} />
                <InfoCard icon={<Sun className="w-5 h-5 text-yellow-400" />} label="Day Length" value={result.day_length} sub={result.date_time.split('T')[0]} />
                <InfoCard icon={<Moon className="w-5 h-5 text-slate-300" />} label="Moon Phase" value={result.moon_phase} sub={`Phase angle: ${result.moon_phase_angle}°`} />
                <InfoCard icon={<Star className="w-5 h-5 text-indigo-400" />} label="Polaris Altitude" value={result.polaris_altitude > 0 ? `${result.polaris_altitude}°` : 'Not visible'} sub={result.true_north_hint} />
              </div>
            )}
          </div>
        )}

        {/* Stars Tab */}
        {activeView === 'stars' && (
          <div className="space-y-3 max-w-3xl">
            {stars.map(star => (
              <div key={star.name} className="p-4 rounded-xl bg-slate-800/30 border border-slate-700/30">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Star className="w-4 h-4 text-amber-400" />
                    <span className="text-sm font-semibold text-slate-200">{star.name}</span>
                    <span className="text-[10px] px-2 py-0.5 rounded-full bg-slate-700/50 text-slate-500">{star.constellation}</span>
                  </div>
                  <span className="text-xs text-slate-500">mag {star.magnitude}</span>
                </div>
                <p className="text-xs text-slate-500 mt-1">{star.description}</p>
                <p className="text-xs text-indigo-400 mt-2">🧭 {star.usage}</p>
              </div>
            ))}
          </div>
        )}

        {/* Techniques Tab */}
        {activeView === 'techniques' && (
          <div className="space-y-3 max-w-3xl">
            {techniques.map(tech => (
              <button
                key={tech.id}
                onClick={() => setActiveTechnique(tech)}
                className="w-full text-left p-4 rounded-xl bg-slate-800/30 border border-slate-700/30 hover:border-indigo-500/30 hover:bg-indigo-500/5 transition-all group"
              >
                <div className="flex items-center justify-between">
                  <span className="text-sm font-semibold text-slate-200 group-hover:text-indigo-300">{tech.name}</span>
                  <ChevronRight className="w-4 h-4 text-slate-600 group-hover:text-indigo-400" />
                </div>
                <p className="text-xs text-slate-500 mt-1">{tech.description}</p>
                <span className="inline-block mt-2 px-2 py-0.5 rounded-full text-[10px] bg-indigo-500/10 border border-indigo-500/20 text-indigo-400">
                  {hemisphereLabel[tech.hemisphere]}
                </span>
              </button>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

function InfoCard({ icon, label, value, sub }: { icon: React.ReactNode; label: string; value: string; sub: string }) {
  return (
    <div className="p-4 rounded-xl bg-slate-800/30 border border-slate-700/30">
      <div className="flex items-center gap-2 mb-2">
        {icon}
        <span className="text-[10px] font-semibold text-slate-500 uppercase tracking-wider">{label}</span>
      </div>
      <p className="text-lg font-semibold text-slate-100">{value}</p>
      <p className="text-xs text-slate-500 mt-0.5">{sub}</p>
    </div>
  )
}
