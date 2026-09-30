"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

export default function RegisterPage() {
  const router = useRouter();
  const [form, setForm] = useState({
    first_name: "",
    last_name: "",
    email: "",
    role: "customer",
    age: "18",
  });
  const [loading, setLoading] = useState(false);

  const update = (field: string, value: string) =>
    setForm((current) => ({ ...current, [field]: value }));

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    try {
      const response = await fetch("/api/auth/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...form, age: Number(form.age) }),
      });
      const data = await response.json();
      if (!response.ok) throw new Error(data.message || "Registration failed");
      toast.success(data.message || "Account created");
      router.push("/auth/login");
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Registration failed",
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="mx-auto flex min-h-[70vh] max-w-xl items-center justify-center py-12">
      <form
        onSubmit={handleSubmit}
        className="w-full space-y-5 rounded-2xl border border-border bg-surface p-8 shadow-sm"
      >
        <div>
          <p className="text-sm font-semibold uppercase tracking-[0.18em] text-primary-light">
            Book Store Nepal
          </p>
          <h1 className="mt-2 text-3xl font-bold text-primary">
            Create your account
          </h1>
        </div>
        <div className="grid gap-4 sm:grid-cols-2">
          <input
            required
            placeholder="First name"
            value={form.first_name}
            onChange={(e) => update("first_name", e.target.value)}
            className="h-12 rounded-lg border border-border px-4"
          />
          <input
            required
            placeholder="Last name"
            value={form.last_name}
            onChange={(e) => update("last_name", e.target.value)}
            className="h-12 rounded-lg border border-border px-4"
          />
        </div>
        <input
          required
          type="email"
          placeholder="Email"
          value={form.email}
          onChange={(e) => update("email", e.target.value)}
          className="h-12 w-full rounded-lg border border-border px-4"
        />
        <div className="grid gap-4 sm:grid-cols-2">
          <input
            required
            min="13"
            type="number"
            placeholder="Age"
            value={form.age}
            onChange={(e) => update("age", e.target.value)}
            className="h-12 rounded-lg border border-border px-4"
          />
          <select
            value={form.role}
            onChange={(e) => update("role", e.target.value)}
            className="h-12 rounded-lg border border-border bg-white px-4"
          >
            <option value="customer">Customer</option>
            <option value="seller">Seller</option>
          </select>
        </div>
        <button
          disabled={loading}
          className="h-12 w-full rounded-lg bg-primary font-semibold text-white disabled:opacity-60"
        >
          {loading ? "Creating account..." : "Create account"}
        </button>
        <p className="text-center text-sm text-muted">
          Already registered?{" "}
          <a className="font-semibold text-primary" href="/auth/login">
            Sign in
          </a>
        </p>
      </form>
    </main>
  );
}
