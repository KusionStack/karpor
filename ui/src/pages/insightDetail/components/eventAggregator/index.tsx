import React, { useEffect, useRef, useState, useCallback } from 'react'
import {
  Table,
  Select,
  Spin,
  Empty,
  Alert,
  Input,
  Skeleton,
  Button,
  Space,
  message,
  Tooltip,
  Modal,
} from 'antd'
import {
  SearchOutlined,
  CloseOutlined,
  PoweroffOutlined,
  ExpandOutlined,
} from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { formatTime } from '@/utils/tools'
import axios from 'axios'
import classNames from 'classnames'
import Markdown from 'react-markdown'
import { useSelector } from 'react-redux'
import { debounce } from 'lodash'
import aiSummarySvg from '@/assets/ai-summary.svg'

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
  const { t, i18n } = useTranslation()
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string>()
  const [eventType, setEventType] = useState<string>()
  const [hasEvents, setHasEvents] = useState(false)
  const [searchText, setSearchText] = useState('')
  const [diagnosis, setDiagnosis] = useState('')
  const [diagnosisStatus, setDiagnosisStatus] = useState<
    'idle' | 'loading' | 'streaming' | 'complete' | 'error'
  >('idle')
  const [isStreaming, setStreaming] = useState(false)
  const eventSource = useRef<EventSource>()
  const diagnosisEndRef = useRef<HTMLDivElement>(null)
  const abortControllerRef = useRef<AbortController | null>(null)
  const contentRef = useRef<HTMLDivElement>(null)
  const tableContainerRef = useRef<HTMLDivElement>(null)
  const [open, setOpen] = useState(false)

  const { aiOptions } = useSelector((state: any) => state.globalSlice)
  const isAIEnabled = aiOptions?.AIModel && aiOptions?.AIAuthToken

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
          const events: Event[] = JSON.parse(event?.data)

          if (
            events?.length === 1 &&
            events?.[0]?.type === 'Warning' &&
            events?.[0].reason === 'Error'
          ) {
            setError(events[0].message)
            return
          }

          setEvents(events)
          setHasEvents(events?.length > 0 || hasEvents)
        } catch (error) {
          setError(t('EventAggregator.ConnectionError'))
        }
      }

      eventSource.current.onerror = () => {
        setLoading(false)
        setError(t('EventAggregator.ConnectionError'))
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

  const debouncedDiagnose = useCallback(
    debounce(async () => {
      try {
        if (!events.length) {
          message.warning(t('EventAggregator.NoEvents'))
          return
        }

        // Reset diagnosis state
        setDiagnosis('')
        setDiagnosisStatus('loading')
        setStreaming(true)

        // Cancel any existing SSE connection
        if (abortControllerRef.current) {
          abortControllerRef.current.abort()
        }

        // Create new AbortController for this request
        const abortController = new AbortController()
        abortControllerRef.current = abortController

        // Create new fetch request for diagnosis
        const url = `${axios.defaults.baseURL}/rest-api/v1/insight/aggregator/event/diagnosis/stream`

        // Send POST request and handle SSE response
        const response = await fetch(url, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Accept: 'text/event-stream',
          },
          body: JSON.stringify({
            events: events.map(event => ({
              type: event.type,
              reason: event.reason,
              message: event.message,
              count: event.count,
              lastTimestamp: event.lastTimestamp,
              firstTimestamp: event.firstTimestamp,
            })),
            language: i18n.language,
          }),
          signal: abortController.signal,
        })

        if (!response.ok) {
          throw new Error(response.statusText)
        }

        // Create a reader from the response body stream
        const reader = response.body?.getReader()
        const decoder = new TextDecoder()

        if (!reader) {
          throw new Error('No response body')
        }

        // Read the stream
        const processStream = async () => {
          try {
            let buffer = ''
            let streaming = true
            while (streaming) {
              const { done, value } = await reader.read()

              if (done) {
                streaming = false
                setDiagnosisStatus('complete')
                setStreaming(false)
                break
              }

              // Decode the chunk and add to buffer
              buffer += decoder.decode(value, { stream: true })

              // Process complete events in buffer
              const lines = buffer.split('\n\n')
              buffer = lines.pop() || '' // Keep the last incomplete event in buffer

              for (const line of lines) {
                if (!line.trim()) continue

                try {
                  const eventData = line.replace('data: ', '')
                  const diagEvent = JSON.parse(eventData)

                  switch (diagEvent.type) {
                    case 'start':
                      setDiagnosisStatus('streaming')
                      break
                    case 'chunk':
                      setDiagnosis(prev => prev + diagEvent.content)
                      // Scroll to bottom of diagnosis
                      if (diagnosisEndRef.current) {
                        diagnosisEndRef.current.scrollIntoView({
                          behavior: 'smooth',
                        })
                      }
                      break
                    case 'error':
                      streaming = false
                      setDiagnosisStatus('error')
                      setStreaming(false)
                      message.error(diagEvent.content)
                      reader.cancel()
                      return
                    case 'complete':
                      streaming = false
                      setDiagnosisStatus('complete')
                      setStreaming(false)
                      reader.cancel()
                      return
                  }
                } catch (error) {
                  console.error('Failed to parse diagnosis event:', error)
                }
              }
            }
          } catch (error) {
            if (error.name === 'AbortError') {
              console.log('Diagnosis stream aborted')
            } else {
              console.error('Error reading stream:', error)
              setDiagnosisStatus('error')
              setStreaming(false)
              message.error(t('EventAggregator.DiagnosisConnectionError'))
            }
          }
        }

        processStream()
      } catch (error) {
        console.error('Failed to start diagnosis:', error)
        setDiagnosisStatus('error')
        setStreaming(false)
        message.error(t('EventAggregator.FailedToDiagnoseLogs'))
      }
    }, 500),
    [events, t, i18n.language],
  )

  const startDiagnosis = useCallback(() => {
    debouncedDiagnose()
  }, [debouncedDiagnose])

  const isVertical = filteredEvents?.length <= 3
  const events_content_styles: React.CSSProperties = isVertical
    ? { flexDirection: 'column' }
    : { flexDirection: 'row' }
  const events_content_withDiagnosis_style: React.CSSProperties = isVertical
    ? { width: '100%' }
    : { width: 'calc(100% - 424px)' }
  const tableHeight = tableContainerRef?.current
    ? tableContainerRef?.current.getBoundingClientRect()?.height
    : 400
  const events_content_diagnosisPanel_style: React.CSSProperties = isVertical
    ? { width: '100%', height: 300 }
    : { width: 400, height: tableHeight || 400 }

  const contentToTopHeight = contentRef.current?.getBoundingClientRect()?.top
  const dotToTopHeight = diagnosisEndRef.current?.getBoundingClientRect()?.top

  const renderDiagnosisWindow = isDailog => {
    if (diagnosisStatus === 'idle') {
      return null
    }

    return (
      <div
        className={styles.events_content_diagnosisPanel}
        style={
          isDailog
            ? { width: '100%', height: 700 }
            : events_content_diagnosisPanel_style
        }
      >
        <div className={styles.events_content_diagnosisHeader}>
          <Space>
            <div className={styles.events_content_diagnosisHeader_aiIcon}>
              <img src={aiSummarySvg} alt="ai summary" />
            </div>
            {t('EventAggregator.DiagnosisResult')}
          </Space>
          <Space>
            {isStreaming && (
              <Tooltip
                title={t('EventAggregator.StopDiagnosis')}
                placement="bottom"
              >
                <Button
                  type="text"
                  className={styles.events_content_diagnosisHeader_stopButton}
                  icon={<PoweroffOutlined />}
                  onClick={() => {
                    if (abortControllerRef.current) {
                      abortControllerRef.current.abort()
                      setDiagnosisStatus('complete')
                      setStreaming(false)
                    }
                  }}
                />
              </Tooltip>
            )}
            {isDailog ? null : (
              <Button
                type="text"
                icon={<ExpandOutlined />}
                onClick={() => {
                  setOpen(!open)
                }}
              />
            )}
            <Button
              type="text"
              icon={<CloseOutlined />}
              onClick={() => {
                if (isDailog) {
                  setOpen(false)
                } else {
                  if (abortControllerRef.current) {
                    abortControllerRef.current.abort()
                  }
                  setDiagnosisStatus('idle')
                  setDiagnosis('')
                  setStreaming(false)
                }
              }}
            />
          </Space>
        </div>
        <div className={styles.events_content_diagnosisBody}>
          <div
            className={styles.events_content_diagnosisContent}
            ref={contentRef}
          >
            {diagnosisStatus === 'loading' ||
            (diagnosisStatus === 'streaming' && !diagnosis) ? (
              <div className={styles.events_content_diagnosisLoading}>
                <Spin />
                <p>{t('EventAggregator.DiagnosisInProgress')}</p>
              </div>
            ) : diagnosisStatus === 'error' ? (
              <Alert
                type="error"
                message={t('EventAggregator.DiagnosisFailed')}
                description={t('EventAggregator.TryAgainLater')}
              />
            ) : (
              <>
                <Markdown>{diagnosis}</Markdown>
                <div ref={diagnosisEndRef} />
              </>
            )}
          </div>
          {diagnosisStatus === 'streaming' && diagnosis && (
            <div
              className={`${styles.events_content_streamingIndicator} ${dotToTopHeight - contentToTopHeight + 53 - (isVertical ? 300 : 600) >= 0 ? styles.events_content_streamingIndicatorFixed : ''}`}
            >
              <span className={styles.dot}></span>
              <span className={styles.dot}></span>
              <span className={styles.dot}></span>
            </div>
          )}
        </div>
      </div>
    )
  }

  return (
    <div className={styles.events_container}>
      <div className={styles.events_header}>
        <div className={styles.events_toolBar}>
          {hasEvents && (
            <>
              <Input
                className={styles.events_toolBar_searchInput}
                placeholder={t('EventAggregator.SearchEvents')}
                prefix={<SearchOutlined />}
                allowClear
                value={searchText}
                onChange={e => handleSearch(e.target.value)}
              />
              <Select
                className={styles.events_toolBar_typeFilter}
                placeholder={t('EventAggregator.Type')}
                allowClear
                value={eventType}
                onChange={setEventType}
                options={[
                  { value: 'Normal', label: t('EventAggregator.Normal') },
                  { value: 'Warning', label: t('EventAggregator.Warning') },
                ]}
              />
              <Space>
                {isAIEnabled && (
                  <Tooltip title={t('EventAggregator.Diagnose')}>
                    <Button
                      type="text"
                      className={styles.events_toolBar_actionButton}
                      icon={
                        <span className={styles.events_toolBar_magicWand}>
                          ✨
                        </span>
                      }
                      onClick={startDiagnosis}
                      loading={diagnosisStatus === 'streaming'}
                      disabled={!hasEvents || diagnosisStatus === 'streaming'}
                    />
                  </Tooltip>
                )}
              </Space>
            </>
          )}
        </div>
      </div>

      {error && (
        <Alert
          message={error}
          type="error"
          showIcon
          className={styles.events_error}
        />
      )}

      <div className={styles.events_content} style={events_content_styles}>
        <div
          ref={tableContainerRef}
          className={classNames(styles.events_content_tableContainer, {
            [styles.events_content_withDiagnosis]: diagnosisStatus !== 'idle',
          })}
          style={events_content_withDiagnosis_style}
        >
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
        </div>
        {renderDiagnosisWindow(false)}
        <Modal
          styles={{
            content: {
              padding: 0,
              borderRadius: 16,
              background:
                'radial-gradient(ellipse 489px 674px at 6px 0px, #ffffff 0%, #f3efff 100%), radial-gradient(ellipse 587px 672px at 433px 513px, #ece6ff95 0%, #f0e9ff30 100%), radial-gradient(ellipse 346px 396px at 15px 506px, #e9e3ff90 0%, #f0e9ff30 100%), radial-gradient(ellipse 583px 668px at 436px 8px, #ffffff 0%, #f5f0ff30 100%)',
            },
          }}
          centered
          closable={false}
          width={'80%'}
          height={700}
          open={open}
          footer={null}
        >
          <div style={{ overflow: 'auto' }}>{renderDiagnosisWindow(true)}</div>
        </Modal>
      </div>
    </div>
  )
}

export default EventAggregator
