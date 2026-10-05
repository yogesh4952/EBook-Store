"use client";

import React, { useEffect, useState } from "react";
import { AiOutlineShoppingCart } from "react-icons/ai";
import { useCartStore } from "@/store/useCartStore";
import Sidebar from "../cart/Sidebar";

const CartBadge = () => {
  const { count, refresh } = useCartStore();
  const [openSidebar, setOpenSidebar] = useState(false);

  useEffect(() => {
    refresh();
  }, [refresh]);

  return (
    <>
      <div
        className="relative cursor-pointer hover:scale-105 transition-all"
        onClick={() => setOpenSidebar(true)}
        role="button"
        aria-label="Open cart"
      >
        <AiOutlineShoppingCart size={32} />
        <div className="absolute -top-2 -right-2 flex h-5 w-5 items-center justify-center rounded-full bg-error text-[11px] font-bold text-white shadow-sm">
          {count}
        </div>
      </div>
      {openSidebar && <Sidebar onClose={() => setOpenSidebar(false)} />}
    </>
  );
};

export default CartBadge;