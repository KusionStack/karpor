import React from "react";
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
        console.log(item, "====irem====");
        return (
          <div
            className={styles.item}
            key={item.value as React.Key}
            onClick={() => onChange(item.value)}
          >
            <div
              className={`${styles.normal} ${
                current === item.value ? styles.active : ""
              }`}
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
