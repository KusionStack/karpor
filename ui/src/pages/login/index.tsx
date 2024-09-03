import { Button, Input, message, Tooltip } from 'antd'
import React, { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { InfoCircleFilled, QuestionCircleOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { useDispatch } from 'react-redux'
import { setIsLogin } from '@/store/modules/globalSlice'
import { useAxios } from '@/utils/request'

import styles from './styles.module.less'

const Login = () => {
  const { t } = useTranslation()
  const dispatch = useDispatch()
  const navigate = useNavigate()
  const [value, setValue] = useState('')
  const [errorMessage, setErrorMessage] = useState('')

  const { response, refetch } = useAxios({
    url: `/rest-api/v1/authn`,
    option: { params: {} },
    manual: true,
    method: 'GET',
  })

  function handleLogin() {
    if (!value) {
      setErrorMessage(t('InputToken'))
      return
    }
    refetch({
      option: {
        headers: {
          Authorization: `Bearer ${value}`,
        },
        params: {},
      },
    })
  }

  useEffect(() => {
    if (response?.success) {
      message.success(t('LoginSuccess'))
      localStorage.setItem('token', value)
      dispatch(setIsLogin(true))
      setTimeout(() => {
        navigate(-1)
      }, 300)
    } else if (response?.code === 401) {
      setErrorMessage(t('LoginFailedAndCheck'))
    }
  }, [response, dispatch, value, navigate, t])

  return (
    <div className={styles.login_wrapper}>
      <div className={styles.login}>
        <div className={styles.title}>
          <InfoCircleFilled style={{ color: '#2f54eb' }} />
          <h4>{t('UnLoginAndTokenLogin')}</h4>
        </div>
        <div className={styles.content}>
          <div className={styles.label}>
            <Tooltip
              color="#fff"
              title={
                <div
                  style={{
                    color: '#646566',
                  }}
                >
                  <p
                    dangerouslySetInnerHTML={{
                      __html: t('TokenCreationGuide', {
                        linkUrl:
                          'https://www.kusionstack.io/karpor/next/user-guide/how-to-create-token',
                      }),
                    }}
                  />
                </div>
              }
            >
              <span>Token</span>
              <QuestionCircleOutlined
                style={{ margin: '0 5px', cursor: 'pointer' }}
              />
            </Tooltip>
            <span className={styles.token_require}>*</span>
          </div>
          <div className={styles.input_box}>
            <Input
              value={value}
              onChange={event => setValue(event.target.value)}
            />
          </div>
        </div>
        {errorMessage && (
          <div className={styles.error_message}>{errorMessage}</div>
        )}
        <div className={styles.footer}>
          <Button key="login" type="primary" onClick={handleLogin}>
            {t('Login')}
          </Button>
        </div>
      </div>
    </div>
  )
}

export default Login
