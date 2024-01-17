import { Collapse, DatePicker, Tag } from "antd";
import {
  ArrowRightOutlined,
  CaretRightOutlined,
  ClockCircleOutlined,
} from "@ant-design/icons";
import K8sStat from "../k8sStat";
import { SEVERITY_MAP } from "../../../../utils/constants";

import styles from "./style.module.less";

const { RangePicker } = DatePicker;

const K8sEvent = ({ execptionList, rescan, showDrawer, onItemClick }) => {
  const panelStyle: React.CSSProperties = {
    background: "#fff",
    borderRadius: 8,
    border: "1px solid rgba(0,0,0,0.15)",
    marginBottom: 8,
  };

  function getItems() {
    return []?.map((item) => {
      return {
        key: item?.id,
        label: (
          <div className={styles.collapse_panel_title}>
            <div className={styles.left}>
              <Tag bordered={false} color={SEVERITY_MAP?.[item?.level]?.color}>
                {SEVERITY_MAP?.[item?.level]?.text}
              </Tag>
            </div>
            <div className={styles.right}>
              <div className={styles.tight_top}>
                <span className={styles.title}>{item?.title}</span>
                <span>（{9}）</span>
                <span className={styles.time}>
                  <ClockCircleOutlined /> 7h40m
                </span>
              </div>
              <div className={styles.tight_bottom}>
                Message:try to switch on ali-monitor for pod
                huanyu/huanyuprodgz00c0shadow133602185-qswgl-fd7d5
              </div>
            </div>
          </div>
        ),
        children: (
          <div className={styles.collapse_panel_body}>
            <div className={styles.body}>
              <div className={styles.label}>时间触发时间点：</div>
              <div className={styles.value}>
                {item?.timeList?.map((item) => {
                  return <div className={styles.time_block}>{item}</div>;
                })}
              </div>
            </div>
          </div>
        ),
        style: panelStyle,
      };
    });
  }
  return (
    <div className={styles.k8s}>
      <div className={styles.header}>
        <div className={styles.header_left}>
          <K8sStat statData={{ all: 10, high: 5, medium: 3, low: 2 }} />
        </div>
        <div className={styles.header_right}>
          <RangePicker />
        </div>
      </div>
      <div className={styles.body}>
        <div className={styles.events}>
          <Collapse
            bordered={false}
            defaultActiveKey={["1"]}
            expandIcon={({ isActive }) => (
              <CaretRightOutlined rotate={isActive ? 90 : 0} />
            )}
            items={getItems()}
          />
        </div>
      </div>
      <div className={styles.footer} onClick={showDrawer}>
        查看全部事件
        <ArrowRightOutlined />
      </div>
    </div>
  );
};

export default K8sEvent;
