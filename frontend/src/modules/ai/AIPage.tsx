import { useState, useEffect, useRef } from 'react'
import {
  Bot, Send, Square, Loader2, Brain, HardDrive,
  AlertTriangle, CheckCircle, Cpu, Sparkles
} from 'lucide-react'

interface Model {
  id: string
  name: string
  path: string
  size_bytes: number
  size_human: string
}

interface AIStatus {
  module: string
  running: boolean
  port: number
  active_model: string
  models_count: number
  models: Model[]
  sidecar: string
  sidecar_note?: string
}

interface ChatMessage {
  role: 'user' | 'assistant' | 'system'
  content: string
  timestamp?: number
}

export function AIPage() {
  const [status, setStatus] = useState<AIStatus | null>(null)
  const [loading, setLoading] = useState(true)
  const [starting, setStarting] = useState(false)
  const [stopping, setStopping] = useState(false)
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [input, setInput] = useState('')
  const [streaming, setStreaming] = useState(false)
  const chatEndRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLTextAreaElement>(null)

  const fetchStatus = async () => {
    try {
      const res = await fetch('/api/ai/status')
      const data = await res.json()
      setStatus(data)
    } catch (e) {
      console.error('Failed to fetch AI status:', e)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchStatus()
    const interval = setInterval(fetchStatus, 5000)
    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    chatEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const startModel = async (modelId: string) => {
    setStarting(true)
    try {
      await fetch('/api/ai/start', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ model_id: modelId })
      })
      // Wait a bit for the server to start
      await new Promise(r => setTimeout(r, 2000))
      await fetchStatus()
    } catch (e) {
      console.error('Failed to start model:', e)
    } finally {
      setStarting(false)
    }
  }

  const stopModel = async () => {
    setStopping(true)
    try {
      await fetch('/api/ai/stop', { method: 'POST' })
      await new Promise(r => setTimeout(r, 1000))
      await fetchStatus()
      setMessages([])
    } catch (e) {
      console.error('Failed to stop model:', e)
    } finally {
      setStopping(false)
    }
  }

  const sendMessage = async () => {
    if (!input.trim() || streaming) return

    const userMsg: ChatMessage = { role: 'user', content: input.trim(), timestamp: Date.now() }
    const newMessages = [...messages, userMsg]
    setMessages(newMessages)
    setInput('')
    setStreaming(true)

    // Add placeholder for assistant response
    const assistantMsg: ChatMessage = { role: 'assistant', content: '', timestamp: Date.now() }
    setMessages([...newMessages, assistantMsg])

    try {
      const res = await fetch('/api/ai/chat', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          messages: newMessages.map(m => ({ role: m.role, content: m.content })),
          stream: true
        })
      })

      if (!res.ok) {
        const err = await res.json()
        setMessages(prev => {
          const updated = [...prev]
          updated[updated.length - 1] = {
            ...updated[updated.length - 1],
            content: `⚠️ Error: ${err.error || 'Failed to get response'}`
          }
          return updated
        })
        setStreaming(false)
        return
      }

      const reader = res.body?.getReader()
      const decoder = new TextDecoder()
      let fullContent = ''

      if (reader) {
        while (true) {
          const { done, value } = await reader.read()
          if (done) break

          const chunk = decoder.decode(value, { stream: true })
          
          // Parse SSE data lines
          const lines = chunk.split('\n')
          for (const line of lines) {
            if (line.startsWith('data: ')) {
              const data = line.slice(6)
              if (data === '[DONE]') continue
              try {
                const parsed = JSON.parse(data)
                const delta = parsed.choices?.[0]?.delta?.content || ''
                fullContent += delta
                setMessages(prev => {
                  const updated = [...prev]
                  updated[updated.length - 1] = {
                    ...updated[updated.length - 1],
                    content: fullContent
                  }
                  return updated
                })
              } catch {
                // Not valid JSON, might be raw text
                fullContent += data
                setMessages(prev => {
                  const updated = [...prev]
                  updated[updated.length - 1] = {
                    ...updated[updated.length - 1],
                    content: fullContent
                  }
                  return updated
                })
              }
            }
          }
        }
      }
    } catch (e) {
      console.error('Chat error:', e)
      setMessages(prev => {
        const updated = [...prev]
        updated[updated.length - 1] = {
          ...updated[updated.length - 1],
          content: '⚠️ Connection failed. Is the AI model loaded?'
        }
        return updated
      })
    } finally {
      setStreaming(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      sendMessage()
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="w-8 h-8 text-emerald-500 animate-spin" />
      </div>
    )
  }

  const hasSidecar = status?.sidecar && status.sidecar !== 'not found'
  const hasModels = (status?.models_count || 0) > 0
  const isRunning = status?.running || false

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="p-6 pb-4 border-b border-slate-800/50">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-lg bg-violet-500/20 border border-violet-500/30 flex items-center justify-center">
              <Bot className="w-5 h-5 text-violet-400" />
            </div>
            <div>
              <h1 className="text-xl font-semibold text-slate-100">AI Assistant</h1>
              <p className="text-sm text-slate-400">
                Offline AI powered by local language models
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            {isRunning ? (
              <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-emerald-500/10 border border-emerald-500/30">
                <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
                <span className="text-xs text-emerald-400 font-medium">Model Active</span>
              </div>
            ) : (
              <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-slate-700/50 border border-slate-600/30">
                <div className="w-2 h-2 rounded-full bg-slate-500" />
                <span className="text-xs text-slate-400 font-medium">Idle</span>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Main content */}
      <div className="flex-1 flex overflow-hidden">
        {/* Sidebar - Model panel */}
        <div className="w-72 border-r border-slate-800/50 flex flex-col bg-slate-900/30">
          <div className="p-4">
            <h2 className="text-sm font-semibold text-slate-300 uppercase tracking-wider mb-3">
              Models
            </h2>

            {/* Sidecar status */}
            <div className={`p-3 rounded-lg border mb-3 ${
              hasSidecar 
                ? 'bg-emerald-500/5 border-emerald-500/20' 
                : 'bg-amber-500/5 border-amber-500/20'
            }`}>
              <div className="flex items-center gap-2 mb-1">
                {hasSidecar ? (
                  <CheckCircle className="w-4 h-4 text-emerald-400" />
                ) : (
                  <AlertTriangle className="w-4 h-4 text-amber-400" />
                )}
                <span className={`text-xs font-medium ${hasSidecar ? 'text-emerald-400' : 'text-amber-400'}`}>
                  {hasSidecar ? 'Engine Ready' : 'Engine Missing'}
                </span>
              </div>
              {!hasSidecar && (
                <p className="text-xs text-slate-500 mt-1">
                  Place <code className="text-amber-400/80 bg-slate-800 px-1 rounded">llama-server</code> in<br />
                  <code className="text-amber-400/80 bg-slate-800 px-1 rounded text-[10px]">sidecars/llama/windows/</code>
                </p>
              )}
            </div>

            {/* Model list */}
            {hasModels ? (
              <div className="space-y-2">
                {status?.models.map(model => (
                  <div
                    key={model.id}
                    className={`p-3 rounded-lg border transition-all ${
                      status?.active_model === model.id
                        ? 'bg-violet-500/10 border-violet-500/30'
                        : 'bg-slate-800/30 border-slate-700/30 hover:border-slate-600/50'
                    }`}
                  >
                    <div className="flex items-start justify-between">
                      <div className="flex items-center gap-2 min-w-0">
                        <Brain className="w-4 h-4 text-violet-400 shrink-0" />
                        <span className="text-sm text-slate-200 truncate" title={model.name}>
                          {model.id}
                        </span>
                      </div>
                    </div>
                    <div className="flex items-center gap-2 mt-2">
                      <HardDrive className="w-3 h-3 text-slate-500" />
                      <span className="text-xs text-slate-500">{model.size_human}</span>
                    </div>
                    <div className="mt-2">
                      {status?.active_model === model.id ? (
                        <button
                          onClick={stopModel}
                          disabled={stopping}
                          className="w-full px-3 py-1.5 rounded text-xs font-medium bg-red-500/10 text-red-400 border border-red-500/30 hover:bg-red-500/20 transition-colors disabled:opacity-50"
                        >
                          {stopping ? (
                            <span className="flex items-center gap-1 justify-center">
                              <Loader2 className="w-3 h-3 animate-spin" /> Stopping...
                            </span>
                          ) : (
                            <span className="flex items-center gap-1 justify-center">
                              <Square className="w-3 h-3" /> Unload Model
                            </span>
                          )}
                        </button>
                      ) : (
                        <button
                          onClick={() => startModel(model.id)}
                          disabled={starting || isRunning || !hasSidecar}
                          className="w-full px-3 py-1.5 rounded text-xs font-medium bg-violet-500/10 text-violet-400 border border-violet-500/30 hover:bg-violet-500/20 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                          {starting ? (
                            <span className="flex items-center gap-1 justify-center">
                              <Loader2 className="w-3 h-3 animate-spin" /> Loading...
                            </span>
                          ) : (
                            <span className="flex items-center gap-1 justify-center">
                              <Cpu className="w-3 h-3" /> Load Model
                            </span>
                          )}
                        </button>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="p-4 rounded-lg bg-slate-800/20 border border-slate-700/30 text-center">
                <Brain className="w-8 h-8 text-slate-600 mx-auto mb-2" />
                <p className="text-sm text-slate-400 mb-2">No models found</p>
                <p className="text-xs text-slate-500">
                  Place <code className="text-violet-400/80 bg-slate-800 px-1 rounded">.gguf</code> model files in<br />
                  <code className="text-violet-400/80 bg-slate-800 px-1 rounded text-[10px]">content-packs/models-ai/</code>
                </p>
              </div>
            )}
          </div>

          {/* Status info */}
          {isRunning && (
            <div className="mt-auto p-4 border-t border-slate-800/50">
              <div className="p-3 rounded-lg bg-violet-500/5 border border-violet-500/20">
                <div className="flex items-center gap-2 mb-1">
                  <Sparkles className="w-4 h-4 text-violet-400" />
                  <span className="text-xs font-medium text-violet-300">Active Model</span>
                </div>
                <p className="text-sm text-slate-300 truncate">{status?.active_model}</p>
                <p className="text-xs text-slate-500 mt-1">Port: {status?.port}</p>
              </div>
            </div>
          )}
        </div>

        {/* Chat area */}
        <div className="flex-1 flex flex-col">
          {!isRunning ? (
            /* Empty state when no model is loaded */
            <div className="flex-1 flex items-center justify-center p-8">
              <div className="text-center max-w-md">
                <div className="w-20 h-20 rounded-2xl bg-gradient-to-br from-violet-500/20 to-purple-500/20 border border-violet-500/20 flex items-center justify-center mx-auto mb-6">
                  <Bot className="w-10 h-10 text-violet-400" />
                </div>
                <h2 className="text-xl font-semibold text-slate-200 mb-2">
                  Offline AI Assistant
                </h2>
                <p className="text-slate-400 mb-6">
                  {!hasSidecar 
                    ? 'Install llama-server and a GGUF model to enable the AI assistant.'
                    : !hasModels 
                    ? 'Download a GGUF model file and place it in content-packs/models-ai/ to get started.'
                    : 'Select and load a model from the sidebar to start chatting.'
                  }
                </p>
                {!hasSidecar && (
                  <div className="p-4 rounded-lg bg-slate-800/50 border border-slate-700/30 text-left">
                    <p className="text-sm text-slate-300 font-medium mb-2">Quick Setup:</p>
                    <ol className="text-xs text-slate-400 space-y-1 list-decimal list-inside">
                      <li>Download <code className="text-violet-400 bg-slate-800 px-1 rounded">llama-server</code> from llama.cpp releases</li>
                      <li>Place it in <code className="text-violet-400 bg-slate-800 px-1 rounded">sidecars/llama/windows/</code></li>
                      <li>Download a .gguf model (e.g. TinyLlama, Phi-3-mini)</li>
                      <li>Place it in <code className="text-violet-400 bg-slate-800 px-1 rounded">content-packs/models-ai/</code></li>
                      <li>Restart AEGIS</li>
                    </ol>
                  </div>
                )}
              </div>
            </div>
          ) : (
            /* Chat interface */
            <>
              {/* Messages */}
              <div className="flex-1 overflow-y-auto p-6 space-y-4">
                {messages.length === 0 && (
                  <div className="flex items-center justify-center h-full">
                    <div className="text-center">
                      <Sparkles className="w-8 h-8 text-violet-400 mx-auto mb-3" />
                      <p className="text-slate-400">Model loaded. Ask me anything.</p>
                      <p className="text-xs text-slate-500 mt-1">All processing happens locally on your device.</p>
                    </div>
                  </div>
                )}

                {messages.map((msg, i) => (
                  <div
                    key={i}
                    className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}
                  >
                    <div className={`max-w-[80%] rounded-xl px-4 py-3 ${
                      msg.role === 'user'
                        ? 'bg-violet-500/20 border border-violet-500/30 text-slate-100'
                        : 'bg-slate-800/50 border border-slate-700/30 text-slate-200'
                    }`}>
                      {msg.role === 'assistant' && (
                        <div className="flex items-center gap-2 mb-2">
                          <Bot className="w-4 h-4 text-violet-400" />
                          <span className="text-xs text-violet-400 font-medium">AEGIS AI</span>
                        </div>
                      )}
                      <div className="text-sm whitespace-pre-wrap leading-relaxed">
                        {msg.content || (
                          <span className="flex items-center gap-2 text-slate-500">
                            <Loader2 className="w-4 h-4 animate-spin" /> Thinking...
                          </span>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
                <div ref={chatEndRef} />
              </div>

              {/* Input area */}
              <div className="p-4 border-t border-slate-800/50">
                <div className="flex items-end gap-3 max-w-4xl mx-auto">
                  <div className="flex-1 relative">
                    <textarea
                      ref={inputRef}
                      value={input}
                      onChange={(e) => setInput(e.target.value)}
                      onKeyDown={handleKeyDown}
                      placeholder="Ask the AI assistant..."
                      rows={1}
                      className="w-full resize-none bg-slate-800/50 border border-slate-700/50 rounded-xl px-4 py-3 text-sm text-slate-100 placeholder-slate-500 focus:outline-none focus:border-violet-500/50 focus:ring-1 focus:ring-violet-500/20"
                      style={{ minHeight: '44px', maxHeight: '120px' }}
                    />
                  </div>
                  <button
                    onClick={sendMessage}
                    disabled={!input.trim() || streaming}
                    className="shrink-0 w-11 h-11 rounded-xl bg-violet-500 hover:bg-violet-600 disabled:bg-slate-700 disabled:cursor-not-allowed flex items-center justify-center transition-colors"
                  >
                    {streaming ? (
                      <Loader2 className="w-5 h-5 text-white animate-spin" />
                    ) : (
                      <Send className="w-5 h-5 text-white" />
                    )}
                  </button>
                </div>
                <p className="text-center text-xs text-slate-600 mt-2">
                  100% offline · All data stays on your device · Model: {status?.active_model}
                </p>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  )
}
