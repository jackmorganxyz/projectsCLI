export function Footer() {
  return (
    <footer className="w-full max-w-5xl mx-auto px-6 py-12 border-t border-terminal-border">
      <div className="flex flex-col sm:flex-row items-center justify-between gap-4 text-muted text-sm">
        <span className="font-mono">projectsCLI</span>
        <div className="flex items-center gap-6">
          <a
            href="https://github.com/jackmorganxyz/projectsCLI"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-foreground transition-colors"
          >
            GitHub
          </a>
          <a
            href="https://github.com/jackmorganxyz/projectsCLI/blob/main/README_4_HUMANS.md"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-foreground transition-colors"
          >
            Docs 4 Humans
          </a>
          <a
            href="https://github.com/jackmorganxyz/projectsCLI/blob/main/README_4_AGENTS.md"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-foreground transition-colors"
          >
            Docs 4 Agents
          </a>
          <a
            href="https://github.com/jackmorganxyz/projectsCLI/blob/main/LICENSE"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-foreground transition-colors"
          >
            MIT License
          </a>
        </div>
      </div>
    </footer>
  );
}
