import React, { forwardRef, useImperativeHandle, useRef } from 'react'
import { Select } from 'antd'
import { Circle, Rect as GRect, Text as GText } from '@antv/g'
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
} from '@antv/g6'
import { useLocation, useNavigate } from 'react-router-dom'
import queryString from 'query-string'
import { useTranslation } from 'react-i18next'
import Loading from '@/components/loading'
// import transferImg from '@/assets/transfer.png'
// import { ICON_MAP } from '@/utils/images'

import styles from './style.module.less'

class CardNode extends Rect {
  get data() {
    return this.context.model.getNodeLikeDatum(this.id)
  }

  get childrenData() {
    return this.context.model.getChildrenData(this.id)
  }

  // getLabelStyle(attributes) {
  //   const [width, height]: any = this.getSize(attributes)
  //   return {
  //     x: -width / 2 + 8,
  //     y: -height / 2 + 16,
  //     text: this.data.name,
  //     fontSize: 12,
  //     opacity: 0.85,
  //     fill: '#000',
  //     cursor: 'pointer',
  //   }
  // }

  getPriceStyle(attributes) {
    const [width, height]: any = this.getSize(attributes)
    return {
      x: -width / 2 + 8,
      y: height / 2 - 8,
      text: this.data.label,
      fontSize: 16,
      fill: '#000',
      opacity: 0.85,
    }
  }

  drawPriceShape(attributes, container) {
    const priceStyle: any = this.getPriceStyle(attributes)
    this.upsert('price', GText, priceStyle, container)
  }

  getCurrencyStyle(attributes) {
    const [, height]: any = this.getSize(attributes)
    return {
      x: this.shapeMap['price'].getLocalBounds().max[0] + 4,
      y: height / 2 - 8,
      text: this.data.currency,
      fontSize: 12,
      fill: '#000',
      opacity: 0.75,
    }
  }

  drawCurrencyShape(attributes, container) {
    const currencyStyle: any = this.getCurrencyStyle(attributes)
    this.upsert('currency', GText, currencyStyle, container)
  }

  getPercentStyle(attributes) {
    const [width, height]: any = this.getSize(attributes)
    return {
      x: width / 2 - 4,
      y: height / 2 - 8,
      text: `${((Number(this.data.variableValue) || 0) * 100).toFixed(2)}%`,
      fontSize: 12,
      textAlign: 'right',
      fill: '#f50',
    }
  }

  drawPercentShape(attributes, container) {
    const percentStyle: any = this.getPercentStyle(attributes)
    this.upsert('percent', GText, percentStyle, container)
  }

  getTriangleStyle(attributes) {
    const percentMinX = this.shapeMap['percent'].getLocalBounds().min[0]
    const [, height]: any = this.getSize(attributes)
    return {
      fill: '#1677ff',
      x: this.data.variableUp ? percentMinX - 18 : percentMinX,
      y: height / 2 - 16,
      fontFamily: 'iconfont',
      fontSize: 16,
      text: '\ue62d',
      transform: this.data.variableUp ? [] : [['rotate', 180]],
    }
  }

  drawTriangleShape(attributes, container) {
    const triangleStyle: any = this.getTriangleStyle(attributes)
    this.upsert('triangle', Label, triangleStyle, container)
  }

  getVariableStyle(attributes) {
    const [, height]: any = this.getSize(attributes)
    return {
      fill: '#000',
      fontSize: 16,
      opacity: 0.45,
      text: this.data.variableName,
      textAlign: 'right',
      x: this.shapeMap['triangle'].getLocalBounds().min[0] - 4,
      y: height / 2 - 8,
    }
  }

  drawVariableShape(attributes, container) {
    const variableStyle: any = this.getVariableStyle(attributes)
    this.upsert('variable', GText, variableStyle, container)
  }

