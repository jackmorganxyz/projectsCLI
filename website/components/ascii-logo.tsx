export function ASCIILogo() {
  const logo = `                    ╔═══════════════════════════════════════╗
                    ║                                       ║
 ┌─────────────┐    ║   ┌─┐┌─┐┌─┐ ┐┌─┐┌─┐┌┬┐┌─┐          ║
 │ ▓▓▓▓▓▓▓▓▓▓▓ │    ║   │─┘│┬┘│ │ │├┤ │   │ └─┐          ║
 │ ▓ project ▓ │    ║   ┴  ┴└─└─┘└┘└─┘└─┘ ┴ └─┘          ║
 │ ▓▓▓▓▓▓▓▓▓▓▓ │    ║                  CLI                 ║
 │  ═══════════ │    ║                                       ║
 │  my-saas-app │    ║   less chaos, more shipping.          ║
 │  api-server  │    ║   built for humans and agents.        ║
 │  mobile-app  │    ║                                       ║
 │  cli-tool    │    ╚═══════════════════════════════════════╝
 └─────────────┘`;

  return (
    <pre className="text-[8px] sm:text-[10px] md:text-xs leading-tight font-mono text-violet select-none whitespace-pre overflow-x-auto">
      {logo}
    </pre>
  );
}
