"use client";
import { Geist, Geist_Mono } from "next/font/google";
import { usePathname } from "next/navigation";
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

export default function RootLayout({ children }) {
  const pathname = usePathname();
  const [isNotFoundPage, setIsNotFoundPage] = useState(false);

  useEffect(() => {
    const el = document.getElementById("not-found-page");
    setIsNotFoundPage(!!el);
  }, [pathname]);

  const hideSidebar = pathname.startsWith("/login") || isNotFoundPage;

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
