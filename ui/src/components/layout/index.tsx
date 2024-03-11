import React, { memo } from 'react'
import { Divider, Menu, Popover } from 'antd'
import {
  ClusterOutlined,
  FundOutlined,
  QuestionCircleOutlined,
  SearchOutlined,
} from '@ant-design/icons'
import type { MenuProps } from 'antd'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useDispatch, useSelector } from 'react-redux'
import { setServerConfigMode } from '@/store/modules/globalSlice'
import axios from 'axios'
import { useTranslation } from 'react-i18next'
import showPng from '@/assets/show.png'
import languagePng from '@/assets/language.png'

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
  const dispatch = useDispatch()
  const { isReadOnlyMode } = useSelector((state: any) => state.globalSlice)
  const { i18n, t } = useTranslation()

  async function getServerConfigs() {
    const response: any = await axios(`/server-configs`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      params: {},
    })
    dispatch(setServerConfigMode(response?.CoreOptions?.ReadOnlyMode))
  }

  getServerConfigs()

  const menuItems = [
    getItem(t('Search'), '/search', <SearchOutlined />),
    getItem(
      t('SearchResult'),
      '/search/result',
      <SearchOutlined />,
      null,
      null,
      true,
    ),
    getItem(t('Insight'), '/insight', <FundOutlined />, null, null, null, true),
    getItem(t('ClusterManagement'), '/cluster', <ClusterOutlined />),
    getItem(
      t('ClusterDetail'),
      '/insightDetail',
      <SearchOutlined />,
      null,
      null,
      true,
    ),
    getItem(
      'RegisterCluster',
      '/cluster/access',
      <SearchOutlined />,
      null,
      null,
      true,
    ),
    getItem(
      t('RotateCertificate'),
      '/cluster/certificate',
      <SearchOutlined />,
      null,
      null,
      true,
    ),
    getItem(
      t('DataSyncConfiguration'),
      '/reflux',
      <SearchOutlined />,
      null,
      null,
      true,
    ),
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

  function handleMenuClick(event) {
    if (event?.domEvent.metaKey && event?.domEvent.button === 0) {
      const { origin } = window.location
      window.open(`${origin}${event.key}`)
    } else {
      navigate(event.key)
    }
  }

  function goHome() {
    navigate('/')
  }
  function changeToZh() {
    localStorage.setItem('lang', 'zh')
    i18n.changeLanguage('zh')
  }
  function changeToEn() {
    localStorage.setItem('lang', 'en')
    i18n.changeLanguage('en')
  }
  const languageContent = (
    <div className={styles.language_content}>
      <div className={styles.language_content_item} onClick={changeToZh}>
        中文
      </div>
      <div className={styles.language_content_item} onClick={changeToEn}>
        English
      </div>
    </div>
  )

  return (
    <div className={styles.wrapper}>
      <div className={styles.nav}>
        <div className={styles.left}>
          <div className={styles.title} onClick={goHome}>
            <div className={styles.sub_logo}>K</div>
            <div className={styles.text}>Karbour</div>
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
          {isReadOnlyMode && (
            <div className={styles.read_only_mode}>
              <img src={showPng} />
              <span>{t('ReadOnlyMode')}</span>
            </div>
          )}
          <div className={styles.help}>
            <Popover content={languageContent} trigger="click">
              <div className={styles.language}>
                <img src={languagePng} />
              </div>
            </Popover>
          </div>
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
