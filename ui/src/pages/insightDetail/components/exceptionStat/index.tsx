import React from 'react'
import { Tag } from 'antd'
import { useTranslation } from 'react-i18next'
import styles from './style.module.less'

type IProps = {
  statData: any
  onClickTable?: (val: string) => void
  currentKey: string
}

const ExceptionStat = ({ statData, onClickTable, currentKey }: IProps) => {
  const { t } = useTranslation()
  return (
    <div className={styles.exception_stat}>
      <div
        className={`${styles.title} ${currentKey === 'All' ? styles.active : ''}`}
        onClick={() => onClickTable('All')}
      >
        {t('AllIssues')}
        <span className={styles.num}>{statData?.all || 0}</span>
      </div>

      <div
        className={`${styles.title} ${currentKey === 'High' ? styles.active : ''}`}
        onClick={() => onClickTable('High')}
      >
        <Tag color="error">{t('High')}</Tag>
        {t('HighRisk')}
        <span className={styles.num}>{statData?.high || 0}</span>
      </div>

      <div
        className={`${styles.title} ${currentKey === 'Medium' ? styles.active : ''}`}
        onClick={() => onClickTable('Medium')}
      >
        <Tag color="warning">{t('Medium')}</Tag>
        {t('MediumRisk')}
        <span className={styles.num}>{statData?.medium || 0}</span>
      </div>

      <div
        className={`${styles.title} ${currentKey === 'Low' ? styles.active : ''}`}
        onClick={() => onClickTable('Low')}
      >
        <Tag color="success">{t('Low')}</Tag>
        {t('LowRisk')}
        <span className={styles.num}>{statData?.low || 0}</span>
      </div>
    </div>
  )
}

export default ExceptionStat
