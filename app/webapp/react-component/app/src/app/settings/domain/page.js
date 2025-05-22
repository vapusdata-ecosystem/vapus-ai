"use client";
import React, { useState, useEffect } from "react";
import { format } from "date-fns";
import Header from "@/app/components/platform/header";
import SectionHeaders from "@/app/components/section-headers";
import { domainApi } from "@/app/utils/settings-endpoint/domain-api";
import AddUserModal from "./adduser/addUser";

export default function DomainDetails() {
  const [domains, setdomains] = useState(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState("basic-info");
  const [modalOpen, setModalOpen] = useState(false);
  const [currentDomainId, setCurrentDomainId] = useState(null);

  const addUserHandler = (domainId) => {
    setCurrentDomainId(domainId);
    setModalOpen(true);
  };

  const closeModal = () => {
    setModalOpen(false);
  };

  // In your component:
  useEffect(() => {
    const fetchData = async () => {
      try {
        const data = await domainApi.getDomains();
        setdomains(data);
      } catch (error) {
        console.error("Error fetching domain data:", error);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
    alert("Copied to clipboard!");
  };

  // Helper function to convert epoch to formatted date
  const epochConverter = (epoch) => {
    if (epoch === 0) {
      return "--";
    }
    return format(new Date(epoch * 1000), "yyyy-MM-dd");
  };

  // Function to show a specific tab (matching the first example)
  const showTab = (tabId) => {
    setActiveTab(tabId);
  };

  if (loading) {
    return (
      <div className="bg-zinc-800 flex h-screen justify-center items-center">
        <div className="text-white text-xl">Loading domain details...</div>
      </div>
    );
  }

  // Get current domain from the fetched data
  const currentDomain = domains?.output?.domains[0];

  // Get the domainId directly from the currentDomain object
  const domainId = currentDomain?.domainId;

  const headerResourceData = {
    id: currentDomain.domainId,
    name: currentDomain.name || "Unnamed Domain",
    createdAt: currentDomain.resourceBase?.createdAt
      ? parseInt(currentDomain.resourceBase.createdAt) * 1000
      : null,
    createdBy: currentDomain.resourceBase?.createdBy,
    status: currentDomain.resourceBase?.status,
    resourceBase: currentDomain.resourceBase,

    // Create action params for update functionality
    createActionParams: currentDomain.createActionParams || {
      weblink: `/settings/domain/#`,
    },
  };

  return (
    <div className="bg-zinc-800 flex h-screen">
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="Domain Settings"
          hideBackListingLink={true}
          backListingLink="./"
        />

        <div className="flex-grow p-2 w-full text-gray-100">
          {/* Section Headers component */}
          <SectionHeaders
            resourceId={domainId}
            resourceType="domain"
            fetchUrl="/setting-domain.json"
            resourceData={headerResourceData}
            customButton={{
              show: true,
              title: "Add Users",
              onClick: () => addUserHandler(domainId),
            }}
          />

          {/* Modal component */}
          <AddUserModal
            isOpen={modalOpen}
            onClose={closeModal}
            domainId={currentDomainId}
          />

          <div className="overflow-x-auto scrollbar text-gray-100 bg-zinc-800 rounded-lg p-8 shadow-md">
            {/* Tabs - Styled to match first example */}
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
                onClick={() => showTab("networks")}
                className={`px-4 py-2 font-semibold focus:outline-none ${
                  activeTab === "networks"
                    ? "bg-[oklch(0.205_0_0)] text-white rounded-t-[10px]"
                    : ""
                }`}
              >
                Networks
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

            {/* Tab Content - Using the same conditional display logic as first example */}

            {/* The rest of the component remains the same */}
            {/* Basic Info Tab */}
            <div
              id="basic-info"
              className={`mt-2 bg-[#1b1b1b] p-4 ${
                activeTab !== "basic-info" ? "hidden" : ""
              }`}
            >
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div className=" lg:lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Display Name:
                  </p>
                  <p className="break-words p-2">{currentDomain.displayName}</p>
                </div>
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Domain ID:
                  </p>
                  <p className="break-words p-2">{currentDomain.domainId}</p>
                </div>
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Type:
                  </p>
                  <p className="break-words p-2">{currentDomain.domainType}</p>
                </div>
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Catalog:
                  </p>
                  <p className="break-words p-2">
                    {currentDomain.catalogIndex}
                  </p>
                </div>
              </div>
              <br />
              <h3 className="text-xl mb-4 text-[1.25rem] font-bold text-[#f4d1c2] underline">
                Attributes
              </h3>
              {currentDomain.attributes ? (
                <div className="space-y-4">
                  <div className="bg-zinc-700 p-4 rounded-lg shadow-md">
                    <h4 className="text-lg text-gray-100">JWT Params</h4>
                    <div className="flex flex-col sm:flex-row sm:justify-between mt-2">
                      <div>
                        <p className="text-base font-extralight text-[#f4d1c2] block">
                          Secret Name
                        </p>
                        <p className="break-words p-2">
                          {currentDomain.attributes.authnJwtParams
                            ? currentDomain.attributes.authnJwtParams.name
                            : "N/A"}
                        </p>
                      </div>
                      <div>
                        <p className="text-base font-extralight text-[#f4d1c2] block">
                          Signing Algo
                        </p>
                        <p className="break-words p-2">
                          {currentDomain.attributes.authnJwtParams
                            ? currentDomain.attributes.authnJwtParams
                                .signingAlgorithm
                            : "N/A"}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              ) : (
                <p className="text-gray-100">No attributes available.</p>
              )}
            </div>

            {/* Storage Info Tab */}
            <div
              id="storage-info"
              className={`mt-2 bg-[#1b1b1b] p-4 ${
                activeTab !== "storage-info" ? "hidden" : ""
              }`}
            >
              <h3 className="text-xl mb-4 text-[1.25rem] font-bold text-[#f4d1c2] underline">
                Artifact Storage
              </h3>
              {currentDomain.artifactStorage ? (
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Storage Type:
                    </p>
                    <p className="break-words p-2">
                      {currentDomain.artifactStorage.besType}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Storage Service:
                    </p>
                    <p className="break-words p-2">
                      {currentDomain.artifactStorage.besService}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Onboarding Type:
                    </p>
                    <p className="break-words p-2">
                      {currentDomain.artifactStorage.besOnboarding}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Storage Engine:
                    </p>
                    <p className="break-words p-2">
                      {currentDomain.artifactStorage.besEngine}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Address:
                    </p>
                    <p className="break-words p-2">
                      {currentDomain.artifactStorage.netParams.address}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Status:
                    </p>
                    <p className="break-words p-2">
                      <span
                        className={`px-3 py-1 text-sm font-medium ${
                          currentDomain.artifactStorage.status === "ACTIVE"
                            ? "text-green-800 bg-green-100"
                            : "text-red-800 bg-red-100"
                        } rounded-full`}
                      >
                        {currentDomain.artifactStorage.status}
                      </span>
                    </p>
                  </div>
                </div>
              ) : (
                <p className="text-gray-100">
                  No storage information available.
                </p>
              )}

              <h3 className="text-xl mb-4 text-[1.25rem] font-bold text-[#f4d1c2] underline">
                Catalog Details
              </h3>
              {currentDomain.dataCatalog &&
                currentDomain.dataCatalog.map((catalog, index) => (
                  <div
                    key={index}
                    className="grid grid-cols-1 sm:grid-cols-2 gap-4"
                  >
                    <div className="lg:flex items-center">
                      <p className="text-base font-extralight text-[#f4d1c2] block">
                        Name:
                      </p>
                      <p className="break-words p-2">{catalog.name}</p>
                    </div>
                    <div className="lg:flex items-center">
                      <p className="text-base font-extralight text-[#f4d1c2] block">
                        Display Name:
                      </p>
                      <p className="break-words p-2">{catalog.displayName}</p>
                    </div>
                    <div className="lg:flex items-center">
                      <p className="text-base font-extralight text-[#f4d1c2] block">
                        Id:
                      </p>
                      <p className="break-words p-2">{catalog.dataCatalogId}</p>
                    </div>
                    <div className="lg:flex items-center">
                      <p className="text-base font-extralight text-[#f4d1c2] block">
                        Description:
                      </p>
                      <p className="break-words p-2">{catalog.description}</p>
                    </div>
                  </div>
                ))}
            </div>

            {/* Operating Systems Tab */}
            <div
              id="base-os"
              className={`mt-2 bg-[#1b1b1b] p-4 ${
                activeTab !== "base-os" ? "hidden" : ""
              }`}
            >
              <div className="space-y-4">
                {currentDomain.domainArtifacts &&
                currentDomain.domainArtifacts.length > 0 ? (
                  currentDomain.domainArtifacts.map((artifact, index) => (
                    <div
                      key={index}
                      className="bg-zinc-700 p-4 rounded-lg shadow-md"
                    >
                      <h4 className="text-lg text-gray-100">
                        Type: {artifact.artifactType}
                      </h4>
                      <div className="flex flex-col sm:flex-row sm:justify-between mt-2 grid grid-cols-2 gap-6">
                        {artifact.artifacts &&
                          artifact.artifacts.map((item, idx) => (
                            <React.Fragment key={idx}>
                              <div>
                                <p className="text-base font-extralight text-[#f4d1c2] block">
                                  Url
                                </p>
                                <div className="lg:flex items-center">
                                  <p
                                    className="text-gray-100"
                                    style={{
                                      maxWidth: "500px",
                                      overflowWrap: "break-word",
                                    }}
                                  >
                                    {item.artifact}
                                    <button
                                      onClick={() =>
                                        copyToClipboard(item.artifact)
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
                                    className="text-gray-100"
                                    style={{
                                      maxWidth: "500px",
                                      overflowWrap: "break-word",
                                    }}
                                  >
                                    {item.digest}
                                    <button
                                      onClick={() =>
                                        copyToClipboard(item.digest)
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
                                  IsLatest
                                </p>
                                <p className="break-words p-2">
                                  {item.isLatest ? "Yes" : "No"}
                                </p>
                              </div>
                              <div>
                                <p className="text-base font-extralight text-[#f4d1c2] block">
                                  AddedOn
                                </p>
                                <p className="break-words p-2">
                                  {epochConverter(item.addedOn)}
                                </p>
                              </div>
                            </React.Fragment>
                          ))}
                      </div>
                    </div>
                  ))
                ) : (
                  <p className="text-gray-100">
                    No domain artifacts available.
                  </p>
                )}
              </div>
            </div>

            {/* Networks Tab */}
            <div
              id="networks"
              className={`mt-2 bg-[#1b1b1b] p-4 ${
                activeTab !== "networks" ? "hidden" : ""
              }`}
            >
              <h3 className="text-xl mb-4 text-[1.25rem] font-bold text-[#f4d1c2] underline">
                Kubernetes Infra
              </h3>
              <div className="space-y-4">
                {currentDomain.dataProductInfraPlatform &&
                currentDomain.dataProductInfraPlatform.length > 0 ? (
                  currentDomain.dataProductInfraPlatform.map((infra, index) => (
                    <div
                      key={index}
                      className="bg-zinc-700 p-4 rounded-lg shadow-md"
                    >
                      <h4 className="text-lg text-gray-100">
                        Type: {infra.name}
                      </h4>
                      <div className="flex flex-col sm:flex-row sm:justify-between mt-2 grid grid-cols-2 gap-6">
                        <div>
                          <p className="text-base font-extralight text-[#f4d1c2] block">
                            Name
                          </p>
                          <p className="break-words p-2">{infra.name}</p>
                        </div>
                        <div>
                          <p className="text-base font-extralight text-[#f4d1c2] block">
                            Id
                          </p>
                          <p className="break-words p-2">{infra.infraId}</p>
                        </div>
                        <div>
                          <p className="text-base font-extralight text-[#f4d1c2] block">
                            Service
                          </p>
                          <p className="break-words p-2">
                            {infra.infraService}
                          </p>
                        </div>
                        <div>
                          <p className="text-base font-extralight text-[#f4d1c2] block">
                            Service Provider
                          </p>
                          <p className="break-words p-2">
                            {infra.serviceProvider}
                          </p>
                        </div>
                        <div>
                          <p className="text-base font-extralight text-[#f4d1c2] block">
                            Secret Name
                          </p>
                          <p className="break-words p-2">{infra.secretName}</p>
                        </div>
                        <div>
                          <p className="text-base font-extralight text-[#f4d1c2] block">
                            Is Default
                          </p>
                          <p className="break-words p-2">
                            {infra.isDefault ? "Yes" : "No"}
                          </p>
                        </div>
                      </div>
                    </div>
                  ))
                ) : (
                  <p className="text-gray-100">
                    No infrastructure information available.
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
