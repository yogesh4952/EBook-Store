import Footer from "@/components/common/Footer";
import Navbar from "@/components/common/Navbar";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Checkout - Book Store Nepal",
  description:
    "Complete your book order at Book Store Nepal. Choose between eSewa online payment or cash on delivery. Fast delivery across Nepal.",
  robots: {
    index: false,
    follow: false,
  },
};

export default function CheckoutLayout({
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