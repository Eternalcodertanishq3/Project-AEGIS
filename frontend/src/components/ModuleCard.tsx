import {
  BookOpen,
  Map,
  Brain,
  FileText,
  Database,
  Heart,
  Leaf,
  GitBranch,
  Star,
  Radio,
  Lock,
  Activity,
  RefreshCw,
  Navigation,
} from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Switch } from '@/components/ui/switch'
import type { Module, ModuleDomain, ModuleStatus } from '@/types'
import { cn } from '@/lib/utils'

const iconMap: Record<string, React.ElementType> = {
  BookOpen,
  Map,
  Brain,
  FileText,
  Database,
  Heart,
  Leaf,
  GitBranch,
  Star,
  Radio,
  Lock,
  Activity,
  RefreshCw,
  Navigation,
}

const domainColors: Record<ModuleDomain, string> = {
  Knowledge: 'bg-blue-600/15 text-blue-400 border-blue-500/25',
  Survival: 'bg-amber-600/15 text-amber-400 border-amber-500/25',
  Comms: 'bg-purple-600/15 text-purple-400 border-purple-500/25',
  AI: 'bg-cyan-600/15 text-cyan-400 border-cyan-500/25',
  System: 'bg-slate-600/15 text-slate-400 border-slate-500/25',
}

const statusConfig: Record<ModuleStatus, { dot: string; label: string }> = {
  active: { dot: 'status-dot-active', label: 'Active' },
  inactive: { dot: 'status-dot-inactive', label: 'Inactive' },
  unavailable: { dot: 'status-dot-critical', label: 'Unavailable' },
}

interface ModuleCardProps {
  module: Module
  onToggle: (id: string) => void
}

export function ModuleCard({ module, onToggle }: ModuleCardProps) {
  const Icon = iconMap[module.icon] || BookOpen
  const status = statusConfig[module.status]
  const isUnavailable = module.status === 'unavailable'

  return (
    <div
      className={cn(
        'glass-panel p-4 transition-all duration-200 group flex flex-col gap-4 text-left border',
        'hover:border-emerald-500/30 hover:bg-slate-800/80',
        module.enabled ? 'border-emerald-500/40 bg-emerald-950/10' : 'border-slate-700/50',
        isUnavailable && 'opacity-50 cursor-not-allowed'
      )}
    >
      <div className="flex items-start justify-between gap-4 w-full">
        <div
          className={cn(
            'flex items-center justify-center w-10 h-10 rounded-xl transition-all duration-300',
            module.enabled
              ? 'bg-emerald-500/20 text-emerald-400 shadow-[0_0_15px_rgba(16,185,129,0.15)]'
              : 'bg-slate-800 text-slate-500'
          )}
        >
          <Icon className="h-5 w-5" />
        </div>
        <Switch
          checked={module.enabled}
          onCheckedChange={() => onToggle(module.id)}
          disabled={isUnavailable}
          className="data-[state=checked]:bg-emerald-500"
        />
      </div>

      {/* Module name */}
      <h3 className="text-sm font-semibold text-slate-100 mb-1">{module.name}</h3>

      {/* Description */}
      <p className="text-xs text-slate-400 mb-3 leading-relaxed">{module.description}</p>

      {/* Bottom: domain badge + status */}
      <div className="flex items-center justify-between">
        <Badge
          variant="outline"
          className={cn('text-[10px] border', domainColors[module.domain])}
        >
          {module.domain}
        </Badge>
        <div className="flex items-center gap-1.5">
          <span className={cn('status-dot', status.dot)} />
          <span className="text-[11px] text-slate-500">{status.label}</span>
        </div>
      </div>
    </div>
  )
}
