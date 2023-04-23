import { Button, Result } from "antd";

import styles from "./styles.module.scss";
import { useNavigate } from "react-router-dom";

export default function Insight() {
  const navigate = useNavigate();
  function goBack() {
    navigate("/search");
  }

  return (
    <div className={styles.container}>
      <Result
        status="404"
        title="404"
        subTitle="对不起，你访问的页面不存在~"
        extra={
          <Button type="primary" onClick={goBack}>
            返回首页
          </Button>
        }
      />
    </div>
  );
}
