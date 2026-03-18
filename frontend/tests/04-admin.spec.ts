import { test } from '@playwright/test'
import { setAuth, mockAllApis, screenshot } from './helpers/setup'
import { MOCK_USERS } from './fixtures/data'

test.describe('Admin flow', () => {
  test.beforeEach(async ({ page }) => {
    await mockAllApis(page, 'admin')
    // Load a page first so localStorage is accessible, then set auth
    await page.goto('/auth/login')
    await setAuth(page, MOCK_USERS.admin)
  })

  test('dashboard — admin view with stats', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/')
    await page.waitForSelector('.stats-grid', { timeout: 8000 })
    await screenshot(page, '18-dashboard-admin', proj)
  })

  test('users management page', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/admin/users')
    await page.waitForSelector('.n-data-table', { timeout: 8000 })
    await screenshot(page, '19-admin-users', proj)
  })

  test('disciplines management page', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/admin/disciplines')
    await page.waitForSelector('.n-data-table', { timeout: 8000 })
    await screenshot(page, '20-admin-disciplines', proj)
  })

  test('disciplines — create modal open', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/admin/disciplines')
    await page.waitForSelector('.n-data-table', { timeout: 8000 })
    await page.click('button:has-text("Добавить дисциплину")')
    // Wait for modal content (n-card inside the modal) to appear
    await page.waitForSelector('.n-card:has(.n-form)', { timeout: 5000 })
    await screenshot(page, '21-admin-disciplines-modal', proj)
  })

  test('disciplines — edit modal with data', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/admin/disciplines')
    await page.waitForSelector('.n-data-table', { timeout: 8000 })
    // Click first "Изменить" button
    await page.click('button:has-text("Изменить")')
    // Wait for modal content to appear
    await page.waitForSelector('.n-card:has(.n-form)', { timeout: 5000 })
    await screenshot(page, '22-admin-disciplines-edit', proj)
  })

  test('teacher queue — admin sees all', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/teacher/queue')
    await page.waitForSelector('.n-data-table, .n-empty', { timeout: 8000 })
    await screenshot(page, '23-admin-queue', proj)
  })

  test('sidebar — admin full nav', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    await page.waitForSelector('.stats-grid', { timeout: 8000 })
    // On mobile — open sidebar
    if (testInfo.project.name === 'mobile') {
      await page.evaluate(() => {
        const btn = document.querySelector('[data-testid="burger-btn"]') as HTMLElement
        btn?.click()
      })
      await page.waitForTimeout(300)
    }
    await screenshot(page, '24-sidebar-admin', proj)
  })
})
