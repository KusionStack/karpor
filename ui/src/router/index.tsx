import React, { lazy, ReactNode, Suspense } from 'react'
import { Navigate, Outlet, useRoutes } from 'react-router-dom'
import {
  SearchOutlined,
  MacCommandOutlined,
  ApartmentOutlined,
} from '@ant-design/icons'
import Layout from '@/components/layout'
import Loading from '@/components/loading'

const Search = lazy(() => import('@/pages/search'))
const Result = lazy(() => import('@/pages/result'))
const Cluster = lazy(() => import('@/pages/cluster'))
const ClusterAdd = lazy(() => import('@/pages/cluster/add'))
const ClusterCertificate = lazy(() => import('@/pages/cluster/certificate'))
const InsightDetailCluster = lazy(() => import('@/pages/insightDetail/cluster'))
const InsightDetailKind = lazy(() => import('@/pages/insightDetail/kind'))
const InsightDetailNamespace = lazy(
  () => import('@/pages/insightDetail/namespace'),
)
const InsightDetailResource = lazy(
  () => import('@/pages/insightDetail/resource'),
)
const Reflux = lazy(() => import('@/pages/reflux'))
const Insight = lazy(() => import('@/pages/insight'))
const NotFound = lazy(() => import('@/pages/notfound'))

const lazyLoad = (children: ReactNode): ReactNode => {
  return <Suspense fallback={<Loading />}>{children}</Suspense>
}

export interface RouteObject {
  key?: string
  path?: string
  title?: string
  icon?: React.ReactNode
  element: React.ReactNode
  children?: RouteObject[]
  index?: any
}

const router: RouteObject[] = [
  {
    path: '/',
    element: <Layout />,
    children: [
      {
        key: '/search',
        path: '/search',
        element: (
          <>
            <Outlet />
          </>
        ),
        icon: <SearchOutlined />,
        children: [
          {
            index: true,
            title: '搜索',
            element: lazyLoad(<Search />),
          },
          {
            key: 'result',
            path: 'result',
            title: '结果',
            element: lazyLoad(<Result />),
          },
        ],
      },
      {
        key: '/reflux',
        path: '/reflux',
        title: 'reflux',
        element: lazyLoad(<Reflux />),
        icon: <SearchOutlined />,
      },
      {
        key: '/insight',
        path: '/insight',
        title: 'Insight',
        element: lazyLoad(<Insight />),
        icon: <MacCommandOutlined />,
      },
      {
        key: 'insightDetail',
        path: 'insightDetail',
        element: (
          <>
            <Outlet />
          </>
        ),
        // element: lazyLoad(<InsightDetail />),
        children: [
          {
            key: 'cluster',
            path: 'cluster',
            element: lazyLoad(<InsightDetailCluster />),
          },
          {
            key: 'kind',
            path: 'kind',
            element: lazyLoad(<InsightDetailKind />),
          },
          {
            key: 'namespace',
            path: 'namespace',
            element: lazyLoad(<InsightDetailNamespace />),
          },
          {
            key: 'resource',
            path: 'resource',
            element: lazyLoad(<InsightDetailResource />),
          },
        ],
      },
      {
        key: '/cluster',
        path: '/cluster',
        element: (
          <>
            <Outlet />
          </>
        ),
        icon: <ApartmentOutlined />,
        children: [
          {
            index: true,
            title: '集群列表',
            element: lazyLoad(<Cluster />),
          },
          {
            key: 'access',
            path: 'access',
            title: '集群接入',
            element: lazyLoad(<ClusterAdd />),
          },
          {
            key: 'certificate',
            path: 'certificate',
            title: '更新证书',
            element: lazyLoad(<ClusterCertificate />),
          },
        ],
      },
      {
        path: '/',
        title: '',
        element: <Navigate to="/search" replace />,
      },
      {
        path: '*',
        title: '',
        element: <NotFound />,
      },
    ],
  },
]

const WrappedRoutes = () => {
  return useRoutes(router)
}

export default WrappedRoutes
