import React, { useEffect, useState } from 'react'
import { NavLink, useLocation, useNavigate } from 'react-router-dom'
import axios from 'axios'
import queryString from 'query-string'
import { Breadcrumb, Tooltip, message } from 'antd'
import { useTranslation } from 'react-i18next'
import ExceptionDrawer from '../components/exceptionDrawer'
import TopologyMap from '../components/topologyMap'
import KarbourTabs from '@/components/tabs'
import ExceptionList from '../components/exceptionList'
import EventDetail from '../components/eventDetail'
import Yaml from '@/components/yaml'
import K8sEvent from '../components/k8sEvent'
import K8sEventDrawer from '../components/k8sEventDrawer'
import SummaryCard from '../components/summaryCard'
import { capitalized, generateResourceTopologyData } from '@/utils/tools'
import { ICON_MAP } from '@/utils/images'

import styles from './styles.module.less'

const ClusterDetail = () => {
  const { t, i18n } = useTranslation()
  const navigate = useNavigate()
  const location = useLocation()
  const urlParams = queryString.parse(location?.search)
  const { type, apiVersion, cluster, kind, namespace, name, key, from, query } =
    urlParams
  const [drawerVisible, setDrawerVisible] = useState<boolean>(false)
  const [k8sDrawerVisible, setK8sDrawerVisible] = useState<boolean>(false)
  const [currentTab, setCurrentTab] = useState('Topology')
  const [modalVisible, setModalVisible] = useState<boolean>(false)
  const [yamlData, setYamlData] = useState('')
  const [auditList, setAuditList] = useState<any>([])
  const [auditLoading, setAuditLoading] = useState<any>(false)
  const [auditStat, setAuditStat] = useState<any>()
  const [breadcrumbItems, setBreadcrumbItems] = useState([])
  const [summary, setSummary] = useState<any>()
  const [currentItem, setCurrentItem] = useState<any>()
  const [topologyData, setTopologyData] = useState<any>()
  const [topologyLoading, setTopologyLoading] = useState(false)

  const tabsList = [
    { label: t('ResourceTopology'), value: 'Topology' },
    { label: 'YAML', value: 'YAML' },
    { label: t('KubernetesEvents'), value: 'K8s', disabled: true },
  ]

  async function handleTabChange(value: string) {
    setCurrentTab(value)
  }

  async function getAudit(isRescan) {
    setAuditLoading(true)
    const response: any = await axios({
      url: `/rest-api/v1/insight/audit`,
      method: 'GET',
      params: {
        apiVersion,
        kind,
        cluster,
        namespace,
        name,
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
      url: `/rest-api/v1/insight/score`,
      method: 'GET',
      params: {
        cluster,
        apiVersion,
        kind,
        namespace,
        name,
      },
    })
    if (response?.success) {
      setAuditStat(response?.data)
    }
  }
  async function getClusterDetail() {
    const response: any = await axios({
      url: '/rest-api/v1/insight/detail',
      params: {
        cluster,
        apiVersion,
        kind,
        namespace,
        name,
        format: 'yaml',
      },
    })
    if (response?.success) {
      setYamlData(response?.data)
    }
  }

  async function getSummary() {
    const response: any = await axios({
      url: '/rest-api/v1/insight/summary',
      params: {
        cluster,
        apiVersion,
        kind,
        namespace,
        name,
      },
    })
    if (response?.success) {
      setSummary(response?.data)
    } else {
      message.error(response?.message || t('RequestFailedAndTry'))
    }
  }

  async function getTopologyData() {
    setTopologyLoading(true)
    try {
      const response: any = await axios({
        url: '/rest-api/v1/insight/topology',
        method: 'GET',
        params: {
          cluster,
          apiVersion,
          kind,
          namespace,
          name,
        },
      })
      if (response?.success) {
        const tmpData = generateResourceTopologyData(response?.data)
        setTopologyData(tmpData)
      } else {
        message.error(response?.message || t('RequestFailedAndTry'))
      }
      setTopologyLoading(false)
    } catch (error) {
      setTopologyLoading(false)
    }
  }

  useEffect(() => {
    getClusterDetail()
    getAudit(false)
    getAuditScore()
    getSummary()
    getTopologyData()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [from, key, cluster, kind, namespace, name, apiVersion])

  function rescan() {
    getAuditScore()
    getAudit(true)
  }

  function showDrawer() {
    setDrawerVisible(true)
  }

  function showK8sDrawer() {
    setK8sDrawerVisible(true)
  }

  function onItemClick(item) {
    setModalVisible(true)
    setCurrentItem(item)
  }

  function replacePage(item) {
    const obj = { from, type, apiVersion, query }
    const list = ['cluster', 'kind', 'namespace', 'name']
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
    navigate(`/insightDetail/${item}?${urlStringfyParams}`)
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
    const middle = []
    ;['cluster', 'kind', 'namespace', 'name']?.forEach(item => {
      const urlParamsItem = urlParams?.[item]
      if (urlParamsItem) {
        const iconMap = {
          cluster: ICON_MAP?.Kubernetes,
          kind: ICON_MAP?.[urlParamsItem as any] || ICON_MAP.CRD,
          namespace: ICON_MAP.Namespace,
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
            item === 'name' ? (
              <Tooltip title={capitalized(item)}>
                <a style={{ display: 'flex', alignItems: 'center' }}>
                  {urlParamsItem}
                </a>
              </Tooltip>
            ) : (
              <Tooltip title={capitalized(item)}>
                <a
                  onClick={() => replacePage(item)}
                  style={{ display: 'flex', alignItems: 'center' }}
                >
                  <img src={iconMap[item]} style={iconStyle} />
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

  async function onTopologyNodeClick(node: any) {
    const { locator } = node || {}
    const paramsObj = {
      apiVersion: locator?.apiVersion,
      cluster: locator?.cluster,
      kind: locator?.kind,
      namespace: locator?.namespace,
      name: locator?.name,
      from,
      query,
      type: 'resource',
    }
    const urlString = queryString.stringify(paramsObj)
    navigate(`/insightDetail/resource?${urlString}`, { replace: true })
  }

  return (
    <div className={styles.container}>
      <Breadcrumb
        style={{ marginBottom: 20 }}
        separator=">"
        items={breadcrumbItems}
      />
      <div className={styles.module}>
        <SummaryCard auditStat={auditStat} summary={summary} />
        <div className={styles.exception_event}>
          {/* 风险 */}
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
        <div className={styles.tab_header}>
          <KarbourTabs
            list={tabsList}
            current={currentTab}
            onChange={handleTabChange}
          />
        </div>
        {currentTab === 'Topology' && (
          <TopologyMap
            tableName={name as string}
            isResource={true}
            topologyData={topologyData}
            topologyLoading={topologyLoading}
            onTopologyNodeClick={onTopologyNodeClick}
          />
        )}
        {currentTab === 'YAML' && <Yaml data={yamlData} />}
        {currentTab === 'K8s' && (
          <K8sEvent
            rescan={rescan}
            exceptionList={[1, 2, 3, 4, 5]}
            showDrawer={showK8sDrawer}
            onItemClick={onItemClick}
          />
        )}
      </div>
      <ExceptionDrawer
        open={drawerVisible}
        onClose={() => setDrawerVisible(false)}
        exceptionList={auditList}
        exceptionStat={auditStat}
      />
      <K8sEventDrawer
        open={k8sDrawerVisible}
        onClose={() => setK8sDrawerVisible(false)}
      />
      <EventDetail
        open={modalVisible}
        cancel={() => setModalVisible(false)}
        detail={currentItem}
      />
    </div>
  )
}

export default ClusterDetail
