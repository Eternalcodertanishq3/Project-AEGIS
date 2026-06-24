import { useState } from 'react'
import { Header } from '@/components/Header'
import { Sidebar } from '@/components/Sidebar'
import { SystemOverview } from '@/components/SystemOverview'
import { ModuleGrid } from '@/components/ModuleGrid'
import { Footer } from '@/components/Footer'
import { KnowledgePage } from '@/modules/knowledge/KnowledgePage'
import { MapsPage } from '@/modules/maps/MapsPage'
import { AIPage } from '@/modules/ai/AIPage'
import { NotesPage } from '@/modules/notes/NotesPage'
import { MedicalPage } from '@/modules/medical/MedicalPage'
import { DataToolsPage } from '@/modules/datatools/DataToolsPage'
import { SkillTreesPage } from '@/modules/skilltrees/SkillTreesPage'
import { CelestialPage } from '@/modules/celestial/CelestialPage'
import { PlantIDPage } from '@/modules/plantid/PlantIDPage'
import { MeshPage } from '@/modules/mesh/MeshPage'
import { P2PPage } from '@/modules/p2p/P2PPage'
import { SDRPage } from '@/modules/sdr/SDRPage'
import { PeerSyncPage } from '@/modules/peersync/PeerSyncPage'
import { BeaconPage } from '@/modules/beacon/BeaconPage'
import { useSystemProfile } from '@/hooks/useSystemProfile'
import { useModules } from '@/hooks/useModules'

function App() {
  const [activeNav, setActiveNav] = useState('overview')
  const { profile } = useSystemProfile()
  const { modules, toggleModule } = useModules()

  const renderPage = () => {
    switch (activeNav) {
      case 'knowledge-library':
        return <KnowledgePage />
      case 'offline-maps':
        return <MapsPage />
      case 'ai-assistant':
        return <AIPage />
      case 'notes':
        return <NotesPage />
      case 'medical-triage':
        return <MedicalPage />
      case 'data-tools':
        return <DataToolsPage />
      case 'skill-trees':
        return <SkillTreesPage />
      case 'celestial-nav':
        return <CelestialPage />
      case 'plant-fungi-id':
        return <PlantIDPage />
      case 'mesh-messaging':
        return <MeshPage />
      case 'encrypted-p2p':
        return <P2PPage />
      case 'sdr-monitor':
        return <SDRPage />
      case 'local-peer-sync':
        return <PeerSyncPage />
      case 'position-beacon':
        return <BeaconPage />
      case 'overview':
      default:
        return (
          <div className="p-6 md:p-8 space-y-8 max-w-7xl mx-auto">
            {/* Page title */}
            <div className="flex flex-col gap-2">
              <h1 className="text-3xl font-bold tracking-tight text-transparent bg-clip-text bg-gradient-to-r from-slate-100 to-slate-400">
                Command Center
              </h1>
              <p className="text-sm font-medium text-slate-500 max-w-2xl">
                Monitor and manage all AEGIS modules from a single, unified dashboard.
              </p>
            </div>

            {/* System overview */}
            <SystemOverview profile={profile} />

            {/* Module grid */}
            <div className="pt-4">
              <ModuleGrid modules={modules} onToggle={toggleModule} />
            </div>
          </div>
        )
    }
  }

  return (
    <div className="flex flex-col h-screen bg-[#030712] text-slate-100 selection:bg-emerald-500/30">
      {/* Header */}
      <Header profile={profile} />

      {/* Body */}
      <div className="flex flex-1 overflow-hidden relative">
        {/* Sidebar */}
        <Sidebar activeItem={activeNav} onNavigate={setActiveNav} />

        {/* Main content */}
        <main className="flex flex-col flex-1 h-full overflow-y-auto relative z-0 scroll-smooth bg-[#030712]">
          {renderPage()}
        </main>
      </div>

      {/* Footer */}
      <Footer />
    </div>
  )
}

export default App
