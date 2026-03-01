import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "projectsCLI — Terminal-native project manager",
  description:
    "Less chaos, more shipping. A terminal-native project manager for humans and AI agents. Scaffold, organize, and push projects from the command line.",
  openGraph: {
    title: "projectsCLI — Terminal-native project manager",
    description:
      "Less chaos, more shipping. A terminal-native project manager for humans and AI agents.",
    url: "https://github.com/jackmorganxyz/projectsCLI",
    siteName: "projectsCLI",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "projectsCLI — Terminal-native project manager",
    description:
      "Less chaos, more shipping. A terminal-native project manager for humans and AI agents.",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link
          rel="preconnect"
          href="https://fonts.gstatic.com"
          crossOrigin="anonymous"
        />
        <link
          href="https://fonts.googleapis.com/css2?family=Geist:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;700&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className="antialiased">{children}</body>
    </html>
  );
}
