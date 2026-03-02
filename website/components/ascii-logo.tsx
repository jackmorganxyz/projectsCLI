export function ASCIILogo() {
  const logo = `┌─┐┌─┐┌─┐ ┐┌─┐┌─┐┌┬┐┌─┐
│─┘│┬┘│ │ │├┤ │   │ └─┐
┴  ┴└─└─┘└┘└─┘└─┘ ┴ └─┘
            CLI`;

  return (
    <pre className="text-xs sm:text-sm md:text-base leading-tight font-mono text-violet select-none whitespace-pre text-center">
      {logo}
    </pre>
  );
}
