import React from '@alipay/bigfish/react';
import { styled } from '@alipay/bigfish';
import styles from "./style.less";


const Wrapper = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
`;

type Props = {
  current: String;
  list: Array<{ label: String, value: String }>;
  onChange: (val: String) => void
}


const Tabs = ({ current, list, onChange }: Props) => {
  return <Wrapper>
    {
      list?.map(item => {
        return <div className={styles.item} key={item.value as React.Key} onClick={() => onChange(item.value)}>
          <div className={`${styles.normal} ${current === item.value ? styles.active : ''}`}>{item.label}</div>
        </div>
      })
    }
  </Wrapper>
}

export default Tabs;