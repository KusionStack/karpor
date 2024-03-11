import React, { useEffect } from 'react'
import { Gauge } from '@antv/g2plot'
import BigNumber from 'bignumber.js'
import { useTranslation } from 'react-i18next'

type IProps = {
  data: number | string
}

const GaugeChart = ({ data }: IProps) => {
  const { t } = useTranslation()
  useEffect(() => {
    if (!data) {
      return
    }
    const container = document.getElementById('gaugeConatiner')
    const numData = Number(data)
    const color =
      numData < 0.6 ? '#F4664A' : numData < 0.8 ? '#FAAD14' : '#30BF78'
    const gauge = new Gauge(container, {
      percent: numData,
      range: {
        width: 12,
        ticks: [0, 3 / 5, 4 / 5, 1],
        color: ['#F4664A', '#FAAD14', '#30BF78'],
      },
      gaugeStyle: {
        // lineCap: 'round',
      },
      indicator: {
        pointer: {
          style: {
            stroke: color,
          },
        },
        pin: {
          style: {
            stroke: color,
          },
        },
      },
      axis: {
        label: {
          formatter(v) {
            return Number(v) * 100
          },
        },
        subTickLine: {
          count: 10,
        },
      },
      statistic: {
        title: {
          offsetY: -10,
          formatter: () => t('HealthScore'),
          style: {
            color: color,
            fontSize: '16px',
            fontWeight: 'bold',
          },
        },
        content: {
          offsetY: 10,
          formatter: ({ percent }) => `${new BigNumber(percent).times(100)}`,
          style: {
            color: color,
            fontSize: '16px',
            fontWeight: 'bold',
          },
        },
      },
    })

    gauge.render()

    return () => {
      if (gauge) {
        gauge.destroy()
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [data])

  return <div style={{ width: '100%', height: 150 }} id="gaugeConatiner"></div>
}

export default GaugeChart
