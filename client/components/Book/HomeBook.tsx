"use client";

import React from "react";
import Link from "next/link";
import { useMemo } from "react";
import { HiOutlineChevronRight } from "react-icons/hi2";

import BookCard, { type Ibook } from "./BookCard";
import { getCategories, getGenres } from "@/lib/book-utils";

interface SectionProps {
  title: string;
  books: Ibook[];
  emptyMessage?: string;
}

const Section = ({ title, books, emptyMessage }: SectionProps) => {
  if (!books.length) {
    return (
      <div className="rounded-xl border border-border bg-surface px-6 py-10 text-center text-sm text-muted">
        {emptyMessage || "No books available."}
      </div>
    );
  }

  return (
    <section className="relative">
      <div className="mb-3 flex items-baseline justify-between px-1">
        <h2 className="text-lg font-semibold text-primary">{title}</h2>
      </div>

      <div className="relative group">
        <div className="scrollbar-slim flex gap-4 overflow-x-auto pb-6 snap-x snap-mandatory">
          {books.map((book) => (
            <div
              key={book.id}
              className="w-44 shrink-0 snap-start sm:w-48 md:w-52"
            >
              <BookCard book={book} priority />
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

interface HomeBookProps {
  books: Ibook[];
}

const HomeBook = ({ books }: HomeBookProps) => {
  const popularBooks = useMemo(() => {
    return books
      .filter((book) => book.status === "in_stock")
      .sort((a, b) => b.units - a.units)
      .slice(0, 12);
  }, [books]);

  const categoryData = useMemo(() => getCategories(books), [books]);
  const genreList = useMemo(() => getGenres(books), [books]);

  const budgetBooks = useMemo(() => {
    return books.filter(
      (book) => book.status === "in_stock" && book.price < 1000,
    );
  }, [books]);

  const newArrivalBooks = useMemo(() => {
    return books.filter((book) => book.status === "in_stock").slice(0, 12);
  }, [books]);

  return (
    <div className="space-y-10 pb-10 mt-5">
      <section>
        <h1 className="text-2xl font-bold text-primary">Browse Books</h1>
        <p className="mt-1 text-sm text-muted">
          {books.length} titles available
        </p>
      </section>

      <Section title="Popular & Trending" books={popularBooks} />

      <div>
        <div className="mb-3 flex items-baseline justify-between px-1">
          <h2 className="text-lg font-semibold text-primary">
            Browse by Category
          </h2>
          <Link
            href="/book"
            className="flex items-center gap-1 text-sm font-medium text-primary hover:underline"
          >
            View all
            <HiOutlineChevronRight size={16} />
          </Link>
        </div>

        <ul className="flex flex-wrap gap-2">
          {categoryData.map(({ category, count }) => (
            <li key={category}>
              <Link
                href={`/book?category=${encodeURIComponent(category)}`}
                className="block rounded-full border border-border bg-surface px-3 py-1 text-sm font-medium text-primary hover:border-accent hover:bg-accent/10"
              >
                {category} ({count})
              </Link>
            </li>
          ))}
        </ul>
      </div>

      <div>
        <div className="mb-3 flex items-baseline justify-between px-1">
          <h2 className="text-lg font-semibold text-primary">Top Genres</h2>
          <Link
            href="/book"
            className="flex items-center gap-1 text-sm font-medium text-primary hover:underline"
          >
            View all
            <HiOutlineChevronRight size={16} />
          </Link>
        </div>

        <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6">
          {genreList.slice(0, 10).map((genre) => (
            <Link
              key={genre}
              href={`/book?genre=${encodeURIComponent(genre)}`}
              className="rounded-lg border border-border bg-surface p-3 text-center text-sm font-medium text-primary transition-colors hover:border-accent hover:bg-accent/10"
            >
              {genre}
            </Link>
          ))}
        </div>
      </div>

      <Section title="Budget Pick - Under Rs. 1,000" books={budgetBooks} />

      <Section title="New Arrivals" books={newArrivalBooks} />

      <div className="mt-12 flex justify-center">
        <Link
          href="/book"
          className="rounded-lg bg-primary px-6 py-3 text-base font-medium text-white transition-colors hover:bg-primary-light"
        >
          Browse All Books
        </Link>
      </div>
    </div>
  );
};

export default HomeBook;
