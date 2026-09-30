"use client";

import { apiRequest } from "@/lib/api";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

const fields = [
  "title",
  "author_name",
  "genre",
  "category",
  "publication",
  "pages",
  "price",
  "units",
  "cover_page_url",
] as const;

export default function PublishBookPage() {
  const router = useRouter();
  const [form, setForm] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);
  const update = (field: string, value: string) =>
    setForm((current) => ({ ...current, [field]: value }));

  async function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    try {
      await apiRequest("book/publish-book", {
        method: "POST",
        body: JSON.stringify({
          ...form,
          pages: Number(form.pages),
          price: Number(form.price),
          units: Number(form.units),
        }),
      });
      toast.success("Book published");
      router.push("/books");
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Could not publish book",
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="mx-auto max-w-3xl px-4 py-10">
      <h1 className="text-3xl font-bold text-primary">Publish a book</h1>
      <p className="mt-2 text-muted">
        Create a listing for readers to discover.
      </p>
      <form
        onSubmit={submit}
        className="mt-8 grid gap-4 rounded-2xl border border-border bg-surface p-6 sm:grid-cols-2"
      >
        {fields.map((field) => (
          <input
            key={field}
            required={field !== "cover_page_url"}
            type={
              ["pages", "price", "units"].includes(field) ? "number" : "text"
            }
            placeholder={field.replaceAll("_", " ")}
            value={form[field] || ""}
            onChange={(e) => update(field, e.target.value)}
            className="h-12 rounded-lg border border-border px-4"
          />
        ))}
        <button
          disabled={loading}
          className="h-12 rounded-lg bg-primary font-semibold text-white sm:col-span-2"
        >
          {loading ? "Publishing..." : "Publish book"}
        </button>
      </form>
    </main>
  );
}
