import { Progress, Dropdown, Space, Input, Select } from 'antd'
import { DownOutlined, SearchOutlined } from '@ant-design/icons'
import type { MenuProps } from 'antd'
import Card from './card'

import styles from './styles.module.less'
import React from 'react'

const Insight = () => {
  const items: MenuProps['items'] = [
    {
      key: '1',
      label: 'label1',
    },
    {
      key: '2',
      label: 'label2',
      disabled: true,
    },
    {
      key: '3',
      danger: true,
      label: 'a danger item',
    },
  ]

  // eslint-disable-next-line @typescript-eslint/no-empty-function
  function handleChangeSelect() {}

  return (
    <div className={styles.container}>
      <div className={styles.pageTitle}>数据洞察</div>
      <div className={styles.content}>
        <div className={styles.header}>
          <Dropdown menu={{ items }}>
            <div className={styles.dropText} onClick={e => e.preventDefault()}>
              <Space
                style={{
                  color: 'rgba(0,10,26,0.89)',
                  fontWeight: 500,
                  fontSize: 14,
                }}
              >
                Hover me
                <DownOutlined style={{ paddingTop: 3 }} />
              </Space>
            </div>
          </Dropdown>
          <div className={styles.selectTips}>
            <span className={styles.bulb}>💡 </span>点击下拉切换资源类型
          </div>
        </div>
        <div className={styles.stat}>
          <div className={styles.circle}>
            {/* <MemoPiePercent width={96} height={96}/> */}
            <Progress
              type="circle"
              percent={80}
              size={96}
              strokeColor="#59D226"
              trailColor="#2F54EB"
              format={() => '集群'}
            />
          </div>
          <div className={styles.all}>
            <Card title="集群总数" value={5000} />
          </div>
          <div className={styles.symbol}>=</div>
          <div className={styles.exception}>
            <Card title="异常数量" value={44} color="#FF4D4F" />
          </div>
          <div className={styles.symbol}>+</div>
          <div className={styles.health}>
            <Card title="健康数量" value={89} />
          </div>
        </div>
        <div className={styles.toolBar}>
          <div className={styles.left}>集群列表</div>
          <div className={styles.right}>
            <Select
              defaultValue="lucy"
              style={{ width: 124 }}
              onChange={handleChangeSelect}
              options={[
                { value: 'jack', label: 'Jack' },
                { value: 'lucy', label: 'Lucy' },
                { value: 'john', label: 'John' },
                { value: 'disabled', label: 'Disabled', disabled: true },
              ]}
            />
            <Input
              style={{ width: 160, marginLeft: 16 }}
              placeholder="请输入搜索关键字"
              suffix={<SearchOutlined />}
            />
          </div>
        </div>
      </div>
    </div>
  )
}

export default Insight
