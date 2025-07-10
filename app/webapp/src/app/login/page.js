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
  landingPage = "/dashboard",
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
      toast.info("Please create a organization to continue");
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

  // Add mouse move effect
  useEffect(() => {
    // Prevent body scroll
    document.body.style.overflow = 'hidden';
    
    // Add mouse move particle effects
    const handleMouseMove = (e) => {
      if (Math.random() > 0.9) {
        const cursor = document.createElement('div');
        cursor.className = 'fixed w-1 h-1 bg-zinc-400 rounded-full pointer-events-none opacity-70 z-50';
        cursor.style.left = e.clientX + 'px';
        cursor.style.top = e.clientY + 'px';
        cursor.style.animation = 'page-stream 2s linear forwards';
        document.body.appendChild(cursor);
        
        setTimeout(() => {
          if (cursor.parentNode) {
            cursor.remove();
          }
        }, 2000);
      }
    };

    document.addEventListener('mousemove', handleMouseMove);

    // Cleanup
    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.body.style.overflow = ''; // Reset on cleanup
    };
  }, []);

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
    <>
      <style jsx>{`
        @keyframes page-stream {
          0% { transform: translateX(-50px); opacity: 0; }
          50% { opacity: 1; }
          100% { transform: translateX(calc(100vw + 50px)); opacity: 0; }
        }
        
        @keyframes grid-fade {
          0%, 100% { opacity: 0.1; }
          50% { opacity: 0.3; }
        }
        
        .page-stream { animation: page-stream 3s linear infinite; }
        .grid-fade { animation: grid-fade 4s ease-in-out infinite; }
        
        .gradient-text {
          background: #c21500; 
          background: -webkit-linear-gradient(to right, #ffc500, #c21500); 
          background: linear-gradient(to right, #ffc500, #c21500); 
          
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          background-clip: text;
          color: transparent; 
        }

        .gradient-svg {
          fill: url(#gradient-fill);
        }

        .grid-pattern {
          background-image: 
            linear-gradient(rgba(161, 161, 170, 0.1) 1px, transparent 1px),
            linear-gradient(90deg, rgba(161, 161, 170, 0.1) 1px, transparent 1px);
          background-size: 50px 50px;
        }
      `}</style>

      <div className="bg-zinc-800 h-screen w-screen flex flex-col overflow-hidden relative">
        {/* SVG Gradient Definition */}
        <svg width="0" height="0" style={{position: 'absolute'}}>
          <defs>
            <linearGradient id="gradient-fill" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" style={{stopColor:'#ffc500', stopOpacity:1}} />
              <stop offset="100%" style={{stopColor:'#c21500', stopOpacity:1}} />
            </linearGradient>
          </defs>
        </svg>

        {/* Toast Container */}
        <ToastContainerMessage />

        {/* Animated Background Grid */}
        <div className="absolute inset-0 grid-pattern opacity-30"></div>

        {/* Animated Background Elements */}
        <div className="absolute inset-0">
          <div className="absolute top-1/4 left-1/4 w-2 h-2 bg-zinc-400 rounded-full page-stream"></div>
          <div className="absolute top-1/2 left-1/3 w-1 h-1 bg-zinc-500 rounded-full page-stream" style={{animationDelay: '1s'}}></div>
          <div className="absolute top-3/4 left-1/2 w-3 h-3 bg-zinc-300 rounded-full page-stream" style={{animationDelay: '2s'}}></div>
          <div className="absolute top-1/6 right-1/4 w-2 h-2 bg-zinc-400 rounded-full page-stream" style={{animationDelay: '0.5s'}}></div>
          <div className="absolute top-2/3 right-1/3 w-1 h-1 bg-zinc-500 rounded-full page-stream" style={{animationDelay: '1.5s'}}></div>
        </div>

        {/* Main Content */}
        <main className="flex-grow flex flex-col items-center justify-center p-6 relative z-10">
          {/* Logo and Tagline */}
          <div className="text-center mb-10">
            <Image
              src="https://storage.googleapis.com/vapusdata-public/website-images/vapus-white.webp"
              alt="VapusData Logo"
              width={400}
              height={200}
              className="h-20 w-auto mx-auto mb-4"
              style={{ height: "88px", width: "auto" }}
            />

            <h2 className="text-4xl font-bold gradient-text mb-2">
              Welcome to VapusData
            </h2>
            <p className="text-lg text-zinc-300">
              Secured & Governed AI Data Platform
            </p>
          </div>

          {/* Login Box */}
          <div className="w-full max-w-md bg-white/10 backdrop-blur-[10px] border border-white/20 rounded-2xl shadow-2xl p-8 text-center">
            <h3 className="text-2xl font-semibold gradient-text mb-2">Login</h3>
            <p className="text-gray-100 mb-6">
              Sign in to access your {platform} dashboard and start exploring.
            </p>
            <div className="flex justify-center">
              <button
                onClick={handleLogin}
                className="w-full py-2 px-4 bg-orange-700 hover:opacity-90 text-white font-semibold rounded-xl border-0 transition-all duration-300 transform hover:scale-105 shadow-lg"
              >
                Click here to Login
              </button>
            </div>
          </div>

          {/* Feature Highlights Section */}
          <div className="mt-12 grid grid-cols-1 md:grid-cols-3 gap-6 px-4 max-w-5xl">
            <div className="bg-white/10 backdrop-blur-[10px] border border-white/20 transition-all duration-300 hover:bg-white/15 hover:-translate-y-1 hover:shadow-[0_20px_40px_rgba(0,0,0,0.1)] p-6 rounded-2xl shadow-xl text-center">
              <div className="w-12 h-12 bg-zinc-700 rounded-full flex items-center justify-center mx-auto mb-4 shadow-lg">
                <svg className="w-12 h-12 gradient-svg" viewBox="0 0 24 24">
                  <path d="M12 2a10 10 0 1 0 10 10A10 10 0 0 0 12 2zm0 18a8 8 0 1 1 8-8 8 8 0 0 1-8 8z"></path>
                  <path d="M13 8h-2v6h6v-2h-4V8z"></path>
                </svg>
              </div>
              <h4 className="text-lg font-bold gradient-text">
                Simplified & Governed AI
              </h4>
              <p className="text-zinc-300">
                Get hands-on AI experience instantly with Vapus AI studio.
              </p>
            </div>
            
            <div className="bg-white/10 backdrop-blur-[10px] border border-white/20 transition-all duration-300 hover:bg-white/15 hover:-translate-y-1 hover:shadow-[0_20px_40px_rgba(0,0,0,0.1)] p-6 rounded-2xl shadow-xl text-center">
              <div className="w-12 h-12 bg-zinc-700 rounded-full flex items-center justify-center mx-auto mb-4 shadow-lg">
                <svg className="w-12 h-12 gradient-svg" viewBox="0 0 24 24">
                  <path d="M12 1L3 5V11C3 16.55 6.84 21.74 12 23C17.16 21.74 21 16.55 21 11V5L12 1M12 7C13.4 7 14.8 8.6 14.8 10V11C15.4 11 16 11.4 16 12V16C16 16.6 15.6 17 15 17H9C8.4 17 8 16.6 8 16V12C8 11.4 8.4 11 9 11V10C9 8.6 10.6 7 12 7M12 8.2C11.2 8.2 10.2 9 10.2 10V11H13.8V10C13.8 9 12.8 8.2 12 8.2Z"></path>
                </svg>
              </div>
              <h4 className="text-lg font-bold gradient-text">Secure Data Access</h4>
              <p className="text-zinc-300">Your data is safe with industry-standard security.</p>
            </div>
            
            <div className="bg-white/10 backdrop-blur-[10px] border border-white/20 transition-all duration-300 hover:bg-white/15 hover:-translate-y-1 hover:shadow-[0_20px_40px_rgba(0,0,0,0.1)] p-6 rounded-2xl shadow-xl text-center">
              <div className="w-12 h-12 bg-zinc-700 rounded-full flex items-center justify-center mx-auto mb-4 shadow-lg">
                <svg className="w-12 h-12 gradient-svg" viewBox="0 0 24 24">
                  <path d="M4 6H2v14c0 1.1.9 2 2 2h14v-2H4V6zm16-4H8c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm-1 9H9V9h10v2zm-4 4H9v-2h6v2zm4-8H9V5h10v2z"></path>
                </svg>
              </div>
              <h4 className="text-lg font-bold gradient-text">Flexible Data Management</h4>
              <p className="text-zinc-300">Customize and manage data your way with powerful tools.</p>
            </div>
          </div>
        </main>

        {/* organization Modal */}
        {showModal && (
          <div className="fixed inset-0 bg-zinc-600/90 flex justify-center items-center z-50">
            <div className="bg-white rounded-2xl shadow-2xl w-1/3 h-1/3 p-6 relative">
              <h3 className="text-xl font-semibold text-gray-800 mb-4">
                Create organization
              </h3>
              <input
                id="organizationInput"
                type="text"
                placeholder="Enter organization name"
                className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-500 mb-4"
              />
              <button
                id="createorganization"
                className="w-full py-2 px-4 bg-zinc-700 hover:bg-zinc-600 text-white font-semibold rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-300 transition duration-200"
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
        <footer className="bg-black/30 backdrop-blur-sm text-white/90 py-4 text-center mt-12 border-t border-white/20">
          <p>
            &copy; {new Date().getFullYear()} VapusData Platform. All rights
            reserved.
          </p>
          <p>
            <Link
              href="https://vapusdata.com/privacy-policy"
              target="_blank"
              className="text-white/80 hover:text-orange-300 transition-colors duration-200 hover:underline"
            >
              Privacy Policy
            </Link>{" "}
            |{" "}
            <Link
              href="https://vapusdata.com/terms-and-conditions"
              target="_blank"
              className="text-white/80 hover:text-orange-300 transition-colors duration-200 hover:underline"
            >
              Terms of Service
            </Link>
          </p>
        </footer>
      </div>
    </>
  );
}