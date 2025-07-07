"use client";
import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import Header from "@/app/components/platform/header";
import ToastContainerMessage from "@/app/components/notification/customToast";
import LoadingOverlay from "@/app/components/loading/loading";
import { platformApi, platformUpdateApi } from "@/app/utils/settings-endpoint/platform-api";
import { enumsApi } from "@/app/utils/developers-endpoint/enums";
import { modelsRegistryApi } from "@/app/utils/ai-studio-endpoint/models-registry-api";

export default function UpdatePlatformForm() {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState("form");
  const [isLoading, setIsLoading] = useState(true);
  const [enums, setEnums] = useState();
  const [aiModelNodes, setAiModelNodes] = useState([]);
  const [modelMap, setModelMap] = useState({});
  const [selectedGenerativeModelNodeId, setSelectedGenerativeModelNodeId] = useState("");
  const [generativeModelOptions, setGenerativeModelOptions] = useState([]);
  const [selectedGuardrailModelNodeId, setSelectedGuardrailModelNodeId] = useState("");
  const [guardrailModelOptions, setGuardrailModelOptions] = useState([]);
  const [selectedEmbeddingModelNodeId, setSelectedEmbeddingModelNodeId] = useState("");
  const [embeddingModelOptions, setEmbeddingModelOptions] = useState([]);
  const [formData, setFormData] = useState({
    // Profile fields
    logo: "",
    favicon: "",
    
    // dmAccessJwtKeys fields
    dmAccessJwtKeys: {
      name: "",
      publicJwtKey: "",
      privateJwtKey: "",
      vId: "",
      signingAlgorithm: "",
      status: "",
    },
    
    // AI Attributes fields
    aiAttributes: {
      embeddingModelNode: "",
      embeddingModel: "",
      generativeModelNode: "",
      generativeModel: "",
      guardrailModelNode: "",
      guardrailModel: ""
    }
  });

  // Store complete profile data from API
  const [completeProfileData, setCompleteProfileData] = useState({
    addresses: [],
    logo: "",
    description: "",
    moto: "",
    favicon: ""
  });
  
  // Fetch enum data for Signing Algorithm
  useEffect(() => {
    const fetchEnumsData = async () => {
      try {
        setIsLoading(true);
        const response = await enumsApi.getEnums();

        const enumResponses = response.enumResponse || [];
        const encryptionAlgoEnum = enumResponses.find(
          (enumData) => enumData.name === "EncryptionAlgo"
        );

        const encryptionAlgos = encryptionAlgoEnum?.value || [];

        // Store only EncryptionAlgo data
        setEnums({ EncryptionAlgo: encryptionAlgos });

        console.log("Encryption algorithms loaded:", encryptionAlgos);
      } catch (error) {
        console.error("Failed to fetch EncryptionAlgo data:", error);
        toast.error("Failed to fetch configuration data");
      } finally {
        setIsLoading(false);
      }
    };

    fetchEnumsData();
  }, []);

  // Fetch existing platform data
  useEffect(() => {
    const fetchPlatformData = async () => {
      try {
        setIsLoading(true);
        const data = await platformApi.getPlatform();
        
        if (data && data.output) {
          const platformData = data.output;
          
          // Store complete profile data
          setCompleteProfileData({
            addresses: platformData.profile?.addresses || [],
            logo: platformData.profile?.logo || "",
            description: platformData.profile?.description || "",
            moto: platformData.profile?.moto || "",
            favicon: platformData.profile?.favicon || ""
          });
          
          setFormData({
            logo: platformData.profile?.logo || "",
            favicon: platformData.profile?.favicon || "",
            dmAccessJwtKeys: {
              name: platformData.dmAccessJwtKeys?.name || "",
              publicJwtKey: platformData.dmAccessJwtKeys?.publicJwtKey || "",
              privateJwtKey: platformData.dmAccessJwtKeys?.privateJwtKey || "",
              vId: platformData.dmAccessJwtKeys?.vId || "",
              signingAlgorithm: platformData.dmAccessJwtKeys?.signingAlgorithm || "",
              status: platformData.dmAccessJwtKeys?.status || "",
            },
            aiAttributes: {
              embeddingModelNode: platformData.aiAttributes?.embeddingModelNode || "",
              embeddingModel: platformData.aiAttributes?.embeddingModel || "",
              generativeModelNode: platformData.aiAttributes?.generativeModelNode || "",
              generativeModel: platformData.aiAttributes?.generativeModel || "",
              guardrailModelNode: platformData.aiAttributes?.guardrailModelNode || "",
              guardrailModel: platformData.aiAttributes?.guardrailModel || ""
            }
          });

          // Pre-populate dropdowns if values exist
          if (platformData.aiAttributes?.generativeModelNode) {
            populateGenerativeModelDropdown(platformData.aiAttributes.generativeModelNode, false);
          }
          if (platformData.aiAttributes?.guardrailModelNode) {
            populateGuardrailModelDropdown(platformData.aiAttributes.guardrailModelNode, false);
          }
          if (platformData.aiAttributes?.embeddingModelNode) {
            populateEmbeddingModelDropdown(platformData.aiAttributes.embeddingModelNode, false);
          }
        }
      } catch (error) {
        console.error("Error fetching platform data:", error);
        toast.error("Failed to fetch platform data");
      } finally {
        setIsLoading(false);
      }
    };

    fetchPlatformData();
  }, []);

  // Fetch model nodes data
  useEffect(() => {
    const fetchModelNodes = async () => {
      try {
        const response = await modelsRegistryApi.getModelsRegistry();
        const nodes = response.output?.aiModelNodes || [];
        console.log("Model nodes:", nodes);
        setAiModelNodes(nodes);     

        const map = {};
        nodes.forEach((item) => {
          map[item.modelNodeId] = item;
        });
        setModelMap(map);     

        return map;
      } catch (error) {
        console.error("Failed to fetch model nodes:", error);
        toast.error(`Failed to load model nodes: ${error.message}`);
        return {};
      }
    };

    fetchModelNodes();
  }, []);

  // Populate available models based on selected node
  const populateGenerativeModelDropdown = (nodeId, resetModel = true) => {
    setSelectedGenerativeModelNodeId(nodeId);
    const selectedNode = aiModelNodes.find(
      (node) => node.modelNodeId === nodeId
    );
    
    console.log("Selected generative node:", selectedNode);
    
    if (selectedNode?.attributes?.generativeModels) {
      setGenerativeModelOptions(selectedNode.attributes.generativeModels);
    } else {
      setGenerativeModelOptions([]);
    }
    
    setFormData(prev => ({
      ...prev,
      aiAttributes: {
        ...prev.aiAttributes,
        generativeModelNode: nodeId,
        generativeModel: resetModel ? "" : prev.aiAttributes.generativeModel
      }
    }));
  };

  const populateGuardrailModelDropdown = (nodeId, resetModel = true) => {
    setSelectedGuardrailModelNodeId(nodeId);
    const selectedNode = aiModelNodes.find(
      (node) => node.modelNodeId === nodeId
    );
    
    console.log("Selected guardrail node:", selectedNode);
    console.log("Node attributes:", selectedNode?.attributes);
    
    // Check multiple possible property names for guardrail models
    let guardrailModels = [];
    if (selectedNode?.attributes) {
      // Try different possible property names
      guardrailModels = selectedNode.attributes.guardrailModels || 
                       selectedNode.attributes.guardRailModels || 
                       selectedNode.attributes.generativeModels ||
                       [];
    }
    
    console.log("Guardrail models found:", guardrailModels);
    setGuardrailModelOptions(guardrailModels);
    
    setFormData(prev => ({
      ...prev,
      aiAttributes: {
        ...prev.aiAttributes,
        guardrailModelNode: nodeId,
        guardrailModel: resetModel ? "" : prev.aiAttributes.guardrailModel
      }
    }));
  };

  const populateEmbeddingModelDropdown = (nodeId, resetModel = true) => {
    setSelectedEmbeddingModelNodeId(nodeId);
    const selectedNode = aiModelNodes.find(
      (node) => node.modelNodeId === nodeId
    );
    
    console.log("Selected embedding node:", selectedNode);
    console.log("Node attributes:", selectedNode?.attributes);
    
    if (selectedNode?.attributes?.embeddingModels) {
      setEmbeddingModelOptions(selectedNode.attributes.embeddingModels);
    } else {
      setEmbeddingModelOptions([]);
    }
    
    console.log("Embedding models found:", selectedNode?.attributes?.embeddingModels || []);
    
    setFormData(prev => ({
      ...prev,
      aiAttributes: {
        ...prev.aiAttributes,
        embeddingModelNode: nodeId,
        embeddingModel: resetModel ? "" : prev.aiAttributes.embeddingModel
      }
    }));
  };

  // Handle input changes
  const handleInputChange = (e) => {
    const { name, value, type, checked } = e.target;
    const nameParts = name.split(".");

    setFormData((prevData) => {
      const newData = { ...prevData };
      let current = newData;

      // Navigate to nested objects
      for (let i = 0; i < nameParts.length - 1; i++) {
        if (!current[nameParts[i]]) {
          current[nameParts[i]] = {};
        }
        current = current[nameParts[i]];
      }

      // Set the value
      const lastKey = nameParts[nameParts.length - 1];
      if (type === "checkbox") {
        current[lastKey] = checked;
      } else {
        current[lastKey] = value;
      }

      return newData;
    });

    // Also update the complete profile data for logo and favicon changes
    if (name === "logo" || name === "favicon") {
      setCompleteProfileData(prev => ({
        ...prev,
        [name]: value
      }));
    }
  };

  // Handle form submission
  const submitUpdateForm = async (e) => {
    e.preventDefault();

    try {
      setIsLoading(true);
      
      // payload with complete profile data
      const payload = {
        actions: "UPDATE_PROFILE",
        spec: {
          profile: completeProfileData, 
          dmAccessJwtKeys: formData.dmAccessJwtKeys,
          aiAttributes: formData.aiAttributes
        }
      };

      console.log("Submitting platform update:", payload);
      
      const response = await platformUpdateApi.getplatformUpdate(payload);
      
      console.log("Platform updated successfully");
      toast.success("Platform updated successfully!");
      
      // Redirect after successful update
      setTimeout(() => {
        router.push("./");
      }, 1000);
      
    } catch (error) {
      console.error("Error updating platform:", error);
      toast.error(error.message || "Failed to update platform");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="bg-zinc-800 flex h-screen">
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="Update Platform"
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

              {activeTab === "form" && (
                <div id="formSpec">
                  <form
                    id="platformUpdateForm"
                    className="space-y-4 border border-zinc-500 rounded-md text-gray-100 p-4"
                    onSubmit={submitUpdateForm}
                  >
                    {/* Profile Section */}
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      {/* Logo */}
                      <div>
                        <label htmlFor="logo" className="labels block text-sm font-medium mb-2">
                          Logo URL
                        </label>
                        <input
                          id="logo"
                          name="logo"
                          type="text"
                          placeholder="Enter logo URL"
                          className="form-input-field w-full p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm"
                          value={formData.logo}
                          onChange={handleInputChange}
                          suppressHydrationWarning
                        />
                      </div>
                      
                      {/* Favicon */}
                      <div>
                        <label htmlFor="favicon" className="labels block text-sm font-medium mb-2">
                          Favicon URL
                        </label>
                        <input
                          id="favicon"
                          name="favicon"
                          type="text"
                          placeholder="Enter favicon URL"
                          className="form-input-field w-full p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm"
                          value={formData.favicon}
                          onChange={handleInputChange}
                          suppressHydrationWarning
                        />
                      </div>
                    </div>
                  
                    {/* JWT Keys Section */}
                    <details className="border border-zinc-500 p-4 rounded mb-4">
                      <summary className="text-lg font-semibold cursor-pointer">
                        JWT Params
                      </summary>
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        {/* JWT Name */}
                        <div>
                          <label htmlFor="dmAccessJwtKeys.name" className="labels block text-sm font-medium mb-2">
                            JWT Key Name
                          </label>
                          <input
                            id="dmAccessJwtKeys.name"
                            name="dmAccessJwtKeys.name"
                            type="text"
                            placeholder="Enter JWT key name"
                            className="form-input-field w-full p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm"
                            value={formData.dmAccessJwtKeys.name}
                            onChange={handleInputChange}
                            suppressHydrationWarning
                          />
                        </div>
                        
                        {/* Signing Algorithm */}
                        <div>
                          <label htmlFor="dmAccessJwtKeys.signingAlgorithm" className="labels block text-sm font-medium mb-2">
                            Signing Algorithm
                          </label>
                          <select
                            id="dmAccessJwtKeys.signingAlgorithm"
                            name="dmAccessJwtKeys.signingAlgorithm"
                            className="mt-1 w-full flex justify-between items-center bg-zinc-700 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700 rounded-md shadow-sm outline-none border-none"
                            value={formData.dmAccessJwtKeys.signingAlgorithm}
                            onChange={handleInputChange}
                            suppressHydrationWarning
                          >
                            <option value="">Select Signing Algorithm</option>
                            {enums?.EncryptionAlgo?.map((algorithm) => (
                              <option key={algorithm} value={algorithm}>
                                {algorithm}
                              </option>
                            ))}
                          </select>
                        </div>
                        
                        {/* Public JWT Key */}
                        <div className="md:col-span-2">
                          <label htmlFor="dmAccessJwtKeys.publicJwtKey" className="labels block text-sm font-medium mb-2">
                            Public JWT Key
                          </label>
                          <textarea
                            id="dmAccessJwtKeys.publicJwtKey"
                            name="dmAccessJwtKeys.publicJwtKey"
                            rows="3"
                            placeholder="Enter public JWT key"
                            className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-300 p-2 shadow-sm w-full"
                            value={formData.dmAccessJwtKeys.publicJwtKey}
                            onChange={handleInputChange}
                            suppressHydrationWarning
                          />
                        </div>
                        
                        {/* Private JWT Key */}
                        <div className="md:col-span-2">
                          <label htmlFor="dmAccessJwtKeys.privateJwtKey" className="labels block text-sm font-medium mb-2">
                            Private JWT Key
                          </label>
                          <textarea
                            id="dmAccessJwtKeys.privateJwtKey"
                            name="dmAccessJwtKeys.privateJwtKey"
                            rows="3"
                            placeholder="Enter private JWT key"
                            className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-300 p-2 shadow-sm w-full"
                            value={formData.dmAccessJwtKeys.privateJwtKey}
                            onChange={handleInputChange}
                            suppressHydrationWarning
                          />
                        </div>
                      </div>
                    </details>
                 
                    {/* Generative AI Params */}
                    <details className="shadow-sm border border-zinc-500 rounded-md shadow-sm p-4 mb-4">
                      <summary className="text-lg font-semibold text-gray-100 cursor-pointer">
                        Generative AI Params
                      </summary>
                      <fieldset className="rounded mb-4">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                          {/* Generative Model Node */}
                          <div>
                            <label htmlFor="aiAttributes_generativeModelNode" className="labels">
                              Generative Model Node
                            </label>
                            <div className="relative">
                              <select
                                id="aiAttributes_generativeModelNode"
                                name="aiAttributes.generativeModelNode"
                                onChange={(e) => populateGenerativeModelDropdown(e.target.value)}
                                value={formData.aiAttributes.generativeModelNode}
                                className="mt-1 w-full flex justify-between overflow-y-auto scrollbar items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                suppressHydrationWarning
                              >
                                <option value="">Select Model Gateway</option>
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
                      
                          {/* Generative Model */}
                          <div>
                            <label htmlFor="aiAttributes_generativeModel" className="labels">
                              Generative Model
                            </label>
                            <div className="relative">
                              <select
                                id="aiAttributes_generativeModel"
                                name="aiAttributes.generativeModel"
                                value={formData.aiAttributes.generativeModel}
                                onChange={handleInputChange}
                                className="mt-1 w-full flex justify-between overflow-y-auto scrollbar items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                suppressHydrationWarning
                              >
                                <option value="">Select Model</option>
                                {generativeModelOptions.map((model) => (
                                  <option
                                    key={model.modelId}
                                    value={model.modelId}
                                    className="text-sm"
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

                    {/* Embedding Generator AI Params */}
                    <details className="shadow-sm border border-zinc-500 rounded-md shadow-sm p-4 mb-4">
                      <summary className="text-lg font-semibold text-gray-100 cursor-pointer">
                        Embedding Generator AI Params
                      </summary>
                      <fieldset className="rounded mb-4">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                          {/* Embedding Model Node */}
                          <div>
                            <label htmlFor="aiAttributes_embeddingModelNode" className="labels">
                              Embedding Model Node
                            </label>
                            <div className="relative">
                              <select
                                id="aiAttributes_embeddingModelNode"
                                name="aiAttributes.embeddingModelNode"
                                onChange={(e) => populateEmbeddingModelDropdown(e.target.value)}
                                value={formData.aiAttributes.embeddingModelNode}
                                className="mt-1 w-full flex justify-between overflow-y-auto scrollbar items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                suppressHydrationWarning
                              >
                                <option value="">Select Model Gateway</option>
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
                      
                          {/* Embedding Model */}
                          <div>
                            <label htmlFor="aiAttributes_embeddingModel" className="labels">
                              Embedding Model
                            </label>
                            <div className="relative">
                              <select
                                id="aiAttributes_embeddingModel"
                                name="aiAttributes.embeddingModel"
                                value={formData.aiAttributes.embeddingModel}
                                onChange={handleInputChange}
                                className="mt-1 w-full flex justify-between overflow-y-auto scrollbar items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                suppressHydrationWarning
                              >
                                <option value="">Select Model</option>
                                {embeddingModelOptions.length > 0 ? (
                                  embeddingModelOptions.map((model) => (
                                    <option
                                      key={model.modelId}
                                      value={model.modelId}
                                      className="text-sm"
                                    >
                                      {model.modelName}
                                    </option>
                                  ))
                                ) : (
                                  <option disabled className="text-gray-500">
                                    {selectedEmbeddingModelNodeId ? "No models available" : "Select a node first"}
                                  </option>
                                )}
                              </select>
                            </div>
                          </div>
                        </div>
                      </fieldset>
                    </details>

                    {/* Guardrail AI Params */}
                    <details className="shadow-sm border border-zinc-500 rounded-md shadow-sm p-4 mb-4">
                      <summary className="text-lg font-semibold text-gray-100 cursor-pointer">
                        Guardrail AI Params
                      </summary>
                      <fieldset className="rounded mb-4">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                          {/* Guardrail Model Node */}
                          <div>
                            <label htmlFor="aiAttributes_guardrailModelNode" className="labels">
                              Guardrail Model Node
                            </label>
                            <div className="relative">
                              <select
                                id="aiAttributes_guardrailModelNode"
                                name="aiAttributes.guardrailModelNode"
                                onChange={(e) => populateGuardrailModelDropdown(e.target.value)}
                                value={formData.aiAttributes.guardrailModelNode}
                                className="mt-1 w-full flex justify-between overflow-y-auto scrollbar items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                suppressHydrationWarning
                              >
                                <option value="">Select Model Gateway</option>
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
                      
                          {/* Guardrail Model */}
                          <div>
                            <label htmlFor="aiAttributes_guardrailModel" className="labels">
                              Guardrail Model
                            </label>
                            <div className="relative">
                              <select
                                id="aiAttributes_guardrailModel"
                                name="aiAttributes.guardrailModel"
                                value={formData.aiAttributes.guardrailModel}
                                onChange={handleInputChange}
                                className="mt-1 w-full flex justify-between overflow-y-auto scrollbar items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                suppressHydrationWarning
                              >
                                <option value="">Select Model</option>
                                {guardrailModelOptions.length > 0 ? (
                                  guardrailModelOptions.map((model) => (
                                    <option
                                      key={model.modelId}
                                      value={model.modelId}
                                      className="text-sm"
                                    >
                                      {model.modelName}
                                    </option>
                                  ))
                                ) : (
                                  <option disabled className="text-gray-500">
                                    {selectedGuardrailModelNodeId ? "No models available" : "Select a node first"}
                                  </option>
                                )}
                              </select>
                            </div>
                          </div>
                        </div>
                      </fieldset>
                    </details>

                    {/* Submit Button */}
                    <div className="mt-6 flex justify-end">
                      <button
                        type="submit"
                        disabled={isLoading}
                        className={`px-6 py-2 bg-orange-700 text-white rounded-md shadow ${
                          isLoading
                            ? "opacity-50 cursor-not-allowed"
                            : "hover:bg-pink-900"
                        }`}
                        suppressHydrationWarning
                      >
                        {isLoading ? "Updating..." : "Update"}
                      </button>
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