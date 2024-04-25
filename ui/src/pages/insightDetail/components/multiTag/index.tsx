import React, { useState } from 'react'
import { Tag, Popover } from 'antd'
import { useTranslation } from 'react-i18next'
import { DoubleLeftOutlined, DoubleRightOutlined } from '@ant-design/icons'

import styles from './style.module.less'

type IProps = {
  allTags: any
}

const MultiTag = ({ allTags }: IProps) => {
  const { t } = useTranslation()
  const [showAll, setShowAll] = useState(false)
  const defaultMaxCount = 5

  const toggleTags = () => {
    setShowAll(!showAll)
  }

  return (
    <div className={styles.appContainer}>
      {allTags
        ?.slice(0, showAll ? allTags?.length : defaultMaxCount)
        ?.map((tag: any) => {
          const content = (
            <div className={styles.popCard}>
              <div className={styles.item}>
                <span className={styles.label}>Cluster: </span>
                <span className={styles.value}>{tag?.cluster || '--'}</span>
              </div>
              <div className={styles.item}>
                <span className={styles.label}>APIVersion: </span>
                <span className={styles.value}>{tag?.apiVersion || '--'}</span>
              </div>
              <div className={styles.item}>
                <span className={styles.label}>Kind: </span>
                <span className={styles.value}>{tag?.kind || '--'}</span>
              </div>
              <div className={styles.item}>
                <span className={styles.label}>Namespace: </span>
                <span className={styles.value}>{tag?.namespace || '--'}</span>
              </div>
              <div className={styles.item}>
                <span className={styles.label}>Name: </span>
                <span className={styles.value}>{tag?.name || '--'}</span>
              </div>
            </div>
          )
          return (
            <div key={tag?.allName} style={{ display: 'inline-block' }}>
              <Popover content={content}>
                <Tag className={styles.antTag}>{tag?.name}</Tag>
              </Popover>
            </div>
          )
        })}
      {allTags?.length > defaultMaxCount && !showAll && (
        <div className={styles.toggleButton} onClick={toggleTags}>
          <span>
            {t('More')}
            <DoubleLeftOutlined
              style={{ transform: 'rotate(-90deg)', marginLeft: 5 }}
            />
          </span>
        </div>
      )}
      {showAll && (
        <div className={styles.toggleButton} onClick={toggleTags}>
          <span>
            {t('Less')}
            <DoubleRightOutlined
              style={{ transform: 'rotate(-90deg)', marginLeft: 5 }}
            />
          </span>
        </div>
      )}
    </div>
  )
}

export default MultiTag
