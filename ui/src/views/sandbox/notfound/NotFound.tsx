/*
 * Copyright 2017 The Karbour Authors.
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
