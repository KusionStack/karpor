import React, { useEffect, useRef, useState, memo } from 'react'
import { EditorState } from '@codemirror/state'
import {
  EditorView,
  keymap,
  KeyBinding,
  highlightSpecialChars,
  ViewPlugin,
} from '@codemirror/view'
import { sql } from '@codemirror/lang-sql'
import {
  HighlightStyle,
  syntaxHighlighting,
  LanguageSupport,
} from '@codemirror/language'
import { tags } from '@lezer/highlight'
import {
  autocompletion,
  completionKeymap,
  acceptCompletion,
  completeFromList,
  currentCompletions,
  startCompletion,
  completionStatus,
} from '@codemirror/autocomplete'
import { useTranslation } from 'react-i18next'
import { css } from '@emotion/css'
import { Divider, message } from 'antd'
import arrowRight from '@/assets/arrow-right.png'
import {
  defaultKeywords,
  whereKeywords,
  kindCompletions,
  operatorKeywords,
  searchSqlPrefix,
} from '@/utils/constants'
import { useAxios } from '@/utils/request'

import styles from './styles.module.less'

const placeholder = text => {
  return ViewPlugin.fromClass(
    class {
      placeholder: any
      constructor(view) {
        this.placeholder = view.dom.ownerDocument.createElement('div')
        this.placeholder.textContent = text
        this.placeholder.className = 'cm-placeholder'
        if (view.state.doc?.length === 0) {
          view.dom.appendChild(this.placeholder)
        }
      }

      update(update) {
        if (update.docChanged || update.selectionSet) {
          if (update.state.doc?.length === 0) {
            update.view.dom.appendChild(this.placeholder)
          } else {
            if (this.placeholder.parentNode) {
              this.placeholder.parentNode.removeChild(this.placeholder)
            }
          }
        }
      }
    },
  )
}

const placeholderStyle = EditorView.baseTheme({
  '.cm-placeholder': {
    color: '#adb5bd',
    position: 'absolute',
    pointerEvents: 'none',
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    paddingLeft: '10px',
  },
})

const focusHandlerExtension = EditorView.domEventHandlers({
  click: (event, view) => {
    const target: any = event.target
    if (view.dom.contains(target) && view.state.doc?.length === 0) {
      startCompletion(view)
    }
  },
})

type SqlSearchIProps = {
  sqlEditorValue: string
  handleSearch: (val: string) => void
}

function getHistoryList() {
  return localStorage?.getItem('sqlEditorHistory')
    ? JSON.parse(localStorage?.getItem('sqlEditorHistory'))
    : []
}

function deleteHistoryByItem(val: string) {
  const lastHistory: any = localStorage.getItem('sqlEditorHistory')
  const tmp = lastHistory ? JSON.parse(lastHistory) : []
  if (tmp?.length > 0 && tmp?.includes(val)) {
    const newList = tmp?.filter(item => item !== val)
    localStorage.setItem('sqlEditorHistory', JSON.stringify(newList))
  }
}

