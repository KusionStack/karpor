import React, { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Button, Input, Modal, Form, Space, Popconfirm, Select } from 'antd'
import { MinusCircleOutlined, PlusOutlined } from '@ant-design/icons'

interface Field {
  key: 'namespace' | 'annotations' | 'labels'
  value: string | undefined
}

interface RuleObject {
  name: string
  desc?: string
  fields: Field[]
}

type IProps = {
  isEdit: boolean
  open: boolean
  onClose?: () => void
  title?: string
  handleSubmit: (values: RuleObject, callback: () => void) => void
  currentItem: any
  deleteItem: (key: string, callback: () => void) => void
}

const RuleForm = ({
  onClose,
  open,
  title,
  handleSubmit,
  isEdit,
  currentItem,
  deleteItem,
}: IProps) => {
  const [form] = Form.useForm()
  const { t } = useTranslation()
  const [disabledInputs, setDisabledInputs] = useState({})

  useEffect(() => {
    if (currentItem) {
      const obj = {}
      currentItem?.fields?.map((item, index) => {
        obj[index] = item === 'namespace'
      })
      setDisabledInputs(obj)
      form.setFieldsValue({
        name: currentItem?.name,
        description: currentItem?.description,
        fields: currentItem?.fields?.map(item => {
          if (item === 'namespace') {
            return { key: 'namespace', value: undefined }
          } else {
            const [key, ...values] = item?.split('.')
            return {
              key,
              value: values?.join('.'),
            }
          }
        }),
      })
    }
  }, [currentItem, form])

  async function handleOk() {
    try {
      const values = await form.validateFields()
      handleSubmit(values, () => {
        form.resetFields()
        setDisabledInputs({})
      })
    } catch (error) {}
  }

  async function confirm() {
    deleteItem(currentItem?.name, () => {
      form.resetFields()
      setDisabledInputs({})
    })
  }

  const groupByOptoins = [
    {
      label: 'namespace',
      value: 'namespace',
    },
    {
      label: 'annotations',
      value: 'annotations',
    },
    {
      label: 'labels',
      value: 'labels',
    },
  ]

  function onCancel() {
    form.resetFields()
    setDisabledInputs({})
    onClose()
  }

  return (
    <Modal
      width={720}
      title={title}
      onOk={handleOk}
      onCancel={onCancel}
      open={open}
      destroyOnClose
      okText={t('Submit')}
      cancelText={t('Close')}
      footer={(_, { OkBtn, CancelBtn }) => (
        <div style={{ textAlign: 'left' }}>
          {isEdit ? (
            <Popconfirm
              title={t('Delete')}
              description={t('AreYouSureDeleteResourceGroupRule')}
              onConfirm={confirm}
              getPopupContainer={triggerNode => {
                return triggerNode.parentNode as HTMLElement
              }}
            >
              <Button danger style={{ background: 'red', color: '#fff' }}>
                {t('Delete')}
              </Button>
            </Popconfirm>
          ) : (
            <CancelBtn />
          )}
          <span style={{ marginRight: 15 }}></span>
          <OkBtn />
        </div>
      )}
    >
      <Form
        layout="vertical"
        form={form}
        name="form_in_drawer"
        initialValues={{
          fields: [
            {
              key: '',
              value: '',
            },
          ],
        }}
      >
        <Form.Item
          name="name"
          label={t('Name')}
          rules={[
            {
              required: true,
              message: 'Please input the title of collection!',
            },
          ]}
        >
          <Input disabled={isEdit} />
        </Form.Item>
        <Form.Item name="description" label={t('Description')}>
          <Input.TextArea rows={2} />
        </Form.Item>
        <Form.List name="fields">
          {(fields, { add, remove }) => {
            return (
              <>
                {fields.map(({ key, name, ...restField }, index) => {
                  return (
                    <Form.Item
                      label={index === 0 ? 'Group By' : ''}
                      required={true}
                      key={key}
                      style={{ marginBottom: 0 }}
                    >
                      <Space
                        style={{ display: 'flex', marginBottom: 8 }}
                        align="baseline"
                      >
                        <Form.Item
                          {...restField}
                          name={[name, 'key']}
                          rules={[
                            {
                              required: false,
                              message: 'Missing first name',
                            },
                          ]}
                        >
                          <Select
                            options={groupByOptoins}
                            placeholder=""
                            style={{ width: 320 }}
                            onChange={value => {
                              if (value === 'namespace') {
                                const newItems = [
                                  ...form.getFieldValue('fields'),
                                ]
                                newItems[name] = {
                                  ...newItems[name],
                                  value: undefined,
                                }
                                form.setFieldsValue({ fields: newItems })
                              }
                              setDisabledInputs(prev => ({
                                ...prev,
                                [name]: value === 'namespace',
                              }))
                            }}
                          />
                        </Form.Item>
                        <Form.Item
                          {...restField}
                          name={[name, 'value']}
                          shouldUpdate={(prevValues: any, curValues: any) =>
                            prevValues.fields?.[name]?.key !==
                            curValues.fields?.[name]?.key
                          }
                          rules={[{ required: false, message: 'Required' }]}
                        >
                          <Input
                            placeholder=""
                            style={{ width: 320 }}
                            disabled={disabledInputs[name]}
                          />
                        </Form.Item>
                        {fields?.length > 1 && (
                          <MinusCircleOutlined onClick={() => remove(name)} />
                        )}
                      </Space>
                    </Form.Item>
                  )
                })}
                <Form.Item>
                  <Button
                    type="dashed"
                    onClick={() => add()}
                    block
                    icon={<PlusOutlined />}
                  >
                    Add
                  </Button>
                </Form.Item>
              </>
            )
          }}
        </Form.List>
      </Form>
    </Modal>
  )
}

export default RuleForm
