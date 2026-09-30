"use client";

import BookCard from "@/components/Book/BookCard";
import { listBooks, type Book } from "@/lib/api";
import { useEffect, useState } from "react";

export default function BooksPage() {
  const [books, setBooks] = useState<Book[]>([]);
  const [error, setError] = useState("");

  useEffect(() => {
    listBooks()
      .then((result) => setBooks(result.data))
      .catch((reason) => setError(reason.message));
  }, []);

  return (
    <main className="mx-auto max-w-7xl px-4 py-8">
      <div className="mb-8 flex items-end justify-between gap-4">
        <div>
          <p className="text-sm font-semibold uppercase tracking-[0.18em] text-primary-light">
            Explore
          </p>
          <h1 className="mt-2 text-4xl font-bold text-primary">All books</h1>
        </div>
        <a
          href="/seller/books/new"
          className="rounded-lg bg-primary px-4 py-3 text-sm font-semibold text-white"
        >
          Publish a book
        </a>
      </div>
      {error && (
        <p className="rounded-lg bg-red-50 p-4 text-red-700">{error}</p>
      )}
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        {books.map((book) => (
          <BookCard book={book} key={book.id} />
        ))}
      </div>
    </main>
  );
}
