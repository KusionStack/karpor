import React, { useEffect, useRef, useState } from 'react'
import type { LegacyRef } from 'react'
import { Button, message } from 'antd'
import { Resizable } from 're-resizable'
import { useTranslation } from 'react-i18next'
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
  const { t } = useTranslation()
  const yamlRef = useRef<LegacyRef<HTMLDivElement> | undefined>()
  const { data } = props
  const [moduleHeight, setModuleHeight] = useState<number>(500)

  useEffect(() => {
    const yamlStatusJson = yaml2json(data)
    if (yamlRef.current && yamlStatusJson?.data) {
      ;(yamlRef.current as unknown as HTMLElement).innerHTML = hljs.highlight(
        'yaml',
        yaml.dump(yamlStatusJson?.data),
      ).value
    }
  }, [data])

  function copy() {
    const textarea = document.createElement('textarea')
    textarea.value = data
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    message.success(t('CopySuccess'))
    document.body.removeChild(textarea)
  }

  return (
    <div style={{ paddingBottom: 20 }}>
      <Resizable
        defaultSize={{
          height: moduleHeight,
        }}
        onResizeStop={(e, direction, ref, d) => {
          const newModuleHeight = moduleHeight + d.height
          setModuleHeight(newModuleHeight)
        }}
      >
        <div className={styles.yaml_content} style={{ height: props?.height }}>
          <div className={styles.copy}>
            {data && (
              <Button
                type="primary"
                size="small"
                onClick={copy}
                disabled={!data}
              >
                {t('Copy')}
              </Button>
            )}
          </div>
          <div className={styles.yaml_container}>
            <div
              className={styles.yaml_box}
              style={{ height: props?.height }}
              ref={yamlRef as any}
            />
          </div>
        </div>
      </Resizable>
    </div>
  )
}

export default Yaml
