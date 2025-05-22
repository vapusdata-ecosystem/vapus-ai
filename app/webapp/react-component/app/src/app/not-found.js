import Link from "next/link";
import Header from "./components/platform/header";

export default function Error404() {
  return (
    <div id="not-found-page" className="bg-gray-200 flex flex-col h-screen">
      <Header hideBackListingLink={false} backListingLink="./" />

      <div className="flex justify-center items-center h-screen">
        <div className="bg-white p-6 rounded-lg shadow-lg text-center">
          <h1 className="text-4xl font-bold mb-4 text-red-900">
            404 Not Found
          </h1>
          <p className="text-lg mb-4 text-red-900">
            Content not found. Please check the URL and try again.
          </p>
          <Link
            href="/"
            className="bg-black text-white px-4 py-2 rounded hover:bg-pink-900"
          >
            Go to Home
          </Link>
        </div>
      </div>
    </div>
  );
}
