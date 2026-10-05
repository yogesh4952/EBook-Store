import Link from "next/link";
import { HiOutlineChevronLeft } from "react-icons/hi2";

import Addresses from "@/components/addresses/Addresses";

export default function AddressesPage() {
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
        <h1 className="text-2xl font-bold text-primary">Delivery addresses</h1>
      </div>
      <p className="mt-1 text-sm text-muted">
        Manage the addresses you can choose from at checkout.
      </p>

      <div className="mt-8">
        <Addresses />
      </div>
    </main>
  );
}
