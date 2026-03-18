import type { AuthTokens } from '~/types'

export const authApi = {
  register: (data: { email: string; password: string; full_name: string }) =>
    useApi().post('/auth/register', data),
  login: (data: { email: string; password: string }) =>
    useApi().post<AuthTokens>('/auth/login', data),
  refresh: (refresh_token: string) =>
    useApi().post('/auth/refresh', { refresh_token }),
  logout: (refresh_token: string) =>
    useApi().post('/auth/logout', { refresh_token }),
}
