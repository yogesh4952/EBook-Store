import React from "react";
import { CiHeart } from "react-icons/ci";
import { FaBookOpen, FaSearch } from "react-icons/fa";
import CartBadge from "./CartBadge";
import ProfileMenu from "./ProfileMenu";
import { getCategoriesFromApi } from "@/lib/book-api";

const Navbar = async () => {
  const categories = await getCategoriesFromApi();

  return (
    <nav className="w-full min-h-24 grid grid-cols-[1fr_2fr_1fr] items-center gap-4 px-6">
      {/* Logo */}
      <div className="flex items-center gap-4">
        <FaBookOpen size={30} />
        <h1 className="text-primary text-xl font-bold">Book-Store Nepal</h1>
      </div>

      {/* Search Bar */}
      <form
        action="/book"
        method="get"
        role="search"
        className="flex h-11 border border-accent rounded-lg overflow-hidden"
      >
        {/* Input */}
        <input
          type="search"
          name="q"
          id="search_book"
          placeholder="Search by title, author..."
          className="min-w-0 flex-1 bg-surface px-4 outline-none placeholder:text-muted"
        />

        {/* Category */}
        <select
          name="category"
          id="categories"
          aria-label="Filter by category"
          className="max-w-40 border-l border-accent bg-surface px-3 outline-none text-sm"
        >
          <option value="">All Categories</option>
          {categories.map((category, i) => (
            <option key={i} value={category}>
              {category}
            </option>
          ))}
        </select>

        {/* Search Button */}
        <button
          type="submit"
          aria-label="Search books"
          className="w-12 flex items-center justify-center bg-primary text-white hover:opacity-90"
        >
          <FaSearch size={18} />
        </button>
      </form>

      {/* Icons */}
      <div className="flex justify-end gap-5 items-center">
        <CiHeart
          className="cursor-pointer hover:scale-105 transition-all"
          size={32}
        />
        <CartBadge />
        <ProfileMenu />
      </div>
    </nav>
  );
};

export default Navbar;
