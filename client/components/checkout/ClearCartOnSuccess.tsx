"use client";

import { useEffect, useRef } from "react";

import { clearCart } from "@/lib/checkout";

/**
 * Empties the cart once the eSewa payment has been confirmed by the backend.
 *
 * Clearing here rather than at order placement is deliberate: an eSewa order
 * stays PENDING until the provider calls back, and the customer may abandon the
 * payment. Wiping the cart before then would lose the order.
 *
 * The ref guards against React 19 StrictMode double-invoking the effect, which
 * would otherwise fire two requests.
 */
const ClearCartOnSuccess = () => {
  const hasCleared = useRef(false);

  useEffect(() => {
    if (hasCleared.current) return;
    hasCleared.current = true;

    clearCart().catch((error) => {
      // The order is already confirmed and paid, so a failure here is not worth
      // interrupting the user over.
      console.error("Failed to clear cart after payment", error);
    });
  }, []);

  return null;
};

export default ClearCartOnSuccess;