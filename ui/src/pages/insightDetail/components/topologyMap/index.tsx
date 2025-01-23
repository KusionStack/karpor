import React, { forwardRef, useImperativeHandle, useRef } from 'react'
import { Select } from 'antd'
import { Circle, Rect as GRect, Text as GText, Image } from '@antv/g'
import {
  Graph,
  NodeEvent,
  ExtensionCategory,
  register,
  Rect,
  Label,
  Badge,
  CommonEvent,
  CubicHorizontal,
  subStyleProps,
  LabelStyleProps,
  IPointerEvent,
} from '@antv/g6'
import { useLocation, useNavigate } from 'react-router-dom'
import queryString from 'query-string'
import { useTranslation } from 'react-i18next'
import Loading from '@/components/loading'
import transferImg from '@/assets/transfer.png'
import { ICON_MAP } from '@/utils/images'

import styles from './style.module.less'

function getTextWidth(str: string, fontSize: number) {
  const canvas = document.createElement('canvas')
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  const context = canvas.getContext('2d')!
  context.font = `${fontSize}px sans-serif`
  return context.measureText(str).width
}

export function fittingString(str: string, maxWidth: number, fontSize: number) {
  const ellipsis = '...'
  const ellipsisLength = getTextWidth(ellipsis, fontSize)

  if (maxWidth <= 0) {
    return ''
  }

  const width = getTextWidth(str, fontSize)
  if (width <= maxWidth) {
    return str
  }

  let len = str.length
  while (len > 0) {
    const substr = str.substring(0, len)
    const subWidth = getTextWidth(substr, fontSize)

    if (subWidth + ellipsisLength <= maxWidth) {
      return substr + ellipsis
    }

    len--
  }

  return str
}

export function getNodeName(cfg: any, type: string) {
  if (type === 'resource') {
    const [left, right] = cfg?.id?.split(':') || []
    const leftList = left?.split('.')
    const leftListLength = leftList?.length || 0
    const leftLast = leftList?.[leftListLength - 1]
    return `${leftLast}:${right}`
  }
  const list = cfg?.label?.split('.')
  const len = list?.length || 0
  return list?.[len - 1] || ''
}

class CardNode extends Rect {
  get data() {
    return this.context.model.getNodeLikeDatum(this.id)
  }

  get childrenData() {
    return this.context.model.getChildrenData(this.id)
  }

  getIconStyle(attributes): any {
    const [width, height]: any = this.getSize(attributes)
    return {
      x: -width / 2 + 8,
      y: -height / 2 + 8,
      width: 32,
      height: 32,
      src:
        ICON_MAP?.[(this.data?.data?.resourceGroup as any)?.kind as any] ||
        transferImg,
      fontSize: 16,
      fill: '#000',
      opacity: 0.85,
    }
  }

  drawIconShape(attributes, container) {
    const iconStyle: any = this.getIconStyle(attributes)
    this.upsert('data', Image, iconStyle, container)
  }

  getLabelStyle(attributes): LabelStyleProps {
    const [width, height]: any = this.getSize(attributes)
    const text: any = this.data.label
    return {
      x: -width / 2 + 45,
      y: -height / 2 + 16,
      text: text,
      fontSize: 14,
      opacity: 0.85,
      fill: '#1677ff',
      fontWeight: 600,
      cursor: 'pointer',
    }
  }

  getTypeStyle(attributes) {
    const [width, height]: any = this.getSize(attributes)
    return {
      x: -width / 2 + 45,
      y: -height / 2 + 40,
      text: (this.data?.data?.resourceGroup as any)?.kind,
      fontSize: 12,
      fill: '#999',
      opacity: 0.85,
    }
  }

  drawTypeShape(attributes, container) {
    const typeStyle: any = this.getTypeStyle(attributes)
    this.upsert('price', GText, typeStyle, container)
  }

  getLeftBarStyle(attributes) {
    const { radius } = attributes
    const color = '#1677ff'
    const [width, height]: any = this.getSize(attributes)
    return {
      x: -width / 2,
      y: -height / 2,
      width: 3,
      height: height,
      radius: [radius, 0, 0, radius],
      // fill: `linear-gradient(to bottom, ${'#fff'} ${'0%'},${color} ${'50%'}, ${'#fff'} ${'100%'})`,
      fill: color,
      opacity: 0.3,
    }
  }

  drawLeftBarShape(attributes, container) {
    const leftBarStyle = this.getLeftBarStyle(attributes)
    this.upsert('process-bar', GRect, leftBarStyle, container)
  }

