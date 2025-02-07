import React, { forwardRef, useImperativeHandle, useRef, useState } from 'react'
import { Select } from 'antd'
import G6 from '@antv/g6'
import type {
  IG6GraphEvent,
  IGroup,
  ModelConfig,
  IAbstractGraph,
} from '@antv/g6'
import { useLocation, useNavigate } from 'react-router-dom'
import queryString from 'query-string'
import { useTranslation } from 'react-i18next'
import Loading from '@/components/loading'
import transferImg from '@/assets/transfer.png'
import { ICON_MAP } from '@/utils/images'

import styles from './style.module.less'
interface NodeConfig extends ModelConfig {
  data?: {
    name?: string
    count?: number
    resourceGroup?: {
      name: string
      [key: string]: any
    }
  }
  label?: string
  id?: string
  resourceGroup?: {
    name: string
    [key: string]: any
  }
}

interface NodeModel {
  id: string
  name?: string
  label?: string
  resourceGroup?: {
    name?: string
    apiVersion: string
    cluster: string
    kind: string
  }
  data?: {
    count?: number
    resourceGroup?: {
      name?: string
      apiVersion: string
      cluster: string
      kind: string
    }
  }
}

function getTextWidth(str: string, fontSize: number) {
  const canvas = document.createElement('canvas')
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  const context = canvas.getContext('2d')!
  context.font = `${fontSize}px sans-serif`
  return context?.measureText(str)?.width
}

