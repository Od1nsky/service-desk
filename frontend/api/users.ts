import type { User } from '~/types'

export const usersApi = {
  list: () => useApi().get<User[]>('/users'),
  updateRole: (id: string, role: string) =>
    useApi().patch(`/users/${id}/role`, { role }),
  block: (id: string, is_active: boolean) =>
    useApi().patch(`/users/${id}/block`, { is_active }),
}
