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
  const [currentOrgId, setCurrentOrgId] = useState(null);
  const [searchTerm, setSearchTerm] = useState("");

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
  const handleOrganizationLogin = async (organizationId) => {
    try {
      setIsOrgSwitching(true);
      setError(null);
      setIsDropdownOpen(false);

      const payload = {
        organization: organizationId,
        utility: "ORGANIZATION_LOGIN",
      };

      // API call
      const output = await OrgAccessTokenApi.getOrgAccessToken(payload);

      console.log("Organization switched:", output);

      if (output && output.token && output.token.accessToken) {
        deleteCookie("access_token");
        setCookie("access_token", output.token.accessToken);

        toast.success(
          "Organization Switched",
          `Switched to ${organizationId} successfully.`
        );

        setCurrentOrgId(organizationId);

        setTimeout(() => {
          window.location.reload();
        }, 1000);
      } else {
        toast.success(
          "Organization Switched",
          "Organization switched successfully."
        );
        
        setCurrentOrgId(organizationId);
        
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

  // Enhanced helper function to check if organization is current
  const isCurrentOrg = (orgId) => {
    if (currentOrgId) {
      return orgId === currentOrgId;
    }
    
    const contextOrgId = currentGlobalContext?.currentOrganization?.organizationId; 
    return orgId === contextOrgId;
  };

  // Get current organization name for display
  const getCurrentOrgName = () => {
    if (currentOrgId) {
      const currentRole = userRoles.find(role => role.organizationId === currentOrgId);
      return currentRole ? currentRole.organizationId : "Select Organization";
    }
    
    const contextOrgId = currentGlobalContext?.currentOrganization?.organizationId;
    if (contextOrgId) {
      const currentRole = userRoles.find(role => role.organizationId === contextOrgId);
      return currentRole ? currentRole.organizationId : "Select Organization";
    }
    
    return "Select Organization";
  };

  // Filter organizations based on search term
  const filteredRoles = userRoles.filter(role => 
    role.organizationId.toLowerCase().includes(searchTerm.toLowerCase())
  );

  useEffect(() => {
    let isMounted = true;

    const fetchData = async () => {
      try {
        setIsLoading(true);
        setError(null);

        const globalContext = await userGlobalData();
        setCurrentGlobalContext(globalContext);
        
        let currentOrgIdFromContext = null;
        if (globalContext?.currentOrganization?.organizationId) {
          currentOrgIdFromContext = globalContext.currentOrganization.organizationId;
        } else if (globalContext?.organization?.organizationId) {
          currentOrgIdFromContext = globalContext.organization.organizationId;
        } else if (globalContext?.organizationId) {
          currentOrgIdFromContext = globalContext.organizationId;
        } else if (globalContext?.userInfo?.organizationId) {
          currentOrgIdFromContext = globalContext.userInfo.organizationId;
        }
        
        console.log("Detected current org ID:", currentOrgIdFromContext);
        
        if (currentOrgIdFromContext) {
          setCurrentOrgId(currentOrgIdFromContext);
        }
        
        if (!isMounted) return;

        if (globalContext?.userInfo?.userId) {
          const userId = globalContext.userInfo.userId;
          const data = await userProfileApi.getuserProfile(userId);
          console.log("Header User Profile Data:", data);

          if (!isMounted) return;

          if (data.output?.users && data.output.users.length > 0) {
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
    <header className="bg-[#1b1b1b] text-gray-100 h-[66px] border-b border-zinc-500">
      <div id="toast" className="toast"></div>
      <div id="errorMessage" className="errorMessage"></div>
      <div id="infoMessage" className="infoMessage"></div>
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

        {/* Right Section: Organization Dropdown */}
        <div className="flex relative z-10">
          <div className="relative overflow-visible z-10 w-80 border border-zinc-500 rounded-md bg-zinc-900" ref={dropdownRef}>
            {/* Dropdown Button */}
            <button
              onClick={toggleDropdown}
              className="flex items-center justify-between w-full px-4 py-2 bg-zinc-900 text-orange-400 rounded-lg hover:bg-zinc-800 focus:outline-none transition"
              disabled={isOrgSwitching}
            >
              <span className="truncate">
                {getCurrentOrgName()}
              </span>
              {isOrgSwitching ? (
                <svg
                  className="animate-spin ml-2 h-5 w-5"
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
                  className={`w-5 h-5 ml-2 transition-transform ${isDropdownOpen ? 'rotate-180' : ''}`}
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    fillRule="evenodd"
                    d="M5.23 7.21a.75.75 0 011.06.02L10 10.586l3.71-3.356a.75.75 0 111.02 1.1l-4.25 3.84a.75.75 0 01-1.02 0l-4.25-3.84a.75.75 0 01.02-1.06z"
                    clipRule="evenodd"
                  />
                </svg>
              )}
            </button>

            {/* Dropdown Menu */}
            {isDropdownOpen && (
              <div className="absolute z-20 mt-2 w-full bg-zinc-900 rounded-xl shadow-lg border border-zinc-500">
                {/* Search Input */}
                <input
                  type="text"
                  placeholder="Search organizations..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="w-full px-4 py-2 bg-zinc-800 text-gray-200 placeholder-zinc-500 rounded-t-xl focus:outline-none"
                />
                 <div className="h-px bg-zinc-500  "></div>
                {/* Dropdown Options with scroll */}
                <div className="max-h-[150px] overflow-y-auto scrollbar">
                  {isLoading ? (
                    <div className="px-4 py-2 text-center text-gray-400">
                      Loading...
                    </div>
                  ) : error ? (
                    <div className="px-4 py-2 text-center text-red-400">
                      Error: {error}
                    </div>
                  ) : (
                    <ul className="relative z-10">
                      {filteredRoles.length > 0 ? (
                        filteredRoles.map((role, index) => {
                          const isCurrent = isCurrentOrg(role.organizationId);
                          
                          return (
                        <li key={role.organizationId || index} className="group relative z-10 flex justify-between items-center px-4 py-2 hover:bg-zinc-800 text-gray-100 cursor-pointer" onClick={() => handleOrganizationLogin(role.organizationId)} >
                            <span className="truncate max-w-[85%] text-sm">
                              {role.organizationId}
                            </span>
                            <div className="relative flex-shrink-0 ml-2 text-gray-100 hover:text-orange-100">
                              <div className="w-4 h-4 text-gray-300 hover:text-gray-100 cursor-default bg-zinc-900 rounded-full shadow-sm">
                                {isCurrent ? (
                                  <svg className="w-full h-full text-orange-400" fill="currentColor" viewBox="0 0 20 20" >
                                    <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                                  </svg>
                                ) : (
                                  <div className="relative group/tooltip">
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" className="w-full h-full cursor-pointer" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="2" >
                                      <path strokeLinecap="round" strokeLinejoin="round" d="M13 16h-1v-4h-1m1-4h.01M12 2a10 10 0 100 20 10 10 0 000-20z" />
                                    </svg>
                                    {/* Tooltip */}
                                    <div className="absolute right-full top-1/2 transform -translate-y-1/2 mr-2 px-2 py-1 bg-zinc-900 text-white text-xs rounded opacity-0 group-hover/tooltip:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-50">
                                      {role.organizationId}
                                      {/* Tooltip arrow */}
                                      <div className="absolute left-full top-1/2 transform -translate-y-1/2 border-2 border-transparent border-l-gray-900"></div>
                                    </div>
                                  </div>
                                )}
                              </div>
                            </div>
                          </li>
                          );
                        })
                      ) : (
                        <li className="px-4 py-2 text-center text-gray-400">
                          No organizations found
                        </li>
                      )}
                    </ul>
                  )}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;