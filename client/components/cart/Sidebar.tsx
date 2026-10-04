"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { IoClose } from "react-icons/io5";
import { toast } from "sonner";

import CartItem, { type ICartItem } from "./CartItem";
import { addToCart, getCart, removeFromCart } from "@/lib/checkout";

const formatPrice = (value: number) =>
  `Rs. ${new Intl.NumberFormat("en-NP", {
    maximumFractionDigits: 0,
  }).format(value)}`;

/** Applies a new quantity to one line, dropping the row when it hits zero. */
const reduceItem = (
  items: ICartItem[],
  bookId: number,
  quantity: number,
): ICartItem[] =>
  items
    .map((item) =>
      item.data.id === bookId
        ? {
            ...item,
            quantity,
            subtotal: item.data.price * quantity,
          }
        : item,
    )
    .filter((item) => item.quantity > 0);

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
        const cart = await getCart();
        if (controller.signal.aborted) return;
        setCartItems(cart?.items ?? []);
        setTotal(cart?.total ?? 0);
      } catch (error) {
        if (controller.signal.aborted) return;
        toast.error(
          error instanceof Error ? error.message : "Error fetching cart items",
        );
      } finally {
        if (!controller.signal.aborted) setIsLoading(false);
      }
    };

    fetchCart();

    // Re-sync when the user comes back to the tab, so an order placed
    // elsewhere does not leave stale items in the drawer.
    window.addEventListener("focus", fetchCart);

    return () => {
      controller.abort();
      window.removeEventListener("focus", fetchCart);
    };
  }, []);

  /**
   * Applies a new quantity locally, then re-reads the cart so the displayed
   * totals come from the server rather than from guesswork.
   *
   * Going up and going down hit different endpoints: the remove endpoint only
   * decrements, so an increase has to go through add-to-cart. Decrementing sends
   * the *amount to remove*, not the new quantity.
   *
   * The remove endpoint refuses to take more than the cart holds, which keeps
   * the stored count from going negative.
   */
  const handleQuantityChange = async (bookId: number, quantity: number) => {
    const current = cartItems.find((item) => item.data.id === bookId);
    if (!current || quantity === current.quantity) return;

    const previous = cartItems;
    setCartItems(reduceItem(cartItems, bookId, quantity));
    setTotal(
      (prev) =>
        prev - current.data.price * (current.quantity - quantity),
    );

    try {
      if (quantity > current.quantity) {
        await addToCart({ book_id: bookId, quantity: quantity - current.quantity });
      } else {
        await removeFromCart({
          book_id: bookId,
          quantity: current.quantity - quantity,
        });
      }

      const cart = await getCart();
      setCartItems(cart.items ?? []);
      setTotal(cart.total ?? 0);
    } catch (error) {
      setCartItems(previous);
      toast.error(
        error instanceof Error ? error.message : "Could not update the cart",
      );
    } finally {
      setPendingBookId(null);
    }
  };

  const handleRemove = async (bookId: number) => {
    const current = cartItems.find((item) => item.data.id === bookId);
    if (!current) return;

    const previous = cartItems;
    setCartItems(cartItems.filter((item) => item.data.id !== bookId));
    setTotal((prev) => prev - current.subtotal);

    try {
      await removeFromCart({ book_id: bookId, quantity: current.quantity });
      const cart = await getCart();
      setCartItems(cart.items ?? []);
      setTotal(cart.total ?? 0);
    } catch (error) {
      setCartItems(previous);
      toast.error(
        error instanceof Error ? error.message : "Could not remove the item",
      );
    } finally {
      setPendingBookId(null);
    }
  };

  const [pendingBookId, setPendingBookId] = useState<number | null>(null);

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
            <p className="text-sm font-medium text-primary">
              Your cart is empty
            </p>
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
                <CartItem
                  key={item.data.id}
                  item={item}
                  isUpdating={pendingBookId === item.data.id}
                  onQuantityChange={(bookId, quantity) => {
                    setPendingBookId(bookId);
                    handleQuantityChange(bookId, quantity);
                  }}
                  onRemove={(bookId) => {
                    setPendingBookId(bookId);
                    handleRemove(bookId);
                  }}
                />
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
