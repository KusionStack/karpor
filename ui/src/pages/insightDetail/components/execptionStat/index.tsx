import React from 'react'
import { Tag } from 'antd'
import styles from './style.module.less'

type IProps = {
  statData: any
  onClickTable?: (val: string) => void
  currentKey: string
}

const ExecptionStat = ({ statData, onClickTable, currentKey }: IProps) => {
  return (
    <div className={styles.execption_stat}>
      <div
        className={`${styles.title} ${currentKey === 'All' ? styles.active : ''}`}
        onClick={() => onClickTable('All')}
      >
        全部异常事件
        <span className={styles.num}>{statData?.all || 0}</span>
      </div>
      <div
        className={`${styles.title} ${currentKey === 'High' ? styles.active : ''}`}
        onClick={() => onClickTable('High')}
      >
        <Tag color="error">高</Tag>
        高风险
        <span className={styles.num}>{statData?.high || 0}</span>
      </div>
      <div
        className={`${styles.title} ${currentKey === 'Medium' ? styles.active : ''}`}
        onClick={() => onClickTable('Medium')}
      >
        <Tag color="warning">中</Tag>
        中风险
        <span className={styles.num}>{statData?.medium || 0}</span>
      </div>
      <div
        className={`${styles.title} ${currentKey === 'Low' ? styles.active : ''}`}
        onClick={() => onClickTable('Low')}
      >
        <Tag color="success">低</Tag>
        低风险
        <span className={styles.num}>{statData?.low || 0}</span>
      </div>
    </div>
  )
}

export default ExecptionStat
