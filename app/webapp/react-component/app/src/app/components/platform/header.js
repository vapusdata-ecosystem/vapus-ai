"use client";
import { useState, useEffect, useRef } from "react";
import Link from "next/link";

const Header = ({
  sectionHeader = "",
  hideBackListingLink = false,
  backListingLink = "",
  globalContext = null,
}) => {
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
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

  const handleDomainLogin = (domainId) => {
    const urlObj = new URL(window.location.href);
    const pathname = urlObj.pathname;
    window.location.href = `/ui/auth/domain/${domainId}?redirect=${pathname}`;
  };

  // Function to limit string length similar to the template's limitletters function
  const limitLetters = (str, maxLength) => {
    return str.length > maxLength ? str.substring(0, maxLength) + "..." : str;
  };

  return (
    <header className="bg-[#1b1b1b] text-gray-100 h-16">
      <div id="toast" className="toast"></div>
      <div id="errorMessage" className="errorMessage"></div>
      <div id="infoMessage" className="infoMessage"></div>
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
            className="flex items-center text-orange-700 px-4 py-2 rounded-lg hover:bg-zinc-800 hover:text-gray-100"
            suppressHydrationWarning
          >
            My Domains
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
          </button>

          {isDropdownOpen && (
            <div className="absolute right-0 top-full mt-2 w-64 p-2 max-h-64 overflow-y-auto border border-zinc-500 rounded-lg z-10 bg-[#1b1b1b] text-gray-100 shadow-lg">
              <ul className="py-2">
                {globalContext?.userInfo?.domainRoles?.map((role) => (
                  <li key={role.domainId}>
                    <a
                      onClick={() => handleDomainLogin(role.domainId)}
                      className={`block rounded-md shadow-sm px-2 py-1 m-1 bg-[#1b1b1b] text-gray-100 hover:bg-zinc-800 hover:text-gray-100 text-xs cursor-pointer w-full ${
                        role.domainId === globalContext?.currentDomain?.domainId
                          ? "bg-zinc-800"
                          : ""
                      }`}
                      data-domain={role.domainId}
                      style={{ width: "400px" }}
                    >
                      <strong>
                        {globalContext?.domainMap?.[role.domainId]} (
                        {limitLetters(role.domainId, 15)})
                      </strong>
                    </a>
                  </li>
                ))}
              </ul>
            </div>
          )}
        </div>
      </div>
    </header>
  );
};

export default Header;
