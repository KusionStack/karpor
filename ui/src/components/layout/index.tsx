import React, { memo, useEffect } from 'react'
import { Divider, Menu, Dropdown } from 'antd'
import {
  ClusterOutlined,
  FundOutlined,
  QuestionCircleOutlined,
  SearchOutlined,
  CaretDownOutlined,
  GithubOutlined,
} from '@ant-design/icons'
import type { MenuProps } from 'antd'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useDispatch, useSelector } from 'react-redux'
import { setServerConfigMode } from '@/store/modules/globalSlice'
import axios from 'axios'
import { useTranslation } from 'react-i18next'
import showPng from '@/assets/show.png'
import logo from '@/assets/img/logo.svg'
import languageSvg from '@/assets/translate_language.svg'
import { Languages, LanguagesMap } from '@/utils/constants'

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
    const response: any = await axios.get(`/server-configs`)
    if (response) {
      dispatch(setServerConfigMode(response?.CoreOptions?.ReadOnlyMode))
    }
  }

  useEffect(() => {
    getServerConfigs()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

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
    getItem(t('Insight'), '/insight', <FundOutlined />),
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

  function handleChangeLanguage(lang) {
    localStorage.setItem('lang', lang)
    i18n.changeLanguage(lang)
  }

  const languageItems: MenuProps['items'] = Languages.map(item => ({
    label: (
      <div
        className={styles.language_content_item}
        onClick={() => handleChangeLanguage(item?.value)}
      >
        {item?.label}
      </div>
    ),
    key: item?.value,
  }))

  return (
    <div className={styles.wrapper}>
      <div className={styles.nav}>
        <div className={styles.left}>
          <div className={styles.title} onClick={() => navigate('/')}>
            <div className={styles.sub_logo}>
              <img src={logo} />
            </div>
            <div className={styles.text}>Karpor</div>
          </div>
          <div>
            <Divider type="vertical" />
          </div>
          <Menu
            style={{ flex: 1, border: 'none', fontSize: 13 }}
            mode="horizontal"
            selectedKeys={[pathname]}
            items={getMenuItems()}
            onClick={handleMenuClick}
          />
        </div>
        <div className={styles.right}>
          {isReadOnlyMode && (
            <div className={styles.read_only_mode}>
              <img className={styles.read_only_mode_img} src={showPng} />
              <span>{t('ReadOnlyMode')}</span>
            </div>
          )}
          <div className={styles.help}>
            <Dropdown menu={{ items: languageItems }}>
              <a
                onClick={e => e.preventDefault()}
                className={styles.help_container}
              >
                <img src={languageSvg} />
                <span className={styles.help_text}>
                  {LanguagesMap?.[i18n.language || 'en']}
                </span>
                <CaretDownOutlined style={{ color: '#646566' }} />
              </a>
            </Dropdown>
          </div>
          <div className={styles.help}>
            <a
              target="_blank"
              href="https://github.com/KusionStack/karpor"
              rel="noreferrer"
            >
              <GithubOutlined style={{ color: '#646566' }} />
            </a>
          </div>
          <div className={styles.help}>
            <a
              target="_blank"
              href="https://kusionstack.io/karpor"
              rel="noreferrer"
            >
              <QuestionCircleOutlined style={{ color: '#646566' }} />
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
