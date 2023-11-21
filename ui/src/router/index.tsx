import { lazy, ReactNode, Suspense } from "react";
import { Navigate } from "react-router-dom";
import {
  SearchOutlined,
  MacCommandOutlined,
  ApartmentOutlined,
} from "@ant-design/icons";
import { Spin } from "antd";
import Layout from "../components/Layout";

const Search = lazy(() => import("../views/search/Search"));
const Result = lazy(() => import("../views/result/Result"));
const Cluster = lazy(() => import("../views/cluster/Cluster"));
const ClusterDetail = lazy(
  () => import("../views/cluster-detail/ClusterDetail")
);
const Insight = lazy(() => import("../views/insight/Insight"));
const Detail = lazy(() => import("../views/detail/Detail"));
const NotFound = lazy(() => import("../views/notfound/NotFound"));
const InsightGrid = lazy(() => import("../views/insightGrid/InsightGrid"));

const lazyLoad = (children: ReactNode): ReactNode => {
  return <Suspense fallback={<Spin />}>{children}</Suspense>;
};

export interface RouteObject {
  key?: string;
  path: string;
  title?: string;
  icon?: React.ReactNode;
  element: React.ReactNode;
  children?: RouteObject[];
  istopmenu?: boolean;
}

const router: RouteObject[] = [
  {
    path: "/",
    element: <Layout />,
    children: [
      {
        key: "/search",
        path: "/search",
        title: "search",
        element: lazyLoad(<Search />),
        icon: <SearchOutlined />,
        istopmenu: true,
      },
      {
        key: "/result",
        path: "/result",
        title: "",
        element: lazyLoad(<Result />),
        // icon: <DesktopOutlined />,
        istopmenu: false,
      },
      {
        key: "/detail",
        path: "/detail",
        title: "",
        element: lazyLoad(<Detail />),
        // icon: <SendOutlined />,
        istopmenu: false,
      },
      {
        key: "/insight",
        path: "/insight",
        title: "Insight",
        element: lazyLoad(<Insight />),
        icon: <MacCommandOutlined />,
        istopmenu: true,
      },
      {
        key: "/cluster",
        path: "/cluster",
        title: "Cluster",
        element: lazyLoad(<Cluster />),
        icon: <ApartmentOutlined />,
        istopmenu: true,
      },
      {
        key: "/clusterDetail",
        path: "/clusterDetail",
        title: "",
        element: lazyLoad(<ClusterDetail />),
        // icon: <SettingOutlined />,
        istopmenu: false,
      },
      {
        key: "/insightGrid",
        path: "/insightGrid",
        title: "",
        element: lazyLoad(<InsightGrid />),
        // icon: <SettingOutlined />,
        istopmenu: false,
      },
      {
        path: "/",
        title: "",
        element: <Navigate to="/search" replace />,
        istopmenu: false,
      },
      {
        path: "*",
        title: "",
        element: <NotFound />,
        istopmenu: false,
      },
    ],
  },
];

export default router;
