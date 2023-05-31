import React from "react";
import {
  Chart,
  Interval,
  Axis,
  Tooltip,
  Coordinate,
  Legend,
  View,
  Annotation,
  getTheme,
} from "bizcharts";

function Ring({
  data = [],
  content = { percent: "100%" },
  intervalConfig = {},
}) {
  const brandFill = getTheme().colors10[0];
  return (
    <Chart placeholder={false} height={200} padding="auto" autoFit>
      <Legend visible={false} />
      {/* 绘制图形 */}
      <View
        data={data}
        scale={{
          percent: {
            formatter: (val) => {
              return (val * 100).toFixed(2) + "%";
            },
          },
        }}
      >
        <Coordinate type="theta" innerRadius={0.75} />
        <Interval
          position="percent"
          adjust="stack"
          // color="type"
          // color={["type", ["rgba(100, 100, 255, 0.6)", "#eee"]]}
          // color={["type", [brandFill, "#eee"]]}
          color={["type", ['#000', "#eee"]]}
          size={16}
          // style={{ fillOpacity: 0.6 }}
          // label={['type', {offset: 40}]}
          {...intervalConfig}
        />
        <Annotation.Text
          position={["50%", "50%"]}
          content={content?.percent}
          style={{
            lineHeight: 240,
            fontSize: 24,
            fill: '#000',
            textAlign: "center",
          }}
        />
      </View>
    </Chart>
  );
}

const myData = [
  { type: "已完成", percent: 1 },
  { type: "待完成", percent: 0 },
];
const myContent = {
  percent: "100%",
};

export default function Progress() {
  return (
    <div style={{ padding: 10, width: 200, height: 200 }}>
      <Ring data={myData} content={myContent as any} />
    </div>
  );
}
