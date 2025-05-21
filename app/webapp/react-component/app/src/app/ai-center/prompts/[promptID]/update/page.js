"use client";
import { useState, useEffect, useRef } from "react";
import { useRouter } from "next/navigation";
import { use } from "react";
import Sidebar from "@/app/components/platform/main-sidebar";
import Header from "@/app/components/platform/header";
import {
  PromptsApi,
  promptsUpdateFormApi,
} from "@/app/utils/ai-studio-endpoint/prompts-api";
import LoadingOverlay from "@/app/components/loading/loading";
import { toast } from "react-toastify";
import ToastContainerMessage from "@/app/components/notification/customToast";
import AIToolCallPopup from "@/app/studios/aitoolcallpoppup";

export default function PromptDetailsUpdate({ params }) {
  const unwrappedParams = use(params);
  const prompt_id = unwrappedParams?.promptID
    ? String(unwrappedParams.promptID).trim()
    : "";

  const router = useRouter();
  const [activeTab, setActiveTab] = useState("form");
  const [promptData, setPromptData] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [variableTags, setVariableTags] = useState([]);
  const [labelTags, setLabelTags] = useState([]);
  const [promptTypeTags, setPromptTypeTags] = useState([]);
  const [isToolCallModalOpen, setIsToolCallModalOpen] = useState(false);

  const userMessageRef = useRef(null);
  const systemMessageRef = useRef(null);

  // Form state
  const [formData, setFormData] = useState({
    name: "",
    promptTypes: [],
    spec: {
      userMessage: "",
      systemMessage: "",
      variables: [],
      sample: {
        inputText: "",
        response: "",
      },
      tools: [],
    },
    labels: [],
  });

  const openToolCallModal = () => {
    setIsToolCallModalOpen(true);
  };

  const closeToolCallModal = () => {
    setIsToolCallModalOpen(false);
  };

  const handleAddTool = (toolData) => {
    setFormData((prevData) => {
      const newTools = [toolData];
      return {
        ...prevData,
        spec: {
          ...prevData.spec,
          tools: newTools,
        },
      };
    });
  };

  // Fetch initial data on component mount
  useEffect(() => {
    const fetchData = async () => {
      try {
        setIsLoading(true);

        if (!prompt_id) {
          console.error("No prompt ID provided");
          toast.error("No prompt ID provided for update");
          return;
        }

        console.log("Fetching prompt data", prompt_id);
        const response = await PromptsApi.getPromptsId(prompt_id);
        console.log("Prompt data fetched:", response);
        const prompt =
          response.output?.[0] || response.data?.[0] || response[0];

        if (!prompt) {
          throw new Error("No prompt data found");
        }

        setPromptData(prompt);
        console.log("Prompt data:", prompt);
        console.log("tools:", prompt.spec?.tools);

        // Set initial form data
        setFormData({
          name: prompt.name || "",
          promptTypes: prompt.promptTypes || [],
          spec: {
            userMessage: prompt.spec?.userMessage || "",
            systemMessage: prompt.spec?.systemMessage || "",
            variables: prompt.spec?.variables || [],
            sample: {
              inputText: prompt.spec?.sample?.inputText || "",
              response: prompt.spec?.sample?.response || "",
            },
            tools: prompt.spec?.tools || [],
          },
          labels: prompt.labels || [],
          promptId: prompt_id,
        });

        // Set tags
        setVariableTags(prompt.spec?.variables || []);
        setLabelTags(prompt.labels || []);
        setPromptTypeTags(prompt.promptTypes || []);
      } catch (error) {
        console.error("Error fetching prompt data:", error);
        toast.error("Failed to fetch prompt data");
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, [prompt_id]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    if (name.includes(".")) {
      const [parent, child] = name.split(".");
      setFormData((prev) => ({
        ...prev,
        [parent]: {
          ...prev[parent],
          [child]: value,
        },
      }));
    } else {
      setFormData((prev) => ({
        ...prev,
        [name]: value,
      }));
    }
  };

  const handleSampleInputChange = (e) => {
    const { name, value } = e.target;
    const fieldName = name.split(".")[2];

    setFormData((prev) => ({
      ...prev,
      spec: {
        ...prev.spec,
        sample: {
          ...prev.spec.sample,
          [fieldName]: value,
        },
      },
    }));
  };

  const handleTagInput = (tagType, value) => {
    if (value.trim() === "") return;

    switch (tagType) {
      case "variable":
        setVariableTags((prev) => [...prev, value.trim()]);
        setFormData((prev) => ({
          ...prev,
          spec: {
            ...prev.spec,
            variables: [...prev.spec.variables, value.trim()],
          },
        }));
        break;
      case "label":
        setLabelTags((prev) => [...prev, value.trim()]);
        setFormData((prev) => ({
          ...prev,
          labels: [...prev.labels, value.trim()],
        }));
        break;
      case "promptType":
        setPromptTypeTags((prev) => [...prev, value.trim()]);
        setFormData((prev) => ({
          ...prev,
          promptTypes: [...prev.promptTypes, value.trim()],
        }));
        break;
      default:
        break;
    }
  };

  const removeTag = (tagType, indexToRemove) => {
    switch (tagType) {
      case "variable":
        setVariableTags((prev) =>
          prev.filter((_, index) => index !== indexToRemove)
        );
        setFormData((prev) => ({
          ...prev,
          spec: {
            ...prev.spec,
            variables: prev.spec.variables.filter(
              (_, index) => index !== indexToRemove
            ),
          },
        }));
        break;
      case "label":
        setLabelTags((prev) =>
          prev.filter((_, index) => index !== indexToRemove)
        );
        setFormData((prev) => ({
          ...prev,
          labels: prev.labels.filter((_, index) => index !== indexToRemove),
        }));
        break;
      case "promptType":
        setPromptTypeTags((prev) =>
          prev.filter((_, index) => index !== indexToRemove)
        );
        setFormData((prev) => ({
          ...prev,
          promptTypes: prev.promptTypes.filter(
            (_, index) => index !== indexToRemove
          ),
        }));
        break;
      default:
        break;
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      setIsLoading(true);

      const dataToSubmit = {
        promptId: prompt_id,
        name: formData.name,
        promptTypes: formData.promptTypes,
        spec: {
          variables: formData.spec.variables,
          userMessage: formData.spec.userMessage,
          systemMessage: formData.spec.systemMessage,
          sample: formData.spec.sample,
          tools: formData.spec.tools,
        },
        labels: formData.labels,
      };

      console.log("Submitting update with payload:", dataToSubmit);
      const payload = {
        spec: dataToSubmit,
      };

      const output = await promptsUpdateFormApi.getPromptsUpdteForm(payload);

      console.log("Resource updated:", output);

      toast.success("Prompt updated successfully");
      setTimeout(() => {
        if (prompt_id) {
          router.push(`/ai-center/prompts/${prompt_id}`);
        }
      }, 1000);
    } catch (error) {
      console.error("Error submitting form:", error);
      toast.error("Failed to update prompt");
    } finally {
      setIsLoading(false);
    }
  };

  // Variable dropdown handler for textareas (unchanged)
  useEffect(() => {
    const setupTextareaHandler = (textareaRef, fieldName) => {
      if (!textareaRef.current) return;

      const textarea = textareaRef.current;
      let lastCount = 0;
      let dropdown = null;
      let isInserting = false;

      const handleInput = () => {
        if (isInserting) return;

        const currentCount = (textarea.value.match(/\{\{/g) || []).length;
        if (currentCount > lastCount) {
          showOptionsDropdown(textarea);
        }
        lastCount = currentCount;
      };

      const showOptionsDropdown = (element) => {
        closeDropdown();

        // Check if we have variables to show
        if (!formData.spec.variables || formData.spec.variables.length === 0) {
          console.log("No variables defined");
          return;
        }

        // Create dropdown
        dropdown = document.createElement("div");
        dropdown.className =
          "absolute bg-purple-100 border border-gray-300 rounded-md shadow-lg z-50 max-h-48 overflow-y-auto";

        formData.spec.variables.forEach((option) => {
          const div = document.createElement("div");
          div.className =
            "px-4 py-2 hover:bg-purple-200 text-purple-800 cursor-pointer text-sm";
          div.textContent = option;
          div.onclick = (e) => {
            e.stopPropagation();
            insertOption(option, element, fieldName);
          };
          dropdown.appendChild(div);
        });

        // Get accurate cursor position
        const cursorPos = element.selectionStart;
        const cursorCoordinates = getCursorCoordinates(element, cursorPos);

        dropdown.style.position = "absolute";
        dropdown.style.top = `${cursorCoordinates.top + window.scrollY + 5}px`;
        dropdown.style.left = `${
          cursorCoordinates.left + window.scrollX + 5
        }px`;

        document.body.appendChild(dropdown);

        // Close on outside click
        setTimeout(() => {
          document.addEventListener("click", closeDropdown);
        }, 10);
      };

      const insertOption = (option, element, fieldName) => {
        isInserting = true;
        const cursorPos = element.selectionStart;

        const newValue = `${element.value.slice(
          0,
          cursorPos
        )}${option}}}${element.value.slice(cursorPos)}`;

        const fieldParts = fieldName.split(".");
        if (fieldParts.length === 2) {
          // Handle nested field like "spec.userMessage"
          setFormData((prev) => ({
            ...prev,
            [fieldParts[0]]: {
              ...prev[fieldParts[0]],
              [fieldParts[1]]: newValue,
            },
          }));
        } else {
          // Handle simple field name (fallback to original behavior)
          setFormData((prev) => ({
            ...prev,
            spec: {
              ...prev.spec,
              [fieldName]: newValue,
            },
          }));
        }

        element.value = newValue;

        requestAnimationFrame(() => {
          const newPos = cursorPos + option.length + 2;
          element.setSelectionRange(newPos, newPos);
          isInserting = false;
          element.focus();
        });

        closeDropdown();
      };

      const closeDropdown = () => {
        if (dropdown) {
          dropdown.remove();
          dropdown = null;
        }
        document.removeEventListener("click", closeDropdown);
      };

      // Accurate cursor position detection from original code
      const getCursorCoordinates = (element, position) => {
        const range = document.createRange();
        const textNode = document.createTextNode(element.value);
        const mirror = document.createElement("div");
        const style = getComputedStyle(element);
        const rect = element.getBoundingClientRect();

        // Mirror all critical styles
        mirror.style.cssText = `
          position: absolute;
          white-space: pre-wrap;
          font: ${style.font};
          left: ${rect.left}px;
          top: ${rect.top}px;
          width: ${element.clientWidth}px;
          padding: ${style.padding};
          margin: ${style.margin};
          border: ${style.border};
          box-sizing: border-box;
          visibility: hidden;
          overflow-wrap: break-word;
          z-index: -9999;
        `;

        // Sync scroll state
        mirror.scrollTop = element.scrollTop;
        mirror.scrollLeft = element.scrollLeft;
        mirror.appendChild(textNode);
        document.body.appendChild(mirror);

        // Set range position
        range.setStart(textNode, Math.min(position, textNode.length));
        range.setEnd(textNode, Math.min(position, textNode.length));
        const caretRect = range.getBoundingClientRect();
        const coordinates = {
          top: caretRect.bottom,
          left: caretRect.left,
        };
        document.body.removeChild(mirror);
        return coordinates;
      };

      // Attach event listener to textarea
      textarea.addEventListener("input", handleInput);
      return () => {
        textarea.removeEventListener("input", handleInput);
        closeDropdown();
      };
    };

    // Set up handlers for both textareas
    const userMessageCleanup = userMessageRef.current
      ? setupTextareaHandler(userMessageRef, "spec.userMessage")
      : () => {};

    const systemMessageCleanup = systemMessageRef.current
      ? setupTextareaHandler(systemMessageRef, "spec.systemMessage")
      : () => {};

    return () => {
      userMessageCleanup();
      systemMessageCleanup();
    };
  }, [formData.spec.variables, setFormData]);

  return (
    <div className="bg-zinc-800 flex h-screen">
      <Sidebar />

      <div className="overflow-y-auto  scrollbar h-screen w-full">
        <Header
          sectionHeader="Update Prompt"
          hideBackListingLink={false}
          backListingLink="./"
        />
        <LoadingOverlay isLoading={isLoading} />
        <ToastContainerMessage />

        <div className="flex-grow p-4 overflow-y-auto w-full">
          <section id="grids" className="space-y-2">
            <div className="max-w-6xl mx-auto bg-[#1b1b1b] shadow rounded-lg p-2">
              <div className="text-gray-100 mb-2 flex justify-center">
                {/* <button
                  id="yamlSpecB"
                  className={`whitespace-nowrap border-b-2 py-2 px-2 text-md font-medium focus:outline-none ${
                    activeTab === "yaml"
                      ? "border-orange-700 text-orange-700 font-semibold"
                      : "border-black"
                  }`}
                  onClick={() => setActiveTab("yaml")}
                >
                  YAML
                </button> */}
                <button
                  id="formSpecB"
                  className={`whitespace-nowrap border-b-2 py-2 px-2 text-md font-medium focus:outline-none ml-4 ${
                    activeTab === "form"
                      ? "border-orange-700 text-orange-700 font-semibold"
                      : "border-black"
                  }`}
                  onClick={() => setActiveTab("form")}
                  suppressHydrationWarning
                >
                  Form
                </button>
              </div>

              <div
                id="yamlSpec"
                style={{ display: activeTab === "yaml" ? "block" : "none" }}
              >
                {/* YAML editor would go here */}
              </div>

              <div
                id="formSpec"
                style={{ display: activeTab === "form" ? "block" : "none" }}
              >
                <form
                  id="vapusPromptSpec"
                  className="space-y-2 border border-zinc-500 rounded-md text-gray-100 p-2"
                  onSubmit={handleSubmit}
                >
                  <fieldset className="p-4 rounded">
                    <div className="grid grid-cols-2 md:grid-cols-2 gap-4 mb-4">
                      {/* Name */}
                      <div>
                        <label htmlFor="name" className="labels">
                          Name
                        </label>
                        <input
                          id="name"
                          name="name"
                          value={formData.name}
                          onChange={handleInputChange}
                          type="text"
                          placeholder="Enter name"
                          className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                          suppressHydrationWarning
                        />
                      </div>

                      {/* promptTypes */}
                      <div>
                        <label
                          htmlFor="promptTypes-input-tag"
                          className="labels"
                        >
                          Prompt Types
                        </label>

                        <input
                          type="text"
                          id="promptTypes-input-tag"
                          placeholder="Enter Prompt Types"
                          suppressHydrationWarning
                          className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                          onKeyDown={(e) => {
                            if (e.key === "Enter") {
                              e.preventDefault();
                              handleTagInput("promptType", e.target.value);
                              e.target.value = "";
                            }
                          }}
                        />

                        <ul
                          id="promptTypes-tags"
                          className="flex flex-wrap gap-2 mt-1"
                        >
                          {promptTypeTags.map((tag, index) => (
                            <li
                              key={index}
                              className="text-purple-800 bg-purple-100 px-2 py-1 rounded-lg flex items-center"
                            >
                              {tag}
                              <span
                                className="ml-2 cursor-pointer"
                                onClick={() => removeTag("promptType", index)}
                              >
                                &times;
                              </span>
                            </li>
                          ))}
                        </ul>
                      </div>
                    </div>

                    <fieldset className="rounded mb-4 border border-zinc-600 rounded-md shadow-sm p-4">
                      {/* variables */}
                      <div className="col-span-2">
                        <label htmlFor="variable-input-tag" className="labels">
                          Variables
                        </label>
                        <div>
                          <input
                            type="text"
                            id="variable-input-tag"
                            placeholder="Enter Variables"
                            className="w-full rounded-md bg-zinc-700 text-gray-100 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-2 w-full mt-2"
                            suppressHydrationWarning
                            onKeyDown={(e) => {
                              if (e.key === "Enter") {
                                e.preventDefault();
                                handleTagInput("variable", e.target.value);
                                e.target.value = "";
                              }
                            }}
                          />

                          <ul id="tags" className="flex flex-wrap gap-2 mt-1">
                            {variableTags.map((tag, index) => (
                              <li
                                key={index}
                                className="text-purple-800 bg-purple-100 px-2 py-1 rounded-lg flex items-center"
                              >
                                {tag}
                                <span
                                  className="ml-2 cursor-pointer"
                                  onClick={() => removeTag("variable", index)}
                                >
                                  &times;
                                </span>
                              </li>
                            ))}
                          </ul>
                        </div>
                      </div>

                      <div className="grid grid-cols-1 mt-6 md:grid-cols-1 gap-4">
                        {/* userMessage */}
                        <div className="col-span-2">
                          <label htmlFor="myTextarea-1" className="labels">
                            User Message
                          </label>
                          <textarea
                            name="spec.userMessage"
                            id="myTextarea-1"
                            ref={userMessageRef}
                            rows="4"
                            value={formData.spec.userMessage}
                            onChange={handleInputChange}
                            placeholder="Enter User Message"
                            className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-300 p-2 shadow-sm w-full"
                            suppressHydrationWarning
                          ></textarea>
                        </div>
                      </div>

                      <div className="grid grid-cols-1 mt-6 md:grid-cols-1 gap-4">
                        {/* systemMessage */}
                        <div className="col-span-2">
                          <label htmlFor="myTextarea-2" className="labels">
                            System Message
                          </label>
                          <textarea
                            name="spec.systemMessage"
                            ref={systemMessageRef}
                            id="myTextarea-2"
                            rows="4"
                            value={formData.spec.systemMessage}
                            onChange={handleInputChange}
                            placeholder="Enter System Message"
                            className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-300 p-2 shadow-sm w-full"
                          ></textarea>
                        </div>
                      </div>

                      <fieldset className="rounded mt-6 mb-4">
                        <legend className="text-lg font-semibold text-gray-100">
                          Sample
                        </legend>
                        <div className="grid grid-cols-2 md:grid-cols-2 gap-4">
                          {/* inputText */}
                          <div>
                            <label
                              htmlFor="sample-input-text"
                              className="labels"
                            >
                              Input Text
                            </label>
                            <input
                              id="sample-input-text"
                              name="spec.sample.inputText"
                              value={formData.spec.sample.inputText}
                              onChange={handleSampleInputChange}
                              type="text"
                              placeholder="Enter Input Text"
                              suppressHydrationWarning
                              className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                            />
                          </div>
                          {/* response */}
                          <div>
                            <label htmlFor="sample-response" className="labels">
                              Response
                            </label>
                            <input
                              id="sample-response"
                              name="spec.sample.response"
                              value={formData.spec.sample.response}
                              onChange={handleSampleInputChange}
                              type="text"
                              placeholder="Enter Response"
                              className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                              suppressHydrationWarning
                            />
                          </div>
                        </div>
                      </fieldset>

                      {/* Label */}
                      <div
                        id="Labels"
                        className="grid grid-cols-1 md:grid-cols-1 gap-4"
                      >
                        <div className="col-span-2">
                          <label htmlFor="label-input-tag" className="labels">
                            Labels
                          </label>

                          <input
                            type="text"
                            id="label-input-tag"
                            placeholder="Enter Labels"
                            className="w-full rounded-md bg-zinc-700 text-gray-100 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-2 w-full mt-2"
                            suppressHydrationWarning
                            onKeyDown={(e) => {
                              if (e.key === "Enter") {
                                e.preventDefault();
                                handleTagInput("label", e.target.value);
                                e.target.value = "";
                              }
                            }}
                          />

                          <ul
                            id="label-tags"
                            className="flex flex-wrap gap-2 mt-1"
                          >
                            {labelTags.map((tag, index) => (
                              <li
                                key={index}
                                className="text-yellow-800 bg-yellow-100 px-2 py-1 rounded-lg flex items-center"
                              >
                                {tag}
                                <span
                                  className="ml-2 cursor-pointer"
                                  onClick={() => removeTag("label", index)}
                                >
                                  &times;
                                </span>
                              </li>
                            ))}
                          </ul>
                        </div>
                      </div>

                      {/* Modal */}
                      <div id="toolCall" className="flex mt-2 mb-2">
                        <button
                          id="addToolCallButton"
                          type="button"
                          onClick={openToolCallModal}
                          className="bg-orange-700 w-full px-2 py-2 text-sm rounded-lg focus:outline-none cursor-pointer"
                          suppressHydrationWarning
                        >
                          Add Tool Call
                        </button>
                        <AIToolCallPopup
                          isOpen={isToolCallModalOpen}
                          onClose={closeToolCallModal}
                          onAddTool={handleAddTool}
                          existingTools={
                            formData.spec.tools.length > 0
                              ? formData.spec.tools[0]
                              : null
                          }
                          editingIndex={0}
                        />
                      </div>
                    </fieldset>
                  </fieldset>

                  {/* Submit Button */}
                  <div className="mt-4 flex justify-end space-x-2">
                    <button
                      type="submit"
                      className={`px-6 py-2 bg-orange-700 text-white rounded-md shadow hover:bg-pink-900 ${
                        isLoading ? "hidden" : ""
                      }`}
                    >
                      Submit
                    </button>
                    <button
                      type="button"
                      className={`px-6 py-2 bg-orange-700 text-white rounded-md shadow hover:bg-pink-900 ${
                        isLoading ? "" : "hidden"
                      }`}
                      suppressHydrationWarning
                    >
                      <svg
                        className="animate-spin h-6 w-6 text-white"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                      >
                        <circle
                          className="opacity-25"
                          cx="12"
                          cy="12"
                          r="10"
                          stroke="currentColor"
                          strokeWidth="4"
                        ></circle>
                        <path
                          className="opacity-75"
                          fill="currentColor"
                          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
                        ></path>
                      </svg>
                    </button>
                  </div>
                </form>
              </div>
            </div>
          </section>
        </div>
      </div>
    </div>
  );
}
