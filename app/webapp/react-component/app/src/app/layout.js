"use client";
import { Geist, Geist_Mono } from "next/font/google";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import "./globals.css";
import Sidebar from "./components/platform/main-sidebar";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

//  Loading Component
const ModernLoader = () => {
  return (
    <div className="flex items-center justify-center min-h-screen bg-[#1b1b1b]">
      <div className="flex flex-col items-center space-y-6">
        <div className="relative">
          <div className="w-16 h-16 border-4 border-slate-200 dark:border-slate-700 rounded-full animate-pulse"></div>
          <div className="absolute top-0 left-0 w-16 h-16 border-4 border-transparent border-t-orange-700 border-r-orange-700 rounded-full animate-spin"></div>
          <div className="absolute top-1/2 left-1/2 w-2 h-2 bg-orange-700 rounded-full transform -translate-x-1/2 -translate-y-1/2 animate-ping"></div>
        </div>

        {/* Loading text with typewriter effect */}
        <div className="flex items-center space-x-1">
          <span className="text-slate-600 dark:text-slate-400 font-medium">
            Loading
          </span>
          <div className="flex space-x-1">
            <div
              className="w-1 h-1 bg-orange-700 rounded-full animate-bounce"
              style={{ animationDelay: "0ms" }}
            ></div>
            <div
              className="w-1 h-1 bg-orange-700 rounded-full animate-bounce"
              style={{ animationDelay: "150ms" }}
            ></div>
            <div
              className="w-1 h-1 bg-orange-700 rounded-full animate-bounce"
              style={{ animationDelay: "300ms" }}
            ></div>
          </div>
        </div>

        {/* Progress bar */}
        <div className="w-48 h-1 bg-slate-200 dark:bg-slate-700 rounded-full overflow-hidden">
          <div
            className="h-full bg-gradient-to-r from-orange-700 to-orange-500 rounded-full animate-pulse"
            style={{
              width: "60%",
              animation: "loading-progress 2s ease-in-out infinite",
            }}
          ></div>
        </div>
      </div>

      <style jsx>{`
        @keyframes loading-progress {
          0% {
            width: 0%;
          }
          50% {
            width: 70%;
          }
          100% {
            width: 100%;
          }
        }
      `}</style>
    </div>
  );
};

export default function RootLayout({ children }) {
  const pathname = usePathname();
  const router = useRouter();
  const [isAuthenticated, setIsAuthenticated] = useState(null);
  const [isNotFoundPage, setIsNotFoundPage] = useState(false);

  const isLoginPage = pathname.startsWith("/login");
  const hideSidebar = isLoginPage || isNotFoundPage;

  // Get cookie value
  const getCookie = (name) => {
    if (typeof document === "undefined") return null;
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    return parts.length === 2 ? parts.pop().split(";").shift() : null;
  };

  // Authentication check
  useEffect(() => {
    if (isLoginPage) {
      setIsAuthenticated(true);
      return;
    }

    const token = getCookie("access_token");
    if (!token) {
      router.replace("/login");
    } else {
      setIsAuthenticated(true);
    }
  }, [pathname, isLoginPage, router]);

  useEffect(() => {
    const checkFor404 = () => {
      // Check if the current page has a not-found indicator
      const notFoundElement =
        document.getElementById("not-found-page") ||
        document.querySelector('[data-testid="not-found"]') ||
        document.querySelector(".not-found") ||
        document.querySelector('[class*="not-found"]') ||
        document.querySelector('[class*="404"]');

      // Check if page title indicates 404
      const titleIndicates404 =
        document.title.toLowerCase().includes("404") ||
        document.title.toLowerCase().includes("not found");

      // Check URL patterns that might indicate 404
      const urlPattern404 = /\/(404|not-found)$/i.test(pathname);

      return !!(notFoundElement || titleIndicates404 || urlPattern404);
    };

    setIsNotFoundPage(checkFor404());

    // Use MutationObserver to detect DOM changes
    const observer = new MutationObserver(() => {
      setIsNotFoundPage(checkFor404());
    });

    observer.observe(document.body, {
      childList: true,
      subtree: true,
      attributes: true,
      attributeFilter: ["class", "id", "data-testid"],
    });

    // Also check after a small delay to catch delayed renders
    const timeoutId = setTimeout(() => {
      setIsNotFoundPage(checkFor404());
    }, 100);

    return () => {
      observer.disconnect();
      clearTimeout(timeoutId);
    };
  }, [pathname]);

  // Show  loading while checking authentication
  if (isAuthenticated === null && !isLoginPage) {
    return (
      <html lang="en">
        <head>
          <link
            rel="stylesheet"
            href="https://cdn.datatables.net/1.13.1/css/jquery.dataTables.min.css"
          />
        </head>
        <body
          className={`${geistSans.variable} ${geistMono.variable} antialiased`}
        >
          <ModernLoader />
        </body>
      </html>
    );
  }

  return (
    <html lang="en">
      <head>
        <link
          rel="stylesheet"
          href="https://cdn.datatables.net/1.13.1/css/jquery.dataTables.min.css"
        />
      </head>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <div className="flex min-h-screen">
          {!hideSidebar && <Sidebar />}
          <main className="flex-1">{children}</main>
        </div>
      </body>
    </html>
  );
}
