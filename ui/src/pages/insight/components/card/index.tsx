import React from 'react'
import { Divider, Popover, Tag, Tooltip } from 'antd'
import { useTranslation } from 'react-i18next'

import styles from './styles.module.less'

type IProps = {
  allTags?: any[]
  handleClick?: (group: any, title: string) => void
  group?: any
}

const CardContent = ({ allTags, handleClick, group }: IProps) => {
  const { t } = useTranslation()
  const title = group?.title

  const tagStyle: React.CSSProperties = {
    maxWidth: '100%',
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    display: 'inline-block',
    marginBottom: 5,
  }

  const content = (
    <>
      {allTags?.map((item, index) => (
        <div key={item?.key}>
          <Tag
            style={{
              ...tagStyle,
              marginBottom: index === allTags?.length - 1 ? 0 : 5,
            }}
          >{`${item?.key}: ${item?.value}`}</Tag>
        </div>
      ))}
    </>
  )

  return (
    <div className={styles.insight_card}>
      <div className={styles.content}>
        <div
          className={styles.content_left}
          onClick={() => handleClick(group, title)}
        >
          {title?.slice(0, 1)}
        </div>
        <div className={styles.content_right}>
          <Tooltip title={title}>
            <div
              className={styles.title}
              onClick={() => handleClick(group, title)}
            >
              {title}
            </div>
          </Tooltip>
          <div className={styles.tag_container}>
            {allTags?.slice(0, 2)?.map(item => (
              <Popover
                key={item?.key}
                content={
                  <span>
                    {item?.key}: {item?.value}
                  </span>
                }
              >
                <Tag style={tagStyle}>{`${item?.key}: ${item?.value}`}</Tag>
              </Popover>
            ))}
          </div>
        </div>
      </div>
      <div className={styles.footer}>
        <Popover content={content}>
          <div className={styles.btn_item}>{t('AllTags')}</div>
        </Popover>
        <Divider type="vertical" />
        <div
          className={styles.btn_item}
          onClick={() => handleClick(group, title)}
        >
          {t('View')}
        </div>
      </div>
    </div>
  )
}

export default CardContent
