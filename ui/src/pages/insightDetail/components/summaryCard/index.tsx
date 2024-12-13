import React, { useEffect, useRef, useState } from 'react'
import { Progress, message, Tooltip } from 'antd'
import { useLocation } from 'react-router-dom'
import queryString from 'query-string'
import { useTranslation } from 'react-i18next'
import { getDataType } from '@/utils/tools'
import dayjs from 'dayjs'

import styles from './styles.module.less'

// Types for metrics data
interface MetricPoint {
  timestamp: string // ISO string
  value: number
}

interface ResourceMetrics {
  points: MetricPoint[]
}

// Unit conversion functions
const convertUnits = {
  // Convert millicore to core
  cpuToCore: (millicore: number) => millicore / 1000,
  // Convert bytes to GB
  bytesToGB: (bytes: number) => bytes / (1024 * 1024 * 1024),
}

// Format number with appropriate precision
const formatNumber = {
  // Format value and remove trailing zeros
  formatDecimal: (value: number) => {
    const fixed = value.toFixed(1)
    return fixed.endsWith('.0') ? fixed.slice(0, -2) : fixed
  },
  // Format CPU value
  cpu: (value: number) => {
    if (value < 0.1) {
      return value.toFixed(1)
    }
    return formatNumber.formatDecimal(value)
  },
  // Format memory value
  memory: (value: number) => {
    if (value < 1) {
      return value.toFixed(1)
    }
    return formatNumber.formatDecimal(value)
  },
}

// Format value with unit
const formatValueWithUnit = (current: number, total: number, unit: string) => {
  const currentFormatted =
    unit === 'Core' ? formatNumber.cpu(current) : formatNumber.memory(current)
  const totalFormatted =
    unit === 'Core' ? formatNumber.cpu(total) : formatNumber.memory(total)
  return `${currentFormatted}/${totalFormatted} ${unit}`
}

interface ResourceTrendChartProps {
  current: number
  total: number
  unit: string
  data: { timestamp: number; value: number }[]
}

const mockData = [
  { timestamp: 1, value: 45 },
  { timestamp: 2, value: 65 },
  { timestamp: 3, value: 40 },
  { timestamp: 4, value: 70 },
  { timestamp: 5, value: 35 },
  { timestamp: 6, value: 55 },
  { timestamp: 7, value: 45 },
  { timestamp: 8, value: 60 },
]

const ResourceTrendChart = ({
  current,
  total,
  unit,
  data,
}: ResourceTrendChartProps) => {
  // Use mock data if no actual data is available
  const displayData = data?.length > 0 ? data : mockData
  const isLoading = data?.length === 0

  // Format display value
  const formatValue = (value: number) => {
    if (unit === 'Core') {
      return formatNumber.cpu(value)
    }
    return formatNumber.memory(value)
  }

  // Calculate normalized data for the chart
  const normalizedData = displayData.map(point => {
    const value = point.value
    if (data?.length > 0) {
      // For real data, use total as the base
      return (value / total) * 90 + 5
    } else {
      // For mock data, use the maximum value in the data as the base
      return (value / Math.max(...displayData.map(p => p.value))) * 85 + 5
    }
  })

  // Generate SVG path points
  const points = normalizedData.map((value, index) => {
    const x = (index / Math.max(normalizedData.length - 1, 1)) * 100
    const y = 100 - value
    return `${x},${y}`
  })

  // Generate smooth curve
  const path =
    points.length > 1
      ? points.reduce((acc, point, index) => {
          if (index === 0) {
            return `M ${point}`
          }
          const prevPoint = points[index - 1].split(',')
          const [x, y] = point.split(',')
          const cp1x =
            Number(prevPoint[0]) + (Number(x) - Number(prevPoint[0])) * 0.5
          const cp1y = Number(prevPoint[1])
          const cp2x = Number(x) - (Number(x) - Number(prevPoint[0])) * 0.5
          const cp2y = Number(y)
          return `${acc} C ${cp1x},${cp1y} ${cp2x},${cp2y} ${x},${y}`
        }, '')
      : ''

  const width = 100
  const height = 100

  return (
    <div className={styles.trend_chart}>
      <div className={styles.chart_container}>
        {displayData.length > 0 && (
          <svg
            className={styles.chart_svg}
            viewBox={`0 0 ${width} ${height}`}
            preserveAspectRatio="none"
          >
            <defs>
              <linearGradient
                id="areaGradient"
                x1="0%"
                y1="0%"
                x2="0%"
                y2="100%"
              >
                <stop offset="0%" stopColor="#597ef7" stopOpacity="0.3" />
                <stop offset="100%" stopColor="#597ef7" stopOpacity="0.05" />
              </linearGradient>
            </defs>
            <path
              className={styles.area}
              d={`${path} L ${width},${height} L 0,${height} Z`}
            />
            <path className={styles.line} d={path} />
            {normalizedData.map((value, index) => (
              <g key={index} className={styles.point_group}>
                <circle
                  cx={`${(index / Math.max(normalizedData.length - 1, 1)) * 100}%`}
                  cy={`${100 - value}%`}
                  className={styles.point_halo}
                />
                <circle
                  cx={`${(index / Math.max(normalizedData.length - 1, 1)) * 100}%`}
                  cy={`${100 - value}%`}
                  className={styles.point}
                />
              </g>
            ))}
          </svg>
        )}
        <div className={styles.value_text}>
          {isLoading ? (
            <span className={styles.loading}>Loading...</span>
          ) : (
            <>
              <span className={styles.current}>{formatValue(current)}</span>
              <span className={styles.unit}>{unit}</span>
              <span className={styles.hover_text}>
                /{formatValue(total)}
                {unit}
              </span>
            </>
          )}
        </div>
      </div>
    </div>
  )
}

