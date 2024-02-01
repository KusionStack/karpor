import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { ArrowLeftOutlined, UploadOutlined } from '@ant-design/icons'
import { Form, Input, Space, Button, Upload, message } from 'antd'
import type { UploadProps } from 'antd'
import axios from 'axios'
import { useSelector } from 'react-redux'
import { HOST } from '../../../utils/request'
import Yaml from '../../../components/yaml'

import styles from './styles.module.less'

const { TextArea } = Input

const AddCluster = () => {
  const { isReadOnlyMode } = useSelector((state: any) => state.globalSlice)
  const [form] = Form.useForm()
  const navigate = useNavigate()
  const [yamlContent, setYamlContent] = useState('')
  const [loading, setLoading] = useState(false)

  async function onFinish(values: any) {
    if (isReadOnlyMode) {
      return
    }
    const tmp = {
      ...values,
    }
    // /rest-api/v1/cluster/config/validate
    setLoading(true)
    const validateResponse: any = await axios.post(
      '/rest-api/v1/cluster/config/validate',
      {
        kubeConfig: tmp?.kubeConfig,
      },
    )
    if (validateResponse?.success) {
      const response: any = await axios({
        url: `/rest-api/v1/cluster/${tmp?.name}`,
        method: 'POST',
        data: tmp,
      })
      if (response?.success) {
        message.success('验证成功并提交，即将跳转到列表页面')
        navigate(-1)
      } else {
        message.error(response?.message || '验证成功但提交失败')
      }
    } else {
      message.error(
        validateResponse?.message ||
          'KubeConfig 不符合要求，请上传合法的证书内容',
      )
      setLoading(false)
    }
  }

  function goBack() {
    navigate(-1)
  }

  // const [radioValue, setRadioValue] = useState("file")
  // const onRadioChange = (e: RadioChangeEvent) => {
  //   setRadioValue(e.target.value);
  // };

  const uploadProps: UploadProps = {
    disabled: isReadOnlyMode,
    name: 'file',
    accept: '.yaml,.yml,.json,.kubeconfig,.kubeconf',
    action: `${HOST}/rest-api/v1/cluster/config/file`,
    headers: {
      authorization: 'authorization-text',
    },
    withCredentials: true,
    maxCount: 1,
    // fileList: fileList,
    showUploadList: {
      showRemoveIcon: false,
      removeIcon: false,
      showPreviewIcon: false,
    },
    onPreview: () => false,
    onChange(info) {
      if (info.file.status === 'done') {
        if (info?.file?.response?.success) {
          message.success(`${info.file.name}上传成功`)
          // setFileList([{
          //   filename: info?.file?.response?.data?.fileName,
          //   uid: '1',
          //   name: info?.file?.response?.data?.fileName,
          //   status: 'done',
          //   url: info?.file?.response?.data?.fileName
          // }])
          form.setFieldsValue({
            kubeConfig: info?.file?.response?.data?.content,
          })
          setYamlContent(info?.file?.response?.data?.content)
        } else {
          message.error(
            info?.file?.response?.message ||
              '文件只支持.yaml, .yml, .json, .kubeConfig, .kubeconf',
          )
          // setFileList([])
          // form.setFieldsValue({
          //   kubeConfig: undefined
          // })
          // setYamlContent('')
        }
      } else if (info.file.status === 'error') {
      }
    },
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <ArrowLeftOutlined style={{ marginRight: 10 }} onClick={goBack} />
        集群接入
      </div>
      <div className={styles.content}>
        <Form
          form={form}
          layout="vertical"
          onFinish={onFinish}
          style={{ width: 400 }}
          initialValues={{
            type: 'file',
          }}
        >
          <Form.Item name="name" label="集群名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item
            name="displayName"
            label="显示名称"
            rules={[{ required: false }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="description"
            label="集群描述"
            rules={[{ required: false }]}
          >
            <TextArea autoSize={{ minRows: 3 }} />
          </Form.Item>
          <Form.Item
            label="kubeConfig"
            name="kubeConfig"
            rules={[{ required: true, message: '配置文件不能为空' }]}
          >
            <Upload {...uploadProps}>
              <Button disabled={isReadOnlyMode} icon={<UploadOutlined />}>
                上传 KubeConfig 配置文件
              </Button>
            </Upload>
          </Form.Item>
          {/* <Form.Item
          rules={[{ required: true }]}
          name="kubeConfig"
          label="Kubeconfig"
          valuePropName="fileList"
          extra="上传yaml文件"
        >
          <UploadConfig />
        </Form.Item> */}
          <Form.Item>
            <Space>
              <Button
                disabled={isReadOnlyMode}
                type="primary"
                htmlType="submit"
                loading={loading}
              >
                验证并接入
              </Button>
              <Button htmlType="button" onClick={goBack}>
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
        <div className={styles.right}>
          {yamlContent && <Yaml data={yamlContent} />}
        </div>
      </div>
    </div>
  )
}

export default AddCluster
