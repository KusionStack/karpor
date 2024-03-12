import React, { useEffect, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
// import { useTranslation } from 'react-i18next'
import { EditorState } from '@codemirror/state'
import {
  EditorView,
  keymap,
  KeyBinding,
  highlightSpecialChars,
  ViewPlugin,
  // ViewUpdate,
} from '@codemirror/view'
import { sql } from '@codemirror/lang-sql'
import {
  // defaultHighlightStyle,
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
  // closeCompletion,
  currentCompletions,
} from '@codemirror/autocomplete'
import { css } from '@emotion/css'
import arrowRight from '@/assets/arrow-right.png'
// import '@codemirror/view/dist'

import './style.css'

const sqlKeywords = [
  'SELECT',
  'FROM',
  'WHERE',
  'VALUES',
  'AS',
  'JOIN',
  'ON',
  'GROUP BY',
  'ORDER BY',
  'HAVING',
  'DISTINCT',
  'PRIMARY KEY',
  'OR',
  'AND',
]

// 假设这是你的自定义语句数组
const customStatements = ['KIND', 'MY_CUSTOM_QUERY', 'ANOTHER_CUSTOM_STATEMENT']

// // 自定义补全函数
const customCompletion = context => {
  const { state, pos } = context
  const beforeCursor = state.doc.sliceString(0, pos)

  // 检查是否满足特定条件 "KIND="
  if (/(KIND=)|(kind=)|(kind = )$/.test(beforeCursor)) {
    // 提供特定的补全列表
    const customCompletions = ['cu1', 'cu2', 'cu3'].map(completion => ({
      label: completion,
      type: 'constant', // 根据需要调整类型
    }))
    return { from: pos, options: customCompletions }
  }
  const word = context.matchBefore(/\w*/)

  console.log(word, 'word')

  // 当没有匹配到word，或者没有明确的补全请求时（例如，不是通过Ctrl+Space触发的），返回null
  if (!word || (word.from === word.to && !context.explicit)) {
    return null
  }

  const options = [...sqlKeywords, ...customStatements]
    .filter(kw => kw.toLowerCase().startsWith(word.text.toLowerCase()))
    .map(stmt => ({ label: stmt, type: 'custom' }))
  console.log(options, '===options===')
  // 如果没有合适的匹配项，也返回null让默认的补全介入
  if (options.length === 0) {
    return null
  }

  return completeFromList(options)(context)
}

// 创建一个插件来显示占位符
const placeholder = text => {
  return ViewPlugin.fromClass(
    class {
      placeholder: any
      constructor(view) {
        this.placeholder = view.dom.ownerDocument.createElement('div')
        this.placeholder.textContent = text
        this.placeholder.className = 'cm-placeholder'
        // 只在编辑器内容为空时显示占位符
        if (view.state.doc.length === 0) {
          view.dom.appendChild(this.placeholder)
        }
      }

      update(update) {
        if (update.docChanged || update.selectionSet) {
          if (update.state.doc.length === 0) {
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

// 自定义占位符样式
const placeholderStyle = EditorView.baseTheme({
  '.cm-placeholder': {
    color: '#adb5bd',
    position: 'absolute',
    pointerEvents: 'none',
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    padding: '0 20px',
  },
})

const SqlEditor: React.FC = () => {
  const editorRef = useRef<any>(null)
  const navigate = useNavigate()
  // const { t } = useTranslation()

  useEffect(() => {
    if (editorRef.current) {
      // 添加CSS样式去除虚线框
      // editorRef.current.style.outline = 'none'
      // 查找编辑区并去除虚线框
      const contentEditableElement = editorRef.current.querySelector(
        '.cm-content',
      ) as any
      if (contentEditableElement) {
        contentEditableElement.style.outline = 'none'
      }

      // 自定义键绑定来使用 Tab 键选择自动完成项
      const customCompletionKeymap: KeyBinding[] = [
        { key: 'Tab', run: acceptCompletion },
        // { key: 'Tab', run: completionKeymap.Enter.run },
        // {
        //   key: 'Enter',
        //   run: view => {
        //     // 如果当前有活跃的补全，关闭它
        //     if ((view.state.field(autocompletion as any) as any).open) {
        //       closeCompletion(view)
        //       return true // 阻止进一步处理 Enter 键
        //     } else {
        //       // 如果没有活跃的补全，执行你的操作
        //       const content = view.state.doc.toString()
        //       console.log(content) // 或者其它操作
        //       return false // 不阻止原本的 Enter 键操作
        //     }
        //   },
        //   preventDefault: true,
        //   stopPropagation: true,
        // },
        // 禁用 Enter 键的默认选择行为
        // {
        //   key: 'Enter',
        //   run: (view: EditorView) => {
        //     console.log(
        //       view,
        //       (view.state.field(autocompletion as any, false) as any)?.active,
        //       '===(view.state.field(autocompletion as any, false) as any)?.active==',
        //     )
        //     if (
        //       (view.state.field(autocompletion as any, false) as any)?.active
        //     ) {
        //       console.log('是否正确')
        //       return acceptCompletion(view) // 如果自动完成激活，接受当前选择
        //     }
        //     // 如果自动完成未激活，阻止 Enter 键的默认行为
        //     return true // 返回 true 阻止默认行为
        //   },
        //   preventDefault: true,
        // },
        // ...completionKeymap.filter(
        //   ({ key }) => key !== 'Enter' && key !== 'Tab',
        // ),
      ]
      // console.log(completionKeymap, '==completionKeymap==')
      // 创建一个新的键绑定数组，将自定义键绑定放在默认绑定前面

      const overrideKeymap = keymap.of(
        customCompletionKeymap.concat(
          completionKeymap.filter(b => b.key !== 'Enter'),
        ),
      )
      console.log(overrideKeymap, '===overrideKeymap===')
      // 创建一个新的高亮样式
      const mySQLHighlightStyle = HighlightStyle.define([
        { tag: tags.keyword, color: 'blue' },
      ])

      // 创建一个新的键绑定来替换默认的 completionKeymap
      // const myKeymap = keymap.of([
      //   // 通过返回 true 来阻止 Enter 键的默认补全行为
      //   {
      //     key: 'Enter',
      //     run: () => true,
      //     shift: acceptCompletion,
      //   },
      //   // 保留其他默认的补全键绑定，而不是 Enter 键
      //   ...Object.keys(completionKeymap)
      //     .filter(key => key !== 'Enter' && key !== 'Shift-Enter')
      //     .map(key => ({ key, run: completionKeymap[key].run })),
      //   // 添加一个绑定，使 Tab 键用于补全
      //   { key: 'Tab', run: acceptCompletion },
      // ])

      const startState = EditorState.create({
        doc: '',
        extensions: [
          placeholder('Search using SQL ......'),
          placeholderStyle,
          new LanguageSupport(sql() as any),
          highlightSpecialChars(),
          // syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
          syntaxHighlighting(mySQLHighlightStyle),
          // EditorView.lineWrapping,
          autocompletion({
            override: [customCompletion],
          }),
          autocompletion(),
          // keymap.of([...completionKeymap]),
          // keymap.of(customCompletionKeymap), // 使用自定义键绑定
          overrideKeymap,
          EditorView.domEventHandlers({
            keydown: (event, view) => {
              if (event.key === 'Enter') {
                // 获取当前的自动补全状态
                const completions = currentCompletions(view.state)
                // 检查是否有活跃的补全提示
                if (!completions || completions.length === 0) {
                  // 补全提示不活跃时获取当前文档内容
                  const content = view.state.doc.toString()
                  console.log(content) // 执行你希望的操作
                  // 防止默认的 Enter 键行为
                  event.preventDefault()
                  return true
                }
                // 如果有补全提示，保持默认行为
              }
              return false
            },
          }),
          // myKeymap,
          EditorState.allowMultipleSelections.of(false),
          EditorView.updateListener.of(update => {
            if (update.docChanged) {
              // 检查是否有换行发生
              if (update.state.doc.lines > 1) {
                // 如果有换行，撤销改变
                update.view.dispatch({
                  changes: {
                    from: update.startState.doc.length,
                    to: update.state.doc.length,
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
  }, [])

  // 用来获取编辑器内容的函数
  const getContent = () => {
    if (editorRef.current?.view) {
      const content = editorRef.current.view.state.doc.toString()
      console.log(content)
      return content
    }
    return ''
  }

  function handleClick() {
    const inputValue = getContent()
    console.log(inputValue)
    navigate(`/search/result?query=${inputValue}&pattern=sql`)
  }

  return (
    <div
      style={{
        position: 'relative',
        background: '#f1f1f1',
        border: '10px solid #f1f1f1',
        borderRadius: 16,
      }}
    >
      <div
        style={{ width: '100%' }}
        // style={{ width: 300, height: 50, border: '1px solid red' }}
        ref={editorRef}
        className={css`
          .cm-editor .cm-content {
            border: 1px solid #d8d8d8;
            height: 42px;
            width: 1000px;
            border-radius: 16px;
            line-height: 42px;
            font-size: 16px;
            padding: 0 10px;
            overflow-x: auto; // 允许水平滚动
            white-space: pre; // 禁止自动换行
            background-color: #fff !important;
          }
          .cm-editor .cm-content:focus {
            outline: none !important;
          }
          .cm-line {
            width: 100%;
            overflow-x: auto; // 允许水平滚动
            white-space: pre; // 禁止自动换行
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
            // font-weight: bold !important;
            font-size: 18px !important;
          }
          .cm-tooltip-autocomplete .cm-completion[aria-selected='true'],
          .cm-tooltip-autocomplete .cm-completion:hover {
            background-color: #f5f5f5 !important;
            color: white !important;
          }
          .cm-tooltip-autocomplete .cm-completionDetail {
            font-size: 0.8em !important;
            margin-left: 10px !important;
            color: #666 !important;
          }

          .cm-tooltip-autocomplete .cm-completionLabel .cm-completionText {
            padding: 20px !important;
          }
          .cm-tooltip-autocomplete .cm-completionLabel .cm-completionDetail {
            padding: 20px !important;
          }
          .cm-tooltip.cm-tooltip-autocomplete {
            border-radius: 10px !important;
            border: 1px solid #fafafa !important;
          }
          .cm-tooltip.cm-tooltip-autocomplete > ul > li {
            background-color: #f3f3f3 !important;
            padding: 10px 0 !important;
            border-radius: 10px !important;
          }

          .cm-tooltip.cm-tooltip-autocomplete > ul > li[aria-selected='true'],
          .cm-tooltip.cm-tooltip-autocomplete > ul > li:hover {
            background-color: #2f54eb !important;
            color: white !important;
          }
        `}
      />
      <div
        onClick={handleClick}
        style={{
          width: 40,
          height: 28,
          textAlign: 'center',
          position: 'absolute',
          right: 10,
          top: 6,
          background: '#2f54eb',
          borderRadius: 16,
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          cursor: 'pointer',
        }}
      >
        <img src={arrowRight} style={{ width: 16, height: 16 }} />
      </div>
    </div>
  )
}

export default SqlEditor
