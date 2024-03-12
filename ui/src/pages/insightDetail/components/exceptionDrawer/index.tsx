import React, { useEffect, useState } from 'react'
import { Collapse, Drawer, Empty, Input, Pagination, Tag } from 'antd'
import ExceptionStat from '../exceptionStat'
import { CaretRightOutlined, SearchOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { truncationPageData } from '../../../../utils/tools'

import MutiTag from '../mutiTag'
import { SEVERITY_MAP } from '../../../../utils/constants'

import styles from './style.module.less'

const DEFALUT_PAGE_SIZE = 10

type IProps = {
  exceptionStat?: any
  exceptionList: any
  onClose: () => void
  open: boolean
}

const ExceptionDrawer = ({ open, onClose, exceptionList }: IProps) => {
  const [pageParams, setPageParams] = useState({
    pageNo: 1,
    pageSize: DEFALUT_PAGE_SIZE,
    total: 0,
  })
  const [searchValue, setSearchValue] = useState('')
  const [currentKey, setCurrentKey] = useState('All')
  const [showPageData, setShowPageData] = useState([])

  const { t } = useTranslation()

  useEffect(() => {
    if (currentKey === 'All') {
      let tmp: any = []
      if (!searchValue) {
        tmp = exceptionList?.issueGroups
      } else {
        const newValue = searchValue?.toLowerCase().trim()?.split(' ')
        const issueGroups = exceptionList?.issueGroups
        if (newValue?.length === 1) {
          issueGroups?.forEach((item: any) => {
            if (item?.issue?.title?.toLowerCase()?.includes(newValue?.[0])) {
              tmp.push(item)
            }
          })
        } else {
          issueGroups?.forEach((item: any) => {
            if (
              newValue?.every((innerValue: string) =>
                item?.issue?.title?.toLowerCase()?.includes(innerValue),
              )
            ) {
              tmp.push(item)
            }
          })
        }
      }
      const pageList = truncationPageData({
        list: tmp,
        page: pageParams?.pageNo,
        pageSize: pageParams?.pageSize,
      })
      setShowPageData(pageList)
      setPageParams({
        ...pageParams,
        total: tmp?.length,
      })
    } else {
      const tmp = exceptionList?.issueGroups?.filter(
        (item: any) => item?.issue?.severity === currentKey,
      )
      let filterTmp: any = []
      if (!searchValue) {
        filterTmp = tmp
      } else {
        const newValue = searchValue?.toLowerCase().trim()?.split(' ')
        const issueGroups = exceptionList?.issueGroups
        if (newValue?.length === 1) {
          issueGroups?.forEach((item: any) => {
            if (item?.issue?.title?.toLowerCase()?.includes(newValue?.[0])) {
              filterTmp.push(item)
            }
          })
        } else {
          issueGroups?.forEach((item: any) => {
            if (
              newValue?.every((innerValue: string) =>
                item?.issue?.title?.toLowerCase()?.includes(innerValue),
              )
            ) {
              filterTmp.push(item)
            }
          })
        }
      }
      const pageList = truncationPageData({
        list: filterTmp,
        page: pageParams?.pageNo,
        pageSize: pageParams?.pageSize,
      })
      setPageParams({
        ...pageParams,
        total: filterTmp?.length,
      })
      setShowPageData(pageList)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    currentKey,
    exceptionList?.issueGroups,
    pageParams?.pageNo,
    pageParams?.pageSize,
    searchValue,
  ])

  // useEffect(() => {
  //   let tmp = [];
  //   if (!searchValue) {
  //     tmp = exceptionList?.issueGroups
  //   } else {
  //     exceptionList?.issueGroups?.forEach(item => {
  //       if (item?.issue?.title?.includes(searchValue)) {
  //         tmp.push(item)
  //       }
  //     })
  //   }
  //   const pageList = truncationPageData({ list: tmp, page: pageParams?.pageNo, pageSize: pageParams?.pageSize })
  //   setShowPageData(pageList)
  // }, [exceptionList?.issueGroups, pageParams?.pageNo, pageParams?.pageSize, searchValue]);

  function onSearch(event) {
    const val = event.target.value
    setSearchValue(val)
    setPageParams({
      ...pageParams,
      pageNo: 1,
    })
  }

  function handleChangePage(page, pageSize) {
    setPageParams({
      ...pageParams,
      pageNo: page,
      pageSize,
    })
  }

  const panelStyle: React.CSSProperties = {
    background: '#fff',
    borderRadius: 8,
    border: '1px solid rgba(0,0,0,0.15)',
    marginBottom: 8,
  }

  function getItems() {
    return showPageData?.map(item => {
      const uniqueKey = `${item?.issue?.title}_${item?.issue?.message}_${item?.issue?.scanner}_${item?.issue?.severity}`
      const locatorsNames = item?.locators?.map(item => {
        return {
          ...item,
          allName: `${item?.cluster || ''} ${item?.apiVersion || ''} ${item?.kind || ''} ${item?.namespace || ''} ${item?.name || ''} `,
        }
      })
      return {
        key: uniqueKey,
        label: (
          <div className={styles.collapse_panel_title}>
            <div>
              <Tag color={SEVERITY_MAP?.[item?.issue?.severity]?.color}>
                {t(SEVERITY_MAP?.[item?.issue?.severity]?.text)}
              </Tag>
            </div>
            <div className={styles.right}>
              <div className={styles.title}>{item?.issue?.title}</div>
              <div className={styles.right_bottom}>
                <div className={styles.label}>{t('Description')}: </div>
                <div className={styles.value}>{item?.issue?.message}</div>
              </div>
            </div>
          </div>
        ),
        children: (
          <div className={styles.collapse_panel_body}>
            <div className={styles.header}>
              <div className={styles.item}>
                <div className={styles.label}>{t('IssueSource')}:&nbsp;</div>
                <div className={styles.value}>
                  {item?.issue?.scanner || '--'}
                </div>
              </div>
              <div className={styles.item}>
                <div className={styles.label}>
                  {t('NumberOfOccurrences')}:&nbsp;
                </div>
                <div className={styles.value}>{item?.locators?.length}</div>
              </div>
            </div>
            <div className={`${styles.row_item}`}>
              <div className={styles.label}>{t('Description')}:&nbsp;</div>
              <div className={styles.value}>{item?.issue?.message || '--'}</div>
            </div>
            <div className={styles.body}>
              <div className={styles.label}>{t('RelatedResources')}: </div>
              <div className={styles.value}>
                <MutiTag allTags={locatorsNames} />
              </div>
            </div>
          </div>
        ),
        style: panelStyle,
      }
    })
  }

  function onClickTable(key) {
    setCurrentKey(key)
    setSearchValue('')
    setPageParams({
      ...pageParams,
      pageNo: 1,
    })
  }

  return (
    <Drawer
      width={1000}
      title={t('IssuesDetail')}
      open={open}
      onClose={onClose}
    >
      <div className={styles.exception_drawer}>
        <ExceptionStat
          currentKey={currentKey}
          statData={{
            all: exceptionList?.issueTotal,
            high: exceptionList?.bySeverity?.High,
            medium: exceptionList?.bySeverity?.Medium,
            low: exceptionList?.bySeverity?.Low,
          }}
          onClickTable={onClickTable}
        />
        <div className={styles.tool_bar}>
          <div className={styles.search}>
            <Input
              placeholder={t('FilterByName')}
              suffix={
                <SearchOutlined
                  className="site-form-item-icon"
                  style={{ color: '#999' }}
                />
              }
              allowClear
              style={{ width: 200 }}
              onChange={onSearch}
            />
          </div>
        </div>
        {showPageData && showPageData?.length > 0 ? (
          <>
            <div className={styles.events}>
              <Collapse
                bordered={false}
                expandIcon={({ isActive }) => (
                  <CaretRightOutlined rotate={isActive ? 90 : 0} />
                )}
                style={{ background: '#fff' }}
                items={getItems()}
              />
            </div>
            <div style={{ textAlign: 'right', marginTop: 16 }}>
              <Pagination
                total={pageParams?.total}
                showTotal={(total, range) =>
                  `${range[0]}-${range[1]} ${t('Total')} ${total}`
                }
                pageSize={pageParams?.pageSize}
                current={pageParams?.pageNo}
                onChange={handleChangePage}
                showSizeChanger
              />
            </div>
          </>
        ) : (
          <div
            style={{
              height: 400,
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center',
            }}
          >
            <Empty />
          </div>
        )}
      </div>
    </Drawer>
  )
}

export default ExceptionDrawer
