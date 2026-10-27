import Footer from "@/components/common/Footer";
import Navbar from "@/components/common/Navbar";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Delivery Addresses - Book Store Nepal",
  description:
    "Manage your delivery addresses for fast book delivery across Nepal. Save and update your shipping addresses for Book Store Nepal orders.",
  robots: {
    index: false,
    follow: false,
  },
};

export default function AddressesLayout({
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