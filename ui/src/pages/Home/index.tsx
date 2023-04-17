
import { styled, useModel } from '@alipay/bigfish';
import React, { useState, useRef, useEffect } from '@alipay/bigfish/react';
import { PageContainer } from '@alipay/tech-ui';
import { Button, Input, Select, Space } from 'antd';
import { FontColorsOutlined, ExpandOutlined, HeatMapOutlined } from "@ant-design/icons"
import classnames from 'classnames'
import Highlighter from "react-highlight-words";
import Tabs from "@/components/Tabs";
import { basicSyntaxColumns } from "@/utils/constants";
import { history } from '@alipay/bigfish';
// import CMEditor from "./components/CodeMirrirEditor";

import styles from "./style.less";

const { Search } = Input;


const options = [
  {
    value: 'filter',
    label: 'Filter(s)',
  },
  {
    value: 'AI',
    label: 'AI Suggestion',
  },
  {
    value: 'SQL',
    label: 'SQL',
  },
];

const Wrapper = styled.div`
  padding-top: 100px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
`;

export const HighlightedLabel: React.FunctionComponent<{ label: string; matches: Set<number>; offset?: number }> = ({
  label,
  matches,
  offset = 0,
}) => {
  const spans: [number, number, boolean][] = []
  let currentStart = 0
  let currentEnd = 0
  let currentMatch = false

  for (let index = 0; index <= label.length; index++) {
    currentEnd = index

    const match = matches.has(index + offset)
    if (currentMatch !== match || index === label.length) {
      // close previous span
      spans.push([currentStart, currentEnd, currentMatch])
      currentStart = index
      currentMatch = match
    }
  }

  return (
    <span>
      {spans.map(([start, end, match]) => {
        const value = label.slice(start, end)
        return match ? (
          <span key={offset + start} className={styles.match}>
            {value}
          </span>
        ) : (
          value
        )
      })}
    </span>
  )
}

const FilterOption: React.FunctionComponent<{ option: Option }> = ({ option }) => (
  <span className={classnames(styles.filterOption, 'search-filter-keyword')}>
    {option.matches ? <HighlightedLabel label={option.label} matches={option.matches} /> : option.label}
    <span className={styles.separator}>:</span>
  </span>
)

const FilterValueOption: React.FunctionComponent<{ option: Option }> = ({ option }) => {
  const label = option.label
  const separatorIndex = label.indexOf(':')
  const field = label.slice(0, separatorIndex)
  const value = label.slice(separatorIndex + 1)

  return (
    <span className={styles.filterOption}>
      <span className={styles.searchFilterKeyword}>
        {field}
        <span className='search-filter-separator'>:</span>
      </span>
      {option.matches ? <HighlightedLabel label={value} matches={option.matches} /> : value}
    </span>
  )
}

