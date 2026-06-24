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
          <div className="p-6 space-y-6 max-w-7xl">
            {/* Page title */}
            <div>
              <h1 className="text-xl font-semibold text-slate-100">Command center</h1>
              <p className="text-sm text-slate-400 mt-1">
                Monitor and manage all AEGIS modules from a single dashboard
              </p>
            </div>

            {/* System overview */}
            <SystemOverview profile={profile} />

            {/* Module grid */}
            <ModuleGrid modules={modules} onToggle={toggleModule} />
          </div>
        )
    }
  }

  return (
    <div className="flex flex-col h-screen bg-slate-950 text-slate-100">
      {/* Header */}
      <Header profile={profile} />

      {/* Body */}
      <div className="flex flex-1 overflow-hidden">
        {/* Sidebar */}
        <Sidebar activeItem={activeNav} onNavigate={setActiveNav} />

        {/* Main content */}
        <main className="flex-1 overflow-y-auto">
          {renderPage()}
        </main>
      </div>

      {/* Footer */}
      <Footer />
    </div>
  )
}

export default App
