"use client";

import { useEffect, useState } from "react";
import React from "react";

// We'll import CodeMirror only on the client side
let CodeMirror;
let editor;

export default function YamlEditorClient() {
  const [yamlContent, setYamlContent] = useState("");
  const [modalParams, setModalParams] = useState({});
  const [apiMethod, setApiMethod] = useState("");

  useEffect(() => {
    // Import all CodeMirror dependencies
    import("codemirror/lib/codemirror.css");
    import("codemirror/theme/dracula.css");
    import("codemirror/addon/lint/lint.css");

    // Import CodeMirror and its modules
    Promise.all([
      import("codemirror"),
      import("codemirror/mode/yaml/yaml"),
      import("codemirror/addon/edit/matchbrackets"),
      import("codemirror/addon/selection/active-line"),
      import("codemirror/addon/hint/show-hint"),
      import("codemirror/addon/display/autorefresh"),
    ]).then(([codeMirrorModule]) => {
      CodeMirror = codeMirrorModule.default;

      if (!editor && document.getElementById("yamlEditor")) {
        initCodeMirror();
      }
    });

    // Cleanup function
    return () => {
      if (editor) {
        // Clean up editor if needed
      }
    };
  }, []);

  function initCodeMirror() {
    if (!document.getElementById("yamlEditor")) return;

    editor = CodeMirror(document.getElementById("yamlEditor"), {
      mode: "yaml",
      lineNumbers: true,
      lineWrapping: true,
      theme: "dracula",
      viewportMargin: Infinity,
      lint: true,
      gutters: ["CodeMirror-linenumbers", "CodeMirror-lint-markers"],
      indentUnit: 2,
      tabSize: 2,
      autoCloseBrackets: true,
      matchBrackets: true,
      styleActiveLine: true,
      highlightSelectionMatches: { showToken: true, annotateScrollbar: true },
    });
  }

  function loadFileContent(event) {
    const file = event.target.files[0];
    if (file && editor) {
      const reader = new FileReader();
      reader.onload = function (e) {
        editor.setValue(e.target.result);
      };
      reader.readAsText(file);
    }
  }

  function openYAMLedModal(apiUrl, tokenKey, contentDiv, method) {
    console.log(apiUrl);
    document.getElementById("yamlModal")?.classList.remove("hidden");
    setModalParams({ apiUrl, tokenKey });
    if (method) {
      setApiMethod(method);
    }

    if (!editor) {
      initCodeMirror();
    }

    if (contentDiv) {
      const val = document.getElementById(contentDiv)?.innerText;
      if (val && editor) {
        editor.setValue(val);
      }
    }
  }

  function closeYAMLedModal() {
    document.getElementById("yamlModal")?.classList.add("hidden");
    if (editor) {
      editor.setValue("");
    }
  }

  return (
    <div className="bg-zinc-800 p-4 rounded-lg shadow-lg w-full overflow-y-auto scrollbar relative">
      <div className="flex flex-col items-center space-y-2">
        <label
          htmlFor="fileInput"
          className="cursor-pointer inline-block text-sm px-2 py-2 bg-orange-700 text-white shadow-md font-semibold rounded-lg hover:bg-pink-900 transition"
        >
          Upload YAML Spec File
        </label>
        <input
          type="file"
          id="fileInput"
          accept=".yaml, .yml"
          onChange={loadFileContent}
          className="mb-2 w-full px-2 py-2 border rounded-lg focus:outline-none hidden"
        />
      </div>

      <div hidden id="modalParams"></div>

      <label
        htmlFor="yamlEditor"
        className="block text-gray-100 font-medium mb-2"
        id="yamlSpecTitle"
      >
        Edit YAML Content Here
      </label>

      <div
        id="yamlEditor"
        className="w-full h-auto border-2 rounded-lg focus:outline-none font-mono"
      ></div>

      <div className="mt-4 flex justify-end space-x-2">
        <button
          id="submitYaml"
          type="button"
          className="px-4 py-2 bg-orange-700 text-white rounded-lg hover:bg-pink-900 shadow-lg cursor-pointer"
        >
          Create
        </button>
      </div>
    </div>
  );
}
