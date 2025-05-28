"use client";
import Header from "../components/platform/header";

export default function Dashboard({ backListingLink = "./" }) {
  return (
    <div className="bg-zinc-800 flex h-screen">
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="Dashboard"
          hideBackListingLink={true}
          backListingLink={backListingLink}
        />

        <div className="flex-grow p-4 overflow-y-auto scrollbar text-gray-100">
          <div className="container mx-auto p-2 space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-4 gap-2">
              <div className="border border-zinc-500 rounded-lg shadow-lg p-2 flex items-center justify-center">
                <div>
                  <h2 className="text-xl font-semibold text-center">NA</h2>
                </div>
              </div>

              <div className="border border-zinc-500 rounded-lg shadow-lg p-2 flex items-center justify-center">
                <div>
                  <h2 className="text-xl font-semibold text-center">
                    Data Sources
                  </h2>
                </div>
              </div>
              <div className="border border-zinc-500 rounded-lg shadow-lg p-2 flex items-center justify-center">
                <div>
                  <h2 className="text-xl font-semibold text-center">NA</h2>
                </div>
              </div>
              <div className="border border-zinc-500 rounded-lg shadow-lg p-2 flex items-center justify-center">
                <div>
                  <h2 className="text-xl font-semibold text-center">Users</h2>
                </div>
              </div>
              <div className="border border-zinc-500 rounded-lg shadow-lg p-2 flex items-center justify-center col-span-2">
                <div>
                  <h2 className="text-xl font-semibold text-center">
                    Inbound Data Queries
                  </h2>
                  <table className="table-auto w-full mt-2 text-xs border-2 border-zinc-500">
                    <thead className="bg-zinc-900 divide-y divide-zinc-500 break-words text-sm font-medium text-gray-500 uppercase">
                      <tr>
                        <th className="px-4 py-2">Failed</th>
                        <th className="px-4 py-2">Successfull</th>
                        <th className="px-4 py-2">Total</th>
                      </tr>
                    </thead>
                    <tbody className="bg-zinc-800 divide-y divide-zinc-900 text-sm">
                      <tr>
                        <td className="px-4 py-2 text-yellow-500">0%</td>
                        <td className="px-4 py-2 text-yellow-500">0%</td>
                        <td className="px-4 py-2 text-yellow-500">0%</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
              <div className="border border-zinc-500 rounded-lg shadow-lg p-2 flex items-center justify-center col-span-2">
                <div>
                  <h2 className="text-xl font-semibold text-center">
                    Outbound Data Queries
                  </h2>
                  <table className="table-auto w-full mt-2 text-xs border-2 border-zinc-500">
                    <thead className="bg-zinc-900 divide-y divide-zinc-500 break-words text-sm font-medium text-gray-500 uppercase">
                      <tr>
                        <th className="px-4 py-2">Failed</th>
                        <th className="px-4 py-2">Successfull</th>
                        <th className="px-4 py-2">Total</th>
                      </tr>
                    </thead>
                    <tbody className="bg-zinc-800 divide-y divide-zinc-900 text-sm">
                      <tr>
                        <td className="px-4 py-2 text-yellow-500">0%</td>
                        <td className="px-4 py-2 text-yellow-500">0%</td>
                        <td className="px-4 py-2 text-yellow-500">0%</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-3 gap-2">
              <div className="rounded-lg shadow-lg p-2 border border-zinc-500">
                <h3 className="text-lg font-semibold mb-4">
                  Dataproducts Heat Map
                </h3>
                <canvas id="dataproductMetrics" className="w-full"></canvas>
              </div>
              <div className="rounded-lg shadow-lg p-2 border border-zinc-500">
                <h3 className="text-lg font-semibold mb-4">
                  Data Sources Compliance Heat Map
                </h3>
                <canvas id="dataSourceCompliances" className="w-full"></canvas>
              </div>

              <div className="rounded-lg shadow-lg p-2 border border-zinc-500">
                <h3 className="text-lg font-semibold mb-4">VDC Deployments </h3>
                <canvas id="deploymentMetrics" className="w-full"></canvas>
              </div>
              {/* <div className="rounded-lg shadow-lg p-2">
        <h3 className="text-lg font-semibold mb-4">Users Growth</h3>
        <canvas id="userMetrics" className="w-full"></canvas>
      </div>  */}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
