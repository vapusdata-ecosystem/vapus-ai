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
import {
  bedrockGuardrailsApi,
  enumsApi,
} from "@/app/utils/developers-endpoint/enums";
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

  // Dynamic form sections state - Arrays instead of just counts
  const [topics, setTopics] = useState([
    { topic: "", description: "", samples: [] },
  ]);
  const [wordEntries, setWordEntries] = useState([
    { words: "", fileLocation: "" },
  ]);
  const [sensitiveEntries, setSensitiveEntries] = useState([
    { piiType: "", action: "", regex: "" },
  ]);

  // Models data state
  const [aiModelNodes, setAiModelNodes] = useState([]);
  const [modelMap, setModelMap] = useState({});
  const [modelOptions, setModelOptions] = useState([]);
  const [selectedModelNodeId, setSelectedModelNodeId] = useState("");
  const [view, setView] = useState("form");

  // Guardrail selection state
  const [selectedGuardrails, setSelectedGuardrails] = useState([]);
  const [guardrailType, setGuardrailType] = useState("vapus");
  const [guardrailProviders, setGuardrailProviders] = useState([]);
  const [guardrailTypes, setGuardrailTypes] = useState({});
  const [bedrockGuardrails, setBedrockGuardrails] = useState([]);
  const [pangeaGuardrails, setPangeaGuardrails] = useState([]);
  const [mistralGuardrails, setMistralGuardrails] = useState([]);

  // UI state
  const [loading, setLoading] = useState(true);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [enums, setEnums] = useState({
    AIGuardrailScanMode: [],
    ResourceScope: [],
    GuardRailLevels: [],
  });

  // Handle guardrail type change
  const handleGuardrailChange = async (e) => {
    const selectedType = e.target.value;
    setGuardrailType(selectedType);
    setSelectedGuardrails([]);

    // Fetch guardrails data for all three types
    if (
      selectedType === "bedrock" ||
      selectedType === "pangea" ||
      selectedType === "mistral"
    ) {
      try {
        setIsLoading(true);
        const response = await bedrockGuardrailsApi.getBedrockGuardrailsApi();

        // Extract data from the new JSON structure
        const bedrockData = response.bedrockOutput || [];
        const pangeaData = response.pangeaOutput || [];
        const mistralData = response.mistralOutput || [];

        setBedrockGuardrails(bedrockData);
        setPangeaGuardrails(pangeaData);
        setMistralGuardrails(mistralData);

        console.log("Bedrock guardrails:", bedrockData);
        console.log("Pangea guardrails:", pangeaData);
        console.log("Mistral guardrails:", mistralData);
      } catch (error) {
        console.error("Failed to fetch guardrails:", error);
        toast.error("Failed to fetch guardrails");
        setBedrockGuardrails([]);
        setPangeaGuardrails([]);
        setMistralGuardrails([]);
      } finally {
        setIsLoading(false);
      }
    } else {
      // Reset all guardrail arrays for other types
      setBedrockGuardrails([]);
      setPangeaGuardrails([]);
      setMistralGuardrails([]);
    }
  };

  // Handle guardrail selection
  const handleGuardrailSelection = (guardrail) => {
  setSelectedGuardrails((prev) => {
    const isSelected = prev.some(item => item.id === guardrail.id);
    if (isSelected) {
      // Remove the guardrail
      const filtered = prev.filter(item => item.id !== guardrail.id);
      return filtered;
    } else {
      // Add the guardrail
      const newSelection = [...prev, guardrail];
      return newSelection;
    }
  });
};
  // For Mistral and Pangea in the JSX:
  const isGuardrailSelected = (guardrail) => {
  const isSelected = selectedGuardrails.some(item => item.id === guardrail.id);
  console.log(`Checking if ${guardrail.name || guardrail.id} is selected:`, isSelected);
  return isSelected;
};

  // Fetch enums data
  useEffect(() => {
    const fetchEnumsData = async () => {
      try {
        setIsLoading(true);
        const response = await enumsApi.getEnums();
        const enumResponses = response.enumResponse || [];
        const guardrailTypesData = response.guardrailTypes || {};

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
        setGuardrailTypes(guardrailTypesData);
        const providers = Object.keys(guardrailTypesData);
        setGuardrailProviders(providers);

        console.log("Guardrail Types Data:", guardrailTypesData);
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

        // Fetch guardrail data and model nodes in parallel
        const [guardrailResponse, modelMap] = await Promise.all([
          GuardrailApi.getGuardrailId(guardrail_id),
          fetchModelNodes(),
        ]);

        const guardrailData = guardrailResponse.output?.[0];

        if (guardrailData) {
          setGuardrail(guardrailData);

          // Initialize form arrays with existing data
          if (guardrailData.topics && guardrailData.topics.length > 0) {
            setTopics(guardrailData.topics);
          } else {
            setTopics([{ topic: "", description: "", samples: [] }]);
          }

          if (guardrailData.words && guardrailData.words.length > 0) {
            setWordEntries(guardrailData.words);
          } else {
            setWordEntries([{ words: "", fileLocation: "" }]);
          }

          if (
            guardrailData.sensitiveDataset &&
            guardrailData.sensitiveDataset.length > 0
          ) {
            setSensitiveEntries(guardrailData.sensitiveDataset);
          } else {
            setSensitiveEntries([{ piiType: "", action: "", regex: "" }]);
          }

          // Set model node selection
          if (guardrailData.guardModel?.modelNodeId) {
            setSelectedModelNodeId(guardrailData.guardModel.modelNodeId);
            // Set model options for selected node
            const selectedNode = modelMap[guardrailData.guardModel.modelNodeId];
            if (selectedNode?.attributes?.generativeModels) {
              setModelOptions(selectedNode.attributes.generativeModels);
            }
          }

          // Initialize partner data if exists
          if (guardrailData.partner && guardrailData.partner.length > 0) {
            const partnerData = guardrailData.partner[0];
            if (partnerData.bedrock && partnerData.bedrock.length > 0) {
              setGuardrailType("bedrock");
              setSelectedGuardrails(partnerData.bedrock);
              // Fetch all guardrails data
              try {
                const response =
                  await bedrockGuardrailsApi.getBedrockGuardrailsApi();
                setBedrockGuardrails(response.bedrockOutput || []);
                setPangeaGuardrails(response.pangeaOutput || []);
                setMistralGuardrails(response.mistralOutput || []);
              } catch (error) {
                console.error("Failed to fetch guardrails:", error);
              }
            } else if (partnerData.mistral && partnerData.mistral.length > 0) {
              setGuardrailType("mistral");
              setSelectedGuardrails(partnerData.mistral);
              // Fetch all guardrails data
              try {
                const response =
                  await bedrockGuardrailsApi.getBedrockGuardrailsApi();
                setBedrockGuardrails(response.bedrockOutput || []);
                setPangeaGuardrails(response.pangeaOutput || []);
                setMistralGuardrails(response.mistralOutput || []);
              } catch (error) {
                console.error("Failed to fetch guardrails:", error);
              }
            } else if (partnerData.pangea && partnerData.pangea.length > 0) {
              setGuardrailType("pangea");
              setSelectedGuardrails(partnerData.pangea);
              // Fetch all guardrails data
              try {
                const response =
                  await bedrockGuardrailsApi.getBedrockGuardrailsApi();
                setBedrockGuardrails(response.bedrockOutput || []);
                setPangeaGuardrails(response.pangeaOutput || []);
                setMistralGuardrails(response.mistralOutput || []);
              } catch (error) {
                console.error("Failed to fetch guardrails:", error);
              }
            }
          }
        } else {
          setError("No guardrails found");
          toast.error("Guardrail not found");
        }
      } catch (error) {
        console.error("Error fetching guardrail data:", error);
        setError(`Failed to load guardrail: ${error.message}`);
        toast.error(`Failed to load guardrail: ${error.message}`);
      } finally {
        setLoading(false);
      }
    };
    fetchGuardrailsData();
  }, [guardrail_id]);

  // Fetch model nodes data
  const fetchModelNodes = async () => {
    try {
      const response = await modelsRegistryApi.getModelsRegistry();
      const nodes = response.output?.aiModelNodes || [];
      setAiModelNodes(nodes);

      const map = {};
      nodes.forEach((item) => {
        map[item.modelNodeId] = item;
      });
      setModelMap(map);

      return map;
    } catch (error) {
      console.error("Failed to fetch model nodes:", error);
      setError(`Failed to load model nodes: ${error.message}`);
      return {};
    }
  };

  // Populate available models based on selected node
  const populateModelDropdown = (nodeId) => {
    setSelectedModelNodeId(nodeId);
    const selectedNode = aiModelNodes.find(
      (node) => node.modelNodeId === nodeId
    );
    if (selectedNode?.attributes?.generativeModels) {
      setModelOptions(selectedNode.attributes.generativeModels);
    } else {
      setModelOptions([]);
    }
  };

  // Topic management functions
  const addTopic = () => {
    setTopics([...topics, { topic: "", description: "", samples: [] }]);
  };

  const removeTopic = (index) => {
    if (topics.length > 1) {
      const newTopics = topics.filter((_, i) => i !== index);
      setTopics(newTopics);
    }
  };

  const updateTopic = (index, field, value) => {
    const newTopics = [...topics];
    if (field === "samples") {
      newTopics[index][field] = value
        .split(",")
        .map((s) => s.trim())
        .filter((s) => s);
    } else {
      newTopics[index][field] = value;
    }
    setTopics(newTopics);
  };

  // Word entry management functions
  const addWordEntry = () => {
    setWordEntries([...wordEntries, { words: "", fileLocation: "" }]);
  };

  const removeWordEntry = (index) => {
    if (wordEntries.length > 1) {
      const newEntries = wordEntries.filter((_, i) => i !== index);
      setWordEntries(newEntries);
    }
  };

  const updateWordEntry = (index, field, value) => {
    const newEntries = [...wordEntries];
    if (field === "words") {
      newEntries[index][field] = value
        .split(",")
        .map((w) => w.trim())
        .filter((w) => w);
    } else {
      newEntries[index][field] = value;
    }
    setWordEntries(newEntries);
  };

  // Sensitive entry management functions
  const addSensitiveEntry = () => {
    setSensitiveEntries([
      ...sensitiveEntries,
      { piiType: "", action: "", regex: "" },
    ]);
  };

  const removeSensitiveEntry = (index) => {
    if (sensitiveEntries.length > 1) {
      const newEntries = sensitiveEntries.filter((_, i) => i !== index);
      setSensitiveEntries(newEntries);
    }
  };

  const updateSensitiveEntry = (index, field, value) => {
    const newEntries = [...sensitiveEntries];
    newEntries[index][field] = value;
    setSensitiveEntries(newEntries);
  };

  // Convert numeric strings function
  const convertNumericStrings = (obj) => {
    if (obj === null || typeof obj !== "object") {
      return obj;
    }

    if (Array.isArray(obj)) {
      return obj.map(convertNumericStrings);
    }

    const result = {};
    for (const key in obj) {
      if (Object.prototype.hasOwnProperty.call(obj, key)) {
        if (key === "minConfidence" || key === "maxConfidence") {
          result[key] = parseFloat(obj[key] || 0);
        } else if (typeof obj[key] === "object" && obj[key] !== null) {
          result[key] = convertNumericStrings(obj[key]);
        } else {
          result[key] = obj[key];
        }
      }
    }

    return result;
  };

  // Submit update form
  const submitUpdateForm = async (event) => {
    event.preventDefault();

    try {
      setIsLoading(true);
      const formData = new FormData(event.target);

      // Build the data object from form
      const dataObj = {
        displayName:
          formData.get("spec.displayName") || guardrail.displayName || "",
        name: formData.get("spec.name") || guardrail.name || "",
        minConfidence: parseFloat(
          formData.get("spec.minConfidence") || guardrail.minConfidence || 0
        ),
        maxConfidence: parseFloat(
          formData.get("spec.maxConfidence") || guardrail.maxConfidence || 0
        ),
        description:
          formData.get("spec.description") || guardrail.description || "",
        failureMessage:
          formData.get("spec.failureMessage") || guardrail.failureMessage || "",
        scanMode:
          formData.get("spec.scanMode") ||
          guardrail.scanMode ||
          "BIDIRECTIONAL",
        resourceBase: {
          ...guardrail.resourceBase,
          scope:
            formData.get("spec.resourceBase.scope") ||
            guardrail.resourceBase?.scope ||
            "ORGANIZATION_SCOPE",
        },
        contents: {
          hateSpeech:
            formData.get("spec.contents.hateSpeech") ||
            guardrail.contents?.hateSpeech ||
            "NONE",
          insults:
            formData.get("spec.contents.insults") ||
            guardrail.contents?.insults ||
            "NONE",
          sexual:
            formData.get("spec.contents.sexual") ||
            guardrail.contents?.sexual ||
            "NONE",
          threats:
            formData.get("spec.contents.threats") ||
            guardrail.contents?.threats ||
            "NONE",
          misconduct:
            formData.get("spec.contents.misconduct") ||
            guardrail.contents?.misconduct ||
            "NONE",
        },
        topics: topics,
        words: wordEntries,
        sensitiveDataset: sensitiveEntries,
        guardModel: {
          modelNodeId:
            formData.get("spec.guardModel.modelNodeId") ||
            guardrail.guardModel?.modelNodeId ||
            "",
          modelId:
            formData.get("spec.guardModel.modelId") ||
            guardrail.guardModel?.modelId ||
            "",
        },
        guardrailId: guardrail_id,
        eligibleModelNodes: guardrail.eligibleModelNodes || [],
        partner: [],
      };

      // Prepare partner data
      if (selectedGuardrails && selectedGuardrails.length > 0) {
        const partnerObj = {};

        if (guardrailType === "bedrock") {
          partnerObj.bedrock = selectedGuardrails;
        } else if (guardrailType === "mistral") {
          partnerObj.mistral = selectedGuardrails;
        } else if (guardrailType === "pangea") {
          partnerObj.pangea = selectedGuardrails;
        }

        if (Object.keys(partnerObj).length > 0) {
          dataObj.partner = [partnerObj];
        }
      }

      const payload = { spec: dataObj };

      console.log("Update Payload:", payload);

      const output = await GuardrailUpdateFormApi.getGuardrailUpdateForm(
        payload
      );

      console.log("Resource updated:", output);

      const resourceInfo = output.result;
      if (resourceInfo) {
        toast.success(
          `${resourceInfo.resource} Resource updated successfully.`
        );
        setTimeout(() => {
          router.push(`/ai-center/guardrails/${guardrail_id}`);
        }, 1000);
      } else {
        toast.success("Resource updated successfully.");
        setTimeout(() => {
          router.push(`/ai-center/guardrails/${guardrail_id}`);
        }, 1000);
      }
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
        <LoadingOverlay isLoading={isLoading} isOverlay={true}/>

        <div className="flex-grow p-4 overflow-y-auto w-full">
          <section className="space-y-2">
            <div className="max-w-6xl mx-auto bg-[#1b1b1b] shadow rounded-lg p-2">
              <div className="text-gray-100 mb-2 flex justify-center">
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
                  {/* YAML content would go here */}
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
                        {/* Guardrails Type */}
                        <div>
                          <label
                            htmlFor="spec_scanMode"
                            className="block text-sm font-medium mb-1"
                          >
                            Guardrail Type
                          </label>
                          <select
                            id="spec_guardrailType"
                            name="spec.guardrailType"
                            className="w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium rounded border border-zinc-700 focus:ring-1 focus:ring-orange-700 focus:outline-none"
                            value={guardrailType}
                            onChange={handleGuardrailChange}
                            suppressHydrationWarning
                          >
                            {guardrailProviders.map((provider) => (
                              <option key={provider} value={provider}>
                                {strTitle(provider)}
                              </option>
                            ))}
                          </select>
                        </div>
                      </div>

                      {/* Consolidated Guardrails Section */}
                      {(guardrailType === "vapus" ||
                        guardrailType === "bedrock" ||
                        guardrailType === "mistral" ||
                        guardrailType === "pangea") && (
                        <div className="shadow-sm border border-zinc-500 rounded-md shadow-sm p-4 mb-4">
                          {/* Dynamic Legend based on guardrail type */}
                          <legend className="text-xl font-bold text-gray-100">
                            {guardrailType === "vapus"
                              ? "Vapus Guardrails"
                              : `${guardrailType} Available Guardrails`}
                          </legend>

                          {/* Vapus-specific content */}
                          {guardrailType === "vapus" && (
                            <>
                              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
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
                                {/* Description */}
                                <div>
                                  <label
                                    htmlFor="spec_Description"
                                    className="labels"
                                  >
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
                                  <label
                                    htmlFor="spec_scanMode"
                                    className="labels"
                                  >
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
                                      guardrail.resourceBase?.scope ||
                                      "ORGANIZATION_SCOPE"
                                    }
                                    className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                  >
                                    <option value="ORGANIZATION_SCOPE">
                                      Select Scope
                                    </option>
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
                                        defaultValue={
                                          guardrail.contents?.insults || ""
                                        }
                                        className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                      >
                                        <option value="">
                                          Select Insults Level
                                        </option>
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
                                        defaultValue={
                                          guardrail.contents?.sexual || ""
                                        }
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
                                        defaultValue={
                                          guardrail.contents?.threats || ""
                                        }
                                        className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                      >
                                        <option value="">
                                          Select Threats Level
                                        </option>
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
                                    {topics.map((topic, index) => (
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
                                          value={topic.topic}
                                          onChange={(e) =>
                                            updateTopic(
                                              index,
                                              "topic",
                                              e.target.value
                                            )
                                          }
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
                                          value={topic.description}
                                          onChange={(e) =>
                                            updateTopic(
                                              index,
                                              "description",
                                              e.target.value
                                            )
                                          }
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
                                          value={
                                            Array.isArray(topic.samples)
                                              ? topic.samples.join(",")
                                              : ""
                                          }
                                          onChange={(e) =>
                                            updateTopic(
                                              index,
                                              "samples",
                                              e.target.value
                                            )
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
                                    ))}
                                  </div>
                                  <AddButton
                                    onClick={addTopic}
                                    name="+ Add Topic"
                                  />
                                </fieldset>
                              </details>

                              {/* Words Section */}
                              <details className="shadow-sm border border-zinc-500 rounded-md shadow-sm p-4 mb-4">
                                <summary className="text-lg font-semibold text-gray-100 cursor-pointer">
                                  Words
                                </summary>
                                <fieldset className="rounded mb-4">
                                  <div id="wordsContainer">
                                    {wordEntries.map((wordEntry, index) => (
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
                                              value={
                                                Array.isArray(wordEntry.words)
                                                  ? wordEntry.words.join(",")
                                                  : ""
                                              }
                                              onChange={(e) =>
                                                updateWordEntry(
                                                  index,
                                                  "words",
                                                  e.target.value
                                                )
                                              }
                                              placeholder="Enter words (comma separated)"
                                              className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                            />
                                          </div>
                                          <div>
                                            <label className="labels">
                                              File Location
                                            </label>
                                            <input
                                              type="text"
                                              value={wordEntry.fileLocation}
                                              onChange={(e) =>
                                                updateWordEntry(
                                                  index,
                                                  "fileLocation",
                                                  e.target.value
                                                )
                                              }
                                              placeholder="Enter file location"
                                              className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                            />
                                          </div>
                                        </div>
                                        {index > 0 && (
                                          <RemoveButton
                                            onClick={() =>
                                              removeWordEntry(index)
                                            }
                                          />
                                        )}
                                      </div>
                                    ))}
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
                                    {sensitiveEntries.map(
                                      (sensitiveEntry, index) => (
                                        <div
                                          key={`sensitive-${index}`}
                                          className="sensitive-entry border p-3 rounded mb-3"
                                        >
                                          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                                            <div>
                                              <label className="labels">
                                                PII Type
                                              </label>
                                              <input
                                                type="text"
                                                value={sensitiveEntry.piiType}
                                                onChange={(e) =>
                                                  updateSensitiveEntry(
                                                    index,
                                                    "piiType",
                                                    e.target.value
                                                  )
                                                }
                                                placeholder="Enter PII type"
                                                className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                              />
                                            </div>
                                            <div>
                                              <label className="labels">
                                                Action
                                              </label>
                                              <input
                                                type="text"
                                                value={sensitiveEntry.action}
                                                onChange={(e) =>
                                                  updateSensitiveEntry(
                                                    index,
                                                    "action",
                                                    e.target.value
                                                  )
                                                }
                                                placeholder="Enter action"
                                                className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                              />
                                            </div>
                                            <div>
                                              <label className="labels">
                                                Regex
                                              </label>
                                              <input
                                                type="text"
                                                value={sensitiveEntry.regex}
                                                onChange={(e) =>
                                                  updateSensitiveEntry(
                                                    index,
                                                    "regex",
                                                    e.target.value
                                                  )
                                                }
                                                placeholder="Enter regex pattern"
                                                className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
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
                                      )
                                    )}
                                  </div>
                                  <AddButton
                                    name="+ Add Sensitive Data Entry"
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
                                      <label className="labels">
                                        Model Node ID
                                      </label>
                                      <div className="relative">
                                        <select
                                          id="spec_guardModel_modelNodeId"
                                          name="spec.guardModel.modelNodeId"
                                          onChange={(e) =>
                                            populateModelDropdown(
                                              e.target.value
                                            )
                                          }
                                          value={selectedModelNodeId}
                                          className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                        >
                                          <option value="">
                                            Select Model Gateway
                                          </option>
                                          {aiModelNodes.map((node) => (
                                            <option
                                              key={node.modelNodeId}
                                              value={node.modelNodeId}
                                              className="text-sm"
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
                                          {modelOptions.map((model) => (
                                            <option
                                              key={model.modelId}
                                              value={model.modelId}
                                              className="text-sm "
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
                            </>
                          )}

                          {/* Common guardrail selection for Bedrock, Mistral, and Pangea */}
                          {(guardrailType === "bedrock" ||
                            guardrailType === "mistral" ||
                            guardrailType === "pangea") && (
                            <>
                              <label className="labels">
                                Select Guardrails
                              </label>
                              <div className="mt-2">
                                <div className="rounded mb-4 max-h-36 scrollbar overflow-y-auto">
                                  <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-4 gap-4">
                                    {/* Bedrock guardrails */}
                                    {guardrailType === "bedrock" &&
                                      bedrockGuardrails.map((guardrail) => {
                                        const isSelected =
                                          selectedGuardrails.some(
                                            (item) =>
                                              typeof item === "object" &&
                                              item.id === guardrail.id
                                          );

                                        return (
                                          <div
                                            key={guardrail.id}
                                            className={`border rounded-md p-3 cursor-pointer transition-all duration-200 flex items-center justify-between ${
                                              isSelected
                                                ? "border-orange-700 bg-zinc-700"
                                                : "border-zinc-600 hover:border-orange-700 hover:bg-zinc-700"
                                            }`}
                                            onClick={() =>
                                              handleGuardrailSelection(
                                                guardrail
                                              )
                                            }
                                          >
                                            <div className="text-sm font-medium">
                                            {strTitle(guardrail.name || guardrail.Name || `Guardrail ${guardrail.id}`)}
                                            </div>
                                            {isSelected && (
                                              <div className="text-orange-500">
                                                <svg
                                                  xmlns="http://www.w3.org/2000/svg"
                                                  width="16"
                                                  height="16"
                                                  viewBox="0 0 24 24"
                                                  fill="none"
                                                  stroke="currentColor"
                                                  strokeWidth="2"
                                                  strokeLinecap="round"
                                                  strokeLinejoin="round"
                                                >
                                                  <polyline points="20 6 9 17 4 12"></polyline>
                                                </svg>
                                              </div>
                                            )}
                                          </div>
                                        );
                                      })}

                                    {/* Pangea guardrails */}
                                    {guardrailType === "pangea" &&
                                      pangeaGuardrails.map((guardrail) => {
                                        const isSelected =
                                          selectedGuardrails.some(
                                            (item) =>
                                              typeof item === "object" &&
                                              item.id === guardrail.id
                                          );

                                        return (
                                          <div
                                            key={guardrail.id}
                                            className={`border rounded-md p-3 cursor-pointer transition-all duration-200 flex items-center justify-between ${
                                              isSelected
                                                ? "border-orange-700 bg-zinc-700"
                                                : "border-zinc-600 hover:border-orange-700 hover:bg-zinc-700"
                                            }`}
                                            onClick={() =>
                                              handleGuardrailSelection(
                                                guardrail
                                              )
                                            }
                                          >
                                            <div className="text-sm font-medium">
                                             {strTitle(guardrail.name || guardrail.Name || `Guardrail ${guardrail.id}`)}
                                            </div>
                                            {isSelected && (
                                              <div className="text-orange-500">
                                                <svg
                                                  xmlns="http://www.w3.org/2000/svg"
                                                  width="16"
                                                  height="16"
                                                  viewBox="0 0 24 24"
                                                  fill="none"
                                                  stroke="currentColor"
                                                  strokeWidth="2"
                                                  strokeLinecap="round"
                                                  strokeLinejoin="round"
                                                >
                                                  <polyline points="20 6 9 17 4 12"></polyline>
                                                </svg>
                                              </div>
                                            )}
                                          </div>
                                        );
                                      })}

                                    {/* Mistral guardrails */}
                                    {guardrailType === "mistral" &&
                                      mistralGuardrails.map((guardrail) => {
                                        const isSelected =
                                          selectedGuardrails.some(
                                            (item) =>
                                              typeof item === "object" &&
                                              item.id === guardrail.id
                                          );

                                        return (
                                          <div
                                            key={guardrail.id}
                                            className={`border rounded-md p-3 cursor-pointer transition-all duration-200 flex items-center justify-between ${
                                              isSelected
                                                ? "border-orange-700 bg-zinc-700"
                                                : "border-zinc-600 hover:border-orange-700 hover:bg-zinc-700"
                                            }`}
                                            onClick={() =>
                                              handleGuardrailSelection(
                                                guardrail
                                              )
                                            }
                                          >
                                            <div className="text-sm font-medium">
                                             {strTitle(guardrail.name || guardrail.Name || `Guardrail ${guardrail.id}`)}
                                            </div>
                                            {isSelected && (
                                              <div className="text-orange-500">
                                                <svg
                                                  xmlns="http://www.w3.org/2000/svg"
                                                  width="16"
                                                  height="16"
                                                  viewBox="0 0 24 24"
                                                  fill="none"
                                                  stroke="currentColor"
                                                  strokeWidth="2"
                                                  strokeLinecap="round"
                                                  strokeLinejoin="round"
                                                >
                                                  <polyline points="20 6 9 17 4 12"></polyline>
                                                </svg>
                                              </div>
                                            )}
                                          </div>
                                        );
                                      })}
                                  </div>
                                </div>
                                {/* Hidden input to store selected guardrails for form submission */}
                                <input
                                  type="hidden"
                                  name="selectedGuardrails"
                                  value={JSON.stringify(selectedGuardrails)}
                                />
                              </div>
                            </>
                          )}
                        </div>
                      )}

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
