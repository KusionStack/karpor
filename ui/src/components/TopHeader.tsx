/*
 * Copyright The Karbour Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
