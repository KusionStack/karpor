import React, { useEffect, useRef, useState } from 'react'
import { NavLink, useLocation } from 'react-router-dom'
import queryString from 'query-string'
import { Breadcrumb, Tooltip } from 'antd'
import { useTranslation } from 'react-i18next'
import KarporTabs from '@/components/tabs'
import Yaml from '@/components/yaml'
import { capitalized, generateTopologyData } from '@/utils/tools'
import { insightTabsList } from '@/utils/constants'
import Kubernetes from '@/assets/kubernetes.png'
import SourceTable from '../components/sourceTable'
import ExceptionList from '../components/exceptionList'
import EventDetail from '../components/eventDetail'
import ExceptionDrawer from '../components/exceptionDrawer'
import TopologyMap from '../components/topologyMap'
import K8sEvent from '../components/k8sEvent'
import K8sEventDrawer from '../components/k8sEventDrawer'
import SummaryCard from '../components/summaryCard'
import { useAxios, isHighAvailability } from '@/utils/request'

import styles from './styles.module.less'

const ClusterDetail = () => {
  const location = useLocation()
  const { t, i18n } = useTranslation()
  const urlParams = queryString.parse(location?.search)
  const { type, cluster, kind, namespace, name, key, from, query, apiVersion } =
    urlParams
  const [drawerVisible, setDrawerVisible] = useState<boolean>(false)
  const [k8sDrawerVisible, setK8sDrawerVisible] = useState<boolean>(false)
  const [currentTab, setCurrentTab] = useState('Topology')
  const [modalVisible, setModalVisible] = useState<boolean>(false)
  const [tableQueryStr, setTableQueryStr] = useState('')
  const [yamlData, setYamlData] = useState('')
  const [highAvailabilityYamlData, setHighAvailabilityYamlData] = useState('')
  const [auditList, setAuditList] = useState<any>([])
  const [auditStat, setAuditStat] = useState<any>()
  const [tableName, setTableName] = useState('Pod')
  const [breadcrumbItems, setBreadcrumbItems] = useState([])
  const [summary, setSummary] = useState<any>()
  const [currentItem, setCurrentItem] = useState<any>()
  const [multiTopologyData, setMultiTopologyData] = useState<any>()
  const [selectedCluster, setSelectedCluster] = useState<any>()
  const [clusterOptions, setClusterOptions] = useState<string[]>([])
  const [tabList, setTabList] = useState(insightTabsList)

  const drawRef = useRef(null)

  useEffect(() => {
    if (selectedCluster) {
      const result = generateTopologyData(multiTopologyData?.[selectedCluster])
      const tmp = result?.nodes?.find((item: any) => {
        if (item?.id) {
          const kindTmp = item?.id?.split('.')
          const len = kindTmp?.length
          const lastKindTmp = kindTmp?.[len - 1]
          return lastKindTmp === 'Pod'
        } else {
          return false
        }
      })
      if (tmp) {
        const tmpData = tmp?.id?.split('.')
        const len = tmpData?.length
        const kindVersion = tmpData?.[len - 2]
        const queryStr = `select * from resources where cluster = '${cluster}' and apiVersion = '${kindVersion}' and kind = 'Pod'`
        setTableQueryStr(queryStr)
      }
    }
  }, [selectedCluster, cluster, apiVersion, multiTopologyData])

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
          cluster,
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
        },
      },
    })
  }

  const { response: clusterDetailResponse, refetch: clusterDetailRefetch } =
    useAxios({
      url: `/rest-api/v1/cluster/${cluster}`,
      method: 'GET',
    })

  const {
    response: highAvailabilityClusterDetailResponse,
    refetch: highAvailabilityClusterDetailRefetch,
  } = useAxios({
    url: `/rest-api/v1/cluster/${cluster}/agentYml`,
    method: 'GET',
  })

  useEffect(() => {
    if (clusterDetailResponse?.success) {
      setYamlData(clusterDetailResponse?.data)
    }
    if (highAvailabilityClusterDetailResponse?.success) {
      setHighAvailabilityYamlData(highAvailabilityClusterDetailResponse?.data)
    }
  }, [clusterDetailResponse])

  function getClusterDetail() {
    clusterDetailRefetch({
      url: `/rest-api/v1/cluster/${cluster}`,
      option: {
        params: {
          format: 'yaml',
        },
      },
    })
  }

  function getHighAvailabilityClusterDetail() {
    highAvailabilityClusterDetailRefetch({
      url: `/rest-api/v1/cluster/${cluster}/agentYml`,
      option: {
        params: {
          format: 'yaml',
        },
      },
    })
  }

  const { response: summaryResponse, refetch: summaryRefetch } = useAxios({
    url: '/rest-api/v1/insight/summary',
    method: 'GET',
    manual: true,
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
    manual: true,
  })

  useEffect(() => {
    if (topologyDataResponse?.success) {
      setMultiTopologyData(topologyDataResponse?.data)
    }
  }, [topologyDataResponse])

  function getTopologyData() {
    topologyDataRefetch({
      option: {
        params: {
          cluster,
        },
      },
    })
  }

  useEffect(() => {
    if (multiTopologyData) {
      const clusterKeys = Object.keys(multiTopologyData)
      setClusterOptions(clusterKeys)
      setSelectedCluster(clusterKeys?.[0])
    }
  }, [multiTopologyData])

  useEffect(() => {
    if (selectedCluster && currentTab === 'Topology') {
      const topologyData =
        multiTopologyData &&
        selectedCluster &&
        generateTopologyData(multiTopologyData?.[selectedCluster])
      drawRef.current?.drawGraph(topologyData)
    }
  }, [multiTopologyData, selectedCluster, currentTab])

  useEffect(() => {
    if (isHighAvailability) {
      getHighAvailabilityClusterDetail()
    } else {
      getClusterDetail()
    }

    getAudit(false)
    getAuditScore()
    getSummary()
    getTopologyData()
    if (type === 'kind' && kind) {
      setTableName(kind as any)
    }
    if (isHighAvailability) {
      const initialTabList = [...insightTabsList]
      if (!initialTabList.find(tab => tab.value === 'AgentYaml')) {
        initialTabList.push({ value: 'AgentYaml', label: 'Agent Yaml' })
      }
      setTabList(initialTabList)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [kind, type])

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

  function onTopologyNodeClick(node) {
    const { resourceGroup } = node?.data || {}
    setTableName(resourceGroup?.kind)
    const sqlStr = `select * from resources where cluster = '${resourceGroup?.cluster}' and apiVersion = '${resourceGroup?.apiVersion}' and kind = '${resourceGroup?.kind}'`
    setTableQueryStr(sqlStr)
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
    const result = [
      first,
      {
        key: cluster,
        title: (
          <Tooltip title={capitalized('cluster')}>
            <div style={{ display: 'flex', alignItems: 'center' }}>
              <img
                src={Kubernetes}
                style={{ width: 14, height: 14, marginRight: 2 }}
              />
              {cluster}
            </div>
          </Tooltip>
        ),
      },
    ]
    setBreadcrumbItems(result)
  }

  useEffect(() => {
    getBreadcrumbs()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [from, key, cluster, kind, namespace, name, i18n?.language])

  function handleChangeCluster(val) {
    setSelectedCluster(val)
  }

  function renderTabPane() {
    if (currentTab === 'Topology') {
      const topologyData =
        multiTopologyData &&
        selectedCluster &&
        generateTopologyData(multiTopologyData?.[selectedCluster])
      if (topologyData?.nodes?.length > 0) {
        return (
          <>
            <TopologyMap
              ref={drawRef}
              tableName={tableName}
              selectedCluster={selectedCluster}
              handleChangeCluster={handleChangeCluster}
              clusterOptions={clusterOptions}
              topologyLoading={topologyLoading}
              onTopologyNodeClick={onTopologyNodeClick}
            />
            <SourceTable queryStr={tableQueryStr} tableName={tableName} />
          </>
        )
      }
    }
    if (currentTab === 'YAML') {
      return <Yaml data={yamlData || ''} />
    }
    if (currentTab === 'AgentYaml') {
      return <Yaml data={highAvailabilityYamlData || ''} />
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
            list={tabList?.filter(item => item?.value !== 'Events')}
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
