import React, { useEffect, useState } from 'react'
import { Empty, Button, Input, message } from 'antd'
import axios from 'axios'
import {
  SearchOutlined,
  SortAscendingOutlined,
  SortDescendingOutlined,
  ArrowDownOutlined,
  ArrowUpOutlined,
} from '@ant-design/icons'
import { useSelector } from 'react-redux'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import KarporTabs from '@/components/tabs'
import Loading from '@/components/loading'
import ClusterCard from './components/clusterCard'
import healthPng from '@/assets/health_green.png'
import exceptionalPng from '@/assets/exceptional.png'
import clusterPng from '@/assets/cluster_outlind.png'

import styles from './styles.module.less'

const Cluster = () => {
  const navigate = useNavigate()
  const { isReadOnlyMode } = useSelector((state: any) => state.globalSlice)
  const [pageData, setPageData] = useState<any>([])
  const [showPageData, setShowPageData] = useState<any>([])
  const [loading, setloading] = useState(false)
  const [summary, setSummary] = useState<any>()
  const [sortParams, setSortParams] = useState<any>({
    orderBy: 'name',
    isAsc: true,
  })
  const { t } = useTranslation()
  const [searchValue, setSearchValue] = useState('')
  const [lastDetail, setLastDetail] = useState<any>()
  const [scale, setScale] = useState<any>(1)

  useEffect(() => {
    const handleResize = () => {
      const innerWidth = window.innerWidth
      if (innerWidth >= 1440) {
        setScale(1.2)
      } else {
        setScale(1)
      }
    }
    handleResize()
    window.addEventListener('resize', handleResize)
    return () => {
      window.removeEventListener('resize', handleResize)
    }
  }, [])

  async function getClusterSummary() {
    setloading(true)
    const response: any = await axios(`/rest-api/v1/clusters?summary=true`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      params: {},
    })
    if (response?.success) {
      setSummary(response?.data)
    } else {
      message.error(response?.message || '请求失败，请重试')
    }
    setloading(false)
  }

  async function getPageData(params) {
    const response: any = await axios(`/rest-api/v1/clusters`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      params: {
        orderBy: params?.orderBy,
        ...(params?.isAsc ? { ascending: true } : { descending: true }),
      },
    })
    if (response?.success) {
      setPageData(response?.data?.items)
    } else {
      message.error(response?.message || '请求失败，请重试')
    }
  }

  function getShowPageData(allData, currentTabKey) {
    let result: any
    if (currentTabKey === 'all') {
      result = allData
    } else {
      const exceptionList = []
      const healthyList = []
      allData?.forEach(item => {
        if (summary?.unhealthyClusters?.includes(item?.metadata?.name)) {
          exceptionList.push(item)
        }
        if (summary?.healthyClusters?.includes(item?.metadata?.name)) {
          healthyList.push(item)
        }
      })
      result = currentTabKey === 'healthy' ? healthyList : exceptionList
    }
    return result
  }

  useState(() => {
    getClusterSummary()
    getPageData(sortParams)
  })

  const join = () => {
    if (isReadOnlyMode) return
    navigate('/cluster/access')
  }

  async function handleSubmit(values, callback: () => void) {
    if (isReadOnlyMode) return
    const response: any = await axios({
      url: `/rest-api/v1/cluster/${lastDetail?.metadata?.name}`,
      method: 'PUT',
      data: values,
    })
    if (response?.success) {
      message.success(t('UpdateSuccess'))
      callback()
      getClusterSummary()
      getPageData(sortParams)
    } else {
      message.error(response?.message || t('RequestFailedAndTry'))
    }
  }

  const [currentTab, setCurrentTab] = useState('all')
  const [triangleLeftOffestIndex, setTriangleLeftOffestIndex] = useState(0)

  function handleTabChange(value: string, index: number) {
    setTriangleLeftOffestIndex(index)
    setCurrentTab(value)
    const res = getShowPageData(pageData, value)
    setShowPageData(res)
  }

  useEffect(() => {
    const res = getShowPageData(pageData, currentTab)
    if (!searchValue) {
      setShowPageData(res)
    } else {
      const newValue = searchValue?.toLowerCase().trim()?.split(' ')
      const newShowPageData = []
      if (newValue?.length === 1) {
        res?.forEach((item: any) => {
          if (item?.metadata?.name?.toLowerCase()?.includes(newValue?.[0])) {
            newShowPageData.push(item)
          }
        })
      } else {
        res?.forEach((item: any) => {
          if (
            newValue?.every((innerValue: string) =>
              item?.metadata?.name?.toLowerCase()?.includes(innerValue),
            )
          ) {
            newShowPageData.push(item)
          }
        })
      }
      setShowPageData(newShowPageData)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchValue, pageData])

  const iconStyle: React.CSSProperties = {
    width: 24,
    height: 24,
    display: 'flex',
    justifyContent: 'center',
    lineHeight: '24px',
    alignItems: 'center',
    marginRight: 10,
    borderRadius: 6,
  }

  const tabStyle: React.CSSProperties = {
    display: 'flex',
    alignItems: 'center',
    fontSize: 14,
    fontWeight: 400,
  }

  const numberStyle = { paddingLeft: 10, fontSize: 24 }

  async function deleteItem(item) {
    if (isReadOnlyMode) return
    const response: any = await axios({
      url: `/rest-api/v1/cluster/${item?.metadata?.name}`,
      method: 'DELETE',
    })
    if (response?.success) {
      message.success(t('DeletedSuccess'))
      getPageData(sortParams)
      getClusterSummary()
    } else {
      message.error(response?.message || t('RequestFailedAndTry'))
    }
  }

  function goCertificate(item) {
    if (isReadOnlyMode) {
      return
    }
    navigate(
      `/cluster/certificate?cluster=${item?.metadata?.name}&apiVersion=${item?.apiVersion}`,
    )
  }

  function goDetailPage(item) {
    navigate(
      `/insightDetail/cluster?from=cluster&type=cluster&cluster=${item?.metadata?.name}&apiVersion=${item?.apiVersion}`,
    )
  }

  const tabsList = [
    {
      label: (
        <div style={tabStyle}>
          <img src={clusterPng} style={iconStyle} />
          {t('AllClusters')}
          <span style={numberStyle}>{summary?.totalCount}</span>
        </div>
      ),
      value: 'all',
    },
    {
      label: (
        <div style={tabStyle}>
          <img src={healthPng} style={iconStyle} />
          {t('HealthyClusters')}
          <span style={numberStyle}>{summary?.healthyCount}</span>
        </div>
      ),
      value: 'healthy',
    },
    {
      label: (
        <div style={tabStyle}>
          <img src={exceptionalPng} style={iconStyle} />
          {t('UnhealthyClusters')}
          <span style={numberStyle}>{summary?.unhealthyCount}</span>
        </div>
      ),
      value: 'exception',
    },
  ]

  const orderIconStyle = {
    marginLeft: 0,
  }

  function handleSort(key) {
    setSortParams({
      orderBy: key,
      isAsc: !sortParams?.isAsc,
    })
    getPageData({
      orderBy: key,
      isAsc: !sortParams?.isAsc,
    })
  }

  return (
    <div className={styles.container}>
      <div className={styles.action_bar}>
        <h4 className={styles.title}>{t('ClusterManagement')}</h4>
        {pageData && pageData?.length > 0 && (
          <Button type="primary" onClick={join} disabled={isReadOnlyMode}>
            {t('RegisterCluster')}
          </Button>
        )}
      </div>
      {loading ? (
        <div className={styles.loading_container}>
          <Loading />
        </div>
      ) : !pageData || !pageData?.length ? (
        <div className={styles.empty_content}>
          <div
            className={styles.empty_data}
            style={{ transform: `scale(${scale})` }}
          >
            <div className={styles.left}>
              <div className={styles.nodate}>{t('EmptyCluster')}</div>
              <div className={styles.tip}>
                {t('ClusterRequiresKubeConfigConfigurationFileAccess')}
              </div>
              <Button type="primary" onClick={join} disabled={isReadOnlyMode}>
                {t('RegisterCluster')}
              </Button>
            </div>
            <div className={styles.right}>
              <Empty />
            </div>
          </div>
        </div>
      ) : (
        <div className={styles.content}>
          <div className={styles.stat}>
            <KarporTabs
              list={tabsList}
              current={currentTab}
              onChange={handleTabChange}
              boxStyle={{ width: '100%' }}
              itemStyle={{ width: '33%' }}
            />
          </div>
          <div
            className={`${styles.page_content} ${styles[`page_content_${triangleLeftOffestIndex}`]}`}
          >
            <div className={styles.tool_bar}>
              <Input
                value={searchValue}
                onChange={event => setSearchValue(event.target.value)}
                style={{ width: 160, marginRight: 16 }}
                placeholder={t('PleaseEnterKeywords')}
                allowClear
                suffix={<SearchOutlined />}
              />
              <Button
                type="link"
                style={{ color: '#646566' }}
                onClick={() => handleSort('name')}
              >
                {t('SortByName')}
                {sortParams?.orderBy === 'name' &&
                  (sortParams?.isAsc ? (
                    <SortDescendingOutlined style={orderIconStyle} />
                  ) : (
                    <SortAscendingOutlined style={orderIconStyle} />
                  ))}
              </Button>
              <Button
                type="link"
                style={{ color: '#646566' }}
                onClick={() => handleSort('timestamp')}
              >
                {t('SortByTime')}
                {sortParams?.orderBy === 'timestamp' &&
                  (sortParams?.isAsc ? (
                    <ArrowUpOutlined style={orderIconStyle} />
                  ) : (
                    <ArrowDownOutlined style={orderIconStyle} />
                  ))}
              </Button>
            </div>
            {showPageData && showPageData?.length > 0 ? (
              <div className={styles.page_list}>
                {showPageData?.map((item: any, index: number) => (
                  <ClusterCard
                    key={`${item?.name}_${index}`}
                    item={item}
                    deleteItem={deleteItem}
                    goDetailPage={goDetailPage}
                    goCertificate={goCertificate}
                    setLastDetail={setLastDetail}
                    handleSubmit={handleSubmit}
                    customStyle={
                      showPageData?.length - 1 === index
                        ? {}
                        : {
                            borderBottom: '1px solid rgb(0 10 26 / 5%)',
                          }
                    }
                  />
                ))}
              </div>
            ) : (
              <div className={styles.empty_data}>
                <Empty />
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

export default Cluster
