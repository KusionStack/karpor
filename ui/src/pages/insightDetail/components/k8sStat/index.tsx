import React from 'react'
import { Badge } from 'antd'
import { useTranslation } from 'react-i18next'
import classNames from 'classnames'

import styles from './style.module.less'

type IProps = {
  statData: any
}

const K8sStat = ({ statData }: IProps) => {
  const { t } = useTranslation()
  return (
    <div className={styles.exception_stat}>
      <div className={styles.title}>
        {t('AllIssues')}
        <span className={styles.num}>{statData?.all}</span>
      </div>
      <div className={classNames(styles.title, styles.height)}>
        <Badge status="error" text={t('Exception')}></Badge>
        <span className={styles.num}>{statData?.high || '--'}</span>
      </div>
      <div className={styles.title}>
        <Badge status="warning" text={t('Warning')}></Badge>
        <span className={styles.num}>{statData?.medium || '--'}</span>
      </div>
      <div className={styles.title}>
        <Badge status="success" text={t('Normal')}></Badge>
        <span className={styles.num}>{statData?.low || '--'}</span>
      </div>
    </div>
  )
}

export default K8sStat