function fittingString(str: string, maxWidth: number, fontSize: number) {
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

function getNodeName(cfg: NodeConfig, type: string) {
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

interface OverviewTooltipProps {
  type: string
  itemWidth: number
  hiddenButtonInfo: {
    x: number
    y: number
    e?: IG6GraphEvent
  }
  open: boolean
}

const OverviewTooltip: React.FC<OverviewTooltipProps> = ({
  type,
  hiddenButtonInfo,
}) => {
  const model = hiddenButtonInfo?.e?.item?.get('model') as NodeModel
  const boxStyle: any = {
    background: '#fff',
    border: '1px solid #f5f5f5',
    position: 'absolute',
    top: hiddenButtonInfo?.y || -500,
    left: hiddenButtonInfo?.x + 14 || -500,
    transform: 'translate(-50%, -100%)',
    zIndex: 5,
    padding: '6px 12px',
    borderRadius: 8,
    boxShadow: '0 2px 8px rgba(0,0,0,0.15)',
  }

  const {
    data: { resourceGroup },
  } = model

  const itemStyle = {
    color: '#333',
    fontSize: 14,
    whiteSpace: 'nowrap',
  }

  const labelStyle = {
    color: '#8a8a8a',
    fontSize: 14,
    whiteSpace: 'nowrap',
  }

  return (
    <div style={boxStyle}>
      <div style={itemStyle}>
        {type === 'cluster' ? model?.label : model?.id}
      </div>
      <div style={itemStyle}>
        <span style={labelStyle}>apiVersion: </span>
        {resourceGroup?.apiVersion}
      </div>
      <div style={itemStyle}>
        <span style={labelStyle}>cluster: </span>
        {resourceGroup?.cluster}
      </div>
      <div style={itemStyle}>
        <span style={labelStyle}>kind: </span>
        {resourceGroup?.kind}
      </div>
    </div>
  )
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
    tableName,
    selectedCluster,
    clusterOptions,
    handleChangeCluster,
  } = props
  const { t } = useTranslation()
  const containerRef = useRef(null)
  const graphRef = useRef<IAbstractGraph | null>(null)
  const location = useLocation()
  const { from, type, query } = queryString.parse(location?.search)
  const navigate = useNavigate()
  const [tooltipopen, setTooltipopen] = useState(false)
  const [itemWidth, setItemWidth] = useState<number>(100)
  const [hiddenButtontooltip, setHiddenButtontooltip] = useState<{
    x: number
    y: number
    e?: IG6GraphEvent
  }>({ x: -500, y: -500, e: undefined })

  function handleMouseEnter(evt) {
    graphRef.current?.setItemState(evt.item, 'hoverState', true)
    const bbox = evt.item.getBBox()
    const point = graphRef.current?.getCanvasByPoint(bbox.centerX, bbox.minY)
    if (bbox) {
      setItemWidth(bbox.width)
    }
    setHiddenButtontooltip({ x: point.x, y: point.y - 5, e: evt })
    setTooltipopen(true)
  }

  const handleMouseLeave = (evt: IG6GraphEvent) => {
    graphRef.current?.setItemState(evt.item, 'hoverState', false)
    setTooltipopen(false)
  }

  G6.registerNode(
    'card-node',
    {
      draw(cfg: NodeConfig, group: IGroup) {
        const displayName = getNodeName(cfg, type as string)
        const count = cfg.data?.count
        const nodeWidth = type === 'cluster' ? 240 : 200

        // Create main container
        const rect = group.addShape('rect', {
          attrs: {
            x: 0,
            y: 0,
            width: nodeWidth,
            height: 48,
            radius: 6,
            fill: '#ffffff',
            stroke: '#e6f4ff',
            lineWidth: 1,
            shadowColor: 'rgba(0,0,0,0.06)',
            shadowBlur: 8,
            shadowOffsetX: 0,
            shadowOffsetY: 2,
            cursor: 'pointer',
          },
          name: 'node-container',
        })

        // Add side accent
        group.addShape('rect', {
          attrs: {
            x: 0,
            y: 0,
            width: 3,
            height: 48,
            radius: [3, 0, 0, 3],
            fill: '#1677ff',
            opacity: 0.4,
          },
          name: 'node-accent',
        })

        // Add Kubernetes icon
        const iconSize = 32
        const kind = cfg?.data?.resourceGroup?.kind || ''
        group.addShape('image', {
          attrs: {
            x: 16,
            y: (48 - iconSize) / 2,
            width: iconSize,
            height: iconSize,
            img: ICON_MAP[kind as keyof typeof ICON_MAP] || ICON_MAP.Kubernetes,
          },
          name: 'node-icon',
        })

        // Add title text
        group.addShape('text', {
          attrs: {
            x: 52,
            y: 24,
            text: fittingString(displayName || '', 100, 14),
            fontSize: 14,
            fontWeight: 500,
            fill: '#1677ff',
            cursor: 'pointer',
            textBaseline: 'middle',
            fontFamily:
              '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial',
          },
          name: 'node-label',
        })

        if (typeof count === 'number') {
          const textWidth = getTextWidth(`${count}`, 12)
          const circleSize = Math.max(textWidth + 12, 20)
          const circleX = 170
          const circleY = 24

          // Add count background
          group.addShape('circle', {
            attrs: {
              x: circleX,
              y: circleY,
              r: circleSize / 2,
              fill: '#f0f5ff',
            },
            name: 'count-background',
          })

          // Add count text
          group.addShape('text', {
            attrs: {
              x: circleX,
              y: circleY,
              text: `${count}`,
              fontSize: 12,
              fontWeight: 500,
              fill: '#1677ff',
              textAlign: 'center',
              textBaseline: 'middle',
              fontFamily:
                '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial',
            },
            name: 'count-text',
          })
        }

        if (type === 'cluster') {
          const iconTransferSize = 20
          group.addShape('image', {
            attrs: {
              x: 210,
              y: 14,
              width: iconTransferSize,
              height: iconTransferSize,
              img: transferImg,
            },
            name: 'transfer-icon',
          })
        }
        return rect
      },

      afterDraw(cfg: NodeConfig, group: IGroup) {
        const transferIcon = group.find(
          element => element.get('name') === 'transfer-icon',
        )
        if (transferIcon) {
          transferIcon.on('mouseenter', evt => {
            evt.defaultPrevented = true
            evt.stopPropagation()
            transferIcon.attr('cursor', 'pointer')
          })
          transferIcon.on('mouseleave', () => {
            transferIcon.attr('cursor', '')
          })
          transferIcon.on('click', evt => {
            evt.defaultPrevented = true
            evt.stopPropagation()
            const resourceGroup = cfg?.data?.resourceGroup
            const objParams = {
              from,
              type: 'kind',
              cluster: resourceGroup?.cluster,
              apiVersion: resourceGroup?.apiVersion,
              kind: resourceGroup?.kind,
              query,
            }
            const urlStr = queryString.stringify(objParams)
            navigate(`/insightDetail/kind?${urlStr}`)
          })
        }
      },
    },
    'single-node',
  )

  G6.registerEdge(
    'running-edge',
    {
      afterDraw(cfg, group) {
        const shape = group?.get('children')[0]
        if (!shape) return

        // Get the path shape
        const startPoint = shape.getPoint(0)

        // Create animated circle
        const circle = group.addShape('circle', {
          attrs: {
            x: startPoint.x,
            y: startPoint.y,
            fill: '#1677ff',
            r: 2,
            opacity: 0.8,
          },
          name: 'running-circle',
        })

        // Add movement animation
        circle.animate(
          ratio => {
            const point = shape.getPoint(ratio)
            return {
              x: point.x,
              y: point.y,
            }
          },
          {
            repeat: true,
            duration: 2000,
          },
        )
      },
      setState(name, value, item) {
        const shape = item.get('keyShape')
        if (name === 'hover') {
          shape?.attr('stroke', value ? '#1677ff' : '#c2c8d1')
          shape?.attr('lineWidth', value ? 2 : 1)
          shape?.attr('strokeOpacity', value ? 1 : 0.7)
        }
      },
    },
    'cubic', // Extend from built-in cubic edge
  )

  function initGraph() {
    const container = containerRef.current
    const width = container?.scrollWidth
    const height = container?.scrollHeight
    const toolbar = new G6.ToolBar()
    return new G6.Graph({
      container,
      width,
      height,
      fitCenter: true,
      plugins: [toolbar],
      enabledStack: true,
      modes: {
        default: ['drag-canvas', 'drag-node', 'click-select'],
      },
      animate: true,
      layout: {
        type: 'dagre',
        rankdir: 'LR',
        align: 'UL',
        nodesep: 10,
        ranksep: 40,
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
      defaultNode: {
        type: 'card-node',
        size: [200, 60],
        style: {
          fill: '#fff',
          stroke: '#e5e6e8',
          radius: 4,
          shadowColor: 'rgba(0,0,0,0.05)',
          shadowBlur: 4,
          shadowOffsetX: 0,
          shadowOffsetY: 2,
          cursor: 'pointer',
        },
      },
      defaultEdge: {
        type: 'running-edge',
        style: {
          radius: 10,
          offset: 5,
          endArrow: {
            path: G6.Arrow.triangle(4, 6, 0),
            d: 0,
            fill: '#c2c8d1',
          },
          stroke: '#c2c8d1',
          lineWidth: 1,
          strokeOpacity: 0.7,
          curveness: 0.5,
        },
        labelCfg: {
          autoRotate: true,
          style: {
            fill: '#86909c',
            fontSize: 12,
          },
        },
      },
      edgeStateStyles: {
        hover: {
          lineWidth: 2,
        },
      },
      nodeStateStyles: {
        selected: {
          stroke: '#1677ff',
          shadowColor: 'rgba(22,119,255,0.12)',
          fill: '#f0f5ff',
          opacity: 0.8,
        },
        hoverState: {
          stroke: '#1677ff',
          shadowColor: 'rgba(22,119,255,0.12)',
          fill: '#f0f5ff',
          opacity: 0.8,
        },
        clickState: {
          stroke: '#1677ff',
          shadowColor: 'rgba(22,119,255,0.12)',
          fill: '#f0f5ff',
          opacity: 0.8,
        },
      },
    })
  }

  function setHightLight() {
    graphRef.current.getNodes().forEach(node => {
      const model: any = node.getModel()
      const displayName = getNodeName(model, type as string)
      const isHighLight =
        type === 'resource'
          ? model?.data?.resourceGroup?.name === tableName
          : displayName === tableName
      if (isHighLight) {
        graphRef.current?.setItemState(node, 'selected', true)
      }
    })
  }

  function drawGraph(topologyData) {
    if (topologyData) {
      if (type === 'resource') {
        graphRef.current?.destroy()
        graphRef.current = null
      }
      if (!graphRef.current) {
        graphRef.current = initGraph()
        graphRef.current?.read(topologyData)

        setHightLight()

        graphRef.current?.on('node:click', evt => {
          const node = evt.item
          const model = node.getModel()
          setTooltipopen(false)
          graphRef.current?.getNodes().forEach(n => {
            graphRef.current?.setItemState(n, 'selected', false)
          })
          graphRef.current?.setItemState(node, 'selected', true)
          onTopologyNodeClick?.(model)
        })

        graphRef.current?.on('node:mouseenter', evt => {
          const node = evt.item
          if (
            !graphRef.current
              ?.findById(node.getModel().id)
              ?.hasState('selected')
          ) {
            graphRef.current?.setItemState(node, 'hover', true)
          }
          handleMouseEnter(evt)
        })

        graphRef.current?.on('node:mouseleave', evt => {
          handleMouseLeave(evt)
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
        graphRef.current.changeData(topologyData)
        setTimeout(() => {
          graphRef.current.fitCenter()
        }, 100)
        setHightLight()
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
        {tooltipopen ? (
          <OverviewTooltip
            type={type as string}
            itemWidth={itemWidth}
            hiddenButtonInfo={hiddenButtontooltip}
            open={tooltipopen}
          />
        ) : null}
      </div>
    </div>
  )
})

export default TopologyMap
