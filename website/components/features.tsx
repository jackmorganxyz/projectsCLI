const FEATURES = [
  {
    icon: "┌─┐\n│▓│\n└─┘",
    title: "Scaffold in Seconds",
    description:
      "One command creates a full project structure with docs, tasks, context, and code directories. Opinionated so you don't have to think.",
  },
  {
    icon: "╔═╗\n║$║\n╚═╝",
    title: "One-Command GitHub",
    description:
      "projects push handles git init, staging, commit, repo creation, and push. Zero to GitHub in one command.",
  },
  {
    icon: "┌──┐\n│{}│\n└──┘",
    title: "Human + Agent",
    description:
      "Beautiful TUI for humans, clean JSON output for agents and scripts. Auto-detects pipes and switches modes.",
  },
  {
    icon: "┌──┐\n│//│\n└──┘",
    title: "Multi-Account",
    description:
      "Organize projects by GitHub account with folders. Work and personal repos, neatly separated.",
  },
];

export function Features() {
  return (
    <section className="w-full max-w-4xl mx-auto px-6">
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        {FEATURES.map((feature) => (
          <div
            key={feature.title}
            className="feature-card rounded-lg border border-terminal-border bg-terminal p-6"
          >
            <pre className="text-violet text-xs leading-tight mb-3 select-none">
              {feature.icon}
            </pre>
            <h3 className="text-foreground font-sans text-lg font-semibold mb-2">
              {feature.title}
            </h3>
            <p className="text-muted text-sm leading-relaxed">
              {feature.description}
            </p>
          </div>
        ))}
      </div>
    </section>
  );
}
