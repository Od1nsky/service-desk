import { test } from '@playwright/test'
import { setAuth, mockAllApis, screenshot } from './helpers/setup'
import { MOCK_USERS } from './fixtures/data'

test.describe('Teacher flow', () => {
  test.beforeEach(async ({ page }) => {
    await mockAllApis(page, 'teacher')
    // Load a page first so localStorage is accessible, then set auth
    await page.goto('/auth/login')
    await setAuth(page, MOCK_USERS.teacher)
  })

  test('dashboard — teacher view', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/')
    await page.waitForSelector('.stats-grid', { timeout: 8000 })
    await screenshot(page, '12-dashboard-teacher', proj)
  })

  test('teacher queue — all grades', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/teacher/queue')
    await page.waitForSelector('.n-data-table, .n-empty', { timeout: 8000 })
    await screenshot(page, '13-teacher-queue', proj)
  })

  test('grade detail — teacher can assign and change status', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/grades/grade-2') // pending grade, no teacher
    await page.waitForSelector('.grade-layout', { timeout: 8000 })
    await screenshot(page, '14-grade-detail-teacher-pending', proj)
  })

  test('grade detail — teacher with status dropdown', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/grades/grade-1') // graded grade
    await page.waitForSelector('.grade-layout', { timeout: 8000 })
    // Open status select (the one in header-actions for status change)
    await page.click('.header-actions .n-select').catch(async () => {
      await page.click('.n-select')
    })
    await page.waitForSelector('.n-base-select-option', { timeout: 3000 }).catch(() => {})
    await screenshot(page, '15-grade-status-dropdown', proj)
    await page.keyboard.press('Escape')
  })

  test('grade detail — with comments', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/grades/grade-1')
    await page.waitForSelector('.comment-list, .empty-comments', { timeout: 8000 })
    await screenshot(page, '16-grade-comments-teacher', proj)
  })

  test('comment form — internal note checkbox visible for teacher', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/grades/grade-1')
    await page.waitForSelector('.comment-form', { timeout: 8000 })
    await page.fill('.comment-form textarea', 'Проверил работу, обнаружена ошибка в третьей задаче')
    await screenshot(page, '17-comment-form-teacher', proj)
  })

  test('create grade form — empty', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/grades/create')
    await page.waitForSelector('.n-form', { timeout: 8000 })
    await screenshot(page, '07-create-grade-empty', proj)
  })

  test('create grade form — filled', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/grades/create')
    await page.waitForSelector('.n-form', { timeout: 8000 })
    // Fill task_name field
    await page.fill('.n-form .n-input input, .n-form .n-input textarea', 'Контрольная работа по теме «Пределы»')
    await screenshot(page, '08-create-grade-filled', proj)
  })
})
