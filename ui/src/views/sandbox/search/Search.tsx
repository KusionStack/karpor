import React, { useState } from "react";
import { Input, Select, Space } from "antd";
import {useNavigate} from "react-router-dom";
import { HeatMapOutlined } from "@ant-design/icons";
import { basicSyntaxColumns } from "../../../utils/constants";
import KarbourTabs from "../../../components/Tabs/index";

import styles from "./styles.module.scss";

const { Search } = Input;

const options = [
  {
    value: "filter",
    label: "Filter(s)",
  },
  {
    value: "AI",
    label: "AI Suggestion",
  },
  {
    value: "SQL",
    label: "SQL",
  },
];

const SearchPage = () => {
  const navigate = useNavigate();
  const tabsList = [
    { label: "Code search basics", value: "code" },
    { label: "Search query examples", value: "examples" },
  ];
  const [currentTab, setCurrentTab] = useState<string>(tabsList?.[0].value);
  const [searchType, setSearchType] = useState<string>("filter");

  const [inputValue, setInputValue] = useState("");

  const handleTabChange = (value: string) => {
    setCurrentTab(value);
  };

  // const divInputRef = useRef();

  const addItem = (queryExample: string) => {
    const val = inputValue ? `${inputValue} ${queryExample}` : queryExample;
    setInputValue(val);
  };

  const renderItem = (query: string) => {
    return renderHight(query);
  };

  const renderFilter = () => {
    console.log(currentTab, "===currentTab===")
    return (
      <>
        <KarbourTabs
          list={tabsList}
          current={currentTab}
          onChange={handleTabChange}
        />
        <div className={styles.box}>
          {basicSyntaxColumns?.map((item: any, index) => {
            return (
              <div className={styles.basicSyntax} key={`${index + 1}`}>
                {item?.map((queryExample: any, inIndex: number) => {
                  return (
                    <div
                      className={styles.itemBox}
                      key={`${queryExample.title}_${inIndex}`}
                    >
                      <div className={styles.title}>{queryExample.title}</div>
                      {queryExample.queryExamples.map(
                        (example: any, innerIdx: number) => {
                          return (
                            <div
                              className={styles.childrenItemBox}
                              key={`${example.id}_${innerIdx}`}
                              onClick={() => addItem(example.query)}
                            >
                              <span className={styles.childrenItem}>
                                {renderItem(example.query)}
                              </span>
                              {/* <FilterValueOption option={{ label: example.query, matches: false }} /> */}
                            </div>
                          );
                        }
                      )}
                    </div>
                  );
                })}
              </div>
            );
          })}
        </div>
      </>
    );
  };

  function renderHight(str: string) {
    if (!str) return;
    const strList = str.split("\n").join(" ").split(" ");
    console.log(strList, "===strList===");
    let list = strList?.map((item, index) => {
      let tmp;
      let itemTmp = item.trim();
      if (itemTmp.indexOf(":") > -1) {
        const tmp1 = itemTmp.split(":");
        tmp = (
          <span className={styles.txtSpan} key={index + 1}>
            <span style={{ color: "#1890ff" }}>{tmp1[0]}</span>:{tmp1[1]}
          </span>
        );
      } else if (itemTmp === "OR" || itemTmp === "AND" || itemTmp === "NOT") {
        tmp = (
          <span
            key={index + 1}
            className={styles.txtSpan}
            style={{ color: "#a305e1", margin: "0 5px" }}
          >
            {itemTmp}
          </span>
        );
      } else {
        tmp = (
          <span key={index + 1} className={styles.txtSpan}>
            {itemTmp}
          </span>
        );
      }
      return tmp;
    });
    return list;
  }

  const handleSearch = (val: string) => {
    navigate(`/result?keyword=${val}`);
  };

  const handleInputChange = (event: any) => {
    setInputValue(event.target.value);
  };

  return (
    <div
      style={{
        paddingTop: "100px",
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <div className={styles.logoBox}>
        <HeatMapOutlined />
        <div className={styles.title}>Karbour</div>
      </div>

      <Space.Compact>
        <Select
          defaultValue="filter"
          options={options}
          style={{ width: 150 }}
          onChange={(val) => setSearchType(val)}
          value={searchType}
        />
        <Search
          placeholder="input search text"
          allowClear
          style={{ width: 600 }}
          value={inputValue}
          onChange={handleInputChange}
          onSearch={handleSearch}
        />
      </Space.Compact>
      {/* {
        renderAI()
      } */}
      <div className={styles.content}>
        {searchType === "filter" && renderFilter()}
      </div>
    </div>
  );
};

export default SearchPage;
