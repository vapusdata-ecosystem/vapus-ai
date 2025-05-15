"use client";
import Header from "@/app/components/platform/header";
import LoadingOverlay from "@/app/components/loading/loading";
import Sidebar from "@/app/components/platform/main-sidebar";
import ToastContainerMessage from "@/app/components/notification/customToast";
import AIToolCallPopup from "@/app/studios/aitoolcallpoppup";
import { promptsFormApi } from "@/app/utils/ai-studio-endpoint/prompts-api";
import { useRouter } from "next/navigation";
import { useState, useEffect, useRef } from "react";
import { toast } from "react-toastify";

export default function PromptForm() {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState("form");
  const [isLoading, setIsLoading] = useState(false);
  const [isToolCallModalOpen, setIsToolCallModalOpen] = useState(false);
  const [formData, setFormData] = useState({
    name: "",
    promptTypes: [],
    spec: {
      variables: [],
      userMessage: "",
      systemMessage: "",
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

  // Fix for handleAddTool in PromptForm.js
  const handleAddTool = (toolData) => {
    console.log("Tool data received:", toolData);

    setFormData((prevData) => {
      const updatedTools = [...prevData.spec.tools, toolData];
      console.log("Updated tools array:", updatedTools);

      return {
        ...prevData,
        spec: {
          ...prevData.spec,
          tools: updatedTools,
        },
      };
    });
  };

  // Refs for variable management
  const variableInputRef = useRef(null);
  const promptTypesInputRef = useRef(null);
  const labelInputRef = useRef(null);
  const userMessageRef = useRef(null);
  const systemMessageRef = useRef(null);

  // Handle form submission
  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      setIsLoading(true);

      const dataToSubmit = {
        name: formData.name,
        promptTypes: formData.promptTypes,
        spec: {
          variables: formData.spec.variables,
          userMessage: formData.spec.userMessage,
          systemMessage: formData.spec.systemMessage,
          sample: {
            inputText: formData.spec.sample.inputText,
            response: formData.spec.sample.response,
          },
          tools: formData.spec.tools,
        },
        labels: formData.labels,
      };

      console.log("Form Data:", dataToSubmit);

      // payload
      const payload = {
        spec: dataToSubmit,
      };

      const output = await promptsFormApi.getPromptsForm(payload);
      const resourceInfo = output.result;

      if (resourceInfo) {
        toast.success("Prompt created successfully.");
        setTimeout(() => {
          router.push(`/ai-center/prompts`);
        }, 1000);
      } else {
        toast.success("Resource Created", "Resource created successfully.");
      }
    } catch (error) {
      console.error("Error sending API request:", error);
      toast.error("Failed to create Prompt. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  const addTag = (tagType, tagContent) => {
    if (!tagContent.trim()) return;

    // Tag limit check - max 4 tags as in original code
    const currentTags =
      tagType === "variables"
        ? formData.spec.variables
        : tagType === "promptTypes"
        ? formData.promptTypes
        : formData.labels;

    if (currentTags.length >= 4) {
      toast.error("Tag Limit, Field cannot be greater than 4");
      return;
    }

    setFormData((prev) => {
      if (tagType === "variables") {
        return {
          ...prev,
          spec: {
            ...prev.spec,
            variables: [...prev.spec.variables, tagContent.trim()],
          },
        };
      } else if (tagType === "promptTypes") {
        return {
          ...prev,
          promptTypes: [...prev.promptTypes, tagContent.trim()],
        };
      } else if (tagType === "labels") {
        return {
          ...prev,
          labels: [...prev.labels, tagContent.trim()],
        };
      }
      return prev;
    });
  };

  const removeTag = (tagType, indexToRemove) => {
    setFormData((prev) => {
      if (tagType === "variables") {
        return {
          ...prev,
          spec: {
            ...prev.spec,
            variables: prev.spec.variables.filter(
              (_, index) => index !== indexToRemove
            ),
          },
        };
      } else if (tagType === "promptTypes") {
        return {
          ...prev,
          promptTypes: prev.promptTypes.filter(
            (_, index) => index !== indexToRemove
          ),
        };
      } else if (tagType === "labels") {
        return {
          ...prev,
          labels: prev.labels.filter((_, index) => index !== indexToRemove),
        };
      }
      return prev;
    });
  };

  // Variable dropdown handler for textareas
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

        // Update React state
        setFormData((prev) => ({
          ...prev,
          spec: {
            ...prev.spec,
            [fieldName]: newValue,
          },
        }));
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
      ? setupTextareaHandler(userMessageRef, "userMessage")
      : () => {};

    const systemMessageCleanup = systemMessageRef.current
      ? setupTextareaHandler(systemMessageRef, "systemMessage")
      : () => {};

    return () => {
      userMessageCleanup();
      systemMessageCleanup();
    };
  }, [formData.spec.variables, setFormData]);

  // Handle tag input keydown events
  const handleTagInputKeyDown = (e, type) => {
    if (e.key === "Enter") {
      e.preventDefault();
      addTag(type, e.target.value);
      e.target.value = "";
    }
  };

  return (
    <div className="bg-zinc-800 flex h-screen">
      {/* Notification component */}
      <Sidebar />
      {/* Main content */}
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="Build Your Prompt"
          hideBackListingLink={false}
          backListingLink="./"
        />
        <ToastContainerMessage />
        <LoadingOverlay isLoading={isLoading} />
        <div className="flex-grow p-4 overflow-y-auto w-full">
          <section id="grids" className="space-y-2">
            <div className="max-w-6xl mx-auto bg-[#1b1b1b] shadow rounded-lg p-2">
              <div className="border-b border-zinc-500 mb-2 flex justify-center">
                {/* <button
                  onClick={() => setActiveTab("yaml")}
                  className={`whitespace-nowrap text-gray-100 border-b-2 py-2 px-2 text-md font-medium focus:outline-none ${
                    activeTab === "yaml"
                      ? "border-orange-700 text-orange-700 font-semibold"
                      : "border-black"
                  }`}
                >
                  YAML
                </button> */}
                <button
                  onClick={() => setActiveTab("form")}
                  className={`whitespace-nowrap text-gray-100 border-b-2 py-2 px-2 text-md font-medium focus:outline-none ml-4 ${
                    activeTab === "form"
                      ? "border-orange-700 text-orange-700 font-semibold"
                      : "border-black"
                  }`}
                  suppressHydrationWarning
                >
                  Form
                </button>
              </div>

              {activeTab === "yaml" ? (
                <div id="yamlSpec">
                  <div className="p-4 bg-zinc-700 rounded text-gray-100">
                    YAML Editor placeholder - would need the actual YAML editor
                    component
                  </div>
                </div>
              ) : (
                <div id="formSpec">
                  <form
                    onSubmit={handleSubmit}
                    className="space-y-2 border border-zinc-500 rounded-md text-gray-100 p-2"
                  >
                    <fieldset className="p-4 rounded">
                      <div className="grid grid-cols-2 md:grid-cols-2 gap-4 mb-4">
                        {/* Name */}
                        <div>
                          <label className="labels">Name</label>
                          <input
                            id="name"
                            name="name"
                            type="text"
                            placeholder="Enter name"
                            className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                            value={formData.name}
                            onChange={(e) =>
                              setFormData({ ...formData, name: e.target.value })
                            }
                            suppressHydrationWarning
                          />
                        </div>

                        {/* promptTypes */}
                        <div>
                          <label className="labels">Prompt Types</label>

                          <input
                            type="text"
                            ref={promptTypesInputRef}
                            placeholder="Enter Prompt Types"
                            className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                            onKeyDown={(e) =>
                              handleTagInputKeyDown(e, "promptTypes")
                            }
                            suppressHydrationWarning
                          />

                          <ul className="flex flex-wrap gap-2 mt-1">
                            {formData.promptTypes.map((tag, index) => (
                              <li
                                key={index}
                                className="text-purple-800 bg-purple-100 px-2 py-1 rounded-lg flex items-center"
                              >
                                {tag}
                                <button
                                  type="button"
                                  className="ml-2 text-purple-800 cursor-pointer"
                                  onClick={() =>
                                    removeTag("promptTypes", index)
                                  }
                                >
                                  ×
                                </button>
                              </li>
                            ))}
                          </ul>
                        </div>
                      </div>

                      <fieldset className="rounded mb-4 border border-zinc-600 rounded-md shadow-sm p-4">
                        <legend className="text-lg font-semibold text-gray-100">
                          Spec
                        </legend>

                        {/* Variables */}
                        <div className="col-span-2">
                          <label className="labels">Variable</label>
                          <div>
                            <input
                              type="text"
                              ref={variableInputRef}
                              placeholder="Enter Variables"
                              className="form-input-field rounded-md p-2 bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 w-full mt-2"
                              onKeyDown={(e) =>
                                handleTagInputKeyDown(e, "variables")
                              }
                              suppressHydrationWarning
                            />
                            <ul className="flex flex-wrap gap-2 mt-1">
                              {formData.spec.variables.map((tag, index) => (
                                <li
                                  key={index}
                                  className="text-purple-800 bg-purple-100 px-2 py-1 rounded-lg flex items-center"
                                >
                                  {tag}
                                  <button
                                    type="button"
                                    className="ml-2 text-purple-800 cursor-pointer"
                                    onClick={() =>
                                      removeTag("variables", index)
                                    }
                                  >
                                    ×
                                  </button>
                                </li>
                              ))}
                            </ul>
                          </div>
                        </div>

                        <div className="grid grid-cols-1 mt-6 md:grid-cols-1 gap-4">
                          {/* userMessage */}
                          <div className="col-span-2">
                            <label className="labels">User Message</label>
                            <textarea
                              id="myTextarea-2"
                              ref={userMessageRef}
                              name="userMessage"
                              rows="4"
                              placeholder="Enter User Message"
                              className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-300 p-2 shadow-sm w-full"
                              value={formData.spec.userMessage}
                              onChange={(e) =>
                                setFormData({
                                  ...formData,
                                  spec: {
                                    ...formData.spec,
                                    userMessage: e.target.value,
                                  },
                                })
                              }
                            />
                          </div>
                        </div>

                        <div className="grid grid-cols-1 mt-6 md:grid-cols-1 gap-4">
                          {/* systemMessage */}
                          <div className="col-span-2">
                            <label className="labels">System Message</label>
                            <textarea
                              id="myTextarea-1"
                              ref={systemMessageRef}
                              name="systemMessage"
                              rows="4"
                              placeholder="Enter System Message"
                              className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-300 p-2 shadow-sm w-full"
                              value={formData.spec.systemMessage}
                              onChange={(e) =>
                                setFormData({
                                  ...formData,
                                  spec: {
                                    ...formData.spec,
                                    systemMessage: e.target.value,
                                  },
                                })
                              }
                            />
                          </div>
                        </div>

                        <fieldset className="rounded mt-6 mb-4">
                          <legend className="text-lg font-semibold text-gray-100">
                            Sample
                          </legend>
                          <div className="grid grid-cols-2 md:grid-cols-2 gap-4">
                            {/* inputText */}
                            <div>
                              <label className="labels">Input Text</label>
                              <input
                                name="inputText"
                                type="text"
                                placeholder="Enter Input Text"
                                className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                value={formData.spec.sample.inputText}
                                suppressHydrationWarning
                                onChange={(e) =>
                                  setFormData({
                                    ...formData,
                                    spec: {
                                      ...formData.spec,
                                      sample: {
                                        ...formData.spec.sample,
                                        inputText: e.target.value,
                                      },
                                    },
                                  })
                                }
                              />
                            </div>
                            {/* response */}
                            <div>
                              <label className="labels">Response</label>
                              <input
                                name="response"
                                type="text"
                                placeholder="Enter Response"
                                className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                value={formData.spec.sample.response}
                                suppressHydrationWarning
                                onChange={(e) =>
                                  setFormData({
                                    ...formData,
                                    spec: {
                                      ...formData.spec,
                                      sample: {
                                        ...formData.spec.sample,
                                        response: e.target.value,
                                      },
                                    },
                                  })
                                }
                              />
                            </div>
                          </div>
                        </fieldset>

                        {/* Labels */}
                        <div
                          id="Labels"
                          className="grid grid-cols-1 md:grid-cols-1 gap-4"
                        >
                          <div className="col-span-2">
                            <label className="labels">Labels</label>

                            <input
                              type="text"
                              ref={labelInputRef}
                              placeholder="Enter Labels"
                              className="form-input-field rounded-md p-2 bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 w-full mt-2"
                              onKeyDown={(e) =>
                                handleTagInputKeyDown(e, "labels")
                              }
                              suppressHydrationWarning
                            />

                            <ul className="flex flex-wrap gap-2 mt-1">
                              {formData.labels.map((tag, index) => (
                                <li
                                  key={index}
                                  className="text-yellow-800 bg-yellow-100 px-2 py-1 rounded-lg flex items-center"
                                >
                                  {tag}
                                  <button
                                    type="button"
                                    className="ml-2 text-yellow-800 cursor-pointer"
                                    onClick={() => removeTag("labels", index)}
                                  >
                                    ×
                                  </button>
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
                          />
                        </div>
                      </fieldset>
                    </fieldset>

                    {/* Submit Button */}
                    <div className="mt-4 flex justify-end space-x-2">
                      {!isLoading ? (
                        <button
                          type="submit"
                          className="px-6 py-2 bg-orange-700 text-white rounded-md shadow hover:bg-pink-900"
                          suppressHydrationWarning
                        >
                          Submit
                        </button>
                      ) : (
                        <button
                          type="button"
                          className="px-6 py-2 bg-orange-700 text-white rounded-md shadow"
                          disabled
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
                      )}
                    </div>
                  </form>
                </div>
              )}
            </div>
          </section>
        </div>
      </div>
    </div>
  );
}
