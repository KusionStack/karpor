import { Button, Modal, Tag } from "antd";
import styles from "./style.module.less";
import MutiTag from "../mutiTag";
import { SEVERITY_MAP } from "../../../../utils/constants";

const EventDetail = ({ open, detail, cancel, onOk }) => {
  const locatorsNames = detail?.locators?.map((item) => {
    return {
      ...item,
      allName: `${item?.cluster || ""} ${item?.apiVersion || ""} ${item?.kind || ""} ${item?.namespace || ""} ${item?.name || ""} `,
    };
  });
  return (
    <Modal
      title="异常事件详情"
      open={open}
      destroyOnClose
      maskClosable
      onCancel={cancel}
      footer={[
        <Button key="closebBtn" onClick={cancel}>
          关闭
        </Button>,
      ]}
    >
      <div className={styles.container}>
        <div className={styles.title}>
          <Tag color={SEVERITY_MAP?.[detail?.issue?.severity]?.color}>
            {SEVERITY_MAP?.[detail?.issue?.severity]?.text}
          </Tag>
          {detail?.issue?.title || "--"}
        </div>
        <div className={styles.content}>
          <div className={styles.desc}>
            <div className={styles.item}>
              <div className={styles.label}>事件来源：</div>
              <div className={styles.value}>{detail?.issue?.scanner}</div>
            </div>
            <div className={styles.item}>
              <div className={styles.label}>发生次数：</div>
              <div className={styles.value}>{detail?.locators?.length}次</div>
            </div>
            <div
              className={styles.item}
              style={{ width: "100%", alignItems: "baseline" }}
            >
              <div className={styles.label}>描述信息：</div>
              <div className={styles.value}>
                <div className={styles.value}>{detail?.issue?.message}</div>
              </div>
            </div>
            {/* <div className={styles.item}>
              <div className={styles.label}>创建时间：</div>
              <div className={styles.value}>2023-11-24 20:33</div>
            </div> */}
          </div>
          <div className={styles.footer}>
            <div className={styles.footer_title}>相关资源：</div>
            <div className={styles.soultion}>
              <MutiTag allTags={locatorsNames} />
            </div>
          </div>
        </div>
      </div>
    </Modal>
  );
};

export default EventDetail;
