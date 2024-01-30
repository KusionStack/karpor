import React, { useState } from 'react'
import { Tag, Popover } from 'antd'
import { DoubleLeftOutlined, DoubleRightOutlined } from '@ant-design/icons'

import styles from './style.module.less'

type IProps = {
  allTags: any
}

const MutiTag = ({ allTags }: IProps) => {
  // 设置展示全部标签的状态，默认为false
  const [showAll, setShowAll] = useState(false)
  // 默认标签的最大展示数量
  const defaultMaxCount = 5

  // 创建一个函数来切换展示状态
  const toggleTags = () => {
    setShowAll(!showAll)
  }

  return (
    <div className={styles.appContainer}>
      {allTags
        .slice(0, showAll ? allTags.length : defaultMaxCount)
        .map((tag: any) => {
          const content = (
            <div className={styles.popCard}>
              <div className={styles.item}>
                <span className={styles.label}>cluster：</span>
                <span className={styles.value}>{tag?.cluster || '--'}</span>
              </div>
              <div className={styles.item}>
                <span className={styles.label}>apiVersion：</span>
                <span className={styles.value}>{tag?.apiVersion || '--'}</span>
              </div>
              <div className={styles.item}>
                <span className={styles.label}>kind：</span>
                <span className={styles.value}>{tag?.kind || '--'}</span>
              </div>
              <div className={styles.item}>
                <span className={styles.label}>namespace：</span>
                <span className={styles.value}>{tag?.namespace || '--'}</span>
              </div>
              <div className={styles.item}>
                <span className={styles.label}>name：</span>
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
      {/* 当标签数多于5个且showAll为false时，显示更多按钮 */}
      {allTags.length > defaultMaxCount && !showAll && (
        <div className={styles.toggleButton} onClick={toggleTags}>
          <span>
            更多
            <DoubleLeftOutlined
              style={{ transform: 'rotate(-90deg)', marginLeft: 5 }}
            />
          </span>
        </div>
      )}
      {/* 当showAll为true时，显示收起按钮 */}
      {showAll && (
        <div className={styles.toggleButton} onClick={toggleTags}>
          <span>
            收起
            <DoubleRightOutlined
              style={{ transform: 'rotate(-90deg)', marginLeft: 5 }}
            />
          </span>
        </div>
      )}
    </div>
  )
}

export default MutiTag
