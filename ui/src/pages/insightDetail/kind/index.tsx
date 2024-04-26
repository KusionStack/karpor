import React, { useEffect, useState } from 'react'
import { NavLink, useLocation, useNavigate } from 'react-router-dom'
import axios from 'axios'
import queryString from 'query-string'
import { Breadcrumb, Tooltip, message } from 'antd'
import { useTranslation } from 'react-i18next'
import ExceptionDrawer from '../components/exceptionDrawer'
import SourceTable from '../components/sourceTable'
import ExceptionList from '../components/exceptionList'
import EventDetail from '../components/eventDetail'
import SummaryCard from '../components/summaryCard'
import { ICON_MAP } from '@/utils/images'
import { capitalized } from '@/utils/tools'

import styles from './styles.module.less'

const ClusterDetail = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const urlParams = queryString.parse(location?.search)
  const { type, apiVersion, cluster, kind, namespace, name, key, from, query } =
    urlParams
  const [drawerVisible, setDrawerVisible] = useState<boolean>(false)
  const [modalVisible, setModalVisible] = useState<boolean>(false)
  const [tableQueryStr] = useState(
    `select * from resources where cluster = '${cluster}' and apiVersion='${apiVersion}' and kind = '${kind}'`,
  )
  const [auditList, setAuditList] = useState<any>([])
  const [auditLoading, setAuditLoading] = useState<any>(false)
  const [auditStat, setAuditStat] = useState<any>()
  const [tableName, setTableName] = useState(kind as any)
  const [breadcrumbItems, setBreadcrumbItems] = useState([])
  const [summary, setSummary] = useState<any>()
  const [currentItem, setCurrentItem] = useState<any>()
  const { t, i18n } = useTranslation()

  async function getAudit(isRescan) {
    setAuditLoading(true)
    const response: any = await axios({
      url: '/rest-api/v1/insight/audit',
      method: 'GET',
      params: {
        apiVersion,
        kind,
        cluster,
        ...(isRescan ? { forceNew: true } : {}),
      },
    })
    setAuditLoading(false)
    if (response?.success) {
      setAuditList(response?.data)
    } else {
      message.error(response?.message || t('RequestFailedAndTry'))
    }
  }
  async function getAuditScore() {
    const response: any = await axios({
      url: '/rest-api/v1/insight/score',
      method: 'GET',
      params: {
        cluster,
        apiVersion,
        kind,
      },
    })
    if (response?.success) {
      setAuditStat(response?.data)
    }
  }

  async function getSummary() {
    const response: any = await axios({
      url: '/rest-api/v1/insight/summary',
      params: {
        cluster,
        apiVersion,
        kind,
      },
    })
    if (response?.success) {
      setSummary(response?.data)
    } else {
      message.error(response?.message || t('RequestFailedAndTry'))
    }
  }

  useEffect(() => {
    getAudit(false)
    getAuditScore()
    getSummary()
    setTableName(kind as any)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [kind, type])

  function rescan() {
    getAuditScore()
    getAudit(true)
  }

  function showDrawer() {
    setDrawerVisible(true)
  }

  function onItemClick(item) {
    setModalVisible(true)
    setCurrentItem(item)
  }

  function replacePage(item) {
    const obj = { from, type, apiVersion, query }
    const list = ['cluster', 'kind']
    for (let i = 0; i < list?.length; i++) {
      if (list[i] === item) {
        obj[list[i]] = urlParams[list[i]]
        obj.type = item
        break
      } else {
        obj[list[i]] = urlParams[list[i]]
      }
    }
    if (from === 'result') {
      obj.query = query
    }
    const urlStringfyParams = queryString.stringify(obj)
    navigate(`/insightDetail/${item}?${urlStringfyParams}`, { replace: true })
  }

  function getBreadcrumbs() {
    let first
    if (from === 'cluster') {
      first = {
        title: <NavLink to={'/cluster'}>{t('ClusterManagement')}</NavLink>,
      }
    }
    if (from === 'result') {
      first = {
        title: (
          <NavLink to={`/search/result?query=${query}`}>
            {t('SearchResult')}
          </NavLink>
        ),
      }
    }
    if (from === 'insight') {
      first = {
        title: <NavLink to={`/insight`}>{t('Insight')}</NavLink>,
      }
    }
    const middle = []
    ;['cluster', 'kind']?.forEach(item => {
      const urlParamsItem = urlParams?.[item]
      if (urlParamsItem) {
        const iconMap = {
          cluster: ICON_MAP?.Kubernetes,
          kind: ICON_MAP?.[urlParamsItem as any] || ICON_MAP.CRD,
        }
        const iconStyle = {
          width: 14,
          height: 14,
          marginRight: 2,
        }
        middle.push({
          key: item,
          label: urlParamsItem,
          title:
            item === 'kind' ? (
              <Tooltip title={capitalized(item)}>
                <a style={{ display: 'flex', alignItems: 'center' }}>
                  <img src={iconMap?.[item]} style={iconStyle} />
                  {urlParamsItem}
                </a>
              </Tooltip>
            ) : (
              <Tooltip title={capitalized(item)}>
                <a
                  onClick={() => replacePage(item)}
                  style={{ display: 'flex', alignItems: 'center' }}
                >
                  <img src={iconMap?.[item]} style={iconStyle} />
                  {urlParamsItem}
                </a>
              </Tooltip>
            ),
        })
      }
    })
    middle[middle?.length - 1] = {
      label: middle[middle?.length - 1]?.lebel,
      title: middle[middle?.length - 1]?.title,
    }
    const result = [first, ...middle]
    setBreadcrumbItems(result)
  }

  useEffect(() => {
    getBreadcrumbs()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [from, key, cluster, kind, namespace, name, i18n?.language])

  return (
    <div className={styles.container}>
      <Breadcrumb
        style={{ marginBottom: 20 }}
        separator=">"
        items={breadcrumbItems}
      />
      <ExceptionDrawer
        open={drawerVisible}
        onClose={() => setDrawerVisible(false)}
        exceptionList={auditList}
        exceptionStat={auditStat}
      />
      <EventDetail
        open={modalVisible}
        cancel={() => setModalVisible(false)}
        detail={currentItem}
      />
      <div className={styles.module}>
        <SummaryCard auditStat={auditStat} summary={summary} />
        <div className={styles.exception_event}>
          <ExceptionList
            auditLoading={auditLoading}
            rescan={rescan}
            exceptionList={auditList}
            exceptionStat={auditStat}
            showDrawer={showDrawer}
            onItemClick={onItemClick}
          />
        </div>
      </div>

      {/* 拓扑图 */}
      <div className={styles.tab_content}>
        <SourceTable queryStr={tableQueryStr} tableName={tableName} />
      </div>
    </div>
  )
}

export default ClusterDetail
