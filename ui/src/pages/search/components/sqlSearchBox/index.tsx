import React, { useRef, useEffect } from 'react'
import { EditorState, Compartment } from '@codemirror/state'
import { EditorView, keymap } from '@codemirror/view'
import { basicSetup } from '@codemirror/basic-setup'
import { sql } from '@codemirror/lang-sql'

const SqlSearchBox = () => {
  const editorRef = useRef(null)
  const viewRef = useRef(null)

  useEffect(() => {
    if (editorRef.current) {
      // 初始化编辑器状态
      const startState = EditorState.create({
        doc: '',
        extensions: [
          basicSetup, // 基础设置（包括行号、括号匹配等）
          keymap.of([]), // 如果需要，你可以在这里添加快捷键
          sql(), // SQL 语法支持
          EditorView.lineWrapping, // 禁用换行
          new Compartment().of(EditorView.editable.of(false)), // 设置为不可编辑
        ],
      })

      // 初始化编辑器视图
      const view = new EditorView({
        state: startState,
        parent: editorRef.current,
      })

      viewRef.current = view

      // 返回清理函数以销毁编辑器视图
      return () => {
        view.destroy()
      }
    }
  }, [])

  return <div ref={editorRef}></div>
}

export default SqlSearchBox
