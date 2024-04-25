import React from 'react'
import { Divider, Popover, Tag, Tooltip } from 'antd'
import { useTranslation } from 'react-i18next'

import styles from './styles.module.less'

const CardContent = ({ allTags, handleClick, group }) => {
  const { t } = useTranslation()
  const title = group?.title

  const tagStyle = {
    maxWidth: '100%',
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    display: 'inline-block',
    marginBottom: 5,
  }

  const content = (
    <div style={{ width: 300 }}>
      {allTags?.map(item => {
        return (
          <Tag
            key={`${item?.key}`}
            style={tagStyle}
          >{`${item?.key}:${item?.value}`}</Tag>
        )
      })}
    </div>
  )

  return (
    <div className={styles.insight_card}>
      <div className={styles.content}>
        <div className={styles.content_left}>{title?.slice(0, 1)}</div>
        <div className={styles.content_right}>
          <Tooltip title={title}>
            <div className={styles.title}>{title}</div>
          </Tooltip>
          <div className={styles.tag_container}>
            {allTags?.slice(0, 2)?.map(item => {
              return (
                <Tag
                  style={tagStyle}
                  key={`${item?.key}`}
                  title={`${item?.key}: ${item?.value}`}
                >
                  {`${item?.key}: ${item?.value}`}
                </Tag>
              )
            })}
          </div>
        </div>
      </div>
      <div className={styles.footer}>
        <Popover content={content}>
          <div className={styles.btn_item}>{t('AllTags')}</div>
        </Popover>
        <div>
          <Divider type="vertical" />
        </div>
        <div
          className={styles.btn_item}
          onClick={() => handleClick(group, title)}
        >
          {t('Check')}
        </div>
      </div>
    </div>
  )
}

export default CardContent
