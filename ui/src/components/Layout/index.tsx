import React, { memo, useEffect, useMemo, useState } from "react";
import { Layout, Menu, MenuProps } from "antd";
import { HeatMapOutlined } from "@ant-design/icons";
import { matchRoutes, Outlet, useLocation, NavLink } from "react-router-dom";
import router, { RouteObject } from "../../router";
import styles from "./style.module.less";

const { Content, Header } = Layout;

type MenuItem = Required<MenuProps>["items"][number];

function getItem(
  label: React.ReactNode,
  key: React.Key,
  icon?: React.ReactNode,
  children?: MenuItem[],
  type?: "group"
): MenuItem {
  return {
    key,
    icon,
    children,
    label,
    type,
  } as MenuItem;
}

const LayoutPage = () => {
  const [defaultSelectedKeys, setSelectKey] = useState<string[]>(["1"]);
  const [defaultOpenKeys, setDefaultOpenKeys] = useState<string[]>(["1"]);
  const [initial, setInitial] = useState(false);
  const location = useLocation();

  const items: MenuItem[] = useMemo(() => {
    let menuItems = [];
    router &&
      (router[0]?.children || []).forEach((item: RouteObject) => {
        // console.log(item, "===item===");
        const navItem = getItem(
          <NavLink to={item.path as string}>{item.title}</NavLink>,
          item.path as string,
          item.icon
        );
        if (item.istopmenu) {
          menuItems.push(navItem);
        }
      });
    return menuItems;
  }, []);

  useEffect(() => {
    const routes = matchRoutes(router, { pathname: location.pathname });
    // console.log('routes', routes, location.pathname);
    const pathArr: string[] = [];
    if (routes && routes.length) {
      routes.forEach((item) => {
        const path = item.route.path;
        if (path === location.pathname) {
          pathArr.push(path);
        }
      });
    }
    setSelectKey(pathArr);
    setDefaultOpenKeys(pathArr);
    setInitial(true);
  }, [location]);

  if (!initial) {
    return null;
  }

  return (
    <Layout>
      <Header className={styles["top-container"]}>
        <div className={styles.logo}>
          <HeatMapOutlined />
        </div>
        <Menu
          style={{ flex: 1, border: "none" }}
          mode="horizontal"
          defaultSelectedKeys={defaultSelectedKeys}
          defaultOpenKeys={defaultOpenKeys}
          selectedKeys={defaultSelectedKeys}
          items={items}
        />
      </Header>
      <Content className={styles.container}>
        <Outlet />
      </Content>
    </Layout>
  );
};

export default memo(LayoutPage);
