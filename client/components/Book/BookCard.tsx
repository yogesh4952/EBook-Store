"use client";

import Image from "next/image";
import { useState } from "react";
import { HiOutlineShoppingCart } from "react-icons/hi2";
import { toast } from "sonner";

import Sidebar from "../cart/Sidebar";
import { addToCart } from "@/lib/checkout";
import { refreshCart } from "@/store/useCartStore";

export interface Ibook {
  id: number;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
  title: string;
  author_name: string;
  genre: string;
  category: string;
  pages: number;
  publication: string;
  price: number;
  units: number;
  status: "in_stock" | "out_of_stock";
  cover_page_url: string;
  seller_id: number | null;
  seller?: {
    id: number;
    created_at: string;
    updated_at: string;
    deleted_at: string | null;
    user_id: number;
    seller_number: number;
    user?: {
      id: number;
      created_at: string;
      updated_at: string;
      deleted_at: string | null;
      first_name: string;
      last_name: string;
      email: string;
      role: "seller";
      phone_number: string;
    };
  };
}

interface BookCardProps {
  book: Ibook;
  priority?: boolean;
}

const LOW_STOCK_THRESHOLD = 5;

const formatPrice = (value: number) =>
  `Rs. ${new Intl.NumberFormat("en-NP", { maximumFractionDigits: 0 }).format(value)}`;

const BookCard = ({ book, priority }: BookCardProps) => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [isAdding, setIsAdding] = useState(false);

  const isOutOfStock = book.status === "out_of_stock" || book.units === 0;
  const isLowStock = !isOutOfStock && book.units <= LOW_STOCK_THRESHOLD;
  const sellerName = book.seller?.user?.first_name;

  const handleAddToCart = async () => {
    if (isAdding || isOutOfStock) return;

    setIsAdding(true);
    try {
      await addToCart({ book_id: book.id, quantity: 1 });
      refreshCart();
      setIsSidebarOpen(true);
      toast.success(`${book.title} added to cart`);
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Could not add to cart",
      );
    } finally {
      setIsAdding(false);
    }
  };

  return (
    <>
      <article className="group flex flex-col overflow-hidden rounded-xl border border-border bg-surface shadow-sm transition-all duration-200 hover:-translate-y-0.5 hover:border-accent hover:shadow-md">
        <div className="relative aspect-[3/4] w-full overflow-hidden bg-accent/20">
          <Image
            src={book.cover_page_url}
            alt={`Cover of ${book.title}`}
            fill
            sizes="(max-width: 640px) 50vw, (max-width: 1024px) 33vw, 20vw"
            priority={priority}
            className="object-cover transition-transform duration-300 group-hover:scale-[1.03]"
          />

          <span className="absolute left-2 top-2 rounded-full bg-background/90 px-2 py-0.5 text-[11px] font-medium text-primary backdrop-blur-sm">
            {book.genre}
          </span>

          {isOutOfStock && (
            <span className="absolute inset-0 flex items-center justify-center bg-text/60 text-sm font-semibold text-white">
              Out of stock
            </span>
          )}
        </div>

        <div className="flex flex-1 flex-col p-3">
          <h3
            title={book.title}
            className="line-clamp-2 min-h-9 text-sm font-semibold leading-tight text-primary"
          >
            {book.title}
          </h3>

          <p className="mt-0.5 truncate text-xs text-muted">
            {book.author_name}
          </p>

          <dl className="mt-2 space-y-0.5 text-[11px] leading-tight text-muted">
            <div className="flex gap-1">
              <dt className="shrink-0">Pages</dt>
              <dd className="truncate font-medium text-text tabular-nums">
                {book.pages}
              </dd>
            </div>
            <div className="flex gap-1">
              <dt className="shrink-0">Publisher</dt>
              <dd className="truncate font-medium text-text" title={book.publication}>
                {book.publication}
              </dd>
            </div>
            {sellerName && (
              <div className="flex gap-1">
                <dt className="shrink-0">Seller</dt>
                <dd className="truncate font-medium text-text">{sellerName}</dd>
              </div>
            )}
          </dl>

          <div className="mt-auto pt-3">
            <div className="flex items-baseline justify-between gap-2">
              <span className="text-base font-bold text-primary tabular-nums">
                {formatPrice(book.price)}
              </span>
              {!isOutOfStock && (
                <span
                  className={`text-[11px] tabular-nums ${
                    isLowStock ? "font-medium text-error" : "text-muted"
                  }`}
                >
                  {isLowStock ? `Only ${book.units} left` : `${book.units} in stock`}
                </span>
              )}
            </div>

            <button
              type="button"
              onClick={handleAddToCart}
              disabled={isAdding || isOutOfStock}
              aria-label={`Add ${book.title} to cart`}
              className="mt-2 flex w-full items-center justify-center gap-1.5 rounded-lg bg-primary px-3 py-2 text-xs font-medium text-white transition-colors hover:bg-primary-light focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary disabled:cursor-not-allowed disabled:opacity-50"
            >
              <HiOutlineShoppingCart size={14} />
              {isOutOfStock
                ? "Unavailable"
                : isAdding
                  ? "Adding…"
                  : "Add to cart"}
            </button>
          </div>
        </div>
      </article>

      {isSidebarOpen && <Sidebar onClose={() => setIsSidebarOpen(false)} />}
    </>
  );
};

export default BookCard;