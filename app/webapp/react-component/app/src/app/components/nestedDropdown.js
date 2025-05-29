"use client";
import { useState, useEffect, useRef } from "react";
import { modelsRegistryApi } from "../utils/ai-studio-endpoint/models-registry-api";

const NestedDropdown = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [selectedModel, setSelectedModel] = useState("");
  const [selectedModelDisplay, setSelectedModelDisplay] = useState("");
  const [aiModelNodes, setAIModelNodes] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [expandedNodeId, setExpandedNodeId] = useState(null);
  const dropdownRef = useRef(null);

  // Fetch AI model data from API

  useEffect(() => {
    const fetchAIModels = async () => {
      try {
        setIsLoading(true);
        const data = await modelsRegistryApi.getModelsRegistry();

        setAIModelNodes(data.output?.aiModelNodes || []);
        setIsLoading(false);
      } catch (err) {
        console.error("Error fetching AI models:", err);
        setError(err.message);
        setIsLoading(false);
      }
    };

    fetchAIModels();
  }, []);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  const handleToggleDropdown = () => {
    setIsOpen(!isOpen);
  };

  const handleToggleNode = (nodeId, event) => {
    event.stopPropagation();
    setExpandedNodeId(expandedNodeId === nodeId ? null : nodeId);
  };

  const handleSelectModel = (modelNodeId, modelName) => {
    setSelectedModel(`${modelNodeId}||${modelName}`);
    setSelectedModelDisplay(modelName);
    setIsOpen(false);
  };

  return (
    <div className="relative mb-2 mt-2" ref={dropdownRef}>
      <button
        onClick={handleToggleDropdown}
        className="vapus-dropdown-toggle border border-gray-400 cursor-pointer rounded-lg bg-zinc-800 text-sm text-white px-4 py-2 w-full text-left flex justify-between items-center"
      >
        {selectedModelDisplay || "Select AI Models"}
        <svg
          className={`w-4 h-4 transition-transform ${
            isOpen ? "rotate-180" : ""
          }`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M19 9l-7 7-7-7"
          />
        </svg>
      </button>

      {isOpen && (
        <div className="vapus-dropdown-menu absolute z-10 mt-1 w-full rounded-xl bg-[oklch(0.37_0.013_285.805)] text-[rgb(235,235,235)] shadow-lg border border-[oklch(0.552_0.016_285.938)] max-h-60 overflow-y-auto scrollbar">
          {isLoading ? (
            <div className="text-center py-4">Loading...</div>
          ) : error ? (
            <div className="text-center text-red-500 py-4">{error}</div>
          ) : aiModelNodes.length === 0 ? (
            <div className="text-center py-4">No models available</div>
          ) : (
            aiModelNodes.map((node) => (
              <div key={node.modelNodeId}>
                <div
                  className="vapus-dropdown-item parent bg-[#cecfcf] text-[#010101] cursor-pointer text-xs px-6 py-2 my-1 mx-2 rounded-xl font-bold flex justify-between items-center"
                  onClick={(e) => handleToggleNode(node.modelNodeId, e)}
                >
                  {node.name}
                  <svg
                    className={`w-4 h-4 transition-transform ${
                      expandedNodeId === node.modelNodeId ? "rotate-180" : ""
                    }`}
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M19 9l-7 7-7-7"
                    />
                  </svg>
                </div>

                {expandedNodeId === node.modelNodeId &&
                  node.attributes?.generativeModels?.map((model) => (
                    <div
                      key={`${node.modelNodeId}-${model.modelName}`}
                      className="vapus-dropdown-item child text-xs px-6 py-2 hover:bg-[oklch(0.21_0.006_285.885)] hover:text-[#FFFFFF] rounded-lg cursor-pointer ml-2"
                      onClick={() =>
                        handleSelectModel(node.modelNodeId, model.modelName)
                      }
                    >
                      {model.modelName}
                    </div>
                  ))}
              </div>
            ))
          )}
        </div>
      )}

      <input type="hidden" id="aiModel" value={selectedModel} />
    </div>
  );
};

export default NestedDropdown;
