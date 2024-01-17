import ReactDOM from 'react-dom';
import G6 from '@antv/g6';
import type { IAbstractGraph, IG6GraphEvent } from '@antv/g6';
import type { Point } from '@antv/g-base/lib/types';
import styles from "./style.module.less";
import { memo, useLayoutEffect, useRef, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import queryString from "query-string";
import {
  Rect,
  Group,
  createNodeFromReact,
  appenAutoShapeListener,
  Image,
  Text,
} from '@antv/g6-react-node';
import Loading from "../../../../components/loading";
import transferPng from "../../../../assets/transfer.png";
import PointButton from "./PointButton";

const TextCopy: any = Text;

// 获取文本的长度
function getTextSize(str: string, maxWidth: number, fontSize: number) {
  let width = G6.Util.getTextSize(str, fontSize)[0];
  return width > maxWidth ? maxWidth : width;
}

// 文本超出隐藏 (字段, 最大长度, 字体大小)
function fittingString(str: any, maxWidth: number, fontSize: number) {
  const ellipsis = '...';
  const ellipsisLength = G6.Util.getTextSize(ellipsis, fontSize)[0];
  let currentWidth = 0;
  let res = str;
  const pattern = new RegExp('[\u4E00-\u9FA5]+'); // distinguish the Chinese charactors and letters
  str?.split('')?.forEach((letter, i) => {
    if (currentWidth > maxWidth - ellipsisLength) return;
    if (pattern?.test(letter)) {
      // Chinese charactors
      currentWidth += fontSize;
    } else {
      // get the width of single letter according to the fontSize
      currentWidth += G6.Util.getLetterWidth(letter, fontSize);
    }
    if (currentWidth > maxWidth - ellipsisLength) {
      res = `${str?.substr(0, i)}${ellipsis}`;
    }
  });
  return res;
};

type propsType = {
  value: Record<string, any>[];
  open?: boolean;
  hiddenButtonInfo?: any;
  itemWidth?: number;
  type: string;
};

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
    margin: "10px 5px"
  }
  return (
    <div style={boxStyle}>
      <div style={itemStyle}>{props?.type === 'cluster' ? model?.label : model?.id}</div>
      {/* <div style={itemStyle}>{model?.label}</div> */}
    </div>
  );
});



type IProps = {
  topologyData: any;
  topologyLoading?: boolean;
  onTopologyNodeClick?: (node: any) => void;
  isResource?: boolean;
  tableName?: string;
}

