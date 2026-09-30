import Link from "next/link";

interface PageProps {
  searchParams: Promise<{ orderId?: string }>;
}

export default async function PaymentSuccessPage({
  searchParams,
}: PageProps) {
  const { orderId } = await searchParams;

  return (
    <main className="flex min-h-screen items-center justify-center px-6">
      <div className="w-full max-w-md rounded-2xl bg-surface p-10 text-center shadow-lg">
        <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-success/10">
          <svg
            width="32"
            height="32"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            className="text-success"
          >
            <circle cx="12" cy="12" r="10" />
            <path d="M8 12l3 3 5-6" />
          </svg>
        </div>

        <h1 className="mt-6 text-3xl font-bold text-text">
          Thank you for your order!
        </h1>

        <p className="mt-3 leading-7 text-muted">
          Your payment was successful and your order is confirmed.
        </p>

        {orderId && (
          <div className="mt-6 rounded-lg bg-background px-4 py-3">
            <p className="text-sm text-muted">Order reference</p>
            <p className="mt-1 font-mono text-lg font-semibold text-text">
              {orderId}
            </p>
          </div>
        )}

        <div className="mt-8 space-y-3">
          <Link
            href="/orders"
            className="block h-12 w-full rounded-lg bg-primary text-white transition hover:opacity-90"
          >
            View my orders
          </Link>

          <Link
            href="/"
            className="block h-12 w-full rounded-lg border border-border bg-surface text-text transition hover:bg-background"
          >
            Continue shopping
          </Link>
        </div>
      </div>
    </main>
  );
}
