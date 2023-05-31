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

import { useEffect, useRef } from "react";
import Card from "./Card";
import { Descriptions } from "antd";
import styles from "./styles.module.scss";
import { FileTextTwoTone } from "@ant-design/icons";

type IProps = {
  data: {
    title: string;
    list: Array<{ title: string; desc: string }>;
  };
  handleClick: () => void;
};

const Overview = ({ data, handleClick }: IProps) => {
  return (
    <div className={styles.overview}>
      <div className={styles.yamlBtn} onClick={handleClick}>
        <FileTextTwoTone style={{fontSize: 20}}/>
      </div>
      {/* <div className={styles.title}>{data?.title}</div> */}
      <div className={styles.content}>
        <Descriptions
          column={1}
          contentStyle={{ fontSize: 16, marginBottom: 10 }}
          labelStyle={{ fontSize: 16 }}
        >
          {data?.list?.map(
            (item: { title: string; desc: string }, index: number) => {
              return (
                <Descriptions.Item key={index} label={item?.desc}>
                  {item?.title}
                </Descriptions.Item>
              );
            }
          )}
        </Descriptions>
      </div>
    </div>
  );
};

export default Overview;
