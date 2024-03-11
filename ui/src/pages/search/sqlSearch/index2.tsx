import React, { useEffect, useRef } from 'react'
import { EditorState } from '@codemirror/state'
import {
  EditorView,
  keymap,
  KeyBinding,
  highlightSpecialChars,
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
} from '@codemirror/autocomplete'
import { css } from '@emotion/css'
// import '@codemirror/view/dist'

import './style.css'

const sqlKeywords = [
  'SELECT',
  'FROM',
  'WHERE',
  'INSERT',
  'INTO',
  'VALUES',
  'UPDATE',
  'SET',
  'DELETE',
  'CREATE',
  'TABLE',
  'INDEX',
  'VIEW',
  'AS',
  'JOIN',
  'ON',
  'GROUP BY',
  'ORDER BY',
  'HAVING',
  'DISTINCT',
  'ALTER',
  'DROP',
  'PRIMARY KEY',
]

// 假设这是你的自定义语句数组
const customStatements = ['MY_CUSTOM_QUERY', 'ANOTHER_CUSTOM_STATEMENT']

// // 自定义补全函数
const customCompletion = context => {
  const word = context.matchBefore(/\w*/)

  console.log(word, 'word')

  // 当没有匹配到word，或者没有明确的补全请求时（例如，不是通过Ctrl+Space触发的），返回null
  if (!word || (word.from === word.to && !context.explicit)) {
    return null
  }

  const options = [...sqlKeywords, ...customStatements]
    .filter(kw => !kw.startsWith(word.text))
    .map(stmt => ({ label: stmt, type: 'custom' }))
  console.log(options, '===options===')
  // 如果没有合适的匹配项，也返回null让默认的补全介入
  if (options.length === 0) {
    return null
  }

  return completeFromList(options)(context)
}

const SqlEditor: React.FC = () => {
  const editorRef = useRef<HTMLDivElement>(null)

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
        // 禁用 Enter 键的默认选择行为
        {
          key: 'Enter',
          run: (view: EditorView) => {
            console.log(
              (view.state.field(autocompletion as any, false) as any)?.active,
              '===(view.state.field(autocompletion as any, false) as any)?.active==',
            )
            if (
              (view.state.field(autocompletion as any, false) as any)?.active
            ) {
              console.log('sdsadasdas')
              return acceptCompletion(view) // 如果自动完成激活，接受当前选择
            }
            // 如果自动完成未激活，阻止 Enter 键的默认行为
            return true // 返回 true 阻止默认行为
          },
          preventDefault: true,
        },
        // ...completionKeymap.filter(
        //   ({ key }) => key !== 'Enter' && key !== 'Tab',
        // ),
      ]
      console.log(completionKeymap, '==completionKeymap==')
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
      const startState = EditorState.create({
        doc: '',
        extensions: [
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

      return () => {
        view.destroy()
      }
    }
  }, [])

  return (
    <div
      // style={{ width: 300, height: 50, border: '1px solid red' }}
      ref={editorRef}
      className={css`
        .cm-editor .cm-content {
          border: 1px solid #d8d8d8;
          height: 40px;
          width: 800px;
          border-radius: 8px;
          line-height: 40px;
          font-size: 16px;
          padding: 0 10px;
          overflow-x: auto; // 允许水平滚动
          white-space: pre; // 禁止自动换行
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
      `}
    />
  )
}

export default SqlEditor
