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
