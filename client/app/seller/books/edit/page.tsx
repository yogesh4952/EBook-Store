"use client";

import { apiRequest } from "@/lib/api";
import { useState } from "react";
import { toast } from "sonner";

export default function EditBookPage() {
  const [bookId, setBookId] = useState("");
  const [sellerId, setSellerId] = useState("");
  const [title, setTitle] = useState("");
  const [price, setPrice] = useState("");
  const [units, setUnits] = useState("");
  const [loading, setLoading] = useState(false);

  async function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    try {
      await apiRequest("book/update-book", {
        method: "PATCH",
        body: JSON.stringify({
          book_id: Number(bookId),
          seller_id: Number(sellerId),
          ...(title && { title }),
          ...(price && { price: Number(price) }),
          ...(units && { units: Number(units) }),
        }),
      });
      toast.success("Book updated");
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Could not update book",
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="mx-auto max-w-2xl px-4 py-10">
      <h1 className="text-3xl font-bold text-primary">Update a book</h1>
      <form
        onSubmit={submit}
        className="mt-8 space-y-4 rounded-2xl border border-border bg-surface p-6"
      >
        <input
          required
          type="number"
          placeholder="Book ID"
          value={bookId}
          onChange={(e) => setBookId(e.target.value)}
          className="h-12 w-full rounded-lg border border-border px-4"
        />
        <input
          required
          type="number"
          placeholder="Seller ID"
          value={sellerId}
          onChange={(e) => setSellerId(e.target.value)}
          className="h-12 w-full rounded-lg border border-border px-4"
        />
        <input
          placeholder="New title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className="h-12 w-full rounded-lg border border-border px-4"
        />
        <div className="grid gap-4 sm:grid-cols-2">
          <input
            type="number"
            placeholder="New price"
            value={price}
            onChange={(e) => setPrice(e.target.value)}
            className="h-12 rounded-lg border border-border px-4"
          />
          <input
            type="number"
            placeholder="New units"
            value={units}
            onChange={(e) => setUnits(e.target.value)}
            className="h-12 rounded-lg border border-border px-4"
          />
        </div>
        <button
          disabled={loading}
          className="h-12 w-full rounded-lg bg-primary font-semibold text-white"
        >
          {loading ? "Updating..." : "Update book"}
        </button>
      </form>
    </main>
  );
}
