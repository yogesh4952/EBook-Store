import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Order Confirmed - Book Store Nepal",
  description:
    "Your book order has been successfully placed at Book Store Nepal. Thank you for your purchase.",
  robots: {
    index: false,
    follow: false,
  },
};

export default function PaymentSuccessLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return <>{children}</>;
}