const TopologyMap = ({ onTopologyNodeClick, topologyData, topologyLoading, isResource, tableName }: IProps) => {
  const ref = useRef(null);
  const graphRef = useRef<any>();
  let graph: IAbstractGraph | null = null;
  const location = useLocation();
  const { from, type, query } = queryString.parse(location?.search)
  const navigate = useNavigate()
  const [tooltipopen, setTooltipopen] = useState(false); //悬停是否显示
  const [itemWidth, setItemWidth] = useState<number>(100); //节点宽
  const [hiddenButtontooltip, setHiddenButtontooltip] = useState<{
    x: number;
    y: number;
    e?: IG6GraphEvent;
  }>({ x: -500, y: -500, e: undefined });


  function getName(cfg: any) {
    if (type === 'resource') {
      const [left, right] = cfg?.id?.split(":");
      const leftList = left?.split('.');
      const leftListLength = leftList?.length;
      const leftLast = leftList?.[leftListLength - 1];
      return `${leftLast}:${right}`;
    }
    if (type === 'cluster' || type === 'namespace') {
      const list = cfg?.label?.split(".");
      const len = list?.length;
      return list?.[len - 1];
    }
  }

  // G6.registerNode('card-node', {
  //   draw: (cfg: any, group: any) => {
  //     const [width, height] = cfg?.size;
  //     const x = -width / 2;
  //     const y = -height / 2;
  //     // const color = cfg?.data?.count < 100 ? '#f4664a' : '#30bf78';
  //     const color = '#5B8FF9'
  //     // const color = '#C0C5D7'
  //     const keyShape = group.addShape('rect', {
  //       attrs: {
  //         x,
  //         y,
  //         width: 200,
  //         height: 60,
  //         stroke: color,
  //         // fill: '#fff',
  //         fill: '#C6E5FF',
  //         radius: 6,
  //         lineWidth: 2,
  //       },
  //       name: 'card-node-keyshape'
  //     })
  //     if (type === 'cluster' || type === 'namespace') {
  //       // const titleRect = group.addShape('rect', {
  //       //   attrs: {
  //       //     x,
  //       //     y,
  //       //     width: 200,
  //       //     height: 30,
  //       //     stroke: color,
  //       //     fill: '#fff',
  //       //     radius: 2
  //       //   },
  //       //   name: 'card-node-title-keyshape'
  //       // })
  //       // const titleIconRect = group.addShape('image', {
  //       //   attrs: {
  //       //     x: 4,
  //       //     y: 6,
  //       //     img: transferPng,
  //       //     width: 18,
  //       //     height: 18,
  //       //     cursor: 'poniter',
  //       //   },
  //       //   name: 'card-node-icon-keyshape'
  //       // })
  //       const titleLabelRect = group.addShape('text', {
  //         attrs: {
  //           text: fittingString(getName(cfg), 150, 16),
  //           x: x + 5,
  //           y: y + 8,
  //           width: getTextSize(getName(cfg), 150, 16),
  //           fill: '#00287E',
  //           textBaseline: 'top',
  //           fontSize: 16,
  //           cursor: 'pointer',
  //           fontWeight: 'bold',
  //         },
  //         name: 'card-node-label-keyshape'
  //       })
  //       const countText = group.addShape('text', {
  //         attrs: {
  //           text: cfg?.data?.count,
  //           x: x + 5,
  //           y: y + 40,
  //           fill: '#00287E',
  //           textBaseline: 'top',
  //           fontSize: 16,
  //           fontWeight: 700,
  //         },
  //         name: 'card-node-count-keyshape'
  //       })
  //       if (type === 'cluster') {
  //         const transferRectImg = group.addShape('image', {
  //           attrs: {
  //             x: x + 165,
  //             y: y + 20,
  //             img: transferPng,
  //             width: 24,
  //             height: 24,
  //             cursor: 'pointer',
  //           },
  //           name: 'card-node-transfer-keyshape'
  //         })
  //       }
  //       return keyShape;
  //     } else {
  //       // 单资源
  //       const titleLabelRect = group.addShape('text', {
  //         attrs: {
  //           text: fittingString(getName(cfg), 190, 16),
  //           x: x + 5,
  //           y: y + 20,
  //           width: getTextSize(getName(cfg), 190, 16),
  //           fill: '#646566',
  //           textBaseline: 'top',
  //           fontSize: 16,
  //           cursor: 'pointer',
  //         },
  //         name: 'card-node-label-keyshape'
  //       })
  //       return keyShape;
  //     }
  //   }
  // }, 'rect')

  function handleTransfer(evt, cfg) {
    evt.defaultPrevented = true;
    evt.stopPropagation();
    // const model = evt?.item?.get('model');
    const locator = cfg?.data?.locator;
    // 跳转到kind详情页
    const objParams = {
      from,
      type: 'kind',
      cluster: locator?.cluster,
      apiVersion: locator?.apiVersion,
      kind: locator?.kind,
      query
    }
    const urlStr = queryString.stringify(objParams);
    navigate(`/insightDetail/kind?${urlStr}`)
  }

  function handleMouseEnter(evt, cfg) {
    const model = evt?.item?.get('model');
    graph.setItemState(evt.item, 'hoverState', true);
    const { x, y } = graph?.getCanvasByPoint(model.x, model.y) as Point;
    const node = graph?.findById(model.id)?.getBBox();
    if (node) {
      setItemWidth(node?.maxX - node?.minX);
    }
    setHiddenButtontooltip({ x, y, e: evt });
    setTooltipopen(true);
  }
  function handleMouseLeave(evt, cfg) {
    graph.setItemState(evt.item, 'hoverState', false);
    setTooltipopen(false);
  }

  function handleClickNode(cfg) {
    onTopologyNodeClick(cfg)
  }




  const Card = ({ cfg }) => {
    const displayName = fittingString(getName(cfg), 190, 16);

    const isHighLight = type === 'resource' ? cfg?.locator?.name === tableName : displayName === tableName;
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
                  onMouseOver={(evt) => handleMouseEnter(evt, cfg)}
                  onMouseLeave={(evt) => handleMouseLeave(evt, cfg)}
                  width={getTextSize(getName(cfg), 190, 16)}
                // style={{
                //   width: getTextSize(getName(cfg), 190, 16),
                //   fill: '#000',
                //   fontSize: 16,
                // }}
                >
                  {displayName}
                </PointButton>
              </Rect>
              {
                (type === 'cluster' || type === 'namespace') &&
                <Rect>
                  <TextCopy
                    onClick={(event) => handleMouseEnter(event, cfg)}
                    style={{
                      fill: '#000',
                      fontSize: 16,
                    }}
                  >
                    {`${cfg?.data?.count}`}
                  </TextCopy>
                </Rect>
              }
            </Group>
            {
              (type === 'cluster') && <Rect>
                <Image
                  onClick={(event) => handleTransfer(event, cfg)}
                  style={{
                    cursor: 'pointer',
                    img: transferPng,
                    width: 20,
                    height: 20,
                  }}
                />
              </Rect>
            }
          </Rect>
        </Rect>
      </Group>
    );
  };

  G6.registerNode('card-node', createNodeFromReact(Card));

  G6.registerEdge('custom-polyline', {
    getPath(points) {
      const [sourcePoint, endPoint] = points;
      const x = (sourcePoint.x + endPoint.x) / 2;
      const y1 = sourcePoint.y;
      const y2 = endPoint.y;
      const path = [
        ['M', sourcePoint.x, sourcePoint.y],
        ['L', x, y1],
        ['L', x, y2],
        ['L', endPoint.x, endPoint.y]
      ]
      return path;
    },
    afterDraw(cfg, group) {
      const keyshape = group.find(ele => ele.get('name') === 'edge-shape');
      const style = keyshape.attr();
      const halo = group.addShape('path', {
        attrs: {
          ...style,
          lineWidth: 8,
          // color: '#2F54EB',
          opacity: 0.3
        },
        name: 'edge-halo',
      })
      halo.hide();
    },
    afterUpdate(cfg, item) {
      const group = item.getContainer();
      const keyshape = group.find(ele => ele.get('name') === 'edge-shape');
      const halo = group?.find(ele => ele.get('name') === 'edge-halo');
      const path = keyshape.attr('path');
      halo.attr('path', path);
    },
    setState(name, value, item) {
      const group = item.getContainer();
      if (name === 'hover') {
        const halo = group?.find(ele => ele.get('name') === 'edge-halo');
        if (value) {
          halo.show()
        } else {
          halo.hide();
        }
      }
    }
  }, 'cubic')

  useLayoutEffect(() => {
    if (topologyData) {
      (async () => {
        const container = document.getElementById('overviewContainer');
        const width = container?.scrollWidth || 800;
        const height = container?.scrollHeight || 400;
        const toolbar = new G6.ToolBar();
        if (!graph) {
          // const minimap = new G6.Minimap({});
          // eslint-disable-next-line
          graphRef.current = graph = new G6.Graph({
            container: ReactDOM.findDOMNode(ref.current) as HTMLElement,
            width,
            height,
            fitCenter: true,
            fitView: true,
            // renderer: 'svg',
            // fitViewPadding: [ 20, 40, 50, 20 ]
            // fitViewPadding: [30, 30, 30, 30],
            fitViewPadding: 20,
            plugins: [toolbar],
            enabledStack: true,
            modes: {
              // drag-canvas 拖拽画布  drag-node 拖拽节点 zoom-canvas 可缩放  click-select 点选节点  'scroll-canvas' 左右 上下滚动
              default: ['drag-canvas', 'drag-node', "click-select"],
            },
            layout: {
              // center: [width / 2, height / 2],
              type: "dagre",
              rankdir: "LR",
              align: "UL",
              // controlPoints: true,
              // ranksep: 20,
              // nodesep: 20,
              nodesepFunc: () => 1,
              ranksepFunc: () => 1
            },
            defaultNode: {
              // type: "cardNode",
              type: 'card-node',
              size: [240, 45],
              // linkPoints: {
              //   left: true,
              //   right: true,
              //   size: 5,
              // },
              // anchorPoints: [
              //   [0, 0.5],
              //   [1, 0.5]
              // ]
            },

            defaultEdge: {
              // 弧线边 quadratic 曲线边 cubic 直线边 line 折线边 polyline
              // type: 'custom-polyline',
              // type: 'line',
              // // color: '#5B8FF9',
              sourceAnchor: 1,
              targetAnchor: 0,
              // style: {
              //   radius: 5,
              //   offset: 20,
              //   endArrow: true,
              //   lineWidth: 2,
              //   stroke: '#C0C5D7',
              // },
              type: 'polyline',
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
                // stroke: '#2F54EB',
              }

            },
            nodeStateStyles: {
              selected: {
                // stroke: 'transparent',
                stroke: '#2F54EB',
                lineWidth: 2,
              },
              'hoverState': {
                // stroke: '#2F54EB',
                lineWidth: 3,
                // fill: '#f50',
              },
              'clickState': {
                stroke: '#2F54EB',
                lineWidth: 2,
              }
            }
          });
          // G6.Util.processParallelEdges(topologyData?.edges);
          graph.read(topologyData);
          appenAutoShapeListener(graph);
          if (topologyData?.nodes?.length < 5) {
            graph?.zoomTo(1.5, { x: width / 2, y: height / 2 }, true, { duration: 10 });
            setTimeout(() => {
              // if (graph) {
              //   graph?.fitCenter?.(); // 居中 直接不生效  使用延时
              // }
              if (graphRef?.current) {
                graphRef?.current?.fitCenter();
              }
            }, 100);
          }
          // graph.zoomTo(0.75, { x: 0, y: 0 }, true, { duration: 10 });
          // graph.on("node:click", (evt) => {
          //   // graph.setItemState(evt.item, 'clickState', true);
          //   const model = evt?.item?.get('model');
          //   onTopologyNodeClick(model)
          // })
          graph.on('card-node-transfer-keyshape:click', (evt) => {
            const model = evt?.item?.get('model');
            evt.defaultPrevented = true;
            evt.stopPropagation();
            const locator = model?.data?.locator;
            // 跳转到kind详情页
            const objParams = {
              from,
              type: 'kind',
              cluster: locator?.cluster,
              apiVersion: locator?.apiVersion,
              kind: locator?.kind,
              query
            }
            const urlStr = queryString.stringify(objParams);
            navigate(`/insightDetail/kind?${urlStr}`)
            // graph.setItemState(evt.item, 'clickState', true);
          })
          // graph.on('node:mouseenter', (evt) => {
          //   // const model = evt?.item?.get('model');
          //   // graph.setItemState(evt.item, 'hoverState', true);
          //   // const { x, y } = graph?.getCanvasByPoint(model.x, model.y) as Point;
          //   // const node = graph?.findById(model.id)?.getBBox();
          //   // if (node) {
          //   //   setItemWidth(node?.maxX - node?.minX);
          //   // }
          //   // setHiddenButtontooltip({ x, y, e: evt });
          //   // setTooltipopen(true);
          // })
          // graph.on('node:mouseleave', (evt) => {
          //   graph.setItemState(evt.item, 'hoverState', false);
          //   setTooltipopen(false);
          //   // graph.setItemState(evt.item, 'clickState', false);
          // })
          // graph.on('card-node-label-keyshape:mouseenter', (evt) => {
          //   const model = evt?.item?.get('model');
          //   graph.setItemState(evt.item, 'hoverState', true);
          //   const { x, y } = graph?.getCanvasByPoint(model.x, model.y) as Point;
          //   const node = graph?.findById(model.id)?.getBBox();
          //   if (node) {
          //     setItemWidth(node?.maxX - node?.minX);
          //   }
          //   setHiddenButtontooltip({ x, y, e: evt });
          //   setTooltipopen(true);
          // })
          // graph.on('card-node-label-keyshape:mouseleave', (evt) => {
          //   graph.setItemState(evt.item, 'hoverState', false);
          //   setTooltipopen(false);
          // })
          graph.on('edge:mouseenter', (evt) => {
            graph.setItemState(evt.item, 'hover', true);
          })
          graph.on('edge:mouseleave', (evt) => {
            graph.setItemState(evt.item, 'hover', false);
          })

          // setTimeout(() => {
          //   if (graph) {
          //     graph?.fitCenter(); // 居中 直接不生效  使用延时
          //   }
          // }, 200);

          if (typeof window !== "undefined") {
            window.onresize = () => {
              if (!graph || graph.get("destroyed")) return;
              if (!container || !container.scrollWidth || !container.scrollHeight)
                return;
              graph.changeSize(container.scrollWidth, container.scrollHeight);
            };
          }
        }
      })();
    }
    return () => {
      try {
        if (graph) {
          graph.destroy(); //清除画布;
          graphRef.current = null;
        }
      } catch (error) { }
    };
    // eslint-disable-next-line
  }, [topologyData, tableName]);



  return (
    <div className={styles.g6Topology} style={{ height: isResource ? 450 : 400 }}>
      {
        topologyLoading
          ? <Loading />
          : (
            <div ref={ref} id="overviewContainer" className={styles.overviewG6}>
              {tooltipopen ? (
                <OverviewTooltip
                  type={type as string}
                  value={
                    // (ruleCheckStatisticsResponseData?.result &&
                    //   hiddenButtontooltip?.e?.item?.getModel()?.name &&
                    //   ruleCheckStatisticsResponseData?.result[
                    //   hiddenButtontooltip?.e?.item?.getModel()?.name as string
                    //   ]) ||
                    []
                  }
                  itemWidth={itemWidth}
                  hiddenButtonInfo={hiddenButtontooltip}
                  open={tooltipopen}
                />
              ) : null}
            </div>
          )
      }
    </div>
  )
}

export default TopologyMap;
