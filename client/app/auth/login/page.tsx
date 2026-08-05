"use client";
import { useState } from "react";
import {
  FaBookmark,
  FaFacebook,
  FaGoogle,
  FaHeart,
  FaJediOrder,
} from "react-icons/fa";
import { FaPeopleGroup, FaShield } from "react-icons/fa6";
import { toast } from "sonner";

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = async () => {
    try {
      const input = { email, password };
      const result = await fetch("/api/login", {
        method: "POST",
        body: JSON.stringify(input),
        headers: {
          "Content-type": "application/json",
        },
      });

      const data = await result.json();
      if (!result.ok) {
        toast.error(data.message || result.statusText || "Error while login");
        return;
      }

      toast.success("Login successful");
    } catch (error) {
      toast.error("Internal server error");
      console.error(error);
    }
  };
  return (
    <div className="grid min-h-[750px] grid-cols-2 gap-8">
      {/* LEFT SIDE */}
      <div
        className="relative overflow-hidden rounded-3xl bg-cover bg-center"
        style={{ backgroundImage: "url('/login.jpg')" }}
      >
        {/* Dark Overlay */}
        <div className="absolute inset-0 bg-black/45" />

        {/* Content */}
        <div className="relative z-10 flex h-full flex-col justify-between p-10 text-white">
          {/* Top */}
          <div>
            <h1 className="text-5xl font-bold leading-tight">
              Welcome
              <br />
              Back!
            </h1>

            <p className="mt-4 text-lg text-gray-200">
              Login to continue your reading journey and discover your next
              favourite book.
            </p>

            {/* Features */}
            <div className="mt-12 space-y-6">
              <div className="flex items-center gap-4">
                <FaBookmark className="text-xl text-amber-300" />
                <span>Save your favourite books</span>
              </div>

              <div className="flex items-center gap-4">
                <FaJediOrder className="text-xl text-amber-300" />
                <span>Track your orders easily</span>
              </div>

              <div className="flex items-center gap-4">
                <FaHeart className="text-xl text-amber-300" />
                <span>Get personalized recommendations</span>
              </div>

              <div className="flex items-center gap-4 ">
                <FaPeopleGroup className="text-xl text-amber-300" />
                <span>Join thousands of passionate readers</span>
              </div>
            </div>
          </div>

          {/* Bottom */}
          <div className="bg-white/80 text-black rounded px-4 py-2 ">
            <p className="text-5xl text-primary">"</p>
            <p>A room without books is like a body without soul.</p>
            <p className="text-primary">- Marcus Tullius Cicero</p>
          </div>
        </div>
      </div>

      {/* RIGHT SIDE */}
      <div className="flex flex-col justify-center">
        <h1 className="text-4xl font-bold">Login to your account</h1>

        <p className="mt-3 text-gray-500">
          Enter your credentials to access your account.
        </p>

        {/* FORM */}
        <div className="mt-10 space-y-6">
          <div>
            <label className="mb-2 block font-medium">Email</label>

            <input
              type="email"
              value={email}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setEmail(e.target.value)
              }
              className="h-12 w-full rounded-lg border border-gray-300 px-4 outline-none transition focus:border-primary"
            />
          </div>

          <div>
            <label className="mb-2 block font-medium">Password</label>

            <input
              type="password"
              value={password}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setPassword(e.target.value)
              }
              className="h-12 w-full rounded-lg border border-gray-300 px-4 outline-none transition focus:border-primary"
            />
          </div>
        </div>

        <p className="mt-4 cursor-pointer text-right font-semibold text-primary hover:underline">
          Forgot Password?
        </p>

        <button
          className="mt-8 h-12 w-full rounded-lg bg-primary text-white transition hover:opacity-90"
          onClick={() => handleLogin()}
        >
          Login
        </button>

        <p className="my-8 text-center text-gray-500">or continue with</p>

        <div className="space-y-4">
          <button className="flex h-12 w-full items-center justify-center gap-3 rounded-lg border border-gray-300 bg-white shadow-sm transition hover:cursor-pointer hover:bg-gray-50">
            <FaGoogle />
            Continue with Google
          </button>

          <button className="flex h-12 w-full items-center justify-center gap-3 rounded-lg border border-gray-300 bg-white shadow-sm transition hover:cursor-pointer hover:bg-gray-50">
            <FaFacebook />
            Continue with Facebook
          </button>
        </div>

        <p className="mt-10 flex items-center justify-center gap-2 text-sm text-gray-500">
          <FaShield />
          We never share your data with anyone.
        </p>
      </div>
    </div>
  );
};

export default Login;
