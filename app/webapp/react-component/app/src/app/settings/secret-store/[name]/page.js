"use client";

import { useState, useEffect } from "react";
import { use } from "react";
import Header from "@/app/components/platform/header";
import SectionHeaders from "@/app/components/section-headers";
import {
  secretStoreApi,
  secretStoreArchiveApi,
} from "@/app/utils/settings-endpoint/secret-store-api";

export default function SecretDetailsPage({ params }) {
  console.log("my params", params);
  const unwrappedParams = use(params);
  const secretName = unwrappedParams?.name
    ? String(unwrappedParams.name).trim()
    : "";

  const [secretDetails, setSecretDetails] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [activeTab, setActiveTab] = useState("spec");

  const epochConverterFull = (epochTime) => {
    if (!epochTime || epochTime === "0") return "N/A";
    return new Date(parseInt(epochTime) * 1000).toLocaleString();
  };

  const stringCheck = (str) => {
    return str && str !== "" ? str : "N/A";
  };

  useEffect(() => {
    const fetchSecretDetails = async () => {
      try {
        const response = await secretStoreApi.getSecretStoreName(secretName);
        console.log("secretDetails", response);
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
          setSecretDetails({
            SecretStore: response.output[0],
          });
        } else if (response.output && !Array.isArray(response.output)) {
          setSecretDetails({
            SecretStore: response.output,
          });
        } else {
          console.error(
            "Data does not contain expected output format:",
            response
          );
          setError("Unexpected data format received from server");
        }

        setLoading(false);
      } catch (err) {
        console.error("Error fetching secret details:", err);
        setError(err.message);
        setLoading(false);
      }
    };

    if (secretName) {
      fetchSecretDetails();
    } else {
      setError("No secret name provided");
      setLoading(false);
    }
  }, [secretName]);

  const showTab = (tabId) => {
    setActiveTab(tabId);
  };

  if (loading) {
    return (
      <div className="bg-zinc-800 flex h-screen justify-center items-center">
        <div className="text-white text-xl">Loading secret details...</div>
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

  if (!secretDetails) {
    return (
      <div className="bg-zinc-800 flex h-screen justify-center items-center">
        <div className="text-white text-xl">Secret not found</div>
      </div>
    );
  }

  const secret = secretDetails.SecretStore;
  // Create a data object for SectionHeaders
  const apiServices = {
    secret: {
      archive: secretStoreArchiveApi.getSecretStoreArchive,
      delete: secretStoreArchiveApi.getSecretStoreArchive,
    },
  };
  const headerResourceData = {
    id: secret.name,
    name: secret.name || "Unnamed Secret",
    createdAt: secret.resourceBase?.createdAt
      ? parseInt(secret.resourceBase.createdAt) * 1000
      : null,
    createdBy: secret.resourceBase?.createdBy,
    status: secret.resourceBase?.status,
    resourceBase: secret.resourceBase,
    resourceType: "secret",

    createActionParams: secret.createActionParams || {
      weblink: `./${secret.name}/update`,
    },
    // Add YAML spec for download button
    yamlSpec: secret.yamlSpec || JSON.stringify(secret, null, 2),
  };

  return (
    <div className="bg-zinc-800 flex h-screen">
      <div className="overflow-y-auto h-screen w-full">
        <Header
          sectionHeader="Secret  Details"
          hideBackListingLink={false}
          backListingLink="./"
        />

        <div className="flex-grow p-2 w-full text-gray-100">
          <SectionHeaders
            resourceId={secret.name}
            resourceType="secret"
            resourceData={headerResourceData}
            apiServices={apiServices}
          />

          {/* Tabs */}
          <div className="overflow-x-auto scrollbar text-gray-100 bg-zinc-800 rounded-lg p-8 shadow-md">
            <div className="flex border-b border-zinc-500 font-semibold text-gray-50">
              <button
                onClick={() => showTab("spec")}
                className={`px-4 py-2 focus:outline-none ${
                  activeTab === "spec"
                    ? "bg-[oklch(0.205_0_0)] text-white rounded-t-[10px]"
                    : ""
                }`}
              >
                Spec
              </button>
            </div>

            {/* Tab Content */}
            <div
              id="spec"
              className={`mt-2 bg-[#1b1b1b] p-4 ${
                activeTab !== "spec" ? "hidden" : ""
              }`}
            >
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
                {/* Name */}
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Name:
                  </p>
                  <p className="p-2">{stringCheck(secret.name)}</p>
                </div>

                {/* Secret Type */}
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Secret Type:
                  </p>
                  <p className="p-2">{secret.secretType}</p>
                </div>

                {/* Provider */}
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Provider:
                  </p>
                  <p className="p-2">{secret.provider}</p>
                </div>

                {/* Organization */}
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Organization:
                  </p>
                  <p className="p-2">{secret.resourceBase?.organization}</p>
                </div>

                {/* Data */}
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Data:
                  </p>
                  <p className="p-2">{secret.data || "N/A"}</p>
                </div>

                {/* Description */}
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Description:
                  </p>
                  <p className="p-2">{secret.description}</p>
                </div>

                {/* ExpireAt */}
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Expire At:
                  </p>
                  <p className="p-2">{epochConverterFull(secret.expireAt)}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