  getCountStyle(attributes) {
    const color = '#1677ff'
    const [width, height]: any = this.getSize(attributes)
    const circleHeight = 24
    return {
      x: width / 2 - 28,
      y: -height / 2 + 12,
      width: circleHeight,
      height: circleHeight,
      radius: [
        circleHeight / 2,
        circleHeight / 2,
        circleHeight / 2,
        circleHeight / 2,
      ],
      fill: color,
      opacity: 0.2,
    }
  }

  drawCountShape(attributes, container) {
    const countStyle = this.getCountStyle(attributes)
    this.upsert('countCircle', GRect, countStyle, container)
  }

  getCountTextStyle(attributes) {
    const color = '#1677ff'
    const [width, height]: any = this.getSize(attributes)
    return {
      x: width / 2 - 26,
      y: -height / 2 + 34,
      text: `${this.data?.data?.count}`,
      // text: (this.data?.data?.resourceGroup as any)?.kind,
      fontSize: 12,
      fill: color,
      opacity: 0.85,
    }
  }

  drawCountTextShape(attributes, container) {
    const countStyle: any = this.getCountTextStyle(attributes)
    this.upsert('count', GText, countStyle, container)
  }

  getKeyStyle(attributes) {
    const keyStyle = super.getKeyStyle(attributes)
    return {
      ...keyStyle,
      fill: '#fff',
      lineWidth: 1,
      stroke: '#fff',
      shadowColor: 'rgba(0,0,0,0.05)',
      shadowBlur: 4,
      shadowOffsetX: 0,
      shadowOffsetY: 2,
    }
  }

  render(attributes = this.parsedAttributes, container) {
    console.log(this.data, '=====render this.data====')

    super.render(attributes, container)
    this.drawIconShape(attributes, container)
    this.drawTypeShape(attributes, container)
    this.drawLeftBarShape(attributes, container)
    if (typeof this.data?.data?.count === 'number') {
      this.drawCountShape(attributes, container)
      this.drawCountTextShape(attributes, container)
    }
  }
}

class FlyMarkerCubic extends CubicHorizontal {
  getMarkerStyle(attributes) {
    return {
      r: 2,
      fill: '#1677ff',
      opacity: 0.7,
      offsetPath: this.shapeMap.key,
      ...subStyleProps(attributes, 'marker'),
    }
  }

  onCreate() {
    const marker = this.upsert(
      'marker',
      Circle,
      this.getMarkerStyle(this.attributes),
      this,
    )
    marker.animate([{ offsetDistance: 0 }, { offsetDistance: 1 }], {
      duration: 3000,
      iterations: Infinity,
    })
  }
}

register(ExtensionCategory.NODE, 'card-node', CardNode)
register(ExtensionCategory.EDGE, 'running-cubic', FlyMarkerCubic)

export interface NodeModel {
  id: string
  name?: string
  label?: string
  resourceGroup?: {
    name: string
  }
  data?: {
    count?: number
    resourceGroup?: {
      name: string
    }
  }
}

type IProps = {
  topologyData?: any
  topologyLoading?: boolean
  onTopologyNodeClick?: (node: any) => void
  isResource?: boolean
  tableName?: string
  handleChangeCluster?: (val: any) => void
  selectedCluster?: string
  clusterOptions?: string[]
}

