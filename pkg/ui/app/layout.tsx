import type { Metadata } from "next";

import "./globals.css";

export const metadata: Metadata = {
  title: "deployment-demo",
  description: "Deployment lifecycle demonstration",
};

export default function RootLayout({
  children,
}: {
  readonly children: React.ReactNode;
}): React.ReactElement {
  return (
    <html lang="en" data-theme="green">
      <body>{children}</body>
    </html>
  );
}
