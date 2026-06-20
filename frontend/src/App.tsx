import { useState } from 'react'
import { Header } from '@/components/Header'
import { Sidebar } from '@/components/Sidebar'
import { SystemOverview } from '@/components/SystemOverview'
import { ModuleGrid } from '@/components/ModuleGrid'
import { Footer } from '@/components/Footer'
import { KnowledgePage } from '@/modules/knowledge/KnowledgePage'
import { MapsPage } from '@/modules/maps/MapsPage'
import { AIPage } from '@/modules/ai/AIPage'
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
