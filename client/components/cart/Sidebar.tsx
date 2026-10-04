"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { IoClose } from "react-icons/io5";
import { toast } from "sonner";

import CartItem, {
  type ICartItem,
  type ICartResponse,
} from "./CartItem";

const formatPrice = (value: number) =>
  `Rs. ${new Intl.NumberFormat("en-NP", {
    maximumFractionDigits: 0,
  }).format(value)}`;

// TODO: replace once the cart module owns the state
const Sidebar = ({ onClose }: { onClose?: () => void }) => {
  const [cartItems, setCartItems] = useState<ICartItem[]>([]);
  const [total, setTotal] = useState(0);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") onClose?.();
    };
    document.addEventListener("keydown", onKeyDown);
    return () => document.removeEventListener("keydown", onKeyDown);
  }, [onClose]);

  useEffect(() => {
    const controller = new AbortController();

    const fetchCart = async () => {
      try {
        // Trailing slash matters: /api/cart 301-redirects, and the redirect
        // drops the Authorization header before it reaches Gin.
        const resp = await fetch("/api/proxy/cart/", {
          method: "GET",
          credentials: "include",
          cache: "no-store",
          signal: controller.signal,
        });

        const payload = await resp.json();

        if (!resp.ok) {
          toast.error(payload.message ?? "Error fetching cart items");
          return;
        }

        const data = payload.data as ICartResponse;
        setCartItems(data?.items ?? []);
        setTotal(data?.total ?? 0);
      } catch (error) {
        if (controller.signal.aborted) return;
        toast.error("Internal Server Error!");
        console.error(error);
      } finally {
        if (!controller.signal.aborted) setIsLoading(false);
      }
    };

    fetchCart();

    return () => controller.abort();
  }, []);

  const itemCount = cartItems.reduce((sum, item) => sum + item.quantity, 0);

  return (
    <div className="fixed inset-0 z-50">
      <div
        onClick={onClose}
        className="absolute inset-0 bg-text/20 backdrop-blur-sm"
      />

      <aside className="absolute inset-y-0 right-0 flex w-full max-w-sm flex-col border-l border-border bg-surface shadow-2xl">
        <div className="flex items-center justify-between border-b border-border px-5 py-4">
          <h2 className="text-lg font-semibold text-primary">Your cart</h2>
          <button
            onClick={onClose}
            aria-label="Close cart"
            className="rounded-lg p-2 text-muted transition-colors hover:bg-accent/20 hover:text-primary"
          >
            <IoClose size={20} />
          </button>
        </div>

        {isLoading ? (
          <div className="flex flex-1 flex-col items-center justify-center gap-3 px-5 py-4">
            <span className="size-6 animate-spin rounded-full border-2 border-border border-t-primary" />
            <p className="text-sm text-muted">Loading your cart…</p>
          </div>
        ) : cartItems.length === 0 ? (
          <div className="flex flex-1 flex-col items-center justify-center gap-2 px-5 py-4 text-center">
            <p className="text-sm font-medium text-primary">Your cart is empty</p>
            <p className="text-sm text-muted">
              Browse the store and add a book to get started.
            </p>
            <Link
              href="/"
              onClick={onClose}
              className="mt-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-light"
            >
              Continue shopping
            </Link>
          </div>
        ) : (
          <>
            <ul className="flex-1 overflow-y-auto px-5">
              {cartItems.map((item) => (
                <CartItem key={item.data.id} item={item} />
              ))}
            </ul>

            <div className="border-t border-border px-5 py-4">
              <div className="flex items-center justify-between text-sm">
                <span className="text-muted">
                  Subtotal ({itemCount} {itemCount === 1 ? "item" : "items"})
                </span>
                <span className="font-semibold text-primary tabular-nums">
                  {formatPrice(total)}
                </span>
              </div>
              <p className="mt-1 text-xs text-muted">
                Shipping and taxes calculated at checkout.
              </p>
              <Link
                href="/checkout"
                onClick={onClose}
                className="mt-3 block w-full rounded-lg bg-primary px-4 py-3 text-center text-sm font-medium text-white transition-colors hover:bg-primary-light"
              >
                Proceed to checkout
              </Link>
            </div>
          </>
        )}
      </aside>
    </div>
  );
};

export default Sidebar;