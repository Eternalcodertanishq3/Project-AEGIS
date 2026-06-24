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
        const [profData, powerData] = await Promise.all([
          apiFetch<any>('/system/profile'),
          apiFetch<any>('/system/power')
        ])
        
        if (!cancelled) {
          setProfile({
            hardwareTier: (profData.tier.charAt(0).toUpperCase() + profData.tier.slice(1)) as any,
            os: profData.os,
            arch: profData.arch,
            cpuCores: profData.cpu_cores,
            ramTotal: `${(profData.total_ram_mb / 1024).toFixed(1)} GB`,
            ramUsed: 'Unknown',
            ramPercent: 0,
            batteryPercent: powerData.battery_percent,
            batteryCharging: powerData.status === 'charging' || powerData.status === 'ac_power',
            isOnline: false,
            uptime: 'Unknown',
            hostname: profData.hostname
          })
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
