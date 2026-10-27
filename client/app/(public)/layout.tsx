import Footer from "@/components/common/Footer";
import Navbar from "@/components/common/Navbar";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Home - Book Store Nepal | Online Bookstore in Nepal",
  description:
    "Welcome to Book Store Nepal - your trusted online bookstore. Browse thousands of books from trusted sellers, enjoy fast delivery across Nepal with cash on delivery and eSewa payment. Find bestsellers, Nepali books, and more.",
  keywords: [
    "online bookstore nepal",
    "buy books nepal",
    "book store nepal",
    "nepali books",
    "english books nepal",
    "bestsellers nepal",
    "bookshop nepal",
    "kathmandu book store",
    "books delivery nepal",
  ],
  openGraph: {
    title: "Book Store Nepal - Online Bookstore in Nepal",
    description:
      "Welcome to Book Store Nepal - your trusted online bookstore. Browse thousands of books from trusted sellers, enjoy fast delivery across Nepal with cash on delivery and eSewa payment.",
    url: "https://book.yogeshthakuri.com/",
    siteName: "Book Store Nepal",
    images: [
      {
        url: "https://book.yogeshthakuri.com/Banner.png",
        width: 1200,
        height: 630,
        alt: "Book Store Nepal - Online Bookstore in Nepal",
      },
    ],
    locale: "en_NP",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "Book Store Nepal - Online Bookstore in Nepal",
    description:
      "Welcome to Book Store Nepal - your trusted online bookstore. Browse thousands of books from trusted sellers, enjoy fast delivery across Nepal.",
    images: ["https://book.yogeshthakuri.com/Banner.png"],
  },
};

export default function PublicLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="min-h-full flex flex-col">
      <Navbar />
      <main className="flex-1">{children}</main>
      <Footer />
    </div>
  );
}