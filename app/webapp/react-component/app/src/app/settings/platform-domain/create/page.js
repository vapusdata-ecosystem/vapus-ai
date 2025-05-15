"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import Sidebar from "@/app/components/platform/main-sidebar";
import Header from "@/app/components/platform/header";
import YamlEditorClient from "@/app/components/formcomponets/ymal";
import ToastContainerMessage from "@/app/components/notification/customToast";
import LoadingOverlay from "@/app/components/loading/loading";
import { platformCreateApi } from "@/app/utils/settings-endpoint/platform-domain";

export default function CreatePlatform() {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState("form");
  const [isLoading, setIsLoading] = useState(false);
  const [formData, setFormData] = useState({
    name: "",
    displayName: "",
  });

  // Handle input changes
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    const fieldName = name.replace("spec.", "");

    setFormData((prevData) => ({
      ...prevData,
      [fieldName]: value,
    }));
  };

  // Handle form submission
  const submitCreateForm = async (e) => {
    e.preventDefault();

    if (!formData.name.trim() || !formData.displayName.trim()) {
      toast.error("Please fill in all required fields");
      return;
    }

    try {
      setIsLoading(true);
      const payload = {
        spec: {
          name: formData.name,
          displayName: formData.displayName,
        },
      };

      console.log("Submitting platform data:", payload);
      // Call API
      const response = await platformCreateApi.getplatformCreate(payload);

      console.log("Platform created:", response);
      toast.success("Platform created successfully!");
      router.push("./");
    } catch (error) {
      console.error("Error creating platform:", error);
      toast.error(error.message || "Failed to create platform");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="bg-zinc-800 flex h-screen">
      <Sidebar />
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="Create Platform"
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
                      <legend className="text-xl font-bold text-gray-100">
                        Spec
                      </legend>

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
                            value={formData.name}
                            onChange={handleInputChange}
                            suppressHydrationWarning
                          />
                        </div>
                        {/* Display Name */}
                        <div>
                          <label htmlFor="spec_displayName" className="labels">
                            Display Name
                          </label>
                          <input
                            id="spec_displayName"
                            name="spec.displayName"
                            type="text"
                            placeholder="Enter Display name"
                            className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                            value={formData.displayName}
                            onChange={handleInputChange}
                            suppressHydrationWarning
                          />
                        </div>
                      </div>

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
                          {isLoading ? "Creating..." : "Submit"}
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
