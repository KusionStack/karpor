import React, { useEffect, useRef, useState } from 'react'
import { NavLink, useLocation } from 'react-router-dom'
import queryString from 'query-string'
import { Breadcrumb, Tooltip } from 'antd'
import { useTranslation } from 'react-i18next'
import { generateTopologyData } from '@/utils/tools'
import { insightTabsList } from '@/utils/constants'
import KarporTabs from '@/components/tabs'
import { useAxios } from '@/utils/request'
import ExceptionDrawer from '../components/exceptionDrawer'
import TopologyMap from '../components/topologyMap'
import SourceTable from '../components/sourceTable'
import ExceptionList from '../components/exceptionList'
import EventDetail from '../components/eventDetail'
import SummaryCard from '../components/summaryCard'

import styles from './styles.module.less'

const GroupDetail = () => {
  const location = useLocation()
  const { i18n } = useTranslation()
  const urlParams = queryString.parse(location?.search)
  const resultUrlParams: any = getUrlParams()

  const [drawerVisible, setDrawerVisible] = useState<boolean>(false)
  const [currentTab, setCurrentTab] = useState('Topology')
  const [modalVisible, setModalVisible] = useState<boolean>(false)
  const [tableQueryStr, setTableQueryStr] = useState('')
  const [auditList, setAuditList] = useState<any>([])
  const [auditStat, setAuditStat] = useState<any>()
  const [tableName, setTableName] = useState('Pod')
  const [breadcrumbItems, setBreadcrumbItems] = useState([])
  const [summary, setSummary] = useState<any>()
  const [currentItem, setCurrentItem] = useState<any>()
  const [multiTopologyData, setMultiTopologyData] = useState<any>()
  const [selectedCluster, setSelectedCluster] = useState<any>()
  const [clusterOptions, setClusterOptions] = useState<string[]>([])

  const drawRef = useRef(null)

  function getUrlParams() {
    const obj = {}
    Object.keys(urlParams)?.forEach(key => {
      if (key?.startsWith('labels') || key?.startsWith('annotations')) {
        const [first, last]: any = key?.split('__')
        if (obj?.[first]) {
          obj[first][last] = urlParams?.[key]
        } else {
          obj[first] = {}
          obj[first][last] = urlParams?.[key]
        }
      } else {
        obj[key] = urlParams?.[key]
      }
    })
    return obj
  }

  function generateSqlParams(sqlParamsObj = {}) {
    return Object.entries(sqlParamsObj)
      ?.map(([k, v]) => ` ${k} = '${v || ''}' `)
      ?.join('and')
  }

  useEffect(() => {
    if (selectedCluster) {
      if (selectedCluster === 'ALL') {
        const result = generateAllClusterTopologyData()
        const tmp: any = result?.nodes?.find((item: any) => {
          if (item?.id) {
            const kindTmp = item?.id?.split('.')
            const len = kindTmp?.length
            const lastKindTmp = kindTmp?.[len - 1]
            if (lastKindTmp === tableName) {
              return true
            } else {
              return false
            }
          } else {
            return false
          }
        })
        if (tmp) {
          const sqlParams = generateSqlParams({
            ...tmp?.data?.resourceGroup,
            ...generateUrlSqlParams(),
          })
          const queryStr = `select * from resources where ${sqlParams} `
          setTableQueryStr(queryStr)
        }
      } else {
        const result = generateTopologyData(
          multiTopologyData?.[selectedCluster],
        )
        const tmp = result?.nodes?.find((item: any) => {
          if (item?.id) {
            const kindTmp = item?.id?.split('.')
            const len = kindTmp?.length
            const lastKindTmp = kindTmp?.[len - 1]
            if (lastKindTmp === tableName) {
              return true
            } else {
              return false
            }
          } else {
            return false
          }
        })
        if (tmp) {
          const sqlParams = generateSqlParams({
            ...tmp?.data?.resourceGroup,
            cluster: resultUrlParams?.cluster || selectedCluster,
            ...generateUrlSqlParams(),
          })
          const queryStr = `select * from resources where ${sqlParams} `
          setTableQueryStr(queryStr)
        }
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    selectedCluster,
    resultUrlParams?.cluster,
    resultUrlParams?.apiVersion,
    multiTopologyData,
  ])

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
          apiVersion: resultUrlParams?.apiVersion,
          cluster: resultUrlParams?.cluster,
          kind: resultUrlParams?.kind,
          namespace: resultUrlParams?.namespace,
          name: resultUrlParams?.name,
          labels: resultUrlParams?.labels
            ? getObjectUrlParams(resultUrlParams?.labels)
            : undefined,
          annotations: resultUrlParams?.annotations
            ? getObjectUrlParams(resultUrlParams?.annotations)
            : undefined,
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
          apiVersion: resultUrlParams?.apiVersion,
          cluster: resultUrlParams?.cluster,
          kind: resultUrlParams?.kind,
          namespace: resultUrlParams?.namespace,
          name: resultUrlParams?.name,
          labels: resultUrlParams?.labels
            ? getObjectUrlParams(resultUrlParams?.labels)
            : undefined,
          annotations: resultUrlParams?.annotations
            ? getObjectUrlParams(resultUrlParams?.annotations)
            : undefined,
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
          apiVersion: resultUrlParams?.apiVersion,
          cluster: resultUrlParams?.cluster,
          kind: resultUrlParams?.kind,
          namespace: resultUrlParams?.namespace,
          name: resultUrlParams?.name,
          labels: resultUrlParams?.labels
            ? getObjectUrlParams(resultUrlParams?.labels)
            : undefined,
          annotations: resultUrlParams?.annotations
            ? getObjectUrlParams(resultUrlParams?.annotations)
            : undefined,
        },
      },
    })
  }

  function getObjectUrlParams(obj) {
    const list = Object.entries(obj)?.map(([k, v]) => `${k}=${v}`)
    return list?.join(',')
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
          apiVersion: resultUrlParams?.apiVersion,
          cluster: resultUrlParams?.cluster,
          kind: resultUrlParams?.kind,
          namespace: resultUrlParams?.namespace,
          name: resultUrlParams?.name,
          labels: resultUrlParams?.labels
            ? getObjectUrlParams(resultUrlParams?.labels)
            : undefined,
          annotations: resultUrlParams?.annotations
            ? getObjectUrlParams(resultUrlParams?.annotations)
            : undefined,
        },
      },
    })
  }

  useEffect(() => {
    if (multiTopologyData) {
      const clusterKeys = Object.keys(multiTopologyData)
      setClusterOptions(['ALL', ...clusterKeys])
      setSelectedCluster(selectedCluster || 'ALL')
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [multiTopologyData])

  useEffect(() => {
    getAudit(false)
    getAuditScore()
    getSummary()
    getTopologyData()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [resultUrlParams?.kind])

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

  function generateUrlSqlParams() {
    const { labels, annotations } = resultUrlParams || {}
    const urlSqlParams = {}
    if (labels) {
      Object.entries(labels)?.forEach(([k, v]) => {
        const key = '`' + `labels.${k}` + '`'
        urlSqlParams[key] = v
      })
    }
    if (annotations) {
      Object.entries(annotations)?.forEach(([k, v]) => {
        const key = '`' + `annotations.${k}` + '`'
        urlSqlParams[key] = v
      })
    }
    return urlSqlParams
  }

  const selectedClusterRef = useRef()
  selectedClusterRef.current = selectedCluster
  function onTopologyNodeClick(node) {
    const { resourceGroup } = node?.data || {}
    setTableName(resourceGroup?.kind)
    const sqlParams = generateSqlParams({
      ...resourceGroup,
      ...(selectedClusterRef?.current === 'ALL'
        ? {}
        : { cluster: selectedClusterRef?.current || resultUrlParams?.cluster }),
      ...generateUrlSqlParams(),
    })
    const sqlStr = `select * from resources where ${sqlParams}`
    setTableQueryStr(sqlStr)
  }

  function getBreadcrumbs() {
    const first = {
      title: (
        <NavLink to={`/insight?activeTabKey=${resultUrlParams?.ruleKey}`}>
          {resultUrlParams?.ruleKey}
        </NavLink>
      ),
    }

    const result = [
      first,
      {
        key: resultUrlParams?.title,
        title: (
          <Tooltip title={resultUrlParams?.title}>
            <div style={{ display: 'flex', alignItems: 'center' }}>
              {resultUrlParams?.title}
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
  }, [
    resultUrlParams?.from,
    resultUrlParams?.cluster,
    resultUrlParams?.kind,
    resultUrlParams?.namespace,
    resultUrlParams?.name,
    i18n?.language,
  ])

  function handleChangeCluster(val) {
    setSelectedCluster(val)
  }

  function generateAllClusterTopologyData() {
    const objNodes = {}
    const objEdges = {}
    if (multiTopologyData) {
      Object.values(multiTopologyData)?.forEach(clusterData => {
        const currentClusterData = generateTopologyData(clusterData)
        if (currentClusterData) {
          currentClusterData?.nodes?.forEach(current => {
            if (!objNodes?.[current?.id]) {
              objNodes[current?.id] = current
            } else {
              objNodes[current?.id].data.count += current?.data?.count
            }
          })
          currentClusterData?.edges?.forEach(current => {
            if (!objEdges?.[`${current?.source}__${current?.target}`]) {
              objEdges[`${current?.source}__${current?.target}`] = current
            }
          })
        }
      })
    }
    const nodes = Object.values(objNodes)
    const edges = Object.values(objEdges)
    return {
      nodes,
      edges,
    }
  }

  useEffect(() => {
    let topologyData
    if (selectedCluster) {
      if (selectedCluster === 'ALL' && multiTopologyData) {
        topologyData = generateAllClusterTopologyData()
      } else {
        topologyData =
          multiTopologyData &&
          selectedCluster &&
          generateTopologyData(multiTopologyData?.[selectedCluster])
      }
      drawRef.current?.drawGraph(topologyData)
    }
  }, [multiTopologyData, selectedCluster])

  function renderTabPane() {
    if (currentTab === 'Topology') {
      let topologyData
      if (selectedCluster === 'ALL') {
        topologyData = generateAllClusterTopologyData()
      } else {
        topologyData =
          multiTopologyData &&
          selectedCluster &&
          generateTopologyData(multiTopologyData?.[selectedCluster])
      }
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
            list={insightTabsList?.filter(
              item => item?.value !== 'YAML' && item?.value !== 'Events',
            )}
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
      <EventDetail
        open={modalVisible}
        cancel={() => setModalVisible(false)}
        detail={currentItem}
      />
    </div>
  )
}

export default GroupDetail
