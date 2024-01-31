import React, { useEffect, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { ArrowLeftOutlined, UploadOutlined } from '@ant-design/icons'
import { Form, Input, Space, Button, Upload, Radio, message } from 'antd'
import type { RadioChangeEvent } from 'antd'
import ReactDiffViewer, { DiffMethod } from 'react-diff-viewer-continued'
import axios from 'axios'
import queryString from 'query-string'
import { useSelector } from 'react-redux'
import yaml from 'js-yaml'
import { yaml2json } from '@/utils/tools'
// import 'react-diff-viewer-continued/es/styles/index.css';
import Yaml from '../../../components/yaml'
import { HOST } from '../../../utils/request'

import styles from './styles.module.less'

const { TextArea } = Input

type UploadConfigProps = {
  onChange: (val) => void
  fileList: any
}

export const UploadConfig = (props: UploadConfigProps) => {
  const [radioValue, setRadioValue] = useState('file')

  const onRadioChange = (e: RadioChangeEvent) => {
    setRadioValue(e.target.value)
    props?.onChange({
      type: e.target.value,
      value: '',
    })
  }

  const handleTextAreaChange = event => {
    props?.onChange({
      ...props?.fileList,
      value: event.target.value,
    })
  }

  return (
    <div>
      <div style={{ marginBottom: 15 }}>
        <Radio.Group onChange={onRadioChange} value={radioValue}>
          <Radio value="file">文件配置</Radio>
          <Radio value="yaml">输入yaml</Radio>
        </Radio.Group>
      </div>
      {radioValue === 'file' ? (
        <Upload name="logo" action="/upload.do">
          <Button icon={<UploadOutlined />}>上传配置文件</Button>
        </Upload>
      ) : (
        <TextArea onChange={handleTextAreaChange} />
      )}
    </div>
  )
}

const ClusterCertificate = () => {
  const [form] = Form.useForm()
  const navigate = useNavigate()
  const { isReadOnlyMode } = useSelector((state: any) => state.globalSlice)
  const location = useLocation()
  const { cluster } = queryString.parse(location?.search)
  const [newYamlContent, setNewYamlContent] = useState<any>()
  const [loading, setLoading] = useState(false)
  const [lastYamlContent, setLastYamlContent] = useState('')
  const [lastYamlContentJson, setLastYamlContentJson] = useState<any>()

  function onCancel() {
    form.resetFields()
    navigate(-1)
  }

  async function getClusterDetail() {
    const response: any = await axios({
      url: `/rest-api/v1/cluster/${cluster}`,
      method: 'GET',
      params: {
        format: 'yaml',
      },
    })
    if (response?.success) {
      setLastYamlContent(response?.data)
      setLastYamlContentJson(yaml2json(response?.data)?.data)
    } else {
      message.error(response?.message || '请求失败，请重试')
    }
  }

  useEffect(() => {
    if (cluster) {
      getClusterDetail()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [cluster])

  async function onFinish() {
    if (isReadOnlyMode) {
      return
    }
    if (!newYamlContent?.content) {
      message.warning('请上传新的 KubeConfig 配置文件')
    } else {
      setLoading(true)
      const validateResponse: any = await axios.post(
        '/rest-api/v1/cluster/config/validate',
        {
          kubeConfig: newYamlContent?.content,
        },
      )
      if (validateResponse?.success) {
        const response: any = await axios({
          url: `/rest-api/v1/cluster/${cluster}`,
          method: 'PUT',
          data: {
            kubeConfig: newYamlContent?.content,
          },
        })
        if (response?.success) {
          setLoading(false)
          message.success('验证成功并提交，3s 后将跳转到集群管理页面')
          setTimeout(() => {
            navigate('/cluster')
          }, 3000)
        } else {
          message.error(response?.message || '验证成功但提交失败')
          setLoading(false)
        }
      } else {
        message.error(
          validateResponse?.message || 'KubeConfig 不符合要求，请验证后上传',
        )
        setLoading(false)
      }
    }
  }

  function goBack() {
    navigate(-1)
  }

  const uploadProps: any = {
    disabled: isReadOnlyMode,
    name: 'file',
    accept: '.yaml,.yml,.json,.kubeconfig,.kubeconf',
    action: `${HOST}/rest-api/v1/cluster/config/file`,
    headers: {
      authorization: 'authorization-text',
    },
    method: 'POST',
    data: {
      name: lastYamlContentJson?.metadata?.name,
      description: lastYamlContentJson?.spec?.description,
      displayName: lastYamlContentJson?.spec?.displayName,
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
      if (info.file.status !== 'uploading') {
      }
      if (info.file.status === 'done') {
        if (info?.file?.response?.success) {
          message.success(`${info.file.name}上传成功`)
          form.setFieldsValue({
            kubeConfig: info?.file?.response?.data?.content,
          })
          setNewYamlContent({
            content: info?.file?.response?.data?.content,
            sanitizedClusterContent:
              info?.file?.response?.data?.sanitizedClusterContent,
          })
        } else {
          message.error(
            info?.file?.response?.message ||
              '文件只支持.yaml, .yml, .json, .kubeconfig, .kubeconf',
          )
          // form.setFieldsValue({
          //   kubeConfig: undefined
          // })
          // setNewYamlContent('')
        }
      } else if (info.file.status === 'error') {
        message.error(`${info.file.name} file upload failed.`)
      }
    },
  }

  const newStyles = {
    variables: {
      dark: {
        highlightBackground: '#fefed5',
        highlightGutterBackground: '#ffcd3c',
      },
    },
    line: {
      padding: '10px 2px',
      '&:hover': {
        background: '#a26ea1',
      },
    },
  }

  const oldYaml = yaml.dump(lastYamlContent)
  const newYaml = yaml.dump(newYamlContent?.sanitizedClusterContent)

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <ArrowLeftOutlined
          style={{ marginRight: 10 }}
          onClick={() => goBack()}
        />
        更新证书
      </div>
      <div className={styles.content}>
        <Form
          form={form}
          onFinish={onFinish}
          layout="vertical"
          style={{ width: 600, display: 'flex' }}
          initialValues={{
            type: 'file',
          }}
        >
          <Form.Item
            name="kubeConfig"
            rules={[{ required: true, message: '该配置内容不能为空' }]}
          >
            <Upload {...uploadProps}>
              <Button disabled={isReadOnlyMode} icon={<UploadOutlined />}>
                上传 KubeConfig 配置文件
              </Button>
            </Upload>
          </Form.Item>
          <Form.Item style={{ marginLeft: 20 }}>
            <Space>
              <Button
                disabled={isReadOnlyMode}
                type="primary"
                htmlType="submit"
                loading={loading}
              >
                {loading ? '验证提交中' : '验证并更新'}
              </Button>
              <Button htmlType="button" onClick={onCancel}>
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
        {newYamlContent?.sanitizedClusterContent ? (
          <div className={styles.config_content}>
            {/* <div className={styles.title}>左侧为原配置信息，右侧为新配置信息</div> */}
            <div className={styles.diff_container}>
              <ReactDiffViewer
                leftTitle="原配置信息"
                rightTitle="新配置信息"
                styles={newStyles}
                oldValue={oldYaml}
                newValue={newYaml}
                splitView={true}
                useDarkTheme={false}
                compareMethod={DiffMethod.LINES}
              />
            </div>
          </div>
        ) : (
          <div className={styles.config_content}>
            <div className={styles.title}>原配置信息</div>
            <Yaml data={lastYamlContent} height="100%" />
          </div>
        )}
      </div>
    </div>
  )
}

export default ClusterCertificate
