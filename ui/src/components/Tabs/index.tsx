import React, { memo } from "react";
import classNames from "classnames";
import styles from "./style.module.less";

type Props = {
  current: string;
  list: Array<{ label: string; value: string }>;
  onChange: (val: string) => void;
};

const KarbourTabs = ({ current, list, onChange }: Props) => {
  return (
    <div className={styles['tab-container']}>
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

export default memo(KarbourTabs);
