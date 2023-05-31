/*
 * Copyright The Karbour Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React, { useEffect } from "react";
import { register } from "@antv/x6-react-shape";
import { Graph, Node, Path, Cell, Edge, Platform, StringExt } from "@antv/x6";
import { Selection } from "@antv/x6-plugin-selection";
import classnames from "classnames";
import insertCss from "insert-css";
import { Tooltip, Dropdown, Space } from "antd";
import kubernetes from "../../../assets/kubernetes.png";

import styles from "./styles.module.scss";

import "./style.css";

const data = {
  nodes: [
    {
      id: "node-0",
      shape: "data-processing-dag-node",
      x: 0,
      y: 100,
      ports: [
        {
          position: "right",
          id: "node-0-out",
          group: "out",
        },
        {
          position: "top",
          id: "node-0-top",
          group: "top",
        },
      ],
      data: {
        name: "Service",
        type: "INPUT",
        checkStatus: "sucess",
      },
    },
    {
      id: "node-1",
      shape: "data-processing-dag-node",
      x: 250,
      y: 100,
      ports: [
        {
          id: "node-1-in",
          group: "in",
        },
        {
          id: "node-1-out",
          group: "out",
        },
      ],
      data: {
        name: "Endpoints",
        type: "FILTER",
      },
    },
    {
      id: "node-3",
      shape: "data-processing-dag-node",
      x: 500,
      y: 100,
      ports: [
        {
          id: "node-3-bottom",
          group: "bottom",
        },
      ],
      data: {
        name: "OperationJob",
        type: "JOIN",
      },
    },
    {
      id: "node-4",
      shape: "data-processing-dag-node",
      x: 750,
      y: 100,
      ports: [
        {
          id: "node-4-in",
          group: "in",
        },
        {
          id: "node-4-out",
          group: "out",
        },
      ],
      data: {
        name: "InPlaceSet",
        type: "OUTPUT",
      },
    },
    {
      id: "node-5",
      shape: "data-processing-dag-node",
      x: 1000,
      y: 100,
      ports: [
        {
          id: "node-5-in",
          group: "in",
        },
      ],
      data: {
        name: "CafeDeployment",
        type: "INPUT",
      },
    },
    {
      id: "node-6",
      shape: "data-processing-dag-node",
      x: 0,
      y: 200,
      ports: [
        {
          id: "node-6-out",
          group: "out",
        },
      ],
      data: {
        name: "Deployment",
        type: "INPUT",
      },
    },
    {
      id: "node-7",
      shape: "data-processing-dag-node",
      x: 500,
      y: 300,
      ports: [
        {
          id: "node-7-top",
          group: "top",
        },
      ],
      data: {
        name: "ResourceDecoration",
        type: "INPUT",
      },
    },
    {
      id: "node-8",
      shape: "data-processing-dag-node",
      x: 250,
      y: 200,
      ports: [
        {
          id: "node-8-in",
          group: "in",
        },
        {
          id: "node-8-out",
          group: "out",
        },
      ],
      data: {
        name: "ReplicaSet",
        type: "INPUT",
      },
    },
    {
      id: "node-9",
      shape: "data-processing-dag-node",
      x: 500,
      y: 200,
      ports: [
        {
          id: "node-9-in",
          group: "in",
        },
        {
          id: "node-9-out",
          group: "out",
        },
        {
          id: "node-9-top",
          group: "top",
        },
        {
          id: "node-9-bottom",
          group: "bottom",
        },
      ],
      data: {
        name: "Pod",
        type: "INPUT",
      },
    },
    {
      id: "node-10",
      shape: "data-processing-dag-node",
      x: 750,
      y: 200,
      ports: [
        {
          id: "node-6-in",
          group: "in",
        },
      ],
      data: {
        name: "RuleSet",
        type: "INPUT",
      },
    },
    {
      id: "node-11",
      shape: "data-processing-dag-node",
      x: 250,
      y: 300,
      ports: [
        {
          id: "node-11-out",
          group: "out",
        },
      ],
      data: {
        name: "StatefulSet",
        type: "INPUT",
      },
    },
  ],
  edges: [
    {
      id: "edge-0",
      source: {
        cell: "node-0",
        port: "node-0-out",
      },
      target: {
        cell: "node-1",
        port: "node-1-in",
      },
      shape: "data-processing-curve",
      zIndex: -1,
      data: {
        source: "node-0",
        target: "node-1",
      },
    },
    {
      id: "edge-1",
      source: {
        cell: "node-1",
        port: "node-1-out",
      },
      target: {
        cell: "node-9",
        port: "node-9-in",
      },
      shape: "data-processing-curve",
      zIndex: -1,
      data: {
        source: "node-1",
        target: "node-9",
      },
    },
    {
      id: "edge-2",
      source: {
        cell: "node-6",
        port: "node-6-out",
      },
      target: {
        cell: "node-8",
        port: "node-8-in",
      },
      shape: "data-processing-curve",
      zIndex: -1,
      data: {
        source: "node-6",
        target: "node-8",
      },
    },
    {
      id: "edge-3",
      source: {
        cell: "node-8",
        port: "node-8-out",
      },
      target: {
        cell: "node-9",
        port: "node-9-in",
      },
      shape: "data-processing-curve",
      zIndex: -1,
      data: {
        source: "node-8",
        target: "node-9",
      },
    },
    {
      id: "edge-4",
      source: {
        cell: "node-11",
        port: "node-11-out",
      },
      target: {
        cell: "node-9",
        port: "node-9-in",
      },
      shape: "data-processing-curve",
      zIndex: -1,
      data: {
        source: "node-11",
        target: "node-9",
      },
    },
    {
      id: "edge-5",
      source: {
        cell: "node-3",
        port: "node-3-bottom",
      },
      target: {
        cell: "node-9",
        port: "node-9-top",
      },
      shape: "data-processing-curve",
      zIndex: -1,
      data: {
        source: "node-3",
        target: "node-9",
      },
    },
    {
      id: "edge-6",
      source: {
        cell: "node-5",
        port: "node-5-in",
      },
      target: {
        cell: "node-4",
        port: "node-4-out",
      },
      shape: "data-processing-curve",
      zIndex: -1,
      data: {
        source: "node-5",
        target: "node-4",
      },
    },
    {
      id: "edge-7",
      source: {
        cell: "node-4",
        port: "node-4-in",
      },
      target: {
        cell: "node-9",
        port: "node-9-out",
      },
      shape: "data-processing-curve",
      zIndex: -1,
      data: {
        source: "node-4",
        target: "node-9",
      },
    },
    {
      id: "edge-8",
      source: {
        cell: "node-10",
        port: "node-10-in",
      },
      target: {
        cell: "node-9",
        port: "node-9-out",
      },
      shape: "data-processing-curve",
      zIndex: -1,
      data: {
        source: "node-10",
        target: "node-9",
      },
    },
    {
      id: "edge-9",
      source: {
        cell: "node-7",
        port: "node-7-top",
      },
      target: {
        cell: "node-9",
        port: "node-9-bottom",
      },
      shape: "data-processing-curve",
      zIndex: -1,
      data: {
        source: "node-7",
        target: "node-9",
      },
    },
  ],
};

// 节点类型
enum NodeType {
  INPUT = "INPUT", // 数据输入
  FILTER = "FILTER", // 数据过滤
  JOIN = "JOIN", // 数据连接
  UNION = "UNION", // 数据合并
  AGG = "AGG", // 数据聚合
  OUTPUT = "OUTPUT", // 数据输出
  TOP = "TOP",
  BOTTOM = "BOTTOM",
  LEFT = "LEFT",
  RIGHT = "RIGHT",
}

// 元素校验状态
enum CellStatus {
  DEFAULT = "default",
  SUCCESS = "success",
  ERROR = "error",
}

// 节点位置信息
interface Position {
  x: number;
  y: number;
}

// 加工类型列表
const PROCESSING_TYPE_LIST = [
  {
    type: "FILTER",
    name: "数据筛选",
  },
  {
    type: "JOIN",
    name: "数据连接",
  },
  {
    type: "UNION",
    name: "数据合并",
  },
  {
    type: "AGG",
    name: "数据聚合",
  },
  {
    type: "OUTPUT",
    name: "数据输出",
  },
  {
    type: "LEFT",
    name: "LEFT",
  },
  {
    type: "RIGHT",
    name: "RIGHT",
  },
  {
    type: "TOP",
    name: "TOP",
  },
  {
    type: "BOTTOM",
    name: "BOTTOM",
  },
];

// 不同节点类型的icon
const NODE_TYPE_LOGO = {
  INPUT: kubernetes,
  // 'https://mdn.alipayobjects.com/huamei_f4t1bn/afts/img/A*RXnuTpQ22xkAAAAAAAAAAAAADtOHAQ/original', // 数据输入
  FILTER: kubernetes,
  // 'https://mdn.alipayobjects.com/huamei_f4t1bn/afts/img/A*ZJ6qToit8P4AAAAAAAAAAAAADtOHAQ/original', // 数据筛选
  JOIN: kubernetes,
  // 'https://mdn.alipayobjects.com/huamei_f4t1bn/afts/img/A*EHqyQoDeBvIAAAAAAAAAAAAADtOHAQ/original', // 数据连接
  UNION: kubernetes,
  // 'https://mdn.alipayobjects.com/huamei_f4t1bn/afts/img/A*k4eyRaXv8gsAAAAAAAAAAAAADtOHAQ/original', // 数据合并
  AGG: kubernetes,
  // 'https://mdn.alipayobjects.com/huamei_f4t1bn/afts/img/A*TKG8R6nfYiAAAAAAAAAAAAAADtOHAQ/original', // 数据聚合
  OUTPUT: kubernetes,
  // 'https://mdn.alipayobjects.com/huamei_f4t1bn/afts/img/A*zUgORbGg1HIAAAAAAAAAAAAADtOHAQ/original', // 数据输出
};

// const AlgoNode = () => {
//   const label = "success";
//   const status = "success";
//   return <>123</>
//   return (
//     <div className={`node ${status}`}>
//       <img src={image.logo} />
//       <span className="label">{label}</span>
//       <span className="status">
//         {status === "success" && <img src={image.success} />}
//         {status === "failed" && <img src={image.failed} />}
//         {status === "running" && <img src={image.running} />}
//       </span>
//     </div>
//   );
// };

/**
 * 根据起点初始下游节点的位置信息
 * @param node 起始节点
 * @param graph
 * @returns
 */
