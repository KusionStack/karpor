import { defineConfig, devices } from '@playwright/test'

const PORT = 8000

export default defineConfig({
  testDir: './e2e',
  timeout: 60 * 1000,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  workers: 1,
  reporter: process.env.CI ? [['list'], ['html', { open: 'never' }]] : 'list',
  use: {
    baseURL: `http://localhost:${PORT}`,
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
  },
  projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'] } }],
  // Start the dev server. The tests mock every backend call, so no karpor
  // server, Elasticsearch or cluster is needed.
  webServer: {
    command: 'npm run dev',
    url: `http://localhost:${PORT}`,
    reuseExistingServer: !process.env.CI,
    timeout: 300 * 1000,
    env: { BROWSER: 'none' },
  },
})
