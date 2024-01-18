import { SearchOutlined } from "@ant-design/icons";
import { Button, Input, Space, Table, message } from "antd";
import axios from "axios";
import queryString from "query-string";
import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";

import styles from "./style.module.less";

type IProps = {
  queryStr: string;
  data?: any[];
  searchKey?: string;
  pagination?: any;
  tableName?: string;
};

const defaultSearchParams = {
  current: 1,
  pageSize: 10,
  total: 0,
};

const SourceTable = ({ queryStr, tableName }: IProps) => {
  const location = useLocation();
  const navigate = useNavigate();
  const [pageParams, setPageParams] = useState(defaultSearchParams);
  const [tableData, setTableData] = useState([]);
  const urlSearchParams = queryString?.parse(location?.search);

  function goResourcePage(record) {
    const params = {
      from: urlSearchParams?.from,
      cluster: urlSearchParams?.cluster,
      name: record?.object?.metadata?.name,
      kind: record?.object?.kind,
      apiVersion: record?.object?.apiVersion,
      namespace: record?.object?.metadata?.namespace,
      query: urlSearchParams?.query,
    };
    const urlParams = queryString?.stringify(params);
    navigate(`/insightDetail/resource?type=resource&${urlParams}`);
  }

  const columns = [
    {
      dataIndex: "name",
      key: "name",
      title: "名称",
      render: (_, record) => {
        return (
          <Button type="link" onClick={() => goResourcePage(record)}>
            {record?.object?.metadata?.name}
          </Button>
        );
      },
    },
    {
      dataIndex: "namespace",
      key: "namespace",
      title: "Namespace",
      render: (_, record) => {
        return record?.object?.metadata?.namespace;
      },
    },
    {
      dataIndex: "apiVersion",
      key: "apiVersion",
      title: "apiVersion",
      render: (_, record) => {
        return record?.object?.apiVersion;
      },
    },
    {
      dataIndex: "kind",
      key: "kind",
      title: "kind",
      // render: (text) => text === 'success' ? <Badge status="success" text="健康" /> : <Badge status="error" text="异常" />
      render: (_, record) => {
        return record?.object?.kind;
      },
    },
  ];

  async function queryTableData(params) {
    const { current, pageSize } = pageParams;
    const response: any = await axios.get(
      `/rest-api/v1/search?query=${queryStr}&pattern=sql&page=${params?.current || current}&pageSize=${params?.pageSize || pageSize}`,
    );
    if (response?.success) {
      setTableData(response?.data?.items || []);
      setPageParams({
        ...params,
        total: response?.data?.total,
      });
    } else {
      message.error(response?.message);
    }
  }

  useEffect(() => {
    if (queryStr) {
      queryTableData({ current: 1, pageSize: pageParams?.pageSize });
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [queryStr]);

  function handleTableChange({ current, pageSize }) {
    queryTableData({ current, pageSize });
  }

  return (
    <div>
      <div className={styles.table_header}>
        <div className={styles.table_title}>
          {tableName || "--"}
          {urlSearchParams?.type === "kind" ? null : (
            <span className={styles.tips}>
              💡 可在上方选择资源后在这里查看对应的详情
            </span>
          )}
        </div>
        <Space style={{ marginBottom: 10 }}>
          <Input
            disabled
            placeholder="请输入名称搜索"
            suffix={<SearchOutlined />}
          />
        </Space>
      </div>
      <Table
        columns={columns}
        dataSource={tableData}
        rowKey={(record) => {
          return `${record?.object?.metadata?.name}_${record?.object?.metadata?.namespace}_${record?.object?.apiVersion}_${record?.object?.kind}`;
        }}
        onChange={handleTableChange}
        pagination={{
          ...pageParams,
          showSizeChanger: true,
        }}
      />
    </div>
  );
};

export default SourceTable;
