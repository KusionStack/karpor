import React, { useState, useEffect } from 'react'
import {
  Pagination,
  Empty,
  Divider,
  Tooltip,
  Input,
  message,
  AutoComplete,
  Space,
  Tag,
} from 'antd'
import { useLocation, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { ClockCircleOutlined, CloseOutlined } from '@ant-design/icons'
import queryString from 'query-string'
import classNames from 'classnames'
import SqlSearch from '@/components/sqlSearch'
import KarporTabs from '@/components/tabs/index'
import {
  cacheHistory,
  deleteHistoryByItem,
  getHistoryList,
  utcDateToLocalDate,
} from '@/utils/tools'
import Loading from '@/components/loading'
import { ICON_MAP } from '@/utils/images'
import { searchSqlPrefix, tabsList } from '@/utils/constants'
import { useAxios } from '@/utils/request'
import { PodStatusCell } from '@/pages/insightDetail/components/sourceTable'
// import useDebounce from '@/hooks/useDebounce'

import styles from './styles.module.less'

const { Search } = Input
const Option = AutoComplete.Option

export const CustomDropdown = props => {
  const { options } = props

  return (
    <div>
      {options.map((option, index) => (
        <div
          key={index}
          style={{ padding: '5px', borderBottom: '1px solid #ccc' }}
        >
          <Option value={option.value}>
            <span>{option.value}</span> -{' '}
            <span style={{ color: '#999' }}>{option.value || <></>}</span>
          </Option>
        </div>
      ))}
    </div>
  )
}

const Result = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const [pageData, setPageData] = useState<any>()
  const urlSearchParams: any = queryString.parse(location.search)
  const [searchType, setSearchType] = useState<string>(urlSearchParams?.pattern)
  const [searchParams, setSearchParams] = useState({
    pageSize: 20,
    page: 1,
    query: urlSearchParams?.query || '',
    total: 0,
  })
  const [naturalValue, setNaturalValue] = useState('')
  const [sqlValue, setSqlValue] = useState('')
  const [naturalOptions, setNaturalOptions] = useState(
    getHistoryList('naturalHistory') || [],
  )

  const { t } = useTranslation()

  function cacheNaturalHistory(key, val) {
    const result = cacheHistory(key, val)
    setNaturalOptions(result)
  }

  useEffect(() => {
    if (searchType === 'natural') {
      setNaturalValue(urlSearchParams?.query)
      handleNaturalSearch(urlSearchParams?.query)
    }
    if (searchType === 'sql') {
      setSqlValue(urlSearchParams?.query)
      handleSqlSearch(urlSearchParams?.query)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    if (urlSearchParams?.pattern) {
      setSearchType(urlSearchParams?.pattern)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [urlSearchParams?.pattern, urlSearchParams?.query])

  function handleTabChange(value: string) {
    setSearchType(value)
    const urlString = queryString.stringify({
      pattern: value,
      query:
        value === 'natural' ? naturalValue : value === 'sql' ? sqlValue : '',
    })
    navigate(`${location?.pathname}?${urlString}`, { replace: true })
  }

  function handleChangePage(page: number, pageSize: number) {
    getPageData({
      ...searchParams,
      page,
      pageSize,
    })
  }

  const { response, refetch, loading } = useAxios({
    url: '/rest-api/v1/search',
    manual: true,
  })

  useEffect(() => {
    if (response?.success) {
      const objParams = {
        ...urlSearchParams,
        pattern: 'sql',
        query: response?.successParams?.query || searchParams?.query,
      }
      if (searchType === 'natural') {
        let sqlVal
        if (response?.data?.sqlQuery?.includes('WHERE')) {
          sqlVal = `where ${response?.data?.sqlQuery?.split(' WHERE ')?.[1]}`
        }
        if (response?.data?.sqlQuery?.includes('where')) {
          sqlVal = `where ${response?.data?.sqlQuery?.split(' where ')?.[1]}`
        }
        setSearchType('sql')
        setSqlValue(sqlVal)
      }
      setPageData(response?.data?.items || {})
      const urlString = queryString.stringify(objParams)
      navigate(`${location?.pathname}?${urlString}`, { replace: true })
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [response])

  function getPageData(params) {
    const pattern =
      searchType === 'natural' ? 'nl' : searchType === 'sql' ? 'sql' : ''
    const query =
      searchType === 'natural'
        ? params?.query
        : searchType === 'sql'
          ? `${searchSqlPrefix} ${params?.query}`
          : ''
    refetch({
      option: {
        params: {
          pattern,
          query,
          page: params?.page || searchParams?.page,
          pageSize: params?.pageSize || searchParams?.pageSize,
        },
      },
    })

    setSearchParams({
      ...searchParams,
      ...params,
      total: response?.data?.total,
    })
  }

  function handleSqlSearch(inputValue) {
    setSqlValue(inputValue)
    setSearchParams({
      ...searchParams,
      query: inputValue,
    })
    getPageData({
      ...searchParams,
      query: inputValue,
      page: 1,
    })
  }

  const handleClick = (item: any, key: string) => {
    const nav = key === 'name' ? 'resource' : key
    const objParams = {
      from: 'result',
      deleted: item?.deleted,
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
      deleted: item?.deleted,
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

  function handleNaturalAutoCompleteChange(val) {
    setNaturalValue(val)
  }

  function handleNaturalSearch(value) {
    if (!value && !naturalValue) {
      message.warning(t('CannotBeEmpty'))
      return
    }
    cacheNaturalHistory('naturalHistory', value)
    getPageData({
      pageSize: searchParams?.pageSize,
      page: 1,
      query: value,
    })
  }

  function renderEmpty() {
    return (
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
    )
  }

  function renderLoading() {
    return (
      <div
        style={{
          height: 500,
          display: 'flex',
          justifyContent: 'center',
        }}
      >
        <Loading />
      </div>
    )
  }

  function renderListContent() {
    return (
      <>
        <div className={styles.stat}>
          <div>
            {t('AboutInSearchResult')}&nbsp;
            {searchParams?.total}&nbsp;
            {t('SearchResult')}
          </div>
        </div>
        {pageData?.map((item: any, index: number) => {
          return (
            <div className={styles.card} key={`${item?.name}_${index}`}>
              {item?.deleted && (
                <div className={styles.delete_tag}>
                  <Tag color="error">{t('Delete')}</Tag>
                </div>
              )}
              <div className={styles.left}>
                <img
                  src={ICON_MAP?.[item?.object?.kind] || ICON_MAP.CRD}
                  alt="icon"
                />
              </div>
              <div className={styles.right}>
                <div
                  className={styles.top}
                  onClick={() => handleTitleClick(item, item?.object?.kind)}
                >
                  <span style={{ marginRight: 8 }}>
                    {item?.object?.metadata?.name || '--'}
                  </span>
                  {item?.object?.kind === 'Pod' && (
                    <PodStatusCell
                      cluster={item?.cluster}
                      namespace={item?.object?.metadata?.namespace}
                      name={item?.object?.metadata?.name}
                    />
                  )}
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
                  <div className={classNames(styles.item, styles.disable)}>
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
                  <div className={classNames(styles.item, styles.disable)}>
                    <ClockCircleOutlined />
                    <Tooltip title={t('CreateTime')}>
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
              `${range?.[0]}-${range?.[1]} ${t('Total')} ${total} `
            }
            pageSize={searchParams?.pageSize}
            current={searchParams?.page}
            onChange={handleChangePage}
          />
        </div>
      </>
    )
  }

  const handleDelete = val => {
    deleteHistoryByItem('naturalHistory', val)
    const list = getHistoryList('naturalHistory') || []
    setNaturalOptions(list)
  }

  const renderOption = val => {
    return (
      <Space
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        <span>{val}</span>
        <CloseOutlined
          onClick={event => {
            event?.stopPropagation()
            handleDelete(val)
          }}
        />
      </Space>
    )
  }

  const tmpOptions = naturalOptions?.map(val => ({
    value: val,
    label: renderOption(val),
  }))

  return (
    <div className={styles.container}>
      <div className={styles.searchTab}>
        <KarporTabs
          list={tabsList}
          current={searchType}
          onChange={handleTabChange}
        />
      </div>
      {searchType === 'sql' && (
        <SqlSearch
          sqlEditorValue={(sqlValue || urlSearchParams?.query) as string}
          handleSqlSearch={handleSqlSearch}
        />
      )}
      {searchType === 'natural' && (
        <div className={styles.search_codemirror_container}>
          <AutoComplete
            style={{ width: '100%' }}
            size="large"
            options={tmpOptions}
            value={naturalValue}
            onChange={handleNaturalAutoCompleteChange}
            filterOption={(inputValue, option) => {
              if (option?.value) {
                return (
                  (option?.value as string)
                    ?.toUpperCase()
                    .indexOf(inputValue.toUpperCase()) !== -1
                )
              }
            }}
          >
            <Search
              size="large"
              placeholder={`${t('SearchByNaturalLanguage')}...`}
              enterButton
              onSearch={handleNaturalSearch}
            />
          </AutoComplete>
        </div>
      )}
      <div className={styles.content}>
        {loading
          ? renderLoading()
          : pageData && pageData?.length > 0
            ? renderListContent()
            : renderEmpty()}
      </div>
    </div>
  )
}

export default Result
