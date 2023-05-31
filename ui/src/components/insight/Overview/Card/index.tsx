import styles from "./styles.module.scss";

type IProps = {
  data: any
}

export default function Card({ data }: IProps) {
  return <div className={styles.card}>
    <div className={styles.left}>{data?.desc}</div>
    <div className={styles.right}>{data?.title}</div>
  </div>
}
