import { Descriptions } from "antd";
import styles from "./styles.module.less";
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
      <div className={styles["yaml-btn"]} onClick={handleClick}>
        <FileTextTwoTone style={{ fontSize: 20 }} />
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
