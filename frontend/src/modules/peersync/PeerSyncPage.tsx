import { useState, useEffect } from 'react'
import { RefreshCw, Loader2, Plus, Trash2, Server, Globe } from 'lucide-react'

interface Peer { id: string; name: string; address: string; node_id: string; status: string; last_sync?: string; added_at: string }
interface SyncStatus { node_id: string; peer_count: number; online_peers: number; last_sync?: string; sync_enabled: boolean; listen_port: number }
interface ContentManifest { node_id: string; hostname: string; modules: { module_id: string; module_name: string; item_count: number; description: string }[]; total_items: number; last_updated: string }

export function PeerSyncPage() {
  const [status, setStatus] = useState<SyncStatus | null>(null)
  const [manifest, setManifest] = useState<ContentManifest | null>(null)
  const [peers, setPeers] = useState<Peer[]>([])
  const [showAdd, setShowAdd] = useState(false)
  const [newName, setNewName] = useState('')
  const [newAddr, setNewAddr] = useState('')
  const [loading, setLoading] = useState(true)

  const refresh = () => {
    Promise.all([
      fetch('/api/sync/status').then(r => r.json()),
      fetch('/api/sync/manifest').then(r => r.json()),
      fetch('/api/sync/peers').then(r => r.json()),
    ]).then(([s, m, p]) => { setStatus(s); setManifest(m); setPeers(p.peers || []) }).finally(() => setLoading(false))
  }

  useEffect(() => { refresh() }, [])

  const addPeer = async () => {
    if (!newName.trim() || !newAddr.trim()) return
    await fetch('/api/sync/peers', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ name: newName, address: newAddr }) })
    setNewName(''); setNewAddr(''); setShowAdd(false); refresh()
  }

  const removePeer = async (id: string) => {
    await fetch(`/api/sync/peers/${id}`, { method: 'DELETE' }); refresh()
  }

  if (loading) return <div className="flex items-center justify-center h-full"><Loader2 className="w-8 h-8 text-teal-500 animate-spin" /></div>

  return (
    <div className="flex flex-col h-full">
      <div className="p-6 pb-4 border-b border-slate-800/50">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-teal-500/20 border border-teal-500/30 flex items-center justify-center">
            <RefreshCw className="w-5 h-5 text-teal-400" />
          </div>
          <div>
            <h1 className="text-xl font-semibold text-slate-100">Local Peer Sync</h1>
            <p className="text-sm text-slate-400">Share content with nearby AEGIS nodes over LAN</p>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-6 space-y-6 max-w-3xl">
        {/* Status cards */}
        {status && (
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
            <StatusCard label="Node ID" value={status.node_id} />
            <StatusCard label="Peers" value={`${status.peer_count} known`} />
            <StatusCard label="Online" value={`${status.online_peers} peers`} />
            <StatusCard label="Listen Port" value={`:${status.listen_port}`} />
          </div>
        )}

        {/* Content Manifest */}
        {manifest && (
          <div>
            <h2 className="text-sm font-semibold text-slate-300 uppercase tracking-wider mb-3">📦 Your Content Manifest</h2>
            <div className="p-4 rounded-xl bg-slate-800/30 border border-slate-700/30">
              <div className="flex items-center justify-between mb-3">
                <span className="text-sm text-slate-300">{manifest.total_items} items available to share</span>
                <span className="text-[10px] text-slate-500">{new Date(manifest.last_updated).toLocaleString()}</span>
              </div>
              <div className="space-y-2">
                {manifest.modules.map(mod => (
                  <div key={mod.module_id} className="flex items-center justify-between p-2 rounded-lg bg-slate-800/30">
                    <div className="flex items-center gap-2">
                      <Server className="w-3 h-3 text-teal-400" />
                      <span className="text-xs text-slate-300">{mod.module_name}</span>
                    </div>
                    <span className="text-xs text-slate-500">{mod.item_count} items</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}

        {/* Peers */}
        <div>
          <div className="flex items-center justify-between mb-3">
            <h2 className="text-sm font-semibold text-slate-300 uppercase tracking-wider">🌐 Known Peers</h2>
            <button onClick={() => setShowAdd(!showAdd)}
              className="flex items-center gap-1 px-3 py-1.5 rounded-lg bg-teal-500 hover:bg-teal-600 text-white text-xs font-semibold transition-colors">
              <Plus className="w-3 h-3" /> Add Peer
            </button>
          </div>

          {showAdd && (
            <div className="p-4 mb-3 rounded-xl bg-slate-800/30 border border-teal-500/30 space-y-2">
              <input type="text" value={newName} onChange={e => setNewName(e.target.value)} placeholder="Peer name (e.g. 'Base Camp Node')"
                className="w-full px-3 py-2 rounded-lg bg-slate-800/50 border border-slate-700/50 text-sm text-slate-200 focus:outline-none focus:border-teal-500/50" />
              <input type="text" value={newAddr} onChange={e => setNewAddr(e.target.value)} placeholder="Address (e.g. 192.168.1.100:8080)"
                className="w-full px-3 py-2 rounded-lg bg-slate-800/50 border border-slate-700/50 text-sm text-slate-200 font-mono focus:outline-none focus:border-teal-500/50" />
              <button onClick={addPeer} className="px-4 py-2 rounded-lg bg-teal-500 hover:bg-teal-600 text-white text-sm font-semibold transition-colors">Add</button>
            </div>
          )}

          {peers.length === 0 ? (
            <div className="text-center py-8 rounded-xl bg-slate-800/20 border border-slate-700/20">
              <Globe className="w-8 h-8 text-slate-700 mx-auto mb-2" />
              <p className="text-sm text-slate-600">No peers configured yet</p>
              <p className="text-xs text-slate-700 mt-1">Add another AEGIS node's IP address to start syncing</p>
            </div>
          ) : (
            <div className="space-y-2">
              {peers.map(peer => (
                <div key={peer.id} className="flex items-center justify-between p-4 rounded-xl bg-slate-800/30 border border-slate-700/30 group">
                  <div>
                    <div className="flex items-center gap-2">
                      <div className={`w-2 h-2 rounded-full ${peer.status === 'online' ? 'bg-emerald-400' : 'bg-slate-600'}`} />
                      <span className="text-sm font-semibold text-slate-200">{peer.name}</span>
                    </div>
                    <span className="text-xs font-mono text-slate-500 ml-4">{peer.address}</span>
                  </div>
                  <button onClick={() => removePeer(peer.id)}
                    className="opacity-0 group-hover:opacity-100 p-2 text-red-400 hover:text-red-300 transition-opacity">
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

function StatusCard({ label, value }: { label: string; value: string }) {
  return (
    <div className="p-3 rounded-xl bg-slate-800/30 border border-slate-700/30">
      <p className="text-[10px] font-semibold text-slate-500 uppercase tracking-wider">{label}</p>
      <p className="text-sm font-semibold text-slate-200 mt-1 font-mono truncate">{value}</p>
    </div>
  )
}
