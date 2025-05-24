"use client";
import { useState, useEffect, use } from "react";
import Header from "@/app/components/platform/header";
import SectionHeaders from "@/app/components/section-headers";
import {
  GuardrailApi,
  GuardrailArchiveApi,
} from "@/app/utils/ai-studio-endpoint/guardrails-api";
import { strTitle } from "@/app/components/JS/common";

export default function GuardrailDetailsPage({ params }) {
  const resolvedParams = use(params);
  const guardrailId = resolvedParams.guardrailId;

  const [guardrailDetails, setGuardrailDetails] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [activeTab, setActiveTab] = useState("spec");

  useEffect(() => {
    const fetchGuardrailDetails = async () => {
      try {
        const response = await GuardrailApi.getGuardrailId(guardrailId);

        if (!response.output) {
          console.error("Data does not contain output property:", response);
          setError("Unexpected data format received from server");
          setLoading(false);
          return;
        }

        // Find the guardrail with matching ID in the output array
        const foundGuardrail = Array.isArray(response.output)
          ? response.output.find((item) => item.guardrailId === guardrailId)
          : null;

        if (foundGuardrail) {
          setGuardrailDetails(foundGuardrail);
        } else {
          setError(`Guardrail not found for the given ID: ${guardrailId}`);
        }

        setLoading(false);
      } catch (err) {
        console.error("Error fetching guardrail details:", err);
        setError(err.message);
        setLoading(false);
      }
    };

    if (guardrailId) {
      fetchGuardrailDetails();
    } else {
      setError("No guardrail ID provided");
      setLoading(false);
    }
  }, [guardrailId]);

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);

    // Create a simple toast notification
    const toast = document.createElement("div");
    toast.textContent = "Copied to clipboard";
    toast.style.position = "fixed";
    toast.style.bottom = "20px";
    toast.style.left = "50%";
    toast.style.transform = "translateX(-50%)";
    toast.style.backgroundColor = "rgba(0, 0, 0, 0.7)";
    toast.style.color = "white";
    toast.style.padding = "10px 20px";
    toast.style.borderRadius = "5px";
    toast.style.zIndex = "1000";

    document.body.appendChild(toast);

    setTimeout(() => {
      document.body.removeChild(toast);
    }, 1000);
  };

  const stringCheck = (str) => {
    return str ? str : "N/A";
  };

  // function to limit words in text
  const limitWords = (text, limit) => {
    if (!text) return "N/A";
    const words = text.split(" ");
    if (words.length <= limit) return text;
    return words.slice(0, limit).join(" ") + "...";
  };

  if (loading) {
    return (
      <div className="bg-zinc-800 flex h-screen justify-center items-center">
        <div className="text-white text-xl">Loading guardrail details...</div>
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

  if (!guardrailDetails) {
    return (
      <div className="bg-zinc-800 flex h-screen justify-center items-center">
        <div className="text-white text-xl">Guardrail not found</div>
      </div>
    );
  }

  // Prepare data for SectionHeaders component with action-related fields

  const apiServices = {
    guardrail: {
      archive: GuardrailArchiveApi.getGuardrailArchive,
      delete: GuardrailArchiveApi.getGuardrailArchive,
    },
  };

  const headerResourceData = {
    id: guardrailDetails.guardrailId,
    name: guardrailDetails.name || guardrailDetails.displayName,
    createdAt: guardrailDetails.resourceBase?.createdAt
      ? parseInt(guardrailDetails.resourceBase.createdAt) * 1000
      : null,
    createdBy: guardrailDetails.resourceBase?.createdBy,
    status: guardrailDetails.resourceBase?.status,
    resourceBase: guardrailDetails.resourceBase,
    resourceType: "guardrail",
    // Sample createActionParams if needed
    createActionParams: guardrailDetails.createActionParams || {
      weblink: `/ai-center/guardrails/${guardrailDetails.guardrailId}/update`,
    },

    yamlSpec:
      guardrailDetails.yamlSpec || JSON.stringify(guardrailDetails, null, 2),
  };

  return (
    <div className="bg-zinc-800 flex h-screen">
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="AI Guardrail Details"
          hideBackListingLink={false}
          backListingLink="./"
        />

        <div className="flex-grow p-2 w-full text-gray-100">
          {/* Using the reusable SectionHeaders component with action data */}
          <SectionHeaders
            resourceId={guardrailDetails.guardrailId}
            resourceData={headerResourceData}
            resourceType="guardrail"
            apiServices={apiServices}
          />

          <div className="overflow-x-auto scrollbar text-gray-100 bg-zinc-800 rounded-lg p-8 shadow-md">
            {/* Tabs */}
            {/* <div className="flex border-b border-zinc-500">
              <button
                onClick={() => setActiveTab("spec")}
                className={`px-4 py-2 font-semibold focus:outline-none ${
                  activeTab === "spec"
                    ? "bg-[oklch(0.205_0_0)] text-white rounded-t-[10px]"
                    : "text-gray-400"
                }`}
              >
                Spec
              </button>
              <button
                onClick={() => setActiveTab("schema")}
                className={`px-4 py-2 font-semibold focus:outline-none ${
                  activeTab === "schema"
                    ? "bg-[oklch(0.205_0_0)] text-white rounded-t-[10px]"
                    : "text-gray-400"
                }`}
              >
                Schema
              </button>
            </div> */}

            {/* Tab Content */}
            <div
              id="spec"
              className={`tab-content mt-2 bg-[#1b1b1b] p-4 rounded-lg shadow-md ${
                activeTab !== "spec" ? "hidden" : ""
              }`}
            >
              {/* <h3 className="text-xl mb-4 font-bold text-[#f4d1c2] underline">
                Basic Information
              </h3> */}
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Name:
                  </p>
                  <p className="break-words p-2">
                    {stringCheck(
                      guardrailDetails.name || guardrailDetails.displayName
                    )}
                  </p>
                </div>
                <div className=" lg:flex items-center">
                  <p
                    className="text-base font-extralight text-[#f4d1c2] block cursor-pointer"
                    onClick={() =>
                      copyToClipboard(guardrailDetails.guardrailId)
                    }
                  >
                    Id:
                  </p>
                  <p className="break-words p-2">
                    {stringCheck(guardrailDetails.guardrailId)}
                  </p>
                </div>
                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Labels:
                  </p>
                  <p className="break-words p-2">
                    {guardrailDetails.resourceBase?.labels &&
                    guardrailDetails.resourceBase.labels.length > 0
                      ? guardrailDetails.resourceBase.labels.map(
                          (label, index) => (
                            <span
                              key={index}
                              className="px-3 py-1 text-sm font-medium rounded-full text-yellow-800 bg-yellow-100 mr-2"
                            >
                              {label}
                            </span>
                          )
                        )
                      : "N/A"}
                  </p>
                </div>
                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Owners:
                  </p>
                  <div className="break-words p-2">
                    <ul className="list-disc ml-5 text-gray-200">
                      {guardrailDetails.resourceBase?.owners &&
                      guardrailDetails.resourceBase.owners.length > 0
                        ? guardrailDetails.resourceBase.owners.map(
                            (owner, index) => <li key={index}>{owner}</li>
                          )
                        : "N/A"}
                    </ul>
                  </div>
                </div>
                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Description:
                  </p>
                  <p
                    className="text-gray-200 break-words p-2 cursor-pointer"
                    onClick={() =>
                      copyToClipboard(guardrailDetails.description)
                    }
                  >
                    {limitWords(guardrailDetails.description, 30)}
                  </p>
                </div>
                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Failure Message:
                  </p>
                  <p
                    className="text-gray-200 break-words p-2 cursor-pointer"
                    onClick={() =>
                      copyToClipboard(guardrailDetails.failureMessage)
                    }
                  >
                    {limitWords(guardrailDetails.failureMessage, 30)}
                  </p>
                </div>
                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Domain:
                  </p>
                  <p className="break-words p-2">
                    {stringCheck(guardrailDetails.resourceBase?.domain)}
                  </p>
                </div>
                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Scan Mode:
                  </p>
                  <p className="break-words p-2">
                    {stringCheck(guardrailDetails.scanMode)}
                  </p>
                </div>
              </div>

              {/* Guard Model Section */}
              {guardrailDetails.guardModel && (
                <div className="mt-4">
                  <h3 className="text-xl mb-4 font-bold text-[#f4d1c2] underline">
                    Guard Model
                  </h3>
                  <div className="w-full bg-zinc-800 p-4 rounded-lg shadow-md text-sm">
                    <table className="min-w-full divide-y divide-zinc-500 text-gray-100 border-2 border-zinc-500">
                      <thead className="bg-zinc-900">
                        <tr>
                          <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                            Model Node Id
                          </th>
                          <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                            Model
                          </th>
                        </tr>
                      </thead>
                      <tbody className="bg-zinc-800 divide-y divide-zinc-500">
                        <tr className="cursor-pointer hover:bg-zinc-600">
                          <td className="px-3 py-3 whitespace-nowrap">
                            {stringCheck(
                              guardrailDetails.guardModel.modelNodeId
                            )}
                          </td>
                          <td className="px-3 py-3 whitespace-nowrap">
                            {stringCheck(guardrailDetails.guardModel.modelId)}
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              )}

              {/* Content Rules Section */}
              <div className="mt-4">
                <h3 className="text-xl mb-4 font-bold text-[#f4d1c2] underline">
                  Content Rules
                </h3>
                <div className="w-full bg-zinc-800 p-4 rounded-lg shadow-md text-sm">
                  <table className="min-w-full divide-y divide-zinc-500 text-gray-100 border-2 border-zinc-500">
                    <thead className="bg-zinc-900">
                      <tr>
                        <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                          Hate Speech
                        </th>
                        <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                          Insults
                        </th>
                        <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                          Sexual
                        </th>
                        <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                          Threats
                        </th>
                        <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                          Misconduct
                        </th>
                      </tr>
                    </thead>
                    <tbody className="bg-zinc-800 divide-y divide-zinc-500">
                      <tr className="cursor-pointer hover:bg-zinc-600">
                        <td className="px-3 py-3 whitespace-nowrap">
                          {strTitle(guardrailDetails.contents?.hateSpeech)}
                        </td>
                        <td className="px-3 py-3 whitespace-nowrap">
                          {strTitle(guardrailDetails.contents?.insults)}
                        </td>
                        <td className="px-3 py-3 whitespace-nowrap">
                          {strTitle(guardrailDetails.contents?.sexual)}
                        </td>
                        <td className="px-3 py-3 whitespace-nowrap">
                          {strTitle(guardrailDetails.contents?.threats)}
                        </td>
                        <td className="px-3 py-3 whitespace-nowrap">
                          {strTitle(guardrailDetails.contents?.misconduct)}
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Topics Rules Section */}
              <div className="mt-4">
                <h3 className="text-xl mb-4 font-bold text-[#f4d1c2] underline">
                  Topics Rules
                </h3>
                <div className="w-full bg-zinc-800 p-4 rounded-lg shadow-md text-sm">
                  <table className="min-w-full divide-y divide-zinc-500 text-gray-100 border-2 border-zinc-500">
                    <thead className="bg-zinc-900">
                      <tr>
                        <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                          Topic
                        </th>
                        <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                          Description
                        </th>
                        <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                          Samples
                        </th>
                      </tr>
                    </thead>
                    <tbody className="bg-zinc-800 divide-y divide-zinc-500">
                      {guardrailDetails.topics &&
                      guardrailDetails.topics.length > 0 ? (
                        guardrailDetails.topics.map((topic, index) => (
                          <tr
                            key={index}
                            className="cursor-pointer hover:bg-zinc-600"
                          >
                            <td className="px-3 py-3 break-words">
                              {topic.topic}
                            </td>
                            <td className="px-3 py-3 break-words">
                              {stringCheck(topic.description)}
                            </td>
                            <td className="px-3 py-3 break-words">
                              <ul className="list-disc ml-5 text-gray-100">
                                {topic.samples &&
                                  topic.samples.map((sample, idx) => (
                                    <li key={idx}>{sample}</li>
                                  ))}
                              </ul>
                            </td>
                          </tr>
                        ))
                      ) : (
                        <tr className="cursor-pointer hover:bg-zinc-600">
                          <td
                            className="px-3 py-3 text-gray-100 whitespace-nowrap"
                            colSpan="3"
                          >
                            No Topics guard rule
                          </td>
                        </tr>
                      )}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Word Rules Section */}
              <div className="mt-4">
                <h3 className="text-xl mb-4 font-bold text-[#f4d1c2] underline">
                  Word Rules
                </h3>
                <div className="w-full bg-zinc-800 p-4 rounded-lg shadow-md text-sm">
                  <table className="min-w-full divide-y divide-zinc-500 text-gray-100 border-2 border-zinc-500">
                    <thead className="bg-zinc-900">
                      <tr>
                        <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                          Words
                        </th>
                        <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                          File Path
                        </th>
                      </tr>
                    </thead>
                    <tbody className="bg-zinc-800 divide-y divide-zinc-500">
                      {guardrailDetails.words &&
                      guardrailDetails.words.length > 0 ? (
                        guardrailDetails.words.map((wordRule, index) => (
                          <tr
                            key={index}
                            className="cursor-pointer hover:bg-zinc-600"
                          >
                            <td className="px-3 py-3 break-words">
                              {wordRule.words
                                ? wordRule.words.join(" | ")
                                : "N/A"}
                            </td>
                            <td className="px-3 py-3 break-words">
                              {stringCheck(wordRule.fileLocation)}
                            </td>
                          </tr>
                        ))
                      ) : (
                        <tr className="cursor-pointer hover:bg-zinc-600">
                          <td
                            className="px-3 py-3 text-gray-100 whitespace-nowrap"
                            colSpan="2"
                          >
                            No Words guard rule
                          </td>
                        </tr>
                      )}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Data Sensitivity Rules Section */}
              <div className="mt-4">
                <h3 className="text-xl mb-4 font-bold text-[#f4d1c2] underline">
                  Data Sensitivity Rules
                </h3>
                <div className="w-full bg-zinc-800 p-4 rounded-lg shadow-md text-sm">
                  <table className="min-w-full divide-y divide-zinc-500 text-gray-100 border-2 border-zinc-500">
                    <thead className="bg-zinc-900">
                      <tr>
                        <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                          PII Type
                        </th>
                        <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                          Regex
                        </th>
                        <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                          Action
                        </th>
                      </tr>
                    </thead>
                    <tbody className="bg-zinc-800 divide-y divide-zinc-500">
                      {guardrailDetails.sensitiveDataset &&
                      guardrailDetails.sensitiveDataset.length > 0 ? (
                        guardrailDetails.sensitiveDataset.map((rule, index) => (
                          <tr
                            key={index}
                            className="cursor-pointer hover:bg-zinc-600"
                          >
                            <td className="px-3 py-3 break-words">
                              {stringCheck(rule.piiType)}
                            </td>
                            <td className="px-3 py-3 break-words">
                              {stringCheck(rule.regex)}
                            </td>
                            <td className="px-3 py-3 break-words">
                              {stringCheck(rule.action)}
                            </td>
                          </tr>
                        ))
                      ) : (
                        <tr className="cursor-pointer hover:bg-zinc-600">
                          <td
                            className="px-3 py-3 text-gray-100 whitespace-nowrap"
                            colSpan="3"
                          >
                            No SensitiveDataset guard rule
                          </td>
                        </tr>
                      )}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>

            {/* Schema Tab Content */}
            <div
              id="schema"
              className={`tab-content mt-2 bg-[#1b1b1b] p-4 rounded-lg shadow-md ${
                activeTab !== "schema" ? "hidden" : ""
              }`}
            >
              <div className="bg-zinc-900 p-4 rounded-md overflow-x-auto">
                <pre className="text-gray-100 whitespace-pre-wrap">
                  {guardrailDetails.schema
                    ? guardrailDetails.schema
                    : "No schema available"}
                </pre>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
