"use client";

interface TerminalWindowProps {
  title?: string;
  children: React.ReactNode;
}

export function TerminalWindow({
  title = "~/projects",
  children,
}: TerminalWindowProps) {
  return (
    <div className="terminal-glow scanlines relative w-full max-w-3xl rounded-xl border border-terminal-border bg-terminal overflow-hidden">
      {/* Title bar */}
      <div className="flex items-center gap-2 px-4 py-3 border-b border-terminal-border bg-terminal">
        <div className="flex gap-2">
          <div className="w-3 h-3 rounded-full bg-red" />
          <div className="w-3 h-3 rounded-full bg-yellow" />
          <div className="w-3 h-3 rounded-full bg-green" />
        </div>
        <span className="flex-1 text-center text-sm font-mono text-muted select-none">
          {title}
        </span>
        <div className="w-[52px]" />
      </div>

      {/* Terminal body */}
      <div className="p-5 md:p-6 font-mono text-sm leading-relaxed overflow-x-auto">
        {children}
      </div>
    </div>
  );
}
