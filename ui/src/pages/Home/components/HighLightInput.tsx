import { useMemo, useState } from "react";
import { createEditor, Transforms, Text } from 'slate';
import { Slate, Editable, withReact } from 'slate-react';

import styles from './styles.less';

const REGEX = /{{(.*?)}}/g;
const REGEX_OP = /(AND|OR|NOT)/g;

const Leaf = ({ attributes, children, leaf }) => {
  const marginRightStyle = leaf.variable ? { marginRight: 5 } : {};
  const marginStyle = leaf.op ? { margin: 5 } : {}
  const style = {...marginRightStyle, ...marginStyle}
  return (
    <span className={`${leaf.variable ? styles.blueText : ''} ${leaf.op ? styles.pinkText : ''}`}
      style={style}
      {...attributes}>
      {children}
    </span>
  )
};


const decorate = ([node, path]) => {
  if (!Text.isText(node)) return [];

  const ranges = [];
  let match = null;

  while ((match = REGEX.exec(node.text)) !== null) {
    ranges.push({
      variable: true,
      anchor: { path, offset: match.index },
      focus: { path, offset: match.index + match[0].length },
    });
  }
  while ((match = REGEX_OP.exec(node.text)) !== null) {
    ranges.push({
      op: true,
      anchor: { path, offset: match.index },
      focus: { path, offset: match.index + match[0].length },
    });
  }

  return ranges;
};


const withSingleLine = (editor) => {
  const { normalizeNode } = editor;

  editor.normalizeNode = ([node, path]) => {
    if (path.length === 0) {
      if (editor.children.length > 1) {
        Transforms.mergeNodes(editor);
      }
    }

    return normalizeNode([node, path]);
  };

  return editor;
};

export default function App({value, onChange}) {
  const editor = useMemo(() => withSingleLine(withReact(createEditor())), []);
  // const [value, setValue] = useState([
  //   {
  //     type: "paragraph",
  //     children: [{ text: "Test input" }, { text: "111111 input" }],
  //   }
  // ]);
  return (
    <Slate
      editor={editor}
      value={value}
      onChange={onChange}
    >
      <Editable
        decorate={decorate}
        renderLeaf={Leaf}
        className={styles.input}
        style={{
          whiteSpace: "pre"
        }}
      />
    </Slate>
  );
}