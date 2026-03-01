"use client";

import { useState, useEffect, useCallback } from "react";

interface Command {
  prompt: string;
  output: string[];
  outputDelay?: number;
  pauseAfter?: number;
}

const COMMANDS: Command[] = [
  {
    prompt: "projects create my-saas-app",
    output: [
      "",
      '\u001b[violet]  Created project "my-saas-app"\u001b[/] — Fresh project, who dis?',
      "  Directory  ~/.projects/projects/my-saas-app",
      "  Created    2025-02-25",
      "",
      '\u001b[muted]Tip: \'projects push <slug>\' handles git init, commit, and GitHub in one step.\u001b[/]',
    ],
    outputDelay: 300,
    pauseAfter: 2200,
  },
  {
    prompt: "projects list",
    output: [
      "",
      "\u001b[header]  Slug              Title              Status    Created       Tags\u001b[/]",
      "  my-saas-app      My SaaS App        active    2025-02-25    next,react",
      "  api-server       API Server         active    2025-02-20    go,api",
      "  mobile-app       Mobile App         paused    2025-02-18    react-native",
      "  cli-tool         CLI Tool           active    2025-02-15    go,cli",
    ],
    outputDelay: 150,
    pauseAfter: 2500,
  },
  {
    prompt: "projects push my-saas-app",
    output: [
      "",
      "\u001b[muted]  Initializing git...\u001b[/]",
      "\u001b[muted]  Staging files...\u001b[/]",
      "\u001b[muted]  Creating commit...\u001b[/]",
      "\u001b[muted]  Creating GitHub repo...\u001b[/]",
      "\u001b[muted]  Pushing to origin...\u001b[/]",
      "",
      "\u001b[emerald]  Shipped it!\u001b[/] github.com/you/my-saas-app",
    ],
    outputDelay: 400,
    pauseAfter: 2200,
  },
  {
    prompt: "projects status",
    output: [
      "",
      "\u001b[violet]  Project Health\u001b[/]",
      "",
      "\u001b[header]  Slug              Status    Git    Remote    Clean\u001b[/]",
      "  my-saas-app      active    yes    yes       \u001b[emerald]clean\u001b[/]",
      "  api-server       active    yes    yes       \u001b[emerald]clean\u001b[/]",
      "  mobile-app       paused    yes    no        \u001b[yellow-text]dirty\u001b[/]",
      "  cli-tool         active    yes    yes       \u001b[emerald]clean\u001b[/]",
    ],
    outputDelay: 100,
    pauseAfter: 2500,
  },
];

function renderStyledLine(line: string) {
  const parts: React.ReactNode[] = [];
  let remaining = line;
  let key = 0;

  while (remaining.length > 0) {
    const match = remaining.match(
      /\u001b\[(violet|emerald|muted|header|yellow-text)\](.*?)\u001b\[\/\]/
    );
    if (!match || match.index === undefined) {
      parts.push(<span key={key++}>{remaining}</span>);
      break;
    }

    if (match.index > 0) {
      parts.push(
        <span key={key++}>{remaining.substring(0, match.index)}</span>
      );
    }

    const colorClass: Record<string, string> = {
      violet: "text-violet",
      emerald: "text-emerald",
      muted: "text-muted",
      header: "text-muted",
      "yellow-text": "text-yellow",
    };

    parts.push(
      <span key={key++} className={colorClass[match[1]] || ""}>
        {match[2]}
      </span>
    );

    remaining = remaining.substring(match.index + match[0].length);
  }

  return parts;
}

type Phase = "typing" | "output" | "pausing";

export function TerminalDemo() {
  const [commandIndex, setCommandIndex] = useState(0);
  const [charIndex, setCharIndex] = useState(0);
  const [outputLines, setOutputLines] = useState<string[]>([]);
  const [phase, setPhase] = useState<Phase>("typing");
  const [history, setHistory] = useState<
    { prompt: string; output: string[] }[]
  >([]);

  const cmd = COMMANDS[commandIndex];

  const advanceCommand = useCallback(() => {
    setHistory((prev) => [
      ...prev,
      { prompt: cmd.prompt, output: outputLines },
    ]);
    const next = (commandIndex + 1) % COMMANDS.length;
    if (next === 0) {
      setHistory([]);
    }
    setCommandIndex(next);
    setCharIndex(0);
    setOutputLines([]);
    setPhase("typing");
  }, [commandIndex, cmd.prompt, outputLines]);

  // Typing phase
  useEffect(() => {
    if (phase !== "typing") return;

    if (charIndex < cmd.prompt.length) {
      const speed = 30 + Math.random() * 40;
      const timeout = setTimeout(() => setCharIndex((c) => c + 1), speed);
      return () => clearTimeout(timeout);
    }

    // Done typing, show output
    const timeout = setTimeout(() => setPhase("output"), 200);
    return () => clearTimeout(timeout);
  }, [phase, charIndex, cmd.prompt.length]);

  // Output phase
  useEffect(() => {
    if (phase !== "output") return;

    if (outputLines.length < cmd.output.length) {
      const delay = cmd.outputDelay || 150;
      const timeout = setTimeout(() => {
        setOutputLines((prev) => [...prev, cmd.output[prev.length]]);
      }, delay);
      return () => clearTimeout(timeout);
    }

    // All output shown, pause
    const timeout = setTimeout(
      () => setPhase("pausing"),
      cmd.pauseAfter || 1500
    );
    return () => clearTimeout(timeout);
  }, [phase, outputLines.length, cmd.output, cmd.outputDelay, cmd.pauseAfter]);

  // Pausing phase
  useEffect(() => {
    if (phase !== "pausing") return;

    const timeout = setTimeout(advanceCommand, 100);
    return () => clearTimeout(timeout);
  }, [phase, advanceCommand]);

  return (
    <div className="text-[13px] md:text-sm">
      {/* History */}
      {history.map((entry, i) => (
        <div key={i} className="mb-1">
          <div>
            <span className="text-violet">$</span>{" "}
            <span className="text-foreground">{entry.prompt}</span>
          </div>
          {entry.output.map((line, j) => (
            <div key={j} className="whitespace-pre">
              {renderStyledLine(line)}
            </div>
          ))}
        </div>
      ))}

      {/* Current command */}
      <div>
        <div>
          <span className="text-violet">$</span>{" "}
          <span className="text-foreground">
            {cmd.prompt.substring(0, charIndex)}
          </span>
          {phase === "typing" && (
            <span className="cursor-blink text-violet">█</span>
          )}
        </div>
        {outputLines.map((line, i) => (
          <div key={i} className="whitespace-pre">
            {renderStyledLine(line)}
          </div>
        ))}
        {phase !== "typing" &&
          outputLines.length >= cmd.output.length && (
            <div className="mt-1">
              <span className="text-violet">$</span>{" "}
              <span className="cursor-blink text-violet">█</span>
            </div>
          )}
      </div>
    </div>
  );
}
