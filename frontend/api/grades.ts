import type { Grade, PaginatedResponse } from '~/types'

export const gradesApi = {
  list: (params?: Record<string, unknown>) =>
    useApi().get<PaginatedResponse<Grade>>('/grades', { params }),
  getById: (id: string) =>
    useApi().get<{ grade: Grade; history: any[] }>(`/grades/${id}`),
  create: (data: {
    task_name: string
    notes: string
    discipline_id: number
    grade_type_id: number
    student_id: string
  }) => useApi().post<Grade>('/grades', data),
  updateStatus: (id: string, status: string) =>
    useApi().patch(`/grades/${id}/status`, { status }),
  assign: (id: string) =>
    useApi().patch(`/grades/${id}/assign`),
  acknowledge: (id: string) =>
    useApi().patch(`/grades/${id}/acknowledge`),
}
