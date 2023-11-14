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

import { useState, useEffect } from "react";
import { Tag } from "antd";
import axios from "axios";
import { useLocation } from "react-router-dom";
import queryString from "query-string";
import styles from "./styles.module.scss";


export default function Detail() {
  const location = useLocation();
  // const [pageData, setPageData] = useState<any>();
  const urlSearchParams = queryString.parse(location.search);
  const query = urlSearchParams?.query;

  const [searchParams, setSearchParams] = useState({
    pageSize: 10,
    page: 1,
  });

  // function handleChangePage(page: number, pageSize: number) {
  //   setSearchParams({
  //     ...searchParams,
  //     page,
  //     pageSize,
  //   });
  // }

  async function getPageData() {
    const data = await axios(`/apis/search.karbour.com/v1beta1/uniresources`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      params: {
        query: encodeURI(query as any),
      },
    });
    // setPageData(data || {});
  }

  useEffect(() => {
    getPageData();
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        <div className={styles.item}>
          <div className={styles.label}>Name</div>
          <div className={styles.val}>coredns-565d847f94-5rbmk</div>
        </div>
        <div className={styles.item}>
          <div className={styles.label}>NameSpace</div>
          <div className={styles.val}>kube-system</div>
        </div>
        <div className={styles.item}>
          <div className={styles.label}>Creation</div>
          <div className={styles.val}>2023-05-16T02:46:57Z</div>
        </div>
        <div className={styles.item}>
          <div className={styles.label}>Labels</div>
          <div className={styles.val}>
            {["k8s-app:kube-dns", "pod-template-hash:565d847f94"]?.map(
              (item, index) => (
                <span key={index} className={styles["label-item"]}>
                  {item}
                </span>
              )
            )}
          </div>
        </div>
        <div className={styles.item}>
          <div className={styles.label}>Controlled By</div>
          <div className={`${styles.val} ${styles.link}`}>
            RsplicaSet: coredns-565d847f94
          </div>
        </div>
        <div className={styles.item}>
          <div className={styles.label}>State</div>
          <div className={styles.val}>
            <span><Tag color="green">running</Tag></span>
          </div>
        </div>
        <div className={styles.item}>
          <div className={styles.label}>Node</div>
          <div className={`${styles.val} ${styles.link}`}>
            {"cluster1-control-plane"}
          </div>
        </div>
        <div className={styles.item}>
          <div className={styles.label}>Host IP</div>
          <div className={`${styles.val} ${styles.node}`}>{"172.18.0.2"}</div>
        </div>
        <div className={styles.item}>
          <div className={styles.label}>Pod IP</div>
          <div className={`${styles.val} ${styles.node}`}>{"10.244.0.3"}</div>
        </div>
      </div>
      {/* {(!pageData?.objects || !pageData?.objects?.length) && <Empty />} */}
    </div>
  );
}
