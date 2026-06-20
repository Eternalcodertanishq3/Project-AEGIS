import { WifiOff, Shield } from 'lucide-react'

export function Footer() {
  return (
    <footer className="border-t border-slate-700/50 bg-slate-950/80 px-4 lg:px-6 py-2.5">
      <div className="flex items-center justify-between text-[11px] text-slate-500">
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-1.5">
            <Shield className="h-3 w-3 text-slate-600" />
            <span className="font-medium">AEGIS</span>
            <span>v0.1.0-alpha</span>
          </div>
          <span className="text-slate-700">|</span>
          <span>Phase 0 — Scaffold</span>
        </div>

        <div className="flex items-center gap-3">
          <div className="flex items-center gap-1.5">
            <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 shadow-[0_0_4px_rgba(16,185,129,0.6)]" />
            <span>Running on port 8080</span>
          </div>
          <span className="text-slate-700">|</span>
          <div className="flex items-center gap-1.5 text-amber-500/80">
            <WifiOff className="h-3 w-3" />
            <span>Offline mode</span>
          </div>
        </div>
      </div>
    </footer>
  )
}
