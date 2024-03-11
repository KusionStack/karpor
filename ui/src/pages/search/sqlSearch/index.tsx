import React, { useEffect } from 'react'
import Editor, { useMonaco, OnMount } from '@monaco-editor/react'
import { languages as monacoLanguages } from 'monaco-editor/esm/vs/editor/editor.api'
import './style.css'

const innerSqlKeywords = [
  { label: 'SELECT', value: 'SELECT' },
  { label: 'FROM', value: 'FROM' },
  { label: 'WHERE', value: 'WHERE' },
  { label: 'INSERT INTO', value: 'INSERT INTO' },
  { label: 'UPDATE', value: 'UPDATE' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'CREATE TABLE', value: 'CREATE TABLE' },
  { label: 'ALTER TABLE', value: 'ALTER TABLE' },
  { label: 'DROP TABLE', value: 'DROP TABLE' },
  { label: 'INDEX', value: 'INDEX' },
  { label: 'VIEW', value: 'VIEW' },
  { label: 'JOIN', value: 'JOIN' },
  { label: 'GROUP BY', value: 'GROUP BY' },
  { label: 'ORDER BY', value: 'ORDER BY' },
  { label: 'HAVING', value: 'HAVING' },
  { label: 'UNION', value: 'UNION' },
  { label: 'INTERSECT', value: 'INTERSECT' },
  { label: 'EXCEPT', value: 'EXCEPT' },
  { label: 'DISTINCT', value: 'DISTINCT' },
  { label: 'IN', value: 'IN' },
  { label: 'LIKE', value: 'LIKE' },
  { label: 'IS NULL', value: 'IS NULL' },
  { label: 'AND', value: 'AND' },
  { label: 'OR', value: 'OR' },
  { label: 'NOT', value: 'NOT' },
  { label: 'CASE', value: 'CASE' },
  { label: 'WHEN', value: 'WHEN' },
  { label: 'THEN', value: 'THEN' },
  { label: 'ELSE', value: 'ELSE' },
  { label: 'END', value: 'END' },
  { label: 'BETWEEN', value: 'BETWEEN' },
  { label: 'EXISTS', value: 'EXISTS' },
  { label: 'CAST', value: 'CAST' },
  { label: 'CONVERT', value: 'CONVERT' },
  { label: 'COUNT', value: 'COUNT' },
  { label: 'MIN', value: 'MIN' },
  { label: 'MAX', value: 'MAX' },
  { label: 'AVG', value: 'AVG' },
  { label: 'SUM', value: 'SUM' },
  { label: 'AS', value: 'AS' },
  { label: 'INTO', value: 'INTO' },
  { label: 'VALUES', value: 'VALUES' },
  { label: 'PRIMARY KEY', value: 'PRIMARY KEY' },
  { label: 'FOREIGN KEY', value: 'FOREIGN KEY' },
  { label: 'REFERENCES', value: 'REFERENCES' },
  { label: 'CONSTRAINT', value: 'CONSTRAINT' },
  { label: 'TRIGGER', value: 'TRIGGER' },
  { label: 'PROCEDURE', value: 'PROCEDURE' },
  { label: 'FUNCTION', value: 'FUNCTION' },
]

const customSqlKeywords = [
  { label: 'CUSTOM1', insertText: 'CUSTOM1' },
  { label: 'CUSTOM2', insertText: 'CUSTOM2' },
  // ...其他自定义关键字
]

const SqlSearchBox: React.FC = () => {
  const monaco = useMonaco()

  useEffect(() => {
    if (monaco) {
      // 注册 SQL 语言的代码完成提供程序
      monaco.languages.registerCompletionItemProvider('sql', {
        provideCompletionItems: () => {
          const suggestions: monacoLanguages.CompletionItem[] = [
            ...innerSqlKeywords.map(keyword => ({
              label: keyword.label,
              kind: monacoLanguages.CompletionItemKind.Keyword,
              insertText: `${keyword.value} `,
              range: null as any, // 由 Monaco 动态计算范围
            })),
            // ...其他内置 SQL 关键字
            // 自定义关键字
            ...customSqlKeywords.map(keyword => ({
              label: keyword.label,
              kind: monacoLanguages.CompletionItemKind.Keyword,
              insertText: `${keyword.insertText} `,
              range: null as any, // 由 Monaco 动态计算范围
            })),
          ]
          return { suggestions }
        },
      })
    }
  }, [monaco])

  const handleEditorDidMount: OnMount = (editor, monaco) => {
    // 可以在这里进一步配置编辑器实例
    console.log(editor, monaco, '===editor, monaco===')
  }

  return (
    <Editor
      className="my-custom-editor"
      width="600px"
      height="40px" // 设置编辑器高度
      defaultLanguage="sql"
      defaultValue="-- 输入 SQL 代码"
      // theme="vs-dark"
      // loading=""
      onMount={handleEditorDidMount}
      options={{
        lineHeight: 40, // 设置行高
        fontSize: 20, // 设置字体大小
        lineNumbers: 'off', // 关闭行号
        glyphMargin: false, // 关闭 glyph 边缘
        folding: false, // 关闭折叠功能
        // 隐藏滚动条，但允许通过鼠标或触摸板滚动
        scrollbar: {
          vertical: 'hidden',
          handleMouseWheel: true,
        },
        scrollBeyondLastLine: false, // 禁止滚动到最后一行之后
        minimap: { enabled: false }, // 关闭小地图
        wordWrap: 'off', // 关闭自动换行
        wrappingIndent: 'none', // 没有换行缩进
        overviewRulerBorder: false, // 去除概览尺边框
        quickSuggestions: true, // 启用快速建议
        suggestOnTriggerCharacters: true, // 触发字符时提供建议
        wordBasedSuggestions: 'currentDocument', // 禁用基于单词的建议（使用自定义建议）
      }}
    />
  )
}

export default SqlSearchBox
