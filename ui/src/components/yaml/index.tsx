import React, { useEffect, useRef, useState, useCallback } from 'react'
import type { LegacyRef } from 'react'
import { Alert, Button, message, Modal, Space, Spin, Tooltip } from 'antd'
import { Resizable } from 're-resizable'
import { useTranslation } from 'react-i18next'
import {
  CopyOutlined,
  CloseOutlined,
  PoweroffOutlined,
  FullscreenExitOutlined,
  FullscreenOutlined,
  ExpandOutlined,
  ReloadOutlined,
} from '@ant-design/icons'
import hljs from 'highlight.js'
import yaml from 'js-yaml'
import 'highlight.js/styles/lightfair.css'
import { yaml2json } from '@/utils/tools'
import { useSelector } from 'react-redux'
import Markdown from 'react-markdown'
import axios from 'axios'
import { FullScreen, useFullScreenHandle } from 'react-full-screen'
import i18n from '@/i18n'
import aiSummarySvg from '@/assets/ai-summary.svg'

import styles from './styles.module.less'

// eslint-disable-next-line @typescript-eslint/no-var-requires
hljs.registerLanguage('yaml', require('highlight.js/lib/languages/yaml'))

type InterpretStatus =
  | 'idle'
  | 'init'
  | 'streaming'
  | 'complete'
  | 'error'
  | 'loading'

type IProps = {
  data: any
  height?: string | number
  onRefresh?: () => void
}

