import { useState, useEffect } from 'react'
import {
  Database, Copy, Check, Loader2, ArrowRightLeft
} from 'lucide-react'

interface OperationInfo {
  id: string
  name: string
  description: string
}

interface TransformResult {
  input: string
  output: string
  operation: string
  error?: string
}

export function DataToolsPage() {
  const [operations, setOperations] = useState<Record<string, OperationInfo[]>>({})
  const [loading, setLoading] = useState(true)
  const [input, setInput] = useState('')
  const [output, setOutput] = useState('')
  const [selectedOp, setSelectedOp] = useState('base64-encode')
  const [error, setError] = useState('')
  const [copied, setCopied] = useState(false)
  const [transforming, setTransforming] = useState(false)

  useEffect(() => {
    const fetchOps = async () => {
      try {
        const res = await fetch('/api/datatools/operations')
        const data = await res.json()
        setOperations(data.operations || {})
      } catch (e) {
        console.error('Failed to fetch operations:', e)
      } finally {
        setLoading(false)
      }
    }
    fetchOps()
  }, [])

  const transform = async () => {
    if (!input.trim()) return
    setTransforming(true)
    setError('')
    try {
      const res = await fetch('/api/datatools/transform', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ operation: selectedOp, input })
      })
      const data: TransformResult = await res.json()
      if (data.error) {
        setError(data.error)
        setOutput('')
      } else {
        setOutput(data.output)
      }
    } catch (e) {
      setError('Transform failed')
    } finally {
      setTransforming(false)
    }
  }

  const copyOutput = async () => {
    try {
      await navigator.clipboard.writeText(output)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch {
      // Fallback for offline/insecure context
      const ta = document.createElement('textarea')
      ta.value = output
      document.body.appendChild(ta)
      ta.select()
      document.execCommand('copy')
      document.body.removeChild(ta)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    }
  }

  const swapIO = () => {
    setInput(output)
    setOutput('')
    setError('')
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="w-8 h-8 text-cyan-500 animate-spin" />
      </div>
    )
  }

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="p-6 pb-4 border-b border-slate-800/50">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-cyan-500/20 border border-cyan-500/30 flex items-center justify-center">
            <Database className="w-5 h-5 text-cyan-400" />
          </div>
          <div>
            <h1 className="text-xl font-semibold text-slate-100">Data Tools</h1>
            <p className="text-sm text-slate-400">
              Encode, decode, hash, and transform data — all offline
            </p>
          </div>
        </div>
      </div>

      <div className="flex-1 flex overflow-hidden">
        {/* Operations sidebar */}
        <div className="w-64 border-r border-slate-800/50 bg-slate-900/30 overflow-y-auto py-3 px-2">
          {Object.entries(operations).map(([group, ops]) => (
            <div key={group} className="mb-4">
              <p className="px-2 mb-1.5 text-[10px] font-semibold uppercase tracking-wider text-slate-500">
                {group}
              </p>
              <div className="space-y-0.5">
                {ops.map(op => (
                  <button
                    key={op.id}
                    onClick={() => { setSelectedOp(op.id); setError('') }}
                    className={`w-full text-left px-3 py-2 rounded-lg text-sm transition-all ${
                      selectedOp === op.id
                        ? 'bg-cyan-500/10 text-cyan-400 border border-cyan-500/20'
                        : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50 border border-transparent'
                    }`}
                  >
                    <span className="font-medium">{op.name}</span>
                    {selectedOp === op.id && (
                      <p className="text-[10px] text-cyan-500/70 mt-0.5">{op.description}</p>
                    )}
                  </button>
                ))}
              </div>
            </div>
          ))}
        </div>

        {/* Transform area */}
        <div className="flex-1 flex flex-col p-6 gap-4">
          {/* Input */}
          <div className="flex-1 flex flex-col">
            <label className="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Input</label>
            <textarea
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder="Enter text to transform..."
              className="flex-1 p-4 rounded-xl bg-slate-800/30 border border-slate-700/30 text-sm text-slate-200 placeholder-slate-600 font-mono resize-none focus:outline-none focus:border-cyan-500/50"
            />
          </div>

          {/* Controls */}
          <div className="flex items-center gap-3 justify-center">
            <button
              onClick={transform}
              disabled={!input.trim() || transforming}
              className="flex items-center gap-2 px-6 py-2.5 rounded-xl bg-cyan-500 hover:bg-cyan-600 disabled:bg-slate-700 text-white text-sm font-semibold transition-colors disabled:cursor-not-allowed"
            >
              {transforming ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : (
                <ArrowRightLeft className="w-4 h-4" />
              )}
              Transform
            </button>
            <button
              onClick={swapIO}
              disabled={!output}
              className="px-3 py-2.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-400 text-sm transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              title="Swap input ↔ output"
            >
              ↕ Swap
            </button>
          </div>

          {/* Error */}
          {error && (
            <div className="p-3 rounded-lg bg-red-500/10 border border-red-500/20 text-sm text-red-400">
              {error}
            </div>
          )}

          {/* Output */}
          <div className="flex-1 flex flex-col">
            <div className="flex items-center justify-between mb-2">
              <label className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Output</label>
              {output && (
                <button
                  onClick={copyOutput}
                  className="flex items-center gap-1 px-2 py-1 rounded text-xs text-slate-400 hover:text-emerald-400 transition-colors"
                >
                  {copied ? (
                    <><Check className="w-3 h-3" /> Copied!</>
                  ) : (
                    <><Copy className="w-3 h-3" /> Copy</>
                  )}
                </button>
              )}
            </div>
            <textarea
              value={output}
              readOnly
              placeholder="Result will appear here..."
              className="flex-1 p-4 rounded-xl bg-slate-800/30 border border-slate-700/30 text-sm text-emerald-400 placeholder-slate-600 font-mono resize-none focus:outline-none"
            />
          </div>
        </div>
      </div>
    </div>
  )
}
