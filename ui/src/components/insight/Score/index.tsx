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

import type { FC } from "react";
import { useEffect, useRef } from "react";
import {
  Chart,
  registerShape,
  Geom,
  Axis,
  Tooltip,
  Interval,
  Interaction,
  Coordinate,
  Annotation,
  Legend
} from "bizcharts";

import styles from "./styles.module.scss";

const PiePercent: FC = () => {
  const data = [
    {
      type: "Low",
      value: 90,
    },
    {
      type: "Medium",
      value: 2,
    },
    {
      type: "High",
      value: 8,
    },
  ]; // 可以通过调整这个数值控制分割空白处的间距，0-1 之间的数值

  const sliceNumber = 0.01; // 自定义 other 的图形，增加两条线

  registerShape("interval", "sliceShape", {
    draw(cfg: any, container: any) {
      const points = cfg.points;
      let path = [];
      path.push(["M", points[0].x, points[0].y]);
      path.push(["L", points[1].x, points[1].y - sliceNumber]);
      path.push(["L", points[2].x, points[2].y - sliceNumber]);
      path.push(["L", points[3].x, points[3].y]);
      path.push("Z");
      path = (this as any).parsePath(path);
      return container.addShape("path", {
        attrs: {
          fill: cfg.color,
          path: path,
        },
      });
    },
  });

  return (
    <div className={styles.score}>
      <Chart data={data} height={300} autoFit>
        <Coordinate type="theta" radius={0.8} innerRadius={0.8} />
        <Axis visible={false} />
        <Tooltip showTitle={false} />
        <Interval
          adjust="stack"
          position="value"
          color="type"
          shape="sliceShape"
        />
        <Legend visible={false} />
        <Interaction type="element-single-selected" />
        <Annotation.Text
          position={["50%", "50%"]}
          content="B-"
          style={{
            lineHeight: 240,
            fontSize: 60,
            fill: "#262626",
            textAlign: "center",
          }}
        />
      </Chart>
      <div className={styles.desc}>
        Health score is the ratio of passing to failing Action Items, weighted by severity
      </div>
    </div>
  );
};

export default PiePercent;
