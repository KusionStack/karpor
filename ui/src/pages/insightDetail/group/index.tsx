import React, { useEffect, useState } from 'react'
import { NavLink, useLocation } from 'react-router-dom'
import axios from 'axios'
import queryString from 'query-string'
import { Breadcrumb, Tooltip, message } from 'antd'
import { useTranslation } from 'react-i18next'
import ExceptionDrawer from '../components/exceptionDrawer'
import TopologyMap from '../components/topologyMap'
import KarporTabs from '@/components/tabs'
import SourceTable from '../components/sourceTable'
import ExceptionList from '../components/exceptionList'
import EventDetail from '../components/eventDetail'
import Yaml from '@/components/yaml'
import K8sEvent from '../components/k8sEvent'
import K8sEventDrawer from '../components/k8sEventDrawer'
import SummaryCard from '../components/summaryCard'
import { generateTopologyData } from '@/utils/tools'
import { insightTabsList } from '@/utils/constants'

import styles from './styles.module.less'

const ClusterDetail = () => {
  const location = useLocation()
  const { t, i18n } = useTranslation()
  const urlParams = queryString.parse(location?.search)
  const resultUrlParams: any = getUrlParams()

  const [drawerVisible, setDrawerVisible] = useState<boolean>(false)
  const [k8sDrawerVisible, setK8sDrawerVisible] = useState<boolean>(false)
  const [currentTab, setCurrentTab] = useState('Topology')
  const [modalVisible, setModalVisible] = useState<boolean>(false)
  const [tableQueryStr, setTableQueryStr] = useState('')
  const [yamlData, setYamlData] = useState('')
  const [auditList, setAuditList] = useState<any>([])
  const [auditLoading, setAuditLoading] = useState<any>(false)
  const [auditStat, setAuditStat] = useState<any>()
  const [tableName, setTableName] = useState('Pod')
  const [breadcrumbItems, setBreadcrumbItems] = useState([])
  const [summary, setSummary] = useState<any>()
  const [currentItem, setCurrentItem] = useState<any>()
  const [multiTopologyData, setMultiTopologyData] = useState<any>()
  const [topologyLoading, setTopologyLoading] = useState(false)
  const [selectedCluster, setSelectedCluster] = useState<any>()
  const [clusterOptions, setClusterOptions] = useState<string[]>([])
  const [groupTabsList, setGroupTabsList] = useState(insightTabsList)

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
      const result = generateTopologyData(multiTopologyData?.[selectedCluster])
      const tmp = result?.nodes?.find((item: any) => {
        if (item?.id) {
          const kindTmp = item?.id?.split('.')
          const len = kindTmp?.length
          const lastKindTmp = kindTmp?.[len - 1]
          if (lastKindTmp === 'Pod') {
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
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    selectedCluster,
    resultUrlParams?.cluster,
    resultUrlParams?.apiVersion,
    multiTopologyData,
  ])

  async function handleTabChange(value: string) {
    setCurrentTab(value)
  }

  async function getAudit(isRescan) {
    setAuditLoading(true)
    const response: any = await axios({
      url: `/rest-api/v1/insight/audit`,
      method: 'GET',
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
    })
    if (response?.success) {
      setAuditStat(response?.data)
    }
  }
  async function getClusterDetail() {
    const response: any = await axios({
      url: '/rest-api/v1/insight/detail',
      method: 'GET',
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
        format: 'yaml',
      },
    })
    if (response?.success) {
      setYamlData(response?.data)
    } else {
      const newGroupTabList = groupTabsList?.filter(
        item => item?.value !== 'YAML',
      )
      setGroupTabsList(newGroupTabList)
    }
  }

  async function getSummary() {
    const response: any = await axios({
      url: '/rest-api/v1/insight/summary',
      method: 'GET',
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
    })
    if (response?.success) {
      setSummary(response?.data)
    } else {
      message.error(response?.message || t('RequestFailedAndTry'))
    }
  }

  function getObjectUrlParams(obj) {
    const list = Object.entries(obj)?.map(([k, v]) => `${k}=${v}`)
    return list?.join(',')
  }

  async function getTopologyData() {
    setTopologyLoading(true)
    try {
      const response: any = await axios({
        url: '/rest-api/v1/insight/topology',
        method: 'GET',
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
      })
      if (response?.success) {
        setMultiTopologyData(response?.data)
      } else {
        message.error(response?.message || t('RequestFailedAndTry'))
      }
      setTopologyLoading(false)
    } catch (error) {
      setTopologyLoading(false)
    }
  }

  useEffect(() => {
    if (multiTopologyData) {
      const clusterKeys = Object.keys(multiTopologyData)
      setClusterOptions(clusterKeys)
      setSelectedCluster(clusterKeys?.[0])
    }
  }, [multiTopologyData])

  useEffect(() => {
    getClusterDetail()
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

  function showK8sDrawer() {
    setK8sDrawerVisible(true)
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

  async function onTopologyNodeClick(node) {
    const { resourceGroup } = node?.data || {}
    setTableName(resourceGroup?.kind)
    const sqlParams = generateSqlParams({
      ...resourceGroup,
      cluster: resultUrlParams?.cluster || selectedCluster,
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
              tableName={tableName}
              selectedCluster={selectedCluster}
              handleChangeCluster={handleChangeCluster}
              clusterOptions={clusterOptions}
              topologyData={topologyData}
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
            list={groupTabsList}
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
