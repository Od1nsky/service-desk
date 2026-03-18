import { test, expect } from '@playwright/test'
import { mockAllApis, screenshot } from './helpers/setup'

test.describe('Auth pages', () => {
  test('login page — desktop', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await mockAllApis(page, 'student')
    await page.goto('/auth/login')
    await page.waitForSelector('.auth-card')
    await screenshot(page, '01-login', proj)
  })

  test('login page — validation errors', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await mockAllApis(page, 'student')
    await page.goto('/auth/login')
    await page.waitForSelector('.auth-card')
    // Click submit without filling form to trigger validation
    await page.click('button:has-text("Войти")')
    await page.waitForTimeout(400)
    await screenshot(page, '02-login-validation', proj)
  })

  test('register page', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await mockAllApis(page, 'student')
    await page.goto('/auth/register')
    await page.waitForSelector('.auth-card')
    await screenshot(page, '03-register', proj)
  })

  test('login page — filled form', async ({ page }, testInfo) => {
    const proj = testInfo.project.name
    await mockAllApis(page, 'student')
    await page.goto('/auth/login')
    await page.waitForSelector('.auth-card')
    // Fill email input (n-input renders an actual input inside)
    await page.fill('.auth-card input[type="text"], .auth-card input:not([type="password"])', 'ivanov@university.ru')
    await page.fill('.auth-card input[type="password"]', 'password123')
    await screenshot(page, '04-login-filled', proj)
  })
})
