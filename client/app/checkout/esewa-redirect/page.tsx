"use client";

import { useEffect, useState } from "react";

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
    items: any[];
    esewa_payload: EsewaPayload | null; // Will be null for COD
  };
}

export default function EsewaRedirect({
  response,
}: {
  response: CheckoutResponse;
}) {
  const [isRedirecting, setIsRedirecting] = useState(true);

  useEffect(() => {
    // 1. Safety check: Do we have the payload?
    if (!response.data.esewa_payload) {
      setIsRedirecting(false);
      return;
    }

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
    document.body.appendChildorm);

  // Small delay to ensure the DOM is ready and the user sees the "Redirecting..." message
  setTimeout(() => {
    form.submit();
  }, 500);
}, [response]);

if (isRedirecting) {
  return (
    <div className="flex flex-col items-center justify-center h-screen bg-gray-50">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mb-4"></div>
      <p className="text-lg font-semibold text-gray-700">
        Securely redirecting to eSewa...
      </p>
      <p className="text-sm text-gray-500 mt-2">
        Please do not close this window.
      </p>
    </div>
  );
}

// Fallback if there was no payload (e.g., user chose COD)
return (
  <div className="p-8 text-center">
    <h1 className="text-2xl font-bold">Order Placed Successfully!</h1>
    <p className="mt-2">Your order code is: {response.data.order_code}</p>
  </div>
);
}
