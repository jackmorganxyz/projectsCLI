import { Nav } from "@/components/nav";
import { TerminalWindow } from "@/components/terminal-window";
import { ASCIILogo } from "@/components/ascii-logo";
import { TerminalDemo } from "@/components/terminal-demo";
import { Features } from "@/components/features";
import { InstallBlock } from "@/components/install-block";
import { Footer } from "@/components/footer";

export default function Home() {
  return (
    <div className="min-h-screen flex flex-col items-center font-sans">
      {/* Nav */}
      <Nav />

      {/* Hero */}
      <section className="w-full flex flex-col items-center px-6 pt-12 md:pt-20 pb-16">
        <div className="animate-fade-in w-full flex justify-center">
          <TerminalWindow title="projectsCLI">
            <div className="mb-6">
              <ASCIILogo />
            </div>
            <div className="border-t border-terminal-border pt-4">
              <TerminalDemo />
            </div>
          </TerminalWindow>
        </div>
      </section>

      {/* Tagline */}
      <section className="animate-fade-in-delay-1 w-full max-w-2xl mx-auto px-6 pb-12 text-center">
        <h1 className="text-foreground font-sans text-3xl md:text-4xl font-bold mb-4 tracking-tight">
          Less chaos, more shipping.
        </h1>
        <p className="text-muted text-lg md:text-xl leading-relaxed">
          A terminal-native project manager built for humans and AI agents.
          Scaffold, organize, and push projects — all from the command line.
        </p>
      </section>

      {/* CTA Buttons */}
      <section className="animate-fade-in-delay-2 w-full flex flex-wrap justify-center gap-3 px-6 pb-20">
        <a
          href="#install"
          className="inline-flex items-center px-6 py-3 rounded-lg bg-violet text-white font-sans text-sm font-medium hover:bg-violet/90 transition-colors"
        >
          Install
        </a>
        <a
          href="https://github.com/jackmorganxyz/projectsCLI"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center px-6 py-3 rounded-lg border border-terminal-border text-foreground font-sans text-sm font-medium hover:border-violet/50 transition-colors"
        >
          GitHub
        </a>
        <a
          href="https://github.com/jackmorganxyz/projectsCLI/blob/main/README_4_HUMANS.md"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center px-6 py-3 rounded-lg border border-terminal-border text-foreground font-sans text-sm font-medium hover:border-emerald/50 transition-colors"
        >
          Docs 4 Humans
        </a>
        <a
          href="https://github.com/jackmorganxyz/projectsCLI/blob/main/README_4_AGENTS.md"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center px-6 py-3 rounded-lg border border-terminal-border text-foreground font-sans text-sm font-medium hover:border-emerald/50 transition-colors"
        >
          Docs 4 Agents
        </a>
      </section>

      {/* Features */}
      <section id="features" className="w-full pb-20">
        <h2 className="text-foreground font-sans text-2xl md:text-3xl font-bold text-center mb-10">
          What It Does
        </h2>
        <Features />
      </section>

      {/* Install */}
      <section id="install" className="w-full pb-20">
        <InstallBlock />
      </section>

      {/* Quickstart */}
      <section className="w-full max-w-3xl mx-auto px-6 pb-20">
        <h2 className="text-foreground font-sans text-2xl md:text-3xl font-bold text-center mb-8">
          Quickstart
        </h2>
        <div className="rounded-xl border border-terminal-border bg-terminal p-5 md:p-6 font-mono text-sm leading-relaxed">
          <div className="space-y-1">
            <div>
              <span className="text-muted"># Install</span>
            </div>
            <div>
              <span className="text-violet">$</span>{" "}
              <span>brew install jackmorganxyz/tap/projects</span>
            </div>
            <div className="pt-2">
              <span className="text-muted"># Create your first project</span>
            </div>
            <div>
              <span className="text-violet">$</span>{" "}
              <span>projects create my-app</span>
            </div>
            <div className="pt-2">
              <span className="text-muted"># See all your projects</span>
            </div>
            <div>
              <span className="text-violet">$</span>{" "}
              <span>projects list</span>
            </div>
            <div className="pt-2">
              <span className="text-muted">
                # Ship it to GitHub in one command
              </span>
            </div>
            <div>
              <span className="text-violet">$</span>{" "}
              <span>projects push my-app</span>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <Footer />
    </div>
  );
}
