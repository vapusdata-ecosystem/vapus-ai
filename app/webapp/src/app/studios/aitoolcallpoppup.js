"use client";

import { useState, useEffect, useRef } from "react";
function useScript(src) {
  const [status, setStatus] = useState(src ? "loading" : "idle");

  useEffect(() => {
    if (!src) {
      setStatus("idle");
      return;
    }

    let script = document.querySelector(`script[src="${src}"]`);

    if (!script) {
      script = document.createElement("script");
      script.src = src;
      script.async = true;
      script.setAttribute("data-status", "loading");
      document.body.appendChild(script);

      const setAttributeFromEvent = (event) => {
        script.setAttribute(
          "data-status",
          event.type === "load" ? "ready" : "error"
        );
      };

      script.addEventListener("load", setAttributeFromEvent);
      script.addEventListener("error", setAttributeFromEvent);
    } else {
      // Get status from existing script
      setStatus(script.getAttribute("data-status"));
    }

    // Set state based on events
    const setStateFromEvent = (event) => {
      setStatus(event.type === "load" ? "ready" : "error");
    };

    script.addEventListener("load", setStateFromEvent);
    script.addEventListener("error", setStateFromEvent);

    return () => {
      if (script) {
        script.removeEventListener("load", setStateFromEvent);
        script.removeEventListener("error", setStateFromEvent);
      }
    };
  }, [src]);

  return status;
}

