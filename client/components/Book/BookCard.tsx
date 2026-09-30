"use client";

import type { Book } from "@/lib/api";
import { addToCart } from "@/lib/cart";
import { useState } from "react";

export type Ibook = Book;

interface BookCardProps {
  book: Ibook;
}

const BookCard = ({ book }: BookCardProps) => {
  const [added, setAdded] = useState(false);

  function handleAddToCart() {
    addToCart(book);
    setAdded(true);
    window.setTimeout(() => setAdded(false), 1400);
  }

  return (
    <div className="group overflow-hidden rounded-xl border border-border bg-surface text-text shadow-sm transition-all duration-300 hover:-translate-y-1 hover:shadow-lg">
      {/* Book Cover */}
      <div
        className="relative h-72 w-full overflow-hidden bg-accent/20 bg-cover bg-center"
        style={{
          backgroundImage: book.cover_page_url
            ? `url(${book.cover_page_url})`
            : undefined,
        }}
      >
        {/* Overlay */}
        <div className="absolute inset-0 bg-black/0 transition-all duration-300 group-hover:bg-black/10" />

        {/* Genre */}
        <span className="absolute left-3 top-3 rounded-full bg-background/90 px-3 py-1 text-xs font-medium text-primary backdrop-blur-sm">
          {book.genre}
        </span>
      </div>

      {/* Content */}
      <div className="p-4">
        <h2 className="line-clamp-2 text-lg font-semibold leading-tight text-primary">
          {book.title}
        </h2>

        <p className="mt-1 text-sm text-muted">by {book.author_name}</p>

        {/* Price */}
        <div className="mt-4 flex items-center justify-between">
          <span className="text-xl font-bold text-primary">
            Rs. {book.price}
          </span>

          <button
            onClick={handleAddToCart}
            className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-light"
          >
            {added ? "Added" : "Add to cart"}
          </button>
        </div>
      </div>
    </div>
  );
};

export default BookCard;
