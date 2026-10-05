import type { ICartItem, ICartResponse } from "@/components/cart/CartItem";

export type PaymentMethod = "COD" | "ESEWA";
export type PaymentStatus = "PAID" | "PENDING" | "REFUND";
export type OrderStatus = "PLACED" | "DELIVERED" | "CANCELLED";

export interface OrderItemResponse {
  book_id: number;
  title: string;
  quantity: number;
  unit_price: number;
  subtotal: number;
}

export interface EsewaPayload {
  amount: string;
  total_amount: string;
  transaction_uuid: string;
  product_code: string;
  product_service_charge: string;
  product_delivery_charge: string;
  tax_amount: string;
  signed_field_names: string;
  signature: string;
  success_url: string;
  failure_url: string;
}

export interface PlaceOrderResponse {
  order_id: number;
  order_code: string;
  total_price: number;
  payment_status: string;
  order_status: string;
  address: string;
  items: OrderItemResponse[];
  esewa_payload?: EsewaPayload | null;
}

export interface PlaceOrderPayload {
  payment_method: PaymentMethod;
  user_address_id: number;
  items: { book_id: number; quantity: number; seller_id: number }[];
}

export interface ApiError {
  success: false;
  message: string;
}

const jsonHeaders = { "Content-Type": "application/json" } as const;

async function readJson<T>(response: Response): Promise<T> {
  const payload = await response.json().catch(() => null);
  if (!response.ok) {
    const message =
      (payload as ApiError | null)?.message ??
      `Request failed (${response.status})`;
    throw new Error(message);
  }
  return payload as T;
}

export async function getCart(): Promise<ICartResponse> {
  // Trailing slash matters: /api/cart 301-redirects, and the redirect drops the
  // Authorization header before it reaches Gin.
  const response = await fetch("/api/proxy/cart/", {
    method: "GET",
    credentials: "include",
    cache: "no-store",
  });
  const payload = await readJson<{ data: ICartResponse }>(response);
  return payload.data;
}

export async function placeOrder(body: PlaceOrderPayload) {
  const response = await fetch("/api/proxy/order/place-order", {
    method: "POST",
    credentials: "include",
    headers: jsonHeaders,
    body: JSON.stringify(body),
  });
  const payload = await readJson<{ data: PlaceOrderResponse }>(response);
  return payload.data;
}

/**
 * Add units of a book to the cart. Quantities are additive, so this is also how
 * the cart increments a line.
 *
 * The server rejects the call when the resulting quantity would exceed the
 * book's remaining units.
 */
export async function addToCart(input: {
  book_id: number;
  quantity: number;
}): Promise<void> {
  const response = await fetch("/api/proxy/cart/add-to-cart", {
    method: "POST",
    credentials: "include",
    headers: jsonHeaders,
    body: JSON.stringify(input),
  });
  await readJson<{ success: boolean; message: string }>(response);
}

export async function clearCart(): Promise<void> {
  const response = await fetch("/api/proxy/cart/clear-cart", {
    method: "POST",
    credentials: "include",
  });
  await readJson<{ success: boolean; message: string }>(response);
}

/**
 * Decrement a book's quantity in the cart, or drop it entirely when the
 * remaining quantity is exactly the amount removed.
 *
 * Two things to know about the endpoint:
 *  - It decrements with HIncrBy rather than deleting the field, so removing
 *    more than the cart holds leaves a negative count that makes GET /cart/
 *    fail to parse. Callers must never send a quantity above what is held.
 *  - It responds 201 with an "added to cart" message, so status and body text
 *    are both misleading on success. Success is read from `success` only.
 */
export async function removeFromCart(input: {
  book_id: number;
  quantity: number;
}): Promise<void> {
  const response = await fetch("/api/proxy/cart/remove-from-cart", {
    method: "POST",
    credentials: "include",
    headers: jsonHeaders,
    body: JSON.stringify(input),
  });
  await readJson<{ success: boolean; message: string }>(response);
}

export interface UserAddress {
  id: number;
  city: string;
  delivery_address: string;
  user_id: number;
  created_at?: string;
  updated_at?: string;
  deleted_at?: string | null;
}

// user_id is deliberately omitted: the backend takes the owner from the JWT and
// currently rejects a body without it. Sending one from the client would let a
// caller attach an address to another account, so it is left out on purpose.
export async function addAddress(input: {
  city: string;
  delivery_address: string;
}): Promise<UserAddress> {
  const response = await fetch("/api/proxy/address/add-address", {
    method: "POST",
    credentials: "include",
    headers: jsonHeaders,
    body: JSON.stringify(input),
  });
  const payload = await readJson<{ data: UserAddress }>(response);
  return payload.data;
}

// The endpoint scopes results to the authenticated user, so no client-side
// filtering is needed. Still normalised to a list so an empty result and a
// single row render through the same code path.
export async function listUserAddresses(): Promise<UserAddress[]> {
  const response = await fetch("/api/proxy/address/", {
    method: "GET",
    credentials: "include",
    cache: "no-store",
  });
  const payload = await readJson<{ data: UserAddress[] | null }>(response);

  return payload.data ?? [];
}

export const cartToOrderItems = (items: ICartItem[]) =>
  items.map((item) => ({
    book_id: item.data.id,
    quantity: item.quantity,
    seller_id: item.data.seller_id ?? 0,
  }));

/**
 * A line item as it comes back inside a listed order.
 *
 * The list endpoint returns the whole `OrderItem` row, which carries the joined
 * `book` record. That is the only place a cover image survives, so `book` is
 * used for the thumbnail and treated as optional for rows written before the
 * preload existed.
 */
export interface OrderLineItem {
  order_id: number;
  book_id: number;
  quantity: number;
  unit_price: number;
  created_at?: string;
  updated_at?: string;
  book?: {
    id: number;
    title: string;
    author_name: string;
    cover_page_url: string;
  };
}

export interface UserOrder {
  id: number;
  created_at: string;
  updated_at: string;
  order_code: string;
  transaction_uuid: string;
  transaction_code: string;
  payment_method: PaymentMethod;
  payment_status: PaymentStatus;
  order_status: OrderStatus | "";
  total_price: number;
  shipping_city: string;
  shipping_delivery_address: string;
  order_items: OrderLineItem[];
}

export async function listUserOrders(): Promise<UserOrder[]> {
  const response = await fetch("/api/proxy/order/list-user-order", {
    method: "GET",
    credentials: "include",
    cache: "no-store",
  });
  const payload = await readJson<{ data: UserOrder[] | null }>(response);

  return payload.data ?? [];
}
