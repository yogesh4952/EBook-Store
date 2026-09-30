"use client";

import { useEffect } from "react";

// This is the type of the payload your Go backend just returned
interface EsewaPayload {
  amount: string;
  tax_amount: string;
  total_amount: string;
  transaction_uuid: string;
  product_code: string;
  product_service_charge: string;
  product_delivery_charge: string;
  success_url: string;
  failure_url: string;
  signed_field_names: string;
  signature: string;
}

interface CheckoutResponse {
  success: boolean;
  message: string;
  data: {
    order_id: number;
    order_code: string;
    total_price: number;
    payment_status: string;
    order_status: string;
    address: string;
    items: {
      book_id: number;
      title: string;
      quantity: number;
      unit_price: number;
      subtotal: number;
    }[];
    esewa_payload: EsewaPayload | null; // Will be null for COD
  };
}

export default function EsewaRedirect({
  response,
}: {
  response?: CheckoutResponse;
}) {
  useEffect(() => {
    if (!response?.data.esewa_payload) return;

    const payload = response.data.esewa_payload;

    // 2. Create a hidden form
    const form = document.createElement("form");
    form.method = "POST";

    form.action = `https://rc-epay.esewa.com.np/api/epay/main/v2/form`;

    // 3. Populate the form with the signed data from your Go backend
    Object.entries(payload).forEach(([key, value]) => {
      const input = document.createElement("input");
      input.type = "hidden";
      input.name = key;
      input.value = value;
      form.appendChild(input);
    });

    // 4. Append to the document and auto-submit
    document.body.appendChild(form);

    // Small delay to ensure the DOM is ready and the user sees the "Redirecting..." message
    setTimeout(() => {
      form.submit();
    }, 500);
  }, [response]);

  if (!response) {
    return (
      <div className="p-8 text-center">
        <h1 className="text-2xl font-bold">Payment status</h1>
        <p className="mt-2 text-muted">No payment response was supplied.</p>
      </div>
    );
  }

  return (
    <div className="p-8 text-center">
      <h1 className="text-2xl font-bold">Order Placed Successfully!</h1>
      <p className="mt-2">Your order code is: {response.data.order_code}</p>
    </div>
  );
}
