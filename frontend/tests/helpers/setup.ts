import type { Page } from '@playwright/test'
import {
  MOCK_USERS,
  MOCK_GRADE_TYPES,
  MOCK_DISCIPLINES,
  MOCK_GRADES,
  MOCK_COMMENTS,
  MOCK_HISTORY,
  MOCK_ALL_USERS,
  MOCK_DASHBOARD,
  MOCK_DISCIPLINE_GRADES,
} from '../fixtures/data'
import type { User } from '../../types'

// A fake but plausible JWT (won't validate, but middleware only checks existence)
const FAKE_TOKEN =
  'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTEiLCJlbWFpbCI6InRlc3RAdW5pdmVyc2l0eS5ydSIsInJvbGUiOiJzdHVkZW50IiwiZXhwIjo5OTk5OTk5OTk5fQ.fake'

/** Set auth tokens + user in localStorage so auth middleware passes */
export async function setAuth(page: Page, user: User) {
  await page.evaluate(
    ({ user, token }) => {
      localStorage.setItem('access_token', token)
      localStorage.setItem('refresh_token', 'fake-refresh-uuid')
      localStorage.setItem('user', JSON.stringify(user))
    },
    { user, token: FAKE_TOKEN },
  )
}

/** Intercept all /api/v1/* requests and return mock data */
export async function mockAllApis(
  page: Page,
  role: 'student' | 'teacher' | 'admin' = 'student',
) {
  const user = MOCK_USERS[role]

  await page.route('**/api/v1/**', async (route) => {
    const url = route.request().url()
    const method = route.request().method()

    // Auth
    if (url.includes('/auth/login')) {
      return route.fulfill({
        json: { access_token: FAKE_TOKEN, refresh_token: 'fake-refresh', user },
      })
    }
    if (
      url.includes('/auth/refresh') ||
      url.includes('/auth/logout') ||
      url.includes('/auth/register')
    ) {
      return route.fulfill({ json: { success: true } })
    }

    // Grade types
    if (url.includes('/grade-types')) {
      return route.fulfill({ json: MOCK_GRADE_TYPES })
    }

    // Disciplines
    if (url.match(/\/disciplines\/\d+/) && method === 'DELETE') {
      return route.fulfill({ status: 204, body: '' })
    }
    if (url.match(/\/disciplines\/\d+/) && (method === 'PUT' || method === 'PATCH')) {
      return route.fulfill({ json: MOCK_DISCIPLINES[0] })
    }
    if (url.includes('/disciplines') && method === 'POST') {
      return route.fulfill({
        status: 201,
        json: { ...MOCK_DISCIPLINES[0], id: 99, name: 'Новая дисциплина' },
      })
    }
    if (url.includes('/disciplines')) {
      return route.fulfill({ json: MOCK_DISCIPLINES })
    }

    // Discipline grades
    if (url.includes('/discipline-grades') && method === 'POST') {
      return route.fulfill({ json: MOCK_DISCIPLINE_GRADES[0] })
    }
    if (url.includes('/discipline-grades')) {
      return route.fulfill({ json: MOCK_DISCIPLINE_GRADES })
    }

    // Stats
    if (url.includes('/stats/dashboard')) {
      return route.fulfill({ json: MOCK_DASHBOARD })
    }

    // Users
    if (url.includes('/users') && method === 'PATCH') {
      return route.fulfill({ json: { success: true } })
    }
    if (url.includes('/users')) {
      return route.fulfill({ json: MOCK_ALL_USERS })
    }

    // Grades - comments
    if (url.match(/\/grades\/[^/]+\/comments/) && method === 'POST') {
      return route.fulfill({
        status: 201,
        json: {
          id: 'new-comment',
          ticket_id: 'grade-1',
          author_id: user.id,
          body: 'Новый комментарий',
          is_internal: false,
          created_at: new Date().toISOString(),
          author: { id: user.id, full_name: user.full_name, email: user.email },
        },
      })
    }
    if (url.match(/\/grades\/[^/]+\/comments/)) {
      return route.fulfill({
        json: MOCK_COMMENTS.filter((c) => !c.is_internal || role !== 'student'),
      })
    }

    // Grade actions
    if (url.match(/\/grades\/[^/]+\/status/)) {
      return route.fulfill({ json: { message: 'status updated' } })
    }
    if (url.match(/\/grades\/[^/]+\/assign/)) {
      return route.fulfill({ json: { message: 'grade assigned' } })
    }
    if (url.match(/\/grades\/[^/]+\/acknowledge/)) {
      return route.fulfill({ json: { message: 'acknowledged' } })
    }

    // Grade by ID
    const gradeByIdMatch = url.match(/\/grades\/([^/?]+)$/)
    if (gradeByIdMatch && method === 'GET') {
      const id = gradeByIdMatch[1]
      const grade = MOCK_GRADES.find((g) => g.id === id) || MOCK_GRADES[0]
      return route.fulfill({ json: { grade, history: MOCK_HISTORY } })
    }

    // Create grade
    if (url.includes('/grades') && method === 'POST') {
      return route.fulfill({ status: 201, json: MOCK_GRADES[0] })
    }

    // Grade list
    if (url.includes('/grades')) {
      const filtered = MOCK_GRADES
      return route.fulfill({
        json: { data: filtered, total: filtered.length, page: 1, limit: 20 },
      })
    }

    // Fallback
    return route.fulfill({ status: 404, json: { error: 'not found' } })
  })
}

/** Screenshot helper — saves to reports/screenshots/<project>/<name>.png */
export async function screenshot(page: Page, name: string, projectName: string = 'desktop') {
  const { mkdirSync } = await import('fs')
  const { resolve, dirname } = await import('path')
  const { fileURLToPath } = await import('url')
  const __filename = fileURLToPath(import.meta.url)
  const __dirname = dirname(__filename)
  const dir = resolve(__dirname, '../../../reports/screenshots', projectName)
  mkdirSync(dir, { recursive: true })
  await page.screenshot({
    path: resolve(dir, `${name}.png`),
    fullPage: true,
    animations: 'disabled',
  })
}
