"use client";

import { placeOrder, type Order } from "@/lib/api";
import { readCart, type CartItem } from "@/lib/cart";
import { useState } from "react";
import { toast } from "sonner";

export default function CheckoutPage() {
  const [addressId, setAddressId] = useState("");
  const [paymentMethod, setPaymentMethod] = useState<"COD" | "ESEWA">("COD");
  const [items] = useState<CartItem[]>(() => readCart());
  const [sellerId, setSellerId] = useState("");
  const [order, setOrder] = useState<Order | null>(null);
  const [loading, setLoading] = useState(false);

  async function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!addressId || !items.length || !sellerId) {
      toast.error("Add a book, address ID, and seller ID before ordering");
      return;
    }
    setLoading(true);
    try {
      const response = await placeOrder({
        payment_method: paymentMethod,
        user_address_id: Number(addressId),
        items: items.map((item) => ({
          book_id: item.book.id,
          seller_id: item.book.seller_id || Number(sellerId),
          quantity: item.quantity,
        })),
      });
      setOrder(response.data);
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Could not place order",
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="mx-auto max-w-2xl px-4 py-10">
      <p className="text-sm font-semibold uppercase tracking-[0.18em] text-primary-light">
        Purchase
      </p>
      <h1 className="mt-2 text-4xl font-bold text-primary">Checkout</h1>
      {order ? (
        <section className="mt-8 rounded-2xl border border-green-200 bg-green-50 p-6">
          <h2 className="text-2xl font-bold text-green-800">Order placed</h2>
          <p className="mt-2 text-green-800">
            Your order code is {order.order_code}.
          </p>
          <p className="mt-2">Total: Rs. {order.total_price}</p>
        </section>
      ) : (
        <form
          onSubmit={submit}
          className="mt-8 space-y-5 rounded-2xl border border-border bg-surface p-6"
        >
          <label className="block text-sm font-semibold text-primary">
            Address ID
            <input
              required
              type="number"
              value={addressId}
              onChange={(e) => setAddressId(e.target.value)}
              className="mt-2 h-12 w-full rounded-lg border border-border px-4"
            />
          </label>
          <label className="block text-sm font-semibold text-primary">
            Seller ID
            <input
              required
              type="number"
              value={sellerId}
              onChange={(e) => setSellerId(e.target.value)}
              className="mt-2 h-12 w-full rounded-lg border border-border px-4"
            />
          </label>
          <p className="rounded-lg bg-background p-4 text-sm text-muted">
            {items.length
              ? `${items.length} book${items.length === 1 ? "" : "s"} ready from your cart.`
              : "Your cart is empty. Add books before checkout."}
          </p>
          <select
            value={paymentMethod}
            onChange={(e) =>
              setPaymentMethod(e.target.value as "COD" | "ESEWA")
            }
            className="h-12 w-full rounded-lg border border-border bg-white px-4"
          >
            <option value="COD">Cash on delivery</option>
            <option value="ESEWA">eSewa</option>
          </select>
          <button
            disabled={loading}
            className="h-12 w-full rounded-lg bg-primary font-semibold text-white"
          >
            {loading ? "Placing order..." : "Place order"}
          </button>
        </form>
      )}
    </main>
  );
}
