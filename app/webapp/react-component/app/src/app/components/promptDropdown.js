"use client";
import { useState, useEffect, useRef } from "react";
import { PromptsApi } from "../utils/ai-studio-endpoint/prompts-api";

const PromptDropdown = ({ onPromptSelect }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [selectedPrompt, setSelectedPrompt] = useState("");
  const [selectedPromptDisplay, setSelectedPromptDisplay] = useState("");
  const [prompts, setPrompts] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const dropdownRef = useRef(null);

  // Fetch prompts data
  useEffect(() => {
    const fetchPrompts = async () => {
      try {
        setIsLoading(true);
        const data = await PromptsApi.getPrompts();
        setPrompts(data.output || []);
        setIsLoading(false);
      } catch (err) {
        console.error("Error fetching prompts:", err);
        setError(err.message);
        setIsLoading(false);
      }
    };

    fetchPrompts();
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

  const handleSelectPrompt = (promptId, promptName) => {
    
    setSelectedPrompt(promptId);
    setSelectedPromptDisplay(promptName);
    setIsOpen(false);
    
    // Send only promptId to parent component
    if (onPromptSelect) {
      onPromptSelect(promptId);
    }
  };

  return (
    <div className="relative mb-2 mt-2" ref={dropdownRef}>
      <button
        onClick={handleToggleDropdown}
        className="vapus-dropdown-toggle border border-gray-400 cursor-pointer rounded-lg bg-zinc-800 text-sm text-white px-4 py-2 w-full text-left flex justify-between items-center"
      >
        {selectedPromptDisplay || "Select Prompt"}
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
          ) : prompts.length === 0 ? (
            <div className="text-center py-4">No prompts available</div>
          ) : (
            prompts.map((prompt) => (
              <div
                key={prompt.promptId}
                className="vapus-dropdown-item child text-xs px-6 py-2 hover:bg-[oklch(0.21_0.006_285.885)] hover:text-[#FFFFFF] rounded-lg cursor-pointer ml-2"
                onClick={() => handleSelectPrompt(prompt.promptId, prompt.name)}
              >
                {prompt.name}
              </div>
            ))
          )}
        </div>
      )}
    </div>
  );
};

export default PromptDropdown;