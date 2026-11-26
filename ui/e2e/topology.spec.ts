import { expect, test } from '@playwright/test'
import type { Page } from '@playwright/test'
import { mockApi } from './mock-api'

const CLUSTER = 'e2e-cluster'

async function mockRoute(page: Page, url: string, data: unknown) {
  await page.route(url, route =>
    route.fulfill({
      contentType: 'application/json',
      body: JSON.stringify({ success: true, data }),
    }),
  )
}

// Insight pages redirect to the login page until the UI knows RBAC is
// disabled, so let a public page persist that flag before opening one.
async function openInsightPage(page: Page, url: string) {
  await page.goto('/search')
  await page.waitForFunction(() => {
    const root = JSON.parse(window.localStorage.getItem('persist:root') || '{}')
    return JSON.parse(root.globalSlice || '{}').isUnsafeMode === true
  })
  await page.goto(url)
}

// The graph is drawn on a canvas and fitted to its center, so hovering the
// center of the canvas hovers the only node of a single-node topology.
async function hoverTopologyNode(page: Page) {
  const canvas = page.locator('canvas').first()
  await expect(canvas).toBeVisible()
  const tooltip = page.getByRole('tooltip')
  await expect(async () => {
    await page.mouse.move(0, 0)
    await canvas.hover()
    await expect(tooltip).toBeVisible({ timeout: 1000 })
  }).toPass()
  return tooltip
}

test.describe('topology', () => {
  test('shows the kind details of a node in the cluster topology', async ({
    page,
  }) => {
    await mockApi(page)
    await mockRoute(page, '**/rest-api/v1/insight/topology**', {
      [CLUSTER]: {
        'apps.v1.Deployment': {
          resourceGroup: {
            cluster: CLUSTER,
            apiVersion: 'apps/v1',
            kind: 'Deployment',
          },
          count: 12,
          relationship: {},
        },
      },
    })

    await openInsightPage(
      page,
      `/insightDetail/cluster?from=insight&cluster=${CLUSTER}&type=cluster`,
    )
    const tooltip = await hoverTopologyNode(page)

    await expect(tooltip.getByText('Deployment', { exact: true })).toBeVisible()
    await expect(tooltip.getByText('apps/v1', { exact: true })).toBeVisible()
    await expect(tooltip.getByText(CLUSTER, { exact: true })).toBeVisible()
    await expect(tooltip.getByText('12', { exact: true })).toBeVisible()
    await expect(tooltip.getByText('Namespace')).toHaveCount(0)
  })

  test('shows the live status of a pod in the resource topology', async ({
    page,
  }) => {
    await mockApi(page)
    const resource = {
      cluster: CLUSTER,
      apiVersion: 'v1',
      kind: 'Pod',
      namespace: 'default',
      name: 'web-0',
    }
    await mockRoute(page, '**/rest-api/v1/insight/topology**', {
      [CLUSTER]: {
        '/v1.Pod:default.web-0': {
          resourceGroup: resource,
          parents: [],
          children: [],
        },
      },
    })
    await mockRoute(page, '**/rest-api/v1/insight/summary**', {
      resource: { ...resource, status: 'CrashLoopBackOff' },
    })

    await openInsightPage(
      page,
      `/insightDetail/resource?from=insight&cluster=${CLUSTER}&apiVersion=v1&kind=Pod&namespace=default&name=web-0&type=resource`,
    )
    const tooltip = await hoverTopologyNode(page)

    await expect(tooltip.getByText('web-0', { exact: true })).toBeVisible()
    await expect(tooltip.getByText('Pod · v1', { exact: true })).toBeVisible()
    await expect(tooltip.getByText(CLUSTER, { exact: true })).toBeVisible()
    await expect(tooltip.getByText('default', { exact: true })).toBeVisible()
    await expect(
      tooltip.getByText('CrashLoopBackOff', { exact: true }),
    ).toBeVisible()
  })
})
