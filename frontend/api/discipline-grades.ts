import type { DisciplineGrade } from '~/types'

export const disciplineGradesApi = {
  list: (studentId?: string) =>
    useApi().get<DisciplineGrade[]>('/discipline-grades', {
      params: studentId ? { student_id: studentId } : {},
    }),
  upsert: (data: { student_id: string; discipline_id: number; extra_score: number; notes?: string; status?: string }) =>
    useApi().post<DisciplineGrade>('/discipline-grades', data),
}
