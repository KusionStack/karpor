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

import React from "react";
import classNames from "classnames";
import styles from "./style.module.scss";

type Props = {
  current: string;
  list: Array<{ label: string; value: string }>;
  onChange: (val: string) => void;
};

const KarbourTabs = ({ current, list, onChange }: Props) => {
  return (
    <div className={styles.tabContainer}>
      {list?.map((item) => {
        return (
          <div
            className={styles.item}
            key={item.value as React.Key}
            onClick={() => onChange(item.value)}
          >
            <div
              className={classNames(styles.normal, { [styles.active]: current === item.value })}
            >
              {item.label}
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default KarbourTabs;
