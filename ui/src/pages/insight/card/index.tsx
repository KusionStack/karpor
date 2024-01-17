import { format_with_regex } from "../../../utils/tools";

import styles from "./styles.module.less";

type IProps = {
  title: string;
  value: number;
  color?: string;
}

const Card = ({ title, value, color }: IProps) => {

  return <div className={styles.wrapper}>
    <div className={styles.top}>{title}</div>
    <div className={styles.bottom} style={{ color }}>{format_with_regex(value)}</div>
  </div>
}

export default Card;
