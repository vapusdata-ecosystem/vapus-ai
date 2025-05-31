"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Head from "next/head";
import { toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import Header from "@/app/components/platform/header";
import YamlEditorClient from "@/app/components/formcomponets/ymal";
import { modelsRegistryApi } from "@/app/utils/ai-studio-endpoint/models-registry-api";
import ToastContainerMessage from "@/app/components/notification/customToast";
import LoadingOverlay from "@/app/components/loading/loading";
import { GuardrailFormApi } from "@/app/utils/ai-studio-endpoint/guardrails-api";
import AddButton from "@/app/components/buttons/addButton";
import RemoveButton from "@/app/components/buttons/removeButton";
import {
  bedrockGuardrailsApi,
  enumsApi,
} from "@/app/utils/developers-endpoint/enums";
import { strTitle } from "@/app/components/JS/common";

export default function CreateGuardrail() {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState("form");
  const [topicCount, setTopicCount] = useState(1);
  const [wordEntryCount, setWordEntryCount] = useState(1);
  const [sensitiveEntryCount, setSensitiveEntryCount] = useState(1);
  const [modelMap, setModelMap] = useState({});
  const [selectedModelNodeId, setSelectedModelNodeId] = useState("");
  const [modelOptions, setModelOptions] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [aiModelNodes, setAiModelNodes] = useState([]);
  const [selectedGuardrails, setSelectedGuardrails] = useState([]);
  const [guardrailType, setGuardrailType] = useState("vapus");
  const [guardrailProviders, setGuardrailProviders] = useState([]);
  const [guardrailTypes, setGuardrailTypes] = useState({});
  const [bedrockGuardrails, setBedrockGuardrails] = useState([]);
  const [pangeaGuardrails, setPangeaGuardrails] = useState([]);
  const [mistralGuardrails, setMistralGuardrails] = useState([]);

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

  const [wordEntries, setWordEntries] = useState([
    { words: "", fileLocation: "" },
  ]);
  const [sensitiveEntries, setSensitiveEntries] = useState([
    { piiType: "", action: "", regex: "" },
  ]);
  const [topics, setTopics] = useState([
    { topic: "", description: "", samples: [] },
  ]);
  const [enums, setEnums] = useState({
    AIGuardrailScanMode: [],
    ResourceScope: [],
    GuardRailLevels: [],
  });

  const handleGuardrailSelection = (guardrail) => {
    setSelectedGuardrails((prev) => {
      // All three types now use object format with id and Name
      const isSelected = prev.some(
        (item) => typeof item === "object" && item.id === guardrail.id
      );

      if (isSelected) {
        return prev.filter(
          (item) => !(typeof item === "object" && item.id === guardrail.id)
        );
      } else {
        return [...prev, guardrail];
      }
    });
  };
  // For Mistral and Pangea in the JSX:
  const isGuardrailSelected = (guardrail) => {
    return selectedGuardrails.some(
      (item) => typeof item === "object" && item.id === guardrail.id
    );
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

        // Set guardrail types and extract provider keys
        setGuardrailTypes(guardrailTypesData);
        const providers = Object.keys(guardrailTypesData);
        setGuardrailProviders(providers);

        console.log("Guardrail Types Data:", guardrailTypesData);

        // You can also console log individual provider types
        console.log("Pangea guardrails:", guardrailTypesData.pangea);
        console.log("Mistral guardrails:", guardrailTypesData.mistral);
        console.log("Bedrock guardrails:", guardrailTypesData.bedrock);
        console.log("Vapus guardrails:", guardrailTypesData.vapus);
      } catch (error) {
        console.error("Failed to fetch enum data:", error);
        toast.error("Failed to fetch configuration data");
      } finally {
        setIsLoading(false);
      }
    };

    fetchEnumsData();
  }, []);

  // fatch api for aiModelNodes dropdown
  useEffect(() => {
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
      } catch (error) {
        console.error("Failed to fetch model nodes:", error);
        setError(error.message);
      } finally {
        setIsLoading(false);
      }
    };

    fetchModelNodes();
  }, []);

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

  const addTopic = () => {
    setTopics([...topics, { topic: "", description: "", samples: [] }]);
    setTopicCount((prevCount) => prevCount + 1);
  };

  const removeTopic = (index) => {
    const newTopics = [...topics];
    newTopics.splice(index, 1);
    setTopics(newTopics);
    setTopicCount((prevCount) => prevCount - 1);
  };

  const addWordEntry = () => {
    setWordEntries([...wordEntries, { words: "", fileLocation: "" }]);
    setWordEntryCount((prevCount) => prevCount + 1);
  };

  const removeWordEntry = (index) => {
    const newEntries = [...wordEntries];
    newEntries.splice(index, 1);
    setWordEntries(newEntries);
    setWordEntryCount((prevCount) => prevCount - 1);
  };

  const addSensitiveEntry = () => {
    setSensitiveEntries([
      ...sensitiveEntries,
      { piiType: "", action: "", regex: "" },
    ]);
    setSensitiveEntryCount((prevCount) => prevCount + 1);
  };

  const removeSensitiveEntry = (index) => {
    const newEntries = [...sensitiveEntries];
    newEntries.splice(index, 1);
    setSensitiveEntries(newEntries);
    setSensitiveEntryCount((prevCount) => prevCount - 1);
  };

  const getFormData = (formData) => {
    const dataObj = {
      spec: {
        displayName: "",
        name: "",
        minConfidence: 0,
        maxConfidence: 0,
        description: "",
        failureMessage: "",
        scanMode: "BIDIRECTIONAL",
        resourceBase: {
          createdAt: "0",
          createdBy: "",
          deletedAt: "0",
          deletedBy: "",
          updatedAt: "0",
          updatedBy: "",
          organization: "",
          account: "",
          status: "",
          owners: [],
          scope: "ORGANIZATION_SCOPE",
          labels: [],
          editors: [],
        },
        contents: {
          hateSpeech: "NONE",
          insults: "NONE",
          sexual: "NONE",
          threats: "NONE",
          misconduct: "NONE",
        },
        topics: [
          {
            topic: "",
            samples: [],
            description: "",
          },
        ],
        words: [
          {
            words: [],
            fileLocation: "",
          },
        ],
        sensitiveDataset: [
          {
            piiType: "",
            action: "",
            regex: "",
          },
        ],
        guardModel: {
          modelNodeId: "",
          modelId: "",
        },
        guardrailId: "",
        schema: "",
        eligibleModelNodes: [],
        partner: [
          {
            bedrock: [],
            mistral: [],
            pangea: [],
          },
        ],
      },
    };

    // Populate with form values
    for (let [key, value] of formData.entries()) {
      if (key.startsWith("spec.")) {
        const path = key.substring(5);
        const parts = path.split(".");

        // Handle array indices in property names, e.g., "topics[0].topic"
        let current = dataObj.spec;
        let processedPath = "";

        for (let i = 0; i < parts.length; i++) {
          const part = parts[i];
          processedPath += part;

          // Handle array notation like topics[0]
          const match = part.match(/([^\[]+)(?:\[(\d+)\])?(.*)$/);
          if (match) {
            const propName = match[1];
            const index = match[2] ? parseInt(match[2]) : null;
            const remaining = match[3];

            if (!current[propName]) {
              if (index !== null) {
                current[propName] = [];
              } else {
                current[propName] = {};
              }
            }

            if (index !== null) {
              while (current[propName].length <= index) {
                if (processedPath.includes("topics")) {
                  current[propName].push({
                    topic: "",
                    samples: [],
                    description: "",
                  });
                } else if (processedPath.includes("words")) {
                  current[propName].push({ words: [], fileLocation: "" });
                } else if (processedPath.includes("sensitiveDataset")) {
                  current[propName].push({
                    piiType: "",
                    action: "",
                    regex: "",
                  });
                } else {
                  current[propName].push({});
                }
              }

              if (i === parts.length - 1) {
                if (remaining) {
                  current[propName][index][remaining.substring(1)] = value;
                } else {
                  current[propName][index] = value;
                }
              } else {
                current = current[propName][index];
              }
            } else {
              if (i === parts.length - 1) {
                current[propName] = value;
              } else {
                current = current[propName];
              }
            }
          }
        }
      }
    }

    return dataObj.spec;
  };

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
          result[key] = parseFloat(obj[key]);
        } else if (typeof obj[key] === "object" && obj[key] !== null) {
          result[key] = convertNumericStrings(obj[key]);
        } else {
          result[key] = obj[key];
        }
      }
    }

    return result;
  };

  const submitCreateForm = async (event) => {
    event.preventDefault();

    try {
      setIsLoading(true);
      const formData = new FormData(event.target);
      let dataObj = getFormData(formData);

      // Process topics samples into arrays
      const topicInputs = document.querySelectorAll(
        '[name^="spec.topics"][name$=".samples"]'
      );

      topicInputs.forEach((input) => {
        const indexMatch = input.name.match(/\[(\d+)\]/);
        if (indexMatch && indexMatch[1]) {
          const index = parseInt(indexMatch[1]);

          const samplesArray = input.value
            ? input.value.split(",").map((item) => item.trim())
            : [];
          dataObj.topics[index].samples = samplesArray;
        }
      });

      // Process words entries into arrays
      const wordEntries = document.querySelectorAll(
        '[name^="spec.words"][name$=".words"]'
      );

      wordEntries.forEach((input) => {
        const indexMatch = input.name.match(/\[(\d+)\]/);
        if (indexMatch && indexMatch[1]) {
          const index = parseInt(indexMatch[1]);

          const wordsArray = input.value
            ? input.value.split(",").map((item) => item.trim())
            : [];
          dataObj.words[index].words = wordsArray;
        }
      });

      dataObj = convertNumericStrings(dataObj);

      // Ensure default values if not set
      if (!dataObj.scanMode) dataObj.scanMode = "BIDIRECTIONAL";
      if (!dataObj.resourceBase || !dataObj.resourceBase.scope) {
        dataObj.resourceBase = dataObj.resourceBase || {};
        dataObj.resourceBase.scope = "ORGANIZATION_SCOPE";
      }

      let partnerData = [];

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
          partnerData.push(partnerObj);
        }
      }

      dataObj.partner = partnerData;
      delete dataObj.guardrailType;

      const payload = { spec: dataObj };

      console.log("Payload:", payload);

      const output = await GuardrailFormApi.getGuardrailForm(payload);

      console.log("Resource created:", output);

      const resourceInfo = output.result;
      if (resourceInfo) {
        toast.success(
          "Guardrail created successfully.",
          `${resourceInfo.resource} Resource created successfully.`
        );
        setTimeout(() => {
          router.push(`/ai-center/guardrails`);
        }, 1000);
      } else {
        toast.success("Resource Created", "Resource created successfully.");
      }
    } catch (error) {
      console.error("Error sending API request:", error);
      toast.error("Failed to create guardrail. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="bg-zinc-800 flex h-screen">
      <Head>
        <title>Create AI Guardrail</title>
        <meta name="description" content="Create a new AI guardrail" />
      </Head>

      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="Create Guardrail"
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
                  onClick={() => setActiveTab("yaml")}
                  className={`whitespace-nowrap border-b-2 py-2 px-2 text-md font-medium focus:outline-none ${
                    activeTab === "yaml"
                      ? "border-orange-700 text-orange-700 font-semibold"
                      : "border-black"
                  }`}
                >
                  YAML
                </button> */}
                <button
                  onClick={() => setActiveTab("form")}
                  className={`whitespace-nowrap border-b-2 py-2 px-2 text-md font-medium focus:outline-none ml-4 ${
                    activeTab === "form"
                      ? "border-orange-700 text-orange-700 font-semibold"
                      : ""
                  }`}
                  suppressHydrationWarning
                >
                  Form
                </button>
              </div>

              {activeTab === "yaml" && (
                <div id="yamlSpec">
                  {/* YAML editor would go here - simplified for this example */}
                  <YamlEditorClient />
                </div>
              )}

              {activeTab === "form" && (
                <div id="formSpec">
                  <form
                    id="dataSourceSpec"
                    className="space-y-2 border border-zinc-500 rounded-md text-gray-100 p-2"
                    onSubmit={submitCreateForm}
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
                            placeholder="Enter name"
                            className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                            required
                            suppressHydrationWarning
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
                                {/* Minimum Confidence */}
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
                                    placeholder="Enter Minimum Confidence"
                                    className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                    required
                                    suppressHydrationWarning
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
                                    placeholder="Enter Maximum Confidence"
                                    className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                    required
                                    suppressHydrationWarning
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
                                    className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-300 p-2 shadow-sm w-full"
                                  />
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
                                    className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-300 p-2 shadow-sm w-full"
                                  />
                                </div>

                                {/* Scan Mode */}
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
                                    className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                    suppressHydrationWarning
                                  >
                                    <option value="">Select Scan Mode</option>
                                    {enums.AIGuardrailScanMode.map(
                                      (scanMode) => (
                                        <option key={scanMode} value={scanMode}>
                                          {strTitle(scanMode)}
                                        </option>
                                      )
                                    )}
                                  </select>
                                </div>

                                {/* Scope */}
                                <div>
                                  <label
                                    htmlFor="spec_scope"
                                    className="labels"
                                  >
                                    Scope
                                  </label>
                                  <select
                                    id="spec_scope"
                                    name="spec.resourceBase.scope"
                                    className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                    suppressHydrationWarning
                                  >
                                    <option value="">Select Scope</option>
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
                                    {/* Hate Speech Level */}
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
                                        className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                        suppressHydrationWarning
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

                                    {/* Insults Level */}
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
                                        className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                        suppressHydrationWarning
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
                                        className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                        suppressHydrationWarning
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
                                        className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                        suppressHydrationWarning
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
                                        className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                        suppressHydrationWarning
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
                                <fieldset className="rounded mb-4">
                                  <div>
                                    {Array.from({ length: topicCount }).map(
                                      (_, index) => (
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
                                            name={`spec.topics[${index}].topic`}
                                            type="text"
                                            placeholder="Enter topic name"
                                            className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full mb-2"
                                            suppressHydrationWarning
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
                                            placeholder="Enter description"
                                            className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-300 p-2 shadow-sm w-full"
                                          />

                                          <label
                                            htmlFor={`spec_Topics_${index}_samples`}
                                            className="labels"
                                          >
                                            Topic Samples
                                          </label>
                                          <input
                                            id={`spec_Topics_${index}_samples`}
                                            name={`spec.topics[${index}].samples`}
                                            type="text"
                                            placeholder="Enter samples (comma separated)"
                                            className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                            suppressHydrationWarning
                                          />
                                          {index > 0 && (
                                            <RemoveButton
                                              onClick={() => removeTopic(index)}
                                            />
                                          )}
                                        </div>
                                      )
                                    )}
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
                                    {Array.from({ length: wordEntryCount }).map(
                                      (_, index) => (
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
                                                placeholder="Enter words (comma separated)"
                                                className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                                suppressHydrationWarning
                                              />
                                            </div>
                                            <div>
                                              <label className="labels">
                                                File Location
                                              </label>
                                              <input
                                                type="text"
                                                name={`spec.words[${index}].fileLocation`}
                                                placeholder="Enter file location"
                                                className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                                suppressHydrationWarning
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
                                      )
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
                                    {sensitiveEntries.map((entry, index) => (
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
                                              name={`spec.sensitiveDataset[${index}].piiType`}
                                              placeholder="Enter PII type"
                                              className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                              suppressHydrationWarning
                                            />
                                          </div>
                                          <div>
                                            <label className="labels">
                                              Action
                                            </label>
                                            <input
                                              type="text"
                                              name={`spec.sensitiveDataset[${index}].action`}
                                              placeholder="Enter action"
                                              className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                              suppressHydrationWarning
                                            />
                                          </div>
                                          <div>
                                            <label className="labels">
                                              Regex
                                            </label>
                                            <input
                                              type="text"
                                              name={`spec.sensitiveDataset[${index}].regex`}
                                              placeholder="Enter regex pattern"
                                              className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                              suppressHydrationWarning
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
                                    ))}
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
                                          value={selectedModelNodeId}
                                          onChange={(e) =>
                                            populateModelDropdown(
                                              e.target.value
                                            )
                                          }
                                          className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                          suppressHydrationWarning
                                        >
                                          <option value="">
                                            Select Model Gateway
                                          </option>
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
                                          className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                          suppressHydrationWarning
                                        >
                                          <option value="">Select Model</option>
                                          {modelOptions.map((model, index) => (
                                            <option
                                              key={`${model.modelId}-${index}`}
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
                                              {guardrail.name}
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
                                              {strTitle(guardrail.name)}
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
                                              {strTitle(guardrail.name)}
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
                          disabled={isLoading}
                          className="px-6 py-2 bg-orange-700 text-white rounded-md shadow hover:bg-pink-900"
                          suppressHydrationWarning
                        >
                          {isLoading ? (
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
