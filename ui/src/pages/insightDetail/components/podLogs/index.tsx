import React, { useEffect, useRef, useState, useCallback } from 'react'
import {
  Select,
  Space,
  Button,
  Alert,
  Badge,
  Tooltip,
  message,
  Spin,
  DatePicker,
  Input,
  InputNumber,
  Modal,
  Switch,
} from 'antd'
import {
  PauseCircleOutlined,
  PlayCircleOutlined,
  ClearOutlined,
  RobotOutlined,
  CloseOutlined,
  PoweroffOutlined,
  DownloadOutlined,
  SettingOutlined,
  HighlightOutlined,
  FilterOutlined,
} from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import yaml from 'js-yaml'
import axios from 'axios'
import styles from './styles.module.less'
import Markdown from 'react-markdown'
import { useSelector } from 'react-redux'
import debounce from 'lodash.debounce'
import dayjs from 'dayjs'

interface LogEntry {
  timestamp: string
  content: string
  error?: string
  id?: string
}

interface PodLogsProps {
  cluster: string
  namespace: string
  podName: string
  yamlData?: string
}

type DiagnosisStatus = 'idle' | 'init' | 'streaming' | 'complete' | 'error'

interface LogSettings {
  since?: string
  sinceTime?: string
  tailLines?: number
  timestamps: boolean
  isShowConnect?: boolean
}

