import styles from "./styles.module.less";

type IProps = {
  data: any;
};

const Card = ({ data }: IProps) => {
  return (
    <div className={styles.card}>
      <div className={styles.left}>{data?.desc}</div>
      <div className={styles.right}>{data?.title}</div>
    </div>
  );
};

export default Card;