const getDownstreamNodePosition = (
  node: Node,
  graph: Graph,
  dx = 250,
  dy = 100
) => {
  // 找出画布中以该起始节点为起点的相关边的终点id集合
  const downstreamNodeIdList: string[] = [];
  graph.getEdges().forEach((edge) => {
    const originEdge = edge.toJSON()?.data;
    if (originEdge.source === node.id) {
      downstreamNodeIdList.push(originEdge.target);
    }
  });
  // 获取起点的位置信息
  const position = node.getPosition();
  let minX = Infinity;
  let maxY = -Infinity;
  graph.getNodes().forEach((graphNode) => {
    if (downstreamNodeIdList.indexOf(graphNode.id) > -1) {
      const nodePosition = graphNode.getPosition();
      // 找到所有节点中最左侧的节点的x坐标
      if (nodePosition.x < minX) {
        minX = nodePosition.x;
      }
      // 找到所有节点中最x下方的节点的y坐标
      if (nodePosition.y > maxY) {
        maxY = nodePosition.y;
      }
    }
  });

  return {
    x: minX !== Infinity ? minX : position.x + dx,
    y: maxY !== -Infinity ? maxY + dy : position.y,
  };
};

// 根据节点的类型获取ports
const getPortsByType = (type: NodeType, nodeId: string) => {
  let ports = [];
  switch (type) {
    case NodeType.TOP:
      ports = [
        {
          id: `${nodeId}-top`,
          group: "top",
        },
      ];
      break;
    case NodeType.BOTTOM:
      ports = [
        {
          id: `${nodeId}-bottom`,
          group: "bottom",
        },
      ];
      break;
    case NodeType.LEFT:
      ports = [
        {
          id: `${nodeId}-in`,
          group: "in",
        },
      ];
      break;
    case NodeType.RIGHT:
      ports = [
        {
          id: `${nodeId}-out`,
          group: "out",
        },
      ];
      break;
    default:
      ports = [
        {
          id: `${nodeId}-in`,
          group: "in",
        },
        {
          id: `${nodeId}-out`,
          group: "out",
        },
      ];
      break;
  }
  return ports;
};

