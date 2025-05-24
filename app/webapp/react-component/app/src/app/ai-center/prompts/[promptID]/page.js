"use client";
import { useState, useEffect } from "react";
import { use } from "react";
import Header from "@/app/components/platform/header";
import SectionHeaders from "@/app/components/section-headers";
import {
  promptArchiveApi,
  PromptsApi,
} from "@/app/utils/ai-studio-endpoint/prompts-api";

export default function PromptDetailsPage({ params }) {
  const unwrappedParams = use(params);
  const promptId = unwrappedParams?.promptID
    ? String(unwrappedParams.promptID).trim()
    : "";

  const [promptDetails, setPromptDetails] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [activeTab, setActiveTab] = useState("spec");

  useEffect(() => {
    const fetchPromptDetails = async () => {
      try {
        const response = await PromptsApi.getPromptsId(promptId);

        if (!response) {
          console.error("No response received from server");
          setError("No response received from server");
          setLoading(false);
          return;
        }
        if (
          response.output &&
          Array.isArray(response.output) &&
          response.output.length > 0
        ) {
          setPromptDetails(response.output[0]);
        } else if (response.output && !Array.isArray(response.output)) {
          setPromptDetails(response.output);
        } else {
          console.error(
            "Data does not contain expected output format:",
            response
          );
          setError("Unexpected data format received from server");
        }

        setLoading(false);
      } catch (err) {
        console.error("Error fetching prompt details:", err);
        setError(err.message);
        setLoading(false);
      }
    };

    if (promptId) {
      fetchPromptDetails();
    } else {
      setError("No prompt ID provided");
      setLoading(false);
    }
  }, [promptId]);

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
    alert("Copied to clipboard!");
  };

  const copyToClipboardUsingElement = (elementId) => {
    const element = document.getElementById(elementId);
    if (element) {
      navigator.clipboard.writeText(element.innerText);
      alert("Copied to clipboard!");
    }
  };

  const showTab = (tabId) => {
    setActiveTab(tabId);
  };

  if (loading) {
    return (
      <div className="bg-zinc-800 flex h-screen justify-center items-center">
        <div className="text-white text-xl">Loading prompt details...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-zinc-800 flex h-screen justify-center items-center">
        <div className="text-red-500 text-xl">Error: {error}</div>
      </div>
    );
  }

  if (!promptDetails) {
    return (
      <div className="bg-zinc-800 flex h-screen justify-center items-center">
        <div className="text-white text-xl">Prompt not found</div>
      </div>
    );
  }

  //  function to safely render tool schemas
  const renderToolSchemas = () => {
    if (
      !promptDetails.spec?.tools ||
      !Array.isArray(promptDetails.spec.tools) ||
      promptDetails.spec.tools.length === 0
    ) {
      return "No tools available";
    }

    return promptDetails.spec.tools.map((tool, index) => {
      if (!tool.schema) {
        if (tool.rawJsonParams) {
          try {
            const parsedParams = JSON.parse(tool.rawJsonParams);
            return (
              <pre key={index} className="s">
                {JSON.stringify(parsedParams, null, 2)}
              </pre>
            );
          } catch (e) {
            return (
              <pre key={index} className="s">
                {tool.rawJsonParams}
              </pre>
            );
          }
        }
        return (
          <pre key={index} className="s">
            No schema available for this tool
          </pre>
        );
      }

      // If schema is present
      try {
        const schemaObj = {
          name: tool.schema.name || "",
          description: tool.schema.description || "",
          parameters:
            typeof tool.schema.parameters === "string"
              ? JSON.parse(tool.schema.parameters)
              : tool.schema.parameters || {},
        };
        return (
          <pre key={index} className="s">
            {JSON.stringify(schemaObj, null, 2)}
          </pre>
        );
      } catch (e) {
        return (
          <pre key={index} className="s">
            Error parsing schema: {e.message}
          </pre>
        );
      }
    });
  };

  // Create a data object for SectionHeaders
  const apiServices = {
    prompt: {
      archive: promptArchiveApi.getPromptArchive,
      delete: promptArchiveApi.getPromptArchive,
    },
  };

  const headerResourceData = {
    id: promptDetails.promptId,
    name: promptDetails.name || "Unnamed Prompt",
    createdAt: promptDetails.resourceBase?.createdAt
      ? parseInt(promptDetails.resourceBase.createdAt) * 1000
      : null,
    createdBy: promptDetails.resourceBase?.createdBy,
    status: promptDetails.resourceBase?.status,
    resourceBase: promptDetails.resourceBase,
    resourceType: "prompt",
    // Create action params for update functionality
    createActionParams: promptDetails.createActionParams || {
      weblink: `./${promptDetails.promptId}/update`,
    },
    // Add YAML spec for download button
    yamlSpec: promptDetails.yamlSpec || JSON.stringify(promptDetails, null, 2),
  };

  return (
    <div className="bg-zinc-800 flex h-screen">
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="AI Model Prompt Details"
          hideBackListingLink={false}
          backListingLink="/ai-center/prompts"
        />

        <div className="flex-grow p-2 w-full">
          {promptDetails && promptDetails.promptId && (
            <SectionHeaders
              resourceId={promptDetails.promptId}
              resourceData={headerResourceData}
              resourceType="prompt"
              apiServices={apiServices}
            />
          )}
          <div className="overflow-x-auto scrollbar text-gray-100 bg-zinc-800 rounded-lg p-8 shadow-md">
            <div className="flex border-b border-zinc-500">
              <button
                onClick={() => showTab("spec")}
                className={`px-4 py-2 font-semibold  focus:outline-none cursor-pointer ${
                  activeTab === "spec"
                    ? "bg-[oklch(0.205_0_0)] text-white rounded-t-[10px]"
                    : ""
                }`}
              >
                Spec
              </button>
              <button
                onClick={() => showTab("dryrunset")}
                className={`px-4 py-2 font-semibold focus:outline-none cursor-pointer ${
                  activeTab === "dryrunset"
                    ? "bg-[oklch(0.205_0_0)] text-white rounded-t-[10px]"
                    : ""
                }`}
              >
                Dry Run
              </button>
            </div>

            {/* Tab Content */}
            <div
              id="spec"
              className={`mt-2 bg-[#1b1b1b] p-4 ${
                activeTab !== "spec" ? "hidden" : ""
              }`}
            >
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block  ">
                    Name:
                  </p>
                  <p className="s p-2">{promptDetails.name || "N/A"}</p>
                </div>
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block  ">
                    Id:
                  </p>
                  <p
                    className="s p-2 cursor-pointer"
                    onClick={() => copyToClipboard(promptDetails.promptId)}
                  >
                    {promptDetails.promptId || "N/A"}
                  </p>
                </div>
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block  ">
                    Prompt Type:
                  </p>
                  <p className="s p-2">
                    {promptDetails.promptTypes &&
                      Array.isArray(promptDetails.promptTypes) &&
                      promptDetails.promptTypes.map((type, index) => (
                        <span
                          key={index}
                          className="px-3 py-1 text-sm font-medium rounded-full text-purple-800 bg-purple-100 mr-2"
                        >
                          {type}
                        </span>
                      ))}
                  </p>
                </div>
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block  ">
                    Preferred Models:
                  </p>
                  <div className="grid grid-cols-2 grid-rows-2 md:grid-cols-2 md:grid-rows-2 lg:grid-cols-4 lg:grid-rows-1 gap-2 p-2 text-center">
                    {promptDetails.preferredModels &&
                      Array.isArray(promptDetails.preferredModels) &&
                      promptDetails.preferredModels.map((model, index) => (
                        <span
                          key={index}
                          className="px-3 py-1 text-sm font-medium rounded-full text-gray-800 bg-gray-100 mr-2 mb-1"
                        >
                          {model}
                        </span>
                      ))}
                  </div>
                </div>
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block  ">
                    Labels:
                  </p>
                  <div className="grid grid-cols-1 md:grid-cols-2 md:grid-rows-2 lg:grid-cols-4 lg:grid-rows-1 gap-2 p-2 text-center">
                    {promptDetails.labels &&
                      Array.isArray(promptDetails.labels) &&
                      promptDetails.labels.map((label, index) => (
                        <span
                          key={index}
                          className="px-3 py-1 text-sm font-medium rounded-full text-yellow-800 bg-yellow-100 mr-2"
                        >
                          {label}
                        </span>
                      ))}
                  </div>
                </div>
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block  ">
                    Owner:
                  </p>
                  <p className="s p-2">{promptDetails.promptOwner || "N/A"}</p>
                </div>
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block  ">
                    Domain:
                  </p>
                  <p className="s p-2">
                    {promptDetails.resourceBase?.domain || "N/A"}
                  </p>
                </div>
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block  ">
                    User Message:
                  </p>
                  <p className="s p-2">
                    {promptDetails.spec?.userMessage || "N/A"}
                  </p>
                </div>
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block  ">
                    System Message:
                  </p>
                  <p className="s p-2">
                    {promptDetails.spec?.systemMessage || "N/A"}
                  </p>
                </div>
                {promptDetails.spec?.sample && (
                  <div>
                    <p className="text-base font-extralight text-[#f4d1c2] block  ">
                      Sample Input:
                    </p>
                    <p className="s p-2">
                      {promptDetails.spec.sample.inputText || "N/A"}
                    </p>
                  </div>
                )}
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block  ">
                    Variables:
                  </p>
                  <p
                    className="flex text-gray-200 s p-2"
                    style={{ maxWidth: "200px", overflowWrap: " " }}
                  >
                    {promptDetails.spec?.variables &&
                    Array.isArray(promptDetails.spec.variables) &&
                    promptDetails.spec.variables.length > 0
                      ? promptDetails.spec.variables.map((variable, index) => (
                          <span
                            key={index}
                            className="px-3 py-1 text-sm font-medium bg-blue-600 rounded-full mr-2"
                          >
                            {variable}
                          </span>
                        ))
                      : "No Variables for this prompt"}
                  </p>
                </div>
              </div>

              <h3 className="text-xl mb-4 text-[1.25rem] font-bold text-[#f4d1c2] underline">
                Tags
              </h3>
              <div className="bg-zinc-800 p-4 rounded-lg shadow-md">
                <table className="min-w-full divide-y divide-zinc-500 border border-zinc-500">
                  <thead className="bg-zinc-900">
                    <tr>
                      <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                        Field
                      </th>
                      <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                        Tag
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-zinc-800 divide-y divide-zinc-500">
                    <tr className="cursor-pointer hover:bg-zinc-600">
                      <td className="px-3 py-3 whitespace-nowrap">Input Tag</td>
                      <td className="px-3 py-3 whitespace-nowrap">
                        {promptDetails.spec?.inputTag || "N/A"}
                      </td>
                    </tr>
                    <tr className="cursor-pointer hover:bg-zinc-600">
                      <td className="px-3 py-3 whitespace-nowrap">
                        Result Tag
                      </td>
                      <td className="px-3 py-3 whitespace-nowrap">
                        {promptDetails.spec?.outputTag || "N/A"}
                      </td>
                    </tr>
                    <tr className="cursor-pointer hover:bg-zinc-600">
                      <td className="px-3 py-3 whitespace-nowrap">
                        Context Tag
                      </td>
                      <td className="px-3 py-3 whitespace-nowrap">
                        {promptDetails.spec?.contextTag || "N/A"}
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            <div
              id="dryrunset"
              className={`mt-2 bg-[#1b1b1b] p-4 ${
                activeTab !== "dryrunset" ? "hidden" : ""
              }`}
            >
              <div className="w-full shadow-lg rounded-lg p-2">
                <div className="flex justify-between items-center border-b border-orange-700 pb-2 mb-4">
                  <h1 className="text-lg font-bold">Rendered Template</h1>
                  <button
                    id="copyJsonBtn"
                    onClick={() => copyToClipboardUsingElement("jsonViewer")}
                    className="hover:text-blue-500 lg:flex items-center"
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-5 w-5 mr-1"
                      viewBox="0 0 24 24"
                      fill="currentColor"
                    >
                      <path d="M8 2C6.895 2 6 2.895 6 4v12c0 1.105.895 2 2 2h8c1.105 0 2-.895 2-2V4c0-1.105-.895-2-2-2H8zM8 4h8v12H8V4zm-4 4H3v8a2 2 0 002 2h6v-2H5V8zm14 0h1v8a2 2 0 01-2 2h-6v2h6a4 4 0 004-4V8z" />
                    </svg>
                    Copy Template
                  </button>
                </div>
                <div
                  id="jsonViewer"
                  className="mt-2 border border-zinc-500 rounded-lg p-4 max-w-full text-sm font-mono overflow-x-auto scrollbar overflow-y-auto scrollbar"
                  style={{ maxHeight: "400px", whiteSpace: "pre-wrap" }}
                >
                  <pre className="s">
                    {promptDetails.template || "No template available"}
                  </pre>
                </div>
              </div>

              {promptDetails.spec?.tools &&
                Array.isArray(promptDetails.spec.tools) &&
                promptDetails.spec.tools.length > 0 && (
                  <div className="w-full shadow-lg rounded-lg p-2 mt-6">
                    <div className="flex justify-between items-center pb-4 mb-4">
                      <h1 className="text-lg font-bold">Tool Schemas</h1>
                      <button
                        id="copyJsonBtn"
                        onClick={() =>
                          copyToClipboardUsingElement("toolSchema")
                        }
                        className="hover:text-blue-500 lg:flex items-center"
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          className="h-5 w-5 mr-1"
                          viewBox="0 0 24 24"
                          fill="currentColor"
                        >
                          <path d="M8 2C6.895 2 6 2.895 6 4v12c0 1.105.895 2 2 2h8c1.105 0 2-.895 2-2V4c0-1.105-.895-2-2-2H8zM8 4h8v12H8V4zm-4 4H3v8a2 2 0 002 2h6v-2H5V8zm14 0h1v8a2 2 0 01-2 2h-6v2h6a4 4 0 004-4V8z" />
                        </svg>
                        Copy Schema
                      </button>
                    </div>
                    <div
                      id="toolSchema"
                      className="mt-2 border border-zinc-500 max-w-full rounded-lg p-4 text-sm font-mono overflow-x-auto scrollbar overflow-y-auto scrollbar"
                      style={{ maxHeight: "400px", whiteSpace: "pre-wrap" }}
                    >
                      {renderToolSchemas()}
                    </div>
                  </div>
                )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
