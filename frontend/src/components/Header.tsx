import {
  Shield,
  Wifi,
  WifiOff,
  Monitor,
} from 'lucide-react'
import type { SystemProfile, HardwareTier } from '@/types'

interface HeaderProps {
  profile: SystemProfile
}

const tierColors: Record<HardwareTier, string> = {
  Minimum: 'bg-amber-600/20 text-amber-400 border-amber-500/30',
  Standard: 'bg-emerald-600/20 text-emerald-400 border-emerald-500/30',
  Optimal: 'bg-cyan-600/20 text-cyan-400 border-cyan-500/30',
}

export function Header({ profile }: HeaderProps) {
  return (
    <header className="sticky top-0 z-50 w-full border-b border-slate-800/50 bg-slate-900/60 backdrop-blur-xl shadow-lg shadow-black/20">
      <div className="flex h-14 items-center justify-between px-4 lg:px-6">
        {/* Left — Branding */}
        <div className="flex items-center gap-3">
          <div className="flex items-center justify-center w-8 h-8 rounded-lg bg-gradient-to-br from-emerald-500/20 to-cyan-500/10 border border-emerald-500/30 shadow-[0_0_15px_rgba(16,185,129,0.15)]">
            <Shield className="h-4.5 w-4.5 text-emerald-400" />
          </div>
          <div className="flex items-center gap-2">
            <span className="text-base font-bold tracking-widest uppercase text-slate-100">
              AEGIS
            </span>
            <span className="hidden sm:inline text-slate-600 text-sm font-normal">
              /
            </span>
            <span className="hidden sm:inline text-xs text-slate-400 font-bold uppercase tracking-wider">
              Command center
            </span>
          </div>
        </div>

        {/* Right — Status badges */}
        <div className="flex items-center gap-3">
          {/* Online/Offline */}
          <div className={`flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold uppercase tracking-wider border ${
            profile.isOnline 
              ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' 
              : 'bg-slate-800/50 text-slate-400 border-slate-700/50'
          }`}>
            {profile.isOnline ? (
              <Wifi className="h-3 w-3" />
            ) : (
              <WifiOff className="h-3 w-3" />
            )}
            {profile.isOnline ? 'Online' : 'Offline'}
          </div>

          {/* Hardware Tier */}
          <div className={`hidden sm:flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold uppercase tracking-wider border ${tierColors[profile.hardwareTier]}`}>
            <Monitor className="h-3 w-3" />
            {profile.hardwareTier}
          </div>

          {/* OS */}
          <div className="hidden md:flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold uppercase tracking-wider bg-slate-800/50 text-slate-400 border border-slate-700/50">
            {profile.os} {profile.arch}
          </div>
        </div>
      </div>
    </header>
  )
}
