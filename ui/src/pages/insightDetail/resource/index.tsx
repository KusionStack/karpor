import React, { useEffect, useRef, useState } from 'react'
import { NavLink, useLocation, useNavigate } from 'react-router-dom'
import queryString from 'query-string'
import { Breadcrumb, Tooltip } from 'antd'
import { useTranslation } from 'react-i18next'
import KarporTabs from '@/components/tabs'
import Yaml from '@/components/yaml'
import { capitalized, generateResourceTopologyData } from '@/utils/tools'
import { insightTabsList } from '@/utils/constants'
import { ICON_MAP } from '@/utils/images'
import { useAxios } from '@/utils/request'
import ExceptionDrawer from '../components/exceptionDrawer'
import TopologyMap from '../components/topologyMap'
import ExceptionList from '../components/exceptionList'
import EventDetail from '../components/eventDetail'
import K8sEvent from '../components/k8sEvent'
import K8sEventDrawer from '../components/k8sEventDrawer'
import SummaryCard from '../components/summaryCard'
import PodLogs from '../components/podLogs'
import EventAggregator from '../components/eventAggregator'

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
  const [currentTab, setCurrentTab] = useState(
    urlParams?.deleted === 'true' ? 'YAML' : 'Topology',
  )
  const [modalVisible, setModalVisible] = useState<boolean>(false)
  const [yamlData, setYamlData] = useState('')
  const [auditList, setAuditList] = useState<any>([])
  const [auditStat, setAuditStat] = useState<any>()
  const [breadcrumbItems, setBreadcrumbItems] = useState([])
  const [summary, setSummary] = useState<any>()
  const [currentItem, setCurrentItem] = useState<any>()
  const [multiTopologyData, setMultiTopologyData] = useState<any>()
  const [selectedCluster, setSelectedCluster] = useState<any>()
  const [clusterOptions, setClusterOptions] = useState<string[]>([])

  const [tabList, setTabList] = useState(insightTabsList)

  const drawRef = useRef(null)

  useEffect(() => {
    const initialTabList = [...insightTabsList]
    if (kind === 'Pod') {
      if (!initialTabList.find(tab => tab.value === 'Logs')) {
        initialTabList.push({ value: 'Logs', label: 'LogAggregator' })
      }
    }
    setTabList(initialTabList)
  }, [kind])

  useEffect(() => {
    if (urlParams?.deleted === 'true') {
      setTabList(prev =>
        prev.map(item => ({
          ...item,
          disabled: item.value === 'Topology' && urlParams?.deleted === 'true',
        })),
      )
    } else {
      setTabList(prev =>
        prev.map(item => ({
          ...item,
          disabled: false,
        })),
      )
    }
  }, [urlParams?.deleted])

  function handleTabChange(value: string) {
    setCurrentTab(value)
  }

  const {
    response: auditResponse,
    refetch: auditRefetch,
    loading: auditLoading,
  } = useAxios({
    url: `/rest-api/v1/insight/audit`,
  })

  useEffect(() => {
    if (auditResponse?.success) {
      setAuditList(auditResponse?.data)
    }
  }, [auditResponse])

  function getAudit(isRescan) {
    auditRefetch({
      option: {
        params: {
          apiVersion,
          kind,
          cluster,
          namespace,
          name,
          ...(isRescan ? { forceNew: true } : {}),
        },
      },
    })
  }

  const { response: auditScoreResponse, refetch: auditScoreRefetch } = useAxios(
    {
      url: '/rest-api/v1/insight/score',
    },
  )

  useEffect(() => {
    if (auditScoreResponse?.success) {
      setAuditStat(auditScoreResponse?.data)
    }
  }, [auditScoreResponse])

  function getAuditScore() {
    auditScoreRefetch({
      option: {
        params: {
          cluster,
          apiVersion,
          kind,
          namespace,
          name,
        },
      },
    })
  }

  const { response: clusterDetailResponse, refetch: clusterDetailRefetch } =
    useAxios({
      url: '/rest-api/v1/insight/score',
      method: 'GET',
    })

  useEffect(() => {
    if (clusterDetailResponse?.success) {
      setYamlData(clusterDetailResponse?.data)
    }
  }, [clusterDetailResponse])

  function getClusterDetail() {
    clusterDetailRefetch({
      url: '/rest-api/v1/insight/detail',
      option: {
        params: {
          cluster,
          apiVersion,
          kind,
          namespace,
          name,
          format: 'yaml',
        },
      },
    })
  }

  const { response: summaryResponse, refetch: summaryRefetch } = useAxios({
    url: '/rest-api/v1/insight/summary',
    method: 'GET',
  })

  useEffect(() => {
    if (summaryResponse?.success) {
      setSummary(summaryResponse?.data)
    }
  }, [summaryResponse])

  function getSummary() {
    summaryRefetch({
      option: {
        params: {
          cluster,
          apiVersion,
          kind,
          namespace,
          name,
        },
      },
    })
  }

  const {
    response: topologyDataResponse,
    refetch: topologyDataRefetch,
    loading: topologyLoading,
  } = useAxios({
    url: '/rest-api/v1/insight/topology',
    method: 'GET',
  })

  useEffect(() => {
    if (topologyDataResponse?.success) {
      const data = topologyDataResponse?.data
      setMultiTopologyData(data)
    }
  }, [topologyDataResponse])

  function getTopologyData() {
    if (urlParams?.deleted === 'true') return
    topologyDataRefetch({
      option: {
        params: {
          cluster,
          apiVersion,
          kind,
          namespace,
          name,
        },
      },
    })
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
    if (from === 'insight') {
      first = {
        title: <NavLink to={'/insight'}>{t('Insight')}</NavLink>,
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

  function onTopologyNodeClick(node: any) {
    const {
      data: { resourceGroup },
    } = node || {}
    const paramsObj = {
      apiVersion: resourceGroup?.apiVersion,
      cluster: resourceGroup?.cluster,
      kind: resourceGroup?.kind,
      namespace: resourceGroup?.namespace,
      name: resourceGroup?.name,
      from,
      query,
      type: 'resource',
    }
    const urlString = queryString.stringify(paramsObj)
    navigate(`/insightDetail/resource?${urlString}`, { replace: true })
  }

  useEffect(() => {
    if (multiTopologyData) {
      const clusterKeys = Object.keys(multiTopologyData)
      setClusterOptions(clusterKeys)
      setSelectedCluster(clusterKeys?.[0])
    }
  }, [multiTopologyData])

  function handleChangeCluster(val) {
    setSelectedCluster(val)
  }

  useEffect(() => {
    if (selectedCluster && currentTab === 'Topology') {
      const topologyData =
        multiTopologyData &&
        selectedCluster &&
        generateResourceTopologyData(multiTopologyData?.[selectedCluster])
      drawRef.current?.drawGraph(topologyData)
    }
  }, [multiTopologyData, selectedCluster, currentTab])

  function renderTabPane() {
    if (currentTab === 'Topology') {
      const topologyData =
        multiTopologyData &&
        selectedCluster &&
        generateResourceTopologyData(multiTopologyData?.[selectedCluster])
      if (topologyData?.nodes?.length > 0) {
        return (
          <TopologyMap
            ref={drawRef}
            tableName={name as string}
            isResource={true}
            selectedCluster={selectedCluster}
            handleChangeCluster={handleChangeCluster}
            clusterOptions={clusterOptions}
            topologyLoading={topologyLoading}
            onTopologyNodeClick={onTopologyNodeClick}
          />
        )
      }
    }
    if (currentTab === 'YAML') {
      return <Yaml data={yamlData || ''} />
    }
    if (currentTab === 'Events') {
      return (
        <EventAggregator
          cluster={cluster as string}
          namespace={namespace as string}
          name={name as string}
          kind={kind as string}
          apiVersion={apiVersion as string}
        />
      )
    }
    if (currentTab === 'Logs' && kind === 'Pod') {
      return (
        <PodLogs
          cluster={cluster as string}
          namespace={namespace as string}
          podName={name as string}
          yamlData={yamlData}
        />
      )
    }
    if (currentTab === 'K8s') {
      return (
        <K8sEvent
          rescan={rescan}
          exceptionList={[1, 2, 3, 4, 5]}
          showDrawer={showK8sDrawer}
          onItemClick={onItemClick}
        />
      )
    }
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

      <div className={styles.tab_content}>
        <div className={styles.tab_header}>
          <KarporTabs
            list={tabList}
            current={currentTab}
            onChange={handleTabChange}
          />
        </div>
        {renderTabPane()}
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
