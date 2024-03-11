import React, { useEffect, useRef, useState } from 'react'
import { EditorState } from '@codemirror/state'
import { EditorView, keymap } from '@codemirror/view'
import { basicSetup } from '@codemirror/basic-setup'
import { autocompletion, completionKeymap } from '@codemirror/autocomplete'

const SearchBoxWithAutocomplete = () => {
  const editorContainerRef = useRef(null)
  const [editorView, setEditorView] = useState(null)

  // 自定义补全函数
  const myCompletions = context => {
    const word = context.matchBefore(/\w*/)
    if (word === null) return null

    // 定义候选词列表
    const completions = [
      { label: 'apple', type: 'keyword' },
      { label: 'banana', type: 'keyword' },
      { label: 'cherry', type: 'keyword' },
      { label: 'date', type: 'keyword' },
      { label: 'elderberry', type: 'keyword' },
    ].filter(item => item.label.startsWith(word.text))

    return {
      from: word.from,
      options: completions,
    }
  }

  useEffect(() => {
    if (editorContainerRef.current && !editorView) {
      // 初始化 CodeMirror 实例
      const startState = EditorState.create({
        doc: '',
        extensions: [
          basicSetup,
          keymap.of(completionKeymap),
          autocompletion({ override: [myCompletions] }),
        ],
      })

      const view = new EditorView({
        state: startState,
        parent: editorContainerRef.current,
      })

      setEditorView(view)
    }

    // 组件卸载时销毁 CodeMirror 实例
    return () => {
      editorView?.destroy()
    }
  }, [editorView])

  // 渲染搜索框的容器
  return <div ref={editorContainerRef} />
}

export default SearchBoxWithAutocomplete
