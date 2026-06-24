import { useState, useEffect, useRef } from 'react'
import { Radio, Loader2, Send, Hash, Wifi, WifiOff } from 'lucide-react'

interface Message {
  id: string; channel: string; sender: string; content: string; timestamp: string; via: string
}
interface Channel {
  id: string; name: string; description: string; message_count: number; last_activity?: string
}
interface MeshStatus {
  lora_connected: boolean; lan_available: boolean; node_id: string
  active_channels: number; total_messages: number; transport_mode: string
}

export function MeshPage() {
  const [status, setStatus] = useState<MeshStatus | null>(null)
  const [channels, setChannels] = useState<Channel[]>([])
  const [activeChannel, setActiveChannel] = useState<string>('general')
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(true)
  const msgEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    Promise.all([
      fetch('/api/mesh/status').then(r => r.json()),
      fetch('/api/mesh/channels').then(r => r.json()),
    ]).then(([s, c]) => {
      setStatus(s); setChannels(c.channels || [])
    }).finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    if (activeChannel) {
      fetch(`/api/mesh/channels/${activeChannel}/messages?limit=100`)
        .then(r => r.json()).then(d => setMessages(d.messages || []))
    }
  }, [activeChannel])

  useEffect(() => { msgEndRef.current?.scrollIntoView({ behavior: 'smooth' }) }, [messages])

  const send = async () => {
    if (!input.trim()) return
    const res = await fetch(`/api/mesh/channels/${activeChannel}/send`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content: input }),
    })
    if (res.ok) {
      const msg = await res.json()
      setMessages(prev => [...prev, msg])
      setInput('')
    }
  }

  if (loading) return <div className="flex items-center justify-center h-full"><Loader2 className="w-8 h-8 text-blue-500 animate-spin" /></div>

  return (
    <div className="flex h-full">
      {/* Channel sidebar */}
      <div className="w-56 border-r border-slate-800/50 flex flex-col">
        <div className="p-4 border-b border-slate-800/50">
          <div className="flex items-center gap-2">
            <Radio className="w-4 h-4 text-blue-400" />
            <span className="text-sm font-semibold text-slate-200">Channels</span>
          </div>
          {status && (
            <div className="flex items-center gap-1.5 mt-2">
              {status.lan_available ? <Wifi className="w-3 h-3 text-emerald-400" /> : <WifiOff className="w-3 h-3 text-red-400" />}
              <span className="text-[10px] text-slate-500 font-mono">{status.node_id}</span>
            </div>
          )}
        </div>
        <div className="flex-1 overflow-y-auto p-2 space-y-1">
          {channels.map(ch => (
            <button key={ch.id} onClick={() => setActiveChannel(ch.id)}
              className={`w-full text-left px-3 py-2 rounded-lg text-sm transition-all ${
                activeChannel === ch.id ? 'bg-blue-500/20 text-blue-300 border border-blue-500/30' : 'text-slate-400 hover:bg-slate-800/50 hover:text-slate-200'
              }`}>
              <div className="flex items-center gap-2">
                <Hash className="w-3 h-3" />
                <span className="font-medium">{ch.name}</span>
              </div>
              {ch.message_count > 0 && <span className="text-[10px] text-slate-600 ml-5">{ch.message_count} msgs</span>}
            </button>
          ))}
        </div>
      </div>

      {/* Message area */}
      <div className="flex-1 flex flex-col">
        <div className="p-4 border-b border-slate-800/50">
          <div className="flex items-center gap-2">
            <Hash className="w-4 h-4 text-slate-500" />
            <span className="text-sm font-semibold text-slate-200">{channels.find(c => c.id === activeChannel)?.name || activeChannel}</span>
          </div>
          <p className="text-xs text-slate-500 mt-0.5">{channels.find(c => c.id === activeChannel)?.description}</p>
        </div>

        <div className="flex-1 overflow-y-auto p-4 space-y-3">
          {messages.length === 0 ? (
            <div className="text-center py-12">
              <Radio className="w-8 h-8 text-slate-700 mx-auto mb-2" />
              <p className="text-sm text-slate-600">No messages yet. Start the conversation.</p>
            </div>
          ) : messages.map(msg => (
            <div key={msg.id} className="flex gap-3">
              <div className="w-8 h-8 rounded-full bg-blue-500/20 border border-blue-500/30 flex items-center justify-center shrink-0">
                <span className="text-xs font-bold text-blue-400">{msg.sender.slice(-2).toUpperCase()}</span>
              </div>
              <div className="flex-1">
                <div className="flex items-center gap-2">
                  <span className="text-xs font-semibold text-blue-400 font-mono">{msg.sender}</span>
                  <span className="text-[10px] text-slate-600">{new Date(msg.timestamp).toLocaleTimeString()}</span>
                  <span className="text-[10px] px-1.5 py-0.5 rounded bg-slate-800 text-slate-500">{msg.via}</span>
                </div>
                <p className="text-sm text-slate-300 mt-0.5">{msg.content}</p>
              </div>
            </div>
          ))}
          <div ref={msgEndRef} />
        </div>

        <div className="p-4 border-t border-slate-800/50">
          <div className="flex gap-2">
            <input type="text" value={input} onChange={e => setInput(e.target.value)}
              onKeyDown={e => e.key === 'Enter' && send()} placeholder={`Message #${activeChannel}...`}
              className="flex-1 px-3 py-2 rounded-lg bg-slate-800/50 border border-slate-700/50 text-sm text-slate-200 placeholder-slate-500 focus:outline-none focus:border-blue-500/50" />
            <button onClick={send} disabled={!input.trim()}
              className="px-4 py-2 rounded-lg bg-blue-500 hover:bg-blue-600 text-white transition-colors disabled:opacity-50">
              <Send className="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
