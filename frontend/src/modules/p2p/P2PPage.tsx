import { useState, useEffect } from 'react'
import { Lock, Loader2, UserPlus, Send, Key, Trash2, Shield } from 'lucide-react'

interface Contact { id: string; alias: string; public_key: string; trusted: boolean; last_seen?: string; added_at: string }
interface SecureMessage { id: string; contact_id: string; direction: string; content: string; encrypted: boolean; timestamp: string }
interface P2PStatus { key_generated: boolean; public_key: string; contact_count: number; message_count: number }

export function P2PPage() {
  const [status, setStatus] = useState<P2PStatus | null>(null)
  const [contacts, setContacts] = useState<Contact[]>([])
  const [activeContact, setActiveContact] = useState<Contact | null>(null)
  const [messages, setMessages] = useState<SecureMessage[]>([])
  const [input, setInput] = useState('')
  const [showAdd, setShowAdd] = useState(false)
  const [newAlias, setNewAlias] = useState('')
  const [newKey, setNewKey] = useState('')
  const [loading, setLoading] = useState(true)

  const refresh = () => {
    Promise.all([
      fetch('/api/p2p/status').then(r => r.json()),
      fetch('/api/p2p/contacts').then(r => r.json()),
    ]).then(([s, c]) => { setStatus(s); setContacts(c.contacts || []) }).finally(() => setLoading(false))
  }

  useEffect(() => { refresh() }, [])

  useEffect(() => {
    if (activeContact) {
      fetch(`/api/p2p/contacts/${activeContact.id}/messages`).then(r => r.json()).then(d => setMessages(d.messages || []))
    }
  }, [activeContact])

  const addContact = async () => {
    if (!newAlias.trim() || !newKey.trim()) return
    const res = await fetch('/api/p2p/contacts', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ alias: newAlias, public_key: newKey }),
    })
    if (res.ok) { setNewAlias(''); setNewKey(''); setShowAdd(false); refresh() }
  }

  const deleteContact = async (id: string) => {
    await fetch(`/api/p2p/contacts/${id}`, { method: 'DELETE' })
    setActiveContact(null); setMessages([]); refresh()
  }

  const sendMsg = async () => {
    if (!input.trim() || !activeContact) return
    const res = await fetch(`/api/p2p/contacts/${activeContact.id}/send`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content: input }),
    })
    if (res.ok) { const msg = await res.json(); setMessages(prev => [...prev, msg]); setInput('') }
  }

  if (loading) return <div className="flex items-center justify-center h-full"><Loader2 className="w-8 h-8 text-purple-500 animate-spin" /></div>

  return (
    <div className="flex h-full">
      {/* Contact sidebar */}
      <div className="w-60 border-r border-slate-800/50 flex flex-col">
        <div className="p-4 border-b border-slate-800/50">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Lock className="w-4 h-4 text-purple-400" />
              <span className="text-sm font-semibold text-slate-200">Encrypted P2P</span>
            </div>
            <button onClick={() => setShowAdd(!showAdd)} className="p-1.5 rounded-lg hover:bg-slate-800 text-slate-400 hover:text-purple-400 transition-colors">
              <UserPlus className="w-4 h-4" />
            </button>
          </div>
          {status && (
            <div className="flex items-center gap-1 mt-2">
              <Key className="w-3 h-3 text-emerald-400" />
              <span className="text-[10px] text-slate-500 font-mono truncate">{status.public_key.slice(0, 24)}...</span>
            </div>
          )}
        </div>

        {showAdd && (
          <div className="p-3 border-b border-slate-800/50 space-y-2 bg-slate-800/20">
            <input type="text" value={newAlias} onChange={e => setNewAlias(e.target.value)} placeholder="Contact alias"
              className="w-full px-2 py-1.5 rounded bg-slate-800/50 border border-slate-700/50 text-xs text-slate-200 focus:outline-none focus:border-purple-500/50" />
            <input type="text" value={newKey} onChange={e => setNewKey(e.target.value)} placeholder="Public key"
              className="w-full px-2 py-1.5 rounded bg-slate-800/50 border border-slate-700/50 text-xs text-slate-200 focus:outline-none focus:border-purple-500/50 font-mono" />
            <button onClick={addContact} className="w-full px-2 py-1.5 rounded bg-purple-500 hover:bg-purple-600 text-white text-xs font-semibold transition-colors">Add Contact</button>
          </div>
        )}

        <div className="flex-1 overflow-y-auto p-2 space-y-1">
          {contacts.length === 0 ? (
            <p className="text-xs text-slate-600 text-center py-4">No contacts yet</p>
          ) : contacts.map(c => (
            <button key={c.id} onClick={() => setActiveContact(c)}
              className={`w-full text-left px-3 py-2 rounded-lg text-sm transition-all group ${
                activeContact?.id === c.id ? 'bg-purple-500/20 text-purple-300 border border-purple-500/30' : 'text-slate-400 hover:bg-slate-800/50'
              }`}>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <Shield className="w-3 h-3 text-purple-400" />
                  <span className="font-medium">{c.alias}</span>
                </div>
                <button onClick={(e) => { e.stopPropagation(); deleteContact(c.id) }}
                  className="opacity-0 group-hover:opacity-100 p-1 text-red-400 hover:text-red-300 transition-opacity">
                  <Trash2 className="w-3 h-3" />
                </button>
              </div>
              <span className="text-[10px] text-slate-600 font-mono ml-5">{c.public_key.slice(0, 16)}...</span>
            </button>
          ))}
        </div>
      </div>

      {/* Chat area */}
      <div className="flex-1 flex flex-col">
        {!activeContact ? (
          <div className="flex-1 flex items-center justify-center">
            <div className="text-center">
              <Lock className="w-10 h-10 text-slate-700 mx-auto mb-3" />
              <p className="text-sm text-slate-500">Select a contact to start an encrypted conversation</p>
              <p className="text-xs text-slate-600 mt-1">All messages are encrypted with Ed25519 keys</p>
            </div>
          </div>
        ) : (
          <>
            <div className="p-4 border-b border-slate-800/50 flex items-center justify-between">
              <div>
                <span className="text-sm font-semibold text-slate-200">{activeContact.alias}</span>
                <div className="flex items-center gap-1 mt-0.5">
                  <Lock className="w-3 h-3 text-emerald-400" />
                  <span className="text-[10px] text-emerald-400">End-to-end encrypted</span>
                </div>
              </div>
            </div>
            <div className="flex-1 overflow-y-auto p-4 space-y-3">
              {messages.length === 0 ? (
                <p className="text-xs text-slate-600 text-center py-8">No messages yet. Send the first encrypted message.</p>
              ) : messages.map(msg => (
                <div key={msg.id} className={`flex ${msg.direction === 'sent' ? 'justify-end' : 'justify-start'}`}>
                  <div className={`max-w-xs px-3 py-2 rounded-xl ${
                    msg.direction === 'sent' ? 'bg-purple-500/20 border border-purple-500/30' : 'bg-slate-800/50 border border-slate-700/30'
                  }`}>
                    <p className="text-sm text-slate-200">{msg.content}</p>
                    <div className="flex items-center gap-1 mt-1">
                      {msg.encrypted && <Lock className="w-2.5 h-2.5 text-emerald-500" />}
                      <span className="text-[10px] text-slate-500">{new Date(msg.timestamp).toLocaleTimeString()}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
            <div className="p-4 border-t border-slate-800/50">
              <div className="flex gap-2">
                <input type="text" value={input} onChange={e => setInput(e.target.value)}
                  onKeyDown={e => e.key === 'Enter' && sendMsg()} placeholder="Encrypted message..."
                  className="flex-1 px-3 py-2 rounded-lg bg-slate-800/50 border border-slate-700/50 text-sm text-slate-200 placeholder-slate-500 focus:outline-none focus:border-purple-500/50" />
                <button onClick={sendMsg} disabled={!input.trim()}
                  className="px-4 py-2 rounded-lg bg-purple-500 hover:bg-purple-600 text-white transition-colors disabled:opacity-50">
                  <Send className="w-4 h-4" />
                </button>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  )
}
