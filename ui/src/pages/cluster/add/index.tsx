import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { ArrowLeftOutlined, UploadOutlined } from '@ant-design/icons'
import { Form, Input, Space, Button, Upload, message, notification } from 'antd'
import type { UploadProps } from 'antd'
import { useTranslation } from 'react-i18next'
import { useSelector } from 'react-redux'
import { HOST, useAxios } from '@/utils/request'
import Yaml from '@/components/yaml'
import { fireConfetti } from '@/utils/confetti'

import styles from './styles.module.less'

const { TextArea } = Input

const RegisterCluster = () => {
  const { t } = useTranslation()
  const [form] = Form.useForm()
  const navigate = useNavigate()
  const { isReadOnlyMode, isUnsafeMode } = useSelector(
    (state: any) => state.globalSlice,
  )
  const [yamlContent, setYamlContent] = useState('')

  const {
    response: validateResponse,
    refetch: validateRefetch,
    loading,
  } = useAxios({
    url: '/rest-api/v1/cluster/config/validate',
    manual: true,
    method: 'POST',
  })

  const { response: addResponse, refetch: addRefetch } = useAxios({
    url: '',
    manual: true,
    method: 'POST',
  })

  useEffect(() => {
    if (validateResponse?.success) {
      validateResponse?.callbackFn && validateResponse?.callbackFn()
    }
  }, [validateResponse])

  useEffect(() => {
    if (addResponse?.success) {
      fireConfetti()
      message.success(t('VerifiedSuccessfullyAndSubmitted'))
      setTimeout(() => {
        navigate(-1)
      }, 1000)
    }
  }, [addResponse, navigate, t])

  function onFinish(values: any) {
    if (isReadOnlyMode) {
      return
    }
    validateRefetch({
      option: {
        data: {
          kubeConfig: values?.kubeConfig,
        },
      },
      callbackFn: () => {
        addRefetch({
          url: `/rest-api/v1/cluster/${values?.name}`,
          option: {
            data: values,
          },
        })
      },
    })
  }

  function goBack() {
    navigate(-1)
  }

  const uploadProps: UploadProps = {
    disabled: isReadOnlyMode,
    name: 'file',
    action: `${HOST}/rest-api/v1/cluster/config/file`,
    headers: {
      Authorization: isUnsafeMode
        ? ''
        : localStorage.getItem('token')
          ? `Bearer ${localStorage.getItem('token')}`
          : '',
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
      } else if (
        info.file.status === 'error' &&
        info.file.response?.code === 403
      ) {
        notification.error({
          message: `${info.file.response?.code}`,
          description: `${info.file.response?.message}`,
        })
      }
    },
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <ArrowLeftOutlined style={{ marginRight: 10 }} onClick={goBack} />
        <h4 className={styles.page_title}>{t('RegisterCluster')}</h4>
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
