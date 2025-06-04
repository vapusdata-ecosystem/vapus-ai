"use client";
import Head from "next/head";
import { useEffect, useState } from "react";
import Script from "next/script";
import Header from "@/app/components/platform/header";
import AIToolCallPopup from "../aitoolcallpoppup";
import NestedDropdown from "@/app/components/nestedDropdown";
import PromptDropdown from "@/app/components/promptDropdown";
import {
  addContextLocally,
  aiInterfaceAction,
} from "@/app/components/JS/ai-chat";
import { toast } from "react-toastify";
import ContextModal from "../ai-context-popup";

export default function AIStudio({ response, globalContext, aiStudioChat }) {
  const [isToolCallModalOpen, setIsToolCallModalOpen] = useState(false);
  const [isContextModalOpen, setIsContextModalOpen] = useState(false);
  
  const openToolCallModal = () => {
    setIsToolCallModalOpen(true);
  };
  
  const closeToolCallModal = () => {
    setIsToolCallModalOpen(false);
  };

  const handleAddTool = (toolData) => {
    setFormData((prevData) => ({
      ...prevData,
      spec: {
        ...prevData.spec,
        Tools: [toolData],
      },
    }));
    toast.success("Tool added successfully");
  };

  useEffect(() => {
    if (aiStudioChat) {
      loadAIStudioChat("aiChatResult");
      const cChatId = aiStudioChat.ChatId;

      if (cChatId !== "") {
        const radio = document.querySelector(
          'input[name="aiInteractionMode"][value="CHAT_MODE"]'
        );
        if (radio) {
          radio.checked = true;
        }
        handleInteractionModeChange("CHAT_MODE");
      } else {
        const radio = document.querySelector(
          'input[name="aiInteractionMode"][value="P2P"]'
        );
        if (radio) {
          radio.checked = true;
        }
        handleInteractionModeChange("P2P");
      }
    }
  }, [aiStudioChat]);

  function handleInteractionModeChange(divId) {
    document.getElementsByName("aiInteractionMode").forEach((element) => {
      if (element.value === divId) {
        document.getElementById(divId).classList.remove("hidden");
      } else {
        document.getElementById(element.value).classList.add("hidden");
      }
    });
  }

  function createNewChat() {
    const currentUrl = window.location.href;
    const urlObj = new URL(currentUrl);
    urlObj.searchParams.set("createNewChat", "true");
    window.location.href = urlObj.toString();
  }

  function closeOptionsBar() {
    const rightbar = document.getElementById("rightbar");
    rightbar.classList.add("hidden");
  }

  function updateTemperatureValue(value) {
    document.getElementById("temperatureValue").value =
      parseFloat(value).toFixed(1);
  }

  function updateSliderValue(value) {
    const numericValue = parseFloat(value);
    if (numericValue >= 0.0 && numericValue <= 1.0) {
      document.getElementById("temperatureSlider").value =
        numericValue.toFixed(1);
    }
  }

  function EnterInput(event) {
    if (event.key === "Enter") {
      submitInput();
    }
  }

  function toggleContextPopup() {
    setIsContextModalOpen(!isContextModalOpen);
  }

  function toggleStream(element) {
    const button =
      element.tagName === "BUTTON" ? element : element.closest("button");
    if (!button) return;
    const isOn = button.getAttribute("aria-checked") === "true";
    const newState = !isOn;
    button.classList.toggle("bg-orange-700", newState);
    button.classList.toggle("bg-gray-300", !newState);
    const indicator = button.querySelector("span");
    if (indicator)
      indicator.style.transform = newState
        ? "translateX(32px)"
        : "translateX(0)";
    const label = document.getElementById("toggleStateLabel");
    if (label) label.textContent = newState ? "Stream ON" : "Stream OFF";
    button.setAttribute("aria-checked", newState.toString());
  }

  function submitInput() {
    const inputArea = document.getElementById("input");
    const textInput = inputArea.value;
    const aiModel = document.getElementById("aiModel").value;
    const temperature = document.getElementById("temperatureValue").value;
    const maxTokens = document.getElementById("maxTokens").value;
    const topk = document.getElementById("topk").value;
    const topP = document.getElementById("topP").value;
    const promptId = document.getElementById("promptId").value;
    const contextType = document.getElementById("contextType").value;
    const contextValue = document.getElementById("contextValue").value;

    let modelName;
    let modelNodeId;
    let splitVal = aiModel.split("||");
    modelNodeId = splitVal[0];
    modelName = splitVal[1];

    if (modelName === "" || modelNodeId === "") {
      return;
    }
    inputArea.value = "";
    inputArea.disabled = true;

    addContextLocally(
      {
        content: textInput,
        userId: globalContext.UserInfo.UserId,
        domain: globalContext.CurrentDomain.DomainId,
      },
      true
    );

    aiInterfaceAction(
      response.ActionParams.API,
      response.ActionParams.ChatAPI,
      globalContext.AccessTokenKey,
      "aiChatResult",
      contextType,
      contextValue,
      modelNodeId,
      promptId,
      textInput,
      temperature,
      topP,
      modelName,
      maxTokens,
      topk
    );

    inputArea.disabled = false;
    document.getElementById("contextType").value = "";
    document.getElementById("contextValue").value = "";
  }

  return (
    <>
      <Head>
        <meta charSet="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>AI Studio</title>

        {/* External Stylesheets */}
        <link
          href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism.min.css"
          rel="stylesheet"
        />
        <link
          href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism-okaidia.min.css"
          rel="stylesheet"
        />
        <link
          href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/styles/atom-one-dark.min.css"
          rel="stylesheet"
        />
        <link
          rel="stylesheet"
          href="https://cdnjs.cloudflare.com/ajax/libs/jsoneditor/9.10.0/jsoneditor.min.css"
        />

        <style jsx global>{`
          .modal {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: rgba(0, 0, 0, 0.5);
            align-items: center;
            justify-content: center;
          }

          .modal.active {
            display: flex;
          }

          .rightbar {
            transition: transform 0.3s ease-in-out;
          }
          .rightbar.hidden {
            transform: translateX(100%);
          }
          .rightbar.visible {
            transform: translateX(0);
          }
        `}</style>
      </Head>

      {/* External Script Tags */}
      <Script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></Script>
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/jsoneditor/9.10.0/jsoneditor.min.js" />
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/highlight.min.js" />
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/languages/go.min.js" />
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/prism.min.js" />
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-python.min.js" />
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-javascript.min.js" />
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-sql.min.js" />
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-go.min.js" />
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/plugins/autoloader/prism-autoloader.min.js" />

      <div className="bg-zinc-800 flex h-screen">
        <div className="h-screen w-full">
          <Header
            sectionHeader="AI Studio Interface"
            hideBackListingLink={true}
            backListingLink="./"
          />

          <button
            id="optionsIcon"
            className="fixed bottom-4 right-4 bg-orange-700 p-3 rounded-full shadow-lg focus:outline-none z-50 cursor-pointer"
            onClick={() => {
              const rightbar = document.getElementById("rightbar");
              const isHidden = rightbar.classList.contains("hidden");
              rightbar.classList.toggle("hidden", !isHidden);
              rightbar.classList.toggle("visible", isHidden);
            }}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth="2"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M4 6h16M4 12h16M4 18h16"
              />
            </svg>
          </button>

          {/* Chat and Editor Section */}
          <div id="aistudiodiv" className="flex flex-1">
            <div
              className="flex bg-zinc-800 p-2 flex-col w-full overflow-y-auto overflow-x-auto scrollbar"
              style={{ height: "calc(100vh - 65px)" }}
            >
              <div
                id="aiChatResult"
                className="ml-8 mr-8 flex-grow rounded-lg bg-gray-50 table-container overflow-y-auto scrollbar"
                style={{
                  height: "500px",
                  background:
                    "radial-gradient(circle, black -25%, transparent 10%) 0 0, radial-gradient(circle, black -25%, transparent 10%) 10px 10px",
                  backgroundSize: "20px 20px",
                }}
              ></div>

              {/* Input Area with Dropdown */}
              <div
                id="userInputGround"
                className="ml-8 mr-8 mt-2 bg-[#1b1b1b] p-1 rounded-lg shadow-md border border-zinc-500"
              >
                <textarea
                  id="input"
                  onKeyDown={EnterInput}
                  className="w-full h-16 p-1 rounded-lg focus:outline-none text-gray-100 overflow-y-auto scrollbar"
                  placeholder="Type your message here..."
                ></textarea>

                <textarea
                  id="promptInput"
                  onKeyDown={EnterInput}
                  className="w-full h-16 p-1 rounded-lg focus:outline-none text-gray-100 hidden"
                  placeholder="Type your message here..."
                ></textarea>

                <div className="flex h-[26px] items-center justify-between">
                  <div className="flex gap-x-1">
                    <button
                      title="View selected dataproduct prompts"
                      aria-label="datafiles"
                      data-testid="datafiles"
                      className="flex cursor-pointer h-6 w-6 items-center justify-center rounded-full transition-colors hover:bg-zinc-900 hover:opacity-70 focus-visible:outline-none focus-visible:outline-black disabled:text-[#f4f4f4] disabled:hover:opacity-100 disabled:dark:bg-token-text-quaternary bg-orange-700 disabled:bg-[#D7D7D7]"
                      onClick={toggleContextPopup}
                    >
                      <svg viewBox="0 0 24 24" className="h-5 w-5 shrink-0">
                        <path
                          fill="currentColor"
                          d="M16.5 6v11.5a4 4 0 0 1-4 4a4 4 0 0 1-4-4V5A2.5 2.5 0 0 1 11 2.5A2.5 2.5 0 0 1 13.5 5v10.5a1 1 0 0 1-1 1a1 1 0 0 1-1-1V6H10v9.5a2.5 2.5 0 0 0 2.5 2.5a2.5 2.5 0 0 0 2.5-2.5V5a4 4 0 0 0-4-4a4 4 0 0 0-4 4v12.5a5.5 5.5 0 0 0 5.5 5.5a5.5 5.5 0 0 0 5.5-5.5V6z"
                        ></path>
                      </svg>
                    </button>
                  </div>

                  <button
                    aria-label="startChat"
                    data-testid="startChat"
                    id="startChat"
                    className="flex cursor-pointer h-6 w-6 items-center justify-center rounded-full transition-colors hover:opacity-70 focus-visible:outline-none focus-visible:outline-black disabled:text-[#f4f4f4] disabled:hover:opacity-100 disabled:dark:bg-token-text-quaternary bg-orange-700 disabled:bg-[#D7D7D7]"
                    onClick={submitInput}
                  >
                    <svg
                      viewBox="0 0 32 32"
                      fill="none"
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-5 w-5"
                    >
                      <path
                        fillRule="evenodd"
                        clipRule="evenodd"
                        d="M15.1918 8.90615C15.6381 8.45983 16.3618 8.45983 16.8081 8.90615L21.9509 14.049C22.3972 14.4953 22.3972 15.2189 21.9509 15.6652C21.5046 16.1116 20.781 16.1116 20.3347 15.6652L17.1428 12.4734V22.2857C17.1428 22.9169 16.6311 23.4286 15.9999 23.4286C15.3688 23.4286 14.8571 22.9169 14.8571 22.2857V12.4734L11.6652 15.6652C11.2189 16.1116 10.4953 16.1116 10.049 15.6652C9.60265 15.2189 9.60265 14.4953 10.049 14.049L15.1918 8.90615Z"
                        fill="currentColor"
                      ></path>
                    </svg>
                  </button>

                  <button
                    aria-label="inprogressChat"
                    data-testid="inprogressChat"
                    id="inprogressChat"
                    className="hidden flex h-6 w-6 items-center justify-center rounded-full transition-colors hover:opacity-70 focus-visible:outline-none focus-visible:outline-black disabled:text-[#f4f4f4] disabled:hover:opacity-100 disabled:dark:bg-token-text-quaternary bg-orange-700 disabled:bg-[#D7D7D7]"
                    onClick={() => cancelStream()}
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-5 w-5"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="2"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                    >
                      <rect
                        x="9"
                        y="9"
                        width="6"
                        height="6"
                        fill="currentColor"
                      />
                    </svg>
                  </button>
                </div>
              </div>
            </div>

            <div className="w-[220px] min-w-[220px] bg-[#1b1b1b] p-2 overflow-y-auto scrollbar text-gray-100">
              <div className="mb-2 space-y-2 text-gray-100">
                <label
                  htmlFor="aiInteractionMode"
                  className="block font-medium text-sm mb-1"
                >
                  Interaction Mode
                </label>
                <div className="flex items-center space-x-4 mt-2"></div>
                <label className="inline-flex items-center cursor-pointer">
                  <input
                    type="radio"
                    name="aiInteractionMode"
                    value="P2P"
                    className="peer hidden"
                    onClick={() => handleInteractionModeChange("P2P")}
                    defaultChecked
                  />
                  <div className=" w-4 h-4 rounded-full border-2 peer-checked:bg-orange-700 transition" />
                  <span className=" ml-2 mr-2 text-sm">Direct</span>
                </label>
                <label className=" inline-flex items-center cursor-pointer">
                  <input
                    type="radio"
                    name="aiInteractionMode"
                    value="CHAT_MODE"
                    className="peer hidden"
                    onClick={() => handleInteractionModeChange("CHAT_MODE")}
                  />
                  <div className="w-4 h-4 rounded-full border-2  peer-checked:bg-orange-700 transition" />
                  <span className="ml-2 text-sm">Chat</span>
                </label>
              </div>

              <div className="mb-2 mt-2 text-gray-100">
                <NestedDropdown />
                <input type="text" id="aiModel" hidden />
              </div>

              <div id="P2P" className="text-gray-100">
                <div className="relative ">
                  <label
                    htmlFor="toggleStateLabel"
                    id="toggleStateLabel"
                    className="block text-sm font-medium mb-1"
                  >
                    Stream
                  </label>
                  <button
                    id="toggleStreamButton"
                    className="relative cursor-pointer inline-flex h-6 w-16 items-center rounded-full bg-gray-300 transition-colors focus:outline-none focus:ring-2 focus:ring-[#ff5f1f]"
                    role="switch"
                    aria-checked="false"
                    onClick={(e) => toggleStream(e.target)}
                  >
                    <span className="absolute left-1 w-6 h-6 bg-white rounded-full shadow-md transition-transform transform"></span>
                  </button>
                </div>

                <div className="mb-2 mt-2">
                  <PromptDropdown />
                  <input type="text" id="promptId" hidden />
                </div>

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
              </div>

              <div id="CHAT_MODE" className="hidden">
                <div className="mb-2 mt-2 ">
                  <a
                    className="flex p-2 bg-orange-700 hover:bg-zinc-900 items-center justify-between text-sm font-medium transition-colors rounded-lg"
                    href="#"
                    onClick={createNewChat}
                  >
                    <span className="flex flex-row items-center justify-start">
                      <div slot="icon" className="w-4 text-neutral-white">
                        <svg
                          viewBox="0 0 24 24"
                          className="h-full w-6"
                          stroke="currentColor"
                        >
                          <path
                            fill="currentColor"
                            d="M5 3c-1.11 0-2 .89-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7h-2v7H5V5h7V3zm12.78 1a.7.7 0 0 0-.48.2l-1.22 1.21l2.5 2.5L19.8 6.7c.26-.26.26-.7 0-.95L18.25 4.2c-.13-.13-.3-.2-.47-.2m-2.41 2.12L8 13.5V16h2.5l7.37-7.38z"
                          ></path>
                        </svg>
                      </div>
                      <span className="ml-2 pl-2 justify-center text-sm">
                        Start New Chat
                      </span>
                    </span>
                  </a>
                </div>

                <div className="mb-2 ">
                  <label
                    htmlFor="chatHistory"
                    className="block font-semibold text-sm mb-1"
                  >
                    Chat History
                  </label>
                  <div className="w-full h-auto overflow-y-auto scrollbar">
                    <ul className="space-y-2">
                      {response?.AIStudioChats?.length > 0 ? (
                        response.AIStudioChats.map(
                          (val, id) =>
                            val.CurrentLog && (
                              <li
                                key={id}
                                className="p-2 text-gray-100 hover:bg-zinc-900 hover:text-white border-b border-zinc-500"
                              >
                                <p className="font-semibold text-xs break-words">
                                  <a
                                    href={`${globalContext?.CurrentUrl}?aiChatId=${val.ChatId}`}
                                  >
                                    {val.CurrentLog.Value.length > 35
                                      ? val.CurrentLog.Value.substring(0, 35) +
                                        "..."
                                      : val.CurrentLog.Value}
                                  </a>
                                </p>
                              </li>
                            )
                        )
                      ) : (
                        <li className="p-2 bg-zinc-900 shadow-sm rounded-md hover:bg-zinc-900">
                          <p className="font-semibold text-xs break-words">
                            No history available
                          </p>
                        </li>
                      )}
                    </ul>
                  </div>
                </div>
              </div>
            </div>

            <div
              id="rightbar"
              className="rightbar fixed top-20 right-0 w-64 h-full bg-[#1b1b1b] shadow-lg hidden rounded-lg text-gray-100"
            >
              <div className="p-4 bg-zinc-900 flex justify-between items-center rounded-md">
                <h2 className="text-lg font-bold">Sidebar</h2>
                <button
                  id="closeSidebar"
                  onClick={closeOptionsBar}
                  className="hover:text-gray-100 focus:outline-none cursor-pointer"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-6 w-6"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    strokeWidth="2"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M6 18L18 6M6 6l12 12"
                    />
                  </svg>
                </button>
              </div>

              <div className="p-4">
                <label
                  htmlFor="temperatureSelect"
                  className="block text-sm font-medium mb-1"
                >
                  Select Temperature
                </label>
                <div
                  id="temperatureSelect"
                  className="flex items-center space-x-4 mt-2"
                >
                  <input
                    id="temperatureSlider"
                    type="range"
                    min="0.0"
                    max="1.0"
                    step="0.1"
                    defaultValue="0.7"
                    className="w-full h-2 bg-orange-700 rounded-lg appearance-none cursor-pointer"
                    onChange={(e) => updateTemperatureValue(e.target.value)}
                  />
                  <input
                    id="temperatureValue"
                    type="number"
                    step="0.1"
                    min="0.0"
                    max="1.0"
                    defaultValue="0.7"
                    className="w-16 px-2 py-1 border border-gray-300 rounded-lg text-center focus:outline-none focus:ring focus:ring-orange-700"
                    onChange={(e) => updateSliderValue(e.target.value)}
                  />
                </div>

                <label
                  htmlFor="maxTokens"
                  className="block text-sm font-medium mt-4 mb-1"
                >
                  Select Max Tokens
                </label>
                <input
                  type="text"
                  id="maxTokens"
                  placeholder="Max Tokens..."
                  className="w-full px-3 py-2 mb-1 text-gray-100 shadow-md rounded-lg focus:outline-none focus:ring text-sm focus:ring-orange-700"
                />

                <label
                  htmlFor="topk"
                  className="block text-sm font-medium mt-4 mb-1"
                >
                  Select Top K
                </label>
                <input
                  type="text"
                  id="topk"
                  placeholder="Top K value between 0 - 1000"
                  className="w-full px-3 py-2 text-gray-100 mb-1 shadow-md rounded-lg focus:outline-none focus:ring text-sm focus:ring-orange-700"
                />

                <label
                  htmlFor="topP"
                  className="block text-sm font-medium mt-4 mb-1"
                >
                  Select Top P
                </label>
                <input
                  type="text"
                  id="topP"
                  placeholder="Top P value between 0.0 - 1.0"
                  className="w-full px-3 py-2 text-gray-100 mb-1 shadow-md rounded-lg focus:outline-none focus:ring text-sm focus:ring-orange-700"
                />
              </div>
            </div>

            <input
              type="text"
              id="currentChatId"
              hidden
            />

            <div id="currentChat" className="hidden">
              {aiStudioChat ? JSON.stringify(aiStudioChat) : ""}
            </div>
          </div>
        </div>

        {/* Context Popup Component */}
        <ContextModal 
          isOpen={isContextModalOpen} 
          onClose={() => setIsContextModalOpen(false)} 
        />
      </div>
    </>
  );
}