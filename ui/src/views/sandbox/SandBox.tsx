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