const PodLogs: React.FC<PodLogsProps> = ({
  cluster,
  namespace,
  podName,
  yamlData,
}) => {
  const { t, i18n } = useTranslation()
  const [container, setContainer] = useState<string>('')
  const [containers, setContainers] = useState<string[]>([])
  const [logs, setLogs] = useState<LogEntry[]>([])
  const [isPaused, setIsPaused] = useState(false)
  const [isConnected, setConnected] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [lastTimestamp, setLastTimestamp] = useState<string | null>(null)
  const [isReconnecting, setIsReconnecting] = useState(false)
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const lastErrorRef = useRef<number>(0)
  const [diagnosisStatus, setDiagnosisStatus] =
    useState<DiagnosisStatus>('idle')
  const [diagnosis, setDiagnosis] = useState('')
  const [isStreaming, setStreaming] = useState(false)
  const [settings, setSettings] = useState<LogSettings>({
    timestamps: true,
    tailLines: 100,
    isShowConnect: false,
  })
  const [showSettings, setShowSettings] = useState(false)
  const [searchText, setSearchText] = useState('')
  const [searchMode, setSearchMode] = useState<'highlight' | 'filter'>(
    'highlight',
  )
  const logsEndRef = useRef<HTMLDivElement>(null)
  const diagnosisEndRef = useRef<HTMLDivElement>(null)
  const eventSourceRef = useRef<EventSource | null>(null)
  const abortControllerRef = useRef<AbortController | null>(null)

  useEffect(() => {
    if (yamlData) {
      try {
        const podSpec = yaml.load(yamlData) as any

        let containerList: string[] = []
        if (podSpec?.spec?.containers) {
          containerList = podSpec.spec.containers.map((c: any) => c.name)
        }

        setContainers(containerList)
        if (containerList.length > 0 && !container) {
          setContainer(containerList[0])
        }
      } catch (error) {
        console.error('Failed to parse pod details:', error)
        setError(t('FailedToParsePodDetails'))
      }
    }
  }, [yamlData, container, t])

  useEffect(() => {
    if (!container || isPaused) {
      return
    }

    const setupEventSource = () => {
      // Clean up previous connection
      if (eventSourceRef.current) {
        eventSourceRef.current.close()
      }

      const params = new URLSearchParams({
        container,
        timestamps: 'true',
      })

      // Add tailLines only for initial connection
      if (settings.tailLines && !isReconnecting) {
        params.append('tailLines', String(settings.tailLines))
      }

      // Add since parameter if specified
      if (settings.since) {
        params.append('since', settings.since)
      }

      // Add sinceTime parameter if specified
      if (settings.sinceTime) {
        params.append('sinceTime', settings.sinceTime)
      }

      // Add last timestamp for reconnection
      if (isReconnecting && lastTimestamp) {
        params.append('sinceTime', lastTimestamp)
      }

      const url = `${axios.defaults.baseURL}/rest-api/v1/insight/aggregator/log/pod/${cluster}/${namespace}/${podName}?${params}`
      const eventSource = new EventSource(url)
      eventSourceRef.current = eventSource

      eventSource.onopen = () => {
        setConnected(true)
        setError(null)
        // Only clear reconnecting flag if connection is stable
        if (reconnectTimeoutRef.current) {
          clearTimeout(reconnectTimeoutRef.current)
          reconnectTimeoutRef.current = null
        }
        reconnectTimeoutRef.current = setTimeout(() => {
          setIsReconnecting(false)
        }, 1000) // Wait for 1 second to ensure connection is stable
      }

      eventSource.onmessage = event => {
        try {
          const logEntry: LogEntry = JSON.parse(event.data)

          if (logEntry.error) {
            setError(logEntry.error)
            return
          }

          // Generate a unique ID for the log entry
          logEntry.id = `${logEntry.timestamp}-${logEntry.content.length}-${logEntry.content.slice(0, 20)}`

          // Track the latest timestamp for reconnection
          if (logEntry.timestamp) {
            setLastTimestamp(logEntry.timestamp)
          }

          setLogs(prev => {
            // Prevent duplicate logs by comparing unique IDs
            const isDuplicate = prev.some(
              existingLog => existingLog.id === logEntry.id,
            )
            if (isDuplicate) {
              return prev
            }
            return [...prev, logEntry]
          })

          // Auto-scroll to bottom for new logs
          if (logsEndRef.current) {
            logsEndRef.current.scrollIntoView({ behavior: 'smooth' })
          }
        } catch (error) {
          console.error('Failed to parse log entry:', error)
        }
      }

      eventSource.onerror = err => {
        console.error('EventSource error:', err)
        const now = Date.now()
        // Prevent rapid reconnection attempts
        if (now - lastErrorRef.current < 1000) {
          return
        }
        lastErrorRef.current = now

        setConnected(false)
        // Only set reconnecting if we're not already in that state
        if (!isReconnecting) {
          setIsReconnecting(true)
        }

        // Clean up and retry connection
        eventSource.close()
        setTimeout(setupEventSource, 1000)
      }
    }

    setupEventSource()

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
      }
      if (eventSourceRef.current) {
        eventSourceRef.current.close()
      }
    }
  }, [cluster, namespace, podName, container, isPaused, settings])

  const { aiOptions } = useSelector((state: any) => state.globalSlice)

  const debouncedDiagnose = useCallback(
    debounce(async () => {
      try {
        if (!logs.length) {
          message.warning(t('LogAggregator.NoLogsSelected'))
          return
        }

        // Reset diagnosis state
        setDiagnosis('')
        setDiagnosisStatus('loading' as DiagnosisStatus)

        // Cancel any existing SSE connection
        if (abortControllerRef.current) {
          abortControllerRef.current.abort()
        }

        // Create new AbortController for this request
        const abortController = new AbortController()
        abortControllerRef.current = abortController

        setStreaming(true)

        // Create new fetch request for diagnosis
        const url = `${axios.defaults.baseURL}/rest-api/v1/insight/aggregator/log/diagnosis/stream`

        // Send POST request and handle SSE response
        fetch(url, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Accept: 'text/event-stream',
          },
          body: JSON.stringify({
            logs: logs.map(log => log.content),
            language: i18n.language,
          }),
          signal: abortController.signal,
        })
          .then(response => {
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
                let streaming = true
                while (streaming) {
                  const { done, value } = await reader.read()

                  if (done) {
                    streaming = false
                    setDiagnosisStatus('complete' as DiagnosisStatus)
                    break
                  }

                  // Decode the chunk and process events
                  const chunk = decoder.decode(value)
                  const events = chunk
                    .split('\n\n')
                    .filter(Boolean)
                    .map(event => event.replace('data: ', ''))

                  for (const event of events) {
                    try {
                      const diagEvent = JSON.parse(event)

                      switch (diagEvent.type) {
                        case 'start':
                          setDiagnosisStatus('streaming' as DiagnosisStatus)
                          break
                        case 'chunk':
                          setDiagnosis(prev => prev + diagEvent.content)
                          break
                        case 'error':
                          streaming = false
                          setDiagnosisStatus('error' as DiagnosisStatus)
                          message.error(diagEvent.content)
                          reader.cancel()
                          break
                        case 'complete':
                          streaming = false
                          setDiagnosisStatus('complete' as DiagnosisStatus)
                          reader.cancel()
                          break
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
                  setDiagnosisStatus('error' as DiagnosisStatus)
                  message.error(t('LogAggregator.DiagnosisConnectionError'))
                }
              }
            }

            processStream()
          })
          .catch(error => {
            if (error.name !== 'AbortError') {
              console.error('Failed to start diagnosis:', error)
              setDiagnosisStatus('error' as DiagnosisStatus)
              message.error(t('LogAggregator.FailedToStartDiagnosis'))
            }
          })
      } catch (error) {
        console.error('Failed to start diagnosis:', error)
        setDiagnosisStatus('error' as DiagnosisStatus)
        message.error(t('LogAggregator.FailedToDiagnoseLogs'))
      } finally {
        setStreaming(false)
      }
    }, 500), // 500ms debounce delay
    [logs, t, i18n.language],
  )

  const handleDiagnose = useCallback(() => {
    debouncedDiagnose()
  }, [debouncedDiagnose])

  const handleDownloadLogs = async () => {
    try {
      // Get visible logs based on current filter/highlight mode
      const visibleLogs = logs.filter(
        log =>
          searchMode === 'highlight' ||
          log.content.toLowerCase().includes(searchText.toLowerCase()),
      )

      // Format logs with timestamps
      const logContent = visibleLogs
        .map(log => `${log.timestamp} ${log.content}`)
        .join('\n')

      // Create blob and download
      const blob = new Blob([logContent], { type: 'text/plain' })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `pod-logs-${podName}-${new Date().toISOString()}.txt`
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
    } catch (error) {
      console.error('Failed to download logs:', error)
      message.error(t('Failed to download logs'))
    }
  }

  const handleSettingsChange = (newSettings: Partial<LogSettings>) => {
    setSettings(prev => ({ ...prev, ...newSettings }))
    // Clear existing logs when settings change
    setLogs([])
    setLastTimestamp(null)

    // Reconnect with new settings
    handleDisconnect()
    handleConnect()
  }

  const handleConnect = () => {
    if (!podName) return

    setError(null)
    setStreaming(true)
    setIsReconnecting(false)

    const params = new URLSearchParams({
      timestamps: 'true',
    })

    // Add tailLines only for initial connection
    if (settings.tailLines && !isReconnecting) {
      params.append('tailLines', String(settings.tailLines))
    }

    // Add since parameter if specified
    if (settings.since) {
      params.append('since', settings.since)
    }

    // Add sinceTime parameter if specified
    if (settings.sinceTime) {
      params.append('sinceTime', settings.sinceTime)
    }

    // Add last timestamp for reconnection
    if (isReconnecting && lastTimestamp) {
      params.append('sinceTime', lastTimestamp)
    }

    const url = `${axios.defaults.baseURL}/rest-api/v1/insight/aggregator/log/pod/${cluster}/${namespace}/${podName}?${params}`
    const eventSource = new EventSource(url)
    eventSourceRef.current = eventSource

    eventSource.onopen = () => {
      setConnected(true)
      setError(null)
      // Only clear reconnecting flag if connection is stable
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
        reconnectTimeoutRef.current = null
      }
      reconnectTimeoutRef.current = setTimeout(() => {
        setIsReconnecting(false)
      }, 1000) // Wait for 1 second to ensure connection is stable
    }

    eventSource.onmessage = event => {
      try {
        const logEntry: LogEntry = JSON.parse(event.data)

        if (logEntry.error) {
          setError(logEntry.error)
          return
        }

        // Generate a unique ID for the log entry
        logEntry.id = `${logEntry.timestamp}-${logEntry.content.length}-${logEntry.content.slice(0, 20)}`

        // Track the latest timestamp for reconnection
        if (logEntry.timestamp) {
          setLastTimestamp(logEntry.timestamp)
        }

        setLogs(prev => {
          // Prevent duplicate logs by comparing unique IDs
          const isDuplicate = prev.some(
            existingLog => existingLog.id === logEntry.id,
          )
          if (isDuplicate) {
            return prev
          }
          return [...prev, logEntry]
        })

        // Auto-scroll to bottom for new logs
        if (logsEndRef.current) {
          logsEndRef.current.scrollIntoView({ behavior: 'smooth' })
        }
      } catch (error) {
        console.error('Failed to parse log entry:', error)
      }
    }

    eventSource.onerror = err => {
      console.error('EventSource error:', err)
      const now = Date.now()
      // Prevent rapid reconnection attempts
      if (now - lastErrorRef.current < 1000) {
        return
      }
      lastErrorRef.current = now

      setConnected(false)
      // Only set reconnecting if we're not already in that state
      if (!isReconnecting) {
        setIsReconnecting(true)
      }

      // Clean up and retry connection
      eventSource.close()
      setTimeout(handleConnect, 1000)
    }
  }

  const handleDisconnect = () => {
    if (eventSourceRef.current) {
      eventSourceRef.current.close()
    }
  }

  const renderSettingsModal = () => (
    <Modal
      title={t('LogAggregator.Settings')}
      open={showSettings}
      onCancel={() => setShowSettings(false)}
      footer={null}
      keyboard={true}
      centered
      maskClosable={true}
      width={400}
    >
      <Space direction="vertical" style={{ width: '100%', padding: '8px 0' }}>
        <div>
          <div style={{ marginBottom: '8px' }}>
            {t('LogAggregator.TailLines')}
          </div>
          <InputNumber
            value={settings.tailLines}
            onChange={value =>
              handleSettingsChange({ tailLines: value as number })
            }
            min={1}
            max={10000}
            style={{ width: '100%' }}
          />
        </div>
        <div>
          <div style={{ marginBottom: '8px' }}>{t('LogAggregator.Since')}</div>
          <Input
            value={settings.since}
            onChange={e => handleSettingsChange({ since: e.target.value })}
            placeholder="e.g. 1h, 2d"
            style={{ width: '100%' }}
          />
        </div>
        <div>
          <div style={{ marginBottom: '8px' }}>
            {t('LogAggregator.SinceTime')}
          </div>
          <DatePicker
            showTime
            value={settings.sinceTime ? dayjs(settings.sinceTime) : null}
            onChange={value =>
              handleSettingsChange({
                sinceTime: value ? value.toISOString() : undefined,
              })
            }
            style={{ width: '100%' }}
          />
        </div>
        <div>
          <div style={{ marginBottom: '8px' }}>
            {t('LogAggregator.SinceTime')}
          </div>
          <Switch
            value={settings.isShowConnect}
            onChange={val => {
              setSettings({ ...settings, isShowConnect: val })
            }}
          />
        </div>
      </Space>
    </Modal>
  )

  const highlightSearchText = (text: string) => {
    if (!searchText) return text

    try {
      const parts = text.split(new RegExp(`(${searchText})`, 'gi'))
      return parts.map((part, i) =>
        part.toLowerCase() === searchText.toLowerCase() ? (
          <span key={i} className={styles.highlight}>
            {part}
          </span>
        ) : (
          part
        ),
      )
    } catch (e) {
      return text
    }
  }

  const filterLogs = (logs: LogEntry[]) => {
    if (!searchText || searchMode === 'highlight') return logs
    return logs.filter(log =>
      log.content.toLowerCase().includes(searchText.toLowerCase()),
    )
  }

  return (
    <div className={styles.podLogs}>
      {error && (
        <Alert className={styles.error} message={error} type="error" showIcon />
      )}
      {container && (
        <div className={styles.toolbar}>
          <Select
            value={container}
            onChange={setContainer}
            placeholder={t('LogAggregator.SelectContainer')}
            options={containers.map(c => ({ value: c, label: c }))}
          />
          <Space
            direction="horizontal"
            size={8}
            style={{ display: 'flex', alignItems: 'center' }}
          >
            <Input.Search
              placeholder={t('LogAggregator.SearchPlaceholder')}
              style={{ width: 180 }}
              value={searchText}
              onChange={e => setSearchText(e.target.value)}
              allowClear
              className={styles.searchInput}
              suffix={
                <Tooltip
                  title={
                    searchMode === 'highlight'
                      ? t('LogAggregator.SearchModeHighlight')
                      : t('LogAggregator.SearchModeFilter')
                  }
                >
                  <span
                    className={`${styles.searchButton} ${
                      searchMode === 'filter' ? styles.active : ''
                    }`}
                    onClick={() =>
                      setSearchMode(mode =>
                        mode === 'highlight' ? 'filter' : 'highlight',
                      )
                    }
                  >
                    {searchMode === 'highlight' ? (
                      <HighlightOutlined style={{ fontSize: '14px' }} />
                    ) : (
                      <FilterOutlined style={{ fontSize: '14px' }} />
                    )}
                  </span>
                </Tooltip>
              }
            />
            <Tooltip
              title={
                isPaused
                  ? t('LogAggregator.ResumeLogs')
                  : t('LogAggregator.PauseLogs')
              }
            >
              <Button
                type="text"
                className={styles.actionButton}
                icon={
                  isPaused ? <PlayCircleOutlined /> : <PauseCircleOutlined />
                }
                onClick={() => setIsPaused(!isPaused)}
              />
            </Tooltip>
            <Tooltip title={t('LogAggregator.ClearLogs')}>
              <Button
                type="text"
                className={styles.actionButton}
                icon={<ClearOutlined />}
                onClick={() => setLogs([])}
              />
            </Tooltip>
            {aiOptions?.AIAuthToken && (
              <Tooltip title={t('LogAggregator.DiagnoseLogs')}>
                <Button
                  type="text"
                  className={styles.actionButton}
                  icon={<span className={styles.magicWand}>✨</span>}
                  onClick={handleDiagnose}
                  disabled={logs.length === 0 || isStreaming}
                />
              </Tooltip>
            )}
            <Tooltip title={t('LogAggregator.DownloadLogs')}>
              <Button
                type="text"
                className={styles.actionButton}
                icon={<DownloadOutlined />}
                onClick={handleDownloadLogs}
              />
            </Tooltip>
            <Tooltip title={t('LogAggregator.Settings')}>
              <Button
                type="text"
                className={styles.actionButton}
                icon={<SettingOutlined />}
                onClick={() => setShowSettings(true)}
              />
            </Tooltip>
            {settings?.isShowConnect && (
              <Tooltip
                title={
                  isConnected
                    ? t('LogAggregator.ConnectedTip', { container })
                    : t('LogAggregator.DisconnectedTip')
                }
              >
                <div className={styles.connectionStatus}>
                  <Badge
                    status={isConnected ? 'success' : 'error'}
                    text={t(
                      isConnected
                        ? 'LogAggregator.Connected'
                        : 'LogAggregator.Disconnected',
                    )}
                  />
                </div>
              </Tooltip>
            )}
          </Space>
        </div>
      )}
      <div className={styles.content}>
        <div className={styles.logsContainer}>
          <div className={styles.logs}>
            {filterLogs(logs).map((log, index) => (
              <div key={index} className={styles.logEntry}>
                <span className={styles.timestamp}>{log.timestamp}</span>
                <span
                  className={log.error ? styles.errorContent : styles.content}
                >
                  {searchMode === 'highlight'
                    ? highlightSearchText(log.content)
                    : log.content}
                </span>
              </div>
            ))}
            <div ref={logsEndRef} />
          </div>
          {diagnosisStatus !== 'idle' && (
            <div className={styles.diagnosisPanel}>
              <div className={styles.diagnosisHeader}>
                <Space>
                  <RobotOutlined />
                  {t('LogAggregator.DiagnosisResult')}
                </Space>
                <Space>
                  {diagnosisStatus === 'streaming' && (
                    <Tooltip
                      title={t('LogAggregator.StopDiagnosis')}
                      placement="bottom"
                    >
                      <Button
                        type="text"
                        className={styles.stopButton}
                        icon={<PoweroffOutlined />}
                        onClick={() => {
                          if (abortControllerRef.current) {
                            abortControllerRef.current.abort()
                            setDiagnosisStatus('complete' as DiagnosisStatus)
                          }
                        }}
                      />
                    </Tooltip>
                  )}
                  <Button
                    type="text"
                    icon={<CloseOutlined />}
                    onClick={() =>
                      setDiagnosisStatus('idle' as DiagnosisStatus)
                    }
                  />
                </Space>
              </div>
              <div className={styles.diagnosisBody}>
                {diagnosisStatus === ('loading' as DiagnosisStatus) ? (
                  <div className={styles.diagnosisContent}>
                    <div className={styles.diagnosisLoading}>
                      <Spin />
                      <span>{t('LogAggregator.PreparingDiagnosis')}</span>
                    </div>
                  </div>
                ) : diagnosisStatus === ('streaming' as DiagnosisStatus) ? (
                  <div className={styles.diagnosisContent}>
                    <Markdown>{diagnosis}</Markdown>
                    <div
                      ref={diagnosisEndRef}
                      style={{ float: 'left', clear: 'both' }}
                    />
                  </div>
                ) : diagnosisStatus === ('error' as DiagnosisStatus) ? (
                  <div className={styles.diagnosisError}>
                    <Alert
                      type="error"
                      message={t('LogAggregator.DiagnosisError')}
                      description={
                        diagnosis || t('LogAggregator.TryAgainLater')
                      }
                    />
                  </div>
                ) : (
                  <div className={styles.diagnosisContent}>
                    <Markdown>{diagnosis}</Markdown>
                  </div>
                )}
              </div>
            </div>
          )}
        </div>
      </div>
      {renderSettingsModal()}
    </div>
  )
}

export default PodLogs
