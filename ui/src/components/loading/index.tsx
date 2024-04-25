import React from 'react'
import { Spin } from 'antd'
import styles from './style.module.less'

const Loading = () => {
  return (
    <div className={styles.loading}>
      <Spin />
    </div>
  )
}

export default Loading
