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

import { useEffect, useRef, useState } from 'react'
import { Button, Input, Tag, Space, AutoComplete } from 'antd'
import { useNavigate } from 'react-router-dom'
import {
  CloseOutlined,
  DeleteOutlined,
  HistoryOutlined,
} from '@ant-design/icons'
import { searchPrefix } from '../../utils/constants'
import KarbourTabs from '../../components/Tabs/index'
// import SqlSearchBox from "./components/sqlSearchBox/index";
import logoJPG from '../../assets/logo.jpg'

import styles from './styles.module.less'
import React from 'react'

const tabsList = [
  { label: '按照关键字搜索', value: 'keyword', disabled: true },
  { label: '按照 SQL 搜索', value: 'sql' },
]

const SearchPage = () => {
  const navigate = useNavigate()
  const [searchType, setSearchType] = useState<string>('sql')
  const [inputValue, setInputValue] = useState('')
  const [open, setOpen] = useState<boolean>(false)
  const [options, setOptions] = useState<{ value: string }[]>([])
  const [optionsCopy, setOptionsCopy] = useState<{ value: string }[]>([])
  const optionsRef = useRef<any>(getHistoryList())

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

  function deleteHistory() {
    localStorage.setItem(`${searchType}History`, '')
  }

  function handlePanelFocus(evt) {
    evt.preventDefault()
    evt.stopPropagation()
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
    <div
      className={styles.search}
      onClick={evt => {
        evt.preventDefault()
        evt.stopPropagation()
        setOpen(false)
      }}
    >
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
      <div className={styles.searchBox}>
        <Space.Compact>
          {/* <Input
            onFocus={(evt) => {
              evt.preventDefault();
              evt.stopPropagation();
              setOpen(true)
            }}
            onClick={(evt) => {
              evt.preventDefault();
              evt.stopPropagation();
            }}
            // onBlur={() => setOpen(false)}
            placeholder={
              searchType === "keyword" ? "支持搜索集群，集群资源（service/pod/cafed）..." : "支持sql语句查询"
            }
            prefix={<SearchOutlined style={{ color: '#999' }} />}
            allowClear
            style={{ width: 600 }}
            value={inputValue}
            onChange={handleInputChange} /> */}
          <Input disabled value={searchPrefix} style={{ width: 180 }} />
          <AutoComplete
            onKeyUp={handleOnkeyUp}
            options={options}
            // onSearch={(text) => }
            placeholder={
              searchType === 'keyword'
                ? '支持搜索集群，集群资源（service/pod/cafed）...'
                : '支持 SQL 语句查询'
            }
            filterOption={(inputValue, option: any) => {
              return (
                option!.value
                  ?.toUpperCase()
                  .indexOf(inputValue.toUpperCase()) !== -1
              )
            }}
            style={{ width: 600 }}
            value={inputValue}
            allowClear={true}
            onChange={handleInputChange}
          />
          <Button type="primary" onClick={handleSearch}>
            搜索
          </Button>
        </Space.Compact>
        {open && (
          <div
            className={styles.searchPanel}
            onFocus={handlePanelFocus}
            onClick={evt => {
              evt.preventDefault()
              evt.stopPropagation()
            }}
          >
            <div className={styles.history}>
              <div className={styles.left}>
                <HistoryOutlined /> 历史记录
              </div>
              <div className={styles.right} onClick={deleteHistory}>
                <DeleteOutlined /> 清空
              </div>
            </div>
            <div className={styles.historyList}>
              {options?.length ? (
                options?.length > 0 &&
                options?.map((item: any) => {
                  return (
                    <Tag onClick={() => setInputValue(item)} key={item}>
                      {item}
                    </Tag>
                  )
                })
              ) : (
                <div>暂无历史记录</div>
              )}
            </div>
            {/* <div className={styles.recommand}>
              <div className={styles.recommandTitle}>🔥 热门推荐</div>
              <div className={styles.list}>{
                recommandList?.map(item => {
                  return <Tag key={item}>{item}</Tag>
                })
              }</div>
            </div> */}
          </div>
        )}
      </div>
      {/* codeMirror */}
      {/* <div>
      <SqlSearchBox/>
        <button onClick={handleSearch}>搜索</button>
      </div> */}
      <div className={styles.examples}>
        {searchType === 'keyword' ? (
          <div className={styles.keywords}>
            <div className={styles.keywordsTitle}>关键字搜索案例</div>
            <div className={styles.item}>
              <Tag bordered={false} color="#fff" style={{ color: '#000' }}>
                "my-application"
              </Tag>
            </div>
            <div className={styles.item}>
              <Tag bordered={false} color="#fff" style={{ color: '#000' }}>
                <span className={styles.keyword}>name:</span>
                /.*my-application.*/kind:pod
              </Tag>
            </div>
            <div className={styles.item}>
              <Tag bordered={false} color="#fff" style={{ color: '#000' }}>
                <span className={styles.keyword}>cluster:</span>xxxkind:service
              </Tag>
            </div>
          </div>
        ) : (
          <div className={styles.sql}>
            <div className={styles.keywordsTitle}>SQL 搜索案例</div>
            <div
              className={styles.item}
              onClick={() => handleClickSql(`where kind='Namespace'`)}
            >
              <Tag bordered={false} color="#fff" style={{ color: '#000' }}>
                <span className={styles.keyword}>select</span> *{' '}
                <span className={styles.keyword}>from</span> resources{' '}
                <span className={styles.keyword}>where </span>kind='Namespace'
              </Tag>
            </div>
            <div
              className={styles.item}
              onClick={() => handleClickSql(`where kind!='Pod'`)}
            >
              <Tag bordered={false} color="#fff" style={{ color: '#000' }}>
                <span className={styles.keyword}>select</span> *{' '}
                <span className={styles.keyword}>from</span> resources{' '}
                <span className={styles.keyword}>where </span>kind!='Pod'
              </Tag>
            </div>
            <div
              className={styles.item}
              onClick={() => handleClickSql(`where namespace='default'`)}
            >
              <Tag bordered={false} color="#fff" style={{ color: '#000' }}>
                <span className={styles.keyword}>select</span> *{' '}
                <span className={styles.keyword}>from</span> resources{' '}
                <span className={styles.keyword}>where </span>
                namespace='default'
              </Tag>
            </div>
            <div
              className={styles.item}
              onClick={() =>
                handleClickSql(`where cluster='democluster' and kind='Pod'`)
              }
            >
              <Tag bordered={false} color="#fff" style={{ color: '#000' }}>
                <span className={styles.keyword}>select</span> *{' '}
                <span className={styles.keyword}>from</span> resources{' '}
                <span className={styles.keyword}>where </span>
                cluster='democluster' and kind='Pod'
              </Tag>
            </div>
            <div
              className={styles.item}
              onClick={() =>
                handleClickSql(`where kind not in ('pod','service')`)
              }
            >
              <Tag bordered={false} color="#fff" style={{ color: '#000' }}>
                <span className={styles.keyword}>select</span> *{' '}
                <span className={styles.keyword}>from</span> resources{' '}
                <span className={styles.keyword}>where </span>kind not in
                ('pod','service')
              </Tag>
            </div>
            <div
              className={styles.item}
              onClick={() =>
                handleClickSql(
                  `where kind='Service' order by object.metadata.creationTimestamp desc`,
                )
              }
            >
              <Tag bordered={false} color="#fff" style={{ color: '#000' }}>
                <span className={styles.keyword}>select</span> *{' '}
                <span className={styles.keyword}>from</span> resources{' '}
                <span className={styles.keyword}>where </span>kind='Service'
                order by object.metadata.creationTimestamp desc
              </Tag>
            </div>
            <div
              className={styles.item}
              onClick={() =>
                handleClickSql(
                  `where kind='Deployment' and object.metadata.creationTimestamp < '2024-01-01T18:00:00Z'`,
                )
              }
            >
              <Tag bordered={false} color="#fff" style={{ color: '#000' }}>
                <span className={styles.keyword}>select</span> *{' '}
                <span className={styles.keyword}>from</span> resources{' '}
                <span className={styles.keyword}>where </span>
                {`kind='Deployment' and object.metadata.creationTimestamp < '2024-01-01T18:00:00Z'`}
              </Tag>
            </div>
            <div
              className={styles.item}
              onClick={() =>
                handleClickSql(
                  `where kind='Pod' and object.metadata.creationTimestamp between '2024-01-01T18:00:00Z' and '2024-01-11T18:00:00Z' order by object.metadata.creationTimestamp`,
                )
              }
            >
              <Tag bordered={false} color="#fff" style={{ color: '#000' }}>
                <span className={styles.keyword}>select</span> *{' '}
                <span className={styles.keyword}>from</span> resources{' '}
                <span className={styles.keyword}>where </span>kind='Pod' and
                object.metadata.creationTimestamp between '2024-01-01T18:00:00Z'
                <br /> and '2024-01-11T18:00:00Z' order by
                object.metadata.creationTimestamp
              </Tag>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export default SearchPage
