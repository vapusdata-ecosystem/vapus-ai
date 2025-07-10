"use client";
import Link from "next/link";
import { useEffect, useState } from "react";

export default function CreateNewButton({
  label = "Add New",
  href = "#",
  className = "",
  position = "top-left",
}) {
  const [isSidebarExpanded, setIsSidebarExpanded] = useState(true);

  useEffect(() => {
    const checkSidebarState = () => {
      // Check if expanded sidebar is visible
      const expandedSidebar = document.querySelector('aside.w-\\[208px\\]');
      const collapsedSidebar = document.querySelector('aside.w-\\[40px\\]');
      
      if (expandedSidebar && !expandedSidebar.classList.contains('hidden')) {
        setIsSidebarExpanded(true);
      } else if (collapsedSidebar && !collapsedSidebar.classList.contains('hidden')) {
        setIsSidebarExpanded(false);
      }
    };

    // Check initially
    checkSidebarState();

    // Create a MutationObserver to watch for class changes
    const observer = new MutationObserver(() => {
      checkSidebarState();
    });

    // Observe the document for changes
    observer.observe(document.body, {
      childList: true,
      subtree: true,
      attributes: true,
      attributeFilter: ['class']
    });

    // Cleanup
    return () => observer.disconnect();
  }, []);

  const getPositionClasses = () => {
    const baseClasses = "fixed z-[1000] text-white text-sm px-3 py-2 rounded-lg bg-orange-700 hover:bg-pink-900 text-lg flex items-center shadow-lg transition-all duration-300";
    
    switch(position) {
      case "top-left":
        // Use exact pixel values for better alignment
        const leftStyle = isSidebarExpanded 
          ? "left-[224px]" // 208px sidebar + 16px gap
          : "left-[56px]";  // 40px sidebar + 16px gap
        return `${baseClasses} top-4 ${leftStyle}`;
      case "top-center":
        return `${baseClasses} top-4 left-1/2 transform -translate-x-1/2`;
      case "top-right":
      default:
        return `${baseClasses} top-4 right-4`;
    }
  };

  return (
    <Link href={href} className={`${getPositionClasses()} ${className}`}>
      <svg
        className="w-4 h-4 mr-1"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          fill="#fff"
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth="2"
          d="M12 4v16m8-8H4"
        ></path>
      </svg>
      {label}
    </Link>
  );
}