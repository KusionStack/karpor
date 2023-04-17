import React, { useState } from '@alipay/bigfish/react';
import { Input, Select, Space, Pagination, } from 'antd';
import { request } from '@alipay/bigfish';
import { searchList } from "@/utils/mockData";
import { useLocation } from "@alipay/bigfish"
import { CaretRightOutlined, RadarChartOutlined } from "@ant-design/icons"
import styles from "./styles.less";

const { Search } = Input;

const options = [
  {
    value: 'filter',
    label: 'Filter(s)',
  },
  {
    value: 'AI Suggestion',
    label: 'AI Suggestion',
  },
  {
    value: 'SQL',
    label: 'SQL',
  },
];


export default () => {
  const location = useLocation();
  const [showPanel, setShowPanel] = useState<Boolean>(false);
  const [pageData, setPageData] = useState([]);
  const [searchValue, setSearchValue] = useState(decodeURIComponent(location.search.split("=")[1]))

  const [searchParams, setSearchParams] = useState({
    pageSize: 10,
    page: 1,
  })


  function handleChangePage(page: number, pageSize: number) {
    console.log(page, pageSize, 'handleChangePage')
    setSearchParams({
      ...searchParams,
      page,
      pageSize,
    })
  }

  async function getPageData() {
    let url = `/apis/search.karbour.com/v1beta1/search/proxy`;
    const res = await request(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      data: {
        searchValue
      }
    })
    console.log(res, "=====res====")
    if (res?.success) {
      setPageData(res?.data || searchList);
    }
  }

  useState(() => {
    getPageData();
  }, [])

  const handleChange = (event) => {
    setSearchValue(event.target.value);
  }

  const handleSearch = (val) => {
    getPageData();
  }



  return (
    <div className={styles.container}>
      <div style={{ width: 850, position: "relative" }}>
        <Space.Compact>
          <Select defaultValue="filter" options={options} style={{ width: 150 }} />
          <Search placeholder="input search text" allowClear style={{ width: 700 }} onFocus={() => setShowPanel(true)} onBlur={() => setShowPanel(false)}
            value={searchValue}
            onChange={handleChange}
            onSearch={handleSearch}
          />
        </Space.Compact>
        {/* {
          showPanel && <div style={{ position: "absolute", background: '#fff', width: '100%', height: 500, zIndex: 1 }}>test</div>
        } */}
      </div>
      <div className={styles.content}>
        {
          searchList?.map((item: any) => {
            return <div className={styles.card} key={item.name}>
              <div><RadarChartOutlined style={{ margin: '0 5px' }} /></div>
              <div className={styles.item}>
                <span className={styles.itemLabel}>ApiVersion</span>
                <span className={styles.itemValue}>{item?.apiVersion}</span>
              </div>
              <CaretRightOutlined style={{ margin: '0 3px', fontSize: 14 }} />
              <div className={styles.item}>
                <span className={styles.itemLabel}>Kind</span>
                <span className={styles.itemValue}>{item?.kind}</span>
              </div>
              {
                item?.metadata?.namespace && <>
                  <CaretRightOutlined style={{ margin: '0 3px', fontSize: 14 }} />
                  <div className={styles.item}>
                    <span className={styles.itemLabel}>NameSpace</span>
                    <span className={styles.itemValue}>{item?.metadata?.namespace}</span>
                  </div>
                </>
              }
              <CaretRightOutlined style={{ margin: '0 3px', fontSize: 14 }} />
              <div className={styles.item}>
                <span className={styles.itemLabel}>Name</span>
                <span className={styles.itemValue}>{item?.metadata?.name}</span>
              </div>
            </div>
          })
        }
        {/* <div style={{ width: 600, height: 600 }}>图片占位</div> */}
      </div>
      <div className={styles.footer}>
        <Pagination
          total={searchList?.length}
          showTotal={(total, range) => `${range[0]}-${range[1]} 共 ${total} 条`}
          pageSize={searchParams?.pageSize}
          current={searchParams?.page}
          onChange={handleChangePage}
        />
      </div>
    </div>
  );
};
