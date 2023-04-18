import {
  SearchOutlined,
  MacCommandOutlined,
  ApartmentOutlined,
  HeatMapOutlined
} from "@ant-design/icons";
import { Layout, Menu} from "antd";
import { useNavigate, useLocation } from "react-router-dom";

const { Header } = Layout;

export default function TopHeader() {
  const nav = useNavigate();
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
    nav(key);
  }

  return (
    <Header style={{display: 'flex', background: '#fff', padding: '0 20px'}}>
      <div className="logo">
        <HeatMapOutlined style={{marginRight: 10}}/>
      </div>
      <Menu
        mode="horizontal"
        selectedKeys={[location.pathname]}
        onClick={handleCkick}
        items={items}
      />
    </Header>
  );
}
