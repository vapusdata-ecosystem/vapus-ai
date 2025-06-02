"use client";
import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import Header from "@/app/components/platform/header";
import ToastContainerMessage from "@/app/components/notification/customToast";
import LoadingOverlay from "@/app/components/loading/loading";
import { platformApi, platformUpdateApi } from "@/app/utils/settings-endpoint/platform-api";

export default function UpdatePlatformForm() {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState("form");
  const [isLoading, setIsLoading] = useState(true);
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
        <LoadingOverlay isLoading={isLoading} />
        
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
                          <input
                            id="dmAccessJwtKeys.signingAlgorithm"
                            name="dmAccessJwtKeys.signingAlgorithm"
                            type="text"
                            placeholder="Enter signing algorithm"
                            className="form-input-field w-full p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm"
                            value={formData.dmAccessJwtKeys.signingAlgorithm}
                            onChange={handleInputChange}
                            suppressHydrationWarning
                          />
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
                     <details className="border border-zinc-500 p-4 rounded mb-4">
                        <summary className="text-lg font-semibold cursor-pointer">
                         Generative AI Params
                        </summary>
                          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                         {/* Generative Model Node */}
                        <div>
                          <label htmlFor="aiAttributes.generativeModelNode" className="labels block text-sm font-medium mb-2">
                            Generative Model Node
                          </label>
                          <input
                            id="aiAttributes.generativeModelNode"
                            name="aiAttributes.generativeModelNode"
                            type="text"
                            placeholder="Enter generative model node"
                            className="form-input-field w-full p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm"
                            value={formData.aiAttributes.generativeModelNode}
                            onChange={handleInputChange}
                            suppressHydrationWarning
                          />
                        </div>
                        {/* Generative Model */}
                        <div>
                          <label htmlFor="aiAttributes.generativeModel" className="labels block text-sm font-medium mb-2">
                            Generative Model
                          </label>
                          <input
                            id="aiAttributes.generativeModel"
                            name="aiAttributes.generativeModel"
                            type="text"
                            placeholder="Enter generative model"
                            className="form-input-field w-full p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm"
                            value={formData.aiAttributes.generativeModel}
                            onChange={handleInputChange}
                            suppressHydrationWarning
                          />
                        </div>
                        </div>
                      </details>

                     {/* Embedding Generator AI Params */}
                         <details className="border border-zinc-500 p-4 rounded mb-4">
                        <summary className="text-lg font-semibold cursor-pointer">
                          Embedding Generator AI Params
                        </summary>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                         {/* Embedding Model Node */}
                        <div>
                          <label htmlFor="aiAttributes.embeddingModelNode" className="labels block text-sm font-medium mb-2">
                            Embedding Model Node
                          </label>
                          <input
                            id="aiAttributes.embeddingModelNode"
                            name="aiAttributes.embeddingModelNode"
                            type="text"
                            placeholder="Enter embedding model node"
                            className="form-input-field w-full p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm"
                            value={formData.aiAttributes.embeddingModelNode}
                            onChange={handleInputChange}
                            suppressHydrationWarning
                          />
                        </div>
                        {/* Embedding Model */}
                        <div>
                          <label htmlFor="aiAttributes.embeddingModel" className="labels block text-sm font-medium mb-2">
                            Embedding Model
                          </label>
                          <input
                            id="aiAttributes.embeddingModel"
                            name="aiAttributes.embeddingModel"
                            type="text"
                            placeholder="Enter embedding model"
                            className="form-input-field w-full p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm"
                            value={formData.aiAttributes.embeddingModel}
                            onChange={handleInputChange}
                            suppressHydrationWarning
                          />
                        </div>
                        </div>
                         </details>

                      {/* Guardrail AI Params */}
                        <details className="border border-zinc-500 p-4 rounded mb-4">
                        <summary className="text-lg font-semibold cursor-pointer">
                         Guardrail AI Params
                        </summary>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                       {/* Guardrail Model Node */}
                        <div>
                          <label htmlFor="aiAttributes.guardrailModelNode" className="labels block text-sm font-medium mb-2">
                            Guardrail Model Node
                          </label>
                          <input
                            id="aiAttributes.guardrailModelNode"
                            name="aiAttributes.guardrailModelNode"
                            type="text"
                            placeholder="Enter guardrail model node"
                            className="form-input-field w-full p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm"
                            value={formData.aiAttributes.guardrailModelNode}
                            onChange={handleInputChange}
                            suppressHydrationWarning
                          />
                        </div>
                        {/* Guardrail Model */}
                        <div>
                          <label htmlFor="aiAttributes.guardrailModel" className="labels block text-sm font-medium mb-2">
                            Guardrail Model
                          </label>
                          <input
                            id="aiAttributes.guardrailModel"
                            name="aiAttributes.guardrailModel"
                            type="text"
                            placeholder="Enter guardrail model"
                            className="form-input-field w-full p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm"
                            value={formData.aiAttributes.guardrailModel}
                            onChange={handleInputChange}
                            suppressHydrationWarning
                          />
                        </div>
                        </div>
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
                        {isLoading ? "Updating..." : "Update "}
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