export default function AIToolCallPopup({
  isOpen,
  onClose,
  onAddTool,
  existingTools = null,
  editingIndex = -1,
}) {
  const [activeTab, setActiveTab] = useState("plainJSON");
  const [toolSchemaEditor, setToolSchemaEditor] = useState(null);
  const [functionName, setFunctionName] = useState("");
  const [functionDescription, setFunctionDescription] = useState("");
  const [plainJSONValue, setPlainJSONValue] = useState("");
  const [editorContent, setEditorContent] = useState(
    JSON.stringify({ schema: "Tool schema will appear here" }, null, 2)
  );
  const editorRef = useRef(null);
  const isEditing = editingIndex !== -1;

  // Load scripts dynamically
  const aceStatus = useScript(
    "https://cdnjs.cloudflare.com/ajax/libs/ace/1.4.12/ace.js"
  );
  const themeStatus = useScript(
    "https://cdnjs.cloudflare.com/ajax/libs/ace/1.4.12/theme-twilight.js"
  );
  const modeStatus = useScript(
    "https://cdnjs.cloudflare.com/ajax/libs/ace/1.4.12/mode-json.js"
  );

  const scriptsLoaded =
    aceStatus === "ready" && themeStatus === "ready" && modeStatus === "ready";

  useEffect(() => {
    if (!isOpen) {
      if (toolSchemaEditor) {
        toolSchemaEditor.destroy();
        setToolSchemaEditor(null);
      }
    }
  }, [isOpen, toolSchemaEditor]);

  // Initialize editor when modal is open and scripts are loaded
  useEffect(() => {
    if (
      typeof window !== "undefined" &&
      isOpen &&
      scriptsLoaded &&
      !toolSchemaEditor &&
      editorRef.current
    ) {
      try {
        const ace = window.ace;
        if (ace) {
          const editor = ace.edit(editorRef.current);
          editor.setTheme("ace/theme/twilight");
          editor.session.setMode("ace/mode/json");
          editor.setOptions({
            fontSize: "14px",
            showPrintMargin: false,
            highlightActiveLine: false,
          });
          editor.setValue(editorContent);

          // Save content when editor changes
          editor.on("change", () => {
            setEditorContent(editor.getValue());
          });

          setToolSchemaEditor(editor);
        }
      } catch (err) {
        console.error("Error initializing editor:", err);
      }
    }
  }, [isOpen, scriptsLoaded, toolSchemaEditor, editorContent]);

  const switchTab = (tab) => {
    setActiveTab(tab);
  };

  const generateJSONSchema = (data) => {
    if (data === null) {
      return { type: "null" };
    }

    if (Array.isArray(data)) {
      if (data.length === 0) {
        return { type: "array", items: {} };
      }
      return { type: "array", items: generateJSONSchema(data[0]) };
    }

    if (typeof data === "object") {
      const schema = {
        type: "object",
        properties: {},
        required: [],
      };

      for (const key in data) {
        if (Object.hasOwnProperty.call(data, key)) {
          schema.properties[key] = generateJSONSchema(data[key]);
          schema.required.push(key);
        }
      }

      return schema;
    }

    if (typeof data === "string") return { type: "string" };
    if (typeof data === "number") return { type: "number" };
    if (typeof data === "boolean") return { type: "boolean" };
    return {};
  };

  const convertToSchema = () => {
    try {
      const jsonData = JSON.parse(plainJSONValue);
      const schema = generateJSONSchema(jsonData);
      const schemaString = JSON.stringify(schema, null, 2);
      setEditorContent(schemaString);
      if (toolSchemaEditor) {
        toolSchemaEditor.setValue(schemaString);
        toolSchemaEditor.clearSelection();
      }

      switchTab("toolSchema");
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    if (existingTools && isOpen) {
      setFunctionName(existingTools.schema?.name || "");
      setFunctionDescription(existingTools.schema?.description || "");
      setPlainJSONValue(existingTools.rawJsonParams || "");

      const parameters =
        typeof existingTools.schema?.parameters === "string"
          ? existingTools.schema.parameters
          : JSON.stringify(
              existingTools.schema?.parameters || {
                schema: "Tool schema will appear here",
              },
              null,
              2
            );

      setEditorContent(parameters);

      // Update the editor if it exists
      if (toolSchemaEditor) {
        toolSchemaEditor.setValue(parameters);
        toolSchemaEditor.clearSelection();
      }
    } else if (isOpen) {
      setFunctionName("");
      setFunctionDescription("");
      setPlainJSONValue("");
      const defaultSchema = JSON.stringify(
        { schema: "Tool schema will appear here" },
        null,
        2
      );
      setEditorContent(defaultSchema);

      if (toolSchemaEditor) {
        toolSchemaEditor.setValue(defaultSchema);
        toolSchemaEditor.clearSelection();
      }
    }
  }, [existingTools, isOpen, toolSchemaEditor]);

  const saveChanges = () => {
    try {
      const toolData = {
        schema: {
          name: functionName,
          parameters: editorContent,
          description: functionDescription,
        },
        rawJsonParams: plainJSONValue,
        autoGenerate: false,
      };
      onAddTool(toolData);
      onClose();
    } catch (err) {
      console.error("Error saving tool schema:", err);
    }
  };

  if (!isOpen) return null;

  return (
    <>
      <div className="fixed inset-0 bg-zinc-600/80 flex items-center justify-center z-50">
        <div className="bg-zinc-800 rounded-lg shadow-lg w-3/4 h-3/4 md:w-1/2 overflow-y-auto scrollbar text-gray-100">
          <div className="bg-zinc-900 px-4 py-2 flex justify-between items-center">
            <h2 className="text-lg font-bold">
              {isEditing ? "Update AI Tool" : "Add AI Tool"}
            </h2>
            <button
              className="text-gray-100 cursor-pointer text-lg"
              type="button"
              onClick={onClose}
            >
              &times;
            </button>
          </div>

          <div className="px-4 py-2 flex justify-between items-center">
            <input
              type="text"
              placeholder="Enter function name..."
              className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
              value={functionName}
              onChange={(e) => setFunctionName(e.target.value)}
            />
          </div>

          <div className="px-4 py-2 flex justify-between items-center">
            <input
              type="text"
              placeholder="Enter description name..."
              className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
              value={functionDescription}
              onChange={(e) => setFunctionDescription(e.target.value)}
            />
          </div>

          <div>
            <div className="flex">
              <button
                type="button"
                className={`w-1/2 py-2 text-center ${
                  activeTab === "plainJSON"
                    ? "text-orange-700"
                    : "hover:text-orange-700 hover:border-orange-700 cursor-pointer"
                }`}
                onClick={() => switchTab("plainJSON")}
              >
                Plain JSON
              </button>
              <button
                type="button"
                className={`w-1/2 py-2 text-center ${
                  activeTab === "toolSchema"
                    ? "text-orange-700"
                    : "hover:text-orange-700 hover:border-orange-700 cursor-pointer"
                }`}
                onClick={() => switchTab("toolSchema")}
              >
                Tool Schema
              </button>
            </div>

            <div className="p-4">
              {/* Plain JSON Tab */}
              <div
                className={`tab-content ${
                  activeTab !== "plainJSON" ? "hidden" : ""
                }`}
              >
                <textarea
                  className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-500 p-2 shadow-sm w-full"
                  rows="9"
                  placeholder="Enter JSON here, and convert it to json Schema by clicking the button below..."
                  value={plainJSONValue}
                  onChange={(e) => setPlainJSONValue(e.target.value)}
                ></textarea>
                <button
                  onClick={convertToSchema}
                  type="button"
                  className="bg-zinc-900 text-gray-100 px-4 py-2 mt-3 text-sm rounded-lg hover:bg-zinc-900 cursor-pointer"
                >
                  Convert to JSON Schema
                </button>
              </div>

              {/* Tool Schema Tab */}
              <div
                className={`tab-content ${
                  activeTab !== "toolSchema" ? "hidden" : ""
                }`}
              >
                <div
                  ref={editorRef}
                  className="w-full scrollbar border border-zinc-600 rounded"
                  style={{ height: "300px" }}
                ></div>
              </div>
            </div>
          </div>

          <div className="px-4 py-2 flex justify-end">
            <button
              className="bg-zinc-900 text-gray-100 px-4 py-2 rounded hover:bg-gray-600 mr-2 cursor-pointer"
              type="button"
              onClick={onClose}
            >
              Cancel
            </button>
            <button
              onClick={saveChanges}
              type="button"
              className="bg-orange-700 text-gray-100 px-4 py-2 rounded hover:bg-zinc-900 cursor-pointer"
            >
              {isEditing ? "Update" : "Add"}
            </button>
          </div>
        </div>
      </div>
    </>
  );
}
