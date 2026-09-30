import Link from "next/link";

const FAILURE_MESSAGES: Record<string, string> = {
  invalid_signature: "The payment response could not be verified. Please try again.",
  order_not_found: "We could not find the order for this payment.",
  order_not_payable: "This order was not eligible for online payment.",
  amount_mismatch: "The amount paid did not match the order total.",
  payment_incomplete: "The payment was not completed.",
  verification_unavailable: "We could not reach the payment provider. Please try again.",
  verification_failed: "The payment could not be verified. Please try again.",
};

interface PageProps {
  searchParams: Promise<{ reason?: string }>;
}

export default async function PaymentFailurePage({
  searchParams,
}: PageProps) {
  const { reason } = await searchParams;
  const message =
    (reason && FAILURE_MESSAGES[reason]) ||
    "The payment could not be completed.";

  return (
    <main className="flex min-h-screen items-center justify-center px-6">
      <div className="w-full max-w-md rounded-2xl bg-surface p-10 text-center shadow-lg">
        <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-error/10">
          <svg
            width="32"
            height="32"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            className="text-error"
          >
            <circle cx="12" cy="12" r="10" />
            <path d="M15 9l-6 6M9 9l6 6" />
          </svg>
        </div>

        <h1 className="mt-6 text-3xl font-bold text-text">Payment failed</h1>

        <p className="mt-3 leading-7 text-muted">{message}</p>

        <p className="mt-4 text-sm text-muted">
          You have not been charged. Your order is still saved, so you can try
          paying again.
        </p>

        <div className="mt-8 space-y-3">
          <Link
            href="/checkout"
            className="block h-12 w-full rounded-lg bg-primary text-white transition hover:opacity-90"
          >
            Try again
          </Link>

          <Link
            href="/"
            className="block h-12 w-full rounded-lg border border-border bg-surface text-text transition hover:bg-background"
          >
            Back to home
          </Link>
        </div>
      </div>
    </main>
  );
}
