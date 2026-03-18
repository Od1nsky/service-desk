import type { DashboardStats } from '~/types'

export const statsApi = {
  dashboard: () => useApi().get<DashboardStats>('/stats/dashboard'),
}
