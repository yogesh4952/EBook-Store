import Footer from "@/components/common/Footer";
import Navbar from "@/components/common/Navbar";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Browse Books - Book Store Nepal | Online Bookstore",
  description:
    "Browse thousands of books at Book Store Nepal. Filter by category, genre, price, and availability. Find the best books in Nepal with fast delivery and cash on delivery or eSewa payment.",
  keywords: [
    "browse books nepal",
    "book catalog nepal",
    "all books nepal",
    "book collection nepal",
    "buy books online nepal",
  ],
  openGraph: {
    title: "Browse Books - Book Store Nepal",
    description:
      "Browse thousands of books at Book Store Nepal. Filter by category, genre, price, and availability.",
    url: "https://book.yogeshthakuri.com/book",
    siteName: "Book Store Nepal",
    images: [
      {
        url: "https://book.yogeshthakuri.com/Banner.png",
        width: 1200,
        height: 630,
        alt: "Browse Books - Book Store Nepal",
      },
    ],
    locale: "en_NP",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "Browse Books - Book Store Nepal",
    description:
      "Browse thousands of books at Book Store Nepal. Filter by category, genre, price, and availability.",
    images: ["https://book.yogeshthakuri.com/Banner.png"],
  },
};

export default function BookLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-full flex flex-col">
      <Navbar />
      <main className="flex-1">{children}</main>
      <Footer />
    </div>
  );
}