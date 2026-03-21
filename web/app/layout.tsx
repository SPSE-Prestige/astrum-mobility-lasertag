import type { Metadata } from "next";
import localFont from "next/font/local";
import { Noto_Sans } from "next/font/google";
import "./globals.css";

const headline = localFont({
  variable: "--font-goldman",
  display: "swap",
  src: [
    {
      path: "./fonts/Goldman/Goldman-Regular.ttf",
      weight: "400",
      style: "normal",
    },
    {
      path: "./fonts/Goldman/Goldman-Bold.ttf",
      weight: "700",
      style: "normal",
    },
  ],
});

const body = Noto_Sans({
  variable: "--font-noto-sans",
  subsets: ["latin-ext"],
  weight: ["400", "500", "600", "700"],
});

export const metadata: Metadata = {
  title: "LaserTag Race Control",
  description: "Real-time race telemetry dashboard for electric laser tag karts.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="cs"
      className={`${headline.variable} ${body.variable} h-full antialiased`}
    >
      <body className="min-h-full flex flex-col">{children}</body>
    </html>
  );
}
