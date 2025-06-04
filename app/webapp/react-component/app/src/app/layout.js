"use client";
import { Geist, Geist_Mono } from "next/font/google";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import "./globals.css";
import Sidebar from "./components/platform/main-sidebar";
import LoadingOverlay from "./components/loading/loading";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

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
      const notFoundElement =
        document.getElementById("not-found-page") ||
        document.querySelector('[data-testid="not-found"]') ||
        document.querySelector(".not-found") ||
        document.querySelector('[class*="not-found"]') ||
        document.querySelector('[class*="404"]');

      const titleIndicates404 =
        document.title.toLowerCase().includes("404") ||
        document.title.toLowerCase().includes("not found");

      const urlPattern404 = /\/(404|not-found)$/i.test(pathname);

      return !!(notFoundElement || titleIndicates404 || urlPattern404);
    };

    setIsNotFoundPage(checkFor404());

    const observer = new MutationObserver(() => {
      setIsNotFoundPage(checkFor404());
    });

    observer.observe(document.body, {
      childList: true,
      subtree: true,
      attributes: true,
      attributeFilter: ["class", "id", "data-testid"],
    });

    const timeoutId = setTimeout(() => {
      setIsNotFoundPage(checkFor404());
    }, 100);

    return () => {
      observer.disconnect();
      clearTimeout(timeoutId);
    };
  }, [pathname]);

  // Show loading while checking authentication
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
          <LoadingOverlay isLoading={true}/>
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