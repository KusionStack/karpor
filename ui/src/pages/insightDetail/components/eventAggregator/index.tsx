import React, { useEffect, useRef, useState } from 'react'
import { Table, Select, Spin, Empty, Alert, Input, Skeleton } from 'antd'
import { SearchOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { formatTime } from '@/utils/tools'
import axios from 'axios'
import classNames from 'classnames'
import styles from './styles.module.less'

interface Event {
  type: string
  reason: string
  message: string
  count: number
  lastTimestamp: string
  firstTimestamp: string
}

interface EventAggregatorProps {
  cluster: string
  namespace: string
  name: string
  kind: string
  apiVersion: string
}

const EventAggregator: React.FC<EventAggregatorProps> = ({
  cluster,
  namespace,
  name,
  kind,
  apiVersion,
}) => {
  const { t } = useTranslation()
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string>()
  const [eventType, setEventType] = useState<string>()
  const [hasEvents, setHasEvents] = useState(false)
  const [searchText, setSearchText] = useState('')
  const eventSource = useRef<EventSource>()

  useEffect(() => {
    if (!cluster || !namespace || !name || !kind || !apiVersion) {
      return
    }

    const connect = () => {
      setLoading(true)
      setError(undefined)

      const params = new URLSearchParams({
        kind,
        apiVersion,
        ...(eventType && { type: eventType }),
      })

      const baseUrl = axios.defaults.baseURL || ''
      const url = `${baseUrl}/rest-api/v1/insight/aggregator/event/${cluster}/${namespace}/${name}?${params}`
      eventSource.current = new EventSource(url)

      eventSource.current.onopen = () => {
        setLoading(false)
      }

      eventSource.current.onmessage = event => {
        try {
          const events: Event[] = JSON.parse(event.data)

          if (
            events.length === 1 &&
            events[0].type === 'Warning' &&
            events[0].reason === 'Error'
          ) {
            setError(events[0].message)
            return
          }

          setEvents(events)
          setHasEvents(events.length > 0 || hasEvents)
        } catch (error) {
          setError(t('EventAggregator.Error'))
        }
      }

      eventSource.current.onerror = () => {
        setLoading(false)
        setError(t('EventAggregator.Error'))
        eventSource.current?.close()
        setTimeout(connect, 5000)
      }
    }

    connect()

    return () => {
      if (eventSource.current) {
        eventSource.current.close()
        eventSource.current = undefined
      }
    }
  }, [cluster, namespace, name, kind, apiVersion, eventType, t])

  const filteredEvents = events.filter(event => {
    const searchLower = searchText.toLowerCase()
    return (
      !searchText ||
      event.type.toLowerCase().includes(searchLower) ||
      event.reason.toLowerCase().includes(searchLower) ||
      event.message.toLowerCase().includes(searchLower)
    )
  })

  const handleSearch = (value: string) => {
    setSearchText(value)
  }

  const columns = [
    {
      title: t('EventAggregator.ColumnType'),
      dataIndex: 'type',
      key: 'type',
      width: 100,
      render: (type: string) => (
        <span
          className={classNames(styles.tag, styles.typeTag, {
            [styles.normal]: type === 'Normal',
            [styles.warning]: type === 'Warning',
          })}
        >
          {type}
        </span>
      ),
    },
    {
      title: t('EventAggregator.ColumnReason'),
      dataIndex: 'reason',
      key: 'reason',
      width: 150,
      render: (reason: string) => <span className={styles.tag}>{reason}</span>,
    },
    {
      title: t('EventAggregator.ColumnMessage'),
      dataIndex: 'message',
      key: 'message',
    },
    {
      title: t('EventAggregator.ColumnTimes'),
      dataIndex: 'count',
      key: 'count',
      width: 100,
      render: (count: number) => (
        <span className={classNames(styles.tag, styles.countTag)}>{count}</span>
      ),
    },
    {
      title: t('EventAggregator.ColumnLastSeen'),
      dataIndex: 'lastTimestamp',
      key: 'lastTimestamp',
      width: 200,
      render: (time: string) => formatTime(time),
    },
  ]

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <div className={styles.toolBar}>
          <Input
            className={styles.searchInput}
            placeholder={t('Search events...')}
            prefix={<SearchOutlined />}
            onChange={e => handleSearch(e.target.value)}
            allowClear
          />
          {hasEvents && (
            <Select
              value={eventType}
              onChange={setEventType}
              allowClear
              placeholder={t('EventAggregator.Type')}
              className={styles.typeFilter}
            >
              <Select.Option value="Normal">
                {t('EventAggregator.Normal')}
              </Select.Option>
              <Select.Option value="Warning">
                {t('EventAggregator.Warning')}
              </Select.Option>
            </Select>
          )}
        </div>
      </div>

      {error && (
        <Alert message={error} type="error" showIcon className={styles.error} />
      )}

      <Spin spinning={loading}>
        {loading ? (
          <Skeleton active paragraph={{ rows: 5 }} />
        ) : events.length > 0 ? (
          <Table
            dataSource={filteredEvents}
            columns={columns}
            rowKey={record =>
              `${record.type}-${record.reason}-${record.message}-${record.count}-${record.lastTimestamp}`
            }
            pagination={false}
            size="small"
          />
        ) : (
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description={t('EventAggregator.NoEvents')}
          />
        )}
      </Spin>
    </div>
  )
}

export default EventAggregator
