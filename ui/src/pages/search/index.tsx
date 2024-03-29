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

import React, { useState } from 'react'
import { Tag } from 'antd'
import { DoubleLeftOutlined, DoubleRightOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'
import KarbourTabs from '@/components/tabs/index'
import logoJPG from '@/assets/logo.jpg'
import SqlEditor from '@/components/sqlSearch'

import styles from './styles.module.less'

const SearchPage = () => {
  const navigate = useNavigate()
  const [searchType, setSearchType] = useState<string>('sql')
  const [sqlEditorValue, setSqlEditorValue] = useState<any>('')

  const [showAll, setShowAll] = useState(false)

  const { t } = useTranslation()

  const tabsList = [
    { label: t('KeywordSearch'), value: 'keyword', disabled: true },
    { label: t('SQLSearch'), value: 'sql' },
  ]

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

  return (
    <div className={styles.container}>
      <div className={styles.search}>
        <div className={styles.title}>
          <img src={logoJPG} width="30%" alt="icon" />
        </div>
        <div className={styles.searchTab}>
          <KarbourTabs
            list={tabsList}
            current={searchType}
            onChange={handleTabChange}
          />
        </div>
        <SqlEditor
          sqlEditorValue={sqlEditorValue}
          handleSearch={handleSearch}
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
            </div>
          ) : (
            <div className={styles.sql}>
              <div
                className={styles.karbour_tag}
                onClick={() => handleClickSql(`where kind='Namespace'`)}
              >
                <span className={styles.keyword}>select</span> *{' '}
                <span className={styles.keyword}>from</span> resources{' '}
                <span className={styles.keyword}>where </span>
                kind='Namespace'
              </div>
              <div
                className={styles.karbour_tag}
                onClick={() => handleClickSql(`where kind!='Pod'`)}
              >
                <span className={styles.keyword}>select</span> *{' '}
                <span className={styles.keyword}>from</span> resources{' '}
                <span className={styles.keyword}>where </span>
                kind!='Pod'
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
                <span className={styles.keyword}>where </span>
                kind not in ('pod','service')
              </div>
              <div
                className={styles.karbour_tag}
                onClick={() =>
                  handleClickSql(
                    `where ${"`annotations.app.kubernetes.io/name` = 'demoapp'"}`,
                  )
                }
              >
                <span className={styles.keyword}>select</span> *{' '}
                <span className={styles.keyword}>from</span> resources{' '}
                <span className={styles.keyword}>where </span>
                `annotations.app.kubernetes.io/name` = 'demoapp'
              </div>
              {!showAll && (
                <div className={styles.toggleButton} onClick={toggleTags}>
                  <span>
                    {t('More')}
                    <DoubleLeftOutlined
                      style={{
                        transform: 'rotate(-90deg)',
                        marginLeft: 5,
                      }}
                    />
                  </span>
                </div>
              )}
              {showAll && (
                <>
                  <div
                    className={styles.karbour_tag}
                    onClick={() =>
                      handleClickSql(
                        `where kind='Service' order by creationTimestamp desc`,
                      )
                    }
                  >
                    <span className={styles.keyword}>select</span> *{' '}
                    <span className={styles.keyword}>from</span> resources{' '}
                    <span className={styles.keyword}>where </span>
                    kind='Service' order by creationTimestamp desc
                  </div>
                  <div
                    className={styles.karbour_tag}
                    onClick={() =>
                      handleClickSql(
                        `where kind='Deployment' and creationTimestamp < '2024-01-01T18:00:00Z'`,
                      )
                    }
                  >
                    <span className={styles.keyword}>select</span> *{' '}
                    <span className={styles.keyword}>from</span> resources{' '}
                    <span className={styles.keyword}>where </span>
                    {`kind='Deployment' and creationTimestamp < '2024-01-01T18:00:00Z'`}
                  </div>
                  <div
                    className={styles.karbour_tag}
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
                  <div className={styles.toggleButton} onClick={toggleTags}>
                    <span>
                      {t('Less')}
                      <DoubleRightOutlined
                        style={{
                          transform: 'rotate(-90deg)',
                          marginLeft: 5,
                        }}
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