const TopologyMap = forwardRef((props: IProps, drawRef) => {
  const {
    onTopologyNodeClick,
    topologyLoading,
    isResource,
    selectedCluster,
    clusterOptions,
    handleChangeCluster,
  } = props
  const { t } = useTranslation()
  const containerRef = useRef(null)
  const graphRef = useRef<any>(null)
  const location = useLocation()
  const { from, type, query } = queryString.parse(location?.search)
  const navigate = useNavigate()
  console.log(from, query, navigate, '==query====')

  function handleMouseEnter(evt) {
    graphRef.current?.setElementState(evt.item, 'active', true)
  }

  const handleMouseLeave = (evt: any) => {
    graphRef.current?.setElementState(evt.item, 'active', false)
  }

  function initGraph(data) {
    const container = containerRef.current
    const width = container?.scrollWidth
    const height = container?.scrollHeight
    return new Graph({
      container,
      width,
      height,
      data,
      padding: [20, 20, 20, 20],
      autoFit: {
        type: 'center',
      },
      plugins: [
        {
          type: 'tooltip',
          getContent: (e, items) => {
            console.log(items, '===items===')
            let result = `<h4>Custom Content</h4>`
            items.forEach(item => {
              result += `<div>Type: ${item.data.resourceGroup?.kind}</div>`
            })
            return result
          },
        },
      ],
      behaviors: [
        'hover-activate',
        'drag-canvas',
        'drag-element',
        'click-select',
      ],
      layout: {
        type: 'dagre',
        rankdir: 'LR',
        align: 'UL',
        // nodesep: 10,
        // ranksep: 40,
        nodesepFunc: () => 1,
        ranksepFunc: () => 1,
        controlPoints: true,
        sortByCombo: false,
        preventOverlap: true,
        nodeSize: [200, 60],
        workerEnabled: true,
        clustering: false,
        clusterNodeSize: [200, 60],
        // Optimize edge layout
        edgeFeedbackStyle: {
          stroke: '#c2c8d1',
          lineWidth: 1,
          strokeOpacity: 0.5,
          endArrow: true,
        },
      },
      node: {
        type: 'card-node',
        style: {
          size: [200, 48],
          fill: '#fff',
          stroke: '#fff',
          radius: 8,
          shadowColor: 'rgba(0,0,0,0.05)',
          cursor: 'pointer',
        },
        state: {
          active: {
            fill: '#f50',
          },
          selected: {
            fill: '#1677ff',
          },
        },
      },
      edge: {
        type: 'running-cubic',
        style: {
          radius: 10,
          offset: 5,
          endArrow: true,
          lineWidth: 1,
          curveness: 0.5,
        },
      },
    })
  }

  function drawGraph(topologyData) {
    if (topologyData) {
      if (type === 'resource') {
        graphRef.current?.destroy()
        graphRef.current = null
      }
      if (!graphRef.current) {
        graphRef.current = initGraph(topologyData)
        graphRef.current?.render()

        graphRef.current?.on(NodeEvent.CLICK, evt => {
          const node = evt.item
          const model = node.getModel()
          graphRef.current?.getNodes().forEach(n => {
            graphRef.current?.setItemState(n, 'selected', false)
          })
          graphRef.current?.setItemState(node, 'selected', true)
          onTopologyNodeClick?.(model)
        })

        graphRef.current?.on(NodeEvent.POINTER_ENTER, evt => {
          console.log(evt, '===evt=====')
          // const node = evt.item
          // return
          // if (
          //   !graphRef.current
          //     ?.findById(node.getModel().id)
          //     ?.hasState('selected')
          // ) {
          //   graphRef.current?.setItemState(node, 'active', true)
          // }
          graphRef.current.updateNodeData([
            {
              id: evt?.target.id,
              style: {
                labelText: 'Hovered',
                fill: 'lightgreen',
                labelFill: 'lightgreen',
              },
            },
          ])
          graphRef.current.draw()
          // handleMouseEnter(evt)
        })

        graphRef.current?.on(NodeEvent.POINTER_LEAVE, evt => {
          handleMouseLeave(evt)
        })

        if (typeof window !== 'undefined') {
          window.onresize = () => {
            console.log(graphRef.current, '==graphRef.current=')
            if (!graphRef.current) return
            if (
              !containerRef ||
              !containerRef.current?.scrollWidth ||
              !containerRef.current?.scrollHeight
            )
              return
            graphRef.current?.setSize(
              containerRef?.current?.scrollWidth,
              containerRef.current?.scrollHeight,
            )
          }
        }
      } else {
        graphRef.current.setData(topologyData)
        graphRef.current.render()
        setTimeout(() => {
          graphRef.current.autoFit()
        }, 100)
      }
    }
  }

  useImperativeHandle(drawRef, () => ({
    drawGraph,
  }))

  return (
    <div
      className={styles.g6_topology}
      style={{ height: isResource ? 450 : 400 }}
    >
      <div className={styles.cluster_select}>
        <Select
          style={{ minWidth: 100 }}
          placeholder=""
          value={selectedCluster}
          onChange={handleChangeCluster}
        >
          {clusterOptions?.map(item => {
            return (
              <Select.Option key={item}>
                {item === 'ALL' ? t('AllClusters') : item}
              </Select.Option>
            )
          })}
        </Select>
      </div>
      <div ref={containerRef} className={styles.g6_overview}>
        <div
          className={styles.g6_loading}
          style={{ display: topologyLoading ? 'block' : 'none' }}
        >
          <Loading />
        </div>
      </div>
    </div>
  )
})

export default TopologyMap
