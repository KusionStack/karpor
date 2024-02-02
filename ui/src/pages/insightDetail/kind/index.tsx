import { useEffect, useState } from 'react'
import { NavLink, useLocation, useNavigate } from 'react-router-dom'
import axios from 'axios'
import queryString from 'query-string'
import { Breadcrumb, message } from 'antd'
import ExecptionDrawer from '../components/execptionDrawer'
import SourceTable from '../components/sourceTable'
import ExecptionList from '../components/execptionList'
import EventDetail from '../components/eventDetail'
import SummaryCard from '../components/summaryCard'

import styles from './styles.module.less'
import React from 'react'

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
      message.error(response?.message || '请求失败，请重试')
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
      message.error(response?.message || '请求失败，请重试')
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
        title: <NavLink to={'/cluster'}>集群管理</NavLink>,
      }
    }
    if (from === 'result') {
      first = {
        title: <NavLink to={`/search/result?query=${query}`}>搜索结果</NavLink>,
      }
    }
    const middle = []
    ;['cluster', 'kind']?.forEach(item => {
      if (urlParams?.[item]) {
        middle.push({
          key: item,
          label: urlParams?.[item],
          title: <a onClick={() => replacePage(item)}>{urlParams?.[item]}</a>,
        })
      }
    })
    middle[middle?.length - 1] = {
      label: middle[middle?.length - 1]?.lebel,
      title: middle[middle?.length - 1]?.label,
    }
    const result = [first, ...middle]
    setBreadcrumbItems(result)
  }

  useEffect(() => {
    getBreadcrumbs()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [from, key, cluster, kind, namespace, name])

  return (
    <div className={styles.container}>
      <Breadcrumb
        style={{ marginBottom: 20 }}
        separator=">"
        items={breadcrumbItems}
      />
      <ExecptionDrawer
        open={drawerVisible}
        onClose={() => setDrawerVisible(false)}
        execptionList={auditList}
        execptionStat={auditStat}
      />
      <EventDetail
        open={modalVisible}
        cancel={() => setModalVisible(false)}
        detail={currentItem}
      />
      <div className={styles.module}>
        <SummaryCard auditStat={auditStat} summary={summary} />
        <div className={styles.execption_event}>
          <ExecptionList
            auditLoading={auditLoading}
            rescan={rescan}
            execptionList={auditList}
            execptionStat={auditStat}
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
