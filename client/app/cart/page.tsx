"use client";

import {
  readCart,
  removeFromCart,
  updateCartQuantity,
  type CartItem,
} from "@/lib/cart";
import Link from "next/link";
import { useEffect, useState } from "react";

export default function CartPage() {
  const [items, setItems] = useState<CartItem[]>([]);

  useEffect(() => {
    const sync = () => setItems(readCart());
    sync();
    window.addEventListener("cart-updated", sync);
    return () => window.removeEventListener("cart-updated", sync);
  }, []);

  const total = items.reduce(
    (sum, item) => sum + item.book.price * item.quantity,
    0,
  );

  return (
    <main className="mx-auto max-w-4xl px-4 py-10">
      <p className="text-sm font-semibold uppercase tracking-[0.18em] text-primary-light">
        Your selection
      </p>
      <h1 className="mt-2 text-4xl font-bold text-primary">Shopping cart</h1>
      {items.length === 0 ? (
        <section className="mt-8 rounded-2xl border border-dashed border-border p-10 text-center">
          <p className="text-muted">Your cart is empty.</p>
          <Link
            href="/books"
            className="mt-5 inline-block rounded-lg bg-primary px-5 py-3 font-semibold text-white"
          >
            Browse books
          </Link>
        </section>
      ) : (
        <div className="mt-8 grid gap-6 lg:grid-cols-[1fr_280px]">
          <section className="space-y-3">
            {items.map((item) => (
              <article
                key={item.book.id}
                className="flex items-center gap-4 rounded-xl border border-border bg-surface p-4"
              >
                <div
                  className="h-20 w-14 shrink-0 rounded bg-accent/30 bg-cover bg-center"
                  style={{
                    backgroundImage: item.book.cover_page_url
                      ? `url(${item.book.cover_page_url})`
                      : undefined,
                  }}
                />
                <div className="min-w-0 flex-1">
                  <h2 className="truncate font-semibold text-primary">
                    {item.book.title}
                  </h2>
                  <p className="text-sm text-muted">Rs. {item.book.price}</p>
                </div>
                <input
                  aria-label={`Quantity for ${item.book.title}`}
                  type="number"
                  min="1"
                  value={item.quantity}
                  onChange={(event) =>
                    updateCartQuantity(item.book.id, Number(event.target.value))
                  }
                  className="h-10 w-16 rounded border border-border px-2"
                />
                <button
                  onClick={() => removeFromCart(item.book.id)}
                  className="text-sm font-semibold text-error"
                >
                  Remove
                </button>
              </article>
            ))}
          </section>
          <aside className="h-fit rounded-xl border border-border bg-surface p-5">
            <h2 className="font-semibold text-primary">Summary</h2>
            <div className="mt-4 flex justify-between border-t border-border pt-4 font-bold">
              <span>Total</span>
              <span>Rs. {total}</span>
            </div>
            <Link
              href="/checkout"
              className="mt-5 block rounded-lg bg-primary px-4 py-3 text-center font-semibold text-white"
            >
              Continue to checkout
            </Link>
          </aside>
        </div>
      )}
    </main>
  );
}
