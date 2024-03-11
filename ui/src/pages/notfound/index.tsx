import React from 'react'
import { Button, Result } from 'antd'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'

import styles from './styles.module.less'

const NotFound = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  function goBack() {
    navigate('/search')
  }

  return (
    <div className={styles.container}>
      <Result
        status="404"
        title="404"
        subTitle={t('SorryThePageYouVisitedDoesNotExist')}
        extra={
          <Button type="primary" onClick={goBack}>
            {t('BackToHome')}
          </Button>
        }
      />
    </div>
  )
}

export default NotFound
