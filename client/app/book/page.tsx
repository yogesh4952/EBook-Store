import React from "react";
import Link from "next/link";
import { HiOutlineChevronLeft } from "react-icons/hi2";

import BookCard, { type Ibook } from "@/components/Book/BookCard";
import { getAllBooks } from "@/lib/book-actions";
import { getGenresFromApi, getCategoriesFromApi } from "@/lib/book-api";

export default async function BookPage({
  searchParams,
}: {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
}) {
  const params = await searchParams;
  const books: Ibook[] = await getAllBooks(1, 100);

  const genreFilter =
    typeof params.genre === "string" ? params.genre : undefined;
  const categoryFilter =
    typeof params.category === "string" ? params.category : undefined;
  const searchQuery =
    typeof params.q === "string" ? params.q.toLowerCase() : undefined;
  const minPrice = params.min ? parseInt(params.min as string) : undefined;
  const maxPrice = params.max ? parseInt(params.max as string) : undefined;
  const stockFilter = params.stock === "in_stock" ? true : undefined;

  const filtered = books.filter((b) => {
    if (genreFilter && b.genre !== genreFilter) return false;
    if (categoryFilter && b.category !== categoryFilter) return false;
    if (searchQuery) {
      const q = searchQuery.toLowerCase();
      if (
        !b.title.toLowerCase().includes(q) &&
        !b.author_name.toLowerCase().includes(q)
      )
        return false;
    }
    if (minPrice !== undefined && b.price < minPrice) return false;
    if (maxPrice !== undefined && b.price > maxPrice) return false;
    if (stockFilter && b.status !== "in_stock") return false;
    return true;
  });

  const categories = await getCategoriesFromApi();
  const genres = await getGenresFromApi();

  return (
    <div className="min-h-screen bg-background px-6 py-10">
      <div className="mx-auto max-w-7xl space-y-8">
        <div className="flex items-center gap-4">
          <Link
            href="/"
            className="flex items-center gap-1 text-sm text-muted hover:text-primary"
          >
            <HiOutlineChevronLeft size={16} /> Home
          </Link>
          <span className="text-muted">/</span>
          <h1 className="text-2xl font-bold text-primary">Browse Books</h1>
        </div>

        <div className="rounded-xl border border-border bg-surface p-4 shadow-sm">
          <form method="get" className="flex flex-wrap gap-3 items-end">
            <div>
              <label
                htmlFor="q"
                className="block text-xs font-medium text-muted mb-1"
              >
                Search
              </label>
              <input
                id="q"
                name="q"
                type="text"
                defaultValue={searchQuery || ""}
                placeholder="Title, author..."
                className="h-9 w-48 rounded-lg border border-border bg-background px-3 text-sm outline-none focus:border-primary"
              />
            </div>
            <div>
              <label
                htmlFor="genre"
                className="block text-xs font-medium text-muted mb-1"
              >
                Genre
              </label>
              <select
                id="genre"
                name="genre"
                defaultValue={genreFilter || ""}
                className="h-9 rounded-lg border border-border bg-background px-2 text-sm outline-none focus:border-primary"
              >
                <option value="">All genres</option>
                {genres.map((g, i) => (
                  <option key={i} value={g}>
                    {g}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label
                htmlFor="category"
                className="block text-xs font-medium text-muted mb-1"
              >
                Category
              </label>
              <select
                id="category"
                name="category"
                defaultValue={categoryFilter || ""}
                className="h-9 rounded-lg border border-border bg-background px-2 text-sm outline-none focus:border-primary"
              >
                <option value="">All categories</option>
                {categories.map((category, index) => (
                  <option key={index} value={category}>
                    {category}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label
                htmlFor="min"
                className="block text-xs font-medium text-muted mb-1"
              >
                Min Rs.
              </label>
              <input
                id="min"
                name="min"
                type="number"
                defaultValue={minPrice || ""}
                className="h-9 w-24 rounded-lg border border-border bg-background px-3 text-sm outline-none focus:border-primary"
              />
            </div>
            <div>
              <label
                htmlFor="max"
                className="block text-xs font-medium text-muted mb-1"
              >
                Max Rs.
              </label>
              <input
                id="max"
                name="max"
                type="number"
                defaultValue={maxPrice || ""}
                className="h-9 w-24 rounded-lg border border-border bg-background px-3 text-sm outline-none focus:border-primary"
              />
            </div>
            <div>
              <label
                htmlFor="stock"
                className="block text-xs font-medium text-muted mb-1"
              >
                Stock
              </label>
              <select
                id="stock"
                name="stock"
                defaultValue={stockFilter ? "in_stock" : ""}
                className="h-9 rounded-lg border border-border bg-background px-2 text-sm outline-none focus:border-primary"
              >
                <option value="">Any</option>
                <option value="in_stock">In stock</option>
              </select>
            </div>
            <button
              type="submit"
              className="h-9 rounded-lg bg-primary px-4 text-sm font-medium text-white hover:bg-primary-light"
            >
              Apply
            </button>
          </form>
        </div>

        <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
          {filtered.map((book) => (
            <BookCard book={book} key={book.id} />
          ))}
        </div>

        {filtered.length === 0 && (
          <p className="text-center text-sm text-muted">
            No books match your filters.
          </p>
        )}

        <div className="text-sm text-muted">
          {filtered.length} result{filtered.length !== 1 ? "s" : ""}
        </div>
      </div>
    </div>
  );
}
