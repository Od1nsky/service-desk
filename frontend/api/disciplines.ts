import type { Discipline } from '~/types'

export const disciplinesApi = {
  list: () => useApi().get<Discipline[]>('/disciplines'),
  create: (data: { name: string; description: string }) =>
    useApi().post<Discipline>('/disciplines', data),
  update: (id: number, data: { name: string; description: string }) =>
    useApi().put<Discipline>(`/disciplines/${id}`, data),
  delete: (id: number) => useApi().delete(`/disciplines/${id}`),
}
