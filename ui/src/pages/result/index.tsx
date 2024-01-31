import { useState, useEffect, useRef } from 'react'
import {
  Pagination,
  Empty,
  Space,
  Button,
  Divider,
  message,
  AutoComplete,
  Tooltip,
  Input,
} from 'antd'
import axios from 'axios'
import { useLocation, useNavigate } from 'react-router-dom'
import { ClockCircleOutlined, CloseOutlined } from '@ant-design/icons'
import queryString from 'query-string'
import KarbourTabs from '../../components/tabs/index'
import { searchPrefix } from '../../utils/constants'

import { utcDateToLocalDate } from '../../utils/tools'
import Loading from '../../components/loading'
import cRole from '../../assets/labeled/c-role-128.png'
import cm from '../../assets/labeled/cm-128.png'
import crb from '../../assets/labeled/crb-128.png'
import cronjob from '../../assets/labeled/cronjob-128.png'
import deploy from '../../assets/labeled/deploy-128.png'
import ds from '../../assets/labeled/ds-128.png'
import ep from '../../assets/labeled/ep-128.png'
import group from '../../assets/labeled/group-128.png'
import hpa from '../../assets/labeled/hpa-128.png'
import ing from '../../assets/labeled/ing-128.png'
import job from '../../assets/labeled/job-128.png'
import limits from '../../assets/labeled/limits-128.png'
import netpol from '../../assets/labeled/netpol-128.png'
import ns from '../../assets/labeled/ns-128.png'
import pod from '../../assets/labeled/pod-128.png'
import psp from '../../assets/labeled/psp-128.png'
import pv from '../../assets/labeled/pv-128.png'
import pvc from '../../assets/labeled/pvc-128.png'
import quota from '../../assets/labeled/quota-128.png'
import rb from '../../assets/labeled/rb-128.png'
import role from '../../assets/labeled/role-128.png'
import rs from '../../assets/labeled/rs-128.png'
import sa from '../../assets/labeled/sa-128.png'
import sc from '../../assets/labeled/sc-128.png'
import secret from '../../assets/labeled/secret-128.png'
import sts from '../../assets/labeled/sts-128.png'
import svc from '../../assets/labeled/svc-128.png'
import user from '../../assets/labeled/user-128.png'
import volume from '../../assets/labeled/vol-128.png'
import nodeIcon from '../../assets/labeled/node-128.png'
import crd from '../../assets/labeled/crd-128.png'

import styles from './styles.module.less'
import React from 'react'

// crd 用来做所有用户自定义的非原生的资源图标
const ICON_MAP = {
  ClusterRole: cRole,
  ConfigMap: cm,
  ClusterRoleBinding: crb,
  CronJob: cronjob,
  Deployment: deploy,
  CafeDeployment: crd,
  DaemonSet: ds,
  Endpoint: ep,
  Group: group,
  HorizontalPodAutoscaler: hpa,
  Ingress: ing,
  Job: job,
  Limits: limits,
  NetworkPolicy: netpol,
  Namespace: ns,
  Pod: pod,
  PodSecurityPolicy: psp,
  PersistentVolume: pv,
  PersistentVolumeClaim: pvc,
  ResourceQuota: quota,
  RoleBinding: rb,
  Role: role,
  ReplicaSet: rs,
  ServiceAccount: sa,
  StorageClass: sc,
  Secret: secret,
  StatefulSet: sts,
  Service: svc,
  User: user,
  Volume: volume,
  Node: nodeIcon,
  InPlaceSet: crd,
  PodDisruptionBudget: crd,
}

const tabsList = [
  { label: '按照关键字搜索', value: 'keyword', disabled: true },
  { label: '按照 SQL 搜索', value: 'sql' },
]