/**
 * 创建节点并添加到画布
 * @param type 节点类型
 * @param graph
 * @param position 节点位置
 * @returns
 */
export const createNode = (
  type: NodeType,
  graph: Graph,
  position?: Position
) => {
  if (!graph) {
    return {};
  }
  let newNode = {};
  const sameTypeNodes = graph
    .getNodes()
    .filter((item) => item.getData()?.type === type);

  const typeName = PROCESSING_TYPE_LIST?.find(
    (item) => item.type === type
  )?.name;
  const id = StringExt.uuid();
  const node = {
    id,
    shape: "data-processing-dag-node",
    x: position?.x,
    y: position?.y,
    ports: getPortsByType(type, id),
    data: {
      name: `${typeName}_${sameTypeNodes.length + 1}`,
      type,
    },
  };
  newNode = graph.addNode(node);
  return newNode;
};

/**
 * 创建边并添加到画布
 * @param source
 * @param target
 * @param graph
 */
const createEdge = (source: string, target: string, graph: Graph) => {
  const edge = {
    id: StringExt.uuid(),
    shape: "data-processing-curve",
    source: {
      cell: source,
      port: `${source}-out`,
    },
    target: {
      cell: target,
      port: `${target}-in`,
    },
    zIndex: -1,
    data: {
      source,
      target,
    },
  };
  if (graph) {
    graph.addEdge(edge);
  }
};

