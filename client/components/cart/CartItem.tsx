"use client";

import Image from "next/image";
import { HiMinus, HiPlus, HiTrash } from "react-icons/hi2";

export interface ICartBook {
  id: number;
  title: string;
  author_name: string;
  genre: string;
  category: string;
  pages: number;
  publication: string;
  price: number;
  units: number;
  cover_page_url: string;
  seller_id: number | null;
  status: "in_stock" | "out_of_stock";
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}

export interface ICartItem {
  data: ICartBook;
  quantity: number;
  subtotal: number;
}

export interface ICartResponse {
  items: ICartItem[];
  total: number;
}

interface CartItemProps {
  item: ICartItem;
  onQuantityChange?: (bookId: number, quantity: number) => void;
  onRemove?: (bookId: number) => void;
  isUpdating?: boolean;
  disabled?: boolean;
}

const formatPrice = (value: number) =>
  `Rs. ${new Intl.NumberFormat("en-NP", {
    maximumFractionDigits: 0,
  }).format(value)}`;

const LOW_STOCK_THRESHOLD = 5;

const CartItem = ({
  item,
  onQuantityChange,
  onRemove,
  isUpdating = false,
  disabled = false,
}: CartItemProps) => {
  const { data: book, quantity, subtotal } = item;
  const maxQuantity = Math.max(book.units, 1);
  const isOutOfStock = book.status === "out_of_stock" || book.units === 0;
  const isLowStock = !isOutOfStock && book.units <= LOW_STOCK_THRESHOLD;
  const isBusy = disabled || isUpdating;

  return (
    <li
      className={`flex gap-3 border-b border-border py-4 last:border-b-0 ${
        isUpdating ? "opacity-60" : ""
      }`}
    >
      <div className="relative h-20 w-16 shrink-0 overflow-hidden rounded-md bg-accent/20">
        <Image
          src={book.cover_page_url}
          alt={`Cover of ${book.title}`}
          fill
          sizes="64px"
          className="object-cover"
        />
      </div>

      <div className="flex min-w-0 flex-1 flex-col">
        <div className="flex items-start justify-between gap-2">
          <div className="min-w-0">
            <h3 className="truncate text-sm font-semibold leading-snug text-primary">
              {book.title}
            </h3>
            <p className="mt-0.5 truncate text-xs text-muted">
              {book.author_name}
            </p>
          </div>

          {onRemove && (
            <button
              type="button"
              onClick={() => onRemove(book.id)}
              disabled={isBusy}
              aria-label={`Remove ${book.title} from cart`}
              className="shrink-0 rounded-md p-2 text-muted transition-colors hover:bg-error/10 hover:text-error focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-error disabled:pointer-events-none disabled:opacity-40"
            >
              <HiTrash size={16} />
            </button>
          )}
        </div>

        <div className="mt-auto flex items-end justify-between gap-3 pt-2">
          <div
            className="flex items-center rounded-lg border border-border"
            role="group"
            aria-label={`Quantity for ${book.title}`}
          >
            <button
              type="button"
              onClick={() => onQuantityChange?.(book.id, quantity - 1)}
              disabled={isBusy || quantity <= 1}
              aria-label="Decrease quantity"
              className="flex h-10 w-10 items-center justify-center rounded-l-lg text-muted transition-colors hover:bg-accent/20 hover:text-primary focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-primary disabled:pointer-events-none disabled:opacity-30"
            >
              <HiMinus size={14} />
            </button>

            <span
              aria-live="polite"
              aria-label={`Quantity ${quantity}`}
              className="min-w-8 text-center text-sm font-medium text-text tabular-nums"
            >
              {quantity}
            </span>

            <button
              type="button"
              onClick={() => onQuantityChange?.(book.id, quantity + 1)}
              disabled={isBusy || quantity >= maxQuantity}
              aria-label="Increase quantity"
              className="flex h-10 w-10 items-center justify-center rounded-r-lg text-muted transition-colors hover:bg-accent/20 hover:text-primary focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-primary disabled:pointer-events-none disabled:opacity-30"
            >
              <HiPlus size={14} />
            </button>
          </div>

          <div className="text-right">
            <p className="text-sm font-semibold text-primary tabular-nums">
              {formatPrice(subtotal)}
            </p>
            {quantity > 1 && (
              <p className="text-xs text-muted tabular-nums">
                {formatPrice(book.price)} each
              </p>
            )}
          </div>
        </div>

        {isOutOfStock ? (
          <p className="mt-2 text-xs font-medium text-error">
            Out of stock — remove to continue
          </p>
        ) : isLowStock ? (
          <p className="mt-2 text-xs text-muted">
            Only {book.units} left in stock
          </p>
        ) : null}
      </div>
    </li>
  );
};

export default CartItem;