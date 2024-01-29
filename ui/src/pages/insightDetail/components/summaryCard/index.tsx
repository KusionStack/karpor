import React from 'react'
import { Popover } from 'antd'
import { useLocation } from 'react-router-dom'
import queryString from 'query-string'
import { filterSize, getDataType } from '@/utils/tools'
import GaugeChart from '../gauge'

import styles from './styles.module.less'

function getContent(data) {
  return (
    <>
      {Object.keys(data)?.map(key => {
        const [displayKey] = key?.split('.')
        return (
          <div className={styles.popoverItem} key={key}>
            <span className={styles.popoverLabel}>{displayKey}:</span>
            <span className={styles.popoverValue}>{data?.[key]}</span>
          </div>
        )
      })}
    </>
  )
}

const PopoverCard = ({ data }: any) => {
  if (!data) {
    return <div className={styles.value}>--</div>
  }
  const type = getDataType(data)
  const showContent = type === 'String' ? <span>{data}</span> : getContent(data)
  return (
    <Popover content={showContent}>
      <div className={styles.value}>
        {type === 'String' ? data : JSON.stringify(data)}
      </div>
    </Popover>
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
            <span className={styles.label}>{displayKey}:</span>
            <span className={styles.value}>{value}</span>
          </div>
        )
      })}
    </>
  )
}
type SummaryCardProps = {
  auditStat: { score: number }
  summary: any
}
const SummaryCard = ({ auditStat, summary }: SummaryCardProps) => {
  const location = useLocation()
  const { type } = queryString.parse(location?.search)
  return (
    <div className={styles.summary_card}>
      <GaugeChart data={(auditStat?.score / 100)?.toFixed(2)} />
      <div className={styles.summary}>
        {type === 'resource' && (
          <div className={styles.field_box}>
            <div className={styles.item}>
              <div className={styles.label}>apiVersion：</div>
              <PopoverCard data={summary?.resource?.apiVersion} />
              {/* <div className={styles.value}>{summary?.resource?.apiVersion || '--'}</div> */}
            </div>
            <div className={styles.item}>
              <div className={styles.label}>cluster：</div>
              {/* <div className={styles.value}>{summary?.resource?.cluster || '--'}</div> */}
              <PopoverCard data={summary?.resource?.cluster} />
            </div>
            <div className={styles.item}>
              <div className={styles.label}>kind：</div>
              {/* <div className={styles.value}>{summary?.resource?.kind || '--'}</div> */}
              <PopoverCard data={summary?.resource?.kind} />
            </div>
            <div className={styles.item}>
              <div className={styles.label}>namespace：</div>
              {/* <div className={styles.value}>{summary?.resource?.namespace || '--'}</div> */}
              <PopoverCard data={summary?.resource?.namespace} />
            </div>
            <div className={styles.item}>
              <div className={styles.label}>name：</div>
              {/* <div className={styles.value}>{summary?.resource?.name || '--'}</div> */}
              <PopoverCard data={summary?.resource?.name} />
            </div>
          </div>
        )}
        {type === 'kind' && (
          <div className={styles.field_box}>
            <div className={styles.item}>
              <div className={styles.label}>cluster：</div>
              {/* <div className={styles.value}>{summary?.cluster || '--'}</div> */}
              <PopoverCard data={summary?.cluster} />
            </div>
            <div className={styles.item}>
              <div className={styles.label}>kind：</div>
              {/* <div className={styles.value}>{summary?.kind || '--'}</div> */}
              <PopoverCard data={summary?.kind} />
            </div>
            <div className={styles.item}>
              <div className={styles.label}>group：</div>
              <div className={styles.value}>{summary?.group || '--'}</div>
            </div>
            <div className={styles.item}>
              <div className={styles.label}>count：</div>
              <div className={styles.value}>{summary?.count || '--'}</div>
            </div>
            <div className={styles.item}>
              <div className={styles.label}>version：</div>
              {/* <div className={styles.value}>{summary?.version || '--'}</div> */}
              <PopoverCard data={summary?.version} />
            </div>
          </div>
        )}
        {type === 'namespace' && (
          <div className={styles.field_box}>
            <div className={styles.item}>
              <div className={styles.label}>cluster：</div>
              {/* <div className={styles.value}>{summary?.cluster || '--'}</div> */}
              <PopoverCard data={summary?.cluster} />
            </div>
            <div className={styles.item}>
              <div className={styles.label}>namespace：</div>
              {/* <div className={styles.value}>{summary?.namespace || '--'}</div> */}
              <PopoverCard data={summary?.namespace} />
            </div>
            <div className={`${styles.item} ${styles.namespaceStat}`}>
              <div className={styles.label}>Top Resources</div>
            </div>
            {summary?.countByGVK ? renderStatistics(summary?.countByGVK) : null}
          </div>
        )}
        {type === 'cluster' && (
          <div className={styles.field_box}>
            <div className={styles.item}>
              <div className={styles.label}>Node：</div>
              <div className={styles.value}>{summary?.nodeCount || '--'}</div>
            </div>
            <div className={styles.item}>
              <div className={styles.label}>Pod：</div>
              <div className={styles.value}>
                {summary?.podsCapacity || '--'}
              </div>
            </div>
            <div className={styles.item}>
              <div className={styles.label}>CPU 容量：</div>
              <div className={styles.value}>{summary?.cpuCapacity || '--'}</div>
            </div>
            <div className={styles.item}>
              <div className={styles.label}>内存容量：</div>
              <div className={styles.value}>
                {summary?.memoryCapacity
                  ? filterSize(summary?.memoryCapacity)
                  : '--'}
              </div>
            </div>
            <div className={styles.item}>
              <div className={styles.label}>Kubernetes 版本：</div>
              {/* <div className={styles.value}>{summary?.serverVersion || '--'}</div> */}
              <PopoverCard data={summary?.serverVersion} />
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export default SummaryCard
