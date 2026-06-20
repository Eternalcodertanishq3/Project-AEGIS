import { useState, useEffect } from 'react'
import type { SystemProfile } from '@/types'
import { defaultSystemProfile } from '@/data/defaults'
import { apiFetch } from './useApi'

export function useSystemProfile() {
  const [profile, setProfile] = useState<SystemProfile>(defaultSystemProfile)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false

    async function fetchProfile() {
      try {
        const data = await apiFetch<SystemProfile>('/system/profile')
        if (!cancelled) {
          setProfile(data)
          setError(null)
        }
      } catch {
        // Fall back to defaults when backend is unavailable
        if (!cancelled) {
          setError('Backend unavailable — using defaults')
        }
      } finally {
        if (!cancelled) {
          setLoading(false)
        }
      }
    }

    fetchProfile()

    // Poll every 30 seconds
    const interval = setInterval(fetchProfile, 30000)

    return () => {
      cancelled = true
      clearInterval(interval)
    }
  }, [])

  return { profile, loading, error }
}
