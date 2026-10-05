"use client";

import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { FaHandHoldingDollar } from "react-icons/fa6";
import { toast } from "sonner";

import EsewaRedirect from "@/components/checkout/EsewaRedirect";
import type { ICartItem } from "@/components/cart/CartItem";
import {
  addAddress,
  cartToOrderItems,
  clearCart,
  getCart,
  listUserAddresses,
  placeOrder,
  type EsewaPayload,
  type PaymentMethod,
  type UserAddress,
} from "@/lib/checkout";

const formatPrice = (value: number) =>
  `Rs. ${new Intl.NumberFormat("en-NP", { maximumFractionDigits: 0 }).format(value)}`;

const PAYMENT_METHODS: {
  value: PaymentMethod;
  label: string;
  description: string;
  logo?: string;
  Icon?: typeof FaHandHoldingDollar;
}[] = [
  {
    value: "ESEWA",
    label: "Pay with eSewa",
    description: "Redirects to eSewa to complete payment.",
    logo: "/esewa_logo.jpg",
  },
  {
    value: "COD",
    label: "Cash on delivery",
    description: "Pay in cash when your order arrives.",
    Icon: FaHandHoldingDollar,
  },
];

export default function CheckoutPage() {
  const router = useRouter();

  const [items, setItems] = useState<ICartItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isPlacing, setIsPlacing] = useState(false);
  const [isSavingAddress, setIsSavingAddress] = useState(false);

  const [city, setCity] = useState("");
  const [deliveryAddress, setDeliveryAddress] = useState("");
  const [addressId, setAddressId] = useState<number | null>(null);
  const [savedAddresses, setSavedAddresses] = useState<UserAddress[]>([]);

  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>("ESEWA");

  const [esewaPayload, setEsewaPayload] = useState<{
    payload: EsewaPayload;
    code: string;
  } | null>(null);

  useEffect(() => {
    const controller = new AbortController();

    const loadCart = async () => {
      try {
        const cart = await getCart();
        if (controller.signal.aborted) return;
        setItems(cart.items ?? []);
      } catch (error) {
        if (controller.signal.aborted) return;
        toast.error(
          error instanceof Error ? error.message : "Could not load your cart",
        );
      } finally {
        if (!controller.signal.aborted) setIsLoading(false);
      }
    };

    loadCart();
    return () => controller.abort();
  }, []);

  useEffect(() => {
    const controller = new AbortController();

    const loadAddresses = async () => {
      try {
        const rows = await listUserAddresses();
        if (controller.signal.aborted) return;
        setSavedAddresses(rows);
        if (rows.length > 0) {
          setAddressId(rows[0].id);
          setCity(rows[0].city);
          setDeliveryAddress(rows[0].delivery_address);
        }
      } catch (error) {
        if (controller.signal.aborted) return;
        console.error(error);
      }
    };

    loadAddresses();
    return () => controller.abort();
  }, []);

  const itemCount = items.reduce((sum, item) => sum + item.quantity, 0);
  const total = items.reduce((sum, item) => sum + item.subtotal, 0);

  const handleSaveAddress = async () => {
    if (!city.trim() || !deliveryAddress.trim()) {
      toast.error("Enter both a city and a delivery address.");
      return;
    }

    setIsSavingAddress(true);
    try {
      const created = await addAddress({
        city: city.trim(),
        delivery_address: deliveryAddress.trim(),
      });
      setAddressId(created.id);
      setSavedAddresses((prev) => [created, ...prev]);
      toast.success("Delivery address saved.");
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Could not save address",
      );
    } finally {
      setIsSavingAddress(false);
    }
  };

  const handlePlaceOrder = async () => {
    if (addressId === null) {
      toast.error("Save your delivery address first.");
      return;
    }

    setIsPlacing(true);
    try {
      const order = await placeOrder({
        payment_method: paymentMethod,
        user_address_id: addressId,
        items: cartToOrderItems(items),
      });

      if (order.esewa_payload) {
        setEsewaPayload({
          payload: order.esewa_payload,
          code: order.order_code,
        });
        return;
      }

      toast.success(`Order ${order.order_code} placed.`);

      // Order is confirmed, so empty the cart. A failure here must not block
      // navigation: the order already exists server-side.
      try {
        await clearCart();
      } catch (error) {
        console.error(error);
        toast.error("Order placed, but we could not clear your cart.");
      }

      await router.push(
        `/orders?placed=${encodeURIComponent(order.order_code)}`,
      );
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Could not place your order",
      );
      setIsPlacing(false);
    }
  };

  if (esewaPayload) {
    return (
      <EsewaRedirect
        payload={esewaPayload.payload}
        orderCode={esewaPayload.code}
      />
    );
  }

  if (isLoading) {
    return (
      <main className="flex min-h-[70vh] items-center justify-center px-6">
        <div className="flex flex-col items-center gap-3">
          <span className="size-7 animate-spin rounded-full border-2 border-border border-t-primary" />
          <p className="text-sm text-muted">Loading your cart…</p>
        </div>
      </main>
    );
  }

  if (items.length === 0) {
    return (
      <main className="flex min-h-[70vh] items-center justify-center px-6">
        <div className="w-full max-w-md rounded-2xl bg-surface p-10 text-center shadow-lg">
          <h1 className="text-2xl font-bold text-primary">
            Nothing to check out
          </h1>
          <p className="mt-2 text-sm text-muted">
            Your cart is empty. Add a book before placing an order.
          </p>
          <Link
            href="/"
            className="mt-6 block rounded-lg bg-primary px-4 py-3 text-sm font-medium text-white transition-colors hover:bg-primary-light"
          >
            Browse books
          </Link>
        </div>
      </main>
    );
  }

  return (
    <main className="mx-auto w-full max-w-6xl px-6 py-10">
      <h1 className="text-3xl font-bold text-primary">Checkout</h1>
      <p className="mt-1 text-sm text-muted">
        {itemCount} {itemCount === 1 ? "item" : "items"} in your cart
      </p>

      <div className="mt-8 grid gap-6 lg:grid-cols-[1fr_360px] lg:items-start">
        <div className="space-y-6">
          {/* Delivery address */}
          <section className="rounded-2xl border border-border bg-surface p-6 shadow-sm">
            <h2 className="text-lg font-semibold text-primary">
              Delivery address
            </h2>

            {savedAddresses.length > 0 && (
              <div className="mt-4 space-y-2">
                <p className="text-sm font-medium text-text">Saved addresses</p>
                {savedAddresses.map((saved) => {
                  const isSelected = addressId === saved.id;
                  return (
                    <label
                      key={saved.id}
                      className={`flex cursor-pointer items-start gap-3 rounded-xl border p-3 transition-colors ${
                        isSelected
                          ? "border-primary bg-accent/10"
                          : "border-border hover:bg-background"
                      }`}
                    >
                      <input
                        type="radio"
                        name="saved_address"
                        checked={isSelected}
                        onChange={() => {
                          setAddressId(saved.id);
                          setCity(saved.city);
                          setDeliveryAddress(saved.delivery_address);
                        }}
                        className="mt-1 accent-[#4B2E20]"
                      />
                      <span className="min-w-0">
                        <span className="block text-sm font-medium text-primary">
                          {saved.city}
                        </span>
                        <span className="mt-0.5 block break-words text-xs text-muted">
                          {saved.delivery_address}
                        </span>
                      </span>
                    </label>
                  );
                })}
              </div>
            )}

            {savedAddresses.length > 0 && (
              <button
                type="button"
                onClick={() => {
                  setAddressId(null);
                  setCity("");
                  setDeliveryAddress("");
                }}
                className="mt-2 text-sm text-muted underline underline-offset-4 transition-colors hover:text-primary"
              >
                Use a different address
              </button>
            )}

            <div className="mt-4 grid gap-4 sm:grid-cols-2">
              <div>
                <label
                  htmlFor="city"
                  className="mb-1.5 block text-sm font-medium text-text"
                >
                  City
                </label>
                <input
                  id="city"
                  name="city"
                  value={city}
                  onChange={(e) => {
                    setCity(e.target.value);
                    setAddressId(null);
                  }}
                  placeholder="Kathmandu"
                  className="h-11 w-full rounded-lg border border-border bg-background px-3 text-sm outline-none transition-colors focus:border-primary"
                />
              </div>

              <div>
                <label
                  htmlFor="delivery_address"
                  className="mb-1.5 block text-sm font-medium text-text"
                >
                  Delivery address
                </label>
                <input
                  id="delivery_address"
                  name="delivery_address"
                  value={deliveryAddress}
                  onChange={(e) => {
                    setDeliveryAddress(e.target.value);
                    setAddressId(null);
                  }}
                  placeholder="Street, landmark"
                  className="h-11 w-full rounded-lg border border-border bg-background px-3 text-sm outline-none transition-colors focus:border-primary"
                />
              </div>
            </div>

            <button
              type="button"
              onClick={handleSaveAddress}
              disabled={isSavingAddress}
              className="mt-4 rounded-lg border border-primary px-4 py-2 text-sm font-medium text-primary transition-colors hover:bg-accent/20 disabled:opacity-50"
            >
              {isSavingAddress ? "Saving…" : "Save address"}
            </button>

            {addressId !== null && (
              <p className="mt-3 text-sm text-success">
                Delivery address selected. Proceed to payment.
              </p>
            )}
          </section>

          {/* Payment method */}
          <section className="rounded-2xl border border-border bg-surface p-6 shadow-sm">
            <h2 className="text-lg font-semibold text-primary">
              Payment method
            </h2>

            <div className="mt-4 space-y-3">
              {PAYMENT_METHODS.map((method) => (
                <label
                  key={method.value}
                  className={`flex cursor-pointer items-start gap-3 rounded-xl border p-4 transition-colors ${
                    paymentMethod === method.value
                      ? "border-primary bg-accent/10"
                      : "border-border hover:bg-background"
                  }`}
                >
                  <input
                    type="radio"
                    name="payment_method"
                    value={method.value}
                    checked={paymentMethod === method.value}
                    onChange={() => setPaymentMethod(method.value)}
                    className="mt-1 accent-[#4B2E20]"
                  />
                  {method.logo ? (
                    <Image
                      src={method.logo}
                      alt=""
                      width={36}
                      height={36}
                      className="mt-0.5 size-9 shrink-0 rounded-lg border border-border object-cover"
                    />
                  ) : (
                    method.Icon && (
                      <method.Icon
                        size={28}
                        className="mt-1 shrink-0 text-primary"
                      />
                    )
                  )}
                  <span>
                    <span className="block text-sm font-medium text-primary">
                      {method.label}
                    </span>
                    <span className="mt-0.5 block text-xs text-muted">
                      {method.description}
                    </span>
                  </span>
                </label>
              ))}
            </div>
          </section>
        </div>

        {/* Order summary */}
        <aside className="rounded-2xl border border-border bg-surface p-6 shadow-sm lg:sticky lg:top-6">
          <h2 className="text-lg font-semibold text-primary">Order summary</h2>

          <ul className="mt-4 space-y-4 border-b border-border pb-4">
            {items.map((item) => (
              <li key={item.data.id} className="flex gap-3">
                <div className="relative h-16 w-12 shrink-0 overflow-hidden rounded-md bg-accent/20">
                  <Image
                    src={item.data.cover_page_url}
                    alt={`Cover of ${item.data.title}`}
                    fill
                    sizes="48px"
                    className="object-cover"
                  />
                </div>
                <div className="min-w-0 flex-1">
                  <p className="truncate text-sm font-medium text-primary">
                    {item.data.title}
                  </p>
                  <p className="text-xs text-muted">
                    {item.quantity} × {formatPrice(item.data.price)}
                  </p>
                </div>
                <p className="text-sm font-semibold text-primary tabular-nums">
                  {formatPrice(item.subtotal)}
                </p>
              </li>
            ))}
          </ul>

          <div className="mt-4 flex items-center justify-between">
            <span className="text-sm text-muted">Total</span>
            <span className="text-xl font-bold text-primary tabular-nums">
              {formatPrice(total)}
            </span>
          </div>

          <button
            type="button"
            onClick={handlePlaceOrder}
            disabled={isPlacing}
            className="mt-5 w-full rounded-lg bg-primary px-4 py-3 text-sm font-medium text-white transition-colors hover:bg-primary-light disabled:cursor-not-allowed disabled:opacity-60"
          >
            {isPlacing
              ? "Placing order…"
              : paymentMethod === "ESEWA"
                ? "Place order & pay with eSewa"
                : "Place order"}
          </button>

          <p className="mt-3 text-xs text-muted">
            {addressId === null
              ? "Save your delivery address to continue."
              : "Prices are recalculated on the server."}
          </p>
        </aside>
      </div>
    </main>
  );
}
