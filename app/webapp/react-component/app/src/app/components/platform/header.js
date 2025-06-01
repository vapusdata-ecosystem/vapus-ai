"use client";
import { useState, useEffect, useRef } from "react";
import Link from "next/link";
import { userGlobalData } from "@/context/GlobalContext";
import { userProfileApi } from "@/app/utils/settings-endpoint/profile-api";
import { OrgAccessTokenApi } from "@/app/utils/org-accessToken-endpoint/organization-api";
import { toast } from "react-toastify";
import ToastContainerMessage from "../notification/customToast";

const Header = ({
  sectionHeader = "",
  hideBackListingLink = false,
  backListingLink = "",
}) => {
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [userRoles, setUserRoles] = useState([]);
  const [currentGlobalContext, setCurrentGlobalContext] = useState(null);
  const [isOrgSwitching, setIsOrgSwitching] = useState(false);

  const dropdownRef = useRef(null);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsDropdownOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const toggleDropdown = (e) => {
    e.stopPropagation();
    setIsDropdownOpen(!isDropdownOpen);
  };

  // Function to set cookie
  const setCookie = (name, value, days = 7) => {
    const expires = new Date();
    expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000);
    document.cookie = `${name}=${value};expires=${expires.toUTCString()};path=/;SameSite=Lax`;
  };

  // Function to delete existing cookie
  const deleteCookie = (name) => {
    document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 UTC;path=/;`;
  };

  // Updated function to handle organization login with API call
  const handleDomainLogin = async (organizationId) => {
    try {
      setIsOrgSwitching(true);
      setError(null);

      const payload = {
        organization: organizationId,
        utility: "ORGANIZATION_LOGIN",
      };

      // API endpoint call
      const output = await OrgAccessTokenApi.getOrgAccessToken(payload);

      // Check if we have the expected response structure
      if (output && output.token && output.token.accessToken) {
        deleteCookie("access_token");
        setCookie("access_token", output.token.accessToken);

        toast.success(
          "Organization Switched",
          `Switched to ${organizationId} successfully.`
        );

        setTimeout(() => {
          // Redirect to current page
          window.location.reload();
        }, 1000);
      } else {
        toast.success(
          "Organization Switched",
          "Organization switched successfully."
        );
        setTimeout(() => {
          window.location.reload();
        }, 1000);
      }
    } catch (error) {
      console.error("Error switching organization:", error);
      toast.error("Organization Switch Failed");
      setError(`Failed to switch organization: ${error.message}`);
    } finally {
      setIsOrgSwitching(false);
    }
  };

  // Function to limit string length similar to the template's limitletters function
  const limitLetters = (str, maxLength) => {
    return str.length > maxLength ? str.substring(0, maxLength) + "..." : str;
  };

  // Helper function to check if organization is current
  const isCurrentOrg = (orgId) => {
    return orgId === currentGlobalContext?.currentOrganization?.organizationId;
  };

  useEffect(() => {
    let isMounted = true;

    const fetchData = async () => {
      try {
        setIsLoading(true);
        setError(null);

        // Get user data from global context
        const globalContext = await userGlobalData();
        setCurrentGlobalContext(globalContext);
        if (!isMounted) return;

        // Check if userId exists
        if (globalContext?.userInfo?.userId) {
          const userId = globalContext.userInfo.userId;
          // Make API call to get user profile with userId
          const data = await userProfileApi.getuserProfile(userId);
          console.log("Header User Profile Data:", data);

          if (!isMounted) return;

          if (data.output?.users && data.output.users.length > 0) {
            // Extract roles from the first user
            const userRolesData = data.output.users[0].roles || [];
            setUserRoles(userRolesData);
            console.log("User Roles:", userRolesData);
          } else {
            console.error("No users found in API response");
            setError("No users found in API response");
          }
        } else {
          console.error("User ID not found in global context");
          setError("User ID not found in global context");
        }
      } catch (error) {
        console.error("Error fetching user data:", error);
        if (isMounted) {
          setError(error.message || "Failed to fetch user data");
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    };

    fetchData();

    return () => {
      isMounted = false;
    };
  }, []);

  return (
    <header className="bg-[#1b1b1b] text-gray-100 h-16">
      <ToastContainerMessage />
      <div className="mx-auto flex justify-between items-center p-2">
        {/* Left Section: Back Button */}
        <div className="flex justify-around gap-10">
          {!hideBackListingLink && (
            <div className="flex justify-between mb-2 items-center rounded-lg p-1">
              <Link
                href={backListingLink}
                className="inline-flex items-center text-orange-700 hover:text-orange-700"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-8 w-8 m-1"
                  viewBox="0 0 36 36"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2.5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <title>Go Back - Contained Arrow</title>
                  <rect x="4" y="4" width="28" height="28" rx="14" />
                  <line x1="23" y1="18" x2="13" y2="18"></line>
                  <polyline points="17 22 13 18 17 14"></polyline>
                </svg>
                Back
              </Link>
            </div>
          )}
        </div>

        {/* Middle Section: Section Header */}
        <div>
          <h1 className="text-xl text-gray-100 font-bold">{sectionHeader}</h1>
        </div>

        {/* Right Section: Domain Dropdown */}
        <div className="flex relative" ref={dropdownRef}>
          <button
            onClick={toggleDropdown}
            className="flex items-center text-orange-700 px-4 py-2 rounded-lg hover:bg-zinc-800 hover:text-gray-100 cursor-pointer"
            suppressHydrationWarning
            disabled={isOrgSwitching}
          >
            My Organization
            {isOrgSwitching ? (
              <svg
                className="animate-spin ml-2 h-4 w-4"
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
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
            ) : (
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="16"
                height="16"
                fill="none"
                viewBox="0 0 16 16"
                className="w-5 h-5 shrink-0"
              >
                <path
                  stroke="currentColor"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="1.5"
                  d="M5.333 6 8 3.334 10.667 6m0 4L8 12.667 5.333 10"
                ></path>
              </svg>
            )}
          </button>

          {isDropdownOpen && (
            <div className="absolute right-0 top-full mt-2 w-70 p-2 border border-zinc-500 rounded-lg z-10 bg-[#1b1b1b] text-gray-100 shadow-lg">
              {isLoading ? (
                <div className="py-2 px-4 text-center text-gray-400">
                  Loading...
                </div>
              ) : error ? (
                <div className="py-2 px-4 text-center text-red-400">
                  Error: {error}
                </div>
              ) : (
                <ul className="py-2">
                  {userRoles.length > 0 ? (
                    userRoles.map((role, index) => {
                      const isCurrent = isCurrentOrg(role.organizationId);
                      return (
                        <li key={role.organizationId || index}>
                          <button
                            onClick={() =>
                              handleDomainLogin(role.organizationId)
                            }
                            disabled={isOrgSwitching}
                            className={`block rounded-md shadow-sm px-2 py-1 m-1 text-gray-100 text-xs cursor-pointer w-full break-words transition-colors ${
                              isCurrent
                                ? "bg-orange-600 hover:bg-orange-700 border border-orange-500"
                                : "bg-zinc-700 hover:bg-zinc-800"
                            } ${
                              isOrgSwitching
                                ? "opacity-50 cursor-not-allowed"
                                : ""
                            }`}
                            data-organization={role.organizationId}
                          >
                            <div className="break-words flex items-center justify-between">
                              <span className="flex items-center">
                                {isCurrent && (
                                  <svg
                                    className="w-3 h-3 mr-2 text-white"
                                    fill="currentColor"
                                    viewBox="0 0 20 20"
                                    xmlns="http://www.w3.org/2000/svg"
                                  >
                                    <path
                                      fillRule="evenodd"
                                      d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                      clipRule="evenodd"
                                    />
                                  </svg>
                                )}
                                {role.organizationId}
                                {isCurrent && (
                                  <span className="ml-2 text-xs opacity-90">
                                    (Current)
                                  </span>
                                )}
                              </span>
                              {isOrgSwitching && (
                                <svg
                                  className="animate-spin h-3 w-3"
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
                                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                  ></path>
                                </svg>
                              )}
                            </div>
                          </button>
                        </li>
                      );
                    })
                  ) : (
                    <li className="py-2 px-4 text-center text-gray-400">
                      No organizations found
                    </li>
                  )}
                </ul>
              )}
            </div>
          )}
        </div>
      </div>
    </header>
  );
};

export default Header;
