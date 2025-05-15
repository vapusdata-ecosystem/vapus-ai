"use client";

import { useState, useEffect, useRef } from "react";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { getAuthService } from "../../../lib/auth";
import Link from "next/link";
import { toast } from "react-toastify";
import ToastContainerMessage from "../components/notification/customToast";

export default function LoginPage({
  platform = "VapusData",
  landingPage = "/settings/domain",
}) {
  const [showModal, setShowModal] = useState(false);
  const router = useRouter();
  const authServiceRef = useRef(null);
  const initialCheckDone = useRef(false);

  // Initialize auth service once
  useEffect(() => {
    authServiceRef.current = getAuthService();
  }, []);

  // Handle URL parameters and authentication check
  useEffect(() => {
    if (!authServiceRef.current) return;

    const params = new URLSearchParams(window.location.search);

    if (params.get("register") === "true") {
      setShowModal(true);
      toast.info("Please create a domain to continue");
    }

    const error = params.get("error");
    if (error) {
      toast.error(decodeURIComponent(error));
    }

    const success = params.get("success");
    if (success) {
      toast.success(decodeURIComponent(success));
    }

    // Only perform the auth check once
    if (!initialCheckDone.current) {
      initialCheckDone.current = true;

      try {
        const isAuthenticated = authServiceRef.current.isAuthenticated();
        if (isAuthenticated) {
          // Get redirect URL from localStorage or use default
          const redirectUrl =
            localStorage.getItem("loginRedirectUrl") || landingPage;
          // Use setTimeout to defer navigation until after render completes
          setTimeout(() => {
            router.push(redirectUrl);
          }, 0);
        }
      } catch (err) {
        console.error("Authentication check error:", err);
        toast.error("Authentication check failed");
      }
    }
  }, [router, landingPage]);

  // Set up token expiration listener separately
  useEffect(() => {
    if (!authServiceRef.current) return;

    const tokenCheckInterval = setInterval(() => {
      if (
        !authServiceRef.current.isAuthenticated() &&
        window.location.pathname !== "/login"
      ) {
        router.push("/login");
      }
    }, 30000);

    return () => {
      clearInterval(tokenCheckInterval);
    };
  }, [router]);

  const closeModal = () => {
    setShowModal(false);
    const newUrl = window.location.pathname;
    window.history.replaceState({}, document.title, newUrl);
  };

  // Api call
  const handleLogin = async (e) => {
    e.preventDefault();
    if (!authServiceRef.current) return;

    try {
      // Use the login method from AuthService
      const result = await authServiceRef.current.login(landingPage);

      if (!result.success) {
        console.error("Login failed:", result.error);
        toast.error(result.error || "Login failed");
      }
    } catch (error) {
      console.error("Login error:", error);
      toast.error(error.message || "Login process failed");
    }
  };

  return (
    <div className="bg-gradient-to-br from-gray-50 to-gray-600 bg-gray-100 flex flex-col min-h-screen ">
      {/* Toast Container */}
      <ToastContainerMessage />

      {/* Main Content */}
      <main className="flex-grow flex flex-col items-center justify-center p-6 ">
        {/* Logo and Tagline */}
        <div className="text-center mb-10">
          <Image
            src="https://storage.googleapis.com/vapusdata-public/website-images/vapus-white.webp"
            alt="VapusData Logo"
            width={400}
            height={200}
            className="h-20 w-auto"
            style={{ height: "88px", width: "auto" }}
          />

          <h2 className="text-4xl font-bold text-gray-800">
            Welcome to VapusData
          </h2>
          <p className="text-lg text-gray-800">
            Secured & Governed AI Data Platform
          </p>
        </div>

        {/* Login Box */}
        <div className="w-full max-w-md bg-white rounded-lg shadow-lg p-8 text-center">
          <h3 className="text-2xl font-semibold text-gray-800 mb-2">Login</h3>
          <p className="text-gray-600 mb-6">
            Sign in to access your {platform} dashboard and start exploring.
          </p>
          <div className="flex justify-center">
            <button
              onClick={handleLogin}
              className="w-full py-2 px-4 bg-black hover:bg-pink-900 text-white font-semibold rounded-lg focus:outline-none focus:ring-2 focus:ring-black focus:ring-opacity-50 transition duration-200"
            >
              Click here to Login
            </button>
          </div>
        </div>

        {/* Feature Highlights Section */}
        <div className="mt-12 grid grid-cols-1 md:grid-cols-3 gap-6 px-4 max-w-5xl">
          <div className="bg-white p-6 rounded-lg shadow-md text-center">
            <svg
              className="w-12 h-12 text-gray-900 mx-auto mb-4"
              fill="currentColor"
              viewBox="0 0 24 24"
            >
              <path d="M12 2a10 10 0 1 0 10 10A10 10 0 0 0 12 2zm0 18a8 8 0 1 1 8-8 8 8 0 0 1-8 8z"></path>
              <path d="M13 8h-2v6h6v-2h-4V8z"></path>
            </svg>
            <h4 className="text-lg font-semibold text-gray-800">
              Simplified & Governed AI
            </h4>
            <p className="text-gray-600">
              Get hands-on AI experience instantly with Vapus AI studio.
            </p>
          </div>
          <div className="bg-white p-6 rounded-lg shadow-md text-center">
            <svg
              className="w-12 h-12 text-gray-900 mx-auto mb-4"
              fill="currentColor"
              viewBox="0 0 24 24"
            >
              <path d="M19.88 18.91l-6.53-5.13 6.53-5.13a1 1 0 0 0 0-1.64 1 1 0 0 0-1.08 0L12 11.82 5.2 7.01a1 1 0 1 0-1.08 1.64L10.6 12l-6.53 5.13a1 1 0 0 0 1.08 1.64L12 12.2l6.8 5.34a1 1 0 0 0 1.08 0 1 1 0 0 0 .08-1.63z"></path>
            </svg>
            <h4 className="text-lg font-semibold text-gray-800">
              Secure Data Access
            </h4>
            <p className="text-gray-600">
              Your data is safe with industry-standard security.
            </p>
          </div>
          <div className="bg-white p-6 rounded-lg shadow-md text-center">
            <svg
              className="w-12 h-12 text-gray-900 mx-auto mb-4"
              fill="currentColor"
              viewBox="0 0 24 24"
            >
              <path d="M3 4h18v2H3zm0 14h18v2H3zm2-9h4v4H5zm14-4h-4v4h4z"></path>
              <path d="M7 11h10v2H7zm0 4h10v2H7z"></path>
            </svg>
            <h4 className="text-lg font-semibold text-gray-800">
              Flexible Data Management
            </h4>
            <p className="text-gray-600">
              Customize and manage data your way with powerful tools.
            </p>
          </div>
        </div>
      </main>

      {/* Domain Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-zinc-600/90 flex justify-center items-center">
          <div className="bg-white rounded-lg shadow-lg w-1/3 h-1/3 p-6 relative">
            <h3 className="text-xl font-semibold text-gray-800 mb-4">
              Create Domain
            </h3>
            <input
              id="domainInput"
              type="text"
              placeholder="Enter domain name"
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-pink-900 mb-4"
            />
            <button
              id="createDomain"
              className="w-full py-2 px-4 bg-black hover:bg-pink-900 text-white font-semibold rounded-lg focus:outline-none focus:ring-2 focus:ring-pink-900 transition duration-200"
            >
              Continue
            </button>
            <button
              onClick={closeModal}
              className="absolute top-2 right-2 text-gray-400 hover:text-gray-600"
            >
              &times;
            </button>
          </div>
        </div>
      )}

      {/* Footer Section */}
      <footer className="bg-black text-gray-200 py-4 text-center mt-12">
        <p>
          &copy; {new Date().getFullYear()} VapusData Platform. All rights
          reserved.
        </p>
        <p>
          <Link
            href="https://vapusdata.com/privacy-policy"
            target="_blank"
            className="text-gray-200 hover:underline"
          >
            Privacy Policy
          </Link>{" "}
          |{" "}
          <Link
            href="https://vapusdata.com/terms-and-conditions"
            target="_blank"
            className="text-gray-200 hover:underline"
          >
            Terms of Service
          </Link>
        </p>
      </footer>
    </div>
  );
}
