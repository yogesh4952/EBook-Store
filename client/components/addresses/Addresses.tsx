"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import {
  HiOutlineMapPin,
  HiOutlinePlus,
  HiOutlineTrash,
} from "react-icons/hi2";
import { toast } from "sonner";

import {
  addAddress,
  listUserAddresses,
  type UserAddress,
} from "@/lib/checkout";

const Addresses = () => {
  const [addresses, setAddresses] = useState<UserAddress[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isAdding, setIsAdding] = useState(false);

  const [city, setCity] = useState("");
  const [deliveryAddress, setDeliveryAddress] = useState("");

  /** Re-reads the list from the server after a write. */
  const reload = async () => {
    try {
      setAddresses(await listUserAddresses());
    } catch (error) {
      toast.error(
        error instanceof Error
          ? error.message
          : "Could not load your addresses",
      );
    }
  };

  // Fetching on mount is an external-system sync, not derived state, so this
  // effect is the right tool despite the lint rule: there is no server render
  // to read the addresses from.
  useEffect(() => {
    const controller = new AbortController();

    listUserAddresses()
      .then((data) => {
        if (controller.signal.aborted) return;
        setAddresses(data);
      })
      .catch((error: unknown) => {
        if (controller.signal.aborted) return;
        toast.error(
          error instanceof Error
            ? error.message
            : "Could not load your addresses",
        );
      })
      .finally(() => {
        if (!controller.signal.aborted) setIsLoading(false);
      });

    return () => controller.abort();
  }, []);

  const resetForm = () => {
    setCity("");
    setDeliveryAddress("");
  };

  const handleAdd = async (event: React.FormEvent) => {
    event.preventDefault();

    if (isAdding) return;

    const trimmedCity = city.trim();
    const trimmedAddress = deliveryAddress.trim();

    if (!trimmedCity || !trimmedAddress) {
      toast.error("City and delivery address are both required");
      return;
    }

    setIsAdding(true);
    try {
      await addAddress({
        city: trimmedCity,
        delivery_address: trimmedAddress,
      });
      await reload();
      resetForm();
      toast.success("Delivery address saved");
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Could not save the address",
      );
    } finally {
      setIsAdding(false);
    }
  };

  /**
   * The backend exposes only list and create, so there is no delete to call.
   * Hiding the row would be a lie that the next page load unmasks, so the
   * button says plainly that the API cannot do this yet.
   */
  const handleDeleteUnavailable = () => {
    toast.info("Deletion is not supported yet", {
      description:
        "The API currently exposes only listing and adding addresses.",
    });
  };

  if (isLoading) {
    return (
      <div className="flex min-h-[60vh] flex-col items-center justify-center gap-3">
        <span className="size-7 animate-spin rounded-full border-2 border-border border-t-primary" />
        <p className="text-sm text-muted">Loading your addresses…</p>
      </div>
    );
  }

  return (
    <div className="grid gap-6 lg:grid-cols-[1fr_380px] lg:items-start">
      {/* Saved addresses */}
      <section>
        <h2 className="text-lg font-semibold text-primary">Saved addresses</h2>

        {addresses.length === 0 ? (
          <div className="mt-4 rounded-2xl border border-border bg-surface px-6 py-12 text-center shadow-sm">
            <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-accent/20">
              <HiOutlineMapPin size={24} className="text-primary" />
            </div>
            <p className="mt-4 text-sm font-medium text-primary">
              No delivery addresses yet
            </p>
            <p className="mt-1 text-sm text-muted">
              Add one using the form, then pick it during checkout.
            </p>
          </div>
        ) : (
          <ul className="mt-4 space-y-3">
            {addresses.map((address, index) => (
              <li
                key={address.id}
                className="flex items-start gap-3 rounded-2xl border border-border bg-surface p-5 shadow-sm transition-colors hover:border-accent"
              >
                <div className="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-accent/20">
                  <HiOutlineMapPin size={16} className="text-primary" />
                </div>

                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <p className="text-sm font-semibold text-primary">
                      {address.city}
                    </p>
                    {index === 0 && (
                      <span
                        title="Checkout preselects the first saved address"
                        className="rounded-full bg-success/10 px-2 py-0.5 text-[11px] font-medium text-success"
                      >
                        Preselected
                      </span>
                    )}
                  </div>
                  <p className="mt-1 break-words text-sm leading-relaxed text-muted">
                    {address.delivery_address}
                  </p>
                </div>

                <button
                  type="button"
                  onClick={handleDeleteUnavailable}
                  aria-label={`Delete address in ${address.city}`}
                  title="Deleting addresses is not supported yet"
                  className="shrink-0 cursor-not-allowed rounded-md p-2 text-muted/50"
                  disabled
                >
                  <HiOutlineTrash size={16} />
                </button>
              </li>
            ))}
          </ul>
        )}
      </section>

      {/* Add form */}
      <aside className="rounded-2xl border border-border bg-surface p-6 shadow-sm lg:sticky lg:top-6">
        <h2 className="flex items-center gap-2 text-lg font-semibold text-primary">
          <HiOutlinePlus size={18} />
          Add an address
        </h2>

        <form onSubmit={handleAdd} className="mt-4 space-y-4">
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
              onChange={(event) => setCity(event.target.value)}
              placeholder="Kathmandu"
              autoComplete="address-level2"
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
            <textarea
              id="delivery_address"
              name="delivery_address"
              value={deliveryAddress}
              onChange={(event) => setDeliveryAddress(event.target.value)}
              placeholder="Street, landmark, house number"
              rows={3}
              autoComplete="street-address"
              className="w-full resize-none rounded-lg border border-border bg-background px-3 py-2.5 text-sm outline-none transition-colors focus:border-primary"
            />
          </div>

          <button
            type="submit"
            disabled={isAdding}
            className="w-full rounded-lg bg-primary px-4 py-3 text-sm font-medium text-white transition-colors hover:bg-primary-light disabled:cursor-not-allowed disabled:opacity-60"
          >
            {isAdding ? "Saving…" : "Save address"}
          </button>

          <p className="text-xs text-muted">
            Addresses are used to calculate delivery for every order you place.
          </p>
        </form>

        <Link
          href="/checkout"
          className="mt-4 block text-center text-sm text-muted underline underline-offset-4 transition-colors hover:text-primary"
        >
          Continue to checkout
        </Link>
      </aside>
    </div>
  );
};

export default Addresses;
