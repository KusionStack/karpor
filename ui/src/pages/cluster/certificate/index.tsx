import React, { useEffect, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { ArrowLeftOutlined, UploadOutlined } from '@ant-design/icons'
import { Form, Space, Button, Upload, message } from 'antd'
import ReactDiffViewer, { DiffMethod } from 'react-diff-viewer-continued'
import { useTranslation } from 'react-i18next'
import axios from 'axios'
import queryString from 'query-string'
import { useSelector } from 'react-redux'
import yaml from 'js-yaml'
import { yaml2json } from '@/utils/tools'
import Yaml from '@/components/yaml'
import { HOST } from '@/utils/request'

import styles from './styles.module.less'

const ClusterCertificate = () => {
  const { t } = useTranslation()
  const [form] = Form.useForm()
  const navigate = useNavigate()
  const location = useLocation()
  const { cluster } = queryString.parse(location?.search)
  const { isReadOnlyMode } = useSelector((state: any) => state.globalSlice)
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
      message.error(response?.message || t('RequestFailedAndTry'))
    }
  }

  useEffect(() => {
    if (cluster) {
      getClusterDetail()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [cluster])

  async function onFinish() {
    if (isReadOnlyMode) return
    if (!newYamlContent?.content) {
      message.warning(t('PleaseUploadNewKubeConfigFile'))
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
          message.success(t('SubmitAnd3STOClusterPage'))
          setTimeout(() => {
            navigate('/cluster')
          }, 3000)
        } else {
          message.error(response?.message || '验证成功但提交失败')
          setLoading(false)
        }
      } else {
        message.error(
          validateResponse?.message ||
            t('KubeConfigDoesNotMeetTheRequirements'),
        )
        setLoading(false)
      }
    }
  }

  function goBack() {
    navigate('/cluster')
  }

  const uploadProps: any = {
    disabled: isReadOnlyMode,
    name: 'file',
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
          message.success(`${info.file.name} ${t('UploadSuccessful')}`)
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
              `${t('TheFileMustBeIn')}.yaml, .yml, .json, .kubeconfig, .kubeconf`,
          )
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
        <div className={styles.header_back} onClick={goBack}>
          <ArrowLeftOutlined style={{ fontSize: 18 }} />
        </div>
        <h4 className={styles.page_title}>{t('RotateCertificate')}</h4>
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
            rules={[
              {
                required: true,
                message: t('TheKubeConfigFileCannotBeEmpty'),
              },
            ]}
          >
            <Upload {...uploadProps}>
              <Button disabled={isReadOnlyMode} icon={<UploadOutlined />}>
                {t('Upload')} KubeConfig {t('ConfigurationFile')}
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
                {loading ? t('SubmitAndValidate') : t('SubmitAndUpdate')}
              </Button>
              <Button htmlType="button" onClick={onCancel}>
                {t('Cancel')}
              </Button>
            </Space>
          </Form.Item>
        </Form>
        {newYamlContent?.sanitizedClusterContent ? (
          <div className={styles.config_content}>
            <div className={styles.diff_container}>
              <ReactDiffViewer
                leftTitle={t('ExistingConfigurations')}
                rightTitle={t('NewConfiguration')}
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
            <div className={styles.title}>{t('ExistingConfigurations')}</div>
            <Yaml data={lastYamlContent} height="100%" />
          </div>
        )}
      </div>
    </div>
  )
}

export default ClusterCertificate
