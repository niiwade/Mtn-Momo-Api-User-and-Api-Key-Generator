"use client";

import { useState } from "react";

interface FormData {
  primaryKey: string;
  secondaryKey: string;
  callbackHost: string;
}

interface GeneratedKeys {
  apiKey: string;
  apiUser: string;
  userId: string;
  callbackHost: string;
  dateTime: string;
  targetEnvironment: string;
  testCommand?: string; // Optional test curl command
  base64Auth?: string; // Base64 encoded auth string (apiUser:apiKey)
}

interface ApiResponse {
  success: boolean;
  message: string;
  data: GeneratedKeys;
}

export default function Home() {
  const [formData, setFormData] = useState<FormData>({
    primaryKey: "",
    secondaryKey: "",
    callbackHost: "",
  });
  const [generatedKeys, setGeneratedKeys] = useState<GeneratedKeys | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");
  const [statusMessage, setStatusMessage] = useState<string | null>(null);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    setStatusMessage("");
    
    try {
      const response = await fetch("http://localhost:8080/api/generate", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          primaryKey: formData.primaryKey,
          secondaryKey: formData.secondaryKey,
          callbackHost: formData.callbackHost,
        }),
      });
      
      const data: ApiResponse = await response.json();
      
      if (!data.success) {
        throw new Error(data.message || "Failed to generate API credentials");
      }
      
      setGeneratedKeys(data.data);
      setStatusMessage(data.message);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "An error occurred while generating API credentials");
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  return (
    <main className="flex min-h-screen flex-col items-center p-8 bg-gradient-to-b from-yellow-50 to-orange-100">
      <div className="w-full max-w-4xl">
        <div className="text-center mb-10">
          <h1 className="text-4xl font-bold text-yellow-600 mb-2">MTN MoMo API Key Generator</h1>
          <p className="text-gray-600">Generate API Keys and API Users for MTN Mobile Money Integration</p>
        </div>

        <div className="bg-white rounded-lg shadow-lg p-8 mb-8">
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="primaryKey" className="block text-sm font-medium text-gray-700 mb-1">
                Subscription Key (Primary Key) <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                id="primaryKey"
                name="primaryKey"
                value={formData.primaryKey}
                onChange={handleInputChange}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-yellow-500"
                placeholder="Enter your MTN MoMo subscription key"
              />
              <p className="text-xs text-gray-500 mt-1">The subscription key from your MTN MoMo developer account</p>
            </div>
            
            <div>
              <label htmlFor="secondaryKey" className="block text-sm font-medium text-gray-700 mb-1">
                Secondary Key (Optional)
              </label>
              <input
                type="text"
                id="secondaryKey"
                name="secondaryKey"
                value={formData.secondaryKey}
                onChange={handleInputChange}
                className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-yellow-500"
                placeholder="Enter your secondary key (if applicable)"
              />
            </div>

            <div>
              <label htmlFor="callbackHost" className="block text-sm font-medium text-gray-700 mb-1">
                Provider Callback Host
              </label>
              <input
                type="text"
                id="callbackHost"
                name="callbackHost"
                value={formData.callbackHost}
                onChange={handleInputChange}
                className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-yellow-500"
                placeholder="example.com"
              />
              <p className="text-xs text-gray-500 mt-1">Your application&apos;s callback host (e.g., example.com)</p>
            </div>
            
            <div>
              <button
                type="submit"
                disabled={loading}
                className={`w-full py-3 px-4 bg-yellow-500 hover:bg-yellow-600 text-white font-semibold rounded-md transition duration-200 ${
                  loading ? "opacity-70 cursor-not-allowed" : ""
                }`}
              >
                {loading ? "Generating..." : "Generate API User & Key"}
              </button>
            </div>
          </form>
        </div>

        {error && (
          <div className="bg-red-50 border-l-4 border-red-500 p-4 mb-8 rounded-md">
            <div className="flex">
              <div className="ml-3">
                <p className="text-sm text-red-700">{error}</p>
              </div>
            </div>
          </div>
        )}

        {generatedKeys && (
          <div className="bg-white rounded-lg shadow-lg p-8">
            <h2 className="text-2xl font-bold text-gray-800 mb-4">Generated MTN MoMo API Credentials</h2>
            
            {statusMessage && (
              <div className={`p-4 mb-6 rounded-md ${statusMessage.includes("registered") ? "bg-green-50 border-l-4 border-green-400" : "bg-yellow-50 border-l-4 border-yellow-400"}`}>
                <div className="flex">
                  <div className="ml-3">
                    <p className="text-sm font-medium text-gray-700">{statusMessage}</p>
                  </div>
                </div>
              </div>
            )}
            
            <div className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">API Key</label>
                <div className="flex">
                  <input
                    type="text"
                    readOnly
                    value={generatedKeys.apiKey}
                    className="flex-grow px-4 py-2 bg-gray-50 border border-gray-300 rounded-l-md focus:outline-none"
                  />
                  <button
                    onClick={() => copyToClipboard(generatedKeys.apiKey)}
                    className="px-4 py-2 bg-gray-200 hover:bg-gray-300 border border-gray-300 rounded-r-md transition duration-200"
                  >
                    Copy
                  </button>
                </div>
                <p className="text-xs text-gray-500 mt-1">Use this as your API Key for authentication</p>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">API User (X-Reference-Id)</label>
                <div className="flex">
                  <input
                    type="text"
                    readOnly
                    value={generatedKeys.apiUser}
                    className="flex-grow px-4 py-2 bg-gray-50 border border-gray-300 rounded-l-md focus:outline-none"
                  />
                  <button
                    onClick={() => copyToClipboard(generatedKeys.apiUser)}
                    className="px-4 py-2 bg-gray-200 hover:bg-gray-300 border border-gray-300 rounded-r-md transition duration-200"
                  >
                    Copy
                  </button>
                </div>
                <p className="text-xs text-gray-500 mt-1">Use this as your API User ID or X-Reference-Id in API calls</p>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Callback Host</label>
                <div className="flex">
                  <input
                    type="text"
                    readOnly
                    value={generatedKeys.callbackHost}
                    className="flex-grow px-4 py-2 bg-gray-50 border border-gray-300 rounded-md focus:outline-none"
                  />
                </div>
                <p className="text-xs text-gray-500 mt-1">Your registered callback host</p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Target Environment</label>
                <div className="flex">
                  <input
                    type="text"
                    readOnly
                    value={generatedKeys.targetEnvironment}
                    className="flex-grow px-4 py-2 bg-gray-50 border border-gray-300 rounded-md focus:outline-none"
                  />
                </div>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Generated At</label>
                <input
                  type="text"
                  readOnly
                  value={new Date(generatedKeys.dateTime).toLocaleString()}
                  className="w-full px-4 py-2 bg-gray-50 border border-gray-300 rounded-md focus:outline-none"
                />
              </div>
              
              {generatedKeys.base64Auth && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Base64 Encoded Auth String</label>
                  <div className="flex">
                    <input
                      type="text"
                      readOnly
                      value={generatedKeys.base64Auth}
                      className="flex-grow px-4 py-2 bg-gray-50 border border-gray-300 rounded-l-md focus:outline-none font-mono text-sm"
                    />
                    <button
                      onClick={() => copyToClipboard(generatedKeys.base64Auth || "")}
                      className="px-4 py-2 bg-gray-200 hover:bg-gray-300 border border-gray-300 rounded-r-md transition duration-200"
                    >
                      Copy
                    </button>
                  </div>
                  <p className="text-xs text-gray-500 mt-1">Use this in your Authorization header: <code className="bg-gray-100 px-1 py-0.5 rounded">Authorization: Basic {generatedKeys.base64Auth}</code></p>
                </div>
              )}

              <div className="p-4 bg-yellow-50 border-l-4 border-yellow-400 rounded-md">
                <div className="flex">
                  <div className="ml-3">
                    <p className="text-sm text-yellow-700">
                      <strong>Important:</strong> Store these credentials securely. The API Key cannot be retrieved again if lost.
                    </p>
                  </div>
                </div>
              </div>
              
              {generatedKeys.testCommand && (
                <div className="mt-6 border-t pt-6">
                  <h3 className="text-lg font-medium text-gray-900 mb-3">Test Your Credentials</h3>
                  <div className="bg-gray-800 rounded-md p-4 overflow-x-auto">
                    <pre className="text-green-400 text-sm whitespace-pre-wrap">{generatedKeys.testCommand}</pre>
                  </div>
                  <button
                    onClick={() => copyToClipboard(generatedKeys.testCommand || "")}
                    className="mt-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-md transition duration-200 text-sm flex items-center"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" />
                    </svg>
                    Copy Test Command
                  </button>
                  <p className="text-xs text-gray-500 mt-2">Run this command to test if your credentials work with the MTN MoMo API</p>
                </div>
              )}
            </div>
          </div>
        )}
        
        <div className="mt-10 text-center text-sm text-gray-600">
          <p>© {new Date().getFullYear()} MTN MoMo API Key Generator | Built with Next.js and Go</p>
        </div>
      </div>
    </main>
  );
}
