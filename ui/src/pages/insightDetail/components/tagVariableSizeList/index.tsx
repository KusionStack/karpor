import React, { useMemo } from 'react'
import { Tag, Popover } from 'antd'
import { VariableSizeList } from 'react-window'
import { getTextSizeByCanvas } from '@/utils/tools'

import styles from './style.module.less'

type IProps = {
  allTags: any
  containerWidth: number
}

const Element = ({ tag }) => {
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
    <div style={{ display: 'inline-block' }}>
      <Popover content={content}>
        <Tag className={styles.antTag}>{tag?.name}</Tag>
      </Popover>
    </div>
  )
}

const TagVariableSizeList = ({ allTags, containerWidth }: IProps) => {
  const convertDataToRows = (data, rowWidth = containerWidth) => {
    const rows = []
    let currentRow = []
    let currentRowWidth = 0
    data.forEach((item, index) => {
      const elementWidth = getTextSizeByCanvas(item?.name, 14) + 20
      if (currentRowWidth + elementWidth > rowWidth) {
        rows.push(currentRow)
        currentRow = []
        currentRowWidth = 0
      }

      currentRow.push(item)
      currentRowWidth += elementWidth

      if (index === data.length - 1) {
        rows.push(currentRow)
      }
    })

    return rows
  }

  // eslint-disable-next-line react-hooks/exhaustive-deps
  const transformedData = useMemo(() => convertDataToRows(allTags), [allTags])
  const itemSize = 30 // lineHeight
  const itemCount = transformedData?.length // row count

  return (
    <VariableSizeList
      height={itemCount * itemSize + 10}
      width={containerWidth}
      itemCount={itemCount}
      itemSize={() => itemSize}
      itemData={transformedData}
    >
      {({ index, style, data }) => {
        const itemsInRow = data?.[index]
        return (
          <div
            style={{
              ...style,
              display: 'flex',
              flexDirection: 'row',
            }}
          >
            {itemsInRow.map(item => {
              return <Element key={item?.allName} tag={item} />
            })}
          </div>
        )
      }}
    </VariableSizeList>
  )
}

export default TagVariableSizeList
