"use client";

import { useEffect, useState } from "react";

import type { EsewaPayload } from "@/lib/checkout";

const ESEWA_FORM_ACTION =
  "https://rc-epay.esewa.com.np/api/epay/main/v2/form";

interface EsewaRedirectProps {
  payload: EsewaPayload;
  orderCode: string;
}

export default function EsewaRedirect({
  payload,
  orderCode,
}: EsewaRedirectProps) {
  const [isRedirecting, setIsRedirecting] = useState(true);

  useEffect(() => {
    const form = document.createElement("form");
    form.method = "POST";
    form.action = ESEWA_FORM_ACTION;

    Object.entries(payload).forEach(([key, value]) => {
      const input = document.createElement("input");
      input.type = "hidden";
      input.name = key;
      input.value = value;
      form.appendChild(input);
    });

    document.body.appendChild(form);
    form.submit();
  }, [payload]);

  if (isRedirecting) {
    return (
      <main className="flex min-h-screen flex-col items-center justify-center gap-3 px-6">
        <span className="size-10 animate-spin rounded-full border-2 border-border border-t-primary" />
        <h1 className="text-lg font-semibold text-primary">
          Redirecting to eSewa…
        </h1>
        <p className="text-sm text-muted">
          Taking you to the secure payment page.
        </p>
        <button
          type="button"
          onClick={() => setIsRedirecting(false)}
          className="mt-2 text-sm text-muted underline underline-offset-4 hover:text-primary"
        >
          Continue manually
        </button>
      </main>
    );
  }

  return (
    <main className="flex min-h-screen flex-col items-center justify-center gap-3 px-6 text-center">
      <h1 className="text-lg font-semibold text-primary">
        Automatic redirect did not start
      </h1>
      <p className="text-sm text-muted">
        Order {orderCode} is saved. If your browser blocked the redirect, open the
        eSewa payment page manually.
      </p>
      <form method="POST" action={ESEWA_FORM_ACTION}>
        {Object.entries(payload).map(([key, value]) => (
          <input key={key} type="hidden" name={key} value={value} />
        ))}
        <button
          type="submit"
          className="mt-2 rounded-lg bg-primary px-4 py-3 text-sm font-medium text-white transition-colors hover:bg-primary-light"
        >
          Go to eSewa
        </button>
      </form>
    </main>
  );
}