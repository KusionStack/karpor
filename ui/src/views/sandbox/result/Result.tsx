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

import { useState, useEffect, useCallback } from "react";
import { Input, Select, Space, Pagination, Empty } from "antd";
import axios from "axios";
import { useLocation } from "react-router-dom";
import { CaretRightOutlined, RadarChartOutlined } from "@ant-design/icons";
import queryString from "query-string";
import Tabs from "../../../components/Tabs";
import { basicSyntaxColumns } from "../../../utils/constants";
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

const tabsList = [
  { label: "Code search basics", value: "code" },
  { label: "Search query examples", value: "examples" },
];

export default function Result() {
  const location = useLocation();
  const [showPanel, setShowPanel] = useState<Boolean>(false);
  const [pageData, setPageData] = useState<any>();
  const urlSearchParams = queryString.parse(location.search);
  const [searchValue, setSearchValue] = useState(
    urlSearchParams?.keyword || ""
  );

  const [currentTab, setCurrentTab] = useState<String>(tabsList?.[0].value);

  const handleTabChange = (value: String) => {
    setCurrentTab(value);
  };

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
        keyword: searchValue,
      },
    });
    setPageData(data || {});
  }

  // const getPageData = useCallback(async () => {
  //   let url = `/apis/search.karbour.com/v1beta1/search/proxy`;
  //   const res = await axios(url, {
  //     method: 'GET',
  //     headers: {
  //       'Content-Type': 'application/json',
  //     },
  //     data: {
  //       searchValue
  //     }
  //   })
  //   setPageData(res?.data || {});
  // }, [])
  // useEffect(() => {
  //   getPageData();
  // }, [getPageData])

  useEffect(() => {
    getPageData();
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const handleChange = (event: any) => {
    setSearchValue(event.target.value);
  };

  const handleSearch = (val: any) => {
    getPageData();
  };

  function renderHight(str: string) {
    if (!str) return;
    const strList = str.split("\n").join(" ").split(" ");
    let list = strList?.map((item, index) => {
      let tmp;
      let itemTmp = item.trim();
      if (itemTmp.indexOf(":") > -1) {
        const tmp1 = itemTmp.split(":");
        tmp = (
          <span key={index} className={styles.txtSpan}>
            <span style={{ color: "#1890ff" }}>{tmp1[0]}</span>:{tmp1[1]}
          </span>
        );
      } else if (itemTmp === "OR" || itemTmp === "AND" || itemTmp === "NOT") {
        tmp = (
          <span
            key={index}
            className={styles.txtSpan}
            style={{ color: "#a305e1", margin: "0 5px" }}
          >
            {itemTmp}
          </span>
        );
      } else {
        tmp = (
          <span key={index} className={styles.txtSpan}>
            {itemTmp}
          </span>
        );
      }
      return tmp;
    });
    return list;
  }

  const addItem = (queryExample: string) => {
    const val = searchValue ? `${searchValue} ${queryExample}` : queryExample;
    setSearchValue(val);
  };

  return (
    <div className={styles.container}>
      <div style={{ width: 850, position: "relative" }}>
        <Space.Compact>
          <Select
            defaultValue="filter"
            options={options}
            style={{ width: 150 }}
          />
          <Search
            placeholder="input search text"
            allowClear
            style={{ width: 700 }}
            onFocus={() => setShowPanel(true)}
            onBlur={() => setShowPanel(false)}
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
        {pageData?.objects?.map((item: any, index: number) => {
          return (
            <div className={styles.card} key={`${item.name}_${index}`}>
              <div>
                <RadarChartOutlined style={{ margin: "0 5px" }} />
              </div>
              <div className={styles.item}>
                <span className={styles.itemLabel}>ApiVersion</span>
                <span className={styles.itemValue}>{item?.apiVersion}</span>
              </div>
              <CaretRightOutlined style={{ margin: "0 3px", fontSize: 14 }} />
              <div className={styles.item}>
                <span className={styles.itemLabel}>Kind</span>
                <span className={styles.itemValue}>{item?.kind}</span>
              </div>
              {item?.metadata?.namespace && (
                <>
                  <CaretRightOutlined
                    style={{ margin: "0 3px", fontSize: 14 }}
                  />
                  <div className={styles.item}>
                    <span className={styles.itemLabel}>NameSpace</span>
                    <span className={styles.itemValue}>
                      {item?.metadata?.namespace}
                    </span>
                  </div>
                </>
              )}
              <CaretRightOutlined style={{ margin: "0 3px", fontSize: 14 }} />
              <div className={styles.item}>
                <span className={styles.itemLabel}>Name</span>
                <span className={styles.itemValue}>{item?.metadata?.name}</span>
              </div>
            </div>
          );
        })}
        {/* <div style={{ width: 600, height: 600 }}>图片占位</div> */}
      </div>
      {
        pageData?.objects && pageData?.objects?.length > 0 && <div className={styles.footer}>
          <Pagination
            total={pageData?.objects?.length}
            showTotal={(total: number, range: any[]) =>
              `${range[0]}-${range[1]} 共 ${total} 条`
            }
            pageSize={searchParams?.pageSize}
            current={searchParams?.page}
            onChange={handleChangePage}
          />
        </div>
      }
      {
        (!pageData?.objects || !pageData?.objects?.length) && <Empty />
      }
    </div>
  );
}
