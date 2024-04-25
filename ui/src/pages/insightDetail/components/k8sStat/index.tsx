import React from 'react'
import { Badge } from 'antd'
import { useTranslation } from 'react-i18next'

import styles from './style.module.less'

type IProps = {
  statData: any
}

const K8sStat = ({ statData }: IProps) => {
  const { t } = useTranslation()
  return (
    <div className={styles.exception_stat}>
      <div className={styles.title}>
        全部<span className={styles.num}>{statData?.all}</span>
      </div>
      <div className={`${styles.title} ${styles.height}`}>
        <Badge status="error" text={t('Exception')}></Badge>
        <span className={styles.num}>{statData?.high || 5}</span>
      </div>
      <div className={styles.title}>
        <Badge status="warning" text={t('Warning')}></Badge>
        <span className={styles.num}>{statData?.medium || 2}</span>
      </div>
      <div className={styles.title}>
        <Badge status="success" text={t('Normal')}></Badge>
        <span className={styles.num}>{statData?.low || 1}</span>
      </div>
    </div>
  )
}

export default K8sStat
