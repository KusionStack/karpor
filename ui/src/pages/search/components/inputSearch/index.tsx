import React, { useState } from 'react'
// import { Controlled as CodeMirror } from '@uiw/react-codemirror';
import 'codemirror/keymap/sublime'
import 'codemirror/theme/monokai.css'

const SearchBox = () => {
  const [value, setValue] = useState('')
  const [suggestions, setSuggestions] = useState([])

  const handleSelectSuggestion = suggestion => {
    setValue(value + suggestion)
    setSuggestions([])
  }

  const highlightKeyword = text => {
    const keywords = ['and', 'or']
    const highlightedText = text.replace(
      new RegExp(keywords.join('|'), 'g'),
      match => `<span style="background-color: yellow;">${match}</span>`,
    )
    return { __html: highlightedText }
  }

  return (
    <div>
      <div dangerouslySetInnerHTML={highlightKeyword(value)} />
      {suggestions.length > 0 && (
        <ul>
          {suggestions.map((suggestion, index) => (
            <li key={index} onClick={() => handleSelectSuggestion(suggestion)}>
              {suggestion}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

export default SearchBox
