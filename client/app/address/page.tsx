"use client";

import { apiRequest } from "@/lib/api";
import { useState } from "react";
import { toast } from "sonner";

export default function AddressPage() {
  const [city, setCity] = useState("");
  const [deliveryAddress, setDeliveryAddress] = useState("");
  const [loading, setLoading] = useState(false);
  async function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    try {
      await apiRequest("address/add-address", {
        method: "POST",
        body: JSON.stringify({ city, delivery_address: deliveryAddress }),
      });
      toast.success("Address added");
      setCity("");
      setDeliveryAddress("");
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Could not add address",
      );
    } finally {
      setLoading(false);
    }
  }
  return (
    <main className="mx-auto max-w-2xl px-4 py-10">
      <h1 className="text-3xl font-bold text-primary">Delivery addresses</h1>
      <form
        onSubmit={submit}
        className="mt-8 space-y-4 rounded-2xl border border-border bg-surface p-6"
      >
        <input
          required
          placeholder="City"
          value={city}
          onChange={(e) => setCity(e.target.value)}
          className="h-12 w-full rounded-lg border border-border px-4"
        />
        <textarea
          required
          placeholder="Delivery address"
          value={deliveryAddress}
          onChange={(e) => setDeliveryAddress(e.target.value)}
          className="min-h-32 w-full rounded-lg border border-border p-4"
        />
        <button
          disabled={loading}
          className="h-12 w-full rounded-lg bg-primary font-semibold text-white"
        >
          {loading ? "Adding..." : "Add address"}
        </button>
      </form>
      <p className="mt-4 text-sm text-muted">
        Use the address ID returned by your account data when placing an order.
      </p>
    </main>
  );
}
