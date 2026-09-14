"use client";

import { useState } from "react";
import EsewaRedirect from "./esewa-redirect/page";

export default function CheckoutPage() {
  const [isLoading, setIsLoading] = useState(false);
  const [checkoutResponse, setCheckoutResponse] = useState<any>(null);

  const handlePlaceOrder = async () => {
    setIsLoading(true);

    try {
      // 1. Call your Go backend
      const response = await fetch(
        "http://localhost:3000/api/proxy/order/place-order",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            // Add your auth token here if needed
            // "Authorization": `Bearer ${token}`
          },
          body: JSON.stringify({
            payment_method: "ESEWA", // Or "COD"
            user_address_id: 1,
            items: [
              {
                book_id: 5,
                quantity: 1,
                seller_id: 1,
              },
            ],
          }),
        },
      );

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || "Failed to place order");
      }

      // 2. Save the response.
      // If it has an esewa_payload, the component below will automatically trigger the redirect.
      setCheckoutResponse(data);
    } catch (error) {
      console.error("Checkout error:", error);
      alert("Something went wrong. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  // 3. CONDITIONAL RENDERING:
  // If we have a response AND it contains an eSewa payload, show the redirect component.
  if (checkoutResponse?.data?.esewa_payload) {
    return <EsewaRedirect response={checkoutResponse} />;
  }

  // 4. If it was COD (no payload), show a normal success message
  if (checkoutResponse?.success && !checkoutResponse?.data?.esewa_payload) {
    return (
      <div className="p-8 text-center">
        <h1 className="text-2xl font-bold text-green-600">
          Order Placed Successfully!
        </h1>
        <p className="mt-2">
          Your order code is: {checkoutResponse.data.order_code}
        </p>
        <p className="mt-1">Payment Method: Cash on Delivery</p>
      </div>
    );
  }

  // 5. Default: Show the actual checkout form
  return (
    <div className="max-w-2xl mx-auto p-8">
      <h1 className="text-3xl font-bold mb-6">Checkout</h1>

      <div className="bg-gray-50 p-6 rounded-lg mb-6">
        <h2 className="text-xl font-semibold mb-4">Order Summary</h2>
        <p>Total: Rs. 50,000</p>
        <p className="text-sm text-gray-600 mt-2">Payment Method: eSewa</p>
      </div>

      <button
        onClick={handlePlaceOrder}
        disabled={isLoading}
        className="w-full bg-blue-600 text-white py-3 px-6 rounded-lg font-semibold hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition"
      >
        {isLoading ? "Processing..." : "Place Order & Pay with eSewa"}
      </button>
    </div>
  );
}
