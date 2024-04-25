import React from 'react'
import { Popover, Progress } from 'antd'
import { useLocation } from 'react-router-dom'
import queryString from 'query-string'
import { useTranslation } from 'react-i18next'
import { filterSize, getDataType } from '@/utils/tools'

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
            <div className={styles.label}>{displayKey}</div>
            <div className={styles.value}>{value}</div>
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
  const { t } = useTranslation()
  const location = useLocation()
  const { type } = queryString.parse(location?.search)
  return (
    <div className={styles.summary_card}>
      <Progress
        size={80}
        type="circle"
        percent={parseInt(`${auditStat?.score}`)}
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
            {/* <div className={`${styles.item} ${styles.namespaceStat}`}>
              <div className={styles.label}>Top Resources</div>
            </div> */}
          </div>
        )}
        {type === 'cluster' && (
          <div className={styles.field_box}>
            <div className={styles.item}>
              <div className={styles.label}>Node </div>
              <div className={styles.value}>{summary?.nodeCount || '--'}</div>
            </div>
            <div className={styles.item}>
              <div className={styles.label}>Pod </div>
              <div className={styles.value}>
                {summary?.podsCapacity || '--'}
              </div>
            </div>
            <div className={styles.item}>
              <div className={styles.label}>CPU {t('Capacity')} </div>
              <div className={styles.value}>{summary?.cpuCapacity || '--'}</div>
            </div>
            <div className={styles.item}>
              <div className={styles.label}>{t('MemoryCapacity')} </div>
              <div className={styles.value}>
                {summary?.memoryCapacity
                  ? filterSize(summary?.memoryCapacity)
                  : '--'}
              </div>
            </div>
            <div className={styles.item}>
              <div className={styles.label}>Kubernetes {t('Version')} </div>
              <PopoverCard data={summary?.serverVersion} />
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
          {summary?.countByGVK ? renderStatistics(summary?.countByGVK) : null}
        </div>
      </div>
    </div>
  )
}

export default SummaryCard
