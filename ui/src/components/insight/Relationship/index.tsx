import { useEffect } from "react";
import { Space } from "antd";
import G6 from "@antv/g6";
import kubernetes from "../../../assets/kubernetes.png";

import styles from "./styles.module.less";

const Relationship = () => {
  function draw() {
    const data = {
      nodes: [
        {
          id: "Service", // String，可选，节点的唯一标识
          label: "Service", // String，节点标签
        },
        {
          id: "Deployment", // String，节点的唯一标识
          label: "Deployment", // String，节点标签
        },
        {
          id: "Endpoints", // String，节点的唯一标识
          label: "Endpoints", // String，节点标签
        },
        {
          id: "ReplicaSet", // String，节点的唯一标识
          label: "ReplicaSet", // String，节点标签
        },
        {
          id: "StatefulSet", // String，节点的唯一标识
          label: "StatefulSet", // String，节点标签
        },
        {
          id: "OperationJob", // String，节点的唯一标识
          label: "OperationJob", // String，节点标签
        },
        {
          id: "Pod", // String，节点的唯一标识
          label: "Pod", // String，节点标签
        },
        {
          id: "ResourceDecoration", // String，节点的唯一标识
          label: "ResourceDecoration", // String，节点标签
        },
        {
          id: "InPlaceSet", // String，节点的唯一标识
          label: "InPlaceSet", // String，节点标签
        },
        {
          id: "RuleSet", // String，节点的唯一标识
          label: "RuleSet", // String，节点标签
        },
        {
          id: "CafeDeployment", // String，节点的唯一标识
          label: "CafeDeployment", // String，节点标签
        },
      ],
      // 边
      edges: [
        {
          source: "Service", // String，必须，起始节点 id
          target: "Endpoints", // String，必须，目标节点 id
        },
        {
          source: "Deployment", // String，必须，起始节点 id
          target: "ReplicaSet", // String，必须，目标节点 id
        },
        {
          source: "Endpoints", // String，必须，起始节点 id
          target: "Pod", // String，必须，目标节点 id
        },
        {
          source: "ReplicaSet", // String，必须，起始节点 id
          target: "Pod", // String，必须，目标节点 id
        },
        {
          source: "StatefulSet", // String，必须，起始节点 id
          target: "Pod", // String，必须，目标节点 id
        },
        {
          source: "OperationJob", // String，必须，起始节点 id
          target: "Pod", // String，必须，目标节点 id
        },
        {
          source: "ResourceDecoration", // String，必须，起始节点 id
          target: "Pod", // String，必须，目标节点 id
        },
        {
          source: "InPlaceSet", // String，必须，起始节点 id
          target: "Pod", // String，必须，目标节点 id
        },
        {
          source: "CafeDeployment", // String，必须，起始节点 id
          target: "InPlaceSet", // String，必须，目标节点 id
        },
        {
          source: "RuleSet", // String，必须，起始节点 id
          target: "Pod", // String，必须，目标节点 id
        },
      ],
    };

    G6.registerNode(
      "card-node",
      {
        options: {
          stroke: "#e6fffb",
          fill: "#e6fffb",
        },
        draw(cfg, group) {
          const styles = this.getShapeStyle(cfg);
          const { labelCfg = {} } = cfg;

          const w = styles.width;
          const h = styles.height;

          const keyShape = group.addShape("rect", {
            attrs: {
              ...styles,
              width: 160,
              height: 40,
              x: -w / 2,
              y: -h / 2,
              fill: "#fff",
              stroke: cfg.label === "Pod" ? "#1890ff" : "#a2a2a2",
              lineWidth: cfg.label === "Pod" ? 2 : 1,
              radius: 4,
            },
          });

          group.addShape("image", {
            attrs: {
              x: 8 - w / 2,
              y: 10 - h / 2,
              width: 20,
              height: 20,
              // img: "https://gw.alipayobjects.com/mdn/rms_f8c6a0/afts/img/A*Q_FQT6nwEC8AAAAAAAAAAABkARQnAQ",
              img: kubernetes,
            },
            // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
            name: "image-shape",
          });

          if (cfg.label) {
            group.addShape("text", {
              attrs: {
                ...labelCfg.style,
                text: cfg.label,
                x: 35 - w / 2,
                y: 28 - h / 2,
                fill: cfg.label === "Pod" ? "#1890ff" : "#333",
              },
            });
          }

          return keyShape;
        },
        update: undefined,
      },
      "rect"
    );

    const container = document.getElementById("container");
    const width = container.scrollWidth;
    const height = container.scrollHeight || 500;
    const graph = new G6.Graph({
      container: "container",
      width,
      height,
      fitView: true,
      modes: {
        default: ["drag-canvas", "drag-node", "zoom-canvas"],
      },
      layout: {
        type: "dagre",
        rankdir: "LR", // 'TB' / 'BT' / 'LR' / 'RL'
        align: "UL", // 'UL' / 'UR' / 'DL' / 'DR' /
        controlPoints: true,
        nodesep: 10,
        ranksep: 10,
        // nodesepFunc: () => 10,
        // ranksepFunc: () => 10,
      },
      defaultNode: {
        size: [200, 40],
        type: "card-node",
        style: {
          fill: "#fff",
          stroke: "#a2a2a2",
          radius: 5,
        },
        labelCfg: {
          style: {
            fill: "#000",
            fontSize: 14,
          },
        },
      },
      defaultEdge: {
        type: "polyline",
        size: 1,
        color: "#a2a2a2",
        style: {
          endArrow: {
            // path: "M 0,0 L 6,2 L 6,-2 Z",
            // path: 'M 0,0 L 6, 1 L 4,0 L 6, -2 Z',
            path: "M 0,0 L 8, 4 L 0,0 L 8, -4 Z",
            fill: "#a2a2a2",
            d: 0,
          },
          // radius: 20,
        },
      },
    });
    graph.data(data);
    graph.zoom(1.5);
    graph.render();

    graph.on("node:mouseenter", (evt) => {
      console.log(evt, "-===mouseenter===");
    });

    graph.on("node:mouseleave", (evt) => {
      console.log(evt, "-===mouseleave===");
    });

    graph.on("node:click", (evt) => {
      console.log(evt, "-===click===");
    });

    if (typeof window !== "undefined") {
      window.onresize = () => {
        if (!graph || graph.get("destroyed")) return;
        if (!container || !container.scrollWidth || !container.scrollHeight)
          return;
        graph.changeSize(container.scrollWidth, container.scrollHeight);
      };
    }
  }

  useEffect(() => {
    draw();
  }, []);

  return (
    <div className={styles.relationship}>
      <Space className={styles.title}>
        <div className={styles.green}></div>
        <div>正常</div>
        <div className={styles.blue} />
        <div>变更中</div>
        <div className={styles.red} />
        <div>异常</div>
        <div className={styles.black} />
        <div>删除中</div>
      </Space>
      <div id="container" style={{ height: 380 }}></div>
    </div>
  );
};

export default Relationship;
