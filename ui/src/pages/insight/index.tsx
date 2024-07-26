import React, { useEffect, useState } from 'react'
import { Button, Empty, Input, Pagination, message } from 'antd'
import {
  EditOutlined,
  SearchOutlined,
  SortAscendingOutlined,
  SortDescendingOutlined,
} from '@ant-design/icons'
import { useSelector } from 'react-redux'
import axios from 'axios'
import { useTranslation } from 'react-i18next'
import queryString from 'query-string'
import { useLocation, useNavigate } from 'react-router-dom'
import Loading from '@/components/loading'
import {
  filterKeywordsOfArray,
  hasDuplicatesOfObjectArray,
  isEmptyObject,
  truncationPageData,
} from '@/utils/tools'
import CardContent from './components/card'
import RuleForm from './components/ruleForm'
import QuotaContent from './components/quotaContent'
import { InsightTabs } from './components/insightTabs'

import styles from './styles.module.less'

const orderIconStyle: React.CSSProperties = {
  marginLeft: 0,
}

const defalutPageParams = {
  pageNo: 1,
  pageSize: 12,
  total: 0,
}

const Insight = () => {
  const { t } = useTranslation()
  const location = useLocation()
  const { isReadOnlyMode } = useSelector((state: any) => state.globalSlice)
  const urlParams = queryString.parse(location?.search)
  const navigate = useNavigate()
  const [tabList, setTabList] = useState([])
  const [activeTabKey, setActiveTabKey] = useState<any>(
    urlParams?.activeTabKey || 'namespace',
  )
  const [open, setOpen] = useState(false)
  const [isEdit, setIsEdit] = useState(false)
  const [currentItem, setCurrentItem] = useState<any>()
  const [pageParams, setPageParams] = useState(defalutPageParams)
  const [sortParams, setSortParams] = useState<any>({
    orderBy: 'name',
    isAsc: true,
  })
  const [allResourcesData, setAllResourcesData] = useState<any>([])
  const [showData, setShowData] = useState<any>([])
  const [statsData, setStatsData] = useState<any>()
  const [keyword, setKeyword] = useState('')
  const [resouresLoading, setResouresLoading] = useState<boolean>(false)

  async function queryStats() {
    const response: any = await axios.get('/rest-api/v1/insight/stats')
    if (response) {
      if (response?.success) {
        if (response?.data?.clusterCount <= 0) {
          message.info(t('NoClusterAndJumpToClusterPage'))
          setTimeout(() => {
            navigate('/cluster')
          }, 2000)
          return
        }
        setStatsData(response?.data)
      } else {
        message.error(response?.message || t('RequestFailedAndTry'))
      }
    }
  }

  async function queryAllRules(isDelete) {
    const response: any = await axios.get('/rest-api/v1/resource-group-rules')
    if (response?.success) {
      const tabList = response?.data
        ?.filter(item => item)
        ?.map(item => ({
          ...item,
          key: item?.name,
          label: item?.name,
        }))
      setTabList(tabList)
      handleClickItem(
        isDelete || !activeTabKey ? tabList?.[0]?.name : activeTabKey,
      )
    } else {
      message.error(response?.message || t('RequestFailedAndTry'))
    }
  }

  async function queryCurrentResources(ruleName) {
    setResouresLoading(true)
    setAllResourcesData({})
    const response: any = await axios.get(
      `/rest-api/v1/resource-groups/${ruleName}`,
    )
    if (response) {
      if (response?.success) {
        const { groups, fields } = response?.data || {}
        const newGroups = groups
          ?.filter(item => item && !isEmptyObject(item))
          ?.map(group => {
            const { title, tags } = getName(group, fields)
            return {
              ...group,
              title,
              tags,
            }
          })
        const newData = {
          fields: response?.data?.fields,
          groups: newGroups,
        }
        setAllResourcesData(newData)
      } else {
        message.error(response?.message || t('RequestFailedAndTry'))
      }
    }
    setResouresLoading(false)
  }

  async function deleteItem(itemKey, callback) {
    const response: any = await axios.delete(
      `/rest-api/v1/resource-group-rule/${itemKey}`,
    )
    if (response?.success) {
      callback && callback()
      setIsEdit(false)
      queryAllRules(true)
      queryStats()
      queryCurrentResources(tabList?.[0]?.name)
      setOpen(false)
    } else {
      message.error(response?.message || t('RequestFailedAndTry'))
    }
  }

  async function handleSubmit(values, callback) {
    const isDuplicates = hasDuplicatesOfObjectArray(values?.fields)
    if (isDuplicates) {
      message.error(t('DuplicateData'))
      return
    }
    const response: any = await axios(`/rest-api/v1/resource-group-rule`, {
      method: isEdit ? 'PUT' : 'POST',
      data: {
        name: values?.name,
        description: values?.description,
        fields: values?.fields?.map(item => {
          return item?.key === 'namespace'
            ? 'namespace'
            : item?.value
              ? `${item?.key}.${item?.value}`
              : item?.key
        }),
      },
    })
    if (response?.success) {
      callback && callback()
      setOpen(false)
      setIsEdit(false)
      queryAllRules(false)
    } else {
      message.error(response?.message || t('RequestFailedAndTry'))
    }
  }

  useEffect(() => {
    let tmp = allResourcesData?.groups
    if (keyword) {
      const keywords = keyword?.toLowerCase()?.trim()?.split(' ')
      tmp = filterKeywordsOfArray(allResourcesData?.groups, keywords, 'title')
    }
    const pageList = truncationPageData({
      list: tmp,
      page: pageParams?.pageNo,
      pageSize: pageParams?.pageSize,
    })
    setShowData({
      fields: allResourcesData?.fields,
      groups: pageList,
    })
    setPageParams({
      ...pageParams,
      total: tmp?.length,
    })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [allResourcesData, pageParams?.pageNo, pageParams?.pageSize, keyword])

  function handleChangePage(page, pageSize) {
    setPageParams({
      ...pageParams,
      pageNo: page,
      pageSize,
    })
  }

  function handleClickItem(key) {
    setPageParams(defalutPageParams)
    setKeyword('')
    setActiveTabKey(key)
    queryCurrentResources(key)
  }

  function handleClickSetting(key) {
    if (key !== activeTabKey) {
      handleClickItem(key)
    }
    setOpen(true)
    setIsEdit(true)
    if (key !== 'namespace') {
      const currentTabData = tabList?.filter((item: any) => item?.name === key)
      setCurrentItem(currentTabData?.[0])
    }
  }

  function handleRuleFormClose() {
    setOpen(false)
    setIsEdit(false)
    setCurrentItem(null)
  }

  function handleSort(key) {
    setSortParams({
      orderBy: key,
      isAsc: !sortParams?.isAsc,
    })
    const groups = allResourcesData?.groups?.sort((a, b) =>
      sortParams?.isAsc
        ? b.title.localeCompare(a.title)
        : a.title.localeCompare(b.title),
    )
    setAllResourcesData({
      fields: allResourcesData?.fields,
      groups,
    })
  }

  function getName(group, fields) {
    const obj = {}
    const labelObj = {}
    const specialList = ['labels', 'annotations']
    specialList?.forEach(item => {
      group?.[item] &&
        Object.keys(group?.[item])?.forEach(innerKey => {
          obj[`${item}.${innerKey}`] = group?.[item]?.[innerKey]
          labelObj[innerKey] = group?.[item]?.[innerKey]
        })
    })
    const nameList = []
    fields?.forEach(item => {
      if (item !== 'namespace') {
        obj?.[item] && nameList.push(obj?.[item])
      } else {
        nameList.push(group?.[item] || '')
        labelObj[item] = group?.[item]
      }
    })
    return {
      title: nameList?.join('-'),
      tags: labelObj,
    }
  }

  function handleCardClick(group, title) {
    const obj: any = {
      ruleKey: activeTabKey,
      from: 'insight',
      title: title,
    }
    Object.keys(group)?.forEach(key => {
      if (key !== 'labels' && key !== 'annotations') {
        obj[key] = group?.[key]
      } else {
        Object.keys(group?.[key] || {})?.forEach(innerKey => {
          obj[`${key}__${innerKey}`] = group?.[key]?.[innerKey]
        })
      }
    })
    const urlParams = queryString?.stringify(obj)
    navigate(`/insightDetail/group?${urlParams}`)
  }

  function renderList() {
    return (
      showData?.groups?.length > 0 &&
      showData?.groups?.map((group: any, index: number) => {
        const renderProps = {
          group,
          allTags: Object.keys(group?.tags || {})?.map(key => ({
            key,
            value: group?.tags?.[key],
          })),
          handleClick: handleCardClick,
        }
        return (
          <div
            key={`${group?.id}_${index + 1}`}
            style={{
              width: '25%',
              display: 'flex',
              justifyContent: 'center',
            }}
          >
            <CardContent {...renderProps} />
          </div>
        )
      })
    )
  }

  function renderSort() {
    return (
      <Button
        type="link"
        style={{ color: '#646566' }}
        onClick={() => handleSort('name')}
      >
        {t('SortByName')}
        {sortParams?.orderBy === 'name' &&
          (sortParams?.isAsc ? (
            <SortDescendingOutlined style={orderIconStyle} />
          ) : (
            <SortAscendingOutlined style={orderIconStyle} />
          ))}
      </Button>
    )
  }

  useEffect(() => {
    queryAllRules(false)
    queryStats()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  function handleChange(event) {
    setKeyword(event?.target.value)
  }

  const tabsItems = tabList
    ?.sort((a, b) => {
      if (a?.name === 'namespace') return -1
      if (b?.name === 'namespace') return 1
      return Date?.parse(a?.createdAt) - Date?.parse(b?.createdAt)
    })
    ?.map(item => {
      const isDisabledEdit = item?.name === 'namespace' || isReadOnlyMode
      return {
        ...item,
        key: item?.name,
        label: item?.name,
        closeIcon: isDisabledEdit ? null : (
          <EditOutlined
            style={{ fontSize: 12, color: `rgba(0, 0, 0, 0.45)` }}
          />
        ),
      }
    })

  function handleChangeTag(key) {
    handleClickItem(key)
  }

  function onEdit(action, key) {
    action === 'add' ? setOpen(true) : handleClickSetting(key)
  }

  return (
    <div className={styles.insight_wrapper}>
      <div className={styles.container}>
        <h4 className={styles.pageTitle}>{t('Insight')}</h4>
        <QuotaContent statsData={statsData} />
        <div className={styles.content}>
          <InsightTabs
            items={tabsItems}
            addIsDiasble={isReadOnlyMode}
            activeKey={activeTabKey}
            handleClickItem={handleChangeTag}
            onEdit={onEdit}
            disabledAdd={isReadOnlyMode}
          />
          <div className={styles.action_bar}>
            <div className={styles.action_bar_left}>
              <Input
                placeholder={t('KeywordSearch')}
                suffix={<SearchOutlined />}
                style={{ width: 260 }}
                value={keyword}
                onChange={handleChange}
                allowClear
              />
            </div>
            <div className={styles.action_bar_right}>
              <div className={styles.action_bar_right_sort}>{renderSort()}</div>
            </div>
          </div>
          {resouresLoading ? (
            <div className={styles.loading_box}>
              <Loading />
            </div>
          ) : (
            <div className={styles.pageList}>
              {renderList()}
              {allResourcesData?.groups &&
                allResourcesData?.groups?.length > 0 && (
                  <div className={styles.footer}>
                    <Pagination
                      total={pageParams?.total}
                      showTotal={(total: number, range: any[]) =>
                        `${range?.[0]}-${range?.[1]} ${t('Total')} ${total} `
                      }
                      pageSize={pageParams?.pageSize}
                      current={pageParams?.pageNo}
                      onChange={handleChangePage}
                    />
                  </div>
                )}
              {(!allResourcesData?.groups ||
                !allResourcesData?.groups?.length) && (
                <div className={styles.noData}>
                  <Empty style={{ marginTop: 30 }} />
                </div>
              )}
            </div>
          )}
        </div>
        <RuleForm
          onClose={handleRuleFormClose}
          open={open}
          handleSubmit={handleSubmit}
          isEdit={isEdit}
          currentItem={currentItem}
          deleteItem={deleteItem}
        />
      </div>
    </div>
  )
}

export default Insight
