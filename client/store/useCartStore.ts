import { create } from "zustand";

import { getCart } from "@/lib/checkout";

interface CartState {
  count: number;
  refresh: () => void;
  loading: boolean;
}

export const useCartStore = create<CartState>((set) => ({
  count: 0,
  loading: false,
  refresh: () => {
    set({ loading: true });
    getCart()
      .then((cart) => {
        const total = cart.items.reduce(
          (sum, item) => sum + item.quantity,
          0,
        );
        set({ count: total });
      })
      .catch(() => set({ count: 0 }))
      .finally(() => set({ loading: false }));
  },
}));

export const refreshCart = () => useCartStore.getState().refresh();
