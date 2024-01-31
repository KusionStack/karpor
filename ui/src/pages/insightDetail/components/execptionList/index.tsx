import { useEffect, useState } from 'react'
import { Button, Empty, Tag } from 'antd'
import styles from './style.module.less'
import ExecptionStat from '../execptionStat'
import { ArrowRightOutlined } from '@ant-design/icons'
import Loading from '../../../../components/loading'
import { SEVERITY_MAP } from '../../../../utils/constants'
import React from 'react'

type IProps = {
  execptionList: any
  rescan: () => void
  showDrawer: () => void
  onItemClick: (val: string) => void
  execptionStat: any
  auditLoading: boolean
}

const ExecptionList = ({
  execptionList,
  rescan,
  showDrawer,
  onItemClick,
  execptionStat,
  auditLoading,
}: IProps) => {
  const [selectedEventId, setSelectedEventId] = useState<any>()
  const [top5List, setTop5list] = useState([])
  const [currentKey, setCurrentKey] = useState('All')

  useEffect(() => {
    if (currentKey === 'All') {
      const defaultTop5List = execptionList?.issueGroups?.slice(0, 5)
      setTop5list(defaultTop5List)
    } else {
      const tmp = execptionList?.issueGroups?.filter(
        (item: any) => item?.issue?.severity === currentKey,
      )
      const top5Tmp = tmp?.slice(0, 5)
      setTop5list(top5Tmp)
    }
  }, [currentKey, execptionList?.issueGroups])

  function onClickTable(key) {
    setCurrentKey(key)
  }
  return (
    <div className={styles.exception}>
      <div className={styles.header}>
        <div className={styles.header_left}>
          <ExecptionStat
            currentKey={currentKey}
            statData={{
              all: execptionStat?.issuesTotal || 0,
              high: execptionStat?.severityStatistic?.High || 0,
              medium: execptionStat?.severityStatistic?.Medium || 0,
              low: execptionStat?.severityStatistic?.Low || 0,
            }}
            onClickTable={onClickTable}
          />
        </div>
        <div className={styles.header_right}>
          <Button onClick={rescan}>重新扫描</Button>
        </div>
      </div>
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
                    <div className={styles.itme_tip}>查看事件详情</div>
                  )}

                  <div className={styles.top}>
                    <div className={styles.left}>
                      <Tag color={SEVERITY_MAP?.[item?.issue?.severity]?.color}>
                        {SEVERITY_MAP?.[item?.issue?.severity]?.text}
                      </Tag>
                      <span className={styles.title}>{item?.issue?.title}</span>
                    </div>
                    <div className={styles.right}>
                      <div>
                        <span>
                          发生&nbsp;
                          <span
                            style={{ fontWeight: 'bold', color: '#646566' }}
                          >
                            {item?.locators?.length}
                          </span>
                          &nbsp;次
                        </span>
                        &nbsp;采集自
                      </div>
                      <div className={styles.tool}>
                        <ArrowRightOutlined />
                        &nbsp;{item?.issue?.scanner}&nbsp;工具
                      </div>
                    </div>
                  </div>
                  <div className={styles.bottom}>
                    <div className={styles.label}>message: </div>
                    <div className={styles.value}>{item?.issue?.message}</div>
                  </div>
                </div>
              )
            })}
            <div className={styles.footer}>
              <span className={styles.btn} onClick={showDrawer}>
                查看全部异常事件
                <ArrowRightOutlined />
              </span>
            </div>
          </>
        ) : (
          <div
            style={{
              height: 372,
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center',
            }}
          >
            <Empty description="无风险异常，再接再厉~" />
          </div>
        )}
      </div>
    </div>
  )
}

export default ExecptionList
