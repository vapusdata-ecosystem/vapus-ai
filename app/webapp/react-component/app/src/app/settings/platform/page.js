"use client";
import React, { useState, useEffect } from "react";
import Header from "@/app/components/platform/header";
import { format } from "date-fns";
import { platformApi } from "@/app/utils/settings-endpoint/platform-api";
import ActionDropdown from "@/app/components/action-dropdown";

export default function PlatformSettings() {
  const [activeTab, setActiveTab] = useState("basic-info");
  const [platformData, setPlatformData] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const data = await platformApi.getPlatform();
        console.log("platformData", data);

        setPlatformData(data);
      } catch (error) {
        console.error("Error fetching platform data:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const showTab = (tabId) => {
    setActiveTab(tabId);
  };

  const copyToClipboard = (text) => {
    navigator.clipboard
      .writeText(text)
      .then(() => {
        alert("Copied to clipboard!");
      })
      .catch((err) => {
        console.error("Failed to copy: ", err);
      });
  };

  // Helper function to convert epoch time to readable date
  const epochConverter = (epoch) => {
    if (epoch === 0) {
      return "--";
    }
    return format(new Date(epoch * 1000), "yyyy-MM-dd");
  };

  // Helper function to limit string length
  const limitLetters = (str, limit) => {
    if (!str) return "";
    return str.length > limit ? str.substring(0, limit) + "..." : str;
  };

  if (loading) {
    return (
      <div className="bg-zinc-800 flex h-screen justify-center items-center">
        <div className="text-white text-xl">Loading platform details...</div>
      </div>
    );
  }

  if (!platformData) {
    return <div className="text-red-500">Error loading platform data</div>;
  }

  const { output: account } = platformData;

  // Create header resource data structure to match DomainDetails

  const responseData = {
    resourceId: "resource-123",
    createActionParams: account.createActionParams || {
      weblink: "./platform/update",
    },
    yamlSpec: account.yamlSpec || JSON.stringify(account, null, 2),
  };

  const globalContextData = {
    AccessTokenKey: "your-access-token-here",
  };

  return (
    <div className="bg-zinc-800 flex h-screen">
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="Platform Settings"
          hideBackListingLink={true}
          backListingLink="./"
        />

        <div className="flex-grow p-2 w-full text-gray-100">
          <div className="overflow-x-auto scrollbar text-gray-100 bg-zinc-800 rounded-lg p-8 shadow-md">
            <div className="flex justify-end">
              <ActionDropdown
                response={responseData}
                globalContext={globalContextData}
              />
            </div>

            <div className="flex border-b border-zinc-500">
              <button
                onClick={() => showTab("basic-info")}
                className={`px-4 py-2 font-semibold focus:outline-none ${
                  activeTab === "basic-info"
                    ? "bg-[oklch(0.205_0_0)] text-white rounded-t-[10px]"
                    : ""
                }`}
              >
                Basic Info
              </button>
              <button
                onClick={() => showTab("storage-info")}
                className={`px-4 py-2 font-semibold focus:outline-none ${
                  activeTab === "storage-info"
                    ? "bg-[oklch(0.205_0_0)] text-white rounded-t-[10px]"
                    : ""
                }`}
              >
                Storages
              </button>
              <button
                onClick={() => showTab("base-os")}
                className={`px-4 py-2 font-semibold focus:outline-none ${
                  activeTab === "base-os"
                    ? "bg-[oklch(0.205_0_0)] text-white rounded-t-[10px]"
                    : ""
                }`}
              >
                Operating Systems
              </button>
            </div>

            {/* Basic Info Tab */}
            <div
              id="basic-info"
              className={`mt-2 bg-[#1b1b1b] p-4 ${
                activeTab !== "basic-info" ? "hidden" : ""
              }`}
            >
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Domain ID:
                  </p>
                  <p className="break-words p-2">{account.accountId}</p>
                </div>
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Type:
                  </p>
                  <p className="break-words p-2">{account.name}</p>
                </div>
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Authentication Method:
                  </p>
                  <p className="break-words p-2">{account.authnMethod}</p>
                </div>
                {account.profile && (
                  <>
                    <div className="lg:flex items-center">
                      <p className="text-base font-extralight text-[#f4d1c2] block">
                        Logo:
                      </p>
                      <p
                        className="break-words p-2 cursor-pointer"
                        onClick={() => copyToClipboard(account.profile.logo)}
                      >
                        {limitLetters(account.profile.logo, 50)}
                      </p>
                    </div>
                    <div className="lg:flex items-center">
                      <p className="text-base font-extralight text-[#f4d1c2] block">
                        Favicon:
                      </p>
                      <p
                        className="break-words p-2 cursor-pointer"
                        onClick={() => copyToClipboard(account.profile.favicon)}
                      >
                        {limitLetters(account.profile.favicon, 50)}
                      </p>
                    </div>
                  </>
                )}
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Description:
                  </p>
                  <p className="break-words p-2">
                    {account.profile &&
                      limitLetters(account.profile.description, 70)}
                  </p>
                </div>
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Status:
                  </p>
                  <p className="break-words p-2">
                    <span
                      className={`px-3 py-1 text-sm font-medium ${
                        account.status === "ACTIVE"
                          ? "text-green-800 bg-green-100"
                          : "text-red-800 bg-red-100"
                      } rounded-full`}
                    >
                      {account.status}
                    </span>
                  </p>
                </div>

                <div className="flex flex-col mt-2">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    User:
                  </p>
                  <ul className="list-disc ml-5 p-2">
                    {account.users &&
                      account.users.map((user, index) => (
                        <li key={index} className="break-words">
                          {user}
                        </li>
                      ))}
                  </ul>
                </div>
              </div>
              <br />
              <h3 className="text-xl mb-4 text-[1.25rem] font-bold text-[#f4d1c2] underline">
                JWT Params
              </h3>
              <div className="space-y-4">
                <div className="bg-zinc-700 p-4 rounded-lg shadow-md">
                  <div className="flex flex-col sm:flex-row sm:justify-between mt-2">
                    <div>
                      <p className="text-base font-extralight text-[#f4d1c2] block">
                        Secret Name
                      </p>
                      <p className="break-words p-2">
                        {account.dmAccessJwtKeys
                          ? account.dmAccessJwtKeys.name
                          : "N/A"}
                      </p>
                    </div>
                    <div>
                      <p className="text-base font-extralight text-[#f4d1c2] block">
                        Signing Algo
                      </p>
                      <p className="break-words p-2">
                        {account.dmAccessJwtKeys
                          ? account.dmAccessJwtKeys.signingAlgorithm
                          : "N/A"}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
              <h3 className="text-xl mb-4 text-[1.25rem] font-bold text-[#f4d1c2] underline">
                Generative AI Params
              </h3>
              <div className="space-y-4">
                <div className="bg-zinc-700 p-4 rounded-lg shadow-md">
                  <div className="flex flex-col sm:flex-row sm:justify-between mt-2">
                    <div>
                      <p className="text-base font-extralight text-[#f4d1c2] block">
                        Generative Model
                      </p>
                      <p className="break-words p-2">
                        {account.aiAttributes.generativeModel}
                      </p>
                    </div>
                    <div>
                      <p className="text-base font-extralight text-[#f4d1c2] block">
                        Generative Model Node
                      </p>
                      <p className="break-words p-2">
                        {account.aiAttributes.generativeModelNode}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
              <h3 className="text-xl mb-4 text-[1.25rem] font-bold text-[#f4d1c2] underline">
                Embedding Generator AI Params
              </h3>
              <div className="space-y-4">
                <div className="bg-zinc-700 p-4 rounded-lg shadow-md">
                  <div className="flex flex-col sm:flex-row sm:justify-between mt-2">
                    <div>
                      <p className="text-base font-extralight text-[#f4d1c2] block">
                        Embedding Model
                      </p>
                      <p className="break-words p-2">
                        {account.aiAttributes.embeddingModel}
                      </p>
                    </div>
                    <div>
                      <p className="text-base font-extralight text-[#f4d1c2] block">
                        Embedding Model Node
                      </p>
                      <p className="break-words p-2">
                        {account.aiAttributes.embeddingModelNode}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
              <h3 className="text-xl mb-4 text-[1.25rem] font-bold text-[#f4d1c2] underline">
                Guardrail AI Params
              </h3>
              <div className="space-y-4">
                <div className="bg-zinc-700 p-4 rounded-lg shadow-md">
                  <div className="flex flex-col sm:flex-row sm:justify-between mt-2">
                    <div>
                      <p className="text-base font-extralight text-[#f4d1c2] block">
                        Guardrail Model
                      </p>
                      <p className="break-words p-2">
                        {account.aiAttributes.guardrailModel}
                      </p>
                    </div>
                    <div>
                      <p className="text-base font-extralight text-[#f4d1c2] block">
                        Guardrail Model Node
                      </p>
                      <p className="break-words p-2">
                        {account.aiAttributes.guardrailModelNode}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Storage Info Tab */}
            <div
              id="storage-info"
              className={`mt-2 bg-[#1b1b1b] p-4 text-sm ${
                activeTab !== "storage-info" ? "hidden" : ""
              }`}
            >
              <h3 className="text-xl mb-2 font-bold text-[#f4d1c2] underline">
                Secret Store
              </h3>
              {account.backendSecretStorage ? (
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Storage Type:
                    </p>
                    <p className="break-words p-2">
                      {account.backendSecretStorage.besType}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Storage Service:
                    </p>
                    <p className="break-words p-2">
                      {account.backendSecretStorage.besService}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Onboarding Type:
                    </p>
                    <p className="break-words p-2">
                      {account.backendSecretStorage.besOnboarding}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Storage Engine:
                    </p>
                    <p className="break-words p-2">
                      {account.backendSecretStorage.besEngine}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Address:
                    </p>
                    <p className="break-words p-2">
                      {account.backendSecretStorage.netParams.address}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Status:
                    </p>
                    <p className="break-words p-2">
                      <span
                        className={`px-3 py-1 text-sm font-medium ${
                          account.backendSecretStorage.status === "ACTIVE"
                            ? "text-green-800 bg-green-100"
                            : "text-red-800 bg-red-100"
                        } rounded-full`}
                      >
                        {account.backendSecretStorage.status}
                      </span>
                    </p>
                  </div>
                </div>
              ) : (
                <p>No secret storage information available.</p>
              )}

              <h3 className="text-xl mb-2 mt-6 text-[1.25rem] font-bold text-[#f4d1c2] underline">
                Backend Storage
              </h3>
              {account.backendDataStorage ? (
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Storage Type:
                    </p>
                    <p className="break-words p-2">
                      {account.backendDataStorage.besType}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Storage Service:
                    </p>
                    <p className="break-words p-2">
                      {account.backendDataStorage.besService}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Onboarding Type:
                    </p>
                    <p className="break-words p-2">
                      {account.backendDataStorage.besOnboarding}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Storage Engine:
                    </p>
                    <p className="break-words p-2">
                      {account.backendDataStorage.besEngine}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Address:
                    </p>
                    <p className="break-words p-2">
                      {account.backendDataStorage.netParams.address}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Status:
                    </p>
                    <p className="break-words p-2">
                      <span
                        className={`px-3 py-1 text-sm font-medium ${
                          account.backendDataStorage.status === "ACTIVE"
                            ? "text-green-800 bg-green-100"
                            : "text-red-800 bg-red-100"
                        } rounded-full`}
                      >
                        {account.backendDataStorage.status}
                      </span>
                    </p>
                  </div>
                </div>
              ) : (
                <p>No backend storage information available.</p>
              )}

              <h3 className="text-xl mb-2 mt-6 text-[1.25rem] font-bold text-[#f4d1c2] underline">
                Artifact Storage
              </h3>
              {account.artifactStorage ? (
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Storage Type:
                    </p>
                    <p className="break-words p-2">
                      {account.artifactStorage.besType}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Storage Service:
                    </p>
                    <p className="break-words p-2">
                      {account.artifactStorage.besService}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Onboarding Type:
                    </p>
                    <p className="break-words p-2">
                      {account.artifactStorage.besOnboarding}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Storage Engine:
                    </p>
                    <p className="break-words p-2">
                      {account.artifactStorage.besEngine}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Address:
                    </p>
                    <p className="break-words p-2">
                      {account.artifactStorage.netParams.address}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Status:
                    </p>
                    <p className="break-words p-2">
                      <span
                        className={`px-3 py-1 text-sm font-medium ${
                          account.artifactStorage.status === "ACTIVE"
                            ? "text-green-800 bg-green-100"
                            : "text-red-800 bg-red-100"
                        } rounded-full`}
                      >
                        {account.artifactStorage.status}
                      </span>
                    </p>
                  </div>
                </div>
              ) : (
                <p>No artifact storage information available.</p>
              )}
            </div>

            {/* Base OS Tab */}
            <div
              id="base-os"
              className={`mt-2 bg-[#1b1b1b] p-4 overflow-y-auto scrollbar max-h-[700px]
  ${activeTab !== "base-os" ? "hidden" : ""}`}
            >
              <div className="space-y-4">
                {account.baseOsArtifacts &&
                account.baseOsArtifacts.length > 0 ? (
                  account.baseOsArtifacts.map((obj, index) => (
                    <div
                      key={index}
                      className="bg-zinc-700 p-4 rounded-lg shadow-md"
                    >
                      <h4 className="text-lg  underline text-[#f4d1c2]">
                        ArtifactType: {obj.artifactType}
                      </h4>
                      <div className="flex flex-col sm:flex-row sm:justify-between mt-2 grid grid-cols-2 gap-6">
                        {obj.artifacts.map((val, idx) => (
                          <React.Fragment key={idx}>
                            <div>
                              <p className="text-base font-extralight text-[#f4d1c2] block">
                                Url
                              </p>
                              <div className="lg:flex items-center">
                                <p
                                  style={{
                                    maxWidth: "500px",
                                    overflowWrap: "break-word",
                                  }}
                                  className="text-gray-100"
                                >
                                  {val.artifact}
                                  <button
                                    onClick={() =>
                                      copyToClipboard(val.artifact)
                                    }
                                    title="Copy"
                                    className="ml-2"
                                  >
                                    <svg
                                      xmlns="http://www.w3.org/2000/svg"
                                      fill="currentColor"
                                      className="w-5 h-5"
                                      viewBox="0 0 24 24"
                                    >
                                      <path d="M13 3H7a2 2 0 0 0-2 2v10h2V5h6V3zm4 4H11a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2zm0 12H11V9h6v10z" />
                                    </svg>
                                  </button>
                                </p>
                              </div>
                            </div>
                            <div>
                              <p className="text-base font-extralight text-[#f4d1c2] block">
                                Digest
                              </p>
                              <div className="lg:flex items-center">
                                <p
                                  style={{
                                    maxWidth: "500px",
                                    overflowWrap: "break-word",
                                  }}
                                  className="text-gray-100"
                                >
                                  {val.digest}
                                  <button
                                    onClick={() => copyToClipboard(val.digest)}
                                    title="Copy"
                                    className="ml-2"
                                  >
                                    <svg
                                      xmlns="http://www.w3.org/2000/svg"
                                      fill="currentColor"
                                      className="w-5 h-5"
                                      viewBox="0 0 24 24"
                                    >
                                      <path d="M13 3H7a2 2 0 0 0-2 2v10h2V5h6V3zm4 4H11a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2zm0 12H11V9h6v10z" />
                                    </svg>
                                  </button>
                                </p>
                              </div>
                            </div>
                            <div>
                              <p className="text-base font-extralight text-[#f4d1c2] block">
                                IsLatest
                              </p>
                              <p className="break-words p-2">
                                {val.isLatest ? "Yes" : "No"}
                              </p>
                            </div>
                            <div>
                              <p className="text-base font-extralight text-[#f4d1c2] block">
                                AddedOn
                              </p>
                              <p className="break-words p-2">
                                {epochConverter(val.addedOn)}
                              </p>
                            </div>
                          </React.Fragment>
                        ))}
                      </div>
                    </div>
                  ))
                ) : (
                  <p className="text-gray-100">
                    No operating system artifacts available.
                  </p>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
