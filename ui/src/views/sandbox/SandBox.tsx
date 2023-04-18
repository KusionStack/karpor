import { BrowserRouter, Route, Routes } from "react-router-dom";
import TopHeader from "../../components/TopHeader";
import Search from "./search/Search";
import Insight from "./insight/Insight";
import Cluster from "./cluster/Cluster";
import Result from "./result/Result";
import NotFound from "./notfound/NotFound";
import styles from "./style.module.scss";

import { Layout } from "antd";

const { Content } = Layout;

export default function SandBox() {
  return (
    <BrowserRouter>
      <Layout className="site-layout">
        <TopHeader></TopHeader>
        <Content className={styles.container}>
          <Routes>
            <Route index element={<Search />} />
            <Route path="/search" element={<Search />} />
            <Route path="/insight" element={<Insight />} />
            <Route path="/cluster" element={<Cluster />} />
            <Route path="/result" element={<Result />} />
            <Route path="*" element={<NotFound />} />
          </Routes>
        </Content>
      </Layout>
    </BrowserRouter>
  );
}
