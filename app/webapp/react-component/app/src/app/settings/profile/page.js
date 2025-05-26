"use client";
import React, { useState, useEffect } from "react";
import Header from "@/app/components/platform/header";
import { userGlobalData } from "@/context/GlobalContext";
import { userProfileApi } from "../../utils/settings-endpoint/profile-api";
import ActionDropdown from "../../components/action-dropdown";

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
      weblink: "./profile/update",
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
          <div className="flex justify-end mb-2">
            <ActionDropdown response={responseData} />
          </div>

          <div className="overflow-x-auto text-gray-100 bg-zinc-800 rounded-lg p-8 shadow-md">
            {/* Profile Header Section */}
            <div className="bg-[#1b1b1b] rounded-lg p-4 shadow-md">
              <div className="flex  items-center text-center space-y-4">
                {/* Profile Avatar */}
                <div className="w-24 h-24 rounded-full border-4 border-gray-600 overflow-hidden bg-zinc-700 flex items-center justify-center">
                  {userData?.profile?.avatar ? (
                    <img
                      src={userData.profile.avatar}
                      alt="Profile Avatar"
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    <div className="text-gray-400 text-center">
                      <svg
                        className="w-12 h-12 mx-auto mb-2"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
                          clipRule="evenodd"
                        />
                      </svg>
                      <span className="text-sm">No Image</span>
                    </div>
                  )}
                </div>

                {/* Display Name */}
                <div>
                  <h2 className="text-2xl font-bold  mb-2">
                    {userData?.displayName || "No Display Name"}
                  </h2>
                  <p>{userData?.email || "No Email"}</p>
                </div>

                {/* Status Badge */}
                <div>
                  <span
                    className={`px-4 py-2 text-sm font-medium rounded-full ${
                      userData?.resourceBase?.status === "ACTIVE"
                        ? "text-green-800 bg-green-100"
                        : "text-red-800 bg-red-100"
                    }`}
                  >
                    {userData?.resourceBase?.status || "N/A"}
                  </span>
                </div>
              </div>
            </div>
            {/* User Information */}
            <div className="tab-content mt-2 bg-[#1b1b1b] rounded-lg p-4">
              {/* User Details Section */}
              <h3 className="text-xl mb-4 font-bold text-[#f4d1c2] underline">
                User Details:
              </h3>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-6">
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    User ID:
                  </p>
                  <p className="break-words p-2">{userData?.userId || "N/A"}</p>
                </div>
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    First Name:
                  </p>
                  <p className="break-words p-2">
                    {userData?.firstName || "N/A"}
                  </p>
                </div>
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Last Name:
                  </p>
                  <p className="break-words p-2">
                    {userData?.lastName || "N/A"}
                  </p>
                </div>
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Onboarding Type:
                  </p>
                  <p className="break-words p-2">
                    {userData?.invitedType || "N/A"}
                  </p>
                </div>
                <div className="lg:flex items-center">
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

              {/* Addresses Section */}
              {userData?.profile?.addresses &&
                userData.profile.addresses.length > 0 && (
                  <>
                    <h3 className="text-xl mb-4 font-bold text-[#f4d1c2] underline">
                      Addresses:
                    </h3>
                    <div className="space-y-4 mb-6">
                      {userData.profile.addresses.map((address, index) => (
                        <div
                          key={index}
                          className="bg-zinc-700 p-4 rounded-lg shadow-md"
                        >
                          <h4 className="text-lg font-semibold mb-3">
                            Address {index + 1}
                          </h4>
                          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 text-sm">
                            {address.street_address1 && (
                              <div>
                                <span className="text-gray-400">
                                  Street Address 1:
                                </span>
                                <p>{address.street_address1}</p>
                              </div>
                            )}
                            {address.street_address2 && (
                              <div>
                                <span className="text-gray-400">
                                  Street Address 2:
                                </span>
                                <p>{address.street_address2}</p>
                              </div>
                            )}
                            {address.city && (
                              <div>
                                <span className="text-gray-400">City:</span>
                                <p>{address.city}</p>
                              </div>
                            )}
                            {address.state && (
                              <div>
                                <span className="text-gray-400">State:</span>
                                <p>{address.state}</p>
                              </div>
                            )}
                            {address.zip_code && (
                              <div>
                                <span className="text-gray-400">Zip Code:</span>
                                <p>{address.zip_code}</p>
                              </div>
                            )}
                            {address.country && (
                              <div>
                                <span className="text-gray-400">Country:</span>
                                <p>{address.country}</p>
                              </div>
                            )}
                            {address.others && (
                              <div className="sm:col-span-2">
                                <span className="text-gray-400">Others:</span>
                                <p>{address.others}</p>
                              </div>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  </>
                )}

              {/* Domain Roles Section */}
              <h3 className="text-xl mb-4 font-bold text-[#f4d1c2] underline">
                Domain Roles:
              </h3>
              <div className="space-y-4">
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
