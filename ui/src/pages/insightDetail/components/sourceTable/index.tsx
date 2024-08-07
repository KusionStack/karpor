import { SearchOutlined } from '@ant-design/icons'
import { Button, Input, Space, Table } from 'antd'
import queryString from 'query-string'
import React, { useCallback, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useLocation, useNavigate } from 'react-router-dom'
import useDebounce from '@/hooks/useDebounce'
import { useAxios } from '@/utils/request'

import styles from './style.module.less'

type IProps = {
  queryStr: string
  data?: any[]
  searchKey?: string
  pagination?: any
  tableName?: string
}

const defaultSearchParams = {
  current: 1,
  pageSize: 10,
  total: 0,
}

const SourceTable = ({ queryStr, tableName }: IProps) => {
  const { t } = useTranslation()
  const location = useLocation()
  const navigate = useNavigate()
  const [pageParams, setPageParams] = useState(defaultSearchParams)
  const [tableData, setTableData] = useState([])
  const urlSearchParams = queryString?.parse(location?.search)
  const [keyword, setKeyword] = useState('')

  const debouncedValue = useDebounce(keyword, 500)

  function goResourcePage(record) {
    const nav = record?.object?.kind === 'Namespace' ? 'namespace' : 'resource'
    const params = {
      from: urlSearchParams?.from,
      type: nav,
      query: urlSearchParams?.query,
      cluster: record?.cluster,
      kind: record?.object?.kind,
      apiVersion: record?.object?.apiVersion,
      ...(nav === 'namespace'
        ? { namespace: record?.object?.metadata?.name }
        : { namespace: record?.object?.metadata?.namespace }),
      ...(nav === 'resource' ? { name: record?.object?.metadata?.name } : {}),
    }
    const urlParams = queryString?.stringify(params)
    navigate(`/insightDetail/${nav}?${urlParams}`)
  }

  const columns = [
    {
      dataIndex: 'name',
      key: 'name',
      title: t('Name'),
      render: (_, record) => (
        <Button type="link" onClick={() => goResourcePage(record)}>
          {record?.object?.metadata?.name}
        </Button>
      ),
    },
    {
      dataIndex: 'cluster',
      key: 'cluster',
      title: 'Cluster',
    },
    {
      dataIndex: 'namespace',
      key: 'namespace',
      title: 'Namespace',
      render: (_, record) => record?.object?.metadata?.namespace,
    },
    {
      dataIndex: 'apiVersion',
      key: 'apiVersion',
      title: 'APIVersion',
      render: (_, record) => record?.object?.apiVersion,
    },
    {
      dataIndex: 'kind',
      key: 'kind',
      title: 'Kind',
      render: (_, record) => record?.object?.kind,
    },
  ]

  const {
    response: tableDataResponse,
    refetch: tableDataRefetch,
    loading,
  } = useAxios({
    url: '',
    manual: true,
  })

  useEffect(() => {
    if (tableDataResponse?.success) {
      setTableData(tableDataResponse?.data?.items || [])
      setPageParams({
        ...tableDataResponse?.successParams,
        total: tableDataResponse?.data?.total,
      })
    }
  }, [tableDataResponse])

  function queryTableData(params) {
    const { current, pageSize } = pageParams
    setTableData([])
    setPageParams({
      ...params,
      total: 0,
    })
    tableDataRefetch({
      url: `/rest-api/v1/search?query=${queryStr}&pattern=sql&page=${params?.current || current}&pageSize=${params?.pageSize || pageSize}&keyword=${params?.keyword || ''}`,
      successParams: params,
    })
  }

  useEffect(() => {
    if (queryStr) {
      setKeyword('')
      queryTableData({ current: 1, pageSize: pageParams?.pageSize })
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [queryStr])

  function handleTableChange({ current, pageSize }) {
    queryTableData({ current, pageSize, keyword: debouncedValue })
  }

  const search = useCallback(() => {
    queryTableData({
      current: 1,
      pageSize: pageParams?.pageSize,
      keyword: debouncedValue,
    })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [debouncedValue])

  useEffect(() => {
    if (debouncedValue) {
      search()
    }
  }, [debouncedValue, search])

  return (
    <div>
      <div className={styles.table_header}>
        <div className={styles.table_title}>
          {tableName || '--'}
          {urlSearchParams?.type === 'kind' ? null : (
            <span className={styles.tips}>
              {t('SelectResourcesAboveToSeeDetailsHere')}
            </span>
          )}
        </div>
        <Space style={{ marginBottom: 10 }}>
          <Input
            value={keyword}
            allowClear
            disabled
            onChange={event => setKeyword(event?.target?.value)}
            placeholder={t('FilterByName')}
            suffix={<SearchOutlined />}
          />
        </Space>
      </div>
      <Table
        loading={loading}
        columns={columns}
        dataSource={tableData}
        rowKey={record => {
          return `${record?.object?.metadata?.name}_${record?.object?.metadata?.namespace}_${record?.object?.apiVersion}_${record?.object?.kind}`
        }}
        onChange={handleTableChange}
        pagination={{
          ...pageParams,
          showSizeChanger: true,
        }}
      />
    </div>
  )
}

export default SourceTable
