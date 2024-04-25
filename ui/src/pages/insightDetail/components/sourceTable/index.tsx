import { SearchOutlined } from '@ant-design/icons'
import { Button, Input, Space, Table, message } from 'antd'
import axios from 'axios'
import queryString from 'query-string'
import React, { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useLocation, useNavigate } from 'react-router-dom'

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
  const [loading, setLoading] = useState(false)

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
      render: (_, record) => {
        return (
          <Button type="link" onClick={() => goResourcePage(record)}>
            {record?.object?.metadata?.name}
          </Button>
        )
      },
    },
    {
      dataIndex: 'namespace',
      key: 'namespace',
      title: 'Namespace',
      render: (_, record) => {
        return record?.object?.metadata?.namespace
      },
    },
    {
      dataIndex: 'apiVersion',
      key: 'apiVersion',
      title: 'APIVersion',
      render: (_, record) => {
        return record?.object?.apiVersion
      },
    },
    {
      dataIndex: 'kind',
      key: 'kind',
      title: 'Kind',
      render: (_, record) => {
        return record?.object?.kind
      },
    },
  ]

  async function queryTableData(params) {
    const { current, pageSize } = pageParams
    setLoading(true)
    const response: any = await axios.get(
      `/rest-api/v1/search?query=${queryStr}&pattern=sql&page=${params?.current || current}&pageSize=${params?.pageSize || pageSize}`,
    )
    if (response?.success) {
      setTableData(response?.data?.items || [])
      setPageParams({
        ...params,
        total: response?.data?.total,
      })
    } else {
      message.error(response?.message)
    }
    setLoading(false)
  }

  useEffect(() => {
    if (queryStr) {
      queryTableData({ current: 1, pageSize: pageParams?.pageSize })
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [queryStr])

  function handleTableChange({ current, pageSize }) {
    queryTableData({ current, pageSize })
  }

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
            disabled
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
