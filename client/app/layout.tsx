import { Inter } from "next/font/google";
import "./globals.css";
import ToasterClient from "./ToasterClient";
import Script from "next/script";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
});

export const metadata = {
  title: "Book Store Nepal",
  description: "Online Book Store",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="bg-background px-6 py-4">
      <body className={inter.variable} suppressHydrationWarning>
        <ToasterClient />
        <Script
          src="https://accounts.google.com/gsi/client"
          strategy="lazyOnload"
        />

        {children}
      </body>
    </html>
  );
}
