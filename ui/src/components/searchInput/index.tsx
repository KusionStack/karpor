import React from 'react'
import { AutoComplete, Input, Space } from 'antd'
import { searchSqlPrefix } from '@/utils/constants'
import arrowRight from '@/assets/arrow-right.png'

import styles from './styles.module.less'

type IProps = {
  options: { value: string }[]
  value: string
  handleInputChange: (val) => void
  handleOnkeyUp: (val: any) => void
  handleSearch: () => void
}

const SearchInput = (props: IProps) => {
  const { handleSearch, handleOnkeyUp, options, value, handleInputChange } =
    props

  return (
    <div className={styles.search_box}>
      <div className={styles.submit} onClick={handleSearch}>
        <img src={arrowRight} />
      </div>
      <Space.Compact>
        <Input
          disabled
          value={searchSqlPrefix}
          style={{
            width: 200,
            fontSize: 16,
            borderRadius: '16px 0 0 16px',
            textAlign: 'center',
          }}
        />
        <AutoComplete
          className={styles.custom_auto_complete}
          size="large"
          onKeyUp={handleOnkeyUp}
          options={options}
          placeholder="Search using SQL ......"
          filterOption={(inputValue, option) =>
            option?.value?.toUpperCase().indexOf(inputValue.toUpperCase()) !==
            -1
          }
          style={{ width: 800, height: 44, fontSize: 24 }}
          value={value}
          allowClear={true}
          onChange={handleInputChange}
        />
      </Space.Compact>
    </div>
  )
}

export default SearchInput
