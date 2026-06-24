import { useState, useCallback, useEffect } from 'react'
import type { Module } from '@/types'
import { defaultModules } from '@/data/defaults'
import { apiFetch } from './useApi'

export function useModules() {
  const [modules, setModules] = useState<Module[]>(defaultModules)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let cancelled = false

    async function fetchModules() {
      try {
        const data = await apiFetch<{modules: Module[], count: number}>('/modules')
        if (!cancelled) {
          setModules(data.modules || [])
        }
      } catch {
        // Fall back to defaults when backend is unavailable
      } finally {
        if (!cancelled) {
          setLoading(false)
        }
      }
    }

    fetchModules()

    return () => {
      cancelled = true
    }
  }, [])

  const toggleModule = useCallback(async (moduleId: string) => {
    setModules(prev =>
      prev.map(m =>
        m.id === moduleId
          ? { ...m, enabled: !m.enabled, status: !m.enabled ? 'active' as const : 'inactive' as const }
          : m
      )
    )

    try {
      await apiFetch(`/modules/${moduleId}/toggle`, { method: 'POST' })
    } catch {
      // Optimistic update — if API fails, the toggle still works locally
    }
  }, [])

  const getModulesByDomain = useCallback((domain: string) => {
    return modules.filter(m => m.domain === domain)
  }, [modules])

  const activeCount = modules.filter(m => m.status === 'active').length

  return { modules, loading, toggleModule, getModulesByDomain, activeCount }
}
