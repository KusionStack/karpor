import { useEffect, memo, useState } from "react";
import { useLocation } from "react-router-dom";
import { Breadcrumb } from "antd";
import queryString from "query-string";
import PiePercent from "../../components/insight/Score";
import IssueTable from "../../components/insight/Issue";
import Overview from "../../components/insight/Overview";
import Relationship from "../../components/insight/Relationship";
import YamlDrawer from "../../components/insight/YamlDrawer";
import { yamlStr } from "../../utils/mock";

import styles from "./styles.module.scss";

const MemoPiePercent = memo(PiePercent);
const MemoIssueTable = memo(IssueTable);
const MemoOverview = memo(Overview);
const MemoYamlDrawer = memo(YamlDrawer);
const MemoRelationship = memo(Relationship);

const relationshipData = {
  nodes: [
    {
      id: "node-0",
      shape: "data-processing-dag-node",
      x: 0,
      y: 100,
      ports: [
        {
          id: "node-0-out",
          group: "out",
        },
      ],
      data: {
        name: "数据输入_1",
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
        name: "数据筛选_1",
        type: "FILTER",
      },
    },
    {
      id: "node-2",
      shape: "data-processing-dag-node",
      x: 250,
      y: 200,
      ports: [
        {
          id: "node-2-out",
          group: "out",
        },
      ],
      data: {
        name: "数据输入_2",
        type: "INPUT",
      },
    },
    {
      id: "node-3",
      shape: "data-processing-dag-node",
      x: 500,
      y: 100,
      ports: [
        {
          id: "node-3-in",
          group: "in",
        },
        {
          id: "node-3-out",
          group: "out",
        },
      ],
      data: {
        name: "数据连接_1",
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
      ],
      data: {
        name: "数据输出_1",
        type: "OUTPUT",
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
        cell: "node-2",
        port: "node-2-out",
      },
      target: {
        cell: "node-3",
        port: "node-3-in",
      },
      shape: "data-processing-curve",
      zIndex: -1,
      data: {
        source: "node-2",
        target: "node-3",
      },
    },
    {
      id: "edge-2",
      source: {
        cell: "node-1",
        port: "node-1-out",
      },
      target: {
        cell: "node-3",
        port: "node-3-in",
      },
      shape: "data-processing-curve",
      zIndex: -1,
      data: {
        source: "node-1",
        target: "node-3",
      },
    },
    {
      id: "edge-3",
      source: {
        cell: "node-3",
        port: "node-3-out",
      },
      target: {
        cell: "node-4",
        port: "node-4-in",
      },
      shape: "data-processing-curve",
      zIndex: -1,
      data: {
        source: "node-3",
        target: "node-4",
      },
    },
  ],
};

const issueDataMock = [
  {
    id: 0,
    title: "Test11",
    severity: "low",
    labels: "SOLUTION,CIS",
  },
  {
    id: 1,
    title: "ABCDD",
    severity: "Medium",
    labels: "CVE,SOLUTION",
  },
  {
    id: 2,
    title: "fffAAAAFGFGGEE",
    severity: "High",
    labels: "CIS,SOLUTION",
  },
  {
    id: 3,
    title: "asdadasdas",
    severity: "Low",
    labels: "kubeaudit",
  },
];

const overviewDataMock = {
  title: "Cluster Version",
  list: [
    {
      title: "cluster1",
      desc: "Cluster",
    },
    {
      title: "karbour-7fb8fd54cf-m64xr",
      desc: "Name",
    },
    {
      title: "elastic",
      desc: "Namespac",
    },
    {
      title: "Pod",
      desc: "Kind",
    },
    {
      title: "v1",
      desc: "ApiVersion",
    },
    {
      title: "2023-05-16T02:46:57Z",
      desc: "CreateTime",
    },
  ],
};

const Insight = () => {
  const location = useLocation();
  const urlSearchParams = queryString.parse(location.search);
  const [issueData, setIssueData] = useState<any[]>([]);
  const [overviewData, setOverviewData] = useState<any>();
  console.log(urlSearchParams, "===urlSearchParams===");
  const breadcrumbList = (urlSearchParams?.query as string)
    ?.split(",")
    ?.map((item) => {
      return { title: item };
    }) || [
    {
      title: "K8S",
    },
    {
      title: "Api Version",
    },
    {
      title: "Kind",
    },
    {
      title: "Pod",
    },
  ];

  useEffect(() => {
    setIssueData(issueDataMock);
    setOverviewData(overviewDataMock);
  }, []);

  function handleSearch(value: string) {
    if (!value) {
      setIssueData(issueDataMock);
    } else {
      const tmp = issueData?.filter((item) =>
        item?.title?.toLowerCase()?.includes(value?.toLowerCase())
      );
      setIssueData(tmp);
    }
  }

  const [visible, setVisible] = useState(false);

  function handleClickYaml() {
    setVisible(true);
  }

  function handleClose() {
    setVisible(false);
  }

  return (
    <div className={styles.container}>
      <div className={styles.bread}>
        <Breadcrumb separator=">" items={breadcrumbList} />
      </div>
      <div className={styles.top}>
        <div className={styles.left}>
          <div className={styles.title}>Relationship</div>
          <MemoRelationship />
        </div>
        <div className={styles.right}>
          <div className={styles.title}>Overview</div>
          <MemoOverview data={overviewData} handleClick={handleClickYaml} />
        </div>
      </div>
      <div className={styles.bottom}>
        <div className={styles.left}>
          <div className={styles.title}>Issue</div>
          <MemoIssueTable data={issueData} handleSearch={handleSearch} />
        </div>
        <div className={styles.right}>
          <div className={styles.title}>Score</div>
          <MemoPiePercent />
        </div>
      </div>
      <MemoYamlDrawer open={visible} onClose={handleClose} data={yamlStr} />
    </div>
  );
};

export default Insight;
