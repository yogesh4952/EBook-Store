import React from "react";
import { AiOutlineShoppingCart } from "react-icons/ai";
import { CiHeart } from "react-icons/ci";
import { FaBookOpen, FaSearch } from "react-icons/fa";
import { MdOutlinePerson3 } from "react-icons/md";

const Navbar = () => {
  return (
    <nav className="w-full min-h-24 grid grid-cols-[1fr_2fr_1fr] items-center gap-4 px-6">
      {/* Logo */}
      <div className="flex items-center gap-4">
        <FaBookOpen size={30} />
        <h1 className="text-primary text-xl font-bold">Book-Store Nepal</h1>
      </div>

      {/* Search Bar */}
      <div className="flex h-11 border border-accent rounded-lg overflow-hidden">
        {/* Input */}
        <input
          type="text"
          name="search_book"
          id="search_book"
          placeholder="Search by title, author..."
          className="flex-1 px-4 outline-none"
        />

        {/* Category */}
        <select
          name="categories"
          id="categories"
          className="w-40 px-3 border-l border-accent bg-white outline-none text-sm"
        >
          <option value="">All Categories</option>
          <option value="romance">Romance</option>
          <option value="sci-fi">Sci-Fi</option>
          <option value="manhwa">Manhwa</option>
          <option value="manga">Manga</option>
        </select>

        {/* Search Button */}
        <button className="w-12 flex items-center justify-center bg-primary text-white hover:opacity-90">
          <FaSearch size={18} />
        </button>
      </div>

      {/* Icons */}
      <div className="flex justify-end gap-5 items-center">
        <CiHeart
          className="cursor-pointer hover:scale-105 transition-all"
          size={32}
        />
        <AiOutlineShoppingCart
          className="cursor-pointer hover:scale-105 transition-all"
          size={32}
        />
        <MdOutlinePerson3
          className="cursor-pointer hover:scale-105 transition-all"
          size={32}
        />
      </div>
    </nav>
  );
};

export default Navbar;
