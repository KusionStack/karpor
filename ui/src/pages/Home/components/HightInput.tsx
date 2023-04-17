/*
 * 开发调试
 */
import React, { useEffect, useRef } from 'react';
import styles from './styles.less';

const keyword = ['repo', 'file', 'AND', 'OR', 'NOT'];
const replaceBlank = (html) => html.replace(/&nbsp;/ig, ' ');

// 获取光标位置
function getCaretCharacterOffsetWithin(element) {
  let caretOffset = 0;
  const doc = element.ownerDocument || element.document;
  const win = doc.defaultView || doc.parentWindow;
  let sel;
  if (typeof win.getSelection !== 'undefined') {
    sel = win.getSelection();
    if (sel.rangeCount > 0) {
      const range = win.getSelection().getRangeAt(0);
      const preCaretRange = range.cloneRange();
      preCaretRange.selectNodeContents(element);
      preCaretRange.setEnd(range.endContainer, range.endOffset);
      caretOffset = preCaretRange.toString().length;
    }
  } else if ((sel = doc.selection) && sel.type != 'Control') {
    const textRange = sel.createRange();
    const preCaretTextRange = doc.body.createTextRange();
    preCaretTextRange.moveToElementText(element);
    preCaretTextRange.setEndPoint('EndToEnd', textRange);
    caretOffset = preCaretTextRange.text.length;
  }
  return caretOffset;
}

// 设置光标位置
function setCaretPosition(element, offset) {
  const range = document.createRange();
  const sel = window.getSelection();

  // select appropriate node
  let currentNode = null;
  let previousNode = null;

  for (let i = 0; i < element.childNodes.length; i++) {
    // save previous node
    previousNode = currentNode;

    // get current node
    currentNode = element.childNodes[i];
    // if we get span or something else then we should get child node
    while (currentNode.childNodes.length > 0) {
      [currentNode] = currentNode.childNodes;
    }

    // calc offset in current node
    if (previousNode != null) {
      offset -= previousNode.length;
    }
    // check whether current node has enough length
    if (offset <= currentNode.length) {
      break;
    }
  }
  // move caret to specified offset
  if (currentNode != null) {
    range.setStart(currentNode, offset);
    range.collapse(true);
    sel.removeAllRanges();
    sel.addRange(range);
  }
}

// 排序方法
function compareWordLength(a, b) {
  if (a.length > b.length) {
    return -1;
  } if (a.length < b.length) {
    return 1;
  }
  return 0;
}

// 高亮关键字
function addKeyWordHighline(oText, keyWords) {
  let returnVal = oText;
  let i = 0;
  let wordReg;
  keyWords.sort(compareWordLength);

  for (i = 0; i < keyWords.length; i++) {
    if (keyWords[i] !== '') {
      wordReg = new RegExp(`(?!<span+>.[^<]*)${keyWords[i]}(?!.[^<]*<\/span>)`, 'g');
      returnVal = returnVal.replace(wordReg, `<span style="color:green;">${keyWords[i]}</span>`);
    }
  }
  return returnVal;
}

const CodeInput = ({
  value,
  disable,
  children,
  onChange,
}) => {
  const ref = useRef();
  // 是否锁定输入
  const isLock = useRef(false);

  const getRef = () => ref && ref.current;

  const onCompositionstart = (e) => {
    isLock.current = true;
  };
  const onCompositionend = (e) => {
    isLock.current = false;
  };

  // 解决中文输入的时候，直接输出英文字母的问题(中文输入期间，不允许输入字符)
  useEffect(() => {
    // 监听中文输入
    const el = getRef();
    el.addEventListener('compositionstart', onCompositionstart, false);
    el.addEventListener('compositionend', onCompositionend, false);
    return () => {
      el.removeEventListener('compositionstart', onCompositionstart, false);
      el.removeEventListener('compositionend', onCompositionend, false);
    };
  }, []);

  const onInput = () => {
    const el = getRef();
    // dom是否为空 || 是否为锁定模式
    if (!el || isLock.current) return;
    // 获取内容
    let text = el.innerHTML;
    // 是否修改了
    if (value !== text) {
      // 获取光标
      const position = getCaretCharacterOffsetWithin(el);
      // 替换空格
      text = replaceBlank(text);
      // 替换关键字
      text = addKeyWordHighline(text, keyword);
      el.innerHTML = text;
      // 更新父组件
      onChange(text);
      // 恢复位置
      setCaretPosition(el, position);
    }
  };

  return (
    <pre
      ref={ref}
      className={styles['formula-input']}
      contentEditable={!disable}
      dangerouslySetInnerHTML={{ __html: value }}
      onInput={onInput}
    >

      {/* {children} */}
    </pre>
  );
};


CodeInput.defaultProps = {
  value: '默认',   // value
  disable: false,   // 是否可用
  children: undefined, // 子元素
  onChange: () => { },
};

export default CodeInput;
