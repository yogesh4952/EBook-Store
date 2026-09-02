"use client";

import { useState } from "react";

const Genre = ["fiction", "non-fiction", "romance", "science"];

const FilterBar = () => {
  const [selectedGenre, setSelectedGenre] = useState("");
  return (
    <div className="sticky top-4 h-fit p-4 rounded shadow-lg">
      <div>
        <label htmlFor="price">Price</label>
        <input type="range" min={200} max={5000} name="" id="" />
      </div>

      {/* Genre */}
      <div>
        {Genre.map((g, i) => (
          <div key={i}>
            <label htmlFor={g}>{g}</label>
            <input
              type="radio"
              value={g}
              name="genre"
              checked={selectedGenre === g}
              onChange={(e) => setSelectedGenre(e.target.value)}
              id={i}
            />
          </div>
        ))}
      </div>
    </div>
  );
};

export default FilterBar;
