import React from "react";

const Banner = () => {
  return (
    <div className="px-10 h-[80vh] w-full bg-[url('/Banner.png')]  bg-cover bg-center  relative overflow-hidden rounded-md text-primary ">
      <div className="flex flex-col justify-center  h-full w-full space-y-2">
        <h1 className="text-4xl tracking-wide">Find Your Next</h1>
        <h1 className="text-5xl tracking-wider font-bold">Great Read</h1>
        <div>
          <p>Thousand of books Trusted sellers.</p>
          <p>Fast delivery across Nepal.</p>
        </div>
        <button className=" mt-2 px-4 py-3 text-xl hover:bg-amber-950 cursor-pointer rounded-xl bg-primary text-white w-fit">
          Explore Books
        </button>
      </div>
    </div>
  );
};

export default Banner;
