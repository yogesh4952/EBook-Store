"use client";

import { apiRequest, type Order } from "@/lib/api";
import { useEffect, useState } from "react";

export default function OrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [error, setError] = useState("");
  useEffect(() => {
    apiRequest<{ data: Order[] }>("order/list-user-order")
      .then((result) => setOrders(result.data || []))
      .catch((reason) => setError(reason.message));
  }, []);
  return (
    <main className="mx-auto max-w-4xl px-4 py-10">
      <p className="text-sm font-semibold uppercase tracking-[0.18em] text-primary-light">
        Account
      </p>
      <h1 className="mt-2 text-4xl font-bold text-primary">Your orders</h1>
      {error && (
        <p className="mt-6 rounded-lg bg-red-50 p-4 text-red-700">{error}</p>
      )}
      <div className="mt-8 space-y-4">
        {orders.map((order) => (
          <article
            key={order.order_code}
            className="rounded-xl border border-border bg-surface p-5"
          >
            <div className="flex flex-wrap justify-between gap-3">
              <h2 className="font-bold text-primary">{order.order_code}</h2>
              <span className="rounded-full bg-accent/30 px-3 py-1 text-sm">
                {order.order_status}
              </span>
            </div>
            <p className="mt-3 text-sm text-muted">{order.address}</p>
            <p className="mt-4 font-semibold">
              Rs. {order.total_price}{" "}
              <span className="ml-3 text-sm font-normal text-muted">
                Payment: {order.payment_status}
              </span>
            </p>
          </article>
        ))}
        {!error && orders.length === 0 && (
          <p className="rounded-lg border border-dashed border-border p-8 text-center text-muted">
            No orders yet.
          </p>
        )}
      </div>
    </main>
  );
}
