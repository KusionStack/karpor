import {
  SearchOutlined,
  MacCommandOutlined,
  ApartmentOutlined,
  HeatMapOutlined,
} from "@ant-design/icons";
import { Layout, Menu } from "antd";
import { memo } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import styles from "./style.module.less";

const { Header } = Layout;

const TopHeader = () => {
  const navigate = useNavigate();
  const location = useLocation();

  const items = [
    {
      key: "/search",
      label: "Search",
      icon: <SearchOutlined />,
    },
    {
      key: "/insight",
      label: "Insight",
      icon: <MacCommandOutlined />,
    },
    {
      key: "/cluster",
      label: "Cluster",
      icon: <ApartmentOutlined />,
    },
  ];

  function handleCkick({ key }: { key: string }) {
    navigate(key);
  }

  const selectedKeys =
    !location.pathname || location.pathname === "/"
      ? "/search"
      : location.pathname;

  return (
    <Header className={styles["top-container"]}>
      <div className={styles.logo}>
        <HeatMapOutlined />
      </div>
      <Menu
        mode="horizontal"
        defaultSelectedKeys={["/search"]}
        selectedKeys={[selectedKeys]}
        onClick={handleCkick}
        items={items}
      />
    </Header>
  );
};

export default memo(TopHeader);
