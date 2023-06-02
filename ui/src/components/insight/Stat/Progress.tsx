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

import { useEffect } from "react";
import { Chart } from "@antv/g2";

export default function Progress() {
  const renderChart = () => {
    const tempNode = document.getElementById("statPie");
    const child = tempNode?.lastElementChild;
    if (child) {
      // 这里防止饼图重复渲染 但是直接操作了dom
      // 如果有更好的方法 可以换掉
      tempNode.removeChild(child);
    }
    const data = [
      { type: "已完成", value: 1 },
      { type: "待完成", value: 0 },
    ];

    // Step 1: 创建 Chart 对象
    const chart: any = new Chart({
      container: "statPie",
      theme: "classic",
      width: 200,
      height: 200,
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
      .scale("color", {
        range: ["#000", "#cd201f"],
      })
      // .label({ text: "type", style: { fontSize: 10, fontWeight: "bold" } })
      // .label({
      //   text: (d, i, data) => (i < data.length - 3 ? d.value : ""),
      //   style: {
      //     fontSize: 9,
      //     dy: 12,
      //   },
      // })
      // .tooltip({
      //   title: "type",
      //   items: ["value"],
      // })
      // .tooltip(['type', 'value'])
      .tooltip({
        title: { channel: "x" },
        items: [{ channel: "y" }],
      })
      .tooltip({
        items: [{ channel: "y", valueFormatter: ".0%" }],
      })
      // .interaction("tooltip", {
      //   // render 回调方法返回一个innerHTML 或者 ReactNode
      //   render: (event, { title, items }) => {
      //     console.log(title, items, "====sdasda");
      //     return `<div>Your custom render content here.</div>`;
      //   },
      // })
      .animate("enter", { type: "waveIn" })
      .legend(false);
    chart.text().style({
      x: "50%", // 百分比
      y: "50%", // 百分比
      text: "100%",
      textAlign: "center",
      fontSize: 32,
      textBaseline: "middle",
      fontWeight: "bold",
    });
    chart.render();
  };

  useEffect(() => {
    renderChart();
  }, []);
  return (
    <div id="statPie" style={{ padding: 10, width: 200, height: 200 }}>
    </div>
  );
}
