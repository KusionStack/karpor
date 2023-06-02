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

import { Drawer } from "antd";
import { default as AnsiUp } from "ansi_up";
import styles from "./styles.module.scss";

type IProps = {
  open: boolean;
  onClose: () => void;
  data: string;
}

export default function YamlDrawer({ open, onClose, data }: IProps) {
  const ansi_up = new AnsiUp();
  const info = ansi_up.ansi_to_html(
    data?.trim()
  )
  return <Drawer width={800} title="Yaml 详情" placement="right" onClose={onClose} open={open}>
    <div
      className={styles.yaml}
      dangerouslySetInnerHTML={{
        __html: info,
      }}
    />
  </Drawer>
}
