"use server";

import type { Ibook } from "@/components/Book/BookCard";

export async function getAllBooks(page = 1, limit = 24): Promise<Ibook[]> {
  const BACKEND_URI= process.env.BACKEND_URL
  let res: Response;
  try {
    res = await fetch(
      `${BACKEND_URI}/api/book/list-books?page=${page}&limit=${limit}`,
      { next: { revalidate: 60 } },
    );
  } catch (error) {
    console.warn("Failed to connect to book API - returning empty array for build", error);
    return [];
  }

  if (!res.ok) {
    console.warn(`Failed to load books (${res.status}) - returning empty array for build`);
    return [];
  }

  const payload = await res.json();
  return payload.data ?? [];
}
