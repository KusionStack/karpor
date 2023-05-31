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

import G6 from '@antv/g6';
import React, { useEffect } from "react";
import styles from "./styles.module.scss";

export default function () {//props为传入的参数

  const data = {
    id: '0',
    label: 'A',
    anchorPoints: [
      [1, 0.5],
    ],
    children: [
      {
        id: '1',
        label: 'B',
        anchorPoints: [
          [1, 0.5],
          [0, 0.5]],
        children: [
          { id: '2', label: 'B1', anchorPoints: [[0, 0.5]] },
          { id: '3', label: 'B2', anchorPoints: [[0, 0.5]] },
        ]
      },
      {
        id: '4',
        label: 'C',
        anchorPoints: [
          [1, 0.5],
          [0, 0.5]],
        children: [
          {
            id: '5',
            label: 'C1',
            anchorPoints: [
              [1, 0.5],
              [0, 0.5]],
            children: [
              { id: '6', label: 'C1-1', },
            ]
          },
          {
            id: '7',
            label: 'C2',
            anchorPoints: [
              [1, 0.5],
              [0, 0.5]],
            children: [
              { id: '8', label: 'C2-2', },
            ]
          },
        ]
      },
      // {
      //   id: '9',
      //   label: 'D',
      //   anchorPoints: [
      //     [1, 0.5],
      //     [0, 0.5]],
      //   children: [
      //     { id: '10', label: 'D1', anchorPoints: [[0, 0.5]] },
      //     { id: '11', label: 'D2', anchorPoints: [[0, 0.5]] },
      //   ]
      // },
      // {
      //   id: '12',
      //   label: 'E',
      //   anchorPoints: [[0, 0.5], [1, 0.5]],
      //   children: [
      //     { id: '13', label: 'E1', anchorPoints: [[0, 0.5]], },
      //     { id: '14', label: 'E2', anchorPoints: [[0, 0.5]], },
      //   ]
      // },
    ]
  }


  const initTree = (data: any) => {
    // 画布宽高
    const width = document.getElementById('container')!.scrollWidth;
    const height = document.getElementById('container')!.scrollHeight || 500;

    G6.Util.traverseTree(data, (d: any) => {
      d.leftIcon = {
        style: {
          fill: '#fafafa',
          stroke: '#fafafa',
        },
        img: require("../../../assets/kubernetes.png"),
      };
      return true;
    });

    G6.registerNode(
      'icon-node',
      {
        options: {
          size: [120, 40],
          stroke: '#91d5ff',
          fill: '#91d5ff',
        },
        draw(cfg: any, group: any) {
          const styles = (this as any).getShapeStyle(cfg);
          const { labelCfg = {} } = cfg;

          const w = styles.width;
          const h = styles.height;

          const keyShape: any = group.addShape('rect', {
            attrs: {
              ...styles,
              x: -w / 2,
              y: -h / 2,
            },
          });

          if (cfg.leftIcon) {
            const { style, img } = cfg.leftIcon;
            group.addShape('rect', {
              attrs: {
                x: 1 - w / 2,
                y: 1 - h / 2,
                width: 38,
                height: styles.height - 2,
                fill: '#8c8c8c',
                ...style,
              },
            });

            group.addShape('image', {
              attrs: {
                x: 8 - w / 2,
                y: 8 - h / 2,
                width: 20,
                height: 20,
                img: img || 'https://g.alicdn.com/cm-design/arms-trace/1.0.155/styles/armsTrace/images/TAIR.png',
              },
              // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
              name: 'image-shape',
            });
          }

          if (cfg.label) {
            group.addShape('text', {
              attrs: {
                ...labelCfg.style,
                text: cfg.label,
                x: 50 - w / 2,
                y: 28 - h / 2,
              },
            });
          }

          return keyShape;
        },
        update: undefined,
      },
      'rect',
    );

    const graph = new G6.TreeGraph({
      // 图的  DOM 容器，可以传入该 DOM 的 id 或者直接传入容器的 HTML 节点对象。
      container: 'container',
      width,
      height,
      fitView: true,
      minZoom: 0.5,
      maxZoom: 1.5,
      modes: {
        default: [
          {
            type: 'collapse-expand',
            onChange: function onChange(item: any, collapsed) {
              const data = item.getModel();
              data.collapsed = collapsed;
              return true;
            },
          },
          'drag-canvas',
          'zoom-canvas',
        ],
      },
      // 默认状态下节点的配置，比如 type, size, color。会被写入的 data 覆盖。
      defaultNode: {
        type: 'icon-node',
        size: [120, 40],
        anchorPoints: [
          [0, 0.5],
          [1, 0.5],
        ],
        style: {
          fill: '#fff',
          stroke: '#d9d9d9',
          radius: 2,
        },
        labelCfg: {
          style: {
            fill: '#000',
            fontSize: 16,
          },
        },
      },
      // 默认状态下边的配置，比如 type, size, color。会被写入的 data 覆盖。
      // defaultEdge: {
      //   // 指定边的类型，可以是内置边的类型名称，也可以是自定义边的名称。默认为 'line'
      //   type: 'polyline',
      //   // 边的样式属性
      //   style: {
      //     // 边的颜色
      //     stroke: '#1890ff',
      //     lineWidth: 2
      //   },
      //   endArrow: {
      //     path: 'M 0,0 L 12, 6 L 9,0 L 12, -6 Z',
      //     fill: '#d9d9d9',
      //     d: -20,
      //   },
      // },
      defaultEdge: {
        type: 'polyline',
        style: {
          stroke: '#d9d9d9',
          endArrow: {
            path: 'M 0,0 L 8, 2 L 5,0 L 8, -2 Z',
            fill: '#d9d9d9',
            d: 0,
          },
        },
      },
      // 布局配置项，使用 type 字段指定使用的布局方式
      layout: {
        // 布局名称
        type: 'compactBox',
        // layout 的方向。
        direction: 'LR', // H / V / LR / RL / TB / BT
        nodesep: 10, // 节点层距离（上下）
        ranksep: 20, // 节点之间距离（左右）
        controlPoints: true,
        // 下面的参数都是一个节点，当存在某些奇葩点节点的时候，可以通过以下控制
        // 节点 id 的回调函数
        getId: function getId(d: any) {
          return d.id;
        },
        // 节点高度的回调函数
        getHeight: function getHeight() {
          return 16;
        },
        // 节点宽度的回调函数
        getWidth: function getWidth() {
          return 16;
        },
        // 节点纵向间距的回调函数
        getVGap: function getVGap() {
          return 40;
        },
        // 节点横向间距的回调函数
        getHGap: function getHGap() {
          return 70;
        },
      },
      // 动画
      animate: true,
    });
    // 设置各个节点样式及其他配置，以及在各个状态下节点的 KeyShape 的样式
    // 该方法必须在 render 之前调用，否则不起作用
    // 使用 graph.node(nodeFn) 配置 > 数据中动态配置 > 实例化图时全局配置
    // graph.node(function (node) {
    //   return {
    //     label: node.label,
    //     size: [120, 40],
    //     labelCfg: {
    //       style: {
    //         textAlign: 'center',
    //       },
    //     },
    //   };
    // });
    // 初始化的图数据
    graph.data(data);
    // 根据提供的数据渲染视图。
    graph.render();
    // 让画布内容适应视口
    graph.fitView();
  }


  useEffect(() => {
    initTree(data)
  }, [])
  return <div className={styles.relationship}>
    <div id="container" />
  </div>;
};
