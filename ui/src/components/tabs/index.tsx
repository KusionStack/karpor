import React, { memo } from 'react'
import classNames from 'classnames'
import { useTranslation } from 'react-i18next'
import { getDataType } from '@/utils/tools'

import styles from './style.module.less'

type Props = {
  current: string
  list: Array<{
    label: string | React.ReactNode
    value: string
    disabled?: boolean
  }>
  onChange: (val: string, index?: number) => void
  itemStyle?: any
  boxStyle?: any
}

const KarporTabs = ({
  current,
  list,
  onChange,
  itemStyle,
  boxStyle,
}: Props) => {
  const { t } = useTranslation()
  return (
    <div className={styles.tab_container} style={boxStyle}>
      {list?.map((item, index) => (
        <div
          className={styles.item}
          key={item.value as React.Key}
          onClick={() => {
            !item?.disabled && onChange(item.value, index)
          }}
          style={{
            ...itemStyle,
            ...(item?.disabled ? { color: '#f1f1f1' } : {}),
          }}
        >
          <div
            className={classNames(styles.normal, {
              [styles.active]: current === item.value,
            })}
            style={
              item?.disabled ? { color: '#999', cursor: 'not-allowed' } : {}
            }
          >
            {getDataType(item?.label) === 'String'
              ? t(item?.label as string)
              : item?.label}
          </div>
        </div>
      ))}
    </div>
  )
}

export default memo(KarporTabs)
