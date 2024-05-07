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

import React, { useEffect, useState } from 'react'
import { Tag } from 'antd'
import { DoubleLeftOutlined, DoubleRightOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'
import KarporTabs from '@/components/tabs/index'
import logoFull from '@/assets/img/logo-full.svg'
import SqlSearch from '@/components/sqlSearch'
import { defaultSqlExamples, tabsList } from '@/utils/constants'

import styles from './styles.module.less'

const SearchPage = () => {
  const navigate = useNavigate()
  const { t } = useTranslation()
  const [searchType, setSearchType] = useState<string>('sql')
  const [sqlEditorValue, setSqlEditorValue] = useState<any>('')
  const [showAll, setShowAll] = useState(false)
  const [scale, setScale] = useState(1)

  const toggleTags = () => {
    setShowAll(!showAll)
  }

  function handleTabChange(value: string) {
    setSearchType(value)
  }

  function handleClickSql(str) {
    setSqlEditorValue(str)
  }

  function handleSearch(inputValue) {
    navigate(`/search/result?query=${inputValue}&pattern=sql`)
  }

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
      `kind='Deployment' and creationTimestamp < '2024-01-01T18:00:00Z'`,
    ]
    return renderSqlExamples(sqlExamples)
  }

  useEffect(() => {
    const handleResize = () => {
      const innerWidth = window.innerWidth
      if (innerWidth >= 1200) {
        setScale(1)
      } else if (innerWidth < 1200 && innerWidth >= 1100) {
        setScale(0.9)
      } else if (innerWidth < 1100 && innerWidth >= 900) {
        setScale(0.8)
      } else {
        setScale(0.6)
      }
    }
    handleResize()
    window.addEventListener('resize', handleResize)
    return () => {
      window.removeEventListener('resize', handleResize)
    }
  }, [])

  return (
    <div className={styles.container} style={{ transform: `scale(${scale})` }}>
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
        <SqlSearch
          sqlEditorValue={sqlEditorValue}
          handleSearch={handleSearch}
        />
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
          ) : (
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
                        `where kind='Pod' and creationTimestamp between '2024-01-01T18:00:00Z' and '2024-01-11T18:00:00Z' order by creationTimestamp`,
                      )
                    }
                  >
                    <span className={styles.keyword}>select</span> *{' '}
                    <span className={styles.keyword}>from</span> resources{' '}
                    <span className={styles.keyword}>where </span>
                    kind='Pod' and creationTimestamp between
                    '2024-01-01T18:00:00Z'
                    <br /> and '2024-01-11T18:00:00Z' order by creationTimestamp
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
