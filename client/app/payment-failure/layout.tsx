import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Payment Failed - Book Store Nepal",
  description:
    "Your payment could not be completed. Try again or use an alternative payment method.",
  robots: {
    index: false,
    follow: false,
  },
};

export default function PaymentFailureLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return <>{children}</>;
}