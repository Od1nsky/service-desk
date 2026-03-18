import axios from 'axios'

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()

  const api = axios.create({
    baseURL: config.public.apiBase,
    headers: { 'Content-Type': 'application/json' },
  })

  api.interceptors.request.use((cfg) => {
    if (import.meta.client) {
      const token = localStorage.getItem('access_token')
      if (token) cfg.headers.Authorization = `Bearer ${token}`
    }
    return cfg
  })

  let refreshing = false
  let queue: Array<(token: string) => void> = []

  api.interceptors.response.use(
    (r) => r,
    async (error) => {
      const original = error.config
      if (error.response?.status === 401 && !original._retry) {
        original._retry = true
        if (refreshing) {
          return new Promise((resolve) => {
            queue.push((token) => {
              original.headers.Authorization = `Bearer ${token}`
              resolve(api(original))
            })
          })
        }
        refreshing = true
        try {
          const rt = localStorage.getItem('refresh_token')
          if (!rt) throw new Error('no refresh token')
          const { data } = await axios.post(`${config.public.apiBase}/auth/refresh`, { refresh_token: rt })
          localStorage.setItem('access_token', data.access_token)
          localStorage.setItem('refresh_token', data.refresh_token)
          queue.forEach((cb) => cb(data.access_token))
          queue = []
          original.headers.Authorization = `Bearer ${data.access_token}`
          return api(original)
        } catch {
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
          localStorage.removeItem('user')
          if (import.meta.client) window.location.href = '/auth/login'
        } finally {
          refreshing = false
        }
      }
      return Promise.reject(error)
    },
  )

  return { provide: { api } }
})
