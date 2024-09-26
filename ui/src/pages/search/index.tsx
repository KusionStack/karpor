/*
 * Copyright The Karpor Authors.
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

import React, { useCallback, useState } from 'react'
import { AutoComplete, Input, message, Space, Tag } from 'antd'
import {
  DoubleLeftOutlined,
  DoubleRightOutlined,
  CloseOutlined,
} from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'
import KarporTabs from '@/components/tabs/index'
import logoFull from '@/assets/img/logo-full.svg'
import SqlSearch from '@/components/sqlSearch'
import { defaultSqlExamples, tabsList } from '@/utils/constants'
import { deleteHistoryByItem, getHistoryList } from '@/utils/tools'

const { Search } = Input

import styles from './styles.module.less'

const SearchPage = () => {
  const navigate = useNavigate()
  const { t } = useTranslation()
  const [searchType, setSearchType] = useState<string>('sql')
  const [sqlEditorValue, setSqlEditorValue] = useState<any>('')
  const [showAll, setShowAll] = useState(false)
  const [naturalOptions, setNaturalOptions] = useState(
    getHistoryList('naturalHistory') || [],
  )

  const toggleTags = () => {
    setShowAll(!showAll)
  }

  function handleTabChange(value: string) {
    setSearchType(value)
  }

  function handleClickSql(str) {
    setSqlEditorValue(str)
  }

  const handleSqlSearch = useCallback(
    inputValue => {
      navigate(`/search/result?query=${inputValue}&pattern=sql`)
    },
    [navigate],
  )

  function renderSqlExamples(data: string[] | null) {
    const sqlExamples = data || defaultSqlExamples
    return sqlExamples?.map(item => (
      <div
        key={item}
        className={styles.karpor_tag}
        onClick={() => handleClickSql(`where ${item}`)}
      >
        <span className={styles.keyword}>select</span> *{' '}
        <span className={styles.keyword}>from</span> resources{' '}
        <span className={styles.keyword}>where </span>
        {item}
      </div>
    ))
  }
  function renderMoreSqlExamples() {
    const sqlExamples = [
      `kind='Service' order by creationTimestamp desc`,
      `kind='Deployment' and creationTimestamp > '2024-01-01T18:00:00Z'`,
    ]
    return renderSqlExamples(sqlExamples)
  }

  function handleNaturalSearch(value) {
    if (!value) {
      message.warning(t('CannotBeEmpty'))
      return
    }
    navigate(`/search/result?query=${value}&pattern=natural`)
  }

  const handleDelete = val => {
    deleteHistoryByItem('naturalHistory', val)
    const list = getHistoryList('naturalHistory') || []
    setNaturalOptions(list)
  }

  const renderOption = val => {
    return (
      <Space
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        <span>{val}</span>
        <CloseOutlined
          onClick={event => {
            event?.stopPropagation()
            handleDelete(val)
          }}
        />
      </Space>
    )
  }

  const tmpOptions = naturalOptions?.map(val => ({
    value: val,
    label: renderOption(val),
  }))

  return (
    <div className={styles.search_container}>
      <div className={styles.search}>
        <div className={styles.title}>
          <img src={logoFull} width="100%" alt="icon" />
        </div>
        <div className={styles.searchTab}>
          <KarporTabs
            list={tabsList}
            current={searchType}
            onChange={handleTabChange}
          />
        </div>

        {searchType === 'sql' && (
          <div className={styles.search_codemirror_container}>
            <SqlSearch
              sqlEditorValue={sqlEditorValue}
              handleSqlSearch={handleSqlSearch}
            />
          </div>
        )}

        {searchType === 'natural' && (
          <div className={styles.search_codemirror_container}>
            <AutoComplete
              style={{ width: '100%' }}
              size="large"
              options={tmpOptions}
              filterOption={(inputValue, option) => {
                if (option?.value) {
                  return (
                    (option?.value as string)
                      ?.toUpperCase()
                      .indexOf(inputValue.toUpperCase()) !== -1
                  )
                }
              }}
            >
              <Search
                size="large"
                placeholder={`${t('SearchByNaturalLanguage')}...`}
                enterButton
                onSearch={handleNaturalSearch}
              />
            </AutoComplete>
          </div>
        )}

        <div className={styles.examples}>
          {searchType === 'keyword' ? (
            <div className={styles.keywords}>
              <div>{t('PopularQueries')}</div>
              <div className={styles.item}>
                <Tag style={{ color: '#000' }}>"my-application"</Tag>
              </div>
              <div className={styles.item}>
                <Tag style={{ color: '#000' }}>
                  <span className={styles.keyword}>name:</span>
                  /.*my-application.*/kind:pod
                </Tag>
              </div>
            </div>
          ) : searchType === 'natural' ? null : (
            <div className={styles.sql}>
              {renderSqlExamples(null)}
              {!showAll && (
                <div className={styles.toggle_button} onClick={toggleTags}>
                  <span>
                    {t('More')}
                    <DoubleLeftOutlined className={styles.toggle_icon} />
                  </span>
                </div>
              )}
              {showAll && (
                <>
                  {renderMoreSqlExamples()}
                  <div
                    className={styles.karpor_tag}
                    onClick={() =>
                      handleClickSql(
                        `where kind='Pod' and creationTimestamp between '2024-01-01T18:00:00Z' and '2025-01-11T18:00:00Z' order by creationTimestamp`,
                      )
                    }
                  >
                    <span className={styles.keyword}>select</span> *{' '}
                    <span className={styles.keyword}>from</span> resources{' '}
                    <span className={styles.keyword}>where </span>
                    kind='Pod' and creationTimestamp between
                    '2024-01-01T18:00:00Z'
                    <br /> and '2025-01-11T18:00:00Z' order by creationTimestamp
                  </div>
                  <div className={styles.toggle_button} onClick={toggleTags}>
                    <span>
                      {t('Less')}
                      <DoubleRightOutlined className={styles.toggle_icon} />
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
