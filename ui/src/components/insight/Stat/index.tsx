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

import type { FC } from "react";
import { useEffect, useRef } from "react";
import Progress from "./Progress";
import { InfoCircleOutlined } from "@ant-design/icons";
import styles from "./styles.module.scss";

const Stat: FC = () => {
  return (
    <div className={styles.stat}>
      <div className={styles.item}>
        <div className={styles.top}>
          CPU
          <div className={styles["top-icon"]}>
            <InfoCircleOutlined style={{ fontSize: 12 }} />
          </div>
        </div>
        <div className={styles.bottom}>4 units</div>
      </div>
      <div className={styles.item}>
        <div className={styles.top}>
          Memory
          <div className={styles["top-icon"]}>
            <InfoCircleOutlined style={{ fontSize: 12 }} />
          </div>
        </div>
        <div className={styles.bottom}>7.77 GB</div>
      </div>
      <div className={styles.item}>
        <div>Pods</div>
        <Progress />
        <div style={{ marginTop: 10 }}>11/11 Requested</div>
      </div>
    </div>
  );
};

export default Stat;
