import React from 'react'
import { Button, Popconfirm } from 'antd'
import { useSelector } from 'react-redux'
import { useTranslation } from 'react-i18next'
import { utcDateToLocalDate } from '@/utils/tools'
import k8sPng from '@/assets/kubernetes.png'
import EditPopForm from '../editPopForm'

import styles from './styles.module.less'

type IProps = {
  item: any
  goDetailPage: (val) => void
  deleteItem: (val) => void
  goCertificate: (val) => void
  setLastDetail: (val) => void
  handleSubmit: (values: any, callback: () => void) => void
  customStyle: React.CSSProperties
}

const ClusterCard = (props: IProps) => {
  const { t } = useTranslation()
  const { isReadOnlyMode } = useSelector((state: any) => state.globalSlice)
  const {
    goDetailPage,
    item,
    deleteItem,
    goCertificate,
    setLastDetail,
    handleSubmit,
    customStyle,
  } = props
  return (
    <div className={styles.card} style={customStyle}>
      <div className={styles.left} onClick={() => goDetailPage(item)}>
        <div className={styles.score}>
          <img src={k8sPng} alt="icon" />
        </div>
        <div className={styles.detail}>
          <div className={styles.top}>
            <div className={styles.name}>
              {item?.spec?.displayName ? (
                <span>
                  {item?.spec?.displayName}
                  <span style={{ color: '#808080' }}>
                    （{item?.metadata?.name}）
                  </span>
                </span>
              ) : (
                <span>{item?.metadata?.name}</span>
              )}
            </div>
          </div>
          <div className={styles.desc}>{item?.spec?.description || '--'}</div>
          <div className={styles.bottom}>
            {item?.metadata?.creationTimestamp
              ? utcDateToLocalDate(item?.metadata?.creationTimestamp)
              : '--'}
          </div>
        </div>
      </div>
      <div className={styles.right}>
        <EditPopForm
          isDisabled={isReadOnlyMode}
          submit={handleSubmit}
          lastDetail={item}
          setLastDetail={setLastDetail}
        />
        <Button
          disabled={isReadOnlyMode}
          style={{ margin: '0 16px' }}
          onClick={() => goCertificate(item)}
        >
          {t('RotateCertificate')}
        </Button>
        <Popconfirm
          disabled={isReadOnlyMode}
          placement="topLeft"
          title={
            <span style={{ display: 'inline-block', width: 200 }}>
              {t('DeleteAndNoLongUpdateResources')}
            </span>
          }
          description=""
          onConfirm={() => deleteItem(item)}
        >
          <Button disabled={isReadOnlyMode}>{t('Delete')}</Button>
        </Popconfirm>
      </div>
    </div>
  )
}

export default ClusterCard
