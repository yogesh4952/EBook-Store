import Link from "next/link";
import { HiOutlineChevronLeft } from "react-icons/hi2";

import Orders from "@/components/orders/Orders";

export default function OrdersPage() {
  return (
    <main className="mx-auto w-full max-w-5xl px-6 py-10">
      <div className="flex items-center gap-4">
        <Link
          href="/"
          className="flex items-center gap-1 text-sm text-muted hover:text-primary"
        >
          <HiOutlineChevronLeft size={16} /> Home
        </Link>
        <span className="text-muted">/</span>
        <h1 className="text-2xl font-bold text-primary">My orders</h1>
      </div>
      <p className="mt-1 text-sm text-muted">
        Track payments and delivery progress for every order you have placed.
      </p>

      <div className="mt-8">
        <Orders />
      </div>
    </main>
  );
}
