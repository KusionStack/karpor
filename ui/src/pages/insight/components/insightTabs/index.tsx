import React from 'react'
import styles from './styles.module.less'
import { PlusOutlined } from '@ant-design/icons'
import classNames from 'classnames'

export const InsightTabs = ({
  items,
  addIsDiasble,
  activeKey,
  handleClickItem,
  onEdit,
  disabledAdd,
}) => {
  function handleActionIcon(event, label) {
    event.preventDefault()
    event.stopPropagation()
    onEdit('edit', label)
  }
  function handleAdd() {
    if (disabledAdd) return
    onEdit('add')
  }
  function handleChangeTab(key) {
    if (activeKey === key) return
    handleClickItem(key)
  }
  return (
    <div className={styles.tabs_wrapper}>
      <div className={styles.tabs_container}>
        <div className={styles.tabs}>
          {items?.map(item => {
            return (
              <div
                key={item?.label}
                className={classNames(styles.tab, {
                  [styles.active_tab]: item?.label === activeKey,
                })}
                onClick={() => handleChangeTab(item?.label)}
              >
                <div className={styles.label}>{item?.label}</div>
                {item?.closeIcon && (
                  <div
                    className={styles.edit_icon}
                    onClick={event => handleActionIcon(event, item?.label)}
                  >
                    {item?.closeIcon}
                  </div>
                )}
              </div>
            )
          })}
        </div>
        <div
          className={styles.add_box}
          style={{ cursor: disabledAdd ? 'not-allowed' : 'pointer' }}
          onClick={handleAdd}
        >
          <PlusOutlined
            style={{ fontSize: 12, color: disabledAdd ? '#999' : '#000' }}
            disabled={addIsDiasble}
          />
        </div>
      </div>
    </div>
  )
}
