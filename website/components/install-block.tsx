"use client";

import { useState } from "react";

const METHODS = [
  {
    label: "Homebrew",
    command: "brew install jackmorganxyz/tap/projects",
  },
  {
    label: "Shell",
    command:
      "curl -sSL https://raw.githubusercontent.com/jackmorganxyz/projectsCLI/main/install.sh | sh",
  },
  {
    label: "Source",
    command:
      "git clone https://github.com/jackmorganxyz/projectsCLI.git && cd projectsCLI && make install",
  },
];

export function InstallBlock() {
  const [activeTab, setActiveTab] = useState(0);
  const [copied, setCopied] = useState(false);

  function handleCopy() {
    navigator.clipboard.writeText(METHODS[activeTab].command);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }

  return (
    <section className="w-full max-w-3xl mx-auto px-6">
      <h2 className="text-foreground font-sans text-2xl md:text-3xl font-bold text-center mb-8">
        Get Started
      </h2>

      <div className="rounded-xl border border-terminal-border bg-terminal overflow-hidden">
        {/* Tabs */}
        <div className="flex border-b border-terminal-border">
          {METHODS.map((method, i) => (
            <button
              key={method.label}
              onClick={() => {
                setActiveTab(i);
                setCopied(false);
              }}
              className={`px-5 py-3 text-sm font-mono transition-colors ${
                activeTab === i
                  ? "text-violet border-b-2 border-violet bg-violet/5"
                  : "text-muted hover:text-foreground"
              }`}
            >
              {method.label}
            </button>
          ))}
        </div>

        {/* Command */}
        <div className="p-5 flex items-center justify-between gap-4">
          <code className="text-sm font-mono text-foreground break-all">
            <span className="text-violet">$</span>{" "}
            {METHODS[activeTab].command}
          </code>
          <button
            onClick={handleCopy}
            className="shrink-0 text-muted hover:text-foreground transition-colors text-sm font-mono"
            title="Copy to clipboard"
          >
            {copied ? (
              <span className="text-emerald">copied!</span>
            ) : (
              <span>[copy]</span>
            )}
          </button>
        </div>
      </div>
    </section>
  );
}
