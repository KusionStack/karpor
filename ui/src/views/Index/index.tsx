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

import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Layout } from "antd";
import TopHeader from "../../components/TopHeader";
import Search from "../search/Search";
import Insight from "../insight/Insight";
import Cluster from "../cluster/Cluster";
import Result from "../result/Result";
import Detail from "../detail/Detail";
import ClusterDetail from "../cluster-detail/ClusterDetail";
import NotFound from "../notfound/NotFound";

import styles from "./styles.module.scss";

const { Content } = Layout;

const SandBox = () => {
  return (
    <BrowserRouter>
      <Layout>
        <TopHeader></TopHeader>
        <Content className={styles.container}>
          <Routes>
            <Route index element={<Search />} />
            <Route path="/search" element={<Search />} />
            <Route path="/insight" element={<Insight />} />
            <Route path="/cluster" element={<Cluster />} />
            <Route path="/result" element={<Result />} />
            <Route path="/detail" element={<Detail />} />
            <Route path="/cluster-detail" element={<ClusterDetail />} />
            <Route path="*" element={<NotFound />} />
          </Routes>
        </Content>
      </Layout>
    </BrowserRouter>
  );
}
export default SandBox

