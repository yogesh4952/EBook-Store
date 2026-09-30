import type { Book } from "./api";

export interface CartItem {
  book: Book;
  quantity: number;
}

const CART_KEY = "ebook-store-cart";

export function readCart(): CartItem[] {
  if (typeof window === "undefined") return [];

  try {
    const value = JSON.parse(localStorage.getItem(CART_KEY) || "[]");
    return Array.isArray(value) ? value : [];
  } catch {
    return [];
  }
}

export function writeCart(items: CartItem[]) {
  localStorage.setItem(CART_KEY, JSON.stringify(items));
  window.dispatchEvent(new Event("cart-updated"));
}

export function addToCart(book: Book) {
  const items = readCart();
  const existing = items.find((item) => item.book.id === book.id);

  if (existing) existing.quantity += 1;
  else items.push({ book, quantity: 1 });

  writeCart(items);
}

export function removeFromCart(bookId: number) {
  writeCart(readCart().filter((item) => item.book.id !== bookId));
}

export function updateCartQuantity(bookId: number, quantity: number) {
  const items = readCart();
  const item = items.find((entry) => entry.book.id === bookId);

  if (!item) return;
  if (quantity < 1) removeFromCart(bookId);
  else {
    item.quantity = quantity;
    writeCart(items);
  }
}