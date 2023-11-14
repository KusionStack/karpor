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
import { useEffect } from "react";

import { Chart } from "@antv/g2";

import styles from "./styles.module.scss";

const PiePercent: FC = () => {

  const renderChart = () => {
    const tempNode = document.getElementById("containerPie");
    const child = tempNode?.lastElementChild;
    if (child) {
      // 这里我做了一个判断  防止饼图重复渲染 但是直接操作了dom
      // 如果有更好的方法 可以换掉
      tempNode.removeChild(child);
    }
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

    // Step 1: 创建 Chart 对象
    const chart: any = new Chart({
      container: "containerPie",
      theme: "classic",
      width: 300,
      height: 300,
    });

    chart.coordinate({ type: "theta", innerRadius: 0.8 });
    // Step 2: 载入数据源
    chart.data(data);
    chart
      .interval()
      .transform({ type: "stackY" })
      .data(data)
      .encode("y", "value")
      .encode("color", "type")
      .style("stroke", "white")
      .style("inset", 1)
      .style("radius", 1)
      // .scale("color", {
      //   range: ["#3b5999", "#f4a460", "#cd201f"],
      // })
      // .label({ text: "type", style: { fontSize: 10, fontWeight: "bold" } })
      // .label({
      //   text: (d, i, data) => (i < data.length - 3 ? d.value : ""),
      //   style: {
      //     fontSize: 9,
      //     dy: 12,
      //   },
      // })
      .tooltip({
        title: "type",
        items: ["value"],
      })
      .interaction("tooltip", {
        // 设置 Tooltip 的位置
        position: "right",
      })
      .interaction("tooltip", {
        // render 回调方法返回一个innerHTML 或者 ReactNode
        render: (event, { title, items }) => {
          return `<div style="color:${items?.[0]?.color}">${title}:&nbsp;&nbsp;&nbsp;${items?.[0]?.value}</div>`;
        },
      })
      .animate("enter", { type: "waveIn" })
      .legend(false);
    chart.text().style({
      x: "50%", // 百分比
      y: "50%", // 百分比
      text: "B-",
      textAlign: "center",
      fontSize: 60,
      textBaseline: "middle",
    });
    chart.render();
  };

  useEffect(() => {
    renderChart();
  }, []);

  return (
    <div className={styles.score}>
      <div
        id="containerPie"
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
        }}
      ></div>
      <div className={styles.desc}>
        Health score is the ratio of passing to failing Action Items, weighted
        by severity
      </div>
    </div>
  );
};

export default PiePercent;
