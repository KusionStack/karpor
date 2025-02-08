import React from 'react'
import { Button, Modal, Tag } from 'antd'
import { useTranslation } from 'react-i18next'
import { SEVERITY_MAP } from '@/utils/constants'
import TagVariableSizeList from '../tagVariableSizeList'
import { IssueGroup } from '../../types'

import styles from './style.module.less'

type IProps = {
  open: boolean
  detail: IssueGroup
  cancel: () => void
}

const EventDetail = ({ open, detail, cancel }: IProps) => {
  const resourceGroupsNames = detail?.resourceGroups?.map(item => {
    return {
      ...item,
      allName: `${item?.cluster || ''} ${item?.apiVersion || ''} ${item?.kind || ''} ${item?.namespace || ''} ${item?.name || ''} `,
    }
  })
  const { t } = useTranslation()
  return (
    <Modal
      centered
      title={t('IssuesDetail')}
      open={open}
      destroyOnClose
      maskClosable
      onCancel={cancel}
      footer={[
        <Button key="closebBtn" onClick={cancel}>
          {t('Close')}
        </Button>,
      ]}
    >
      <div className={styles.container}>
        <div className={styles.title}>
          <Tag color={SEVERITY_MAP?.[detail?.issue?.severity]?.color}>
            {t(SEVERITY_MAP?.[detail?.issue?.severity]?.text)}
          </Tag>
          {detail?.issue?.title || '--'}
        </div>
        <div className={styles.content}>
          <div className={styles.desc}>
            <div className={styles.item}>
              <div className={styles.label}>{t('IssueSource')}:&nbsp;</div>
              <div className={styles.value}>{detail?.issue?.scanner}</div>
            </div>
            <div className={styles.item}>
              <div className={styles.label}>
                {t('NumberOfOccurrences')}:&nbsp;
              </div>
              <div className={styles.value}>
                {detail?.resourceGroups?.length}
              </div>
            </div>
            <div
              className={styles.item}
              style={{ width: '100%', alignItems: 'baseline' }}
            >
              <div className={styles.label}>{t('Description')}:&nbsp;</div>
              <div className={styles.value}>
                <div className={styles.value}>{detail?.issue?.message}</div>
              </div>
            </div>
          </div>
          <div className={styles.footer}>
            <div className={styles.footer_title}>
              {t('RelatedResources')}:&nbsp;
            </div>
            <div className={styles.soultion}>
              <TagVariableSizeList
                allTags={resourceGroupsNames}
                containerWidth={480}
              />
            </div>
          </div>
        </div>
      </div>
    </Modal>
  )
}

export default EventDetail
