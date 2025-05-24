"use client";
import React, { useState, useEffect } from "react";
import Header from "@/app/components/platform/header";
import { userGlobalData } from "@/context/GlobalContext";
import { userProfileApi } from "../utils/settings-endpoint/profile-api";
import ActionDropdown from "../components/action-dropdown";

const UserDetails = () => {
  const [userData, setUserData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [domainMap, setDomainMap] = useState({});
  const [contextData, setContextData] = useState(null);
  useEffect(() => {
    const fetchData = async () => {
      try {
        // Get user data from global context
        const globalContext = await userGlobalData();
        setContextData(globalContext);
        console.log("my data", globalContext);

        // Check if userId exists
        if (globalContext?.userInfo?.userId) {
          const userId = globalContext.userInfo.userId;
          console.log("User ID:", userId);

          // Make API call to get user profile with userId
          const data = await userProfileApi.getuserProfile(userId);
          console.log("data", data);
          setUserData(data.output.users[0]);
          setDomainMap(data.domainMap);
        } else {
          console.error("User ID not found in global context");
        }
      } catch (error) {
        console.error("Error fetching user data:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const formatEpochTime = (epoch) => {
    if (!epoch || epoch === "0") return "N/A";
    return new Date(epoch * 1000).toLocaleString();
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-screen bg-zinc-800">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-white"></div>
      </div>
    );
  }

  // Action button for update
  const responseData = {
    resourceId: "resource-123",
    createActionParams: userData.createActionParams || {
      weblink: "./settings/update",
    },
  };

  return (
    <div className="bg-zinc-800 flex h-screen">
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="Your profile"
          hideBackListingLink={false}
          backListingLink="./"
        />
        <div className="flex-grow p-2 w-full text-gray-100">
          {loading && (
            <div className="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50">
              <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-white"></div>
            </div>
          )}

          {/* Action Dropdown */}
          <div className="flex justify-end">
            <ActionDropdown response={responseData} />
          </div>

          {/* User Information */}
          <div className="overflow-x-auto text-gray-100 bg-zinc-800 rounded-lg p-8 shadow-md">
            <div className="tab-content mt-2 bg-[#1b1b1b] rounded-lg p-4">
              {/* User Overview Section */}
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-6 ">
                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    User ID:
                  </p>
                  <p className="break-words p-2">{userData?.userId || "N/A"}</p>
                </div>
                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Email:
                  </p>
                  <p className="break-words p-2">{userData?.email || "N/A"}</p>
                </div>
                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Display Name:
                  </p>
                  <p className="break-words p-2">
                    {userData?.displayName || "N/A"}
                  </p>
                </div>
                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Status
                  </p>
                  <p className="break-words p-2">
                    <span
                      className={`px-3 py-1 text-sm font-medium rounded-full ${
                        userData?.resourceBase?.status === "ACTIVE"
                          ? "text-green-800 bg-green-100"
                          : "text-red-800 bg-red-100"
                      }`}
                    >
                      {userData?.resourceBase?.status || "N/A"}
                    </span>
                  </p>
                </div>
                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    First Name:
                  </p>
                  <p className="break-words p-2">
                    {userData?.firstName || "N/A"}
                  </p>
                </div>
                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Last Name:
                  </p>
                  <p className="break-words p-2">
                    {userData?.lastName || "N/A"}
                  </p>
                </div>
                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Onboarding Type:
                  </p>
                  <p className="break-words p-2">
                    {userData?.invitedType || "N/A"}
                  </p>
                </div>
                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Platform Role:
                  </p>
                  <p className="break-words p-2">
                    {userData?.platformRoles?.map((role, index) => (
                      <span
                        key={index}
                        className="px-3 py-1 text-sm font-medium text-gray-800 bg-gray-200 rounded-full mr-2"
                      >
                        {role}
                      </span>
                    ))}
                  </p>
                </div>
              </div>

              {/* Domain Roles Section */}
              <h3 className="text-xl mb-4 font-bold text-[#f4d1c2] underline">
                Domain Roles:
              </h3>
              <div className="space-y-4 ">
                {userData?.domainRoles?.map((domain, index) => (
                  <div
                    key={index}
                    className="bg-zinc-700 p-4 rounded-lg shadow-md"
                  >
                    <h4 className="text-md font-semibold">
                      Domain ID: {domain.domainId}{" "}
                      {domainMap[domain.domainId] &&
                        `(${domainMap[domain.domainId]})`}
                    </h4>
                    <div className="flex flex-col sm:flex-row sm:justify-between mt-2">
                      <div>
                        <p className="font-semibold text-gray-400">Roles</p>
                        <ul className="list-disc ml-5">
                          {domain.role?.map((role, roleIndex) => (
                            <li key={roleIndex}>{role}</li>
                          ))}
                        </ul>
                      </div>
                      <div>
                        <p className="font-semibold text-gray-400">
                          Invited On:
                        </p>
                        <p className="break-words p-2">
                          {formatEpochTime(domain.invitedOn)}
                        </p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* YAML Editor Modal Placeholder - you can implement this separately */}
      <div
        id="yamlEditorModal"
        className="hidden fixed inset-0 bg-black bg-opacity-50 z-50"
      >
        {/* Modal content */}
      </div>
    </div>
  );
};

export default UserDetails;
