import { useState } from "react";
import {
  Pagination,
  Badge,
  Tooltip,
  Empty,
  Button,
  Drawer,
  Input,
  Space,
  message,
  Form,
} from "antd";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import styles from "./styles.module.less";

const layout = {
  labelCol: { span: 4 },
  wrapperCol: { span: 20 },
};

const Cluster = () => {
  const [form] = Form.useForm();

  const navigate = useNavigate();
  const [pageData, setPageData] = useState<any>([]);
  const [searchParams, setSearchParams] = useState({
    pageSize: 10,
    page: 1,
  });

  const [visible, setVisible] = useState<boolean>(false);
  const [yamlValue, setYamlValue] = useState<string>(undefined);

  async function getPageData() {
    const data = await axios(`/apis/cluster.karbour.com/v1beta1/clusters`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      params: {},
    });
    setPageData(data || {});
  }

  useState(() => {
    getPageData();
  });

  function handleChangePage(page: number, pageSize: number) {
    setSearchParams({
      ...searchParams,
      page,
      pageSize,
    });
  }

  function handleMore(item: any) {
    console.log(item, "====handleMore====");
  }

  const handleClick = (item) => {
    console.log(item, "===item===");
    let queryStr = "";
    if (item?.metadata?.managedFields?.[0]?.apiVersion) {
      queryStr = `${item?.metadata?.managedFields?.[0]?.apiVersion},${item.metadata?.name}`;
    } else {
      queryStr = `${item.metadata?.name}`;
    }
    navigate(`/cluster-detail?query=${queryStr}`);
  };

  const join = () => {
    setVisible(true);
  };

  function onClose() {
    setVisible(false);
  }

  function handleSubmit() {
    form
      .validateFields()
      .then((values) => {
        console.log(values, "=====sada====");
        message.success("添加成功");
        form.resetFields();
        setVisible(false);
      })
      .catch((info) => {
        console.log("Validate Failed:", info);
      });
  }

  return (
    <div className={styles.container}>
      <div className={styles["action-container"]}>
        <Button type="primary" onClick={join}>
          入驻
        </Button>
      </div>
      <div className={styles.content}>
        {pageData?.items?.map((item: any, index: number) => {
          return (
            <div
              className={styles.card}
              key={`${item.name}_${index}`}
              onClick={() => handleClick(item)}
            >
              <div className={styles.header}>
                <div className={styles["header-left"]}>
                  {item.metadata?.name}
                  <Badge
                    style={{
                      marginLeft: 20,
                      fontWeight: "normal",
                      color: item?.status?.healthy === "true" ? "green" : "red",
                    }}
                    status={
                      item?.status?.healthy === "true" ? "success" : "error"
                    }
                    text={item?.status?.healthy === "true" ? "健康" : "不健康"}
                  />
                </div>
                <div
                  className={styles["header-right"]}
                  onClick={() => handleMore(item)}
                >
                  More
                </div>
              </div>
              <div className={styles["card-body"]}>
                <div className={styles.item}>
                  <div className={styles["item-label"]}>Endpoint: </div>
                  <Tooltip title={item.spec?.access?.endpoint}>
                    <div className={styles["item-value"]}>
                      {item.spec?.access?.endpoint}
                    </div>
                  </Tooltip>
                </div>
                <div className={styles.stat}>
                  <div className={styles.node}>
                    Nodes: {item?.status?.node || "--"}
                  </div>
                  <div className={styles.deloy}>
                    Delay: {item?.status?.delay || "--"}
                  </div>
                </div>
              </div>
            </div>
          );
        })}
      </div>
      {pageData?.items && pageData?.items?.length > 0 && (
        <div className={styles.footer}>
          <Pagination
            total={pageData?.items?.length}
            showTotal={(total, range) =>
              `${range[0]}-${range[1]} 共 ${total} 条`
            }
            pageSize={searchParams?.pageSize}
            current={searchParams?.page}
            onChange={handleChangePage}
          />
        </div>
      )}
      {(!pageData?.items || !pageData?.items?.length) && (
        <Empty style={{ marginTop: 30 }} />
      )}
      <Drawer
        width={750}
        title="证书内容"
        placement="right"
        onClose={onClose}
        open={visible}
        extra={
          <Space>
            <Button onClick={onClose}>取消</Button>
            <Button type="primary" onClick={handleSubmit}>
              提交
            </Button>
          </Space>
        }
      >
        <Form {...layout} form={form} name="kubeconfigForm">
          <Form.Item
            name={"name"}
            label="集群名称"
            rules={[{ required: true, message: "集群名称不能为空！" }]}
          >
            <Input style={{ width: 300 }} />
          </Form.Item>
          <Form.Item
            name={"kubeconfig"}
            label="kubeconfig"
            rules={[{ required: true, message: "kubeconfig不能为空！" }]}
          >
            <Input.TextArea autoSize={{ minRows: 7 }} />
          </Form.Item>
        </Form>
      </Drawer>
    </div>
  );
};

export default Cluster;
