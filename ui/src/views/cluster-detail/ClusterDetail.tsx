import { useEffect, memo, useState } from "react";
import { useLocation } from "react-router-dom";
import { Breadcrumb } from "antd";
import queryString from "query-string";
import PiePercent from "../../components/insight/Score";
import IssueTable from "../../components/insight/Issue";
import Overview from "../../components/insight/Overview";
import Stat from "../../components/insight/Stat";
import YamlDrawer from "../../components/insight/YamlDrawer";
import { yamlStr } from "../../utils/mock";

import styles from "./styles.module.less";

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
      title: "1.22.1",
      desc: "Kubernetes Version",
    },
    {
      title: "3",
      desc: "Node",
    },
    {
      title: "2023-05-16T02:46:57Z",
      desc: "CreateTime",
    },
  ],
};

const ClusterDetail = () => {
  const [issueData, setIssueData] = useState<any[]>([]);
  const [overviewData, setOverviewData] = useState<any>();
  const location = useLocation();
  const urlSearchParams = queryString.parse(location.search);
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
        title: "Cluster",
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

  const MemoPiePercent = memo(PiePercent);
  const MemoIssueTable = memo(IssueTable);
  const MemoOverview = memo(Overview);
  const MemoStat = memo(Stat);
  const MemoYamlDrawer = memo(YamlDrawer);

  return (
    <div className={styles.container}>
      <div className={styles.bread}>
        <Breadcrumb separator=">" items={breadcrumbList} />
      </div>
      <div className={styles.top}>
        <div className={styles.left}>
          <div className={styles.title}>Info</div>
          <MemoStat />
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

export default ClusterDetail;
