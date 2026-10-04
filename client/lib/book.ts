import type { Ibook } from "@/components/Book/BookCard";

const PAGE_SIZE = 24;

export async function ListBooks(page = 1, limit = PAGE_SIZE): Promise<Ibook[]> {
  const res = await fetch(
    `http://localhost:8080/api/book/list-books?page=${page}&limit=${limit}`,
    { next: { revalidate: 60 } },
  );

  if (!res.ok) {
    throw new Error(`Failed to load books (${res.status})`);
  }

  const payload = await res.json();
  return payload.data ?? [];
}