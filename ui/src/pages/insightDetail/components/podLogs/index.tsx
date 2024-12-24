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
} from 'antd'
import {
  PauseCircleOutlined,
  PlayCircleOutlined,
  ClearOutlined,
  RobotOutlined,
  CloseOutlined,
  PoweroffOutlined,
} from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import yaml from 'js-yaml'
import axios from 'axios'
import styles from './styles.module.less'
import Markdown from 'react-markdown'
import { useSelector } from 'react-redux'
import { debounce } from 'lodash'

interface LogEntry {
  timestamp: string
  content: string
  error?: string
}

interface PodLogsProps {
  cluster: string
  namespace: string
  podName: string
  yamlData?: string
}

type DiagnosisStatus = 'idle' | 'init' | 'streaming' | 'complete' | 'error'

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
  const [diagnosisStatus, setDiagnosisStatus] =
    useState<DiagnosisStatus>('idle')
  const [diagnosis, setDiagnosis] = useState('')
  const [isStreaming, setStreaming] = useState(false)
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

    // Clean up previous connection
    if (eventSourceRef.current) {
      eventSourceRef.current.close()
      setLogs([]) // Clear logs when switching containers or reconnecting
    }

    const url = `${axios.defaults.baseURL}/rest-api/v1/insight/aggregator/log/pod/${cluster}/${namespace}/${podName}?container=${container}`
    const eventSource = new EventSource(url)
    eventSourceRef.current = eventSource

    eventSource.onopen = () => {
      setConnected(true)
      setError(null)
    }

    eventSource.onmessage = event => {
      try {
        const logEntry: LogEntry = JSON.parse(event.data)

        if (logEntry.error) {
          setError(logEntry.error)
          return
        }

        setLogs(prev => [...prev, logEntry])

        // Auto-scroll to bottom
        if (logsEndRef.current) {
          logsEndRef.current.scrollIntoView({ behavior: 'smooth' })
        }
      } catch (error) {
        console.error('Failed to parse log entry:', error)
      }
    }

    eventSource.onerror = err => {
      console.error('EventSource error:', err)
      setConnected(false)
      // SSE will automatically reconnect, no manual handling needed
    }

    return () => {
      eventSource.close()
    }
  }, [cluster, namespace, podName, container, isPaused])

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

  // Auto scroll to the bottom of diagnosis content
  useEffect(() => {
    if (diagnosisEndRef.current && diagnosisStatus === 'streaming') {
      // Use requestAnimationFrame to ensure scrolling after the next frame render
      requestAnimationFrame(() => {
        diagnosisEndRef.current?.scrollIntoView({
          behavior: 'auto', // Use auto for faster response
          block: 'end', // Ensure scrolling to bottom
        })
      })
    }
  }, [diagnosis]) // Only listen to diagnosis changes

  // Clean up diagnosis connection on unmount
  useEffect(() => {
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort()
      }
    }
  }, [])

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
          </Space>
        </div>
      )}
      <div className={styles.content}>
        <div className={styles.logsContainer}>
          <div className={styles.logs}>
            {logs.map((log, index) => (
              <div key={index} className={styles.logEntry}>
                <span className={styles.timestamp}>{log.timestamp}</span>
                <span
                  className={log.error ? styles.errorContent : styles.content}
                >
                  {log.content}
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
    </div>
  )
}

export default PodLogs
