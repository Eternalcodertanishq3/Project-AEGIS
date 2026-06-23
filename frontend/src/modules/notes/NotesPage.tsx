import { useState, useEffect } from 'react'
import {
  FileText, Plus, Save, Trash2, Pin, PinOff,
  Loader2, Search, ArrowLeft
} from 'lucide-react'

interface Note {
  id: string
  title: string
  content: string
  tags: string
  created_at: string
  updated_at: string
  pinned: boolean
}

export function NotesPage() {
  const [notes, setNotes] = useState<Note[]>([])
  const [loading, setLoading] = useState(true)
  const [activeNote, setActiveNote] = useState<Note | null>(null)
  const [editing, setEditing] = useState(false)
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [tags, setTags] = useState('')
  const [saving, setSaving] = useState(false)
  const [search, setSearch] = useState('')

  const fetchNotes = async () => {
    try {
      const res = await fetch('/api/notes')
      const data = await res.json()
      setNotes(data.notes || [])
    } catch (e) {
      console.error('Failed to fetch notes:', e)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchNotes() }, [])

  const createNote = async () => {
    setSaving(true)
    try {
      const res = await fetch('/api/notes', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title: 'Untitled Note', content: '', tags: '' })
      })
      const note = await res.json()
      await fetchNotes()
      setActiveNote(note)
      setTitle(note.title)
      setContent(note.content)
      setTags(note.tags)
      setEditing(true)
    } catch (e) {
      console.error('Failed to create note:', e)
    } finally {
      setSaving(false)
    }
  }

  const saveNote = async () => {
    if (!activeNote) return
    setSaving(true)
    try {
      const res = await fetch(`/api/notes/${activeNote.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, content, tags })
      })
      const updated = await res.json()
      setActiveNote(updated)
      await fetchNotes()
    } catch (e) {
      console.error('Failed to save note:', e)
    } finally {
      setSaving(false)
    }
  }

  const deleteNote = async (id: string) => {
    try {
      await fetch(`/api/notes/${id}`, { method: 'DELETE' })
      if (activeNote?.id === id) {
        setActiveNote(null)
        setEditing(false)
      }
      await fetchNotes()
    } catch (e) {
      console.error('Failed to delete note:', e)
    }
  }

  const togglePin = async (id: string) => {
    try {
      await fetch(`/api/notes/${id}/pin`, { method: 'POST' })
      await fetchNotes()
    } catch (e) {
      console.error('Failed to toggle pin:', e)
    }
  }

  const openNote = (note: Note) => {
    setActiveNote(note)
    setTitle(note.title)
    setContent(note.content)
    setTags(note.tags)
    setEditing(true)
  }

  const filteredNotes = notes.filter(n =>
    n.title.toLowerCase().includes(search.toLowerCase()) ||
    n.content.toLowerCase().includes(search.toLowerCase()) ||
    n.tags.toLowerCase().includes(search.toLowerCase())
  )

  const formatDate = (dateStr: string) => {
    try {
      return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit'
      })
    } catch { return dateStr }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="w-8 h-8 text-emerald-500 animate-spin" />
      </div>
    )
  }

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="p-6 pb-4 border-b border-slate-800/50">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-lg bg-amber-500/20 border border-amber-500/30 flex items-center justify-center">
              <FileText className="w-5 h-5 text-amber-400" />
            </div>
            <div>
              <h1 className="text-xl font-semibold text-slate-100">Field Notes</h1>
              <p className="text-sm text-slate-400">{notes.length} notes · Stored locally in SQLite</p>
            </div>
          </div>
          <button
            onClick={createNote}
            disabled={saving}
            className="flex items-center gap-2 px-4 py-2 rounded-lg bg-amber-500/10 text-amber-400 border border-amber-500/30 hover:bg-amber-500/20 transition-colors text-sm font-medium"
          >
            <Plus className="w-4 h-4" /> New Note
          </button>
        </div>
      </div>

      <div className="flex-1 flex overflow-hidden">
        {/* Sidebar — Note list */}
        <div className="w-72 border-r border-slate-800/50 flex flex-col bg-slate-900/30">
          <div className="p-3">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500" />
              <input
                type="text"
                placeholder="Search notes..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="w-full pl-9 pr-3 py-2 rounded-lg bg-slate-800/50 border border-slate-700/50 text-sm text-slate-200 placeholder-slate-500 focus:outline-none focus:border-amber-500/50"
              />
            </div>
          </div>
          <div className="flex-1 overflow-y-auto px-2 pb-2 space-y-1">
            {filteredNotes.length === 0 && (
              <div className="text-center py-8">
                <FileText className="w-8 h-8 text-slate-600 mx-auto mb-2" />
                <p className="text-sm text-slate-500">
                  {search ? 'No notes found' : 'No notes yet'}
                </p>
              </div>
            )}
            {filteredNotes.map(note => (
              <button
                key={note.id}
                onClick={() => openNote(note)}
                className={`w-full text-left p-3 rounded-lg border transition-all ${
                  activeNote?.id === note.id
                    ? 'bg-amber-500/10 border-amber-500/30'
                    : 'bg-slate-800/20 border-slate-700/30 hover:border-slate-600/50'
                }`}
              >
                <div className="flex items-start justify-between">
                  <span className="text-sm font-medium text-slate-200 truncate flex-1">
                    {note.pinned && <Pin className="w-3 h-3 text-amber-400 inline mr-1" />}
                    {note.title || 'Untitled'}
                  </span>
                </div>
                <p className="text-xs text-slate-500 mt-1 line-clamp-2">
                  {note.content.substring(0, 80) || 'Empty note'}
                </p>
                <p className="text-[10px] text-slate-600 mt-1">{formatDate(note.updated_at)}</p>
              </button>
            ))}
          </div>
        </div>

        {/* Editor */}
        <div className="flex-1 flex flex-col">
          {!editing ? (
            <div className="flex-1 flex items-center justify-center">
              <div className="text-center">
                <FileText className="w-12 h-12 text-slate-700 mx-auto mb-3" />
                <p className="text-slate-500">Select a note or create a new one</p>
              </div>
            </div>
          ) : (
            <>
              {/* Toolbar */}
              <div className="flex items-center gap-2 p-3 border-b border-slate-800/50">
                <button
                  onClick={() => { setEditing(false); setActiveNote(null) }}
                  className="p-2 rounded-lg text-slate-400 hover:text-slate-200 hover:bg-slate-800 transition-colors"
                >
                  <ArrowLeft className="w-4 h-4" />
                </button>
                <input
                  type="text"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="Note title..."
                  className="flex-1 bg-transparent text-lg font-semibold text-slate-100 placeholder-slate-600 focus:outline-none"
                />
                <button
                  onClick={() => activeNote && togglePin(activeNote.id)}
                  className="p-2 rounded-lg text-slate-400 hover:text-amber-400 hover:bg-slate-800 transition-colors"
                  title={activeNote?.pinned ? 'Unpin' : 'Pin'}
                >
                  {activeNote?.pinned ? <PinOff className="w-4 h-4" /> : <Pin className="w-4 h-4" />}
                </button>
                <button
                  onClick={saveNote}
                  disabled={saving}
                  className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-emerald-500/10 text-emerald-400 border border-emerald-500/30 hover:bg-emerald-500/20 transition-colors text-sm font-medium"
                >
                  {saving ? <Loader2 className="w-4 h-4 animate-spin" /> : <Save className="w-4 h-4" />}
                  Save
                </button>
                <button
                  onClick={() => activeNote && deleteNote(activeNote.id)}
                  className="p-2 rounded-lg text-slate-400 hover:text-red-400 hover:bg-red-500/10 transition-colors"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>

              {/* Tags */}
              <div className="px-4 py-2 border-b border-slate-800/30">
                <input
                  type="text"
                  value={tags}
                  onChange={(e) => setTags(e.target.value)}
                  placeholder="Tags (comma-separated)..."
                  className="w-full bg-transparent text-xs text-slate-400 placeholder-slate-600 focus:outline-none"
                />
              </div>

              {/* Content */}
              <textarea
                value={content}
                onChange={(e) => setContent(e.target.value)}
                placeholder="Start writing... (Markdown supported)"
                className="flex-1 p-4 bg-transparent text-sm text-slate-200 placeholder-slate-600 resize-none focus:outline-none font-mono leading-relaxed"
              />
            </>
          )}
        </div>
      </div>
    </div>
  )
}
