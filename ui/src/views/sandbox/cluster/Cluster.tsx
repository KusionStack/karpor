import { useState } from "react";
import { Pagination, Badge, Tooltip } from "antd";
import axios from "axios";
import { result, clusterList } from "../../../utils/mockData";
import styles from "./styles.module.scss";

export default function Cluster() {
  const [pageData, setPageData] = useState([]);
  const [searchParams, setSearchParams] = useState({
    pageSize: 10,
    page: 1,
  });
  async function getPageData() {
    let url = `/apis/search.karbour.com/v1beta1/search/proxy`;
    setPageData(clusterList?.items as any);
    const res = await axios(url, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      data: {},
    });
    console.log(res, "=====res====");
    setPageData(res?.data || result);
  }

  useState(() => {
    getPageData();
  });

  console.log(pageData, "===pageData===");
  function handleChangePage(page: number, pageSize: number) {
    console.log(page, pageSize, "handleChangePage");
    setSearchParams({
      ...searchParams,
      page,
      pageSize,
    });
  }

  function handleMore(item: any) {
    console.log(item, "====item====");
  }

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        {pageData?.map((item: any, index: number) => {
          return (
            <div className={styles.card} key={`${item.name}_${index}`}>
              <div className={styles.header}>
                <div className={styles.headerLeft}>
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
                  className={styles.headerRight}
                  onClick={() => handleMore(item)}
                >
                  More
                </div>
              </div>
              <div className={styles.cardBody}>
                <div className={styles.item}>
                  <div className={styles.itemLabel}>Endpoint: </div>
                  <Tooltip title={item.spec?.access?.endpoint}>
                    <div className={styles.itemValue}>
                      {item.spec?.access?.endpoint}
                    </div>
                  </Tooltip>
                </div>
                <div className={styles.stat}>
                  <div className={styles.node}>Nodes: {item?.status?.node}</div>
                  <div className={styles.deloy}>
                    Delay: {item?.status?.delay}
                  </div>
                </div>
              </div>
            </div>
          );
        })}
      </div>
      <div className={styles.footer}>
        <Pagination
          total={clusterList?.items?.length}
          showTotal={(total, range) => `${range[0]}-${range[1]} 共 ${total} 条`}
          pageSize={searchParams?.pageSize}
          current={searchParams?.page}
          onChange={handleChangePage}
        />
      </div>
    </div>
  );
}
