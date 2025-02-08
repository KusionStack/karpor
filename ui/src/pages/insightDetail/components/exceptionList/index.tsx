import React, { useEffect, useState, useCallback, useRef } from 'react'
import { Button, Empty, Tag, Space, Tooltip, Alert, Spin } from 'antd'
import { useTranslation } from 'react-i18next'
import {
  ArrowRightOutlined,
  UpOutlined,
  DownOutlined,
  PoweroffOutlined,
  CloseOutlined,
  RedoOutlined,
} from '@ant-design/icons'
import Loading from '@/components/loading'
import { SEVERITY_MAP } from '@/utils/constants'
import ExceptionStat from '../exceptionStat'
import { useSelector } from 'react-redux'
import axios from 'axios'
import { debounce } from 'lodash'
import Markdown from 'react-markdown'
import aiSummarySvg from '@/assets/ai-summary.svg'
import classNames from 'classnames'
import { IssueGroup, IssueResponse } from '../../types'

import styles from './style.module.less'

type IProps = {
  exceptionList: IssueResponse
  rescan: () => void
  showDrawer: () => void
  onItemClick: (val: IssueGroup) => void
  exceptionStat: any
  auditLoading: boolean
}

const ExceptionList = ({
  exceptionList,
  rescan,
  showDrawer,
  onItemClick,
  exceptionStat,
  auditLoading,
}: IProps) => {
  const [selectedEventId, setSelectedEventId] = useState<any>()
  const [top5List, setTop5list] = useState([])
  const [currentKey, setCurrentKey] = useState('All')
  const { t, i18n } = useTranslation()
  const [isShowList, setIsShowList] = useState(true)

  // AI interpret states
  const [interpret, setInterpret] = useState('')
  const [interpretStatus, setInterpretStatus] = useState<
    'idle' | 'loading' | 'streaming' | 'complete' | 'error'
  >('idle')
  const [isStreaming, setStreaming] = useState(false)
  const interpretEndRef = useRef<HTMLDivElement>(null)
  const abortControllerRef = useRef<AbortController | null>(null)

  const { aiOptions } = useSelector((state: any) => state.globalSlice)
  const isAIEnabled = aiOptions?.AIModel && aiOptions?.AIAuthToken
  const exceptionRef = useRef<HTMLDivElement>(null)
  const contentRef = useRef<HTMLDivElement>(null)
  const interpretBodyRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    setIsShowList(
      exceptionList?.issueGroups && exceptionList?.issueGroups?.length > 0,
    )
    if (currentKey === 'All') {
      const defaultTop5List = exceptionList?.issueGroups?.slice(0, 5)
      setTop5list(defaultTop5List)
    } else {
      const tmp = exceptionList?.issueGroups?.filter(
        (item: any) => item?.issue?.severity === currentKey,
      )
      const top5Tmp = tmp?.slice(0, 5)
      setTop5list(top5Tmp)
    }
  }, [currentKey, exceptionList?.issueGroups])

  function onClickTable(key) {
    setCurrentKey(key)
  }
  function handleClick() {
    setIsShowList(!isShowList)
  }
  const debouncedInterpret = useCallback(
    debounce(async () => {
      try {
        if (!top5List?.length) {
          // message.warning(t('ExceptionList.NoIssues'))
          return
        }

        // Reset interpret state
        setInterpret('')
        setInterpretStatus('loading')
        setStreaming(true)

        // Cancel any existing SSE connection
        if (abortControllerRef.current) {
          abortControllerRef.current.abort()
        }

        // Create new AbortController for this request
        const abortController = new AbortController()
        abortControllerRef.current = abortController

        // Create new fetch request for interpret
        const url = `${axios.defaults.baseURL}/rest-api/v1/insight/issue/interpret/stream`

        // Send POST request and handle SSE response
        const response = await fetch(url, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Accept: 'text/event-stream',
          },
          body: JSON.stringify({
            auditData: {
              issueTotal: exceptionList.issueTotal,
              resourceTotal: exceptionList.resourceTotal,
              bySeverity: exceptionList.bySeverity,
              issueGroups: top5List.map(item => ({
                issue: {
                  scanner: item.issue.scanner,
                  severity: severityMap[item.issue.severity] ?? 0,
                  title: item.issue.title,
                  message: item.issue.message,
                },
                resourceGroups: item.resourceGroups,
              })),
            },
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
                setInterpretStatus('complete')
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
                      setInterpretStatus('streaming')
                      break
                    case 'chunk':
                      setInterpret(prev => prev + diagEvent.content)
                      // Scroll to bottom of interpret
                      if (interpretEndRef.current) {
                        const interpretBodyScrollHeight =
                          interpretBodyRef.current.scrollHeight
                        interpretBodyRef.current.scrollTo({
                          top: interpretBodyScrollHeight,
                          behavior: 'smooth',
                        })
                      }
                      break
                    case 'error':
                      streaming = false
                      setInterpretStatus('error')
                      setStreaming(false)
                      // message.error(t('ExceptionList.InterpretConnectionError'))
                      reader.cancel()
                      return
                    case 'complete':
                      streaming = false
                      setInterpretStatus('complete')
                      setStreaming(false)
                      reader.cancel()
                      return
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
              setInterpretStatus('error')
              setStreaming(false)
              // message.error(t('ExceptionList.InterpretConnectionError'))
            }
          }
        }

        processStream()
      } catch (error) {
        console.error('Failed to start interpret:', error)
        setInterpretStatus('error')
        setStreaming(false)
        // message.error(t('ExceptionList.FailedToInterpretLogs'))
      }
    }, 500),
    [top5List, t, i18n.language],
  )

  const startInterpret = useCallback(() => {
    debouncedInterpret()
  }, [debouncedInterpret])

  const contentToTopHeight = contentRef.current?.getBoundingClientRect()?.top
  const dotToTopHeight = interpretEndRef.current?.getBoundingClientRect()?.top

  const renderInterpretWindow = () => {
    if (interpretStatus === 'idle') {
      return null
    }

    return (
      <div
        className={styles.interpret_panel}
        style={{ height: exceptionRef?.current?.offsetHeight }}
      >
        <div className={styles.interpret_header}>
          <Space>
            <div className={styles.interpret_header_aiIcon}>
              <img src={aiSummarySvg} alt="ai summary" />
            </div>
            {t('ExceptionList.InterpretResult')}
          </Space>
          <Space>
            {isStreaming && (
              <Tooltip
                title={t('ExceptionList.StopInterpret')}
                placement="bottom"
              >
                <Button
                  type="text"
                  className={styles.interpret_header_stopButton}
                  icon={<PoweroffOutlined />}
                  onClick={() => {
                    if (abortControllerRef.current) {
                      abortControllerRef.current.abort()
                      setInterpretStatus('complete')
                      setStreaming(false)
                    }
                  }}
                />
              </Tooltip>
            )}
            <Button
              type="text"
              icon={<CloseOutlined />}
              onClick={() => {
                if (abortControllerRef.current) {
                  abortControllerRef.current.abort()
                }
                setInterpretStatus('idle')
                setInterpret('')
                setStreaming(false)
              }}
            />
          </Space>
        </div>
        <div className={styles.interpret_body} ref={interpretBodyRef}>
          <div className={styles.interpret_content} ref={contentRef}>
            {interpretStatus === 'loading' ||
            (interpretStatus === 'streaming' && !interpret) ? (
              <div className={styles.interpret_loading}>
                <Spin />
                <p>{t('ExceptionList.InterpretInProgress')}</p>
              </div>
            ) : interpretStatus === 'error' ? (
              <Alert
                type="error"
                message={t('ExceptionList.InterpretFailed')}
                description={t('ExceptionList.TryAgainLater')}
              />
            ) : (
              <>
                <Markdown>{interpret}</Markdown>
                <div ref={interpretEndRef} />
              </>
            )}
          </div>
          {interpretStatus === 'streaming' && interpret && (
            <div
              className={`${styles.streaming_indicator} ${dotToTopHeight - contentToTopHeight + 57 - exceptionRef.current?.offsetHeight >= 0 ? styles.streaming_indicatorFixed : ''}`}
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

  const severityMap = {
    SAFE: 0,
    LOW: 1,
    MEDIUM: 2,
    HIGH: 3,
    CRITICAL: 5,
  }

  return (
    <div className={styles.exceptionContainer}>
      <div className={styles.exception} ref={exceptionRef}>
        <div className={styles.header}>
          <div className={styles.header_left}>
            <ExceptionStat
              currentKey={currentKey}
              statData={{
                all: exceptionStat?.issuesTotal || 0,
                high: exceptionStat?.severityStatistic?.High || 0,
                medium: exceptionStat?.severityStatistic?.Medium || 0,
                low: exceptionStat?.severityStatistic?.Low || 0,
              }}
              onClickTable={onClickTable}
            />
          </div>
          <div className={styles.header_right}>
            {isAIEnabled && (
              <Tooltip title={t('ExceptionList.Interpret')}>
                <Button
                  type="text"
                  className={styles.header_right_ai_button}
                  icon={<span className={styles.magic_wand}>✨</span>}
                  onClick={startInterpret}
                  loading={interpretStatus === 'streaming'}
                  disabled={
                    !top5List?.length || interpretStatus === 'streaming'
                  }
                />
              </Tooltip>
            )}
            <Tooltip title={t('ExceptionList.Rescan')}>
              <Button
                className={styles.header_right_action}
                icon={<RedoOutlined />}
                onClick={rescan}
                type="text"
              />
            </Tooltip>
            <Tooltip
              title={
                isShowList
                  ? t('ExceptionList.Collapse')
                  : t('ExceptionList.Expand')
              }
            >
              <Button
                className={styles.header_right_action}
                icon={isShowList ? <UpOutlined /> : <DownOutlined />}
                onClick={handleClick}
                type="text"
              />
            </Tooltip>
          </div>
        </div>
        {isShowList && (
          <div className={styles.body}>
            {auditLoading ? (
              <div className={styles.loading_box}>
                <Loading />
              </div>
            ) : top5List && top5List?.length > 0 ? (
              <div className={styles.content_wrapper}>
                <div
                  className={classNames(styles.list_container, {
                    [styles.with_interpret]: interpretStatus !== 'idle',
                  })}
                >
                  {top5List?.map(item => {
                    const uniqueKey = `${item?.issue?.title}_${item?.issue?.message}_${item?.issue?.scanner}_${item?.issue?.severity}`
                    return (
                      <div
                        key={uniqueKey}
                        className={styles.item}
                        onMouseMove={() => setSelectedEventId(uniqueKey)}
                        onMouseOut={() => setSelectedEventId(undefined)}
                        onClick={() => onItemClick(item)}
                      >
                        {selectedEventId === uniqueKey && (
                          <div className={styles.itme_tip}>
                            {t('ExceptionList.ViewIssueDetail')}
                          </div>
                        )}

                        <div className={styles.itme_content}>
                          <div className={styles.left}>
                            <Tag
                              color={
                                SEVERITY_MAP?.[item?.issue?.severity]?.color
                              }
                            >
                              {t(
                                `ExceptionList.${SEVERITY_MAP?.[item?.issue?.severity]?.text}`,
                              )}
                            </Tag>
                          </div>
                          <div className={styles.right}>
                            <div className={styles.right_top}>
                              <div>
                                <span className={styles.title}>
                                  {item?.issue?.title}
                                </span>
                                <span>
                                  {t('ExceptionList.Occur')}&nbsp;
                                  <span
                                    style={{
                                      fontWeight: 'bold',
                                      color: '#646566',
                                    }}
                                  >
                                    {item?.resourceGroups?.length}
                                  </span>
                                  &nbsp;{t('ExceptionList.Times')}
                                </span>
                                <span style={{ color: '#000' }}>&nbsp;|</span>
                                &nbsp;{t('ExceptionList.CollectedFrom')}
                              </div>
                              <div className={styles.tool}>
                                <ArrowRightOutlined />
                                &nbsp;{item?.issue?.scanner}&nbsp;
                                {t('ExceptionList.Tool')}
                              </div>
                            </div>
                            <div className={styles.right_bottom}>
                              <div className={styles.label}>
                                {t('ExceptionList.Description')}:&nbsp;
                              </div>
                              <div className={styles.value}>
                                {item?.issue?.message}
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    )
                  })}
                  <div className={styles.footer}>
                    <span className={styles.btn} onClick={showDrawer}>
                      {t('ExceptionList.CheckAllIssues')}
                      <ArrowRightOutlined />
                    </span>
                  </div>
                </div>
              </div>
            ) : (
              <div className={styles.content_empty}>
                <Empty description={`${t('ExceptionList.NoIssuesFound')}`} />
              </div>
            )}
          </div>
        )}
      </div>
      {renderInterpretWindow()}
    </div>
  )
}

export default ExceptionList
