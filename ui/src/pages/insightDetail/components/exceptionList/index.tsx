import React, { useEffect, useState } from 'react'
import { Button, Empty, Tag } from 'antd'
import { useTranslation } from 'react-i18next'
import { ArrowRightOutlined, UpOutlined, DownOutlined } from '@ant-design/icons'
import Loading from '@/components/loading'
import { SEVERITY_MAP } from '@/utils/constants'
import ExceptionStat from '../exceptionStat'

import styles from './style.module.less'

type IProps = {
  exceptionList: any
  rescan: () => void
  showDrawer: () => void
  onItemClick: (val: string) => void
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
  const { t } = useTranslation()
  const [isShowList, setIsShowList] = useState(true)

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
  const iconStyle = { marginLeft: 5, color: '#646566' }
  return (
    <div className={styles.exception}>
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
          <Button onClick={rescan}>{t('Rescan')}</Button>
          <div className={styles.header_right_action} onClick={handleClick}>
            {isShowList ? (
              <span>
                {t('Collapse')}
                <UpOutlined style={iconStyle} />
              </span>
            ) : (
              <span>
                {t('Expand')}
                <DownOutlined style={iconStyle} />
              </span>
            )}
          </div>
        </div>
      </div>
      {isShowList && (
        <div className={styles.body}>
          {auditLoading ? (
            <div className={styles.loading_box}>
              <Loading />
            </div>
          ) : top5List && top5List?.length > 0 ? (
            <>
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
                        {t('ViewIssueDetail')}
                      </div>
                    )}

                    <div className={styles.itme_content}>
                      <div className={styles.left}>
                        <Tag
                          color={SEVERITY_MAP?.[item?.issue?.severity]?.color}
                        >
                          {t(SEVERITY_MAP?.[item?.issue?.severity]?.text)}
                        </Tag>
                      </div>
                      <div className={styles.right}>
                        <div className={styles.right_top}>
                          <div>
                            <span className={styles.title}>
                              {item?.issue?.title}
                            </span>
                            <span>
                              {t('Occur')}&nbsp;
                              <span
                                style={{ fontWeight: 'bold', color: '#646566' }}
                              >
                                {item?.resourceGroups?.length}
                              </span>
                              &nbsp;{t('Times')}
                            </span>
                            <span style={{ color: '#000' }}>&nbsp;|</span>
                            &nbsp;{t('CollectedFrom')}
                          </div>
                          <div className={styles.tool}>
                            <ArrowRightOutlined />
                            &nbsp;{item?.issue?.scanner}&nbsp;{t('Tool')}
                          </div>
                        </div>
                        <div className={styles.right_bottom}>
                          <div className={styles.label}>
                            {t('Description')}:&nbsp;
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
                  {t('CheckAllIssues')}
                  <ArrowRightOutlined />
                </span>
              </div>
            </>
          ) : (
            <div className={styles.content_empty}>
              <Empty description={`${t('NoIssuesFound')}`} />
            </div>
          )}
        </div>
      )}
    </div>
  )
}

export default ExceptionList
