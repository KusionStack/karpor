import React, { memo } from 'react'
import { Divider, Menu } from 'antd'
import {
  ClusterOutlined,
  FundOutlined,
  QuestionCircleOutlined,
  SearchOutlined,
} from '@ant-design/icons'
import type { MenuProps } from 'antd'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import styles from './style.module.less'

type MenuItem = Required<MenuProps>['items'][number]

function getItem(
  label: React.ReactNode,
  key: React.Key,
  icon?: React.ReactNode,
  children?: MenuItem[],
  type?: 'group',
  hidden?: boolean,
  disabled?: boolean,
): MenuItem {
  return {
    key,
    icon,
    children,
    label,
    type,
    hidden,
    disabled,
  } as MenuItem
}

const LayoutPage = () => {
  const navigate = useNavigate()
  const { pathname } = useLocation()

  const menuItems = [
    getItem('搜索', '/search', <SearchOutlined />),
    getItem('结果', '/search/result', <SearchOutlined />, null, null, true),
    getItem('数据洞察', '/insight', <FundOutlined />, null, null, null, true),
    getItem('集群列表', '/cluster', <ClusterOutlined />),
    getItem('集群详情', '/insightDetail', <SearchOutlined />, null, null, true),
    getItem(
      '集群接入',
      '/cluster/access',
      <SearchOutlined />,
      null,
      null,
      true,
    ),
    getItem(
      '更新证书',
      '/cluster/certificate',
      <SearchOutlined />,
      null,
      null,
      true,
    ),
    getItem('回流配置', '/reflux', <SearchOutlined />, null, null, true),
  ]

  function getKey() {
    return [pathname]
  }

  function getMenuItems() {
    function loop(list) {
      return list
        ?.filter(item => !item?.hidden)
        ?.map(item => {
          if (item?.children) {
            item.children = loop(item?.children)
          }
          return item
        })
    }
    return loop(menuItems)
  }

  function handleMenuClick(e) {
    navigate(e.key)
  }

  function goHome() {
    navigate('/search')
  }

  return (
    <div className={styles.wrapper}>
      <div className={styles.nav}>
        <div className={styles.left}>
          <div className={styles.title} onClick={goHome}>
            <div className={styles.subLogo}>K</div>
            <div className={styles.text}>Karbour 数据门户</div>
          </div>
          <div>
            <Divider type="vertical" />
          </div>
          <Menu
            style={{ flex: 1, border: 'none' }}
            mode="horizontal"
            selectedKeys={getKey()}
            items={getMenuItems()}
            onClick={handleMenuClick}
          />
        </div>
        <div className={styles.right}>
          <div className={styles.help}>
            <a
              target="_blank"
              href="https://github.com/KusionStack/karbour"
              rel="noreferrer"
            >
              <QuestionCircleOutlined style={{ color: '#999' }} />
            </a>
          </div>
        </div>
      </div>
      <div className={styles.content}>
        <div className={styles.right}>
          <div className={styles.right_content}>
            <Outlet />
          </div>
        </div>
      </div>
    </div>
  )
}

export default memo(LayoutPage)
