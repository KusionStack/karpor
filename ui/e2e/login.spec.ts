import { expect, test } from '@playwright/test'
import { mockApi, VALID_TOKEN } from './mock-api'

test.describe('login', () => {
  test.beforeEach(async ({ page }) => {
    await mockApi(page, { rbac: true })
  })

  test('asks for a token before showing a protected page', async ({ page }) => {
    await page.goto('/insight')
    await expect(page).toHaveURL(/\/login$/)
    await expect(
      page.getByText('Not logged in. Please enter token to login'),
    ).toBeVisible()
  })

  test('rejects an empty token', async ({ page }) => {
    await page.goto('/login')
    await page.getByRole('button', { name: 'Login' }).click()
    await expect(page.getByText('Enter Token', { exact: true })).toBeVisible()
  })

  test('rejects an invalid token', async ({ page }) => {
    await page.goto('/login')
    await page.getByRole('textbox').fill('wrong-token')
    await page.getByRole('button', { name: 'Login' }).click()
    await expect(
      page.getByText('Login failed. Please check token validity'),
    ).toBeVisible()
  })

  test('logs in with a valid token and returns to the page', async ({
    page,
  }) => {
    await page.goto('/insight')
    await expect(page).toHaveURL(/\/login$/)

    await page.getByRole('textbox').fill(VALID_TOKEN)
    await page.getByRole('button', { name: 'Login' }).click()

    await expect(page.getByText('Login successful')).toBeVisible()
    await expect(page).toHaveURL(/\/insight$/)
    expect(
      await page.evaluate(() => window.localStorage.getItem('token')),
    ).toBe(VALID_TOKEN)
  })
})