const SqlSearch = memo(({ sqlEditorValue, handleSearch }: SqlSearchIProps) => {
  const editorRef = useRef<any>(null)
  const { t, i18n } = useTranslation()
  const clusterListRef = useRef<any>(null)
  const [clusterList, setClusterList] = useState([])
  const [historyCompletions, setHistoryCompletions] = useState<
    { value: string }[]
  >([])
  const historyCompletionsRef = useRef<any>(getHistoryList())

  function cacheHistory(val: string) {
    const lastHistory: any = localStorage.getItem('sqlEditorHistory')
    const tmp = lastHistory ? JSON.parse(lastHistory) : []
    const newList = [val, ...tmp?.filter(item => item !== val)]
    localStorage.setItem('sqlEditorHistory', JSON.stringify(newList))
    historyCompletionsRef.current = getHistoryList()
    setHistoryCompletions(historyCompletionsRef.current)
  }

  const { response } = useAxios({
    url: '/rest-api/v1/clusters',
    option: {
      params: {
        orderBy: 'name',
        ascending: true,
      },
    },
    manual: false,
    method: 'GET',
  })

  useEffect(() => {
    if (response?.success) {
      clusterListRef.current = response?.data?.items
      setClusterList(response?.data?.items)
    }
  }, [response])

  useEffect(() => {
    getHistoryList()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  function getCustomCompletions(regMatch, cusCompletions, pos) {
    const filterTerm = regMatch[2]
    const customCompletions = cusCompletions
      .filter(completion =>
        completion.toLowerCase().includes(filterTerm.toLowerCase()),
      )
      .map(completion => ({
        label: completion,
        type: 'custom',
        apply: completion,
        boost: 0,
      }))

    const from = pos - filterTerm?.length
    if (customCompletions?.length > 0) {
      return { from, options: customCompletions }
    }
    return null
  }

  useEffect(() => {
    if (editorRef.current) {
      const contentEditableElement = editorRef.current.querySelector(
        '.cm-content',
      ) as any
      if (contentEditableElement) {
        contentEditableElement.style.outline = 'none'
      }
      const customCompletionKeymap: KeyBinding[] = [
        { key: 'Tab', run: acceptCompletion },
      ]
      const overrideKeymap = keymap.of(
        customCompletionKeymap.concat(
          completionKeymap.filter(b => b.key !== 'Enter'),
        ),
      )
      const mySQLHighlightStyle = HighlightStyle.define([
        { tag: tags.keyword, color: 'blue' },
      ])

      const customCompletion = context => {
        const { state, pos } = context
        const beforeCursor = state.doc.sliceString(0, pos)
        if (state.doc?.length === 0) {
          const historyOptions: any[] = historyCompletionsRef?.current?.map(
            record => ({
              label: record,
              type: 'history',
              apply: record,
            }),
          )
          return {
            from: context.pos,
            options: historyOptions,
            filter: false,
          }
        }

        const whereMatch = /\b(where|or|and) (\S*)$/.exec(beforeCursor)
        if (whereMatch) {
          return getCustomCompletions(whereMatch, whereKeywords, pos)
        }

        const kindMatch = /(kind\s*=\s*)(\S*)$/i.exec(beforeCursor)
        if (kindMatch) {
          return getCustomCompletions(kindMatch, kindCompletions, pos)
        }

        if (whereKeywords?.some(item => beforeCursor?.endsWith(`${item} `))) {
          const customCompletions = operatorKeywords.map(completion => ({
            label: completion,
            type: 'custom',
            validFor: () => false,
          }))
          return { from: pos, options: customCompletions }
        }

        const clusterMatch = /(cluster\s*=\s*)(\S*)$/i.exec(beforeCursor)
        if (clusterMatch) {
          const clusterNameList = clusterListRef.current?.map(
            item => `'${item?.metadata?.name}'`,
          )
          return getCustomCompletions(clusterMatch, clusterNameList, pos)
        }

        const word = context?.matchBefore(/\w*/)

        if (!word || (word?.from === word?.to && !context?.explicit)) {
          return null
        }

        const options = defaultKeywords
          .filter(kw => kw.toLowerCase().includes(word.text?.toLowerCase()))
          .map(stmt => ({ label: stmt, type: 'custom' }))
        if (options?.length === 0) {
          return null
        }

        return completeFromList(options)(context)
      }

      const completionPlugin = ViewPlugin.fromClass(
        class {
          constructor(view) {
            this.addDeleteButtons(view)
          }

          update(update) {
            this.addDeleteButtons(update.view)
          }

          addDeleteButtons(view) {
            const compState: any = completionStatus(view.state)
            if (compState === 'active') {
              const completions: any = currentCompletions(view.state)
              setTimeout(() => {
                if (completions?.[0]?.type === 'history') {
                  view.dom
                    .querySelectorAll(
                      '.cm-tooltip.cm-tooltip-autocomplete > ul > li',
                    )
                    .forEach((item, index) => {
                      if (
                        item.querySelector(
                          '.cm-tooltip-autocomplete_item_label',
                        )
                      ) {
                        return
                      }
                      if (item.querySelector('.delete-btn')) {
                        return
                      }
                      const labelSpan = document.createElement('span')
                      labelSpan.className = 'cm-tooltip-autocomplete_item_label'
                      labelSpan.innerText = completions?.[index]?.label
                      item.style.display = 'flex'
                      item.style.justifyContent = 'space-between'
                      item.style.alignItems = 'center'
                      labelSpan.style.flex = '1'
                      labelSpan.style.overflow = 'hidden'
                      labelSpan.style.textOverflow = 'ellipsis'
                      labelSpan.style.whiteSpace = 'nowrap'
                      labelSpan.style.padding = '0 10px'
                      labelSpan.style.fontSize = '14px'
                      const btn = document.createElement('span')
                      btn.innerText = '✕'
                      btn.className = 'delete-btn'
                      btn.style.border = 'none'
                      btn.style.fontSize = '20px'
                      btn.style.display = 'flex'
                      btn.style.justifyContent = 'center'
                      btn.style.alignItems = 'center'
                      btn.style.height = '100%'
                      btn.style.padding = '0 15px'
                      item.innerText = ''
                      item.appendChild(labelSpan)
                      item.appendChild(btn)
                      btn.addEventListener('mousedown', event => {
                        event.preventDefault()
                        event.stopPropagation()
                        const completionOption = completions?.[index]
                        historyCompletionsRef.current =
                          historyCompletionsRef?.current?.filter(
                            item => item !== completionOption?.label,
                          )
                        if (view) {
                          startCompletion(view)
                          deleteHistoryByItem(completionOption?.label)
                        }
                      })
                    })
                }
              }, 0)
            }
          }
        },
      )

      const startState = EditorState.create({
        doc: '',
        extensions: [
          completionPlugin,
          placeholder(`${t('SearchUsingSQL')} ......`),
          placeholderStyle,
          new LanguageSupport(sql() as any),
          highlightSpecialChars(),
          syntaxHighlighting(mySQLHighlightStyle),
          autocompletion({
            override: [customCompletion],
          }),
          focusHandlerExtension,
          autocompletion(),
          overrideKeymap,
          EditorView.domEventHandlers({
            keydown: (event, view) => {
              if (event.key === 'Enter') {
                const completions = currentCompletions(view.state)
                if (!completions || completions?.length === 0) {
                  event.preventDefault()
                  handleClick()
                  return true
                }
              }
              return false
            },
          }),
          EditorState.allowMultipleSelections.of(false),
          EditorView.updateListener.of(update => {
            if (update.docChanged) {
              if (update.state.doc.lines > 1) {
                update.view.dispatch({
                  changes: {
                    from: update.startState.doc?.length,
                    to: update.state.doc?.length,
                  },
                })
              }
            }
          }),
        ],
      })

      const view = new EditorView({
        state: startState,
        parent: editorRef.current,
      })

      editorRef.current.view = view

      return () => {
        if (editorRef.current?.view) {
          // eslint-disable-next-line react-hooks/exhaustive-deps
          editorRef.current?.view?.destroy()
        }
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [editorRef.current, historyCompletions, i18n?.language])

  useEffect(() => {
    if (sqlEditorValue && clusterList && editorRef.current?.view) {
      editorRef.current?.view.dispatch({
        changes: {
          from: 0,
          to: editorRef.current?.view.state.doc?.length,
          insert: sqlEditorValue,
        },
      })
    }
  }, [clusterList, editorRef.current?.view, sqlEditorValue])

  const getContent = () => {
    if (editorRef.current?.view) {
      const content = editorRef.current.view.state.doc.toString()
      return content
    }
    return ''
  }

  function handleClick() {
    const inputValue = getContent()
    if (!inputValue) {
      message.warning(t('PleaseEnterValidSQLStatement'))
      return
    }
    cacheHistory(inputValue)
    handleSearch(inputValue)
  }

  return (
    <div className={styles.karpor_editor_container}>
      <div className={styles.karpor_editor}>
        <div className={styles.karpor_editor_prefix}>{searchSqlPrefix}</div>
        <div className={styles.karpor_editor_content}>
          <div
            style={{ flex: 1, width: 100 }}
            ref={editorRef}
            className={css`
              .cm-editor .cm-scroller {
                padding-left: 10px;
                box-sizing: border-box;
                background: #fff;
                border: none;
                overflow-x: auto;
              }
              .cm-editor .cm-scroller::-webkit-scrollbar {
                display: none;
              }
              .cm-editor .cm-content {
                height: 40px;
                line-height: 40px;
                font-size: 14px;
                padding: 0;
                overflow-x: auto;
                white-space: pre;
                background-color: #fff !important;
                border: none;
              }
              .cm-editor .cm-content::-webkit-scrollbar {
                display: none;
              }
              .cm-editor .cm-content:focus {
                outline: none !important;
              }
              .cm-line {
                height: 40px !important;
                padding: 0;
                overflow-x: auto;
                white-space: pre;
                border: none;
              }
              .CodeMirror-focused {
                outline: none !important;
              }
              .cm-editor.cm-focused {
                outline: none !important;
              }
              .cm-keyword {
                color: blue !important;
              }
              .cm-tooltip-autocomplete .cm-completion {
                background-color: #f5f5f5 !important;
              }
              .cm-tooltip-autocomplete .cm-completionLabel {
                padding: 10px !important;
                font-size: 18px !important;
              }
              .cm-tooltip.cm-tooltip-autocomplete {
                border-radius: 6px !important;
                border: none;
                padding: 10px !important;
                box-sizing: border-box;
              }
              .cm-tooltip.cm-tooltip-autocomplete > ul {
                box-sizing: border-box;
                height: auto;
                max-height: 40vh;
                overflow-y: auto !important;
              }
              .cm-tooltip.cm-tooltip-autocomplete > ul > li {
                background-color: #f5f5f5 !important;
                margin: 5px 0 !important;
                padding: 10px 0 !important;
                border-radius: 6px !important;
                width: auto !important;
                box-sizing: border-box;
              }

              .cm-tooltip.cm-tooltip-autocomplete
                > ul
                > li[aria-selected='true'],
              .cm-tooltip.cm-tooltip-autocomplete > ul > li:hover {
                background-color: #97a9f5 !important;
                color: white !important;
              }
            `}
          />
          <div className={styles.karpor_editor_divider}>
            <Divider type="vertical" />
          </div>
          <div className={styles.karpor_editor_btn_container}>
            <div onClick={handleClick} className={styles.karpor_editor_btn}>
              <img src={arrowRight} />
            </div>
          </div>
        </div>
      </div>
    </div>
  )
})

export default SqlSearch
