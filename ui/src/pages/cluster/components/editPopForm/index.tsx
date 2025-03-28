import React, { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Modal, Form, Input, Button, Space, Select } from 'antd'
import { useSelector } from 'react-redux'

const { TextArea } = Input

const { Option } = Select

type EditFormProps = {
  submit: (val: any, callback: () => void) => void
  cancel: () => void
  lastDetail: any
  open: boolean
}

export const EditForm = ({
  submit,
  cancel,
  lastDetail,
  open,
}: EditFormProps) => {
  const [form] = Form.useForm()
  const { t } = useTranslation()

  const { isHighAvailability } = useSelector((state: any) => state.globalSlice)

  useEffect(() => {
    if (lastDetail) {
      form.setFieldsValue({
        name: lastDetail?.metadata?.name,
        displayName: lastDetail?.spec?.displayName,
        description: lastDetail?.spec?.description,
      })
    }
  }, [lastDetail, form])

  function onFinish(values: any) {
    submit(values, () => {
      onCancel()
    })
  }

  function onCancel() {
    cancel()
  }

  return (
    <Modal open={open} footer={null} centered onCancel={onCancel}>
      <Form form={form} layout="vertical" onFinish={onFinish}>
        <Form.Item
          name="name"
          label={t('ClusterName')}
          rules={[{ required: true }]}
        >
          <Input disabled />
        </Form.Item>
        <Form.Item
          name="displayName"
          label={t('DisplayName')}
          rules={[{ required: false }]}
        >
          <Input />
        </Form.Item>
        {isHighAvailability && (
          <>
            <Form.Item
              name="clusterMode"
              label={t('ClusterMode')}
              rules={[{ required: true }]}
            >
              <Select>
                <Option key="pull" value="pull">
                  pull
                </Option>
                <Option key="push" value="push">
                  push
                </Option>
              </Select>
            </Form.Item>
            <Form.Item
              name="clusterLevel"
              label={t('ClusterLevel')}
              rules={[{ required: true }]}
            >
              <Select>
                <Option key="1" value={1}>
                  1
                </Option>
                <Option key="2" value={2}>
                  2
                </Option>
                <Option key="3" value={3}>
                  3
                </Option>
              </Select>
            </Form.Item>
          </>
        )}
        <Form.Item
          name="description"
          label={t('Description')}
          rules={[{ required: false }]}
        >
          <TextArea autoSize={{ minRows: 3 }} />
        </Form.Item>
        <Form.Item style={{ textAlign: 'right', marginBottom: 0 }}>
          <Space>
            <Button htmlType="button" onClick={onCancel}>
              {t('Cancel')}
            </Button>
            <Button type="primary" htmlType="submit">
              {t('OK')}
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  )
}

type IProps = {
  setLastDetail: (value) => void
  title?: string | React.ReactNode
  submit: (values: any, callback: () => void) => void
  btnType?: 'dashed' | 'link' | 'text' | 'default' | 'primary'
  lastDetail?: any
  isDisabled?: boolean
}

const EditPopForm = ({
  setLastDetail,
  submit,
  btnType,
  lastDetail,
  isDisabled,
}: IProps) => {
  const { t } = useTranslation()
  const [open, setOpen] = useState(false)
  const hide = () => {
    setOpen(false)
  }
  const formProps = {
    open,
    submit,
    cancel: hide,
    lastDetail,
  }
  function handleOpenChange(v) {
    setOpen(v)
    setLastDetail(lastDetail)
  }
  return (
    <>
      <Button
        disabled={isDisabled}
        type={btnType || 'default'}
        onClick={handleOpenChange}
      >
        {t('Edit')}
      </Button>
      <EditForm {...formProps} />
    </>
  )
}

export default EditPopForm
