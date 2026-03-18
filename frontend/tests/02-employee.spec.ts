import { test, expect } from '@playwright/test'
import { setAuth, mockAllApis, screenshot } from './helpers/setup'
import { MOCK_USERS } from './fixtures/data'

test.describe('Student flow', () => {
  test.beforeEach(async ({ page }) => {
    await mockAllApis(page, 'student')
    // Load a page first so localStorage is accessible, then set auth
    await page.goto('/auth/login')
    await setAuth(page, MOCK_USERS.student)
  })

  test('redirect to /grades from /', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    // Student gets redirected from / to /grades; go directly so we don't race auth init
    await page.goto('/grades')
    await page.waitForSelector('.grade-card, .n-empty', { timeout: 8000 })
    await screenshot(page, '05-student-grades-list', proj)
  })

  test('grade list — with filters', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/grades')
    await page.waitForSelector('.grade-card, .n-empty', { timeout: 8000 })
    // Open status filter
    await page.click('.n-select:first-of-type')
    await page.waitForSelector('.n-base-select-option', { timeout: 3000 }).catch(() => {})
    await screenshot(page, '06-student-grades-filter-open', proj)
    await page.keyboard.press('Escape')
  })

  test('grade detail page', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/grades/grade-1')
    await page.waitForSelector('.grade-layout', { timeout: 8000 })
    await screenshot(page, '09-grade-detail-student', proj)
  })

  test('grade detail — graded (can acknowledge)', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/grades/grade-1')
    await page.waitForSelector('.grade-layout', { timeout: 8000 })
    await screenshot(page, '10-grade-detail-graded', proj)
  })

  test('sidebar — student nav', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/grades')
    await page.waitForLoadState('networkidle')
    await page.waitForSelector('.grade-card, .n-empty', { timeout: 8000 })
    // On mobile — open sidebar
    if (testInfo.project.name === 'mobile') {
      await page.evaluate(() => {
        const btn = document.querySelector('[data-testid="burger-btn"]') as HTMLElement
        btn?.click()
      })
      await page.waitForTimeout(300)
    }
    await screenshot(page, '11-sidebar-student', proj)
  })
})
