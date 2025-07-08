"use client";
import { useState, useEffect } from "react";
import Link from "next/link";
import Image from "next/image";
import { usePathname } from "next/navigation";
import { getAuthService } from "../../../../lib/auth";
import { userGlobalData } from "@/context/GlobalContext";
import { userProfileApi } from "@/app/utils/settings-endpoint/profile-api";

const Sidebar = () => {
  // Internal navigation data
  const navMenuMap = [
    // Dashboard
    {
      itemId: "dashboard",
      itemName: "Dashboard",
      url: "/dashboard",
      svg: (
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          className="h-5 w-5 m-1"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <path d="M3 9l9-6 9 6v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
          <path d="M9 22V12h6v10"></path>
        </svg>
      ),
    },
    // Insights
    {
      itemId: "insights",
      itemName: "Insights",
      url: " ",
      svg: (
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          className="h-6 w-6 m-1"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <polyline points="3 17 8 11 13 13 17 8 21 12" />
          <rect
            x="5"
            y="18"
            width="4"
            height="3"
            stroke="none"
            fill="currentColor"
          />
          <rect
            x="10"
            y="16"
            width="4"
            height="5"
            stroke="none"
            fill="currentColor"
          />
          <rect
            x="15"
            y="17"
            width="4"
            height="4"
            stroke="none"
            fill="currentColor"
          />
        </svg>
      ),
      children: [
        {
          itemId: "llms",
          itemName: "LLMS",
          url: "/insights/llms",
          svg: (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              className="h-4 w-4 m-1"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <rect x="2" y="3" width="20" height="14" rx="2" ry="2" />
              <line x1="8" y1="21" x2="16" y2="21" />
              <line x1="12" y1="17" x2="12" y2="21" />
              <circle cx="8" cy="9" r="2" />
              <path d="M16 7v4a2 2 0 0 1-2 2H10" />
            </svg>
          ),
        },
      ],
    },
    // Studios
    {
      itemId: "playground",
      itemName: "playground",
      url: "/studios/ai-studio",
      svg: (
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 20 20"
          className="h-5 w-5 m-1"
        >
          <path
            stroke="currentColor"
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeMiterlimit="22.926"
            strokeWidth="1.5"
            d="m6.514 9.06-3.988-.383 3.217-3.216a3.36 3.36 0 0 1 3.925-.595M10.95 13.55l.377 3.924 3.217-3.216a3.354 3.354 0 0 0 .52-4.06"
          ></path>
          <path
            stroke="currentColor"
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeMiterlimit="22.926"
            strokeWidth="1.5"
            d="M15.005 10.166q-.36.416-.76.817c-1.342 1.341-2.838 2.347-4.315 2.978a1.11 1.11 0 0 1-1.24-.24l-2.416-2.414a1.11 1.11 0 0 1-.24-1.24c.632-1.477 1.638-2.972 2.98-4.314 2.815-2.814 6.309-4.151 8.882-3.65.454 2.33-.599 5.412-2.89 8.063M13.673 2.65l3.556 3.555"
          ></path>
          <path
            fill="currentColor"
            d="M13.56 6.44a1.5 1.5 0 1 1-2.12 2.12 1.5 1.5 0 0 1 2.12-2.12"
          ></path>
          <path
            stroke="currentColor"
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeMiterlimit="22.926"
            strokeWidth="1.5"
            d="M6.015 13.987 2 18M4.42 12.392l-1.357 1.356M7.61 15.581l-1.356 1.356"
          ></path>
        </svg>
      ),
    },
    // AI Center
    {
      itemId: "ai center",
      itemName: "AI Center",
      url: " ",
      svg: (
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          className="h-6 w-6 m-1"
        >
          <polyline points="8 3 3 9 8 15" />
          <polyline points="16 3 21 9 16 15" />
          <circle cx="12" cy="18" r="3" fill="currentColor" />
          <line x1="12" y1="15" x2="12" y2="9" />
        </svg>
      ),
      children: [
        {
          itemId: "models registry",
          itemName: "Models Registry",
          url: "/ai-center/models-registry",
          svg: (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              className="h-4 w-4 m-1"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" />
              <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" />
              <circle cx="12" cy="8" r="2" />
              <path d="M10 16h4" />
              <path d="M12 14v2" />
            </svg>
          ),
        },
        {
          itemId: "prompts",
          itemName: "Prompts",
          url: "/ai-center/prompts",
          svg: (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              className="h-4 w-4 m-1"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
              <polyline points="14,2 14,8 20,8" />
              <line x1="16" y1="13" x2="8" y2="13" />
              <line x1="16" y1="17" x2="8" y2="17" />
              <polyline points="10,9 9,9 8,9" />
            </svg>
          ),
        },
        // {
        //   itemId: "agents",
        //   itemName: "Agents",
        //   url: "/ai-center/agents",
        //   svg: (
        //     <svg
        //       xmlns="http://www.w3.org/2000/svg"
        //       viewBox="0 0 24 24"
        //       className="h-4 w-4 m-1"
        //       fill="none"
        //       stroke="currentColor"
        //       strokeWidth="2"
        //       strokeLinecap="round"
        //       strokeLinejoin="round"
        //     >
        //       <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
        //       <circle cx="9" cy="7" r="4" />
        //       <path d="M22 21v-2a4 4 0 0 0-3-3.87" />
        //       <path d="M16 3.13a4 4 0 0 1 0 7.75" />
        //       <circle cx="12" cy="12" r="1" />
        //       <circle cx="18" cy="12" r="1" />
        //       <circle cx="6" cy="12" r="1" />
        //     </svg>
        //   ),
        // },
        {
          itemId: "guardrails",
          itemName: "Guardrails",
          url: "/ai-center/guardrails",
          svg: (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              className="h-4 w-4 m-1"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
              <path d="M9 12l2 2 4-4" />
            </svg>
          ),
        },
      ],
    },
  ];

  const bottomMenuMap = [
    // Settings
    {
      itemId: "settings",
      itemName: "Settings",
      url: " ",
      svg: (
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="icon icon-tabler icon-tabler-settings h-5 w-5 m-1"
          viewBox="0 0 24 24"
          strokeWidth="2"
          stroke="currentColor"
          fill="none"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <path stroke="none" d="M0 0h24v24H0z" fill="none" />
          <path d="M10.325 4.317a1.724 1.724 0 0 1 3.35 0l.333 1.35a7.04 7.04 0 0 1 1.817 .621l1.308 -.478a1.724 1.724 0 0 1 2.156 2.156l-.478 1.308a7.034 7.034 0 0 1 .621 1.817l1.35 .333a1.724 1.724 0 0 1 0 3.35l-1.35 .333a7.034 7.034 0 0 1 -.621 1.817l.478 1.308a1.724 1.724 0 0 1 -2.156 2.156l-1.308 -.478a7.04 7.04 0 0 1 -1.817 .621l-.333 1.35a1.724 1.724 0 0 1 -3.35 0l-.333 -1.35a7.04 7.04 0 0 1 -1.817 -.621l-1.308 .478a1.724 1.724 0 0 1 -2.156 -2.156l.478 -1.308a7.034 7.034 0 0 1 -.621 -1.817l-1.35 -.333a1.724 1.724 0 0 1 0 -3.35l1.35 -.333a7.034 7.034 0 0 1 .621 -1.817l-.478 -1.308a1.724 1.724 0 0 1 2.156 -2.156l1.308 .478a7.04 7.04 0 0 1 1.817 -.621z" />
          <circle cx="12" cy="12" r="3" />
        </svg>
      ),
      children: [
        {
          itemId: "organization",
          itemName: "organization",
          url: "/settings/organization",
          svg: (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              className="h-4 w-4 m-1"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <circle cx="12" cy="12" r="10" />
              <line x1="2" y1="12" x2="22" y2="12" />
              <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z" />
            </svg>
          ),
        },
        {
          itemId: "platform",
          itemName: "Platform",
          url: "/settings/platform",
          svg: (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              className="h-4 w-4 m-1"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <rect x="2" y="3" width="20" height="14" rx="2" ry="2" />
              <line x1="8" y1="21" x2="16" y2="21" />
              <line x1="12" y1="17" x2="12" y2="21" />
            </svg>
          ),
        },
        {
          itemId: "users",
          itemName: "Users",
          url: "/settings/users",
          svg: (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              className="h-4 w-4 m-1"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
              <circle cx="9" cy="7" r="4" />
              <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
              <path d="M16 3.13a4 4 0 0 1 0 7.75" />
            </svg>
          ),
        },
        {
          itemId: "plugins",
          itemName: "Plugins",
          url: "/settings/plugins",
          svg: (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              className="h-4 w-4 m-1"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M12 2v6l3-3 3 3v4l-3-3-3 3V2z" />
              <path d="M12 17.5L9.5 15 7 17.5 9.5 20l2.5-2.5z" />
              <circle cx="12" cy="12" r="2" />
            </svg>
          ),
        },
        {
          itemId: "platform-organizations",
          itemName: "Platform organizations",
          url: "/settings/platform-organizations",
          svg: (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              className="h-4 w-4 m-1"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z" />
              <circle cx="12" cy="12" r="4" />
            </svg>
          ),
        },
        {
          itemId: "secret store",
          itemName: "Secret Store",
          url: "/settings/secret-store",
          svg: (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              className="h-4 w-4 m-1"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
              <circle cx="12" cy="16" r="1" />
              <path d="M7 11V7a5 5 0 0 1 10 0v4" />
            </svg>
          ),
        },
      ],
    },
    // Developers
    {
      itemId: "developers",
      itemName: "Developers",
      url: " ",
      svg: (
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="h-5 w-5 m-1"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <polyline points="16 18 22 12 16 6"></polyline>
          <polyline points="8 6 2 12 8 18"></polyline>
        </svg>
      ),
      children: [
        {
          itemId: "resources",
          itemName: "Resources",
          url: "/developers/resources",
          svg: (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              className="h-4 w-4 m-1"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
              <polyline points="14,2 14,8 20,8" />
              <line x1="16" y1="13" x2="8" y2="13" />
              <line x1="16" y1="17" x2="8" y2="17" />
              <line x1="10" y1="9" x2="8" y2="9" />
            </svg>
          ),
        },
        {
          itemId: "enums",
          itemName: "Enums",
          url: "/developers/enums",
          svg: (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              className="h-4 w-4 m-1"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M8 6h13" />
              <path d="M8 12h13" />
              <path d="M8 18h13" />
              <path d="M3 6h.01" />
              <path d="M3 12h.01" />
              <path d="M3 18h.01" />
            </svg>
          ),
        },
      ],
    },
  ];

  const [isSidebarExpanded, setIsSidebarExpanded] = useState(true);
  const [openDropdowns, setOpenDropdowns] = useState(() => {
    const initialState = {};
    navMenuMap.forEach((item) => {
      if (item.children) {
        initialState[item.itemId] = true;
      }
    });
    bottomMenuMap.forEach((item) => {
      if (item.children) {
        initialState[item.itemId] = false;
      }
    });
    return initialState;
  });
  const [isHamburgerMenuOpen, setIsHamburgerMenuOpen] = useState(false);
  const [activeNav, setActiveNav] = useState("");
  const [activeSideBar, setActiveSideBar] = useState("");
  const [activeSubmenu, setActiveSubmenu] = useState("");
  const [userData, setUserData] = useState(null);
  const [loading, setLoading] = useState(true);
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

  // clear cookies and redirect to login
  const handleLogout = () => {
    const authService = getAuthService();
    authService.logout();
  };

  const pathname = usePathname();

  const updateActiveStates = (currentPath) => {
    let navFound = false;
    let sidebarFound = false;
    let exactMatches = [];

    // Collect all possible matches
    const allMenuItems = [...navMenuMap, ...bottomMenuMap];

    // Then, collect all exact child matches
    for (const item of allMenuItems) {
      if (item.children) {
        for (const child of item.children) {
          if (
            currentPath === child.url ||
            currentPath.startsWith(child.url + "/")
          ) {
            exactMatches.push({
              parentId: item.itemId,
              childId: child.itemId,
              url: child.url,
              isChild: true,
            });
          }
        }
      }
    }

    if (exactMatches.length > 0) {
      exactMatches.sort((a, b) => b.url.length - a.url.length);
      const bestMatch = exactMatches[0];
      setActiveNav(bestMatch.parentId);
      setActiveSideBar(bestMatch.childId);
      navFound = true;
      sidebarFound = true;
    } else {
      // back to parent matches if no child matches
      for (const item of allMenuItems) {
        if (
          currentPath === item.url ||
          currentPath.startsWith(item.url + "/")
        ) {
          setActiveNav(item.itemId);
          navFound = true;
          break;
        }
      }
    }

    if (!navFound) setActiveNav("");
    if (!sidebarFound) setActiveSideBar("");
  };

  useEffect(() => {
    updateActiveStates(pathname);
  }, [pathname]);

  useEffect(() => {
    if (typeof window !== "undefined") {
      const path = window.location.pathname;
      updateActiveStates(path);
    }
  }, []);

  const toggleSidebar = () => {
    setIsSidebarExpanded(!isSidebarExpanded);
  };

  const toggleDropdown = (itemId) => {
    setOpenDropdowns((prev) => ({
      ...prev,
      [itemId]: !prev[itemId],
    }));
  };

  // Add missing toggleSubmenu function
  const toggleSubmenu = (itemId) => {
    setActiveSubmenu(activeSubmenu === itemId ? "" : itemId);
  };

  const toggleHamburgerMenu = () => {
    setIsHamburgerMenuOpen(!isHamburgerMenuOpen);
  };

  // Close hamburger menu when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      const hamburgerBtn = event.target.closest(
        "#hamburgerSideButton, #hamburgerSideButtonCollapsed"
      );
      const hamburgerMenu = event.target.closest(
        "#hamburgerSideMenu, #hamburgerSideMenuCollapsed"
      );
      if (!hamburgerBtn && !hamburgerMenu) {
        setIsHamburgerMenuOpen(false);
      }
    };
    document.addEventListener("click", handleClickOutside);
    return () => {
      document.removeEventListener("click", handleClickOutside);
    };
  }, []);

  // Close submenu when clicking outside (for collapsed sidebar)
  useEffect(() => {
    const handleClickOutside = (event) => {
      const submenu = event.target.closest(".submenu");
      const menuItem = event.target.closest(".main-item");
      if (!submenu && !menuItem) {
        setActiveSubmenu("");
      }
    };
    document.addEventListener("click", handleClickOutside);
    return () => {
      document.removeEventListener("click", handleClickOutside);
    };
  }, []);

  return (
    <>
      <div className="bg-zinc-800 flex h-screen">
        {/* Expanded Sidebar */}
        <aside
          className={`flex h-screen w-[208px] flex-col bg-[#1b1b1b] border-r border-zinc-500 text-gray-100 shadow-lg transition-all duration-300 ${
            !isSidebarExpanded ? "hidden" : ""
          }`}
        >
          <div className="flex items-center h-[65px]">
            <Link href="/dashboard/">
              <Image
                src="https://storage.googleapis.com/vapusdata-public/website-images/vapus-white.webp"
                alt="Logo"
                width={136}
                height={88}
                className="h-10 ml-4"
              />
            </Link>
           <button
              className=" rounded-lg ml-4 relative group inline-block cursor-pointer"
              onClick={toggleSidebar}
              suppressHydrationWarning
            >
              <svg 
                viewBox="0 0 20 20" 
                fill="currentColor" 
                xmlns="http://www.w3.org/2000/svg" 
                data-rtl-flip="" 
                className="icon w-8 h-8"
              >
                <path d="M6.83496 3.99992C6.38353 4.00411 6.01421 4.0122 5.69824 4.03801C5.31232 4.06954 5.03904 4.12266 4.82227 4.20012L4.62207 4.28606C4.18264 4.50996 3.81498 4.85035 3.55859 5.26848L3.45605 5.45207C3.33013 5.69922 3.25006 6.01354 3.20801 6.52824C3.16533 7.05065 3.16504 7.71885 3.16504 8.66301V11.3271C3.16504 12.2712 3.16533 12.9394 3.20801 13.4618C3.25006 13.9766 3.33013 14.2909 3.45605 14.538L3.55859 14.7216C3.81498 15.1397 4.18266 15.4801 4.62207 15.704L4.82227 15.79C5.03904 15.8674 5.31234 15.9205 5.69824 15.9521C6.01398 15.9779 6.383 15.986 6.83398 15.9902L6.83496 3.99992ZM18.165 11.3271C18.165 12.2493 18.1653 12.9811 18.1172 13.5702C18.0745 14.0924 17.9916 14.5472 17.8125 14.9648L17.7295 15.1415C17.394 15.8 16.8834 16.3511 16.2568 16.7353L15.9814 16.8896C15.5157 17.1268 15.0069 17.2285 14.4102 17.2773C13.821 17.3254 13.0893 17.3251 12.167 17.3251H7.83301C6.91071 17.3251 6.17898 17.3254 5.58984 17.2773C5.06757 17.2346 4.61294 17.1508 4.19531 16.9716L4.01855 16.8896C3.36014 16.5541 2.80898 16.0434 2.4248 15.4169L2.27051 15.1415C2.03328 14.6758 1.93158 14.167 1.88281 13.5702C1.83468 12.9811 1.83496 12.2493 1.83496 11.3271V8.66301C1.83496 7.74072 1.83468 7.00898 1.88281 6.41985C1.93157 5.82309 2.03329 5.31432 2.27051 4.84856L2.4248 4.57317C2.80898 3.94666 3.36012 3.436 4.01855 3.10051L4.19531 3.0175C4.61285 2.83843 5.06771 2.75548 5.58984 2.71281C6.17898 2.66468 6.91071 2.66496 7.83301 2.66496H12.167C13.0893 2.66496 13.821 2.66468 14.4102 2.71281C15.0069 2.76157 15.5157 2.86329 15.9814 3.10051L16.2568 3.25481C16.8833 3.63898 17.394 4.19012 17.7295 4.84856L17.8125 5.02531C17.9916 5.44285 18.0745 5.89771 18.1172 6.41985C18.1653 7.00898 18.165 7.74072 18.165 8.66301V11.3271ZM8.16406 15.995H12.167C13.1112 15.995 13.7794 15.9947 14.3018 15.9521C14.8164 15.91 15.1308 15.8299 15.3779 15.704L15.5615 15.6015C15.9797 15.3451 16.32 14.9774 16.5439 14.538L16.6299 14.3378C16.7074 14.121 16.7605 13.8478 16.792 13.4618C16.8347 12.9394 16.835 12.2712 16.835 11.3271V8.66301C16.835 7.71885 16.8347 7.05065 16.792 6.52824C16.7605 6.14232 16.7073 5.86904 16.6299 5.65227L16.5439 5.45207C16.32 5.01264 15.9796 4.64498 15.5615 4.3886L15.3779 4.28606C15.1308 4.16013 14.8165 4.08006 14.3018 4.03801C13.7794 3.99533 13.1112 3.99504 12.167 3.99504H8.16406C8.16407 3.99667 8.16504 3.99829 8.16504 3.99992L8.16406 15.995Z"></path>
                <text 
                  x="12" 
                  y="11.5" 
                  textAnchor="middle" 
                  fontFamily="Arial, sans-serif" 
                  fontSize="6" 
                  fill="currentColor"
                >
                  ❮
                </text>
              </svg>
             <div className="absolute hidden group-hover:flex bg-gray-700 text-gray-100 text-sm p-2 rounded shadow-lg left-full top-1/2 -translate-y-1/2 ml-2 whitespace-nowrap z-[1010]">
                  Click here to collapse sidebar
              </div>
            </button>
          </div>
          <div className="h-px bg-zinc-500 "></div>

          {/* Main Navigation Items */}
          <nav className="flex flex-col space-y-1 bg-[#1b1b1b] p-1 flex-1 overflow-y-auto scrollbar">
            {navMenuMap.map((main) => (
              <div key={main.itemId} className="main-item">
                <div
                  className={`flex items-center justify-between w-full text-sm p-1 mb-1 cursor-pointer hover:bg-zinc-800 hover:text-gray-100 rounded-md ${
                    activeNav === main.itemId
                      ? "bg-zinc-600 text-gray-100"
                      : "text-gray-100"
                  }`}
                  onClick={() =>
                    main.children ? toggleDropdown(main.itemId) : null
                  }
                >
                  <Link href={main.url} className="flex items-center flex-1">
                    {main.svg}
                    <span>{main.itemName}</span>
                  </Link>
                  {main.children && (
                    <svg
                      className={`w-3 h-3 transition-transform transform ${
                        openDropdowns[main.itemId] ? "rotate-180" : "rotate-0"
                      }`}
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M19 9l-7 7-7-7"
                      ></path>
                    </svg>
                  )}
                </div>

                {/* Children Dropdown Menu with Connector Lines */}
                {main.children && (
                  <div
                    className={`ml-4 overflow-hidden transition-all duration-300 ease-in-out ${
                      openDropdowns[main.itemId]
                        ? "max-h-96 opacity-100"
                        : "max-h-0 opacity-0"
                    }`}
                  >
                    <div className="relative">
                      <div className="space-y-1 ">
                        {main.children.map((sub, index) => (
                          <div key={sub.itemId} className="relative">
                            <Link
                              href={sub.url}
                              className={`block text-xs text-gray-300 p-1 pl-3 rounded hover:bg-zinc-700 hover:text-white transition-colors ${
                                activeSideBar === sub.itemId
                                  ? "bg-white !text-black font-bold"
                                  : ""
                              }`}
                            >
                              <div className="flex items-center">
                                {sub.svg}
                                <span>{sub.itemName}</span>
                              </div>
                            </Link>
                          </div>
                        ))}
                      </div>
                    </div>
                  </div>
                )}
              </div>
            ))}
          </nav>

          {/* Bottom Navigation */}
          <div className="flex flex-col space-y-1 bg-[#1b1b1b] p-1 mt-auto">
            <div className="h-px bg-zinc-500 mx-4 my-2"></div>
            <nav className="flex flex-col space-y-1 bg-[#1b1b1b] p-1">
              {bottomMenuMap.map((main) => (
                <div key={main.itemId} className="main-item">
                  <div
                    className={`flex items-center justify-between w-full text-sm p-1 mb-1 cursor-pointer hover:bg-zinc-800 hover:text-gray-100 rounded-md ${
                      activeNav === main.itemId
                        ? "bg-zinc-600 text-gray-100 font-semibold"
                        : "text-gray-100"
                    }`}
                    onClick={() =>
                      main.children ? toggleDropdown(main.itemId) : null
                    }
                  >
                    <Link href={main.url} className="flex items-center flex-1">
                      {main.svg}
                      <span>{main.itemName}</span>
                    </Link>
                    {main.children && (
                      <svg
                        className={`w-3 h-3 transition-transform transform ${
                          openDropdowns[main.itemId] ? "rotate-180" : "rotate-0"
                        }`}
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth="2"
                          d="M19 9l-7 7-7-7"
                        ></path>
                      </svg>
                    )}
                  </div>

                  {/* Dropdown Menu with Connector Lines */}
                  {main.children && (
                    <div
                      className={`ml-4 overflow-hidden transition-all duration-300 ease-in-out ${
                        openDropdowns[main.itemId]
                          ? "max-h-96 opacity-100"
                          : "max-h-0 opacity-0"
                      }`}
                    >
                      <div className="relative">
                        <div className="space-y-1 ">
                          {main.children.map((sub, index) => (
                            <div key={sub.itemId} className="relative">
                              <Link
                                href={sub.url}
                                className={`block text-xs text-gray-300 p-1 pl-3 rounded hover:bg-zinc-700 hover:text-white transition-colors ${
                                  activeSideBar === sub.itemId
                                    ? "bg-white !text-black font-bold"
                                    : ""
                                }`}
                              >
                                <div className="flex items-center">
                                  {sub.svg}
                                  <span>{sub.itemName}</span>
                                </div>
                              </Link>
                            </div>
                          ))}
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              ))}
            </nav>

            {/* User Profile Section */}
            <nav className="flex flex-col space-y-1 p-1">
              <div className="main-item flex px-2 py-2 shadow-lg rounded-full border border-zinc-500">
                <p className="flex items-center justify-center w-8 h-8 text-xs font-semibold text-black rounded-full bg-white">
                  {userData?.displayName?.charAt(0) || "?"}
                </p>
                <p className="block text-xs text-primary100 font-semibold text-gray-100 p-[10px]">
                  {userData?.displayName?.slice(0, 7) || " "}
                </p>
                <button
                  id="hamburgerSideButton"
                  className="text-orange-700 pl-4 cursor-pointer"
                  onClick={toggleHamburgerMenu}
                  suppressHydrationWarning
                >
                  <svg
                    className="w-8 h-8"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="white"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M4 6h16M4 12h16M4 18h16"
                    ></path>
                  </svg>
                </button>
              </div>

              {/* Hamburger Menu */}
              <div
                id="hamburgerSideMenu"
                className={`absolute bottom-6 left-[160px] sm:left-[185px] md:left-[200px] lg:left-[210px] bg-[#1b1b1b] text-gray-100 shadow-lg rounded-lg w-48 z-10 divide-y divide-gray-300 rounded-md border border-zinc-500 ${
                  isHamburgerMenuOpen ? "block" : "hidden"
                }`}
              >
                <div className="py-2">
                  <p className="block px-4 py-2 text-sm text-primary100 font-semibold">
                    {userData?.displayName || " "}
                  </p>
                </div>
                <ul className="py-2 text-sm text-gray-100">
                  <li>
                    <Link
                      href="/settings/profile"
                      className="flex p-2 white hover:bg-zinc-800 hover:text-gray-100"
                      suppressHydrationWarning
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        className="w-4 mr-2 h-4"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth="2"
                          d="M5.121 17.804A9.969 9.969 0 0112 15a9.969 9.969 0 016.879 2.804M15 11a3 3 0 11-6 0 3 3 0 016 0zM12 3c4.97 0 9 4.03 9 9s-4.03 9-9 9-9-4.03-9-9 4.03-9 9-9z"
                        />
                      </svg>
                      Profile
                    </Link>
                  </li>
                  <li>
                    <button
                      onClick={handleLogout}
                      className="flex p-2 white hover:bg-zinc-800 hover:text-gray-100"
                      suppressHydrationWarning
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        className="h-4 w-4 mr-2"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        strokeWidth="2"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h6a2 2 0 012 2v1"
                        />
                      </svg>
                      Logout
                    </button>
                  </li>
                </ul>
              </div>
            </nav>
          </div>
        </aside>
        {/* collapsed sidebar */}
        <aside
          className={`flex h-screen w-[40px] flex-col bg-[#1b1b1b] border-r border-zinc-500 text-gray-100 shadow-lg transition-all duration-300 ${
            isSidebarExpanded ? "hidden" : ""
          }`}
        >
          <div className="items-center my-2">
            <Link href="/dashboard/">
              <Image
                src="https://storage.googleapis.com/vapusdata-public/website-images/favicon.svg"
                alt="Logo"
                className="ml-1"
                width={32}
                height={32}
              />
            </Link>
           <button
              className="rounded-lg ml-1 mt-2 cursor-pointer"
              onClick={toggleSidebar}
            >
              <svg 
                viewBox="0 0 20 20" 
                fill="currentColor" 
                xmlns="http://www.w3.org/2000/svg" 
                data-rtl-flip="" 
                className="icon w-8 h-8"
              >
                <path d="M6.83496 3.99992C6.38353 4.00411 6.01421 4.0122 5.69824 4.03801C5.31232 4.06954 5.03904 4.12266 4.82227 4.20012L4.62207 4.28606C4.18264 4.50996 3.81498 4.85035 3.55859 5.26848L3.45605 5.45207C3.33013 5.69922 3.25006 6.01354 3.20801 6.52824C3.16533 7.05065 3.16504 7.71885 3.16504 8.66301V11.3271C3.16504 12.2712 3.16533 12.9394 3.20801 13.4618C3.25006 13.9766 3.33013 14.2909 3.45605 14.538L3.55859 14.7216C3.81498 15.1397 4.18266 15.4801 4.62207 15.704L4.82227 15.79C5.03904 15.8674 5.31234 15.9205 5.69824 15.9521C6.01398 15.9779 6.383 15.986 6.83398 15.9902L6.83496 3.99992ZM18.165 11.3271C18.165 12.2493 18.1653 12.9811 18.1172 13.5702C18.0745 14.0924 17.9916 14.5472 17.8125 14.9648L17.7295 15.1415C17.394 15.8 16.8834 16.3511 16.2568 16.7353L15.9814 16.8896C15.5157 17.1268 15.0069 17.2285 14.4102 17.2773C13.821 17.3254 13.0893 17.3251 12.167 17.3251H7.83301C6.91071 17.3251 6.17898 17.3254 5.58984 17.2773C5.06757 17.2346 4.61294 17.1508 4.19531 16.9716L4.01855 16.8896C3.36014 16.5541 2.80898 16.0434 2.4248 15.4169L2.27051 15.1415C2.03328 14.6758 1.93158 14.167 1.88281 13.5702C1.83468 12.9811 1.83496 12.2493 1.83496 11.3271V8.66301C1.83496 7.74072 1.83468 7.00898 1.88281 6.41985C1.93157 5.82309 2.03329 5.31432 2.27051 4.84856L2.4248 4.57317C2.80898 3.94666 3.36012 3.436 4.01855 3.10051L4.19531 3.0175C4.61285 2.83843 5.06771 2.75548 5.58984 2.71281C6.17898 2.66468 6.91071 2.66496 7.83301 2.66496H12.167C13.0893 2.66496 13.821 2.66468 14.4102 2.71281C15.0069 2.76157 15.5157 2.86329 15.9814 3.10051L16.2568 3.25481C16.8833 3.63898 17.394 4.19012 17.7295 4.84856L17.8125 5.02531C17.9916 5.44285 18.0745 5.89771 18.1172 6.41985C18.1653 7.00898 18.165 7.74072 18.165 8.66301V11.3271ZM8.16406 15.995H12.167C13.1112 15.995 13.7794 15.9947 14.3018 15.9521C14.8164 15.91 15.1308 15.8299 15.3779 15.704L15.5615 15.6015C15.9797 15.3451 16.32 14.9774 16.5439 14.538L16.6299 14.3378C16.7074 14.121 16.7605 13.8478 16.792 13.4618C16.8347 12.9394 16.835 12.2712 16.835 11.3271V8.66301C16.835 7.71885 16.8347 7.05065 16.792 6.52824C16.7605 6.14232 16.7073 5.86904 16.6299 5.65227L16.5439 5.45207C16.32 5.01264 15.9796 4.64498 15.5615 4.3886L15.3779 4.28606C15.1308 4.16013 14.8165 4.08006 14.3018 4.03801C13.7794 3.99533 13.1112 3.99504 12.167 3.99504H8.16406C8.16407 3.99667 8.16504 3.99829 8.16504 3.99992L8.16406 15.995Z"></path>
                <text 
                  x="12" 
                  y="11.5" 
                  textAnchor="middle" 
                  fontFamily="Arial, sans-serif" 
                  fontSize="6" 
                  fill="currentColor"
                >
                  ❯
                </text>
              </svg>
            </button>
          </div>

          <div className="h-px bg-zinc-500 mx-4 my-2"></div>

          {/* Collapsed Main Navigation with Tooltips */}
          <nav className="flex flex-col space-y-1 bg-[#1b1b1b]">
            {navMenuMap.map((main) => (
              <div
                key={`collapsed-${main.itemId}`}
                className="main-item relative group"
              >
                <div
                  className={`flex items-center justify-between w-full text-sm mb-1 cursor-pointer hover:bg-zinc-800 rounded-md relative text-gray-100 ${
                    activeNav === main.itemId
                      ? "bg-orange-700 text-gray-100 font-semibold"
                      : "text-gray-100"
                  }`}
                  onClick={() => toggleSubmenu(`collapsed-${main.itemId}`)}
                >
                  {main.children ? (
                    <span className="flex items-center pl-1">{main.svg}</span>
                  ) : (
                    <span className="flex items-center pl-1">
                      <Link href={main.url}>{main.svg}</Link>
                    </span>
                  )}

                  {/* Tooltip */}
                  <div className="absolute hidden group-hover:flex bg-gray-700 text-gray-100 text-sm p-2 rounded shadow-lg left-full ml-2 whitespace-nowrap z-50">
                    {main.itemName}
                  </div>
                </div>

                {main.children && (
                  <div
                    className={`submenu absolute left-[50px] top-0 pl-2 bg-[#1b1b1b] w-44 p-2 z-50 rounded-md border border-zinc-500 ${
                      activeSubmenu === `collapsed-${main.itemId}`
                        ? "block"
                        : "hidden"
                    }`}
                  >
                    {main.children.map((sub) => (
                      <Link
                        key={sub.itemId}
                        href={sub.url}
                        className={`flex items-center gap-2 text-xs text-gray-100 p-2 mt-1 mb-1 hover:bg-zinc-800 shadow-sm shadow-zinc-700 ${
                          activeSideBar === sub.itemId
                            ? "bg-white !text-black font-bold"
                            : ""
                        }`}
                      >
                        {sub.svg}
                        <span>{sub.itemName}</span>
                      </Link>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </nav>

          {/* Collapsed Bottom Navigation with Tooltips */}
          <div className="flex flex-col space-y-1 bg-[#1b1b1b] p-1 mt-auto">
            <div className="h-px bg-zinc-500 mx-4 my-2"></div>

            <nav className="flex flex-col space-y-1 bg-[#1b1b1b]">
              {bottomMenuMap.map((main) => (
                <div
                  key={`collapsed-${main.itemId}`}
                  className="main-item relative group"
                >
                  <div
                    className={`flex items-center justify-between w-full bottom-0 text-sm mb-1 cursor-pointer hover:bg-zinc-800 rounded-md text-gray-100 ${
                      activeNav === main.itemId
                        ? "bg-orange-700 text-gray-100 font-semibold"
                        : "text-gray-100"
                    }`}
                    onClick={() => toggleSubmenu(`collapsed-${main.itemId}`)}
                  >
                    <span className="flex items-center pl-1">{main.svg}</span>

                    {/* Tooltip */}
                    <div className="absolute hidden group-hover:flex bg-gray-700 text-gray-100 text-sm p-2 rounded shadow-lg left-full ml-2 whitespace-nowrap z-50">
                      {main.itemName}
                    </div>
                  </div>

                  <div
                    className={`submenu absolute left-[50px] p-2 bg-[#1b1b1b] w-44 p-2 z-50 rounded-md border border-zinc-500 ${
                      activeSubmenu === `collapsed-${main.itemId}`
                        ? "block"
                        : "hidden"
                    }`}
                    style={{ bottom: "var(--submenu-bottom, 0)" }}
                  >
                    {main.children.map((sub) => (
                      <Link
                        key={sub.itemId}
                        href={sub.url}
                        className={`flex items-center gap-2 text-xs text-gray-100 p-2 mt-1 mb-1 hover:bg-zinc-800 shadow-sm shadow-zinc-700 ${
                          activeSideBar === sub.itemId
                            ? "bg-white !text-black font-bold"
                            : ""
                        }`}
                      >
                        {sub.svg}
                        <span>{sub.itemName}</span>
                      </Link>
                    ))}
                  </div>
                </div>
              ))}
            </nav>

            {/* Collapsed Hamburger Menu */}
            <nav className="flex flex-col space-y-1 bg-[#1b1b1b]">
              <div className="main-item flex shadow-lg rounded-lg">
                <button
                  id="hamburgerSideButtonCollapsed"
                  className="text-black cursor-pointer"
                  onClick={toggleHamburgerMenu}
                >
                  <svg
                    className="w-8 h-8"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="white"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M4 6h16M4 12h16M4 18h16"
                    ></path>
                  </svg>
                </button>
              </div>

              <div
                className={`absolute bottom-6 left-[50px] bg-[#1b1b1b] text-gray-100 shadow-lg rounded-lg w-48 z-10 divide-y divide-gray-300 rounded-md border border-zinc-500 ${
                  isHamburgerMenuOpen ? "block" : "hidden"
                }`}
              >
                <div className="py-2">
                  <p className="pl-2 block text-sm text-primary100 font-semibold">
                    {userData?.displayName || " "}
                  </p>
                </div>
                <ul className="py-2 text-sm text-gray-100">
                  <li>
                    <Link
                      href="/settings/profile"
                      className="flex p-2 white hover:bg-zinc-800 hover:text-gray-100"
                      suppressHydrationWarning
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        className="w-4 mr-2 h-4"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth="2"
                          d="M5.121 17.804A9.969 9.969 0 0112 15a9.969 9.969 0 016.879 2.804M15 11a3 3 0 11-6 0 3 3 0 016 0zM12 3c4.97 0 9 4.03 9 9s-4.03 9-9 9-9-4.03-9-9 4.03-9 9-9z"
                        />
                      </svg>
                      Profile
                    </Link>
                  </li>
                  <li>
                    <button
                      onClick={handleLogout}
                      className="flex p-2 white hover:bg-zinc-800 hover:text-gray-100"
                      suppressHydrationWarning
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        className="h-4 w-4 mr-2"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        strokeWidth="2"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h6a2 2 0 012 2v1"
                        />
                      </svg>
                      Logout
                    </button>
                  </li>
                </ul>
              </div>
            </nav>
          </div>
        </aside>
      </div>
    </>
  );
};

export default Sidebar;
