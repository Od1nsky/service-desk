import type { GradeType } from '~/types'

export const gradeTypesApi = {
  list: () => useApi().get<GradeType[]>('/grade-types'),
}
