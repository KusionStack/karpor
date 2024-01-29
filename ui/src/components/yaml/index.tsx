import React, { useEffect, useRef } from 'react'
import type { LegacyRef } from 'react'
import hljs from 'highlight.js'
import yaml from 'js-yaml'
import 'highlight.js/styles/lightfair.css'
import { yaml2json } from '@/utils/tools'

import styles from './styles.module.less'

// eslint-disable-next-line @typescript-eslint/no-var-requires
hljs.registerLanguage('yaml', require('highlight.js/lib/languages/yaml'))

type IProps = {
  data: any
  height?: string | number
}

const Yaml = (props: IProps) => {
  const yamlRef = useRef<LegacyRef<HTMLDivElement> | undefined>()
  const { data } = props
  useEffect(() => {
    const yamlStatusJson = yaml2json(data)
    if (yamlRef.current && yamlStatusJson?.data) {
      ;(yamlRef.current as unknown as HTMLElement).innerHTML = hljs.highlight(
        'yaml',
        yaml.dump(yamlStatusJson?.data),
      ).value
    }
  }, [data])

  // function copy() {
  //   const textarea = document.createElement('textarea');
  //   textarea.value = data;
  //   document.body.appendChild(textarea);
  //   textarea.select();
  //   document.execCommand('copy');
  //   message.success('复制成功');
  //   document.body.removeChild(textarea);
  // }

  return (
    <div className={styles.yamlContent} style={{ height: props?.height }}>
      {/* <div className={styles.copyContainer}>
        {data && (
          <Button type="primary" size="small" onClick={copy} disabled={!data}>
            复制
          </Button>
        )}
      </div> */}
      <div
        className={styles.yamlBox}
        style={{ height: props?.height }}
        ref={yamlRef as any}
      />
    </div>
  )
}

export default Yaml
