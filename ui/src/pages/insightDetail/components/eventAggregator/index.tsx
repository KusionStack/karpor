import React, { useEffect, useState, useRef } from 'react'
import { Table, Badge, Select, Empty, Alert, Spin } from 'antd'
import { useTranslation } from 'react-i18next'
import { formatTime } from '@/utils/tools'
import axios from 'axios'
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

  const columns = [
    {
      title: t('EventAggregator.ColumnType'),
      dataIndex: 'type',
      key: 'type',
      width: 100,
      render: (type: string) => (
        <Badge status={type === 'Normal' ? 'success' : 'error'} text={type} />
      ),
    },
    {
      title: t('EventAggregator.ColumnReason'),
      dataIndex: 'reason',
      key: 'reason',
      width: 150,
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
        {events.length > 0 && (
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

      {error && (
        <Alert message={error} type="error" showIcon className={styles.error} />
      )}

      <Spin spinning={loading}>
        {events.length > 0 ? (
          <Table
            dataSource={events}
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
