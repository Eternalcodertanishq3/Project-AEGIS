import { useState } from 'react'
import {
  LayoutDashboard,
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
  Cpu,
  ChevronLeft,
  ChevronRight,
  Shield,
  Zap,
} from 'lucide-react'
import { cn } from '@/lib/utils'

interface NavItem {
  id: string
  label: string
  icon: React.ElementType
  group: string
}

const navItems: NavItem[] = [
  { id: 'overview', label: 'Overview', icon: LayoutDashboard, group: 'Main' },
  { id: 'knowledge-library', label: 'Knowledge library', icon: BookOpen, group: 'Knowledge' },
  { id: 'offline-maps', label: 'Offline maps', icon: Map, group: 'Knowledge' },
  { id: 'notes', label: 'Notes', icon: FileText, group: 'Knowledge' },
  { id: 'data-tools', label: 'Data tools', icon: Database, group: 'Knowledge' },
  { id: 'ai-assistant', label: 'AI assistant', icon: Brain, group: 'AI' },
  { id: 'medical-triage', label: 'Medical triage', icon: Heart, group: 'Survival' },
  { id: 'plant-fungi-id', label: 'Plant/Fungi ID', icon: Leaf, group: 'Survival' },
  { id: 'skill-trees', label: 'Skill trees', icon: GitBranch, group: 'Survival' },
  { id: 'celestial-nav', label: 'Celestial navigation', icon: Star, group: 'Survival' },
  { id: 'mesh-messaging', label: 'Mesh messaging', icon: Radio, group: 'Comms' },
  { id: 'encrypted-p2p', label: 'Encrypted P2P', icon: Lock, group: 'Comms' },
  { id: 'sdr-monitor', label: 'SDR monitor', icon: Activity, group: 'Comms' },
  { id: 'local-peer-sync', label: 'Local peer sync', icon: RefreshCw, group: 'Comms' },
  { id: 'position-beacon', label: 'Position beacon', icon: Navigation, group: 'Comms' },
  { id: 'system', label: 'System', icon: Cpu, group: 'System' },
  { id: 'power', label: 'Power budget', icon: Zap, group: 'System' },
]

interface SidebarProps {
  activeItem: string
  onNavigate: (id: string) => void
}

export function Sidebar({ activeItem, onNavigate }: SidebarProps) {
  const [collapsed, setCollapsed] = useState(false)

  const groups = navItems.reduce<Record<string, NavItem[]>>((acc, item) => {
    if (!acc[item.group]) acc[item.group] = []
    acc[item.group].push(item)
    return acc
  }, {})

  return (
    <aside
      className={cn(
        'flex flex-col h-[calc(100vh-3.5rem)] border-r border-slate-800/50 bg-slate-900/60 backdrop-blur-xl transition-all duration-300 z-10 shadow-2xl shadow-black/50',
        collapsed ? 'w-16' : 'w-64'
      )}
    >
      {/* Collapse toggle */}
      <div className="flex items-center justify-end p-2 border-b border-slate-800/50 bg-slate-900/40">
        <button
          onClick={() => setCollapsed(!collapsed)}
          className="p-1.5 rounded-md text-slate-400 hover:text-slate-100 hover:bg-slate-800/80 transition-all hover:scale-105"
          aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
        >
          {collapsed ? (
            <ChevronRight className="h-4 w-4" />
          ) : (
            <ChevronLeft className="h-4 w-4" />
          )}
        </button>
      </div>

      {/* Navigation */}
      <nav className="flex-1 overflow-y-auto py-4 px-3 space-y-6">
        {Object.entries(groups).map(([group, items]) => (
          <div key={group}>
            {!collapsed && (
              <p className="px-2 mb-2 text-[10px] font-bold uppercase tracking-widest text-slate-500/80">
                {group}
              </p>
            )}
            <div className="space-y-1">
              {items.map((item) => {
                const Icon = item.icon
                const isActive = activeItem === item.id
                return (
                  <button
                    key={item.id}
                    onClick={() => onNavigate(item.id)}
                    className={cn(
                      'flex items-center w-full gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-200 group',
                      isActive
                        ? 'bg-gradient-to-r from-emerald-500/20 to-emerald-500/5 text-emerald-400 border border-emerald-500/20 shadow-[inset_0_1px_0_rgba(16,185,129,0.1)]'
                        : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50 border border-transparent',
                      collapsed && 'justify-center px-0'
                    )}
                    title={collapsed ? item.label : undefined}
                  >
                    <Icon className={cn('h-4 w-4 flex-shrink-0 transition-transform group-hover:scale-110 duration-300', isActive ? 'text-emerald-400' : 'text-slate-500 group-hover:text-slate-300')} />
                    {!collapsed && <span className="truncate">{item.label}</span>}
                  </button>
                )
              })}
            </div>
          </div>
        ))}
      </nav>

      {/* Bottom brand */}
      {!collapsed && (
        <div className="p-4 border-t border-slate-800/50 bg-slate-900/40">
          <div className="flex items-center gap-2 text-slate-500/80">
            <Shield className="h-4 w-4" />
            <span className="text-xs font-semibold tracking-wide">AEGIS v0.1.0</span>
          </div>
        </div>
      )}
    </aside>
  )
}
