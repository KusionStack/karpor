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

import { useState } from "react";
import { Pagination, Badge, Tooltip, Empty } from "antd";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import styles from "./styles.module.scss";

export default function Cluster() {
  const navigate = useNavigate();
  const [pageData, setPageData] = useState<any>([]);
  const [searchParams, setSearchParams] = useState({
    pageSize: 10,
    page: 1,
  });
  async function getPageData() {
    const data = await axios(`/apis/cluster.karbour.com/v1beta1/clusters`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      params: {},
    });
    setPageData(data || {});
  }

  useState(() => {
    getPageData();
  });

  function handleChangePage(page: number, pageSize: number) {
    setSearchParams({
      ...searchParams,
      page,
      pageSize,
    });
  }

  function handleMore(item: any) {
    console.log(item, "====handleMore====");
  }

  const handleClick = (item) => {
    console.log(item, "===item===")
    let queryStr = "";
    if (item?.metadata?.managedFields?.[0]?.apiVersion) {
      queryStr = `${item?.metadata?.managedFields?.[0]?.apiVersion},${item.metadata?.name}`;
    } else {
      queryStr = `${item.metadata?.name}`;
    }
    navigate(`/cluster-detail?query=${queryStr}`);
  };

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        {pageData?.items?.map((item: any, index: number) => {
          return (
            <div className={styles.card} key={`${item.name}_${index}`} onClick={() => handleClick(item)}>
              <div className={styles.header}>
                <div className={styles['header-left']}>
                  {item.metadata?.name}
                  <Badge
                    style={{
                      marginLeft: 20,
                      fontWeight: "normal",
                      color: item?.status?.healthy === "true" ? "green" : "red",
                    }}
                    status={
                      item?.status?.healthy === "true" ? "success" : "error"
                    }
                    text={item?.status?.healthy === "true" ? "健康" : "不健康"}
                  />
                </div>
                <div
                  className={styles['header-right']}
                  onClick={() => handleMore(item)}
                >
                  More
                </div>
              </div>
              <div className={styles['card-body']}>
                <div className={styles.item}>
                  <div className={styles['item-label']}>Endpoint: </div>
                  <Tooltip title={item.spec?.access?.endpoint}>
                    <div className={styles['item-value']}>
                      {item.spec?.access?.endpoint}
                    </div>
                  </Tooltip>
                </div>
                <div className={styles.stat}>
                  <div className={styles.node}>
                    Nodes: {item?.status?.node || "--"}
                  </div>
                  <div className={styles.deloy}>
                    Delay: {item?.status?.delay || "--"}
                  </div>
                </div>
              </div>
            </div>
          );
        })}
      </div>
      {
        (pageData?.items && pageData?.items?.length > 0) &&
        <div className={styles.footer}>
          <Pagination
            total={pageData?.items?.length}
            showTotal={(total, range) => `${range[0]}-${range[1]} 共 ${total} 条`}
            pageSize={searchParams?.pageSize}
            current={searchParams?.page}
            onChange={handleChangePage}
          />
        </div>
      }
      {
        (!pageData?.items || !pageData?.items?.length) && <Empty />
      }
    </div>
  );
}