class DataProcessingDagNode extends React.Component<{
  node: Node;
}> {
  state = {
    plusActionSelected: false,
  };

  // 创建下游的节点和边
  createDownstream = (type: NodeType) => {
    const { node } = this.props;
    const { graph } = node.model || {};
    if (graph) {
      // 获取下游节点的初始位置信息
      const position = getDownstreamNodePosition(node, graph);
      // 创建下游节点
      const newNode = createNode(type, graph, position) as any;
      const source = node.id;
      const target = newNode.id;
      // 创建该节点出发到下游节点的边
      createEdge(source, target, graph);
    }
  };

  // 点击添加下游+号
  clickPlusDragMenu = (type: NodeType) => {
    this.createDownstream(type);
    this.setState({
      plusActionSelected: false,
    });
  };

  //  获取+号下拉菜单
  getPlusDagMenu = () => {
    return (
      <ul>
        {PROCESSING_TYPE_LIST.map((item: any) => {
          const content = (
            // eslint-disable-next-line
            <a onClick={() => this.clickPlusDragMenu(item.type)}>
              <i
                className="node-mini-logo"
                style={{ backgroundImage: `url(${NODE_TYPE_LOGO[item.type]})` }}
              />

              <span>{item.name}</span>
            </a>
          );
          return (
            <li className="each-sub-menu" key={item.type}>
              {content}
            </li>
          );
        })}
      </ul>
    );
  };

  // 添加下游菜单的打开状态变化
  onPlusDropdownOpenChange = (value: boolean) => {
    this.setState({
      plusActionSelected: value,
    });
  };

  // 鼠标进入矩形主区域的时候显示连接桩
  onMainMouseEnter = () => {
    const { node } = this.props;
    // 获取该节点下的所有连接桩
    const ports = node.getPorts() || [];
    ports.forEach((port) => {
      node.setPortProp(port.id, "attrs/circle", {
        fill: "#fff",
        stroke: "#85A5FF",
      });
    });
  };

  // 鼠标离开矩形主区域的时候隐藏连接桩
  onMainMouseLeave = () => {
    const { node } = this.props;
    // 获取该节点下的所有连接桩
    const ports = node.getPorts() || [];
    ports.forEach((port) => {
      node.setPortProp(port.id, "attrs/circle", {
        fill: "transparent",
        stroke: "transparent",
      });
    });
  };

  render() {
    const { plusActionSelected } = this.state;
    const { node } = this.props;
    const data = node?.getData();
    const { name, type, status, statusMsg } = data;

    return (
      <div className="data-processing-dag-node">
        <div
          className="main-area"
          style={{
            color: name === "Pod" ? "#fff" : "none",
            border: name === "Pod" ? "2px solid #f4a460" : "none",
            background: name === "Pod" ? "#f4a460" : "none",
          }}
          onMouseEnter={this.onMainMouseEnter}
          onMouseLeave={this.onMainMouseLeave}
        >
          <div className="main-info">
            {/* 节点类型icon */}
            <i
              className="node-logo"
              style={{ backgroundImage: `url(${NODE_TYPE_LOGO[type]})` }}
            />
            <Tooltip title={name} mouseEnterDelay={0.8}>
              <div
                className="ellipsis-row node-name"
                style={{
                  color: name === "Pod" ? "#fff" : "none",
                  fontSize: 16,
                }}
              >
                {name}
              </div>
            </Tooltip>
          </div>

          {/* 节点状态信息 */}
          <div className="status-action">
            {status === CellStatus.ERROR && (
              <Tooltip title={statusMsg}>5</Tooltip>
            )}
            {status === CellStatus.SUCCESS && (
              <i className="status-icon status-icon-success">5</i>
            )}

            {/* 节点操作菜单 */}
            <div
              className="more-action-container"
              style={{
                color: name === "Pod" ? "#fff" : "none",
                fontSize: 16,
              }}
            >
              {name === "Pod" ? 10 : 1}
            </div>
          </div>
        </div>

        {/* 添加下游节点 */}
        {/* {type !== NodeType.OUTPUT && (
          <div className="plus-dag">
            <Dropdown
              dropdownRender={this.getPlusDagMenu}
              overlayClassName="processing-node-menu"
              trigger={["click"]}
              placement="bottom"
              open={plusActionSelected}
              onOpenChange={this.onPlusDropdownOpenChange}
            >
              <i
                className={classnames("plus-action", {
                  "plus-action-selected": plusActionSelected,
                })}
              />
            </Dropdown>
          </div>
        )} */}
      </div>
    );
  }
}

