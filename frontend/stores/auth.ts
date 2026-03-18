import { defineStore } from 'pinia'
import type { User } from '~/types'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    accessToken: '' as string,
    refreshToken: '' as string,
  }),
  getters: {
    isAuthenticated: (s) => !!s.accessToken,
    isAdmin: (s) => s.user?.role === 'admin',
    isTeacher: (s) => s.user?.role === 'teacher' || s.user?.role === 'admin',
    isStudent: (s) => s.user?.role === 'student',
  },
  actions: {
    init() {
      if (!import.meta.client) return
      this.accessToken = localStorage.getItem('access_token') || ''
      this.refreshToken = localStorage.getItem('refresh_token') || ''
      const u = localStorage.getItem('user')
      if (u) this.user = JSON.parse(u)
    },
    setTokens(access: string, refresh: string, user: User) {
      this.accessToken = access
      this.refreshToken = refresh
      this.user = user
      if (import.meta.client) {
        localStorage.setItem('access_token', access)
        localStorage.setItem('refresh_token', refresh)
        localStorage.setItem('user', JSON.stringify(user))
      }
    },
    async logout() {
      try {
        if (this.refreshToken) {
          const { authApi } = await import('~/api/auth')
          await authApi.logout(this.refreshToken)
        }
      } catch {}
      this.accessToken = ''
      this.refreshToken = ''
      this.user = null
      if (import.meta.client) {
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        localStorage.removeItem('user')
      }
      navigateTo('/auth/login')
    },
  },
})
