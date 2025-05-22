"use client";
import React, { useState, useEffect } from "react";
import { use } from "react";
import Header from "@/app/components/platform/header";
import { modelsRegistryApi } from "@/app/utils/ai-studio-endpoint/models-registry-api";
import {
  GuardrailApi,
  GuardrailUpdateFormApi,
} from "@/app/utils/ai-studio-endpoint/guardrails-api";
import { toast } from "react-toastify";
import LoadingOverlay from "@/app/components/loading/loading";
import RemoveButton from "@/app/components/buttons/removeButton";
import AddButton from "@/app/components/buttons/addButton";
import { enumsApi } from "@/app/utils/developers-endpoint/enums";
import { strTitle } from "@/app/components/JS/common";
import { useRouter } from "next/navigation";
import ToastContainerMessage from "@/app/components/notification/customToast";

export default function UpdateGuardrail({ params }) {
  const router = useRouter();
  const unwrappedParams = use(params);
  const guardrail_id = unwrappedParams?.guardrailId
    ? String(unwrappedParams.guardrailId).trim()
    : "";

  // Form state
  const [guardrail, setGuardrail] = useState(null);
  const [view, setView] = useState("form"); // 'form' or 'yaml'

  // Dynamic form sections state
  const [topicCount, setTopicCount] = useState(0);
  const [wordEntryCount, setWordEntryCount] = useState(1);
  const [sensitiveEntryCount, setSensitiveEntryCount] = useState(1);

  // Models data state
  const [aiModelNodes, setAiModelNodes] = useState([]);
  const [modelMap, setModelMap] = useState({});
  const [availableModels, setAvailableModels] = useState([]);

  // UI state
  const [loading, setLoading] = useState(true);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [enums, setEnums] = useState({
    AIGuardrailScanMode: [],
    ResourceScope: [],
    GuardRailLevels: [],
  });

  // Fetch enums data
  useEffect(() => {
    const fetchEnumsData = async () => {
      try {
        setIsLoading(true);
        const response = await enumsApi.getEnums();
        const enumResponses = response.enumResponse || [];
        const newEnums = {
          AIGuardrailScanMode: [],
          ResourceScope: [],
          GuardRailLevels: [],
        };

        enumResponses.forEach((enumData) => {
          if (enumData.name === "AIGuardrailScanMode") {
            newEnums.AIGuardrailScanMode = enumData.value || [];
          } else if (enumData.name === "ResourceScope") {
            newEnums.ResourceScope = enumData.value || [];
          } else if (enumData.name === "GuardRailLevels") {
            newEnums.GuardRailLevels = enumData.value || [];
          }
        });

        setEnums(newEnums);
        console.log("Resource scopes loaded:", newEnums.ResourceScope);
        console.log("Service providers loaded:", newEnums.AIGuardrailScanMode);
        console.log("Guard rail levels loaded:", newEnums.GuardRailLevels);
      } catch (error) {
        console.error("Failed to fetch enum data:", error);
        setError(error.message);
        toast.error("Failed to fetch configuration data");
      } finally {
        setIsLoading(false);
      }
    };

    fetchEnumsData();
  }, []);

  // Fetch guardrail data and model nodes when guardrailId is available
  useEffect(() => {
    if (!guardrail_id) {
      console.error("No guardrail ID provided");
      toast.error("No guardrail ID provided for update");
      return;
    }

    const fetchGuardrailsData = async () => {
      try {
        setLoading(true);
        console.log("Fetching guardrail data:", guardrail_id);
        const response = await GuardrailApi.getGuardrailId(guardrail_id);
        const guardrail = response.output[0];

        if (guardrail) {
          // Set counts based on existing data
          setTopicCount(guardrail.topics ? guardrail.topics.length : 0);
          setWordEntryCount(
            guardrail.words ? Math.max(guardrail.words.length, 1) : 1
          );
          setSensitiveEntryCount(
            guardrail.sensitiveDataset
              ? Math.max(guardrail.sensitiveDataset.length, 1)
              : 1
          );
          const modelMap = await fetchModelNodes();
          setGuardrail(guardrail);
          if (
            guardrail.guardModel?.modelNodeId &&
            modelMap[guardrail.guardModel.modelNodeId]
          ) {
            const selectedNode = modelMap[guardrail.guardModel.modelNodeId];
            if (selectedNode?.attributes?.generativeModels) {
              setAvailableModels(selectedNode.attributes.generativeModels);
            }
          }
        } else {
          setError("No guardrails found");
        }
      } catch (error) {
        console.error("Error fetching guardrail data:", error);
        setError(`Failed to load guardrail: ${error.message}`);
      } finally {
        setLoading(false);
      }
    };

    fetchGuardrailsData();
  }, [guardrail_id]);

  // Fetch model nodes data
  const fetchModelNodes = async (preSelectedNodeId = null) => {
    try {
      setLoading(true);
      const response = await modelsRegistryApi.getModelsRegistry();
      const nodes = response.output?.aiModelNodes || [];
      setAiModelNodes(nodes);
      const map = {};
      nodes.forEach((item) => {
        map[item.modelNodeId] = item;
      });
      setModelMap(map);

      // If we have a pre-selected node ID, populate available models
      if (preSelectedNodeId && map[preSelectedNodeId]) {
        const node = map[preSelectedNodeId];
        if (node?.attributes?.generativeModels) {
          setAvailableModels(node.attributes.generativeModels);
        }
      }

      return map;
    } catch (error) {
      console.error("Failed to fetch model nodes:", error);
      setError(`Failed to load model nodes: ${error.message}`);
      return {};
    } finally {
      setLoading(false);
    }
  };

  // Populate available models based on selected node
  const populateModelDropdown = (nodeId) => {
    if (!nodeId || !modelMap[nodeId]) {
      setAvailableModels([]);
      return;
    }

    const node = modelMap[nodeId];
    if (node?.attributes?.generativeModels) {
      setAvailableModels(node.attributes.generativeModels);
    } else {
      setAvailableModels([]);
    }
  };

  const addTopic = () => {
    setTopicCount((prevCount) => prevCount + 1);
  };

  const removeTopic = (index) => {
    if (topicCount <= 0) return;
    setTopicCount((prevCount) => prevCount - 1);
  };

  const addWordEntry = () => {
    setWordEntryCount((prevCount) => prevCount + 1);
  };

  const removeWordEntry = (index) => {
    if (wordEntryCount <= 1) return;
    setWordEntryCount((prevCount) => prevCount - 1);
  };

  const addSensitiveEntry = () => {
    setSensitiveEntryCount((prevCount) => prevCount + 1);
  };

  const removeSensitiveEntry = (index) => {
    if (sensitiveEntryCount <= 1) return;
    setSensitiveEntryCount((prevCount) => prevCount - 1);
  };

  const submitUpdateForm = async (e) => {
    e.preventDefault();
    setIsLoading(true);

    const dataObj = {
      spec: {
        displayName: "",
        name: "",
        minConfidence: "",
        maxConfidence: "",
        description: "",
        failureMessage: "",
        scanMode: "INVALID_SSP",
        resourceBase: {
          scope: "DOMAIN_SCOPE",
        },
        contents: {
          hateSpeech: "",
          insults: "",
          sexual: "",
          threats: "",
          misconduct: "",
        },
        topics: [],
        words: [],
        sensitiveDataset: [],
        guardModel: {
          modelNodeId: "",
          modelId: "",
        },
        guardrailId: guardrail_id,
      },
    };

    const formData = new FormData(e.target);

    // Process standard fields
    dataObj.spec.displayName = formData.get("spec.displayName") || "";
    dataObj.spec.name = formData.get("spec.name") || "";
    dataObj.spec.minConfidence = formData.get("spec.minConfidence")
      ? Number(formData.get("spec.minConfidence"))
      : "";
    dataObj.spec.maxConfidence = formData.get("spec.maxConfidence")
      ? Number(formData.get("spec.maxConfidence"))
      : "";
    dataObj.spec.description = formData.get("spec.description") || "";
    dataObj.spec.failureMessage = formData.get("spec.failureMessage") || "";
    dataObj.spec.scanMode = formData.get("spec.scanMode") || "INVALID_SSP";
    dataObj.spec.resourceBase.scope =
      formData.get("spec.resourceBase.scope") || "DOMAIN_SCOPE";

    // Process contents
    dataObj.spec.contents.hateSpeech =
      formData.get("spec.contents.hateSpeech") || "";
    dataObj.spec.contents.insults = formData.get("spec.contents.insults") || "";
    dataObj.spec.contents.sexual = formData.get("spec.contents.sexual") || "";
    dataObj.spec.contents.threats = formData.get("spec.contents.threats") || "";
    dataObj.spec.contents.misconduct =
      formData.get("spec.contents.misconduct") || "";

    // Process topics (array)
    for (let i = 0; i < topicCount; i++) {
      const topic = formData.get(`spec.topics[${i}].topic`) || "";
      const description = formData.get(`spec.topics[${i}].description`) || "";
      const samplesStr = formData.get(`spec.topics[${i}].samples`) || "";
      const samples = samplesStr
        ? samplesStr.split(",").map((s) => s.trim())
        : [];

      dataObj.spec.topics.push({
        topic,
        description,
        samples,
      });
    }

    // Process words (array)
    for (let i = 0; i < wordEntryCount; i++) {
      const wordsStr = formData.get(`spec.words[${i}].words`) || "";
      const words = wordsStr ? wordsStr.split(",").map((w) => w.trim()) : [];
      const fileLocation = formData.get(`spec.words[${i}].fileLocation`) || "";

      dataObj.spec.words.push({
        words,
        fileLocation,
      });
    }

    // Process sensitiveDataset (array)
    for (let i = 0; i < sensitiveEntryCount; i++) {
      const piiType = formData.get(`spec.sensitiveDataset[${i}].piiType`) || "";
      const action = formData.get(`spec.sensitiveDataset[${i}].action`) || "";
      const regex = formData.get(`spec.sensitiveDataset[${i}].regex`) || "";

      dataObj.spec.sensitiveDataset.push({
        piiType,
        action,
        regex,
      });
    }

    // Process guard model
    dataObj.spec.guardModel.modelNodeId =
      formData.get("spec.guardModel.modelNodeId") || "";
    dataObj.spec.guardModel.modelId =
      formData.get("spec.guardModel.modelId") || "";

    console.log("Submitting form data:", dataObj);

    const payload = { spec: dataObj.spec };

    try {
      const output = await GuardrailUpdateFormApi.getGuardrailUpdateForm(
        payload
      );

      console.log("Resource updated:", output);

      toast.success("Guardrail updated successfully!");

      setTimeout(() => {
        router.push(`/ai-center/guardrails/${guardrail_id}`);
      }, 1000);
    } catch (error) {
      console.error("Error updating guardrail:", error);
      toast.error(`Failed to update guardrail: ${error.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  // Handle error case
  if (!guardrail) {
    return (
      <div className="bg-zinc-800 flex h-screen">
        <div className="overflow-y-auto scrollbar h-screen w-full">
          <Header />
          <div className="flex-grow p-4 overflow-y-auto w-full">
            <div className="flex justify-center items-center h-64">
              <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-orange-700"></div>
              <p className="ml-3 text-gray-100">Loading guardrail data...</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-zinc-800 flex h-screen">
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="Update Guardrail"
          hideBackListingLink={false}
          backListingLink="./"
        />
        <ToastContainerMessage />
        <LoadingOverlay isLoading={isLoading} />

        <div className="flex-grow p-4 overflow-y-auto w-full">
          <section className="space-y-2">
            <div className="max-w-6xl mx-auto bg-[#1b1b1b] shadow rounded-lg p-2">
              <div className="text-gray-100 mb-2 flex justify-center">
                {/* <button
                  onClick={() => setView("yaml")}
                  className={`whitespace-nowrap border-b-2 py-2 px-2 text-md font-medium focus:outline-none ${
                    view === "yaml"
                      ? "border-orange-700 text-orange-700 font-semibold"
                      : "border-black"
                  }`}
                >
                  YAML
                </button> */}
                <button
                  onClick={() => setView("form")}
                  className={`whitespace-nowrap border-b-2 py-2 px-2 text-md font-medium focus:outline-none ml-4 ${
                    view === "form"
                      ? "border-orange-700 text-orange-700 font-semibold"
                      : "border-black"
                  }`}
                >
                  Form
                </button>
              </div>

              {view === "yaml" ? (
                <div
                  id="yamlSpec"
                  className="p-4 bg-zinc-900 text-gray-100 rounded"
                >
                  {/* <pre className="whitespace-pre-wrap">
                    {JSON.stringify(guardrail, null, 2)}
                  </pre>
                  <div className="mt-4 flex justify-end">
                    <button className="px-6 py-2 bg-orange-700 text-white rounded-md shadow hover:bg-pink-900">
                      Submit
                    </button>
                  </div> */}
                </div>
              ) : (
                <div id="formSpec">
                  <form
                    id="dataSourceSpec"
                    className="space-y-2 border border-zinc-500 rounded-md text-gray-100 p-2"
                    onSubmit={submitUpdateForm}
                  >
                    <fieldset className="p-4 rounded">
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                        {/* Display Name */}
                        <div>
                          <label htmlFor="spec_displayName" className="labels">
                            Display Name
                          </label>
                          <input
                            id="spec_displayName"
                            name="spec.displayName"
                            type="text"
                            defaultValue={guardrail.displayName}
                            placeholder="Enter Display name"
                            className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                          />
                        </div>
                        {/* Name */}
                        <div>
                          <label htmlFor="spec_name" className="labels">
                            Name
                          </label>
                          <input
                            id="spec_name"
                            name="spec.name"
                            type="text"
                            defaultValue={guardrail.name}
                            placeholder="Enter name"
                            className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                          />
                        </div>
                        {/* Minimum Confidance */}
                        <div>
                          <label
                            htmlFor="spec_minConfidence"
                            className="labels"
                          >
                            Minimum Confidence
                          </label>
                          <input
                            id="spec_minConfidence"
                            name="spec.minConfidence"
                            type="number"
                            step="0.001"
                            defaultValue={guardrail.minConfidence}
                            placeholder="Enter Minimum Confidence"
                            className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                          />
                        </div>
                        {/* Maximum Confidence */}
                        <div>
                          <label
                            htmlFor="spec_maxConfidence"
                            className="labels"
                          >
                            Maximum Confidence
                          </label>
                          <input
                            id="spec_maxConfidence"
                            name="spec.maxConfidence"
                            type="number"
                            step="0.001"
                            defaultValue={guardrail.maxConfidence}
                            placeholder="Enter Maximum Confidence"
                            className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                          />
                        </div>
                        {/* DEscription */}
                        <div>
                          <label htmlFor="spec_Description" className="labels">
                            Description
                          </label>
                          <textarea
                            id="spec_Description"
                            name="spec.description"
                            placeholder="Enter Description"
                            rows="3"
                            defaultValue={guardrail.description}
                            className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-300 p-2 shadow-sm w-full"
                          ></textarea>
                        </div>
                        {/* Failure Message */}
                        <div>
                          <label
                            htmlFor="spec_failureMessage"
                            className="labels"
                          >
                            Failure Message
                          </label>
                          <textarea
                            id="spec_failureMessage"
                            name="spec.failureMessage"
                            placeholder="Enter failureMessage"
                            rows="3"
                            defaultValue={guardrail.failureMessage}
                            className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-300 p-2 shadow-sm w-full"
                          ></textarea>
                        </div>
                        {/* Scan Mode*/}
                        <div>
                          <label htmlFor="spec_scanMode" className="labels">
                            Scan Mode
                          </label>
                          <select
                            id="spec_scanMode"
                            name="spec.scanMode"
                            defaultValue={guardrail.scanMode}
                            className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                          >
                            <option value="INVALID_SSP">
                              Select Scan Mode
                            </option>
                            {enums.AIGuardrailScanMode.map((mode) => (
                              <option key={mode} value={mode}>
                                {strTitle(mode)}
                              </option>
                            ))}
                          </select>
                        </div>
                        {/* Scope*/}
                        <div>
                          <label
                            htmlFor="spec_resourceBase_scope"
                            className="labels"
                          >
                            Scope
                          </label>
                          <select
                            id="spec_resourceBase_scope"
                            name="spec.resourceBase.scope"
                            defaultValue={
                              guardrail.resourceBase?.scope || "DOMAIN_SCOPE"
                            }
                            className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                          >
                            <option value="DOMAIN_SCOPE">Select Scope</option>
                            {enums.ResourceScope.map((scope) => (
                              <option key={scope} value={scope}>
                                {strTitle(scope)}
                              </option>
                            ))}
                          </select>
                        </div>
                      </div>

                      {/* Contents Section */}
                      <details className="shadow-sm border border-zinc-500 rounded-md shadow-sm p-4 mb-4">
                        <summary className="text-lg font-semibold text-gray-100 cursor-pointer">
                          Contents
                        </summary>
                        <fieldset className="rounded mb-4">
                          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                            {/* Select Hate Speech Level */}
                            <div>
                              <label
                                htmlFor="spec_ContentsHatespeech"
                                className="labels"
                              >
                                Hate Speech Level
                              </label>
                              <select
                                id="spec_ContentsHatespeech"
                                name="spec.contents.hateSpeech"
                                defaultValue={
                                  guardrail.contents?.hateSpeech || ""
                                }
                                className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                              >
                                <option value="">
                                  Select Hate Speech Level
                                </option>
                                {enums.GuardRailLevels.map((level) => (
                                  <option key={level} value={level}>
                                    {strTitle(level)}
                                  </option>
                                ))}
                              </select>
                            </div>
                            {/*  Insults Level */}
                            <div>
                              <label
                                htmlFor="spec_ContentsInsults"
                                className="labels"
                              >
                                Insults Level
                              </label>
                              <select
                                id="spec_ContentsInsults"
                                name="spec.contents.insults"
                                defaultValue={guardrail.contents?.insults || ""}
                                className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                              >
                                <option value="">Select Insults Level</option>
                                {enums.GuardRailLevels.map((level) => (
                                  <option key={level} value={level}>
                                    {strTitle(level)}
                                  </option>
                                ))}
                              </select>
                            </div>
                            {/* Sexual Harassment Level */}
                            <div>
                              <label
                                htmlFor="spec_ContentsSexual"
                                className="labels"
                              >
                                Sexual Harassment Level
                              </label>
                              <select
                                id="spec_ContentsSexual"
                                name="spec.contents.sexual"
                                defaultValue={guardrail.contents?.sexual || ""}
                                className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                              >
                                <option value="">
                                  Select Sexual Harassment Level
                                </option>
                                {enums.GuardRailLevels.map((level) => (
                                  <option key={level} value={level}>
                                    {strTitle(level)}
                                  </option>
                                ))}
                              </select>
                            </div>
                            {/* Threats Level */}
                            <div>
                              <label
                                htmlFor="spec_ContentsThreats"
                                className="labels"
                              >
                                Threats Level
                              </label>
                              <select
                                id="spec_ContentsThreats"
                                name="spec.contents.threats"
                                defaultValue={guardrail.contents?.threats || ""}
                                className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                              >
                                <option value="">Select Threats Level</option>
                                {enums.GuardRailLevels.map((level) => (
                                  <option key={level} value={level}>
                                    {strTitle(level)}
                                  </option>
                                ))}
                              </select>
                            </div>
                            {/* Misconduct Level */}
                            <div>
                              <label
                                htmlFor="spec_ContentsMisconduct"
                                className="labels"
                              >
                                Misconduct Level
                              </label>
                              <select
                                id="spec_ContentsMisconduct"
                                name="spec.contents.misconduct"
                                defaultValue={
                                  guardrail.contents?.misconduct || ""
                                }
                                className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                              >
                                <option value="">
                                  Select Misconduct Level
                                </option>
                                {enums.GuardRailLevels.map((level) => (
                                  <option key={level} value={level}>
                                    {strTitle(level)}
                                  </option>
                                ))}
                              </select>
                            </div>
                          </div>
                        </fieldset>
                      </details>

                      {/* Topics Section */}
                      <details className="shadow-sm border border-zinc-500 rounded-md shadow-sm p-4 mb-4">
                        <summary className="text-lg font-semibold text-gray-100 cursor-pointer">
                          Topics
                        </summary>
                        <fieldset className="p-4 rounded">
                          <div id="topicsContainer">
                            {Array.from({ length: topicCount }).map(
                              (_, index) => {
                                const topic =
                                  guardrail.topics && guardrail.topics[index]
                                    ? guardrail.topics[index]
                                    : null;
                                return (
                                  <div
                                    key={`topic-${index}`}
                                    className="shadow-sm border border-zinc-500 rounded-md shadow-sm p-4 mb-4"
                                  >
                                    <label
                                      htmlFor={`spec_Topics_${index}_topic`}
                                      className="labels"
                                    >
                                      Topic Name
                                    </label>
                                    <input
                                      id={`spec_Topics_${index}_topic`}
                                      type="text"
                                      name={`spec.topics[${index}].topic`}
                                      defaultValue={topic?.topic || ""}
                                      placeholder="Enter topic name"
                                      className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full mb-2"
                                    />

                                    <label
                                      htmlFor={`spec_Topics_${index}_description`}
                                      className="labels"
                                    >
                                      Topic Description
                                    </label>
                                    <textarea
                                      id={`spec_Topics_${index}_description`}
                                      name={`spec.topics[${index}].description`}
                                      defaultValue={topic?.description || ""}
                                      placeholder="Enter description"
                                      className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-300 p-2 shadow-sm w-full"
                                    ></textarea>

                                    <label
                                      htmlFor={`spec_Topics_${index}_samples`}
                                      className="labels"
                                    >
                                      Topic Samples
                                    </label>
                                    <input
                                      id={`spec_Topics_${index}_samples`}
                                      type="text"
                                      name={`spec.topics[${index}].samples`}
                                      defaultValue={
                                        topic?.samples
                                          ? topic.samples.join(",")
                                          : ""
                                      }
                                      placeholder="Enter samples (comma separated)"
                                      className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                    />
                                    {index > 0 && (
                                      <RemoveButton
                                        onClick={() => removeTopic(index)}
                                      />
                                    )}
                                  </div>
                                );
                              }
                            )}
                          </div>
                          <AddButton onClick={addTopic} name="+ Add Topic" />
                        </fieldset>
                      </details>

                      {/* Words Section */}
                      <details className="shadow-sm border border-zinc-500 rounded-md shadow-sm p-4 mb-4">
                        <summary className="text-lg font-semibold text-gray-100 cursor-pointer">
                          Words
                        </summary>
                        <fieldset className="rounded mb-4">
                          <div id="wordsContainer">
                            {Array.from({ length: wordEntryCount }).map(
                              (_, index) => {
                                const wordEntry =
                                  guardrail.words && guardrail.words[index]
                                    ? guardrail.words[index]
                                    : null;
                                return (
                                  <div
                                    key={`word-${index}`}
                                    className="word-entry border p-3 rounded mb-3"
                                  >
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                      <div>
                                        <label className="labels">
                                          Words (comma separated)
                                        </label>
                                        <input
                                          type="text"
                                          name={`spec.words[${index}].words`}
                                          defaultValue={
                                            wordEntry?.words
                                              ? wordEntry.words.join(",")
                                              : ""
                                          }
                                          placeholder="Enter words (comma separated)"
                                          className=" form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                        />
                                      </div>
                                      <div>
                                        <label className="labels">
                                          File Location
                                        </label>
                                        <input
                                          type="text"
                                          name={`spec.words[${index}].fileLocation`}
                                          defaultValue={
                                            wordEntry?.fileLocation || ""
                                          }
                                          placeholder="Enter file location"
                                          className=" form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                        />
                                      </div>
                                    </div>
                                    {index > 0 && (
                                      <RemoveButton
                                        onClick={() => removeWordEntry(index)}
                                      />
                                    )}
                                  </div>
                                );
                              }
                            )}
                          </div>
                          <AddButton
                            name="+ Add Word Entry"
                            onClick={addWordEntry}
                          />
                        </fieldset>
                      </details>

                      {/* Sensitive Dataset Section */}
                      <details className="shadow-sm border border-zinc-500 rounded-md shadow-sm p-4 mb-4">
                        <summary className="text-lg font-semibold text-gray-100 cursor-pointer">
                          Sensitive Dataset
                        </summary>
                        <fieldset className="rounded mb-4">
                          <div id="sensitiveDatasetContainer">
                            {Array.from({ length: sensitiveEntryCount }).map(
                              (_, index) => {
                                const sensitiveEntry =
                                  guardrail.sensitiveDataset &&
                                  guardrail.sensitiveDataset[index]
                                    ? guardrail.sensitiveDataset[index]
                                    : null;
                                return (
                                  <div
                                    key={`sensitive-${index}`}
                                    className="sensitive-entry  border p-3 rounded mb-3"
                                  >
                                    <div className="grid grid-cols-1  md:grid-cols-3 gap-4">
                                      <div>
                                        <label className="labels">
                                          PII Type
                                        </label>
                                        <input
                                          type="text"
                                          name={`spec.sensitiveDataset[${index}].piiType`}
                                          defaultValue={
                                            sensitiveEntry?.piiType || ""
                                          }
                                          placeholder="Enter PII type"
                                          className=" form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                        />
                                      </div>
                                      <div>
                                        <label className="labels">Action</label>
                                        <input
                                          type="text"
                                          name={`spec.sensitiveDataset[${index}].action`}
                                          defaultValue={
                                            sensitiveEntry?.action || ""
                                          }
                                          placeholder="Enter action"
                                          className=" form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                        />
                                      </div>
                                      <div>
                                        <label className="labels">Regex</label>
                                        <input
                                          type="text"
                                          name={`spec.sensitiveDataset[${index}].regex`}
                                          defaultValue={
                                            sensitiveEntry?.regex || ""
                                          }
                                          placeholder="Enter regex pattern"
                                          className=" form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                        />
                                      </div>
                                    </div>
                                    {index > 0 && (
                                      <RemoveButton
                                        onClick={() =>
                                          removeSensitiveEntry(index)
                                        }
                                      />
                                    )}
                                  </div>
                                );
                              }
                            )}
                          </div>
                          <AddButton
                            name="+ Add Senstive Data Entry"
                            onClick={addSensitiveEntry}
                          />
                        </fieldset>
                      </details>

                      {/* Guard Model Section */}
                      <details className="shadow-sm border border-zinc-500 rounded-md shadow-sm p-4 mb-4">
                        <summary className="text-lg font-semibold text-gray-100 cursor-pointer">
                          Guard Model
                        </summary>
                        <fieldset className="rounded mb-4">
                          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                              <label className="labels">Model Node ID</label>
                              <div className="relative">
                                <select
                                  id="spec_guardModel_modelNodeId"
                                  name="spec.guardModel.modelNodeId"
                                  onChange={(e) =>
                                    populateModelDropdown(e.target.value)
                                  }
                                  defaultValue={
                                    guardrail.guardModel?.modelNodeId || ""
                                  }
                                  className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                >
                                  <option value="">Select Model Gateway</option>
                                  {aiModelNodes.map((node) => (
                                    <option
                                      key={node.modelNodeId}
                                      value={node.modelNodeId}
                                      className="text-sm text-orange-700 hover:text-pink-900"
                                    >
                                      {node.name}
                                    </option>
                                  ))}
                                </select>
                              </div>
                            </div>
                            <div>
                              <label className="labels">Model ID</label>
                              <div className="relative">
                                <select
                                  id="spec_guardModel_modelId"
                                  name="spec.guardModel.modelId"
                                  defaultValue={
                                    guardrail.guardModel?.modelId || ""
                                  }
                                  className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                >
                                  <option value="">Select Model</option>
                                  {availableModels.map((model) => (
                                    <option
                                      key={model.modelId}
                                      value={model.modelId}
                                      className="text-sm text-orange-700 hover:text-pink-900"
                                    >
                                      {model.modelName}
                                    </option>
                                  ))}
                                </select>
                              </div>
                            </div>
                          </div>
                        </fieldset>
                      </details>

                      {/* Submit Button */}
                      <div className="mt-4 flex justify-end space-x-2">
                        <button
                          type="submit"
                          disabled={loading}
                          className="px-6 py-2 bg-orange-700 text-white rounded-md shadow hover:bg-pink-900"
                        >
                          {loading ? (
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
                          ) : (
                            "Submit"
                          )}
                        </button>
                      </div>
                    </fieldset>
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
