import React, { useEffect, useRef } from "react";
import { EditorView } from "@codemirror/view";
import { EditorState } from "@codemirror/state";

function CMEditor() {
  const viewRef = useRef();
  const edContainer = useRef<any>();

  useEffect(() => {
    let onUpdateExt = EditorView.updateListener.of((v: ViewUpdate) => {
      if (v.docChanged) {
        viewRef.current!.dispatch({
          state: v.state,
        });
      }
    });
    const viewRef = new EditorView({
      parent: edContainer.current,
      doc: 'if(true) {return false}',
      // extensions: [basicSetup, javascript(), oneDark, onUpdateExt],
    });

    return () => {
      viewRef?.current?.destroy();
    };
  }, []);

  return <div ref={edContainer} className="container"></div>;
}

export default CMEditor;