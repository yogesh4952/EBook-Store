import { BACKEND_URL } from "./config";

export interface Book {
  id: number;
  title: string;
  author_name: string;
  genre: string;
  category: string;
  pages: number;
  publication: string;
  price: number;
  units: number;
  cover_page_url?: string;
  seller_id?: number;
  status?: string;
  seller?: { user?: { first_name?: string } };
}

export interface Address {
  id?: number;
  city: string;
  delivery_address: string;
  user_id?: number;
}

export interface OrderItem {
  book_id: number;
  title: string;
  quantity: number;
  unit_price: number;
  subtotal: number;
  seller_id?: number;
}

export interface Order {
  order_id?: number;
  order_code: string;
  total_price: number;
  payment_status: string;
  order_status: string;
  address: string;
  items?: OrderItem[];
  esewa_payload?: Record<string, string> | null;
}

export interface User {
  id: number;
  first_name: string;
  last_name: string;
  email: string;
  role: string;
  phone_number?: string;
}

export async function apiRequest<T>(path: string, init?: RequestInit): Promise<T> {
  const normalizedPath = path.replace(/^\//, "");
  const baseUrl =
    typeof window === "undefined"
      ? (process.env.BACKEND_URL || BACKEND_URL) + "/api"
      : "/api/proxy";

  const response = await fetch(`${baseUrl}/${normalizedPath}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...init?.headers,
    },
    cache: "no-store",
  });

  const data = await response.json().catch(() => null);
  if (!response.ok) {
    throw new Error(data?.message || data?.error || "Request failed");
  }

  return data as T;
}

export async function listBooks(page = 1, limit = 20) {
  return apiRequest<{ data: Book[]; total: number; page: number; limit: number }>(
    `book/list-books?page=${page}&limit=${limit}`,
  );
}

export async function placeOrder(payload: {
  payment_method: "COD" | "ESEWA";
  user_address_id: number;
  items: { book_id: number; quantity: number; seller_id: number }[];
}) {
  return apiRequest<{ data: Order; message: string }>("order/place-order", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}