import React, { useEffect, useState } from 'react'
import { Empty, Button, Input, message } from 'antd'
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
import classNames from 'classnames'
import KarporTabs from '@/components/tabs'
import Loading from '@/components/loading'
import ClusterCard from './components/clusterCard'
import healthPng from '@/assets/health_green.png'
import exceptionalPng from '@/assets/exceptional.png'
import clusterPng from '@/assets/cluster_outlind.png'
import { useAxios } from '@/utils/request'

import styles from './styles.module.less'

const Cluster = () => {
  const navigate = useNavigate()
  const { t } = useTranslation()
  const { isReadOnlyMode } = useSelector((state: any) => state.globalSlice)
  const [pageData, setPageData] = useState<any>([])
  const [showPageData, setShowPageData] = useState<any>([])
  const [summary, setSummary] = useState<any>()
  const [sortParams, setSortParams] = useState<any>({
    orderBy: 'name',
    isAsc: true,
  })
  const [searchValue, setSearchValue] = useState('')
  const [lastDetail, setLastDetail] = useState<any>()
  const [scale, setScale] = useState<any>(1)

  const {
    response: summaryResponse,
    loading,
    refetch: summaryRefetch,
  } = useAxios({
    url: `/rest-api/v1/clusters?summary=true`,
    option: { params: {} },
    manual: true,
    method: 'GET',
  })

  const { response: pageResponse, refetch: pageRefetch } = useAxios({
    url: '/rest-api/v1/clusters',
    option: { params: {} },
    manual: true,
    method: 'GET',
  })

  const { response: deleteResponse, refetch: deleteRefetch } = useAxios({
    url: '/rest-api/v1/clusters',
    option: { params: {} },
    manual: true,
    method: 'DELETE',
  })

  const { response: putClusterResponse, refetch: putClusterRefetch } = useAxios(
    {
      url: '/rest-api/v1/clusters',
      option: { data: {} },
      manual: true,
      method: 'PUT',
    },
  )

  useEffect(() => {
    if (summaryResponse?.success) {
      setSummary(summaryResponse?.data)
    }
  }, [summaryResponse])

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

  function getClusterSummary() {
    summaryRefetch()
  }

  function getPageData(params) {
    pageRefetch({
      option: {
        params: {
          orderBy: params?.orderBy,
          ...(params?.isAsc ? { ascending: true } : { descending: true }),
        },
      },
    })
  }

  useEffect(() => {
    if (pageResponse?.success) {
      setPageData(pageResponse?.data?.items)
    }
  }, [pageResponse])

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

  useEffect(() => {
    getClusterSummary()
    getPageData(sortParams)
  }, [])

  const join = () => {
    if (isReadOnlyMode) return
    navigate('/cluster/access')
  }

  function handleSubmit(values, callback: () => void) {
    if (isReadOnlyMode) return
    putClusterRefetch({
      url: `/rest-api/v1/cluster/${lastDetail?.metadata?.name}`,
      option: {
        data: {
          values,
        },
      },
      callbackFn: callback,
    })
  }

  useEffect(() => {
    if (putClusterResponse?.success) {
      message.success(t('UpdateSuccess'))
      putClusterResponse?.callbackFn()
      getClusterSummary()
      getPageData(sortParams)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [putClusterResponse])

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

  function deleteItem(item) {
    if (isReadOnlyMode) return
    deleteRefetch({
      url: `/rest-api/v1/cluster/${item?.metadata?.name}`,
    })
  }

  useEffect(() => {
    if (deleteResponse?.success) {
      message.success(t('DeletedSuccess'))
      getPageData(sortParams)
      getClusterSummary()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [deleteResponse])

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
              <div className={styles.nodata}>{t('EmptyCluster')}</div>
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
            className={classNames(
              styles.page_content,
              styles[`page_content_${triangleLeftOffestIndex}`],
            )}
          >
            <div className={styles.tool_bar}>
              <Input
                value={searchValue}
                onChange={event => setSearchValue(event.target.value)}
                style={{ width: 300, marginRight: 16 }}
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
