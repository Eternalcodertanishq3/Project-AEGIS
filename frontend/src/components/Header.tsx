import {
  Shield,
  Wifi,
  WifiOff,
  Monitor,
} from 'lucide-react'
import { Badge } from '@/components/ui/badge'
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
    <header className="sticky top-0 z-50 w-full border-b border-slate-700/50 bg-slate-950/80 backdrop-blur-md">
      <div className="flex h-14 items-center justify-between px-4 lg:px-6">
        {/* Left — Branding */}
        <div className="flex items-center gap-3">
          <div className="flex items-center justify-center w-8 h-8 rounded-md bg-emerald-600/20 border border-emerald-500/30">
            <Shield className="h-4.5 w-4.5 text-emerald-400" />
          </div>
          <div className="flex items-center gap-2">
            <span className="text-base font-semibold tracking-wide text-slate-100">
              AEGIS
            </span>
            <span className="hidden sm:inline text-slate-500 text-sm font-normal">
              /
            </span>
            <span className="hidden sm:inline text-sm text-slate-400 font-medium">
              Command center
            </span>
          </div>
        </div>

        {/* Right — Status badges */}
        <div className="flex items-center gap-2">
          {/* Online/Offline */}
          <Badge
            variant={profile.isOnline ? 'default' : 'secondary'}
            className="gap-1.5 text-xs"
          >
            {profile.isOnline ? (
              <Wifi className="h-3 w-3" />
            ) : (
              <WifiOff className="h-3 w-3" />
            )}
            {profile.isOnline ? 'Online' : 'Offline'}
          </Badge>

          {/* Hardware Tier */}
          <Badge className={`gap-1.5 text-xs border ${tierColors[profile.hardwareTier]}`}>
            <Monitor className="h-3 w-3" />
            {profile.hardwareTier}
          </Badge>

          {/* OS */}
          <Badge variant="outline" className="hidden md:inline-flex gap-1 text-xs">
            {profile.os} {profile.arch}
          </Badge>
        </div>
      </div>
    </header>
  )
}
