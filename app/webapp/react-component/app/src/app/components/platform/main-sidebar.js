"use client";
import { useState, useEffect } from "react";
import Link from "next/link";
import Image from "next/image";
import { usePathname } from "next/navigation";
import { getAuthService } from "../../../../lib/auth";

const Sidebar = ({ userInfo }) => {
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
        {
          itemId: "agents",
          itemName: "Agents",
          url: "/ai-center/agents",
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
              <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
              <circle cx="9" cy="7" r="4" />
              <path d="M22 21v-2a4 4 0 0 0-3-3.87" />
              <path d="M16 3.13a4 4 0 0 1 0 7.75" />
              <circle cx="12" cy="12" r="1" />
              <circle cx="18" cy="12" r="1" />
              <circle cx="6" cy="12" r="1" />
            </svg>
          ),
        },
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
          itemId: "domain",
          itemName: "Domain",
          url: "/settings/domain",
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
          itemId: "platform-Domain",
          itemName: "Platform Domain",
          url: "/settings/platform-domain",
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
  // Add missing activeSubmenu state
  const [activeSubmenu, setActiveSubmenu] = useState("");

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
            <Link href="/ui/">
              <Image
                src="https://i0.wp.com/blog.vapusdata.com/wp-content/uploads/2024/06/Transparent-e1718715526411.png?resize=300%2C88&ssl=1"
                alt="Logo"
                width={136}
                height={88}
                className="h-10 ml-4"
              />
            </Link>
            <button
              className="bg-orange-700 rounded-lg ml-4 relative group inline-block cursor-pointer"
              onClick={toggleSidebar}
              suppressHydrationWarning
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                className="w-6 h-6"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M11 17l-5-5 5-5" />
                <path d="M18 17l-5-5 5-5" />
              </svg>
              <div className="absolute hidden group-hover:flex bg-gray-700 text-gray-100 text-sm p-2 rounded shadow-lg left-full top-1/2 -translate-y-1/2 ml-2 whitespace-nowrap z-50">
                Click here to collapse sidebar
              </div>
            </button>
          </div>
          <div className="h-px bg-zinc-500 mx-4 my-2"></div>

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
                  A
                </p>
                <p className="block text-xs text-primary100 font-semibold text-gray-100 p-[10px]">
                  Anand K
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
                    {userInfo?.displayName || "Anand Kumar"}
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
            <Link href="/ui/">
              <Image
                src="https://i0.wp.com/blog.vapusdata.com/wp-content/uploads/2024/07/cropped-Transparent-12.png?fit=32%2C32&ssl=1"
                alt="Logo"
                className="ml-1"
                width={32}
                height={32}
              />
            </Link>
            <button
              className="bg-orange-700 rounded-lg ml-2 mt-2 cursor-pointer"
              onClick={toggleSidebar}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="w-6 h-6"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M13 17l5-5-5-5" />
                <path d="M6 17l5-5-5-5" />
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
                    {userInfo?.displayName || "Anand Kumar"}
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