const copyToClipboard = (text: string, t: any) => {
  navigator.clipboard.writeText(text).then(() => {
    message.success(t('CopiedToClipboard'))
  })
}

const PopoverCard = ({ data }: any) => {
  const { t } = useTranslation()

  if (!data) {
    return <div className={styles.value}>--</div>
  }
  const type = getDataType(data)
  const displayText = type === 'String' ? data : JSON.stringify(data)

  const handleClick = () => {
    if (displayText !== '--') {
      copyToClipboard(displayText, t)
    }
  }

  return (
    <div className={styles.value} onClick={handleClick}>
      {displayText}
    </div>
  )
}

interface MetricTooltipProps {
  value: string | number
  tooltipKey: string
}

const MetricTooltip = ({ value, tooltipKey }: MetricTooltipProps) => {
  const { t } = useTranslation()
  return (
    <Tooltip title={t(`Metrics.Tooltips.${tooltipKey}`)}>
      <span>{value}</span>
    </Tooltip>
  )
}

type SummaryCardProps = {
  auditStat: { score: number }
  summary: any
}

const SummaryCard = ({ auditStat, summary }: SummaryCardProps) => {
  const [score, setScore] = useState(0)
  const { t } = useTranslation()
  const location = useLocation()
  const scoreRef = useRef(0)
  const { type } = queryString.parse(location?.search)

  useEffect(() => {
    let timer
    if (auditStat?.score) {
      timer = setInterval(() => {
        scoreRef.current = scoreRef?.current + 10
        if (scoreRef.current > auditStat?.score) {
          clearInterval(timer)
          setScore(auditStat?.score)
        } else {
          setScore(scoreRef?.current)
        }
      }, 100)
    }
    return () => {
      clearInterval(timer)
    }
  }, [auditStat?.score])

  // Convert metrics data for charts
  const getMetricsData = (
    metrics: ResourceMetrics | undefined,
    converter: (value: number) => number,
  ) => {
    if (!metrics?.points?.length) return []

    return metrics.points
      .sort(
        (a, b) => dayjs(a.timestamp).valueOf() - dayjs(b.timestamp).valueOf(),
      )
      .map(point => ({
        timestamp: dayjs(point.timestamp).valueOf(),
        value: converter(point.value),
      }))
  }

  const cpuData = getMetricsData(summary?.cpuMetrics, convertUnits.cpuToCore)
  const memoryData = getMetricsData(
    summary?.memoryMetrics,
    convertUnits.bytesToGB,
  )

  return (
    <div className={styles.summary_card}>
      <Progress
        size={80}
        type="circle"
        percent={parseInt(`${score}`)}
        format={percent => `${percent}`}
        strokeWidth={12}
      />
      <div className={styles.summary}>
        {type === 'resource' && (
          <div className={styles.field_box}>
            <div className={styles.item}>
              <div className={styles.label}>APIVersion </div>
              <PopoverCard data={summary?.resource?.apiVersion} />
            </div>
            <div className={styles.item}>
              <div className={styles.label}>Cluster </div>
              <PopoverCard data={summary?.resource?.cluster} />
            </div>
            <div className={styles.item}>
              <div className={styles.label}>Kind </div>
              <PopoverCard data={summary?.resource?.kind} />
            </div>
            <div className={styles.item}>
              <div className={styles.label}>Namespace </div>
              <PopoverCard data={summary?.resource?.namespace} />
            </div>
            <div className={styles.item}>
              <div className={styles.label}>Name </div>
              <PopoverCard data={summary?.resource?.name} />
            </div>
          </div>
        )}
        {type === 'kind' && (
          <div className={styles.field_box}>
            <div className={styles.item}>
              <div className={styles.label}>Cluster </div>
              <PopoverCard data={summary?.cluster} />
            </div>
            <div className={styles.item}>
              <div className={styles.label}>Kind </div>
              <PopoverCard data={summary?.kind} />
            </div>
            <div className={styles.item}>
              <div className={styles.label}>Group </div>
              <div className={styles.value}>{summary?.group || '--'}</div>
            </div>
            <div className={styles.item}>
              <div className={styles.label}>Count </div>
              <div className={styles.value}>{summary?.count || '--'}</div>
            </div>
            <div className={styles.item}>
              <div className={styles.label}>Version </div>
              <PopoverCard data={summary?.version} />
            </div>
          </div>
        )}
        {type === 'namespace' && (
          <div className={styles.field_box}>
            <div className={styles.item}>
              <div className={styles.label}>Cluster </div>
              <PopoverCard data={summary?.cluster} />
            </div>
          </div>
        )}
        {type === 'cluster' && (
          <div className={styles.field_box}>
            {/* Node metrics */}
            <div className={styles.item}>
              <div className={styles.label}>Node Count</div>
              <div className={styles.value}>
                <MetricTooltip
                  value={summary?.readyNodes || '--'}
                  tooltipKey="ReadyNodes"
                />
                /
                <MetricTooltip
                  value={summary?.nodeCount || '--'}
                  tooltipKey="NodeCount"
                />
              </div>
            </div>

            {/* Resource usage metrics */}
            <div
              className={`${styles.item} ${summary?.metricsEnabled === false ? '' : styles.with_trend}`}
            >
              <div className={styles.label}>CPU</div>
              {summary?.metricsEnabled === false ? (
                <div className={styles.value}>
                  <MetricTooltip
                    value={formatValueWithUnit(
                      convertUnits.cpuToCore(summary?.cpuUsage || 0),
                      convertUnits.cpuToCore(summary?.cpuCapacity || 100000),
                      'Core',
                    )}
                    tooltipKey="CPU"
                  />
                </div>
              ) : (
                <div className={styles.trend}>
                  <ResourceTrendChart
                    current={convertUnits.cpuToCore(summary?.cpuUsage || 0)}
                    total={convertUnits.cpuToCore(
                      summary?.cpuCapacity || 100000,
                    )}
                    unit="Core"
                    data={cpuData}
                  />
                </div>
              )}
            </div>
            <div
              className={`${styles.item} ${summary?.metricsEnabled === false ? '' : styles.with_trend}`}
            >
              <div className={styles.label}>Memory</div>
              {summary?.metricsEnabled === false ? (
                <div className={styles.value}>
                  <MetricTooltip
                    value={formatValueWithUnit(
                      convertUnits.bytesToGB(summary?.memoryUsage || 0),
                      convertUnits.bytesToGB(
                        summary?.memoryCapacity || 8 * 1024 * 1024 * 1024,
                      ),
                      'GB',
                    )}
                    tooltipKey="Memory"
                  />
                </div>
              ) : (
                <div className={styles.trend}>
                  <ResourceTrendChart
                    current={convertUnits.bytesToGB(summary?.memoryUsage || 0)}
                    total={convertUnits.bytesToGB(
                      summary?.memoryCapacity || 8 * 1024 * 1024 * 1024,
                    )}
                    unit="GB"
                    data={memoryData}
                  />
                </div>
              )}
            </div>
            <div className={styles.item}>
              <div className={styles.label}>Pods</div>
              <div className={styles.value}>
                <MetricTooltip
                  value={`${summary?.podsUsage || 0}/${summary?.podsCapacity || 0}`}
                  tooltipKey="Pods"
                />
              </div>
            </div>

            {/* Cluster information */}
            <div className={styles.item}>
              <div className={styles.label}>Kubernetes {t('Version')} </div>
              <div className={styles.value}>
                <MetricTooltip
                  value={summary?.serverVersion || '--'}
                  tooltipKey="ServerVersion"
                />
              </div>
            </div>

            <div className={styles.item}>
              <div className={styles.label}>Metrics Server</div>
              <div className={styles.value}>
                <MetricTooltip
                  value={
                    summary?.metricsEnabled === undefined
                      ? '--'
                      : summary?.metricsEnabled
                        ? t('Enabled')
                        : t('Disabled')
                  }
                  tooltipKey="MetricsServer"
                />
              </div>
            </div>
          </div>
        )}
        <div className={styles.field_box}>
          {summary?.namespace && (
            <div className={styles.item}>
              <div className={styles.label}>Namespace </div>
              <PopoverCard data={summary?.namespace} />
            </div>
          )}
          {summary?.labels && (
            <>
              {Object?.entries(summary?.labels)?.map(([k, v]) => (
                <div key={`${k}_${v}`} className={styles.item}>
                  <div className={styles.label}>{k}</div>
                  <PopoverCard data={v} />
                </div>
              ))}
            </>
          )}
          {summary?.annotations && (
            <>
              {Object?.entries(summary?.annotations)?.map(([k, v]) => (
                <div key={`${k}_${v}`} className={styles.item}>
                  <div className={styles.label}>{k}</div>
                  <PopoverCard data={v} />
                </div>
              ))}
            </>
          )}
          {summary?.countByGVK ? renderStatistics(summary?.countByGVK) : null}
        </div>
      </div>
    </div>
  )
}

function renderStatistics(GVKData: any) {
  const sortList = Object.entries(GVKData)
    ?.map(([key, value]) => ({ key, value }))
    ?.sort((a: any, b: any) => b?.value - a?.value)
  return (
    <>
      {sortList?.map(({ key, value }: any) => {
        const [displayKey] = key?.split('.')
        return (
          <div className={styles.item} key={key}>
            <div className={styles.label}>{displayKey}</div>
            <div className={styles.value}>{value}</div>
          </div>
        )
      })}
    </>
  )
}

export default SummaryCard
