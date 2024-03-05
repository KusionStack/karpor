import React, { memo, useLayoutEffect, useRef, useState } from 'react'
import ReactDOM from 'react-dom'
import G6 from '@antv/g6'
import type { IAbstractGraph, IG6GraphEvent } from '@antv/g6'
import type { Point } from '@antv/g-base/lib/types'
import { useLocation, useNavigate } from 'react-router-dom'
import queryString from 'query-string'
import {
  Rect,
  Group,
  createNodeFromReact,
  appenAutoShapeListener,
  Image,
  Text,
} from '@antv/g6-react-node'
import Loading from '../../../../components/loading'
import transferPng from '../../../../assets/transfer.png'
import PointButton from './PointButton'

import styles from './style.module.less'

const TextCopy: any = Text

// 获取文本的长度
function getTextSize(str: string, maxWidth: number, fontSize: number) {
  const width = G6.Util.getTextSize(str, fontSize)[0]
  return width > maxWidth ? maxWidth : width
}

// 文本超出隐藏 (字段, 最大长度, 字体大小)
function fittingString(str: any, maxWidth: number, fontSize: number) {
  const ellipsis = '...'
  const ellipsisLength = G6.Util.getTextSize(ellipsis, fontSize)[0]
  let currentWidth = 0
  let res = str
  const pattern = new RegExp('[\u4E00-\u9FA5]+') // distinguish the Chinese charactors and letters
  str?.split('')?.forEach((letter, i) => {
    if (currentWidth > maxWidth - ellipsisLength) return
    if (pattern?.test(letter)) {
      // Chinese charactors
      currentWidth += fontSize
    } else {
      // get the width of single letter according to the fontSize
      currentWidth += G6.Util.getLetterWidth(letter, fontSize)
    }
    if (currentWidth > maxWidth - ellipsisLength) {
      res = `${str?.substr(0, i)}${ellipsis}`
    }
  })
  return res
}

type propsType = {
  value?: Record<string, any>[]
  open?: boolean
  hiddenButtonInfo?: any
  itemWidth?: number
  type: string
}

// eslint-disable-next-line react/display-name
const OverviewTooltip = memo((props: propsType) => {
  const model = props?.hiddenButtonInfo?.e.item?.get('model')
  const boxStyle: any = {
    background: '#fff',
    border: '1px solid #f5f5f5',
    position: 'absolute',
    top: props?.hiddenButtonInfo?.y - 20 || -500,
    left: props?.hiddenButtonInfo?.x + (props?.itemWidth || 100) / 2 || -500, //居中
    zIndex: 5,
    padding: 10,
    borderRadius: 8,
    fontSize: 12,
  }
  const itemStyle = {
    color: '#646566',
    margin: '10px 5px',
  }
  return (
    <div style={boxStyle}>
      <div style={itemStyle}>
        {props?.type === 'cluster' ? model?.label : model?.id}
      </div>
      {/* <div style={itemStyle}>{model?.label}</div> */}
    </div>
  )
})

type IProps = {
  topologyData: any
  topologyLoading?: boolean
  onTopologyNodeClick?: (node: any) => void
  isResource?: boolean
  tableName?: string
}

