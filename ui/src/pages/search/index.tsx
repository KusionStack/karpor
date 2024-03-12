/*
 * Copyright The Karbour Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React, { useEffect, useRef, useState } from 'react'
import { Tag } from 'antd'
import { useNavigate } from 'react-router-dom'
import {
  CloseOutlined,
  DoubleLeftOutlined,
  DoubleRightOutlined,
} from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import KarbourTabs from '@/components/tabs/index'
import logoJPG from '@/assets/logo.jpg'
import SearchInput from '@/components/searchInput'
// import SqlEditor from './sqlSearch'

import styles from './styles.module.less'

const SearchPage = () => {
  const navigate = useNavigate()
  const [searchType, setSearchType] = useState<string>('sql')
  const [inputValue, setInputValue] = useState('')
  const [options, setOptions] = useState<{ value: string }[]>([])
  const [optionsCopy, setOptionsCopy] = useState<{ value: string }[]>([])
  const optionsRef = useRef<any>(getHistoryList())

  const [showAll, setShowAll] = useState(false)

  const { t } = useTranslation()

  const tabsList = [
    { label: t('KeywordSearch'), value: 'keyword', disabled: true },
    { label: t('SQLSearch'), value: 'sql' },
  ]

  // 创建一个函数来切换展示状态
  const toggleTags = () => {
    setShowAll(!showAll)
  }

  function getHistoryList() {
    const historyList: any = localStorage?.getItem(`${searchType}History`)
      ? JSON.parse(localStorage?.getItem(`${searchType}History`))
      : []
    return historyList
  }

  function deleteHistoryByItem(searchType: string, val: string) {
    const lastHistory: any = localStorage.getItem(`${searchType}History`)
    const tmp = lastHistory ? JSON.parse(lastHistory) : []
    if (tmp?.length > 0 && tmp?.includes(val)) {
      const newList = tmp?.filter(item => item !== val)
      localStorage.setItem(`${searchType}History`, JSON.stringify(newList))
    }
  }

  function deleteItem(event, value) {
    event.preventDefault()
    event.stopPropagation()
    deleteHistoryByItem(searchType, value)
    optionsRef.current = getHistoryList()
    setOptionsCopy(optionsRef.current)
  }

  useEffect(() => {
    const tmpOption = optionsRef.current?.map(item => ({
      label: (
        <div className={styles.option_item}>
          <div className={styles.option_item_label}>{item}</div>
          <div
            className={styles.option_item_delete}
            onClick={event => deleteItem(event, item)}
          >
            <CloseOutlined style={{ color: '#808080' }} />
          </div>
        </div>
      ),
      value: item,
    }))
    setOptions(tmpOption)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [optionsCopy])

  const handleTabChange = (value: string) => {
    setSearchType(value)
  }

  function cacheHistory(searchType: string, val: string) {
    const lastHistory: any = localStorage.getItem(`${searchType}History`)
    const tmp = lastHistory ? JSON.parse(lastHistory) : []
    if (tmp?.length > 0 && tmp?.includes(val)) {
      return
    } else {
      const newList = [val, ...tmp]
      localStorage.setItem(`${searchType}History`, JSON.stringify(newList))
      optionsRef.current = getHistoryList()
      setOptionsCopy(optionsRef.current)
    }
  }

  const handleSearch = () => {
    // if (!inputValue) {
    //   message.warning("请输入查询条件")
    //   return
    // }
    if (inputValue) {
      cacheHistory(searchType, inputValue)
    }
    if (searchType.toLocaleUpperCase() === 'sql') {
      navigate(`/search/result?query=${inputValue}&pattern=sql`)
    } else {
      navigate(`/search/result?query=${inputValue}&pattern=sql`)
    }
  }

  const handleInputChange = (value: any) => {
    setInputValue(value)
  }

  function handleClickSql(str) {
    setInputValue(str)
  }

  function handleOnkeyUp(event) {
    if (event?.code === 'Enter' && event?.keyCode === 13) {
      handleSearch()
    }
  }

  return (
    <div className={styles.container}>
      <div className={styles.search}>
        {/* <SqlEditor /> */}
        <div className={styles.title}>
          {/* Hi~ 欢迎来到KarBour */}
          <img src={logoJPG} width="30%" alt="icon" />
        </div>
        {/* <div className={styles.subTitle}>
        你可以通过搜索，匹配集群及其所有资源，帮你轻松管理
      </div> */}
        <div className={styles.searchTab}>
          <KarbourTabs
            list={tabsList}
            current={searchType}
            onChange={handleTabChange}
          />
        </div>
        <SearchInput
          value={inputValue}
          handleSearch={handleSearch}
          handleOnkeyUp={handleOnkeyUp}
          options={options}
          handleInputChange={handleInputChange}
        />
        <div className={styles.examples}>
          {searchType === 'keyword' ? (
            <div className={styles.keywords}>
              <div className={styles.keywordsTitle}>{t('PopularQueries')}</div>
              <div className={styles.item}>
                <Tag style={{ color: '#000' }}>"my-application"</Tag>
              </div>
              <div className={styles.item}>
                <Tag style={{ color: '#000' }}>
                  <span className={styles.keyword}>name:</span>
                  /.*my-application.*/kind:pod
                </Tag>
              </div>
              <div className={styles.item}>
                <Tag style={{ color: '#000' }}>
                  <span className={styles.keyword}>cluster:</span>
                  xxxkind:service
                </Tag>
              </div>
            </div>
          ) : (
            <div className={styles.sql}>
              <div
                className={styles.karbour_tag}
                onClick={() => handleClickSql(`where kind='Namespace'`)}
              >
                <span className={styles.keyword}>select</span> *{' '}
                <span className={styles.keyword}>from</span> resources{' '}
                <span className={styles.keyword}>where </span>kind='Namespace'
              </div>
              <div
                className={styles.karbour_tag}
                onClick={() => handleClickSql(`where kind!='Pod'`)}
              >
                <span className={styles.keyword}>select</span> *{' '}
                <span className={styles.keyword}>from</span> resources{' '}
                <span className={styles.keyword}>where </span>kind!='Pod'
              </div>
              <div
                className={styles.karbour_tag}
                onClick={() => handleClickSql(`where namespace='default'`)}
              >
                <span className={styles.keyword}>select</span> *{' '}
                <span className={styles.keyword}>from</span> resources{' '}
                <span className={styles.keyword}>where </span>
                namespace='default'
              </div>
              <div
                className={styles.karbour_tag}
                onClick={() =>
                  handleClickSql(`where cluster='democluster' and kind='Pod'`)
                }
              >
                <span className={styles.keyword}>select</span> *{' '}
                <span className={styles.keyword}>from</span> resources{' '}
                <span className={styles.keyword}>where </span>
                cluster='democluster' and kind='Pod'
              </div>
              <div
                className={styles.karbour_tag}
                onClick={() =>
                  handleClickSql(`where kind not in ('pod','service')`)
                }
              >
                <span className={styles.keyword}>select</span> *{' '}
                <span className={styles.keyword}>from</span> resources{' '}
                <span className={styles.keyword}>where </span>kind not in
                ('pod','service')
              </div>
              {!showAll && (
                <div className={styles.toggleButton} onClick={toggleTags}>
                  <span>
                    {t('More')}
                    <DoubleLeftOutlined
                      style={{ transform: 'rotate(-90deg)', marginLeft: 5 }}
                    />
                  </span>
                </div>
              )}
              {/* 当showAll为true时，显示收起按钮 */}
              {showAll && (
                <>
                  <div
                    className={styles.karbour_tag}
                    onClick={() =>
                      handleClickSql(
                        `where kind='Service' order by object.metadata.creationTimestamp desc`,
                      )
                    }
                  >
                    <span className={styles.keyword}>select</span> *{' '}
                    <span className={styles.keyword}>from</span> resources{' '}
                    <span className={styles.keyword}>where </span>kind='Service'
                    order by object.metadata.creationTimestamp desc
                  </div>
                  <div
                    className={styles.karbour_tag}
                    onClick={() =>
                      handleClickSql(
                        `where kind='Deployment' and object.metadata.creationTimestamp < '2024-01-01T18:00:00Z'`,
                      )
                    }
                  >
                    <span className={styles.keyword}>select</span> *{' '}
                    <span className={styles.keyword}>from</span> resources{' '}
                    <span className={styles.keyword}>where </span>
                    {`kind='Deployment' and object.metadata.creationTimestamp < '2024-01-01T18:00:00Z'`}
                  </div>
                  <div
                    className={styles.karbour_tag}
                    onClick={() =>
                      handleClickSql(
                        `where kind='Pod' and object.metadata.creationTimestamp between '2024-01-01T18:00:00Z' and '2024-01-11T18:00:00Z' order by object.metadata.creationTimestamp`,
                      )
                    }
                  >
                    <span className={styles.keyword}>select</span> *{' '}
                    <span className={styles.keyword}>from</span> resources{' '}
                    <span className={styles.keyword}>where </span>kind='Pod' and
                    object.metadata.creationTimestamp between
                    '2024-01-01T18:00:00Z'
                    <br /> and '2024-01-11T18:00:00Z' order by
                    object.metadata.creationTimestamp
                  </div>
                  <div className={styles.toggleButton} onClick={toggleTags}>
                    <span>
                      {t('Less')}
                      <DoubleRightOutlined
                        style={{ transform: 'rotate(-90deg)', marginLeft: 5 }}
                      />
                    </span>
                  </div>
                </>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default SearchPage
