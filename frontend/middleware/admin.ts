export default defineNuxtRouteMiddleware(() => {
  if (!import.meta.client) return
  const u = localStorage.getItem('user')
  if (!u) return navigateTo('/auth/login')
  const user = JSON.parse(u)
  if (user.role !== 'admin') return navigateTo('/')
})