const TopologyMap = ({
  onTopologyNodeClick,
  topologyData,
  topologyLoading,
  isResource,
  tableName,
}: IProps) => {
  const ref = useRef(null)
  const graphRef = useRef<any>()
  let graph: IAbstractGraph | null = null
  const location = useLocation()
  const { from, type, query } = queryString.parse(location?.search)
  const navigate = useNavigate()
  const [tooltipopen, setTooltipopen] = useState(false) //悬停是否显示
  const [itemWidth, setItemWidth] = useState<number>(100) //节点宽
  const [hiddenButtontooltip, setHiddenButtontooltip] = useState<{
    x: number
    y: number
    e?: IG6GraphEvent
  }>({ x: -500, y: -500, e: undefined })

  function getName(cfg: any) {
    if (type === 'resource') {
      const [left, right] = cfg?.id?.split(':')
      const leftList = left?.split('.')
      const leftListLength = leftList?.length
      const leftLast = leftList?.[leftListLength - 1]
      return `${leftLast}:${right}`
    }
    if (type === 'cluster' || type === 'namespace') {
      const list = cfg?.label?.split('.')
      const len = list?.length
      return list?.[len - 1]
    }
  }

  function handleTransfer(evt, cfg) {
    evt.defaultPrevented = true
    evt.stopPropagation()
    // const model = evt?.item?.get('model');
    const locator = cfg?.data?.locator
    // 跳转到kind详情页
    const objParams = {
      from,
      type: 'kind',
      cluster: locator?.cluster,
      apiVersion: locator?.apiVersion,
      kind: locator?.kind,
      query,
    }
    const urlStr = queryString.stringify(objParams)
    navigate(`/insightDetail/kind?${urlStr}`)
  }

  function handleMouseEnter(evt) {
    const model = evt?.item?.get('model')
    graph.setItemState(evt.item, 'hoverState', true)
    const { x, y } = graph?.getCanvasByPoint(model.x, model.y) as Point
    const node = graph?.findById(model.id)?.getBBox()
    if (node) {
      setItemWidth(node?.maxX - node?.minX)
    }
    setHiddenButtontooltip({ x, y, e: evt })
    setTooltipopen(true)
  }
  function handleMouseLeave(evt) {
    graph.setItemState(evt.item, 'hoverState', false)
    setTooltipopen(false)
  }

  function handleClickNode(cfg) {
    onTopologyNodeClick(cfg)
  }

  const Card = ({ cfg }: any) => {
    const displayName = fittingString(getName(cfg), 190, 16)

    const isHighLight =
      type === 'resource'
        ? cfg?.locator?.name === tableName
        : displayName === tableName
    return (
      <Group draggable>
        <Rect
          style={{
            width: 250,
            height: 'auto',
            fill: isHighLight ? '#fff' : '#C6E5FF',
            shadowColor: '#eee',
            shadowBlur: 30,
            radius: [8],
            justifyContent: 'center',
            padding: [10, 0],
            stroke: '#C6E5FF',
          }}
          draggable
        >
          <Rect
            onClick={() => handleClickNode(cfg)}
            style={{
              cursor: 'pointer',
              stroke: 'transparent',
              fill: 'transparent',
              flexDirection: 'row',
              justifyContent: 'space-between',
              alignItems: 'center',
              margin: [0, 15],
            }}
          >
            <Group>
              <Rect
                onClick={() => handleClickNode(cfg)}
                style={{
                  stroke: 'transparent',
                  fill: 'transparent',
                  margin: [0, 10, 10, 0],
                }}
              >
                <PointButton
                  onClick={() => handleClickNode(cfg)}
                  onMouseOver={evt => handleMouseEnter(evt)}
                  onMouseLeave={evt => handleMouseLeave(evt)}
                  width={getTextSize(getName(cfg), 190, 16)}
                >
                  {displayName}
                </PointButton>
              </Rect>
              {(type === 'cluster' || type === 'namespace') && (
                <Rect>
                  <TextCopy
                    onClick={event => handleMouseEnter(event)}
                    style={{
                      fill: '#000',
                      fontSize: 16,
                    }}
                  >
                    {`${cfg?.data?.count}`}
                  </TextCopy>
                </Rect>
              )}
            </Group>
            {type === 'cluster' && (
              <Rect>
                <Image
                  onClick={event => handleTransfer(event, cfg)}
                  style={{
                    cursor: 'pointer',
                    img: transferPng,
                    width: 20,
                    height: 20,
                  }}
                />
              </Rect>
            )}
          </Rect>
        </Rect>
      </Group>
    )
  }

  G6.registerNode('card-node', createNodeFromReact(Card))

  G6.registerEdge(
    'custom-polyline',
    {
      getPath(points) {
        const [sourcePoint, endPoint] = points
        const x = (sourcePoint.x + endPoint.x) / 2
        const y1 = sourcePoint.y
        const y2 = endPoint.y
        const path = [
          ['M', sourcePoint.x, sourcePoint.y],
          ['L', x, y1],
          ['L', x, y2],
          ['L', endPoint.x, endPoint.y],
        ]
        return path
      },
      afterDraw(cfg, group) {
        const keyshape = group.find(ele => ele.get('name') === 'edge-shape')
        const style = keyshape.attr()
        const halo = group.addShape('path', {
          attrs: {
            ...style,
            lineWidth: 8,
            // color: '#2F54EB',
            opacity: 0.3,
          },
          name: 'edge-halo',
        })
        halo.hide()
      },
      afterUpdate(cfg, item) {
        const group = item.getContainer()
        const keyshape = group.find(ele => ele.get('name') === 'edge-shape')
        const halo = group?.find(ele => ele.get('name') === 'edge-halo')
        const path = keyshape.attr('path')
        halo.attr('path', path)
      },
      setState(name, value, item) {
        const group = item.getContainer()
        if (name === 'hover') {
          const halo = group?.find(ele => ele.get('name') === 'edge-halo')
          if (value) {
            halo.show()
          } else {
            halo.hide()
          }
        }
      },
    },
    'cubic',
  )

  useLayoutEffect(() => {
    if (topologyData) {
      ;(async () => {
        const container = document.getElementById('overviewContainer')
        const width = container?.scrollWidth || 800
        const height = container?.scrollHeight || 400
        const toolbar = new G6.ToolBar()
        if (!graph) {
          // const minimap = new G6.Minimap({});
          // eslint-disable-next-line
          graphRef.current = graph = new G6.Graph({
            // eslint-disable-next-line react/no-find-dom-node
            container: ReactDOM.findDOMNode(ref.current) as HTMLElement,
            width,
            height,
            fitCenter: true,
            fitView: true,
            // renderer: 'svg',
            fitViewPadding: 20,
            plugins: [toolbar],
            enabledStack: true,
            modes: {
              // drag-canvas 拖拽画布  drag-node 拖拽节点 zoom-canvas 可缩放  click-select 点选节点  'scroll-canvas' 左右 上下滚动
              default: ['drag-canvas', 'drag-node', 'click-select'],
            },
            layout: {
              type: 'dagre',
              rankdir: 'LR',
              align: 'UL',
              nodesepFunc: () => 1,
              ranksepFunc: () => 1,
            },
            defaultNode: {
              type: 'card-node',
              size: [240, 45],
            },

            defaultEdge: {
              type: 'polyline',
              sourceAnchor: 1,
              targetAnchor: 0,
              //线条样式
              style: {
                radius: 10,
                offset: 20,
                endArrow: true,
                lineWidth: 2,
                stroke: '#C0C5D7',
              },
            },
            edgeStateStyles: {
              hover: {
                lineWidth: 6,
              },
            },
            nodeStateStyles: {
              selected: {
                stroke: '#2F54EB',
                lineWidth: 2,
              },
              hoverState: {
                lineWidth: 3,
              },
              clickState: {
                stroke: '#2F54EB',
                lineWidth: 2,
              },
            },
          })
          graph.read(topologyData)
          appenAutoShapeListener(graph)
          if (topologyData?.nodes?.length < 5) {
            graph?.zoomTo(1.5, { x: width / 2, y: height / 2 }, true, {
              duration: 10,
            })
            setTimeout(() => {
              if (graphRef?.current) {
                graphRef?.current?.fitCenter()
              }
            }, 100)
          }
          graph.on('card-node-transfer-keyshape:click', evt => {
            const model = evt?.item?.get('model')
            evt.defaultPrevented = true
            evt.stopPropagation()
            const locator = model?.data?.locator
            // 跳转到kind详情页
            const objParams = {
              from,
              type: 'kind',
              cluster: locator?.cluster,
              apiVersion: locator?.apiVersion,
              kind: locator?.kind,
              query,
            }
            const urlStr = queryString.stringify(objParams)
            navigate(`/insightDetail/kind?${urlStr}`)
          })
          graph.on('edge:mouseenter', evt => {
            graph.setItemState(evt.item, 'hover', true)
          })
          graph.on('edge:mouseleave', evt => {
            graph.setItemState(evt.item, 'hover', false)
          })
          if (typeof window !== 'undefined') {
            window.onresize = () => {
              if (!graph || graph.get('destroyed')) return
              if (
                !container ||
                !container.scrollWidth ||
                !container.scrollHeight
              )
                return
              graph.changeSize(container.scrollWidth, container.scrollHeight)
            }
          }
        }
      })()
    }
    return () => {
      try {
        if (graph) {
          graph.destroy() //清除画布;
          graphRef.current = null
        }
      } catch (error) {}
    }
    // eslint-disable-next-line
  }, [topologyData, tableName])

  return (
    <div
      className={styles.g6_topology}
      style={{ height: isResource ? 450 : 400 }}
    >
      {topologyLoading ? (
        <Loading />
      ) : (
        <div ref={ref} id="overviewContainer" className={styles.g6_overview}>
          {tooltipopen ? (
            <OverviewTooltip
              type={type as string}
              itemWidth={itemWidth}
              hiddenButtonInfo={hiddenButtontooltip}
              open={tooltipopen}
            />
          ) : null}
        </div>
      )}
    </div>
  )
}

export default TopologyMap
