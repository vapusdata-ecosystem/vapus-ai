"use client";
import React, { useState, useEffect } from "react";
import { format } from "date-fns";
import Header from "@/app/components/platform/header";
import SectionHeaders from "@/app/components/section-headers";
import { domainApi } from "@/app/utils/settings-endpoint/organization-api";
import AddUserModal from "./adduser/addUser";
import { strTitle } from "@/app/components/JS/common";
import LoadingOverlay from "@/app/components/loading/loading";

export default function OrganizationDetails() {
  const [organizations, setdomains] = useState(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState("basic-info");
  const [modalOpen, setModalOpen] = useState(false);

  // Simplified - we don't need currentDomainId state since we can get it from organizations
  const addUserHandler = () => {
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
        console.error("Error fetching Organization data:", error);
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

  // Get current Organization from the fetched data
  const currentOrganization = organizations?.output?.organizations[0];

  // Get the organizationId directly from the currentOrganization object
  const organizationId = currentOrganization?.organizationId;

  const headerResourceData = currentOrganization ? {
    id: currentOrganization.organizationId,
    name: currentOrganization.name || "Unnamed Organization",
    createdAt: currentOrganization.resourceBase?.createdAt
      ? parseInt(currentOrganization.resourceBase.createdAt) * 1000
      : null,
    createdBy: currentOrganization.resourceBase?.createdBy,
    status: currentOrganization.resourceBase?.status,
    resourceBase: currentOrganization.resourceBase,

    // Create action params for update functionality
    createActionParams: currentOrganization.createActionParams || {
      weblink: `/settings/organization/#`,
    },
  } : null;

  return (
    <div className="bg-zinc-800 flex h-screen relative">
      <LoadingOverlay 
        isLoading={loading} 
        text="Loading Organization details"
        size="default"
        isOverlay={true}
        className="absolute bg-zinc-800 inset-0 z-10"
      />
      
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="Organization Settings"
          hideBackListingLink={true}
          backListingLink="./"
        />

        <div className="flex-grow p-2 w-full text-gray-100">
          {/* Section Headers component */}
          {currentOrganization && (
            <SectionHeaders
              resourceId={organizationId}
              resourceType="Organization"
              fetchUrl="/setting-Organization.json"
              resourceData={headerResourceData}
              customButton={{
                show: true,
                title: "Add Users",
                onClick: () => addUserHandler(organizationId),
              }}
            />
          )}

          {/* Modal component */}
          <AddUserModal
            isOpen={modalOpen}
            onClose={closeModal}
            organizationId={organizationId}
          />

          {currentOrganization && (
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
              </div>

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
                    <p className="break-words p-2">
                      {currentOrganization.displayName}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Organization ID:
                    </p>
                    <p className="break-words p-2">
                      {currentOrganization.organizationId}
                    </p>
                  </div>
                  <div className="lg:flex items-center">
                    <p className="text-base font-extralight text-[#f4d1c2] block">
                      Type:
                    </p>
                    <p className="break-words p-2">
                      {strTitle(currentOrganization.organizationType)}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}