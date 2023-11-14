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
import { Input, Select, Space, Pagination, Empty } from "antd";
import axios from "axios";
import { useLocation, useNavigate } from "react-router-dom";
import { CaretRightOutlined, RadarChartOutlined } from "@ant-design/icons";
import queryString from "query-string";
import styles from "./styles.module.scss";

const { Search } = Input;

const options = [
  {
    value: "filter",
    label: "Filter(s)",
  },
  {
    value: "AI Suggestion",
    label: "AI Suggestion",
  },
  {
    value: "SQL",
    label: "SQL",
  },
];

// const tabsList = [
//   { label: "Code search basics", value: "code" },
//   { label: "Search query examples", value: "examples" },
// ];

export default function Result() {
  const location = useLocation();
  const navigate = useNavigate();
  // const [showPanel, setShowPanel] = useState<Boolean>(false);
  const [pageData, setPageData] = useState<any>();
  const urlSearchParams = queryString.parse(location.search);
  const [searchValue, setSearchValue] = useState(urlSearchParams?.query || "");
  const [searchType, setSearchType] = useState<string>(urlSearchParams?.patternType === 'SQL' ? "SQL" : "filter");

  // const [currentTab, setCurrentTab] = useState<String>(tabsList?.[0].value);

  // const handleTabChange = (value: String) => {
  //   setCurrentTab(value);
  // };

  const [searchParams, setSearchParams] = useState({
    pageSize: 10,
    page: 1,
  });

  function handleChangePage(page: number, pageSize: number) {
    setSearchParams({
      ...searchParams,
      page,
      pageSize,
    });
  }

  async function getPageData() {
    const data = await axios(`/apis/search.karbour.com/v1beta1/uniresources`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      params: {
        query: encodeURI(searchValue as any),
        ...(searchType === 'SQL' ? { patternType: "sql", } : {})
      },
    });
    setPageData(data || {});
  }

  useEffect(() => {
    getPageData();
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const handleChange = (event: any) => {
    setSearchValue(event.target.value);
  };

  const handleSearch = (val: any) => {
    getPageData();
  };

  // function renderHight(str: string) {
  //   if (!str) return;
  //   const strList = str.split("\n").join(" ").split(" ");
  //   let list = strList?.map((item, index) => {
  //     let tmp;
  //     let itemTmp = item.trim();
  //     if (itemTmp.indexOf(":") > -1) {
  //       const tmp1 = itemTmp.split(":");
  //       tmp = (
  //         <span key={index} className={styles.txtSpan}>
  //           <span style={{ color: "#1890ff" }}>{tmp1[0]}</span>:{tmp1[1]}
  //         </span>
  //       );
  //     } else if (itemTmp === "OR" || itemTmp === "AND" || itemTmp === "NOT") {
  //       tmp = (
  //         <span
  //           key={index}
  //           className={styles.txtSpan}
  //           style={{ color: "#a305e1", margin: "0 5px" }}
  //         >
  //           {itemTmp}
  //         </span>
  //       );
  //     } else {
  //       tmp = (
  //         <span key={index} className={styles.txtSpan}>
  //           {itemTmp}
  //         </span>
  //       );
  //     }
  //     return tmp;
  //   });
  //   return list;
  // }

  // const addItem = (queryExample: string) => {
  //   const val = searchValue ? `${searchValue} ${queryExample}` : queryExample;
  //   setSearchValue(val);
  // };

  const handleClick = (item) => {
    let queryStr = "";
    if (item?.metadata?.namespace) {
      queryStr = `${item?.apiVersion},${item?.kind},${item?.metadata?.namespace},${item?.metadata?.name}`;
    } else {
      queryStr = `${item?.apiVersion},${item?.kind},${item?.metadata?.name}`;
    }
    navigate(`/insight?query=${queryStr}`);
  };

  function handleChangeType(val) {
    setSearchType(val);
  }

  return (
    <div className={styles.container}>
      <div style={{ width: 850, position: "relative" }}>
        <Space.Compact>
          <Select
            defaultValue="filter"
            value={searchType}
            options={options}
            style={{ width: 150 }}
            onChange={handleChangeType}
          />
          <Search
            placeholder="input search text"
            allowClear
            style={{ width: 700 }}
            // onFocus={() => setShowPanel(true)}
            // onBlur={() => setShowPanel(false)}
            value={searchValue as any}
            onChange={handleChange}
            onSearch={handleSearch}
          />
        </Space.Compact>
        {/* {
          showPanel &&  <div className={styles.panel} onMouseMove={() => setShowPanel(true)}>
            {renderFilter()}
          </div>
        } */}
      </div>
      <div className={styles.content}>
        {pageData?.items?.map((item: any, index: number) => {
          return (
            <div
              className={styles.card}
              key={`${item.name}_${index}`}
              onClick={() => handleClick(item)}
            >
              <div>
                <RadarChartOutlined style={{ margin: "0 5px" }} />
              </div>
              <div className={styles.item}>
                <span className={styles["item-label"]}>ApiVersion</span>
                <span className={styles["item-value"]}>{item?.apiVersion}</span>
              </div>
              <CaretRightOutlined style={{ margin: "0 3px", fontSize: 14 }} />
              <div className={styles.item}>
                <span className={styles["item-label"]}>Kind</span>
                <span className={styles["item-value"]}>{item?.kind}</span>
              </div>
              {item?.metadata?.namespace && (
                <>
                  <CaretRightOutlined
                    style={{ margin: "0 3px", fontSize: 14 }}
                  />
                  <div className={styles.item}>
                    <span className={styles["item-label"]}>NameSpace</span>
                    <span className={styles["item-value"]}>
                      {item?.metadata?.namespace}
                    </span>
                  </div>
                </>
              )}
              <CaretRightOutlined style={{ margin: "0 3px", fontSize: 14 }} />
              <div className={styles.item}>
                <span className={styles["item-label"]}>Name</span>
                <span className={styles["item-value"]}>
                  {item?.metadata?.name}
                </span>
              </div>
            </div>
          );
        })}
        {/* <div style={{ width: 600, height: 600 }}>图片占位</div> */}
      </div>
      {pageData?.items && pageData?.items?.length > 0 && (
        <div className={styles.footer}>
          <Pagination
            total={pageData?.items?.length}
            showTotal={(total: number, range: any[]) =>
              `${range[0]}-${range[1]} 共 ${total} 条`
            }
            pageSize={searchParams?.pageSize}
            current={searchParams?.page}
            onChange={handleChangePage}
          />
        </div>
      )}
      {(!pageData?.items || !pageData?.items?.length) && <Empty />}
    </div>
  );
}
