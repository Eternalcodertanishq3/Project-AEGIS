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
        'glass-panel p-4 transition-all duration-200 group',
        'hover:border-slate-600/60 hover:bg-slate-800/60',
        module.enabled && 'border-emerald-500/20 glow-emerald',
        isUnavailable && 'opacity-50'
      )}
    >
      {/* Top row: icon + toggle */}
      <div className="flex items-start justify-between mb-3">
        <div
          className={cn(
            'flex items-center justify-center w-9 h-9 rounded-lg transition-colors',
            module.enabled
              ? 'bg-emerald-600/15 border border-emerald-500/25'
              : 'bg-slate-800/60 border border-slate-700/40'
          )}
        >
          <Icon
            className={cn(
              'h-4.5 w-4.5 transition-colors',
              module.enabled ? 'text-emerald-400' : 'text-slate-500'
            )}
          />
        </div>
        <Switch
          checked={module.enabled}
          onCheckedChange={() => onToggle(module.id)}
          disabled={isUnavailable}
          className="data-[state=checked]:bg-emerald-600"
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
