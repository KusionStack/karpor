import React, { useEffect, useRef } from 'react'
import type { LegacyRef } from 'react'
import hljs from 'highlight.js'
import yaml from 'js-yaml'
import 'highlight.js/styles/lightfair.css'
import { yaml2json } from '../../utils/tools'

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

  return (
    <div className={styles.yaml_content} style={{ height: props?.height }}>
      <div
        className={styles.yaml_box}
        style={{ height: props?.height }}
        ref={yamlRef as any}
      />
    </div>
  )
}

export default Yaml
