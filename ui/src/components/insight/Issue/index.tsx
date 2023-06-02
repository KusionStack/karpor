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

import { Table, Tag, Input, Button, Tooltip, Badge } from "antd";
import { useState } from "react";
import styles from "./styles.module.scss";

type IProps = {
  data: any[];
  handleSearch: (val: string) => void;
};

const EnhancerTable = (props: IProps) => {
  const [searchValue, setSearchValue] = useState("");

  const columns: any[] = [
    {
      title: "TITLE",
      dataIndex: "title",
      key: "title",
    },
    {
      title: "SEVERITY",
      dataIndex: "severity",
      key: "severity",
      defaultSortOrder: "descend",
      sorter: (a: any, b: any) => a.severity.length - b.severity.length,
      render: (text: string) => {
        const levelMap = {
          low: "#3b5999",
          medium: "#F4A460",
          high: "#cd201f",
        };
        let color = levelMap[text.toLowerCase() as "low" | "medium" | "high"];
        return (
          <Tag color={color} key={text}>
            {text.toUpperCase()}
          </Tag>
        );
      },
    },
    {
      title: "LABELS",
      dataIndex: "labels",
      key: "labels",
      render: (text: string) => {
        const levelMap = {
          CVE: "#2db7f5",
          CIS: "#87d068",
          KUBEAUDIT: "#108ee9",
        };
        const list = text?.split(",");
        return list?.map((item, index) => {
          let color =
            levelMap[item.toUpperCase() as "CVE" | "CIS" | "KUBEAUDIT"];
          if (item.toUpperCase() === "SOLUTION") {
            return (
              <Tooltip
                placement="topLeft"
                title={"SOLUTION"}
                key={`${item}_${index}`}
              >
                <Tag color={color} key={item} style={{ fontSize: 12 }}>
                  {item.toUpperCase()}
                </Tag>
              </Tooltip>
            );
          }
          return (
            <Tag
              color={color}
              key={`${item}_${index}`}
              style={{ fontSize: 12 }}
            >
              {item.toUpperCase()}
            </Tag>
          );
        });
      },
    },
  ];

  function handleChange(event: any) {
    setSearchValue(event.target.value);
  }

  function handleSearch() {
    props?.handleSearch(searchValue);
  }

  return (
    <div className={styles.enhancer}>
      <div className={styles.search}>
        <Input
          value={searchValue}
          onChange={handleChange}
          style={{ width: 300, marginRight: 10 }}
        />
        <Button type="primary" onClick={handleSearch}>
          Search
        </Button>

        <div className={styles.stat}>
          <div className={styles.item}>
            <Badge count={2} size="small" offset={[-5, 1]}>
              <Tag color="#3b5999">Low</Tag>
            </Badge>
          </div>
          <div className={styles.item}>
            <Badge count={1} size="small" offset={[-5, 1]}>
              <Tag color="#F4A460">Medium</Tag>
            </Badge>
          </div>
          <div className={styles.item}>
            <Badge count={1} size="small" offset={[-5, 1]}>
              <Tag color="#cd201f">High</Tag>
            </Badge>
          </div>
        </div>
      </div>
      <Table
        dataSource={props?.data}
        columns={columns}
        rowKey={"id"}
        pagination={{
          pageSize: 5,
        }}
      />
    </div>
  );
};

export default EnhancerTable;
