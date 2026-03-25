export type Role = 'student' | 'teacher' | 'admin'
export type GradeStatus = 'pending' | 'graded' | 'not_certified' | 'certified' | 'closed'

export interface User {
  id: string
  email: string
  full_name: string
  role: Role
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface UserInfo {
  id: string
  full_name: string
  email: string
}

export interface Discipline {
  id: number
  name: string
  description: string
  created_at: string
}

export interface GradeType {
  id: number
  name: string
  color: string
  max_score: number
}

export interface Grade {
  id: string
  number: number
  task_name: string
  notes: string
  status: GradeStatus
  student_id: string
  teacher_id: string | null
  discipline_id: number
  grade_type_id: number
  score: number | null
  created_at: string
  updated_at: string
  graded_at: string | null
  student?: UserInfo
  teacher?: UserInfo
  discipline?: Discipline
  grade_type?: GradeType
}

export interface Comment {
  id: string
  ticket_id: string
  author_id: string
  body: string
  is_internal: boolean
  created_at: string
  author?: UserInfo
}

export interface GradeHistory {
  id: string
  ticket_id: string
  changed_by: string
  old_value: string
  new_value: string
  changed_at: string
  user?: UserInfo
}

export type DisciplineGradeStatus = 'pending' | 'certified' | 'not_certified'

export interface DisciplineGrade {
  id?: string
  student_id: string
  discipline_id: number
  teacher_id?: string | null
  extra_score: number | null
  notes: string
  status: DisciplineGradeStatus
  works_total: number | null
  total: number | null
  student?: UserInfo
  discipline?: Discipline
  teacher?: UserInfo
  updated_at?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
}

export interface AuthTokens {
  access_token: string
  refresh_token: string
  user: User
}

export interface TeacherStat {
  id: string
  full_name: string
  count: number
}

export interface DashboardStats {
  total_grades: number
  pending_grades: number
  avg_score: number | null
  grades_by_status: Partial<Record<GradeStatus, number>>
  top_teachers: TeacherStat[]
}
