import { defineStore } from 'pinia'
import type { Grade, GradeHistory } from '~/types'
import { gradesApi } from '~/api/grades'

export const useGradesStore = defineStore('grades', {
  state: () => ({
    list: [] as Grade[],
    total: 0,
    page: 1,
    limit: 20,
    current: null as Grade | null,
    history: [] as GradeHistory[],
    loading: false,
    filters: {
      status: '',
      grade_type_id: null as number | null,
      discipline_id: null as number | null,
    },
  }),
  actions: {
    async fetchList() {
      this.loading = true
      try {
        const params: Record<string, unknown> = { page: this.page, limit: this.limit }
        if (this.filters.status) params.status = this.filters.status
        if (this.filters.grade_type_id) params.grade_type_id = this.filters.grade_type_id
        if (this.filters.discipline_id) params.discipline_id = this.filters.discipline_id
        const { data } = await gradesApi.list(params)
        this.list = data.data
        this.total = data.total
      } finally {
        this.loading = false
      }
    },
    async fetchById(id: string) {
      this.loading = true
      try {
        const { data } = await gradesApi.getById(id)
        this.current = data.grade
        this.history = data.history
      } finally {
        this.loading = false
      }
    },
    setPage(p: number) { this.page = p; this.fetchList() },
    setFilters(f: Partial<typeof this.filters>) {
      this.filters = { ...this.filters, ...f }
      this.page = 1
      this.fetchList()
    },
  },
})