const Result = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const [pageData, setPageData] = useState<any>()
  const urlSearchParams = queryString.parse(location.search)
  const [searchType, setSearchType] = useState<string>('sql')
  const [searchParams, setSearchParams] = useState({
    pageSize: 20,
    page: 1,
    query: urlSearchParams?.query || '',
    total: 0,
  })
  const [options, setOptions] = useState<{ value: string }[]>([])
  const [loading, setLoading] = useState(false)
  const [optionsCopy, setOptionsCopy] = useState<{ value: string }[]>([])
  const optionsRef = useRef<any>(getHistoryList())

  function getHistoryList() {
    const historyList: any = localStorage?.getItem(`${searchType}History`)
      ? JSON.parse(localStorage?.getItem(`${searchType}History`))
      : []
    return historyList
  }

  function deleteHistoryByItem(searchType: string, val: string) {
    const lastHistory: any = localStorage.getItem(`${searchType}History`)
    const tmp = lastHistory ? JSON.parse(lastHistory) : []
    if (tmp?.length > 0 && tmp?.includes(val)) {
      const newList = tmp?.filter(item => item !== val)
      localStorage.setItem(`${searchType}History`, JSON.stringify(newList))
    }
  }

  function deleteItem(event, value) {
    event.preventDefault()
    event.stopPropagation()
    deleteHistoryByItem(searchType, value)
    optionsRef.current = getHistoryList()
    setOptionsCopy(optionsRef.current)
  }

  useEffect(() => {
    const tmpOption = optionsRef.current?.map(item => ({
      label: (
        <div className={styles.option_item}>
          <div className={styles.option_item_label}>{item}</div>
          <div
            className={styles.option_item_delete}
            onClick={event => deleteItem(event, item)}
          >
            <CloseOutlined style={{ color: '#808080' }} />
          </div>
        </div>
      ),
      value: item,
    }))
    setOptions(tmpOption)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [optionsCopy])

  function handleTabChange(value: string) {
    setSearchType(value)
  }

  function handleChangePage(page: number, pageSize: number) {
    getPageData({
      ...searchParams,
      page,
      pageSize,
    })
  }

  async function getPageData(params) {
    setLoading(true)
    const response: any = await axios('/rest-api/v1/search', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      params: {
        query: `${searchPrefix} ${params?.query || searchParams?.query}`, // encodeURIComponent(searchValue as any),
        ...(searchType === 'sql' ? { pattern: 'sql' } : {}),
        page: params?.page || searchParams?.page,
        pageSize: params?.pageSize || searchParams?.pageSize,
      },
    })
    if (response?.success) {
      setPageData(response?.data?.items || {})
      setSearchParams({
        ...searchParams,
        ...params,
        total: response?.data?.total,
      })
      const objParams = {
        ...urlSearchParams,
        query: params?.query || searchParams?.query,
      }
      const urlString = queryString.stringify(objParams)
      navigate(`${location?.pathname}?${urlString}`, { replace: true })
    } else {
      message.error(response?.message || '请求失败，请重试')
    }
    setLoading(false)
  }

  useEffect(() => {
    getPageData(searchParams)
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  function handleInputChange(value: any) {
    setSearchParams({
      ...searchParams,
      query: value,
    })
  }

  function cacheHistory(searchType: string, val: string) {
    const lastHistory: any = localStorage.getItem(`${searchType}History`)
    const tmp = lastHistory ? JSON.parse(lastHistory) : []
    if (tmp?.length > 0 && tmp?.includes(val)) {
      return
    } else {
      const newList = [val, ...tmp]
      localStorage.setItem(`${searchType}History`, JSON.stringify(newList))
      optionsRef.current = getHistoryList()
      setOptionsCopy(optionsRef.current)
    }
  }

  function handleSearch() {
    // if (!searchParams?.query) {
    //   message.warning("请输入查询条件")
    //   return
    // }
    if (searchParams?.query) {
      cacheHistory(searchType, searchParams?.query as any)
    }
    getPageData({
      ...searchParams,
      page: 1,
    })
  }

  const handleClick = (item: any, key: string) => {
    const nav = key === 'name' ? 'resource' : key
    const objParams = {
      from: 'result',
      cluster: item?.cluster,
      apiVersion: item?.object?.apiVersion,
      type: key,
      kind: item?.object?.kind,
      namespace: item?.object?.metadata?.namespace,
      name: item?.object?.metadata?.name,
      query: searchParams?.query,
    }
    const urlParams = queryString.stringify(objParams)
    navigate(`/insightDetail/${nav}?${urlParams}`)
  }

  const handleTitleClick = (item: any, kind: string) => {
    const nav = kind === 'Namespace' ? 'namespace' : 'resource'
    const objParams = {
      from: 'result',
      cluster: item?.cluster,
      apiVersion: item?.object?.apiVersion,
      type: nav,
      kind: item?.object?.kind,
      ...(nav === 'namespace'
        ? { namespace: item?.object?.metadata?.name }
        : { namespace: item?.object?.metadata?.namespace }),
      ...(nav === 'resource' ? { name: item?.object?.metadata?.name } : {}),
      query: searchParams?.query,
    }
    const urlParams = queryString.stringify(objParams)
    navigate(`/insightDetail/${nav}?${urlParams}`)
  }

  function handleOnkeyUp(event) {
    if (event?.code === 'Enter' && event?.keyCode === 13) {
      handleSearch()
    }
  }

  return (
    <div className={styles.container}>
      <div className={styles.searchTab}>
        <KarbourTabs
          list={tabsList}
          current={searchType}
          onChange={handleTabChange}
        />
      </div>
      <div style={{ width: 850, position: 'relative' }}>
        <Space.Compact>
          <Input disabled value={searchPrefix} style={{ width: 180 }} />
          <AutoComplete
            onKeyUp={handleOnkeyUp}
            options={options}
            // onSearch={(text) => }
            placeholder={
              searchType === 'keyword'
                ? '支持搜索集群，集群资源（service/pod/cafed）...'
                : '支持 SQL 语句查询'
            }
            filterOption={(inputValue, option) =>
              option!.value.toUpperCase().indexOf(inputValue.toUpperCase()) !==
              -1
            }
            style={{ width: 600 }}
            value={searchParams?.query}
            allowClear={true}
            onChange={handleInputChange}
          />
          <Button type="primary" onClick={handleSearch}>
            搜索
          </Button>
        </Space.Compact>
      </div>
      <div className={styles.content}>
        {loading ? (
          <div
            style={{ height: 500, display: 'flex', justifyContent: 'center' }}
          >
            <Loading />
          </div>
        ) : pageData && pageData?.length > 0 ? (
          <>
            {/* 汇总 */}
            <div className={styles.stat}>
              <div>约&nbsp;{searchParams?.total}&nbsp;条搜索结果</div>
            </div>
            {pageData?.map((item: any, index: number) => {
              return (
                <div className={styles.card} key={`${item?.name}_${index}`}>
                  <div className={styles.left}>
                    <img
                      src={ICON_MAP?.[item?.object?.kind] || crd}
                      alt="icon"
                    />
                  </div>
                  <div className={styles.right}>
                    <div
                      className={styles.top}
                      onClick={() => handleTitleClick(item, item?.object?.kind)}
                    >
                      {item?.object?.metadata?.name || '--'}
                    </div>
                    <div className={styles.bottom}>
                      <div
                        className={styles.item}
                        onClick={() => handleClick(item, 'cluster')}
                      >
                        <span className={styles.api_icon}>Cluster</span>
                        <span className={styles.label}>
                          {item?.cluster || '--'}
                        </span>
                      </div>
                      <Divider type="vertical" />
                      <div className={`${styles.item} ${styles.disable}`}>
                        <span className={styles.api_icon}>APIVersion</span>
                        <span className={styles.label}>
                          {item?.object?.apiVersion || '--'}
                        </span>
                      </div>
                      <Divider type="vertical" />
                      <div
                        className={styles.item}
                        onClick={() => handleClick(item, 'kind')}
                      >
                        <span className={styles.api_icon}>Kind</span>
                        <span className={styles.label}>
                          {item?.object?.kind || '--'}
                        </span>
                      </div>
                      <Divider type="vertical" />
                      {item?.object?.metadata?.namespace && (
                        <>
                          <div
                            className={styles.item}
                            onClick={() => handleClick(item, 'namespace')}
                          >
                            <span className={styles.api_icon}>Namespace</span>
                            <span className={styles.label}>
                              {item?.object?.metadata?.namespace || '--'}
                            </span>
                          </div>
                          <Divider type="vertical" />
                        </>
                      )}
                      <div className={`${styles.item} ${styles.disable}`}>
                        <ClockCircleOutlined />
                        <Tooltip title="创建时间">
                          <span className={styles.label}>
                            {utcDateToLocalDate(
                              item?.object?.metadata?.creationTimestamp,
                            ) || '--'}
                          </span>
                        </Tooltip>
                      </div>
                    </div>
                  </div>
                </div>
              )
            })}
            <div className={styles.footer}>
              <Pagination
                total={searchParams?.total}
                showTotal={(total: number, range: any[]) =>
                  `${range[0]}-${range[1]} 共 ${total} 条`
                }
                pageSize={searchParams?.pageSize}
                current={searchParams?.page}
                onChange={handleChangePage}
              />
            </div>
          </>
        ) : (
          <div
            style={{
              height: 500,
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center',
            }}
          >
            <Empty />
          </div>
        )}
      </div>
    </div>
  )
}

export default Result
