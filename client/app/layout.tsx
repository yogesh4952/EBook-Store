import { Inter } from "next/font/google";
import "./globals.css";
import ToasterClient from "./ToasterClient";
import Script from "next/script";
import type { Metadata, Viewport } from "next";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
});

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
};

export const metadata: Metadata = {
  title: {
    default: "Book Store Nepal - Online Bookstore in Nepal",
    template: "%s | Book Store Nepal",
  },
  description:
    "Book Store Nepal is your trusted online bookstore in Nepal. Browse thousands of books from trusted sellers, enjoy fast delivery across Nepal with cash on delivery and eSewa payment. Find bestsellers, Nepali books, and more.",
  keywords: [
    "online bookstore nepal",
    "buy books nepal",
    "book store nepal",
    "nepali books",
    "english books nepal",
    "bookshop nepal",
    "kathmandu book store",
    "books delivery nepal",
    "eSewa books",
    "cash on delivery books nepal",
    "nepali literature",
    "bestsellers nepal",
    "book marketplace nepal",
    "buy books online nepal",
    "affordable books nepal",
  ],
  authors: [{ name: "Book Store Nepal", url: "https://book.yogeshthakuri.com" }],
  creator: "Book Store Nepal",
  publisher: "Book Store Nepal",
  formatDetection: {
    email: false,
    telephone: false,
    address: false,
  },
  metadataBase: new URL("https://book.yogeshthakuri.com"),
  alternates: {
    canonical: "/",
    languages: {
      "en-NP": "/",
      "ne-NP": "/",
      en: "/",
      ne: "/",
    },
  },
  openGraph: {
    title: "Book Store Nepal - Online Bookstore in Nepal",
    description:
      "Book Store Nepal is your trusted online bookstore in Nepal. Browse thousands of books from trusted sellers, enjoy fast delivery across Nepal with cash on delivery and eSewa payment.",
    url: "https://book.yogeshthakuri.com",
    siteName: "Book Store Nepal",
    images: [
      {
        url: "https://book.yogeshthakuri.com/og-image.png",
        width: 1200,
        height: 630,
        alt: "Book Store Nepal - Online Bookstore in Nepal",
        type: "image/png",
      },
    ],
    locale: "en_NP",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "Book Store Nepal - Online Bookstore in Nepal",
    description:
      "Book Store Nepal is your trusted online bookstore in Nepal. Browse thousands of books from trusted sellers, enjoy fast delivery across Nepal with cash on delivery and eSewa payment.",
    images: ["https://book.yogeshthakuri.com/og-image.png"],
    creator: "@bookstorenepal",
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      noimageindex: false,
      "max-video-preview": -1,
      "max-image-preview": "large",
    },
  },
  verification: {
    google: "your-google-verification-code-here",
    other: {
      "msvalidate.01": "your-bing-verification-code-here",
    },
  },
};

const jsonLd = {
  "@context": "https://schema.org",
  "@type": "Organization",
  name: "Book Store Nepal",
  description:
    "Book Store Nepal is your trusted online bookstore in Nepal. Browse thousands of books from trusted sellers, enjoy fast delivery across Nepal with cash on delivery and eSewa payment.",
  url: "https://book.yogeshthakuri.com",
  logo: "https://book.yogeshthakuri.com/logo.png",
  contactPoint: {
    "@type": "ContactPoint",
    telephone: "+977-1-XXXXXXX",
    contactType: "customer service",
    areaServed: "NP",
    availableLanguage: ["English", "Nepali"],
  },
  address: {
    "@type": "PostalAddress",
    streetAddress: "New Road, Kathmandu",
    addressLocality: "Kathmandu",
    postalCode: "44600",
    addressCountry: "NP",
  },
  sameAs: [
    "https://www.facebook.com/bookstorenepal",
    "https://www.instagram.com/bookstorenepal",
    "https://x.com/bookstorenepal",
    "https://www.youtube.com/@bookstorenepal",
  ],
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      className="bg-background px-6 py-4"
      suppressHydrationWarning
    >
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="theme-color" content="#4B2E20" />
        <meta name="color-scheme" content="light" />
        <link rel="icon" type="image/svg+xml" href="/favicon.svg" />
        <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png" />
        <link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png" />
        <link rel="manifest" href="/site.webmanifest" />
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link rel="dns-prefetch" href="https://covers.openlibrary.org" />
        <link rel="dns-prefetch" href="https://accounts.google.com" />
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
        />
      </head>
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