  getCollapseStyle(attributes) {
    if (this.childrenData.length === 0) return false
    const { collapsed } = attributes
    const [width]: any = this.getSize(attributes)
    return {
      backgroundFill: '#fff',
      backgroundHeight: 16,
      backgroundLineWidth: 1,
      backgroundRadius: 0,
      backgroundStroke: '#aaa',
      backgroundWidth: 16,
      cursor: 'pointer',
      fill: '#aaa',
      fontSize: 16,
      text: collapsed ? '+' : '-',
      textAlign: 'center',
      textBaseline: 'middle',
      x: width / 2,
      y: 0,
    }
  }

  drawCollapseShape(attributes, container) {
    const collapseStyle: any = this.getCollapseStyle(attributes)
    const btn = this.upsert('collapse', Badge, collapseStyle, container)

    if (btn && !Reflect.has(btn, '__bind__')) {
      Reflect.set(btn, '__bind__', true)
      btn.addEventListener(CommonEvent.CLICK, () => {
        const { collapsed } = this.attributes
        const graph = this.context.graph
        if (collapsed) graph.expandElement(this.id)
        else graph.collapseElement(this.id)
      })
    }
  }

  getProcessBarStyle(attributes) {
    const { rate } = this.data
    const { radius } = attributes
    const color = 'orange'
    // const percent = `${Number(rate) * 100}%`
    const [width, height]: any = this.getSize(attributes)
    return {
      x: -width / 2,
      y: height / 2 - 4,
      width: width,
      height: 4,
      radius: [0, 0, radius, radius],
      fill: `linear-gradient(to right, ${color} ${'20%'}, ${'#aaa'} ${'60%'})`,
    }
  }

  drawProcessBarShape(attributes, container) {
    const processBarStyle = this.getProcessBarStyle(attributes)
    this.upsert('process-bar', GRect, processBarStyle, container)
  }

  getKeyStyle(attributes) {
    const keyStyle = super.getKeyStyle(attributes)
    return {
      ...keyStyle,
      fill: '#fff',
      lineWidth: 1,
      stroke: '#aaa',
    }
  }

  render(attributes = this.parsedAttributes, container) {
    super.render(attributes, container)
    // this.getLabelStyle(attributes)
    this.drawPriceShape(attributes, container)
    this.drawCurrencyShape(attributes, container)
    this.drawPercentShape(attributes, container)
    this.drawTriangleShape(attributes, container)
    this.drawVariableShape(attributes, container)
    this.drawProcessBarShape(attributes, container)
    this.drawCollapseShape(attributes, container)
  }
}

class FlyMarkerCubic extends CubicHorizontal {
  getMarkerStyle(attributes) {
    return {
      r: 4,
      fill: '#1677ff',
      opacity: 0.7,
      // stroke: '#c2c8d1',
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
    graphRef.current?.setItemState(evt.item, 'hoverState', true)
  }

  const handleMouseLeave = (evt: any) => {
    graphRef.current?.setItemState(evt.item, 'hoverState', false)
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
      behaviors: ['drag-canvas', 'zoom-canvas', 'drag-element', 'click-select'],
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
          stroke: '#1677fff',
          radius: 8,
          shadowColor: 'rgba(0,0,0,0.05)',
          cursor: 'pointer',
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
          // const node = evt.item
          // if (
          //   !graphRef.current
          //     ?.findById(node.getModel().id)
          //     ?.hasState('selected')
          // ) {
          //   graphRef.current?.setItemState(node, 'hover', true)
          // }
          // handleMouseEnter(evt)
        })

        graphRef.current?.on(NodeEvent.POINTER_LEAVE, evt => {
          // handleMouseLeave(evt)
        })

        if (typeof window !== 'undefined') {
          window.onresize = () => {
            if (!graphRef.current || graphRef.current?.get('destroyed')) return
            if (
              !containerRef ||
              !containerRef.current?.scrollWidth ||
              !containerRef.current?.scrollHeight
            )
              return
            graphRef.current?.changeSize(
              containerRef?.current?.scrollWidth,
              containerRef.current?.scrollHeight,
            )
          }
        }
      } else {
        graphRef.current.clear()
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
