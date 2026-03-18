import type { Comment } from '~/types'

export const commentsApi = {
  list: (gradeId: string) =>
    useApi().get<Comment[]>(`/grades/${gradeId}/comments`),
  create: (gradeId: string, data: { body: string; is_internal: boolean }) =>
    useApi().post<Comment>(`/grades/${gradeId}/comments`, data),
}
