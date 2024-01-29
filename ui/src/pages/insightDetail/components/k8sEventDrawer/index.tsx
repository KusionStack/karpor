import React from 'react'
import {
  Collapse,
  DatePicker,
  Drawer,
  Input,
  Pagination,
  Select,
  Tag,
} from 'antd'
import {
  CaretRightOutlined,
  ClockCircleOutlined,
  SearchOutlined,
} from '@ant-design/icons'
import K8sStat from '../k8sStat'

import styles from './style.module.less'
import { SEVERITY_MAP } from '../../../../utils/constants'

const { RangePicker } = DatePicker

type K8sEventDrawerProps = {
  open: boolean
  onClose: () => void
}

const K8sEventDrawer = ({ open, onClose }: K8sEventDrawerProps) => {
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  function onSearch() {}

  // eslint-disable-next-line @typescript-eslint/no-empty-function
  function handleChangeSelect() {}

  // eslint-disable-next-line @typescript-eslint/no-empty-function
  function handleChangePage() {}

  const panelStyle: React.CSSProperties = {
    background: '#fff',
    borderRadius: 8,
    border: '1px solid rgba(0,0,0,0.15)',
    marginBottom: 8,
  }

  function getItems() {
    return []?.map(item => {
      return {
        key: item?.id,
        label: (
          <div className={styles.collapse_panel_title}>
            <div className={styles.left}>
              <Tag bordered={false} color={SEVERITY_MAP?.[item?.level]?.color}>
                {SEVERITY_MAP?.[item?.level]?.text}
              </Tag>
            </div>
            <div className={styles.right}>
              <div className={styles.tight_top}>
                <span className={styles.title}>{item?.title}</span>
                <span>（{9}）</span>
                <span className={styles.time}>
                  <ClockCircleOutlined /> 7h40m
                </span>
              </div>
              <div className={styles.tight_bottom}>
                Message:try to switch on monitor for pod foo/fooprod-qswgl-fd7d5
              </div>
            </div>
          </div>
        ),
        children: (
          <div className={styles.collapse_panel_body}>
            <div className={styles.body}>
              <div className={styles.label}>时间触发时间点：</div>
              <div className={styles.value}>
                {item?.timeList?.map(item => {
                  return (
                    <div key={item} className={styles.time_block}>
                      {item}
                    </div>
                  )
                })}
              </div>
            </div>
          </div>
        ),
        style: panelStyle,
      }
    })
  }

  return (
    <Drawer width={1000} title="异常事件" open={open} onClose={onClose}>
      <div className={styles.execption_drawer}>
        <K8sStat statData={{ all: 10, high: 5, medium: 3, low: 2 }} />
        <div className={styles.tool_bar}>
          <div className={styles.search}>
            <Input
              placeholder="请输入名称搜索"
              suffix={
                <SearchOutlined
                  className="site-form-item-icon"
                  style={{ color: '#999' }}
                />
              }
              allowClear
              style={{ width: 200 }}
              onChange={onSearch}
            />
          </div>
          <div className={styles.tool_right}>
            <Select
              defaultValue="lucy"
              style={{ width: 124, marginRight: 16 }}
              onChange={handleChangeSelect}
              options={[
                { value: 'jack', label: 'Jack' },
                { value: 'lucy', label: 'Lucy' },
                { value: 'john', label: 'John' },
                { value: 'disabled', label: 'Disabled', disabled: true },
              ]}
            />
            <RangePicker />
          </div>
        </div>
        <div className={styles.events}>
          <Collapse
            bordered={false}
            defaultActiveKey={['1']}
            expandIcon={({ isActive }) => (
              <CaretRightOutlined rotate={isActive ? 90 : 0} />
            )}
            style={{ background: '#fff' }}
            items={getItems()}
          />
        </div>
        <div style={{ textAlign: 'right', marginTop: 16 }}>
          <Pagination
            total={1000}
            // showTotal={(total, range) =>
            //   `${range[0]}-${range[1]} 共 ${total} 条`
            // }
            pageSize={20}
            current={1}
            onChange={handleChangePage}
            showSizeChanger
          />
        </div>
      </div>
    </Drawer>
  )
}

export default K8sEventDrawer
