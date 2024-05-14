import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { ArrowLeftOutlined, UploadOutlined } from '@ant-design/icons'
import { Form, Input, Space, Button, Upload, message } from 'antd'
import type { UploadProps } from 'antd'
import { useTranslation } from 'react-i18next'
import axios from 'axios'
import { useSelector } from 'react-redux'
import { HOST } from '@/utils/request'
import Yaml from '@/components/yaml'

import styles from './styles.module.less'

const { TextArea } = Input

const RegisterCluster = () => {
  const { isReadOnlyMode } = useSelector((state: any) => state.globalSlice)
  const [form] = Form.useForm()
  const navigate = useNavigate()
  const [yamlContent, setYamlContent] = useState('')
  const [loading, setLoading] = useState(false)
  const { t } = useTranslation()

  async function onFinish(values: any) {
    if (isReadOnlyMode) {
      return
    }
    setLoading(true)
    const validateResponse: any = await axios.post(
      '/rest-api/v1/cluster/config/validate',
      {
        kubeConfig: values?.kubeConfig,
      },
    )
    if (validateResponse?.success) {
      const response: any = await axios({
        url: `/rest-api/v1/cluster/${values?.name}`,
        method: 'POST',
        data: values,
      })
      if (response?.success) {
        message.success(t('VerifiedSuccessfullyAndSubmitted'))
        navigate(-1)
      } else {
        message.error(
          response?.message || t('VerificationSuccessfulButSubmissionFailed'),
        )
      }
    } else {
      message.error(
        validateResponse?.message || t('KubeConfigDoesNotMeetTheRequirements'),
      )
      setLoading(false)
    }
  }

  function goBack() {
    navigate(-1)
  }

  const uploadProps: UploadProps = {
    disabled: isReadOnlyMode,
    name: 'file',
    action: `${HOST}/rest-api/v1/cluster/config/file`,
    headers: {
      authorization: 'authorization-text',
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
      if (info.file.status === 'done') {
        if (info?.file?.response?.success) {
          message.success(`${info.file.name} ${t('UploadSuccessful')}`)
          form.setFieldsValue({
            kubeConfig: info?.file?.response?.data?.content,
          })
          setYamlContent(info?.file?.response?.data?.content)
        } else {
          message.error(
            info?.file?.response?.message ||
              `${t('TheFileMustBeIn')}.yaml, .yml, .json, .kubeConfig, .kubeconf`,
          )
        }
      } else if (info.file.status === 'error') {
      }
    },
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <ArrowLeftOutlined style={{ marginRight: 10 }} onClick={goBack} />
        {t('RegisterCluster')}
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
          <Form.Item
            name="name"
            label={t('ClusterName')}
            rules={[{ required: true }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="displayName"
            label={t('DisplayName')}
            rules={[{ required: false }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="description"
            label={t('Description')}
            rules={[{ required: false }]}
          >
            <TextArea autoSize={{ minRows: 3 }} />
          </Form.Item>
          <Form.Item
            label="kubeConfig"
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
          <Form.Item>
            <Space>
              <Button
                disabled={isReadOnlyMode}
                type="primary"
                htmlType="submit"
                loading={loading}
              >
                {t('VerifyAndSubmit')}
              </Button>
              <Button htmlType="button" onClick={goBack}>
                {t('Cancel')}
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

export default RegisterCluster
