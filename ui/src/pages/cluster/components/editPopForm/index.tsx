import React, { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Popover, Form, Input, Button, Space } from 'antd'

const { TextArea } = Input

type EditFormProps = {
  submit: (val: any, callback: () => void) => void
  cancel: () => void
  lastDetail: any
}

export const EditForm = ({ submit, cancel, lastDetail }: EditFormProps) => {
  const [form] = Form.useForm()
  const { t } = useTranslation()

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
    <div style={{ width: 320, padding: '20px 20px 0 20px' }}>
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
        <Form.Item
          name="description"
          label={t('Description')}
          rules={[{ required: false }]}
        >
          <TextArea autoSize={{ minRows: 3 }} />
        </Form.Item>
        <Form.Item style={{ textAlign: 'right' }}>
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
    </div>
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
  title,
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
    submit,
    cancel: hide,
    lastDetail,
  }
  function handleOpenChange(v) {
    setOpen(v)
    setLastDetail(lastDetail)
  }
  return (
    <Popover
      open={open}
      onOpenChange={handleOpenChange}
      placement="bottomRight"
      title={title || t('Edit')}
      trigger="click"
      content={<EditForm {...formProps} />}
    >
      <Button disabled={isDisabled} type={btnType || 'default'}>
        {t('Edit')}
      </Button>
    </Popover>
  )
}

export default EditPopForm
