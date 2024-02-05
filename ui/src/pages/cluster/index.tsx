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
import KarbourTabs from '../../components/tabs'
import styles from './styles.module.less'
import Loading from '../../components/loading'
import ClusterCard from './components/clusterCard'
import healthPng from '@/assets/health_green.png'
import execptionalPng from '@/assets/execptional.png'
// import clusterPng from '@/assets/cluster.png'
import clusterPng from '@/assets/cluster_outlind.png'

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
  const [searchValue, setSearchValue] = useState('')

  const [lastDetail, setLastDetail] = useState<any>()

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

  // function handleChangePage(page: number, pageSize: number) {
  //   setSearchParams({
  //     ...searchParams,
  //     page,
  //     pageSize,
  //   });
  // }

  const join = () => {
    if (isReadOnlyMode) {
      return
    }
    navigate('/cluster/access')
  }

  async function handleSubmit(values, callback: () => void) {
    if (isReadOnlyMode) {
      return
    }
    const response: any = await axios({
      url: `/rest-api/v1/cluster/${lastDetail?.metadata?.name}`,
      method: 'PUT',
      data: values,
    })
    if (response?.success) {
      message.success('修改成功')
      callback()
      getClusterSummary()
      getPageData(sortParams)
    } else {
      message.error(response?.message || '请求失败，请重试')
    }
  }

  const [currentTab, setCurrentTab] = useState('all')
  // const [radioValue, setRadioValue] = useState('all');
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

  const iconStyle: any = {
    width: 24,
    height: 24,
    display: 'flex',
    justifyContent: 'center',
    lineHeight: '24px',
    alignItems: 'center',
    marginRight: 10,
    borderRadius: 6,
  }

  const tabStyle = {
    display: 'flex',
    alignItems: 'center',
    fontSize: 14,
    fontWeight: 400,
  }

  const numberStyle = { paddingLeft: 10, fontSize: 24 }

  // const radioOptions = [
  //   { label: '全部', value: 'all' },
  //   { label: '我是Owner', value: 'owner' },
  // ]

  // function handleChangeRadio(event: any) {

  //   setRadioValue(event.target.value);
  // }

  async function deleteItem(item) {
    if (isReadOnlyMode) {
      return
    }
    const response: any = await axios({
      url: `/rest-api/v1/cluster/${item?.metadata?.name}`,
      method: 'DELETE',
    })
    if (response?.success) {
      message.success('删除成功')
      getPageData(sortParams)
      getClusterSummary()
    } else {
      message.error(response?.message || '请求失败，请重试')
    }
  }

  function goCertificate(item) {
    if (isReadOnlyMode) {
      return
    }
    navigate(
      `/cluster/certificate?cluster=${item?.metadata?.name}&apiVersion=${item?.apiVersion}`,
    )
    navigate(`/cluster/certificate?cluster=${item?.metadata?.name}`)
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
          全部集群<span style={numberStyle}>{summary?.totalCount}</span>
        </div>
      ),
      value: 'all',
    },
    {
      label: (
        <div style={tabStyle}>
          <img src={execptionalPng} style={iconStyle} />
          异常集群<span style={numberStyle}>{summary?.unhealthyCount}</span>
        </div>
      ),
      value: 'exception',
    },
    {
      label: (
        <div style={tabStyle}>
          <img src={healthPng} style={iconStyle} />
          健康集群<span style={numberStyle}>{summary?.healthyCount}</span>
        </div>
      ),
      value: 'healthy',
    },
    // { label: <div style={tabStyle}><DeleteRowOutlined style={{ ...iconStyle, background: 'rgba(0,10,26, 0.08)', color: '#000A1A' }} />已删除<span style={numberStyle}>3</span></div>, value: "delete" },
  ]

  const orderIconStyle = {
    marginLeft: 0,
    // color: '#2f54eb'
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
      <div className={styles.actionBar}>
        <div className={styles.title}>集群列表</div>
        {pageData && pageData?.length > 0 && (
          <Button type="primary" onClick={join} disabled={isReadOnlyMode}>
            接入集群
          </Button>
        )}
      </div>
      {!pageData || !pageData?.length ? (
        <div className={styles.emptyContent}>
          <div className={styles.emptyData}>
            <div className={styles.left}>
              <div className={styles.nodate}>当前暂无可管理的集群</div>
              <div className={styles.tip}>集群需 KubeConfig 配置文件接入</div>
              <div className={styles.btnBox}>
                <Button type="primary" onClick={join}>
                  接入集群
                </Button>
              </div>
            </div>
            <div className={styles.right}></div>
          </div>
        </div>
      ) : (
        <div className={styles.content}>
          <div className={styles.stat}>
            <KarbourTabs
              list={tabsList}
              current={currentTab}
              onChange={handleTabChange}
              boxStyle={{ width: '100%' }}
              itemStyle={{ width: '33%' }}
            />
          </div>
          <div
            className={`${styles.pageContent} ${styles[`pageContent_${triangleLeftOffestIndex}`]}`}
          >
            <div className={styles.toolBar}>
              <Input
                value={searchValue}
                onChange={event => {
                  setSearchValue(event.target.value)
                }}
                style={{ width: 160, marginRight: 16 }}
                placeholder="请输入搜索关键字"
                allowClear
                suffix={<SearchOutlined />}
              />
              <Button
                type="link"
                style={{ color: '#646566' }}
                onClick={() => handleSort('name')}
              >
                名称排序
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
                时间排序
                {sortParams?.orderBy === 'timestamp' &&
                  (sortParams?.isAsc ? (
                    <ArrowUpOutlined style={orderIconStyle} />
                  ) : (
                    <ArrowDownOutlined style={orderIconStyle} />
                  ))}
              </Button>
              {/* <Radio.Group options={radioOptions} onChange={handleChangeRadio} value={radioValue} optionType="button" buttonStyle="solid" /> */}
              {/* <div className={styles.right}>
                </div> */}
            </div>
            {loading ? (
              <div
                style={{
                  height: 300,
                  display: 'flex',
                  justifyContent: 'center',
                  alignItems: 'center',
                }}
              >
                <Loading />
              </div>
            ) : showPageData && showPageData?.length > 0 ? (
              <div className={styles.pageList}>
                {showPageData?.map((item: any, index: number) => {
                  return (
                    <ClusterCard
                      key={`${item?.name}_${index}`}
                      item={item}
                      deleteItem={deleteItem}
                      goDetailPage={goDetailPage}
                      goCertificate={goCertificate}
                      setLastDetail={setLastDetail}
                      handleSubmit={handleSubmit}
                    />
                  )
                })}
              </div>
            ) : (
              <div
                style={{
                  background: '#fff',
                  borderRadius: 8,
                  height: 500,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}
              >
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
