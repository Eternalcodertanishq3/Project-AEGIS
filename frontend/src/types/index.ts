export type ModuleStatus = 'active' | 'inactive' | 'unavailable'
export type ModuleDomain = 'Knowledge' | 'Survival' | 'Comms' | 'AI' | 'System'
export type HardwareTier = 'Minimum' | 'Standard' | 'Optimal'

export interface Module {
  id: string
  name: string
  description: string
  domain: ModuleDomain
  status: ModuleStatus
  enabled: boolean
  icon: string
}

export interface SystemProfile {
  hardwareTier: HardwareTier
  os: string
  arch: string
  cpuCores: number
  ramTotal: string
  ramUsed: string
  ramPercent: number
  batteryPercent: number
  batteryCharging: boolean
  isOnline: boolean
  uptime: string
  hostname: string
}
