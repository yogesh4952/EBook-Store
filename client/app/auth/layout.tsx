import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Login - Book Store Nepal",
  description:
    "Login to your Book Store Nepal account. Access your orders, saved addresses, and continue your reading journey with thousands of books.",
  robots: {
    index: false,
    follow: false,
  },
};

export default function AuthLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return <>{children}</>;
}