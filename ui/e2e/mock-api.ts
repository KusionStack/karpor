import type { Page, Route } from '@playwright/test'

export const VALID_TOKEN = 'e2e-valid-token'

type MockOptions = {
  // When true the UI asks for a token before it shows protected pages.
  rbac?: boolean
  // Items returned by GET /rest-api/v1/search.
  searchItems?: unknown[]
}

function json(route: Route, body: unknown, status = 200) {
  return route.fulfill({
    status,
    contentType: 'application/json',
    body: JSON.stringify(body),
  })
}

// Stub the karpor server so the UI can be tested on its own.
export async function mockApi(page: Page, options: MockOptions = {}) {
  const { rbac = false, searchItems = [] } = options

  await page.addInitScript(() => {
    window.localStorage.setItem('lang', 'en')
  })

  // Catch-all first, so the specific routes below take precedence.
  await page.route('**/rest-api/**', route =>
    json(route, { success: true, data: {} }),
  )

  await page.route('**/server-configs', route =>
    json(route, {
      Version: 'e2e',
      CoreOptions: {
        ReadOnlyMode: false,
        EnableRBAC: rbac,
        GithubBadge: false,
      },
      AIOptions: {},
    }),
  )

  await page.route('**/rest-api/v1/authn**', route => {
    const auth = route.request().headers()['authorization']
    if (auth === `Bearer ${VALID_TOKEN}`) {
      return json(route, { success: true, data: 'OK' })
    }
    return json(
      route,
      { success: false, code: 401, message: 'unauthorized' },
      401,
    )
  })

  await page.route('**/rest-api/v1/search**', route =>
    json(route, {
      success: true,
      data: {
        items: searchItems,
        total: searchItems.length,
        currentPage: 1,
        pageSize: 20,
      },
    }),
  )
}

export function deployment(name: string, namespace = 'default') {
  return {
    cluster: 'e2e-cluster',
    deleted: false,
    object: {
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      metadata: {
        name,
        namespace,
        creationTimestamp: '2026-01-01T00:00:00Z',
      },
    },
  }
}
