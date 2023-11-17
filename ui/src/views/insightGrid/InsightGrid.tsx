import { useEffect, memo, useState, useRef } from "react";
import { useLocation } from "react-router-dom";
import { Breadcrumb } from "antd";
import ReactGridLayout from "react-grid-layout";
import queryString from "query-string";
import PiePercent from "../../components/insight/Score";
import IssueTable from "../../components/insight/Issue";
import Overview from "../../components/insight/Overview";
import Relationship from "../../components/insight/Relationship";
import YamlDrawer from "../../components/insight/YamlDrawer";
import { insightModuleGrid } from "../../utils/constants";
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
  const [width, setWidth] = useState(1200);
  const [isEdit, setIsEdit] = useState(false);
  const localGrid = localStorage.getItem("gridModule");
  const [gridModule, setGridModule] = useState(
    localGrid ? JSON.parse(localGrid) : insightModuleGrid
  );
  const [layout, setLayout] = useState([]);

  const containerRef = useRef<HTMLDivElement>();

  useEffect(() => {
    if (containerRef?.current) {
      setWidth(containerRef?.current?.offsetWidth);
    }
  }, [containerRef]);

  function resizeFunc() {
    if (containerRef?.current) {
      setWidth(containerRef.current?.offsetWidth);
    }
  }

  useEffect(() => {
    window.addEventListener("resize", resizeFunc);
    return () => {
      window.removeEventListener("resize", resizeFunc);
    };
  }, []);

  function onLayoutChange(currentLayout) {
    setLayout(currentLayout);
  }

  function onLayout() {
    const resetLayout = gridModule?.map((item) => {
      return {
        i: item.key,
        key: item.key,
        isDraggable: isEdit,
        isResizable: isEdit,
        ...(item?.config ?? {}),
      };
    });
    return resetLayout;
  }

  function handleEditGrid() {
    if (isEdit) {
      const tmp = layout?.map((item) => {
        return {
          key: item?.i,
          config: {
            ...item,
            isDraggable: false,
            isResizable: false,
          },
        };
      });
      setGridModule(tmp);
      localStorage.setItem("gridModule", JSON.stringify(tmp));
    } else {
      const tmp = layout?.map((item) => {
        return {
          key: item?.i,
          config: {
            ...item,
            isDraggable: true,
            isResizable: true,
          },
        };
      });
      setGridModule(tmp);
    }
    setIsEdit(!isEdit);
  }

  function handleClickYaml() {
    setVisible(true);
  }

  function handleClose() {
    setVisible(false);
  }

  const MemoPiePercent = memo(PiePercent);
  const MemoIssueTable = memo(IssueTable);
  const MemoOverview = memo(Overview);
  const MemoYamlDrawer = memo(YamlDrawer);
  const MemoRelationship = memo(Relationship);

  const blockMap = {
    relationship: <MemoRelationship />,
    overview: (
      <MemoOverview data={overviewData} handleClick={handleClickYaml} />
    ),
    issue: <MemoIssueTable data={issueData} handleSearch={handleSearch} />,
    score: <MemoPiePercent />,
  };

  return (
    <div className={styles.container} ref={containerRef}>
      <div className={styles.bread}>
        <Breadcrumb separator=">" items={breadcrumbList} />
        <div className={styles.editBtn} onClick={handleEditGrid}>
          {!isEdit ? "编辑布局" : "保存布局"}
        </div>
      </div>
      <ReactGridLayout
        className="react-grid-layout"
        layout={onLayout()}
        autoSize
        cols={24}
        rowHeight={100}
        margin={[14, 14]}
        containerPadding={[0, 0]}
        width={width}
        onLayoutChange={onLayoutChange}
      >
        {(gridModule as any[])?.map((item) => {
          return (
            <div className={styles.item} key={item.key}>
              <div className={styles.title}>{item.key}</div>
              {blockMap[item.key]}
            </div>
          );
        })}
      </ReactGridLayout>
      <MemoYamlDrawer open={visible} onClose={handleClose} data={yamlStr} />
    </div>
  );
};

export default Insight;