const Yaml = (props: IProps) => {
  const { t } = useTranslation()
  const handle = useFullScreenHandle()
  const yamlRef = useRef<LegacyRef<HTMLDivElement> | undefined>()
  const diagnosisContentRef = useRef<HTMLDivElement>(null)
  const interpretEndRef = useRef<HTMLDivElement>(null)
  const contentRef = useRef<HTMLDivElement>(null)
  const observerRef = useRef<MutationObserver | null>(null)
  const { data, onRefresh } = props
  const [moduleHeight, setModuleHeight] = useState<number>(500)
  const [interpretStatus, setInterpretStatus] =
    useState<InterpretStatus>('idle')
  const [interpret, setInterpret] = useState('')
  const [isStreaming, setStreaming] = useState(false)
  const abortControllerRef = useRef<AbortController | null>(null)
  const { aiOptions } = useSelector((state: any) => state.globalSlice)
  const isAIEnabled = aiOptions?.AIModel && aiOptions?.AIAuthToken
  const [open, setOpen] = useState(false)

  useEffect(() => {
    onRefresh?.()
  }, [])

  useEffect(() => {
    const yamlStatusJson = yaml2json(data)
    if (yamlRef.current && yamlStatusJson?.data) {
      ;(yamlRef.current as unknown as HTMLElement).innerHTML = hljs.highlight(
        'yaml',
        yaml.dump(yamlStatusJson?.data),
      ).value
    }
  }, [data])

  // Function to scroll to the bottom of the container
  const scrollToBottom = useCallback(() => {
    if (diagnosisContentRef.current && interpretStatus === 'streaming') {
      const container = diagnosisContentRef.current
      const scrollHeight = container.scrollHeight
      const height = container.clientHeight
      const maxScroll = scrollHeight - height
      container.scrollTo({
        top: maxScroll,
        behavior: 'auto',
      })
    }
  }, [interpretStatus])

  // Watch for content changes
  useEffect(() => {
    if (interpretStatus === 'streaming' && diagnosisContentRef.current) {
      if (observerRef.current) {
        observerRef.current.disconnect()
      }

      const observer = new MutationObserver(() => {
        scrollToBottom()
      })

      observer.observe(diagnosisContentRef.current, {
        childList: true,
        subtree: true,
        characterData: true,
      })

      observerRef.current = observer

      return () => {
        observer.disconnect()
      }
    }
  }, [interpretStatus, scrollToBottom])

  // Scroll when content updates
  useEffect(() => {
    if (interpretStatus === 'streaming') {
      scrollToBottom()
    }
  }, [interpret, scrollToBottom, interpretStatus])

  function copy() {
    const textarea = document.createElement('textarea')
    textarea.value = data
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    message.success(t('CopySuccess'))
    document.body.removeChild(textarea)
  }

  const handleInterpret = async () => {
    try {
      if (!data) {
        message.warning(t('YAML.NoContent'))
        return
      }

      // Reset interpret state
      setInterpret('')
      setInterpretStatus('loading' as InterpretStatus)

      // Cancel any existing SSE connection
      if (abortControllerRef.current) {
        abortControllerRef.current.abort()
      }

      // Create new AbortController for this request
      const abortController = new AbortController()
      abortControllerRef.current = abortController

      setStreaming(true)

      // Create new fetch request for interpret
      const url = `${axios.defaults.baseURL}/rest-api/v1/insight/yaml/interpret/stream`

      // Send POST request and handle SSE response
      fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Accept: 'text/event-stream',
        },
        body: JSON.stringify({
          yaml: data,
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
                  setInterpretStatus('complete' as InterpretStatus)
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
                    const interpretEvent = JSON.parse(event)

                    switch (interpretEvent.type) {
                      case 'start':
                        setInterpretStatus('streaming' as InterpretStatus)
                        break
                      case 'chunk':
                        setInterpret(prev => prev + interpretEvent.content)
                        if (interpretEndRef.current) {
                          interpretEndRef.current.scrollIntoView({
                            behavior: 'smooth',
                          })
                        }
                        break
                      case 'error':
                        streaming = false
                        setInterpretStatus('error' as InterpretStatus)
                        message.error(interpretEvent.content)
                        reader.cancel()
                        break
                      case 'complete':
                        streaming = false
                        setInterpretStatus('complete' as InterpretStatus)
                        reader.cancel()
                        break
                    }
                  } catch (error) {
                    console.error('Failed to parse interpret event:', error)
                  }
                }
              }
            } catch (error) {
              if (error.name === 'AbortError') {
                console.log('Interpret stream aborted')
              } else {
                console.error('Error reading stream:', error)
                setInterpretStatus('error' as InterpretStatus)
                message.error(t('YAML.InterpretConnectionError'))
              }
            }
          }

          processStream()
        })
        .catch(error => {
          if (error.name !== 'AbortError') {
            console.error('Failed to start interpret:', error)
            setInterpretStatus('error' as InterpretStatus)
            message.error(t('YAML.FailedToStartInterpret'))
          }
        })
    } catch (error) {
      console.error('Failed to start interpret:', error)
      setInterpretStatus('error' as InterpretStatus)
      message.error(t('YAML.FailedToInterpret'))
    } finally {
      setStreaming(false)
    }
  }
  const contentToTopHeight = contentRef.current?.getBoundingClientRect()?.top
  const dotToTopHeight = interpretEndRef.current?.getBoundingClientRect()?.top

  function renderAiAnalysis(isDailog) {
    return (
      interpretStatus !== 'idle' && (
        <div
          className={styles.yaml_content_diagnosisPanel}
          style={{
            height: isDailog ? 700 : moduleHeight,
            ...(isDailog ? { width: '100%' } : {}),
          }}
        >
          <div className={styles.yaml_content_diagnosisHeader}>
            <Space>
              <div className={styles.yaml_content_diagnosisHeader_aiIcon}>
                <img src={aiSummarySvg} alt="ai summary" />
              </div>
              {t('YAML.InterpretResult')}
            </Space>
            <Space>
              {interpretStatus === 'streaming' && (
                <Tooltip title={t('YAML.StopInterpret')} placement="bottom">
                  <Button
                    type="text"
                    className={styles.stopButton}
                    icon={<PoweroffOutlined />}
                    onClick={() => {
                      if (abortControllerRef.current) {
                        abortControllerRef.current.abort()
                        setInterpretStatus('complete' as InterpretStatus)
                      }
                    }}
                  />
                </Tooltip>
              )}
              {isDailog ? null : (
                <>
                  <Button
                    type="text"
                    icon={<ExpandOutlined />}
                    onClick={() => {
                      setOpen(!open)
                    }}
                  />
                </>
              )}
              <Button
                type="text"
                icon={<CloseOutlined />}
                onClick={() => {
                  if (isDailog) {
                    setOpen(false)
                  } else {
                    setInterpretStatus('idle' as InterpretStatus)
                  }
                }}
              />
            </Space>
          </div>
          <div
            className={styles.yaml_content_diagnosisBody}
            style={isDailog ? { width: '100%' } : {}}
          >
            <div
              className={styles.yaml_content_diagnosisContent}
              ref={contentRef}
            >
              {interpretStatus === 'loading' ||
              (interpretStatus === 'streaming' && !interpret) ? (
                <div className={styles.yaml_content_diagnosisLoading}>
                  <Spin />
                  <p>{t('EventAggregator.DiagnosisInProgress')}</p>
                </div>
              ) : interpretStatus === 'streaming' ? (
                <>
                  <Markdown>{interpret}</Markdown>
                  <div
                    ref={interpretEndRef}
                    style={{
                      float: 'left',
                      clear: 'both',
                    }}
                  />
                </>
              ) : interpretStatus === 'error' ? (
                <Alert
                  type="error"
                  message={t('EventAggregator.DiagnosisFailed')}
                  description={t('EventAggregator.TryAgainLater')}
                />
              ) : (
                <Markdown>{interpret}</Markdown>
              )}
            </div>
            {interpretStatus === 'streaming' && interpret && (
              <div
                className={`${styles.yaml_content_streamingIndicator} ${dotToTopHeight - contentToTopHeight + 53 - moduleHeight >= 0 ? styles.yaml_content_streamingIndicatorFixed : ''}`}
              >
                <span className={styles.dot}></span>
                <span className={styles.dot}></span>
                <span className={styles.dot}></span>
              </div>
            )}
          </div>
        </div>
      )
    )
  }

  return (
    <div style={{ paddingBottom: 20 }}>
      <Resizable
        defaultSize={{
          height: moduleHeight,
        }}
        onResizeStop={(e, direction, ref, d) => {
          const newModuleHeight = moduleHeight + d.height
          setModuleHeight(newModuleHeight)
        }}
        handleStyles={{
          bottom: {
            bottom: 0,
            height: '6px',
            cursor: 'row-resize',
            background: 'transparent',
            transition: 'background 0.3s ease',
          },
        }}
        handleClasses={{
          bottom: styles.resizeHandle,
        }}
      >
        <div className={styles.yaml_content} style={{ height: props?.height }}>
          <FullScreen handle={handle} className={styles.fullScreenConatiner}>
            <div className={styles.yaml_container}>
              <div className={styles.copy}>
                <Space>
                  {data && (
                    <>
                      {!handle.active && (
                        <Tooltip title={t('LogAggregator.FullScreen')}>
                          <Button
                            type="text"
                            icon={<FullscreenOutlined />}
                            onClick={handle.enter}
                          />
                        </Tooltip>
                      )}
                      {!handle.active && onRefresh && (
                        <Tooltip title={t('YAML.Refresh')}>
                          <Button
                            type="text"
                            onClick={onRefresh}
                            icon={<ReloadOutlined />}
                          />
                        </Tooltip>
                      )}
                      {!handle.active && (
                        <Button
                          type="primary"
                          size="small"
                          onClick={copy}
                          disabled={!data}
                          icon={<CopyOutlined />}
                        >
                          {t('Copy')}
                        </Button>
                      )}
                      {handle.active && (
                        <Tooltip title={t('LogAggregator.FullScreen')}>
                          <Button
                            type="text"
                            icon={<FullscreenExitOutlined />}
                            onClick={handle.exit}
                          />
                        </Tooltip>
                      )}
                    </>
                  )}
                  {isAIEnabled && !handle.active && (
                    <Tooltip title={t('YAML.Interpret')}>
                      <Button
                        type="primary"
                        size="small"
                        icon={<span className={styles.magicWand}>✨</span>}
                        onClick={handleInterpret}
                        disabled={!data || isStreaming}
                      >
                        <span style={{ fontSize: 14, marginTop: 2 }}>
                          {t('YAML.Interpret')}
                        </span>
                      </Button>
                    </Tooltip>
                  )}
                </Space>
              </div>
              <div
                className={styles.yaml_box}
                style={{ height: props?.height }}
                ref={yamlRef as any}
              />
            </div>
          </FullScreen>
          {renderAiAnalysis(false)}
        </div>
      </Resizable>
      <Modal
        styles={{
          content: {
            padding: 0,
            borderRadius: 16,
            background:
              'linear-gradient(135deg, rgba(147, 112, 219, 0.15) 0%, rgba(43, 29, 60, 0.95) 100%)',
          },
        }}
        centered
        closable={false}
        width={'80%'}
        height={700}
        open={open}
        footer={null}
      >
        <div style={{ overflow: 'auto' }}>{renderAiAnalysis(true)}</div>
      </Modal>
    </div>
  )
}

export default Yaml
