import { expect, test } from '@playwright/test'
import { deployment, mockApi } from './mock-api'

test.describe('search', () => {
  test('runs an SQL search and shows the results', async ({ page }) => {
    await mockApi(page, {
      searchItems: [deployment('web'), deployment('worker', 'jobs')],
    })

    await page.goto('/search')
    const editor = page.locator('.cm-content')
    await editor.click()
    await page.keyboard.type("where kind='Deployment'")

    const searchRequest = page.waitForRequest('**/rest-api/v1/search**')
    await page.keyboard.press('Escape')
    await editor.press('Enter')
    const request = await searchRequest

    const params = new URL(request.url()).searchParams
    expect(params.get('pattern')).toBe('sql')
    expect(params.get('query')).toBe(
      "select * from resources where kind='Deployment'",
    )

    await expect(page).toHaveURL(/\/search\/result\?/)
    await expect(page.getByText('Search Result')).toBeVisible()
    await expect(page.getByText('web', { exact: true })).toBeVisible()
    await expect(page.getByText('worker', { exact: true })).toBeVisible()
  })
})
