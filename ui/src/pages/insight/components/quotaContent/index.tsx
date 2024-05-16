import React from 'react'
import insightClusterPng from '@/assets/insight_cluster.png'
import insightResourcesPng from '@/assets/insight_resources.png'
import insightResourceGroupPng from '@/assets/insight_resources_group.png'
import QuotaCard from './quotaCard'

import styles from './styles.module.less'

type IProps = {
  statsData: {
    clusterCount: number
    resourceCount: number
    resourceGroupRuleCount: number
  }
}

const QuotaContent = ({ statsData }: IProps) => {
  return (
    <div className={styles.stat}>
      <div className={styles.item}>
        <QuotaCard
          title="Cluster Count"
          value={statsData?.clusterCount}
          bgColor="#eaf3ed"
          iconNode={
            <img
              src={insightClusterPng}
              style={{
                width: 150,
                height: 150,
                transform: 'rotate(-45deg)',
              }}
            />
          }
        />
      </div>
      <div className={styles.item}>
        <QuotaCard
          title="Total Resources"
          value={statsData?.resourceCount}
          bgColor="#fbf4e7"
          iconNode={
            <img
              src={insightResourcesPng}
              style={{
                width: 150,
                height: 150,
                transform: 'rotate(225deg)',
              }}
            />
          }
        />
      </div>
      <div className={styles.item}>
        <QuotaCard
          title="Total ResourceGroupRules"
          value={statsData?.resourceGroupRuleCount}
          bgColor="#e6f1ff"
          iconNode={
            <img
              src={insightResourceGroupPng}
              style={{
                width: 150,
                height: 150,
                transform: 'rotate(-45deg)',
              }}
            />
          }
        />
      </div>
    </div>
  )
}

export default QuotaContent
