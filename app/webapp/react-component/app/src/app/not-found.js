"use client"
import { useEffect } from 'react';
import Head from 'next/head';

export default function Error404() {

  const goHome = () => {
    window.location.href = '/';
  };

  const goBack = () => {
    if (window.history.length > 1) {
      window.history.back();
    } else {
      goHome();
    }
  };

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

    // Add keyboard navigation
    const handleKeyDown = (e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        goHome();
      }
    };

    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('keydown', handleKeyDown);

    // Cleanup
    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('keydown', handleKeyDown);
      document.body.style.overflow = ''; // Reset on cleanup
    };
  }, []);

  return (
    <>
      <Head>
        <title>404 - Data Not Found</title>
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <script src="https://cdn.tailwindcss.com"></script>
      </Head>

      <div id="not-found-page" className="bg-zinc-800 h-screen w-screen flex items-center justify-center relative">
        
        {/* Grid Background */}
        <div className="absolute inset-0 grid-fade" style={{
          backgroundImage: `
            linear-gradient(rgba(161, 161, 170, 0.1) 1px, transparent 1px),
            linear-gradient(90deg, rgba(161, 161, 170, 0.1) 1px, transparent 1px)
          `,
          backgroundSize: '50px 50px'
        }}></div>
        
        {/* Animated Background Elements */}
        <div className="absolute inset-0">
          <div className="absolute top-1/4 left-1/4 w-2 h-2 bg-zinc-400 rounded-full page-stream"></div>
          <div className="absolute top-1/2 left-1/3 w-1 h-1 bg-zinc-500 rounded-full page-stream" style={{animationDelay: '1s'}}></div>
          <div className="absolute top-3/4 left-1/2 w-3 h-3 bg-zinc-300 rounded-full page-stream" style={{animationDelay: '2s'}}></div>
          <div className="absolute top-1/6 right-1/4 w-2 h-2 bg-zinc-400 rounded-full page-stream" style={{animationDelay: '0.5s'}}></div>
          <div className="absolute top-2/3 right-1/3 w-1 h-1 bg-zinc-500 rounded-full page-stream" style={{animationDelay: '1.5s'}}></div>
        </div>
        
        <div className="text-center z-10 px-6 max-w-4xl mx-auto">
          
          {/* Main 404 Display */}
          <div className="relative mb-8">
            <h1 className="text-6xl sm:text-7xl md:text-8xl lg:text-9xl font-bold mb-4 float-animation" style={{
              background: 'linear-gradient(135deg, #e4e4e7 0%, #a1a1aa 100%)',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
              backgroundClip: 'text'
            }}>404</h1>
          </div>
          
          {/* Error Message */}
          <div className="mb-8">
            <h2 className="text-2xl sm:text-3xl md:text-4xl font-bold text-zinc-100 mb-4">
              Page Not Found
            </h2>
            <p className="text-lg sm:text-xl text-zinc-300 mb-6 leading-relaxed px-4">
              Oops! The page you're looking for seems to have gone offline. 
              <br className="hidden sm:block" />
              Our servers couldn't locate the requested page.
            </p>
          </div>
          
          {/* Action Buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center px-4">
            <button 
              onClick={goHome} 
              className="bg-zinc-700 hover:bg-zinc-600 text-zinc-100 font-semibold py-3 px-8 rounded-lg transform hover:scale-105 transition-all duration-200 shadow-lg hover:shadow-xl border border-zinc-600 hover:border-zinc-500 w-full sm:w-auto"
            >
              <svg className="w-5 h-5 inline-block mr-2" fill="currentColor" viewBox="0 0 20 20">
                <path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"/>
              </svg>
              Return to Home
            </button>
            
            <button 
              onClick={goBack} 
              className="bg-transparent hover:bg-zinc-700 text-zinc-400 hover:text-zinc-100 font-semibold py-3 px-8 rounded-lg transform hover:scale-105 transition-all duration-200 border border-zinc-600 hover:border-zinc-500 w-full sm:w-auto"
            >
              <svg className="w-5 h-5 inline-block mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
              </svg>
              Back
            </button>
          </div>
          
          {/* Footer Message */}
          <div className="mt-8 sm:mt-12 text-zinc-500 text-sm px-4">
            <p>If you believe this is an error, please check your query parameters or contact our support team.</p>
          </div>
          
        </div>
        
        <style jsx>{`
          @keyframes float {
            0%, 100% { transform: translateY(0px) rotate(0deg); }
            25% { transform: translateY(-10px) rotate(1deg); }
            50% { transform: translateY(-5px) rotate(-1deg); }
            75% { transform: translateY(-15px) rotate(2deg); }
          }
          
          @keyframes pulse-glow {
            0%, 100% { box-shadow: 0 0 20px rgba(161, 161, 170, 0.3); }
            50% { box-shadow: 0 0 40px rgba(161, 161, 170, 0.6); }
          }
          
          @keyframes page-stream {
            0% { transform: translateX(-50px); opacity: 0; }
            50% { opacity: 1; }
            100% { transform: translateX(calc(100vw + 50px)); opacity: 0; }
          }
          
          @keyframes grid-fade {
            0%, 100% { opacity: 0.1; }
            50% { opacity: 0.3; }
          }
          
          .float-animation { animation: float 6s ease-in-out infinite; }
          .pulse-glow { animation: pulse-glow 2s ease-in-out infinite; }
          .page-stream { animation: page-stream 3s linear infinite; }
          .grid-fade { animation: grid-fade 4s ease-in-out infinite; }
        `}</style>
      </div>
    </>
  );
}