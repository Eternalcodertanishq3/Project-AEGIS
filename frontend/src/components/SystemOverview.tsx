import {
  Cpu, HardDrive, Monitor, Battery,
  BatteryCharging, Clock, Zap, Activity
} from 'lucide-react'
import type { SystemProfile, HardwareTier } from '@/types'

interface SystemOverviewProps {
  profile: SystemProfile
}

const tierConfig: Record<HardwareTier, { color: string; glow: string; label: string; badgeClasses: string }> = {
  Minimum: {
    color: 'text-amber-400',
    glow: 'glow-amber',
    label: 'Limited modules available',
    badgeClasses: 'bg-amber-500/10 text-amber-400 border-amber-500/20 shadow-[0_0_10px_rgba(245,158,11,0.2)]',
  },
  Standard: {
    color: 'text-emerald-400',
    glow: 'glow-emerald',
    label: 'Standard modules available',
    badgeClasses: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20 shadow-[0_0_10px_rgba(16,185,129,0.2)]',
  },
  Optimal: {
    color: 'text-cyan-400',
    glow: 'glow-cyan',
    label: 'All modules available',
    badgeClasses: 'bg-cyan-500/10 text-cyan-400 border-cyan-500/20 shadow-[0_0_10px_rgba(6,182,212,0.2)]',
  },
}

export function SystemOverview({ profile }: SystemOverviewProps) {
  const tierKey = (profile.hardwareTier.charAt(0).toUpperCase() + profile.hardwareTier.slice(1)) as HardwareTier
  const tier = tierConfig[tierKey] || tierConfig.Standard

  return (
    <div className="glass-panel p-6 flex flex-col h-full bg-slate-900/40 relative overflow-hidden group">
      {/* Ambient background glow */}
      <div className="absolute -top-24 -right-24 w-48 h-48 bg-emerald-500/5 blur-[50px] rounded-full pointer-events-none transition-all group-hover:bg-emerald-500/10" />

      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-xl bg-slate-800/80 border border-slate-700/50 shadow-inner">
            <Activity className="h-5 w-5 text-emerald-400" />
          </div>
          <h2 className="text-base font-bold tracking-tight text-slate-100">Telemetry</h2>
        </div>
      </div>

      {/* Hardware tier */}
      <div className={`rounded-xl border border-slate-700/50 bg-slate-800/40 p-5 mb-6 backdrop-blur-md transition-all ${tier.glow}`}>
        <div className="flex items-center justify-between mb-2">
          <span className="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Hardware Tier</span>
          <div className={`text-xs font-bold px-3 py-1 rounded-full ${tier.badgeClasses}`}>
            {profile.hardwareTier}
          </div>
        </div>
        <p className="text-xs text-slate-400 font-medium">{tier.label}</p>
      </div>

      {/* Stats grid */}
      <div className="grid grid-cols-2 gap-4 flex-1">
        {/* CPU */}
        <div className="rounded-xl border border-slate-700/50 bg-slate-800/40 p-4 hover:bg-slate-800/60 transition-colors">
          <div className="flex items-center gap-2 mb-2">
            <Cpu className="h-4 w-4 text-slate-400" />
            <span className="text-[10px] font-bold text-slate-500 uppercase tracking-widest">CPU</span>
          </div>
          <p className="text-2xl font-bold text-slate-100">{profile.cpuCores}</p>
          <p className="text-[10px] font-medium text-slate-500 uppercase tracking-wider mt-1">Cores Active</p>
        </div>

        {/* RAM */}
        <div className="rounded-xl border border-slate-700/50 bg-slate-800/40 p-4 hover:bg-slate-800/60 transition-colors">
          <div className="flex items-center gap-2 mb-2">
            <HardDrive className="h-4 w-4 text-slate-400" />
            <span className="text-[10px] font-bold text-slate-500 uppercase tracking-widest">Memory</span>
          </div>
          <p className="text-2xl font-bold text-slate-100">{profile.ramUsed}</p>
          <div className="mt-2 w-full bg-slate-900 rounded-full h-1.5 shadow-inner overflow-hidden border border-slate-800">
            <div
              className="h-full bg-gradient-to-r from-emerald-500 to-cyan-400 rounded-full transition-all duration-1000 ease-out"
              style={{ width: `${Math.min(profile.ramPercent, 100)}%` }}
            />
          </div>
          <p className="text-[10px] font-medium text-slate-500 mt-1.5">{profile.ramPercent}% of {profile.ramTotal}</p>
        </div>

        {/* Power */}
        <div className="rounded-xl border border-slate-700/50 bg-slate-800/40 p-4 hover:bg-slate-800/60 transition-colors">
          <div className="flex items-center gap-2 mb-2">
            {profile.batteryCharging ? (
              <BatteryCharging className="h-4 w-4 text-emerald-400 animate-pulse" />
            ) : (
              <Battery className={`h-4 w-4 ${profile.batteryPercent < 20 ? 'text-red-400' : 'text-slate-400'}`} />
            )}
            <span className="text-[10px] font-bold text-slate-500 uppercase tracking-widest">Power</span>
          </div>
          <div className="flex items-baseline gap-1.5">
            <p className="text-2xl font-bold text-slate-100">{profile.batteryPercent}%</p>
            {profile.batteryCharging && <Zap className="h-4 w-4 text-emerald-400" />}
          </div>
          <p className="text-[10px] font-medium text-slate-500 uppercase tracking-wider mt-1">
            {profile.batteryCharging ? 'Charging' : 'Discharging'}
          </p>
        </div>

        {/* System */}
        <div className="rounded-xl border border-slate-700/50 bg-slate-800/40 p-4 hover:bg-slate-800/60 transition-colors">
          <div className="flex items-center gap-2 mb-2">
            <Monitor className="h-4 w-4 text-slate-400" />
            <span className="text-[10px] font-bold text-slate-500 uppercase tracking-widest">OS</span>
          </div>
          <p className="text-base font-bold text-slate-100 capitalize truncate">{profile.os}</p>
          <p className="text-[10px] font-medium text-slate-500 uppercase tracking-wider mt-1">{profile.arch}</p>
        </div>
      </div>

      {/* Footer */}
      <div className="mt-6 pt-4 border-t border-slate-800/50 flex items-center justify-between text-xs text-slate-500 font-mono">
        <div className="flex items-center gap-2">
          <Clock className="h-3.5 w-3.5 text-slate-400" />
          <span>{profile.uptime}</span>
        </div>
        <span className="px-2 py-0.5 rounded-md bg-slate-800/50 text-slate-400 border border-slate-700/50">{profile.hostname}</span>
      </div>
    </div>
  )
}
