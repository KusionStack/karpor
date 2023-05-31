import React from "react";
import { Drawer, Button } from "antd";
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
