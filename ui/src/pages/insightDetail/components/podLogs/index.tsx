import React, { useEffect, useRef, useState } from 'react'
import { Select, Space, Button, Alert, Badge, Tooltip } from 'antd'
import {
  PauseCircleOutlined,
  PlayCircleOutlined,
  ClearOutlined,
} from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import yaml from 'js-yaml'
import styles from './styles.module.less'

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

const PodLogs: React.FC<PodLogsProps> = ({
  cluster,
  namespace,
  podName,
  yamlData,
}) => {
  const { t } = useTranslation()
  const [container, setContainer] = useState<string>('')
  const [containers, setContainers] = useState<string[]>([])
  const [logs, setLogs] = useState<LogEntry[]>([])
  const [isPaused, setPaused] = useState(false)
  const [isConnected, setConnected] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const eventSourceRef = useRef<EventSource | null>(null)
  const logsEndRef = useRef<HTMLDivElement>(null)

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
        setError(t('Failed to parse pod details'))
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

    const url = `/rest-api/v1/insight/aggregator/pod/${cluster}/${namespace}/${podName}/log?container=${container}`
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

  const handlePause = () => {
    setPaused(!isPaused)
  }

  const handleClear = () => {
    setLogs([])
    setError(null)
  }

  return (
    <div className={styles.podLogs}>
      <div className={styles.toolbar}>
        <Space>
          <Select
            value={container}
            onChange={setContainer}
            style={{ width: 200 }}
            placeholder={t('Select container')}
          >
            {containers.map(c => (
              <Select.Option key={c} value={c}>
                {c}
              </Select.Option>
            ))}
          </Select>
          <Tooltip title={isPaused ? t('Resume logs') : t('Pause logs')}>
            <Button
              type={isPaused ? 'default' : 'primary'}
              icon={isPaused ? <PlayCircleOutlined /> : <PauseCircleOutlined />}
              onClick={handlePause}
            />
          </Tooltip>
          <Tooltip title={t('Clear logs')}>
            <Button icon={<ClearOutlined />} onClick={handleClear} />
          </Tooltip>
          <Badge
            status={isConnected ? 'success' : 'error'}
            text={isConnected ? t('Connected') : t('Disconnected')}
          />
        </Space>
      </div>

      {error && (
        <Alert
          type="error"
          message={error}
          className={styles.error}
          closable
          onClose={() => setError(null)}
        />
      )}

      <div className={styles.logsContainer}>
        {logs.map((log, index) => (
          <div key={index} className={styles.logEntry}>
            <span className={styles.content}>{log.content}</span>
          </div>
        ))}
        <div ref={logsEndRef} />
      </div>
    </div>
  )
}

export default PodLogs
