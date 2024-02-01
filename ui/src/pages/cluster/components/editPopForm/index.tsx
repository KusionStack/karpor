import React, { useEffect, useState } from 'react'
import { Popover, Form, Input, Button, Space } from 'antd'

const { TextArea } = Input

type EditFormProps = {
  submit: (val: any, callback: () => void) => void
  cancel: () => void
  lastDetail: any
}

export const EditForm = ({ submit, cancel, lastDetail }: EditFormProps) => {
  const [form] = Form.useForm()

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
    form.resetFields()
    cancel()
  }

  return (
    <div style={{ width: 320, padding: '20px 20px 0 20px' }}>
      <Form form={form} layout="vertical" onFinish={onFinish}>
        <Form.Item name="name" label="集群名称" rules={[{ required: true }]}>
          <Input disabled />
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
        <Form.Item style={{ textAlign: 'right' }}>
          <Space>
            <Button htmlType="button" onClick={onCancel}>
              取消
            </Button>
            <Button type="primary" htmlType="submit">
              确定
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
  cancel: () => void
  btnType?: 'dashed' | 'link' | 'text' | 'default' | 'primary'
  btnStyle?: any
  lastDetail?: any
  isDisabled?: boolean
}

const EditPopForm = ({
  setLastDetail,
  title,
  submit,
  cancel,
  btnType,
  btnStyle,
  lastDetail,
  isDisabled,
}: IProps) => {
  const [open, setOpen] = useState(false)
  const hide = () => {
    setOpen(false)
    cancel()
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
      title={title || '编辑'}
      trigger="click"
      content={<EditForm {...formProps} />}
    >
      <Button
        disabled={isDisabled}
        style={btnStyle}
        type={btnType || 'default'}
      >
        编辑
      </Button>
    </Popover>
  )
}

export default EditPopForm
