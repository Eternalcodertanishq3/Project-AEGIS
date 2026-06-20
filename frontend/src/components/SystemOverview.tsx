import {
  Cpu,
  HardDrive,
  Monitor,
  Battery,
  BatteryCharging,
  Clock,
  Server,
  Zap,
} from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import type { SystemProfile, HardwareTier } from '@/types'

interface SystemOverviewProps {
  profile: SystemProfile
}

const tierConfig: Record<HardwareTier, { color: string; glow: string; label: string }> = {
  Minimum: {
    color: 'text-amber-400',
    glow: 'glow-amber',
    label: 'Limited hardware — some modules unavailable',
  },
  Standard: {
    color: 'text-emerald-400',
    glow: 'glow-emerald',
    label: 'Standard hardware — most modules available',
  },
  Optimal: {
    color: 'text-cyan-400',
    glow: '',
    label: 'Full hardware — all modules available',
  },
}

export function SystemOverview({ profile }: SystemOverviewProps) {
  const tier = tierConfig[profile.hardwareTier]

  return (
    <div className="glass-panel p-5">
      <div className="flex items-center gap-2 mb-4">
        <Server className="h-4 w-4 text-slate-400" />
        <h2 className="text-sm font-semibold text-slate-200">System overview</h2>
      </div>

      {/* Hardware tier */}
      <div className={`rounded-lg border border-slate-700/50 bg-slate-800/30 p-4 mb-4 ${tier.glow}`}>
        <div className="flex items-center justify-between mb-1">
          <span className="text-xs font-medium text-slate-400 uppercase tracking-wide">Hardware tier</span>
          <Badge className={`text-xs ${
            profile.hardwareTier === 'Minimum'
              ? 'bg-amber-600/20 text-amber-400 border-amber-500/30'
              : profile.hardwareTier === 'Standard'
              ? 'bg-emerald-600/20 text-emerald-400 border-emerald-500/30'
              : 'bg-cyan-600/20 text-cyan-400 border-cyan-500/30'
          }`}>
            {profile.hardwareTier}
          </Badge>
        </div>
        <p className="text-xs text-slate-500 mt-1">{tier.label}</p>
      </div>

      {/* Stats grid */}
      <div className="grid grid-cols-2 gap-3">
        {/* CPU */}
        <div className="rounded-lg border border-slate-700/30 bg-slate-800/20 p-3">
          <div className="flex items-center gap-2 mb-1.5">
            <Cpu className="h-3.5 w-3.5 text-slate-500" />
            <span className="text-[11px] font-medium text-slate-500 uppercase">CPU</span>
          </div>
          <p className="text-lg font-semibold text-slate-100">{profile.cpuCores}</p>
          <p className="text-[11px] text-slate-500">cores</p>
        </div>

        {/* RAM */}
        <div className="rounded-lg border border-slate-700/30 bg-slate-800/20 p-3">
          <div className="flex items-center gap-2 mb-1.5">
            <HardDrive className="h-3.5 w-3.5 text-slate-500" />
            <span className="text-[11px] font-medium text-slate-500 uppercase">RAM</span>
          </div>
          <p className="text-lg font-semibold text-slate-100">{profile.ramTotal}</p>
          <div className="mt-1.5 w-full bg-slate-700/50 rounded-full h-1.5">
            <div
              className="bg-emerald-500 h-1.5 rounded-full transition-all duration-500"
              style={{ width: `${Math.min(profile.ramPercent, 100)}%` }}
            />
          </div>
          <p className="text-[11px] text-slate-500 mt-1">{profile.ramUsed} used ({profile.ramPercent}%)</p>
        </div>

        {/* Power */}
        <div className="rounded-lg border border-slate-700/30 bg-slate-800/20 p-3">
          <div className="flex items-center gap-2 mb-1.5">
            {profile.batteryCharging ? (
              <BatteryCharging className="h-3.5 w-3.5 text-emerald-500" />
            ) : (
              <Battery className={`h-3.5 w-3.5 ${profile.batteryPercent < 20 ? 'text-red-500' : 'text-slate-500'}`} />
            )}
            <span className="text-[11px] font-medium text-slate-500 uppercase">Power</span>
          </div>
          <div className="flex items-baseline gap-1">
            <p className="text-lg font-semibold text-slate-100">{profile.batteryPercent}%</p>
            {profile.batteryCharging && (
              <Zap className="h-3 w-3 text-emerald-400" />
            )}
          </div>
          <p className="text-[11px] text-slate-500">{profile.batteryCharging ? 'Charging' : 'On battery'}</p>
        </div>

        {/* System */}
        <div className="rounded-lg border border-slate-700/30 bg-slate-800/20 p-3">
          <div className="flex items-center gap-2 mb-1.5">
            <Monitor className="h-3.5 w-3.5 text-slate-500" />
            <span className="text-[11px] font-medium text-slate-500 uppercase">System</span>
          </div>
          <p className="text-sm font-semibold text-slate-100">{profile.os}</p>
          <p className="text-[11px] text-slate-500">{profile.arch}</p>
        </div>
      </div>

      {/* Uptime & hostname */}
      <div className="mt-3 flex items-center justify-between text-[11px] text-slate-500 px-1">
        <div className="flex items-center gap-1.5">
          <Clock className="h-3 w-3" />
          <span>Uptime: {profile.uptime}</span>
        </div>
        <span className="font-mono">{profile.hostname}</span>
      </div>
    </div>
  )
}
