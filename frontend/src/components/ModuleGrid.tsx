import { ModuleCard } from './ModuleCard'
import type { Module, ModuleDomain } from '@/types'

interface ModuleGridProps {
  modules: Module[]
  onToggle: (id: string) => void
}

const domainOrder: ModuleDomain[] = ['Knowledge', 'AI', 'Survival', 'Comms', 'System']

const domainDescriptions: Record<ModuleDomain, string> = {
  Knowledge: 'Encyclopedias, maps, notes, and reference tools',
  AI: 'Local AI inference and document analysis',
  Survival: 'Medical, botanical, navigation, and skill tools',
  Comms: 'Mesh radio, encrypted messaging, and peer sync',
  System: 'Power management and hardware profiling',
}

export function ModuleGrid({ modules, onToggle }: ModuleGridProps) {
  const grouped = domainOrder
    .map((domain) => ({
      domain,
      description: domainDescriptions[domain],
      modules: modules.filter((m) => m.domain === domain),
    }))
    .filter((g) => g.modules.length > 0)

  const activeCount = modules.filter((m) => m.status === 'active').length

  return (
    <div className="space-y-6">
      {/* Section header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-sm font-semibold text-slate-200">Modules</h2>
          <p className="text-xs text-slate-500 mt-0.5">
            {activeCount} of {modules.length} active
          </p>
        </div>
      </div>

      {/* Domain groups */}
      {grouped.map(({ domain, description, modules: domainModules }) => (
        <div key={domain}>
          <div className="mb-3">
            <h3 className="text-xs font-semibold text-slate-400 uppercase tracking-wider">{domain}</h3>
            <p className="text-[11px] text-slate-600 mt-0.5">{description}</p>
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-3">
            {domainModules.map((module) => (
              <ModuleCard key={module.id} module={module} onToggle={onToggle} />
            ))}
          </div>
        </div>
      ))}
    </div>
  )
}
