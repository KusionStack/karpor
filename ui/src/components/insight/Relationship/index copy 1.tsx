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

  function draw2() {
    const data = {
      id: 'root',
      label: 'root',
      children: [
        {
          id: 'c1',
          label: 'c1',
          children: [
            {
              id: 'c1-1',
              label: 'c1-1',
            },
            {
              id: 'c1-2',
              label: 'c1-2',
              children: [
                {
                  id: 'c1-2-1',
                  label: 'c1-2-1',
                },
                {
                  id: 'c1-2-2',
                  label: 'c1-2-2',
                },
              ],
            },
          ],
        },
        {
          id: 'c2',
          label: 'c2',
        },
        {
          id: 'c3',
          label: 'c3',
          children: [
            {
              id: 'c3-1',
              label: 'c3-1',
            },
            {
              id: 'c3-2',
              label: 'c3-2',
              children: [
                {
                  id: 'c3-2-1',
                  label: 'c3-2-1',
                },
                {
                  id: 'c3-2-2',
                  label: 'c3-2-2',
                },
                {
                  id: 'c3-2-3',
                  label: 'c3-2-3',
                },
              ],
            },
            {
              id: 'c3-3',
              label: 'c3-3',
            },
          ],
        },
      ],
    };

    G6.Util.traverseTree(data, (d: any) => {
      d.leftIcon = {
        style: {
          fill: '#fafafa',
          stroke: '#fafafa',
        },
        img: 'https://gw.alipayobjects.com/mdn/rms_f8c6a0/afts/img/A*Q_FQT6nwEC8AAAAAAAAAAABkARQnAQ',
      };
      return true;
    });

    G6.registerNode(
      'icon-node',
      {
        options: {
          size: [60, 20],
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

          /**
           * leftIcon 格式如下：
           *  {
           *    style: ShapeStyle;
           *    img: ''
           *  }
           */
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
                width: 24,
                height: 24,
                img: img || 'https://g.alicdn.com/cm-design/arms-trace/1.0.155/styles/armsTrace/images/TAIR.png',
              },
              // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
              name: 'image-shape',
            });
          }

          // 如果不需要动态增加或删除元素，则不需要 add 这两个 marker
          // group.addShape('marker', {
          //   attrs: {
          //     x: 40 - w / 2,
          //     y: 52 - h / 2,
          //     r: 6,
          //     stroke: '#73d13d',
          //     cursor: 'pointer',
          //     symbol: EXPAND_ICON,
          //   },
          //   // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
          //   name: 'add-item',
          // });

          // group.addShape('marker', {
          //   attrs: {
          //     x: 80 - w / 2,
          //     y: 52 - h / 2,
          //     r: 6,
          //     stroke: '#ff4d4f',
          //     cursor: 'pointer',
          //     symbol: COLLAPSE_ICON,
          //   },
          //   // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
          //   name: 'remove-item',
          // });

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

    G6.registerEdge('flow-line', {
      draw(cfg: any, group: any) {
        const startPoint = cfg.startPoint;
        const endPoint = cfg.endPoint;

        const { style } = cfg;
        const shape = group.addShape('path', {
          attrs: {
            stroke: style.stroke,
            endArrow: style.endArrow,
            // path: [
            //   ['M', startPoint.x, startPoint.y],
            //   ['L', startPoint.x, (startPoint.y + endPoint.y) / 2],
            //   ['L', endPoint.x, (startPoint.y + endPoint.y) / 2],
            //   ['L', endPoint.x, endPoint.y],
            // ],
            path: [
              ['M', startPoint.x, startPoint.y],
              ['L', startPoint.x, (startPoint.y + endPoint.y) / 2],
              ['L', endPoint.x, (startPoint.y + endPoint.y) / 2],
              ['L', endPoint.x, endPoint.y],
            ],
          },
        });

        return shape;
      },
    });

    const defaultStateStyles = {
      hover: {
        stroke: '#d9d9d9',
        lineWidth: 2,
      },
    };

    const defaultNodeStyle = {
      fill: '#fff',
      stroke: '#d9d9d9',
      radius: 2,
    };

    const defaultEdgeStyle = {
      stroke: '#d9d9d9',
      endArrow: {
        path: 'M 0,0 L 12, 6 L 9,0 L 12, -6 Z',
        fill: '#d9d9d9',
        d: -20,
      },
    };

    const defaultLayout = {
      type: 'compactBox',
      direction: 'LR',
      nodesep: 10, // 节点层距离（上下）
      ranksep: 20, // 节点之间距离（左右）
      controlPoints: true,
      getId: function getId(d: any) {
        return d.id;
      },
      getHeight: function getHeight() {
        return 16;
      },
      getWidth: function getWidth() {
        return 16;
      },
      getVGap: function getVGap() {
        return 60;
      },
      getHGap: function getHGap() {
        return 70;
      },
    };

    const defaultLabelCfg = {
      style: {
        fill: '#000',
        fontSize: 16,
      },
    };

    const container = document.getElementById('g6Container')!;
    const width = container.scrollWidth;
    const height = container.scrollHeight || 500;

    // const minimap = new G6.Minimap({
    //   size: [150, 100],
    // });
    const graph = new G6.TreeGraph({
      container: 'g6Container',
      width,
      height,
      linkCenter: true,
      // plugins: [minimap],
      modes: {
        default: ['drag-canvas', 'zoom-canvas'],
      },
      defaultNode: {
        type: 'icon-node',
        size: [120, 40],
        anchorPoints: [
          [0, 0.5],
          [1, 0.5],
        ],
        style: defaultNodeStyle,
        labelCfg: defaultLabelCfg,
      },
      defaultEdge: {
        type: 'flow-line',
        style: defaultEdgeStyle,
      },
      nodeStateStyles: defaultStateStyles,
      edgeStateStyles: defaultStateStyles,
      layout: defaultLayout,
    });

    graph.data(data);
    graph.render();
    graph.fitView();

    graph.on('node:mouseenter', (evt: any) => {
      const { item } = evt;
      graph.setItemState(item, 'hover', true);
    });

    graph.on('node:mouseleave', (evt: any) => {
      const { item } = evt;
      graph.setItemState(item, 'hover', false);
    });

    if (typeof window !== 'undefined')
      window.onresize = () => {
        if (!graph || graph.get('destroyed')) return;
        if (!container || !container.scrollWidth || !container.scrollHeight) return;
        graph.changeSize(container.scrollWidth, container.scrollHeight);
      };
  }
  useEffect(() => {
    // draw();
    draw2();
  }, [])
  return <div className={styles.relationship}>
    <div id="g6Container" />
  </div>;
};