export default () => {


  const tabsList = [
    { label: 'Code search basics', value: 'code' },
    { label: 'Search query examples', value: 'examples' },
  ]
  const [currentTab, setCurrentTab] = useState<String>(tabsList?.[0].value)
  const [searchType, setSearchType] = useState<String>('filter');

  const [inputValue, setInputValue] = useState('');

  const handleTabChange = (value: String) => {
    setCurrentTab(value);
  }

  // const divInputRef = useRef();

  const addItem = (queryExample: string) => {
    console.log(queryExample, "[===queryExample==")
    const val = inputValue ? `${inputValue} ${queryExample}` : queryExample;
    // console.log(val, val.split(' '), "===val===")
    // const valList = val.split(' ');
    // let str = "";
    // let index = 0;
    //   while (index < valList.length) {
    //     // str += `<span>${valList[index]}</span><span style="width:5px;"> </span>`;
    //     const label = valList[index];
    //     const separatorIndex = label.indexOf(':')
    //     const field = label.slice(0, separatorIndex)
    //     const value = label.slice(separatorIndex + 1)
    //     str += `<span className="filterOption">
    //     <span className="searchFilterKeyword">
    //         ${field}
    //         <span className='search-filter-separator'>:</span>
    //     </span>
    //     ${value}
    // </span>`
    //     index++;
    //   }
    // console.log(str, "===str+++++")
    // divInputRef.current.innerHTML = str;
    // divInputRef.current.innerText = val;
    // setInputValue([...inputValue, queryExample]);
    setInputValue(val);
  }


  const renderItem = (query: string) => {
    return renderHight(query);
    return <Highlighter
      highlightClassName="YourHighlightClass"
      searchWords={["repo:", "file:", "lang:", "AND", "OR", "NOT"]}
      autoEscape={true}
      caseSensitive={true}
      highlightStyle={{ color: '#1890ff', background: '#fff' }}
      textToHighlight={query}
    />
  }

  const renderFilter = () => {
    return <>
      <Tabs list={tabsList} current={currentTab} onChange={handleTabChange} />
      <div className={styles.box}>
        {
          basicSyntaxColumns?.map((item: any, index) => {
            return <div className={styles.basicSyntax} key={index}>
              {
                item?.map((queryExample: any) => {
                  return <div className={styles.itemBox} key={queryExample.title}>
                    <div className={styles.title}>{queryExample.title}</div>
                    {
                      queryExample.queryExamples.map((example: any) => {
                        return <div className={styles.childrenItemBox} key={example.id} onClick={() => addItem(example.query)}>
                          <span className={styles.childrenItem}>{renderItem(example.query)}</span>
                          {/* <FilterValueOption option={{ label: example.query, matches: false }} /> */}
                        </div>
                      })
                    }
                  </div>
                })
              }

            </div>
          })
        }
      </div>
    </>
  }

  // 获取光标位置
  // function getCaretCharacterOffsetWithin(element) {
  //   let caretOffset = 0;
  //   const doc = element.ownerDocument || element.document;
  //   const win = doc.defaultView || doc.parentWindow;
  //   let sel;
  //   if (typeof win.getSelection !== 'undefined') {
  //     sel = win.getSelection();
  //     if (sel.rangeCount > 0) {
  //       const range = win.getSelection().getRangeAt(0);
  //       const preCaretRange = range.cloneRange();
  //       preCaretRange.selectNodeContents(element);
  //       preCaretRange.setEnd(range.endContainer, range.endOffset);
  //       caretOffset = preCaretRange.toString().length;
  //     }
  //   } else if ((sel = doc.selection) && sel.type != 'Control') {
  //     const textRange = sel.createRange();
  //     const preCaretTextRange = doc.body.createTextRange();
  //     preCaretTextRange.moveToElementText(element);
  //     preCaretTextRange.setEndPoint('EndToEnd', textRange);
  //     caretOffset = preCaretTextRange.text.length;
  //   }
  //   return caretOffset;
  // }


  // // 设置光标位置
  // function setCaretPosition(element, offset) {
  //   const range = document.createRange();
  //   const sel = window.getSelection();

  //   // select appropriate node
  //   let currentNode = null;
  //   let previousNode = null;

  //   for (let i = 0; i < element.childNodes.length; i++) {
  //     // save previous node
  //     previousNode = currentNode;

  //     // get current node
  //     currentNode = element.childNodes[i];
  //     // if we get span or something else then we should get child node
  //     while (currentNode.childNodes.length > 0) {
  //       [currentNode] = currentNode.childNodes;
  //     }

  //     // calc offset in current node
  //     if (previousNode != null) {
  //       offset -= previousNode.length;
  //     }
  //     // check whether current node has enough length
  //     if (offset <= currentNode.length) {
  //       break;
  //     }
  //   }
  //   // move caret to specified offset
  //   if (currentNode != null) {
  //     console.log(currentNode.length, offset, "===saasdqweq==")
  //     range.setStart(currentNode, offset);
  //     range.collapse(true);
  //     sel.removeAllRanges();
  //     sel.addRange(range);
  //   }
  // }


  // const handleKeyDown = (e) => {
  //   const val = e.target.innerText;
  //   setInputValue(val);
  // }

  function renderHight(str: string) {
    if (!str) return;
    const strList = str.split('\n').join(' ').split(' ');
    console.log(strList, "===strList===")
    let list = strList?.map((item) => {
      let tmp;
      let itemTmp = item.trim();
      if (itemTmp.indexOf(':') > -1) {
        const tmp1 = itemTmp.split(':');
        tmp = <span className={styles.txtSpan}><span style={{ color: '#1890ff' }}>{tmp1[0]}</span>:{tmp1[1]}</span>;
      } else if (itemTmp === 'OR' || itemTmp === 'AND' || itemTmp === 'NOT') {
        tmp = <span className={styles.txtSpan} style={{ color: '#a305e1', margin: '0 5px' }}>{itemTmp}</span>
      } else {
        tmp = <span className={styles.txtSpan}>{itemTmp}</span>
      }
      return tmp;
    })
    return list;
  }

  // function insertAfter(newEl, targetEl) {
  //   var parentEl = targetEl.parentNode;
  //   if (parentEl.lastChild == targetEl) {
  //     parentEl.appendChild(newEl);
  //   } else {
  //     parentEl.insertBefore(newEl, targetEl.nextSibling);
  //   }
  // }

  // function createCursor() {
  //   var fragment = document.createDocumentFragment();
  //   var wrapSpan = document.createElement('span');
  //   wrapSpan.className = "cursor";
  //   var lineSpan = document.createElement('span');
  //   lineSpan.className = "line";
  //   wrapSpan.appendChild(lineSpan);
  //   fragment.appendChild(wrapSpan);
  //   return fragment;
  // }

  // useEffect(() => {
  //   let txts = document.querySelectorAll(".txtSpan");
  //   let cursor = createCursor();
  //   Array.from(txts).forEach(item => {
  //     item.addEventListener("click", function () {
  //       cursor = document.querySelector(".cursor");
  //       cursor && this.parentNode.removeChild(cursor);
  //       let previousSibling = item.previousSibling
  //       if (previousSibling) {
  //         insertAfter(cursor, previousSibling);
  //       } else {
  //         insertAfter(cursor, item);
  //       }
  //     })
  //   })
  // })

  // const renderAI = () => {
  //   const spanList = renderHight(inputValue);

  //   return <div className={styles.inputWrapper}>
  //     <Select defaultValue="filter" className={styles.sele} options={options} style={{ width: 140 }} onChange={(val) => setSearchType(val)} value={searchType} />
  //     <div className={styles.divInput} contentEditable={true} placeholder='input search text'
  //       ref={divInputRef} onInput={handleKeyDown}
  //     // dangerouslySetInnerHTML={{
  //     //   __html: c.join(''),
  //     // }}
  //     >
  //       {
  //         spanList && spanList?.length > 0 && spanList?.map(item => <>{item}</>)
  //       }
  //       {/* <span className={styles.cursor}>
  //         <span className={styles.line}></span>
  //       </span> */}
  //       {/* {
  //         inputValue && inputValue?.split('@@@')?.map(item => {
  //           return <Highlighter
  //             highlightClassName="YourHighlightClass"
  //             searchWords={["repo:", "file:", "lang:", "AND", "OR", "NOT"]}
  //             autoEscape={true}
  //             caseSensitive={true}
  //             highlightStyle={{ color: '#1890ff', background: '#fff' }}
  //             style={{ marginRight: 5 }}
  //             textToHighlight={item}
  //           />
  //         })
  //       } */}
  //       {/* <CodeMirror
  //         value={content}
  //         options={options1}
  //         onChange={onChange}
  //       /> */}
  //     </div>
  //     <FontColorsOutlined className={styles.iconFont} />
  //     <ExpandOutlined className={styles.iconFont} />
  //   </div>
  // }

  const handleSearch = (val: string) => {
    history.push(`/result?keyword=${val}`)
  }


  const handleInputChange = (event: any) => {
    setInputValue(event.target.value);
  }

  return (
    <Wrapper>
      {/* <CMEditor/> */}
      <div className={styles.logoBox}>
        <HeatMapOutlined />
        <div className={styles.title}>Karbour</div>
      </div>

      <Space.Compact>
        <Select defaultValue="filter" options={options} style={{ width: 150 }} onChange={(val) => setSearchType(val)} value={searchType} />
        <Search placeholder="input search text" allowClear style={{ width: 600 }} value={inputValue} onChange={handleInputChange} onSearch={handleSearch} />
      </Space.Compact>
      {/* {
        renderAI()
      } */}
      <div className={styles.content}>
        {
          searchType === 'filter' && renderFilter()
        }
      </div>
    </Wrapper>
  );
};
