import React from 'react'
import { useTranslation } from 'react-i18next'
import styles from './styles.module.less'

const Reflux = () => {
  const { t } = useTranslation()
  return <div className={styles.container}>{t('InDevelopment')}...</div>
}

export default Reflux
