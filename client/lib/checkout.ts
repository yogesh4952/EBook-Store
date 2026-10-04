import type { ICartItem, ICartResponse } from "@/components/cart/CartItem";

export type PaymentMethod = "COD" | "ESEWA";

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
  // Absent for COD, present for ESEWA.
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
      (payload as ApiError | null)?.message ?? `Request failed (${response.status})`;
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