// insertCss(`

// `);

export default function Relationship() {
  function draw() {
    register({
      shape: "data-processing-dag-node",
      width: 212,
      height: 48,
      component: DataProcessingDagNode as any,
      // port默认不可见
      ports: {
        groups: {
          in: {
            position: {
              name: "left",
              args: {
                dx: 2,
              },
            },
            attrs: {
              circle: {
                r: 4,
                magnet: true,
                stroke: "transparent",
                strokeWidth: 1,
                fill: "transparent",
              },
            },
          },
          out: {
            position: {
              name: "right",
              args: {
                dx: -34,
              },
            },
            attrs: {
              circle: {
                r: 4,
                magnet: true,
                stroke: "transparent",
                strokeWidth: 1,
                fill: "transparent",
              },
            },
          },
          top: {
            position: {
              name: "top",
              args: {
                dx: -10,
              },
            },
            attrs: {
              circle: {
                r: 4,
                magnet: true,
                stroke: "transparent",
                strokeWidth: 1,
                fill: "transparent",
              },
            },
          },
          bottom: {
            position: {
              name: "bottom",
              args: {
                dx: -15,
              },
            },
            attrs: {
              circle: {
                r: 4,
                magnet: true,
                stroke: "transparent",
                strokeWidth: 1,
                fill: "transparent",
              },
            },
          },
        },
      },
    });

    // 注册连线
    // Graph.registerConnector(
    //   "curveConnector",
    //   (sourcePoint, targetPoint, ...param) => {
    //     console.log(param, "===sourcePoint===");
    //     const hgap = Math.abs(targetPoint.x - sourcePoint.x);
    //     const ygap = Math.abs(targetPoint.y - sourcePoint.y);

    //     const path = new Path();
    //     const isRight = sourcePoint.x > targetPoint.x;
    //     const isTop = sourcePoint.y > targetPoint.y;
    //     path.appendSegment(
    //       Path.createSegment(
    //         "M",
    //         isRight ? sourcePoint.x + 4 : sourcePoint.x - 4,
    //         sourcePoint.y
    //       )
    //     );
    //     path.appendSegment(
    //       Path.createSegment("L", sourcePoint.x, sourcePoint.y)
    //     );
    //     // 水平三阶贝塞尔曲线
    //     path.appendSegment(
    //       Path.createSegment(
    //         "C",
    //         sourcePoint.x < targetPoint.x
    //           ? sourcePoint.x + hgap / 2
    //           : sourcePoint.x - hgap / 2,
    //         sourcePoint.y,
    //         sourcePoint.x < targetPoint.x
    //           ? targetPoint.x - hgap / 2
    //           : targetPoint.x + hgap / 2,
    //         targetPoint.y,
    //         isRight ? targetPoint.x + 12 : targetPoint.x - 6,
    //         isRight ? targetPoint.y : targetPoint.y
    //       )
    //     );

    //     path.appendSegment(
    //       Path.createSegment(
    //         "L",
    //         isRight ? targetPoint.x - 4 : targetPoint.x + 2,
    //         isRight ? targetPoint.y : targetPoint.y
    //       )
    //     );

    //     return path.serialize();
    //   },
    //   true
    // );

    Edge.config({
      markup: [
        {
          tagName: "path",
          selector: "wrap",
          attrs: {
            fill: "none",
            cursor: "pointer",
            stroke: "transparent",
            strokeLinecap: "round",
          },
        },
        {
          tagName: "path",
          selector: "line",
          attrs: {
            fill: "none",
            pointerEvents: "none",
          },
        },
      ],
      // connector: { name: "curveConnector" },
      connector: {
        name: "rounded",
        args: {
          radius: 20,
        },
      },
      attrs: {
        wrap: {
          connection: true,
          strokeWidth: 10,
          strokeLinejoin: "round",
        },
        line: {
          connection: true,
          stroke: "#A2B1C3",
          strokeWidth: 1,
          targetMarker: {
            name: "classic",
            size: 8,
          },
        },
      },
    });

    Graph.registerEdge("data-processing-curve", Edge, true);

    const graph: Graph = new Graph({
      container: document.getElementById("container")!,
      panning: {
        enabled: true,
        eventTypes: ["leftMouseDown", "mouseWheel"],
      },

      mousewheel: {
        enabled: true,
        modifiers: "ctrl",
        factor: 1.1,
        maxScale: 1.5,
        minScale: 0.5,
      },
      highlighting: {
        magnetAdsorbed: {
          name: "stroke",
          args: {
            attrs: {
              fill: "#fff",
              stroke: "#31d0c6",
              strokeWidth: 4,
            },
          },
        },
      },
      connecting: {
        snap: true,
        allowBlank: false,
        allowLoop: false,
        highlight: true,
        // anchor: "center",
        // connectionPoint: "anchor",
        sourceAnchor: {
          name: "left",
          args: {
            dx: Platform.IS_SAFARI ? 4 : 8,
          },
        },
        targetAnchor: {
          name: "right",
          args: {
            dx: Platform.IS_SAFARI ? 4 : -8,
          },
        },
        createEdge(...params) {
          console.log(params, "=====params=====");
          return graph.createEdge({
            shape: "data-processing-curve",
            attrs: {
              line: {
                strokeDasharray: "5 5",
              },
            },
            zIndex: -1,
          });
        },
        // 连接桩校验
        validateConnection({ sourceMagnet, targetMagnet }) {
          // 只能从输出链接桩创建连接
          if (
            !sourceMagnet ||
            sourceMagnet.getAttribute("port-group") === "in"
          ) {
            return false;
          }
          // 只能连接到输入链接桩
          if (
            !targetMagnet ||
            targetMagnet.getAttribute("port-group") !== "in"
          ) {
            return false;
          }
          return true;
        },
      },
    });
    graph.use(
      new Selection({
        multiple: true,
        rubberEdge: true,
        rubberNode: true,
        modifiers: "shift",
        rubberband: true,
      })
    );
    // 节点状态列表
    const nodeStatusList = [
      {
        id: "node-0",
        status: "success",
      },
      {
        id: "node-1",
        status: "success",
      },
      {
        id: "node-8",
        status: "success",
      },
      {
        id: "node-3",
        status: "success",
      },
      {
        id: "node-4",
        status: "error",
        statusMsg: "错误信息示例",
      },
    ];

    // 边状态列表
    const edgeStatusList = [
      {
        id: "edge-0",
        status: "success",
      },
      {
        id: "edge-1",
        status: "success",
      },
      {
        id: "edge-2",
        status: "success",
      },
      {
        id: "edge-3",
        status: "success",
      },
    ];

    // 显示节点状态
    const showNodeStatus = () => {
      nodeStatusList.forEach((item: any) => {
        const { id, status, statusMsg } = item;
        const node = graph.getCellById(id);
        const data = node.getData() as CellStatus;
        const obj = Object.assign(data, { status, statusMsg });
        node.setData(obj);
      });
    };

    // 开启边的运行动画
    const excuteAnimate = () => {
      graph.getEdges().forEach((edge) => {
        edge.attr({
          line: {
            stroke: "#3471F9",
          },
        });
        edge.attr("line/strokeDasharray", 5);
        edge.attr("line/style/animation", "running-line 30s infinite linear");
      });
    };

    // 关闭边的动画
    const stopAnimate = () => {
      graph.getEdges().forEach((edge) => {
        edge.attr("line/strokeDasharray", 0);
        edge.attr("line/style/animation", "");
      });
      edgeStatusList.forEach((item) => {
        const { id, status } = item;
        const edge = graph.getCellById(id);
        // if (status === "success") {
        //   // edge.attr('line/stroke', '#52c41a')
        //   edge.attr("line/stroke", "#ff5500");
        // }
        // if (status === "error") {
        //   edge.attr("line/stroke", "#ff4d4f");
        // }
      });
      // 默认选中一个节点
      graph.select("node-9");
    };

    // 初始化节点/边
    const init = (data: Cell.Metadata[]) => {
      graph.fromJSON(data);
      const zoomOptions = {
        padding: {
          left: 10,
          right: 10,
        },
      };
      graph.zoomToFit(zoomOptions);
      // setTimeout(() => {
      //   excuteAnimate();
      // }, 2000);
      // setTimeout(() => {
      //   showNodeStatus();
      //   stopAnimate();
      // }, 3000);
    };

    init(data as any);
    graph.centerContent();
  }

  useEffect(() => {
    draw();
    return () => {
      // graph.dispose()
    };
  }, []);

  return (
    <div className={styles.relationship}>
      <div className={styles.overviewG6}>
        <Space className={styles.title}>
          <div className={styles.green}> </div>
          <div>正常</div>
          <div className={styles.blue} />
          <div>变更中</div>
          <div className={styles.red} />
          <div>异常</div>
          <div className={styles.black} />
          <div>删除中</div>
        </Space>
      </div>
      <div id="container" style={{ width: 1000, height: 400 }}></div>
    </div>
  );
}

// 我们用 insert-css 演示引入自定义样式
// 推荐将样式添加到自己的样式文件中
// 若拷贝官方代码，别忘了 npm install insert-css
