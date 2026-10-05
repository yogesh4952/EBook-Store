"use client";

import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import { HiOutlineShoppingBag } from "react-icons/hi2";
import { toast } from "sonner";

import {
  OrderAddress,
  OrderItems,
  OrderSummaryBar,
  OrderTimeline,
  PaymentHint,
} from "@/components/orders/OrderUI";
import { listUserOrders, type UserOrder } from "@/lib/checkout";

type Filter = "all" | "active" | "delivered" | "cancelled";

const FILTERS: { value: Filter; label: string }[] = [
  { value: "all", label: "All orders" },
  { value: "active", label: "In progress" },
  { value: "delivered", label: "Delivered" },
  { value: "cancelled", label: "Cancelled" },
];

const matches = (order: UserOrder, filter: Filter) => {
  switch (filter) {
    case "active":
      return order.order_status === "PLACED";
    case "delivered":
      return order.order_status === "DELIVERED";
    case "cancelled":
      return order.order_status === "CANCELLED";
    default:
      return true;
  }
};

const Orders = () => {
  const [orders, setOrders] = useState<UserOrder[]>([]);
  const [filter, setFilter] = useState<Filter>("all");
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const controller = new AbortController();

    const fetchOrders = async () => {
      try {
        const data = await listUserOrders();
        if (controller.signal.aborted) return;
        setOrders(data);
      } catch (error) {
        if (controller.signal.aborted) return;
        toast.error(
          error instanceof Error ? error.message : "Could not load your orders",
        );
      } finally {
        if (!controller.signal.aborted) setIsLoading(false);
      }
    };

    fetchOrders();

    // An order placed in another tab shows up on the next focus, so the list
    // does not go stale while the customer browses.
    window.addEventListener("focus", fetchOrders);

    return () => {
      controller.abort();
      window.removeEventListener("focus", fetchOrders);
    };
  }, []);

  // Newest first. The API returns rows in insertion order, which is only
  // incidentally chronological, so sort explicitly.
  const sorted = useMemo(
    () =>
      [...orders].sort(
        (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
      ),
    [orders],
  );

  const visible = useMemo(
    () => sorted.filter((order) => matches(order, filter)),
    [sorted, filter],
  );

  const counts = useMemo(
    () => ({
      all: sorted.length,
      active: sorted.filter((order) => order.order_status === "PLACED").length,
      delivered: sorted.filter((order) => order.order_status === "DELIVERED")
        .length,
      cancelled: sorted.filter((order) => order.order_status === "CANCELLED")
        .length,
    }),
    [sorted],
  );

  const spend = useMemo(
    () =>
      sorted
        .filter((order) => order.payment_status === "PAID")
        .reduce((sum, order) => sum + order.total_price, 0),
    [sorted],
  );

  if (isLoading) {
    return (
      <div className="flex min-h-[60vh] flex-col items-center justify-center gap-3">
        <span className="size-7 animate-spin rounded-full border-2 border-border border-t-primary" />
        <p className="text-sm text-muted">Loading your orders…</p>
      </div>
    );
  }

  if (orders.length === 0) {
    return (
      <div className="rounded-2xl border border-border bg-surface px-6 py-16 text-center shadow-sm">
        <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-accent/20">
          <HiOutlineShoppingBag size={28} className="text-primary" />
        </div>
        <h1 className="mt-6 text-2xl font-bold text-primary">No orders yet</h1>
        <p className="mx-auto mt-2 max-w-sm text-sm text-muted">
          When you place an order it will show up here with its status, payment
          method, and delivery progress.
        </p>
        <Link
          href="/book"
          className="mt-6 inline-block rounded-lg bg-primary px-6 py-3 text-sm font-medium text-white transition-colors hover:bg-primary-light"
        >
          Browse books
        </Link>
      </div>
    );
  }

  return (
    <>
      {/* Summary tiles */}
      <div className="grid gap-4 sm:grid-cols-3">
        <div className="rounded-2xl border border-border bg-surface p-5 shadow-sm">
          <p className="text-xs font-medium tracking-wide text-muted uppercase">
            Orders placed
          </p>
          <p className="mt-2 text-2xl font-bold text-primary tabular-nums">
            {counts.all}
          </p>
        </div>
        <div className="rounded-2xl border border-border bg-surface p-5 shadow-sm">
          <p className="text-xs font-medium tracking-wide text-muted uppercase">
            In progress
          </p>
          <p className="mt-2 text-2xl font-bold text-primary tabular-nums">
            {counts.active}
          </p>
        </div>
        <div className="rounded-2xl border border-border bg-surface p-5 shadow-sm">
          <p className="text-xs font-medium tracking-wide text-muted uppercase">
            Total paid
          </p>
          <p className="mt-2 text-2xl font-bold text-primary tabular-nums">
            {new Intl.NumberFormat("en-NP", {
              maximumFractionDigits: 0,
              style: "currency",
              currency: "NPR",
              currencyDisplay: "narrowSymbol",
            }).format(spend)}
          </p>
        </div>
      </div>

      {/* Filters */}
      <div className="mt-8 flex flex-wrap items-center gap-2">
        {FILTERS.map(({ value, label }) => {
          const isActive = filter === value;

          return (
            <button
              key={value}
              type="button"
              onClick={() => setFilter(value)}
              aria-pressed={isActive}
              className={`rounded-full border px-4 py-1.5 text-sm font-medium transition-colors ${
                isActive
                  ? "border-primary bg-primary text-white"
                  : "border-border bg-surface text-muted hover:border-accent hover:bg-accent/10 hover:text-primary"
              }`}
            >
              {label}
              <span
                className={`ml-1.5 tabular-nums ${
                  isActive ? "text-white/70" : "text-muted/70"
                }`}
              >
                {counts[value]}
              </span>
            </button>
          );
        })}
      </div>

      {/* Orders */}
      {visible.length === 0 ? (
        <div className="mt-6 rounded-xl border border-border bg-surface px-6 py-12 text-center">
          <p className="text-sm font-medium text-primary">
            Nothing in this view
          </p>
          <p className="mt-1 text-sm text-muted">
            Try a different filter to see your other orders.
          </p>
        </div>
      ) : (
        <ul className="mt-6 space-y-5">
          {visible.map((order) => (
            <li
              key={order.id}
              className="overflow-hidden rounded-2xl border border-border bg-surface shadow-sm transition-shadow hover:shadow-md"
            >
              <OrderSummaryBar order={order} />

              <div className="px-5 py-4">
                <OrderTimeline order={order} />
              </div>

              <OrderItems order={order} />

              <div className="flex flex-col gap-2 border-t border-border px-5 py-4 sm:flex-row sm:items-center sm:justify-between">
                <OrderAddress order={order} />
                <PaymentHint method={order.payment_method} />
              </div>
            </li>
          ))}
        </ul>
      )}
    </>
  );
};

export default Orders;
