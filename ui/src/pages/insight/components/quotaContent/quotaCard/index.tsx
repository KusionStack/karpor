import React, { ReactNode } from 'react'
import { format_with_regex } from '@/utils/tools'

import styles from './styles.module.less'

type IProps = {
  title: string
  value: number
  color?: string
  bgColor?: string
  iconNode?: ReactNode
}

const QuotaCard = ({ title, value, color, bgColor, iconNode }: IProps) => {
  return (
    <div className={styles.wrapper} style={{ background: bgColor }}>
      <div className={styles.top}>{title}</div>
      <div className={styles.bottom}>
        <div className={styles.icon}>{iconNode}</div>
        <div className={styles.status}>
          <div className={styles.num} style={{ color }}>
            {value ? format_with_regex(value) : 0}
          </div>
        </div>
      </div>
    </div>
  )
}

export default QuotaCard
