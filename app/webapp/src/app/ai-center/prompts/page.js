// AIPromptsPage.jsx
"use client";
import { useState, useEffect } from "react";
import Header from "../../components/platform/header";
import CreateNewButton from "@/app/components/add-new-button";
import Card from "@/app/components/card";
import { PromptsApi } from "@/app/utils/ai-studio-endpoint/prompts-api";
import LoadingOverlay from "@/app/components/loading/loading";

export default function AIPromptsPage({ backListingLink = "./" }) {
  const [aiPrompts, setAiPrompts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchPrompts = async () => {
      try {
        const response = await PromptsApi.getPrompts();
        const prompts = response.output || response.data || response || [];
        setAiPrompts(prompts);
      } catch (error) {
        console.error("Error fetching prompts:", error);
        setError(error.message);
      } finally {
        setLoading(false);
      }
    };

    fetchPrompts();
  }, []);

  if (loading) {
    return (
      <div className="bg-zinc-800 flex h-screen justify-center items-center relative">
          <LoadingOverlay 
                isLoading={loading} 
                text="Loading plugin details"
                size="default"
                isOverlay={true}
                className="absolute inset-0 z-10 bg-zinc-800"
              />
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-zinc-800 flex h-screen justify-center items-center">
        <div className="text-red-500 text-xl">Error: {error}</div>
      </div>
    );
  }

  return (
    <div className="bg-zinc-800 flex h-screen">
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="AI Model Prompts"
          hideBackListingLink={true}
          backListingLink={backListingLink}
        />

        <div className="flex-grow p-2 w-full">
          <div className="flex justify-between mb-2 items-center p-2">
            <CreateNewButton href="./prompts/create" label="Add New" />
          </div>

          <section id="grids" className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
              {aiPrompts && aiPrompts.length > 0 ? (
                aiPrompts.map((prompt) => (
                  <Card
                    key={prompt.promptId}
                    prompt={{
                      PromptId: prompt.promptId,
                      Name: prompt.name,
                      PromptTypes: prompt.promptTypes,
                      Labels: prompt.labels,
                      PromptOwner: prompt.promptOwner || "Not specified",
                      ResourceBase: {
                        Status: prompt.resourceBase?.status || "INACTIVE",
                      },
                    }}
                    backListingLink={backListingLink}
                  />
                ))
              ) : (
                <div className="text-center text-gray-100 m-20 text-4xl col-span-4">
                  <p>No AI Prompts available.</p>
                </div>
              )}
            </div>
          </section>
        </div>
      </div>
    </div>
  );
}
