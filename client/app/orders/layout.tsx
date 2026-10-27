import Footer from "@/components/common/Footer";
import Navbar from "@/components/common/Navbar";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "My Orders - Book Store Nepal",
  description:
    "Track your book orders and delivery progress at Book Store Nepal. View order history, status, and manage your purchases.",
  robots: {
    index: false,
    follow: false,
  },
};

export default function OrdersLayout({
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