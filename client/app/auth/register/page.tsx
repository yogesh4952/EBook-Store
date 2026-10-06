"use client";

import Link from "next/link";
import { useState } from "react";
import { FaArrowLeft, FaArrowRight, FaGoogle } from "react-icons/fa";
import { FaCircleUser, FaStore } from "react-icons/fa6";
import { toast } from "sonner";

import OtpVerification from "@/components/OtpVerification";

type Role = "customer" | "seller";

const ROLES: {
  value: Role;
  label: string;
  description: string;
  Icon: typeof FaCircleUser;
}[] = [
  {
    value: "customer",
    label: "I'm buying books",
    description: "Browse, order, and track deliveries.",
    Icon: FaCircleUser,
  },
  {
    value: "seller",
    label: "I'm selling books",
    description: "List your books and manage sales.",
    Icon: FaStore,
  },
];

const Register = () => {
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [age, setAge] = useState("");
  const [email, setEmail] = useState("");
  const [phoneNumber, setPhoneNumber] = useState("");
  const [role, setRole] = useState<Role>("customer");

  const [isRegistering, setIsRegistering] = useState(false);
  const [pendingEmail, setPendingEmail] = useState<string | null>(null);

  /**
   * Registration is a two-step handshake: the account is created first, then the
   * OTP proves the address is reachable before a session is issued. So the
   * account is posted, and only on success do we ask for a code.
   */
  const handleGoogleRegister = () => {
    const clientId = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID;
    const redirectUri = `${process.env.NEXT_PUBLIC_FRONTEND_URL}/api/auth/google/callback`;
    const url = `https://accounts.google.com/o/oauth2/v2/auth?client_id=${clientId}&redirect_uri=${encodeURIComponent(redirectUri)}&response_type=code&scope=${encodeURIComponent("openid email profile")}&prompt=select_account`;
    window.location.href = url;
  };

  const handleRegister = async (event: React.FormEvent) => {
    event.preventDefault();

    if (isRegistering) return;

    if (!firstName.trim() || !lastName.trim()) {
      toast.error("Enter your first and last name");
      return;
    }

    if (!age || Number(age) < 13 || Number(age) > 120) {
      toast.error("Enter a valid age (13–120)");
      return;
    }

    if (!email.trim()) {
      toast.error("Enter your email address");
      return;
    }

    if (!phoneNumber.trim()) {
      toast.error("Enter your phone number");
      return;
    }

    setIsRegistering(true);

    try {
      const response = await fetch("/api/auth/register", {
        method: "POST",
        body: JSON.stringify({
          first_name: firstName.trim(),
          last_name: lastName.trim(),
          email: email.trim(),
          phone_number: phoneNumber.trim(),
          role,
          age: Number(age),
        }),
        headers: { "Content-Type": "application/json" },
      });

      const data = await response.json();

      if (!response.ok) {
        toast.error(data.message || "Could not create your account");
        return;
      }

      const responseOtp = await fetch("/api/auth/send-otp", {
        method: "POST",
        body: JSON.stringify({ email: email.trim() }),
        headers: { "Content-Type": "application/json" },
      });

      if (!responseOtp.ok) {
        toast.error(
          "Account created, but we could not send the OTP. Try logging in.",
        );
        return;
      }

      toast.success("Account created — check your inbox for the code");
      setPendingEmail(email.trim());
    } catch {
      toast.error("Internal server error");
    } finally {
      setIsRegistering(false);
    }
  };

  // Step 2: the account exists, so hand off to the shared OTP component.
  if (pendingEmail) {
    return (
      <OtpVerification
        email={pendingEmail}
        onBack={() => setPendingEmail(null)}
      />
    );
  }

  return (
    <div className="grid min-h-[750px] grid-cols-2 gap-8">
      {/* LEFT SIDE */}
      <div
        className="relative overflow-hidden rounded-3xl bg-cover bg-center"
        style={{ backgroundImage: "url('/loginbanner.jpg')" }}
      >
        <div className="absolute inset-0 bg-black/45" />

        <div className="relative z-10 flex h-full flex-col justify-between p-10 text-white">
          <div>
            <h1 className="text-5xl font-bold leading-tight">
              Join us
              <br />
              today!
            </h1>

            <p className="mt-4 text-lg text-gray-200">
              Create an account to order books, save your addresses, and track
              every delivery.
            </p>

            <div className="mt-12 space-y-6">
              <div className="flex items-center gap-4">
                <FaArrowRight className="text-xl text-amber-300" />
                <span>Thousands of titles from trusted sellers</span>
              </div>

              <div className="flex items-center gap-4">
                <FaArrowRight className="text-xl text-amber-300" />
                <span>Cash on delivery and eSewa checkout</span>
              </div>

              <div className="flex items-center gap-4">
                <FaArrowRight className="text-xl text-amber-300" />
                <span>Live order tracking from checkout to door</span>
              </div>
            </div>
          </div>

          <div className="rounded bg-white/80 px-4 py-2 text-black">
            <p className="text-5xl text-primary">&quot;</p>
            <p>A book is a device to ignite the imagination.</p>
            <p className="text-primary">- Alan Bennett</p>
          </div>
        </div>
      </div>

      {/* RIGHT SIDE */}
      <div className="flex flex-col justify-center">
        <h1 className="text-4xl font-bold">Create your account</h1>

        <p className="mt-3 text-gray-500">
          Already have one?{" "}
          <Link
            href="/auth/login"
            className="font-medium text-primary underline underline-offset-4"
          >
            Log in instead
          </Link>
        </p>

        <form onSubmit={handleRegister} className="mt-10 space-y-6">
          <div className="grid gap-6 sm:grid-cols-2">
            <div>
              <label htmlFor="first_name" className="mb-2 block font-medium">
                First name
              </label>
              <input
                id="first_name"
                name="first_name"
                value={firstName}
                onChange={(e) => setFirstName(e.target.value)}
                autoComplete="given-name"
                className="h-12 w-full rounded-lg border border-gray-300 px-4 outline-none transition focus:border-primary"
              />
            </div>

            <div>
              <label htmlFor="last_name" className="mb-2 block font-medium">
                Last name
              </label>
              <input
                id="last_name"
                name="last_name"
                value={lastName}
                onChange={(e) => setLastName(e.target.value)}
                autoComplete="family-name"
                className="h-12 w-full rounded-lg border border-gray-300 px-4 outline-none transition focus:border-primary"
              />
            </div>

            <div className="sm:col-span-2">
              <label htmlFor="age" className="mb-2 block font-medium">
                Age
              </label>
              <input
                id="age"
                name="age"
                type="number"
                min={13}
                max={120}
                value={age}
                onChange={(e) => setAge(e.target.value)}
                className="h-12 w-full rounded-lg border border-gray-300 px-4 outline-none transition focus:border-primary"
              />
            </div>
          </div>

          <div>
            <label htmlFor="email" className="mb-2 block font-medium">
              Email
            </label>
            <input
              id="email"
              name="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              autoComplete="email"
              placeholder="you@example.com"
              className="h-12 w-full rounded-lg border border-gray-300 px-4 outline-none transition focus:border-primary"
            />
          </div>

          <div>
            <label htmlFor="phone_number" className="mb-2 block font-medium">
              Phone number
            </label>
            <input
              id="phone_number"
              name="phone_number"
              type="tel"
              value={phoneNumber}
              onChange={(e) => setPhoneNumber(e.target.value)}
              placeholder="98xxxxxxxx"
              className="h-12 w-full rounded-lg border border-gray-300 px-4 outline-none transition focus:border-primary"
            />
          </div>

          <fieldset>
            <legend className="mb-2 block font-medium">
              How will you use Book-Store Nepal?
            </legend>

            <div className="grid gap-3 sm:grid-cols-2">
              {ROLES.map(({ value, label, description, Icon }) => {
                const isSelected = role === value;

                return (
                  <label
                    key={value}
                    className={`flex cursor-pointer items-start gap-3 rounded-xl border p-4 transition-colors ${
                      isSelected
                        ? "border-primary bg-accent/10"
                        : "border-gray-300 hover:border-accent"
                    }`}
                  >
                    <input
                      type="radio"
                      name="role"
                      value={value}
                      checked={isSelected}
                      onChange={() => setRole(value)}
                      className="sr-only"
                    />

                    <Icon
                      size={20}
                      className={isSelected ? "text-primary" : "text-gray-400"}
                    />

                    <span className="min-w-0">
                      <span className="block text-sm font-medium text-primary">
                        {label}
                      </span>
                      <span className="mt-0.5 block text-xs text-gray-500">
                        {description}
                      </span>
                    </span>
                  </label>
                );
              })}
            </div>
          </fieldset>

          <button
            type="submit"
            disabled={isRegistering}
            className="h-12 w-full rounded-lg bg-primary text-white transition hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {isRegistering ? "Creating account..." : "Create account"}
          </button>

          <button
            type="button"
            onClick={handleGoogleRegister}
            className="flex h-12 w-full items-center justify-center gap-3 rounded-lg border border-gray-300 bg-white shadow-sm transition hover:bg-gray-50"
          >
            <FaGoogle /> Continue with Google
          </button>

          <Link
            href="/auth/login"
            className="flex h-12 w-full items-center justify-center gap-2 rounded-lg border border-gray-300 bg-white shadow-sm transition hover:bg-gray-50"
          >
            <FaArrowLeft />
            Back to login
          </Link>
        </form>
      </div>
    </div>
  );
};

export default Register;
