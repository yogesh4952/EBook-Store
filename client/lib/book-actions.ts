"use server";

import type { Ibook } from "@/components/Book/BookCard";

export async function getAllBooks(page = 1, limit = 24): Promise<Ibook[]> {
  const BACKEND_URI= process.env.BACKEND_URL
  const res = await fetch(
    `${BACKEND_URI}/api/book/list-books?page=${page}&limit=${limit}`,
    { next: { revalidate: 60 } },
  );

  if (!res.ok) {
    throw new Error(`Failed to load books (${res.status})`);
  }

  const payload = await res.json();
  return payload.data ?? [];
}
