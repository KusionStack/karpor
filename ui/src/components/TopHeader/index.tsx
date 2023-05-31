import {
  SearchOutlined,
  MacCommandOutlined,
  ApartmentOutlined,
  HeatMapOutlined
} from "@ant-design/icons";
import { Layout, Menu} from "antd";
import { useNavigate, useLocation } from "react-router-dom";
import styles from "./style.module.scss";

const { Header } = Layout;

export default function TopHeader() {
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

  function handleCkick({ key }: {key: string}) {
    navigate(key);
  }

  const selectedKeys = (!location.pathname || location.pathname === '/') ? '/search' : location.pathname;

  return (
    <Header className={styles.container}>
      <div className={styles.logo}>
        <HeatMapOutlined/>
      </div>
      <Menu
        mode="horizontal"
        defaultSelectedKeys={['/search']}
        selectedKeys={[selectedKeys]}
        onClick={handleCkick}
        items={items}
      />
    </Header>
  );
}
