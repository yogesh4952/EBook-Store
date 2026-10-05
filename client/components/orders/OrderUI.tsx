import type { ReactNode } from "react";
import { FaHandHoldingDollar } from "react-icons/fa6";
import {
  HiOutlineBanknotes,
  HiOutlineCheckCircle,
  HiOutlineClock,
  HiOutlineCreditCard,
  HiOutlineTruck,
  HiOutlineXCircle,
} from "react-icons/hi2";
import Image from "next/image";

import type {
  OrderStatus,
  PaymentMethod,
  PaymentStatus,
  UserOrder,
} from "@/lib/checkout";

export const formatPrice = (value: number) =>
  `Rs. ${new Intl.NumberFormat("en-NP", {
    maximumFractionDigits: 0,
  }).format(value)}`;

type Tone = "success" | "warning" | "error" | "neutral";

const TONE_CLASSES: Record<Tone, string> = {
  success: "bg-success/10 text-success",
  warning: "bg-accent/30 text-primary",
  error: "bg-error/10 text-error",
  neutral: "bg-background text-muted",
};

export const PaymentBadge = ({ method }: { method: PaymentMethod }) => (
  <span className="inline-flex items-center gap-1.5 text-xs text-muted">
    {method === "ESEWA" ? (
      <Image
        src="/esewa_logo.jpg"
        alt=""
        width={16}
        height={16}
        className="size-4 rounded border border-border object-cover"
      />
    ) : (
      <FaHandHoldingDollar size={14} />
    )}
    {method === "ESEWA" ? "eSewa" : "Cash on delivery"}
  </span>
);

const PAYMENT_STATUS: Record<
  PaymentStatus,
  { label: string; tone: Tone; Icon: typeof HiOutlineClock }
> = {
  PAID: { label: "Paid", tone: "success", Icon: HiOutlineCheckCircle },
  PENDING: { label: "Payment pending", tone: "warning", Icon: HiOutlineClock },
  REFUND: { label: "Refunded", tone: "neutral", Icon: HiOutlineBanknotes },
};

const ORDER_STATUS: Record<
  NonNullable<OrderStatus>,
  { label: string; tone: Tone; Icon: typeof HiOutlineClock }
> = {
  PLACED: { label: "Placed", tone: "warning", Icon: HiOutlineClock },
  DELIVERED: { label: "Delivered", tone: "success", Icon: HiOutlineTruck },
  CANCELLED: { label: "Cancelled", tone: "error", Icon: HiOutlineXCircle },
};

const Badge = ({
  tone,
  children,
}: {
  tone: Tone;
  children: ReactNode;
}) => (
  <span
    className={`inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-[11px] font-medium ${TONE_CLASSES[tone]}`}
  >
    {children}
  </span>
);

/** Server-sent status strings are upper case, but guard against drift. */
const normaliseStatus = <T extends string>(
  status: string | T,
  table: Record<T, { label: string; tone: Tone; Icon: typeof HiOutlineClock }>,
) => {
  const key = String(status).toUpperCase() as T;
  return table[key];
};

export const PaymentStatusBadge = ({ status }: { status: PaymentStatus }) => {
  const config = normaliseStatus(status, PAYMENT_STATUS);
  if (!config) return null;

  const { Icon } = config;
  return (
    <Badge tone={config.tone}>
      <Icon size={12} />
      {config.label}
    </Badge>
  );
};

export const OrderStatusBadge = ({
  status,
}: {
  status: OrderStatus | "";
}) => {
  // Orders created before the status column was populated come back as "".
  if (!status) return null;

  const config = normaliseStatus(status, ORDER_STATUS);
  if (!config) return null;

  const { Icon } = config;
  return (
    <Badge tone={config.tone}>
      <Icon size={12} />
      {config.label}
    </Badge>
  );
};

const formatOrderDate = (value: string) => {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return null;

  return date.toLocaleDateString("en-NP", {
    day: "numeric",
    month: "short",
    year: "numeric",
  });
};

/** Place → on the way → delivered, with the reached step highlighted. */
const STEPS = [
  { key: "PLACED", label: "Placed", Icon: HiOutlineCheckCircle },
  { key: "DELIVERED", label: "Delivered", Icon: HiOutlineTruck },
] as const;

const StepIndex: Record<NonNullable<OrderStatus>, number> = {
  PLACED: 0,
  DELIVERED: 1,
  CANCELLED: 0,
};

export const OrderTimeline = ({ order }: { order: UserOrder }) => {
  if (order.order_status === "CANCELLED") {
    return (
      <p className="flex items-center gap-1.5 text-xs text-error">
        <HiOutlineXCircle size={14} />
        This order was cancelled.
      </p>
    );
  }

  const current = StepIndex[order.order_status || "PLACED"] ?? 0;

  return (
    <ol className="flex items-center gap-2">
      {STEPS.map((step, index) => {
        const isDone = index <= current;
        const { Icon } = step;

        return (
          <li key={step.key} className="flex items-center gap-2">
            <span
              className={`flex items-center gap-1.5 rounded-full px-2.5 py-1 text-[11px] font-medium ${
                isDone ? "bg-success/10 text-success" : "bg-background text-muted"
              }`}
            >
              <Icon size={12} />
              {step.label}
            </span>

            {index < STEPS.length - 1 && (
              <span
                aria-hidden
                className={`h-px w-6 ${index < current ? "bg-success/40" : "bg-border"}`}
              />
            )}
          </li>
        );
      })}
    </ol>
  );
};

export const OrderSummaryBar = ({ order }: { order: UserOrder }) => {
  const date = formatOrderDate(order.created_at);
  const units = order.order_items.reduce(
    (sum, item) => sum + item.quantity,
    0,
  );

  return (
    <div className="flex flex-wrap items-center justify-between gap-4 border-b border-border px-5 py-4">
      <div className="min-w-0">
        <p className="font-mono text-sm font-semibold text-primary">
          {order.order_code}
        </p>
        <p className="mt-0.5 text-xs text-muted">
          {date && <span>{date} · </span>}
          {units} {units === 1 ? "item" : "items"}
        </p>
      </div>

      <div className="flex flex-wrap items-center gap-3">
        <OrderStatusBadge status={order.order_status} />
        <PaymentStatusBadge status={order.payment_status} />
        <PaymentBadge method={order.payment_method} />
        <span className="text-base font-bold text-primary tabular-nums">
          {formatPrice(order.total_price)}
        </span>
      </div>
    </div>
  );
};

export const OrderItems = ({ order }: { order: UserOrder }) => (
  <ul className="divide-y divide-border">
    {order.order_items.map((item) => (
      <li key={item.book_id} className="flex items-center gap-3 px-5 py-3">
        <div className="relative h-14 w-11 shrink-0 overflow-hidden rounded-md bg-accent/20">
          {item.book?.cover_page_url && (
            <Image
              src={item.book.cover_page_url}
              alt={`Cover of ${item.book.title}`}
              fill
              sizes="44px"
              className="object-cover"
            />
          )}
        </div>

        <div className="min-w-0 flex-1">
          <p className="truncate text-sm font-medium text-primary">
            {item.book?.title ?? `Book #${item.book_id}`}
          </p>
          {item.book?.author_name && (
            <p className="truncate text-xs text-muted">
              {item.book.author_name}
            </p>
          )}
        </div>

        <p className="shrink-0 text-xs text-muted tabular-nums">
          {item.quantity} × {formatPrice(item.unit_price)}
        </p>

        <p className="w-24 shrink-0 text-right text-sm font-semibold text-primary tabular-nums">
          {formatPrice(item.unit_price * item.quantity)}
        </p>
      </li>
    ))}
  </ul>
);

export const OrderAddress = ({ order }: { order: UserOrder }) => (
  <div className="flex items-start gap-2 border-t border-border px-5 py-4">
    <HiOutlineTruck size={16} className="mt-0.5 shrink-0 text-primary-light" />
    <p className="text-xs leading-relaxed text-muted">
      <span className="font-medium text-text">Delivering to </span>
      {order.shipping_city}, {order.shipping_delivery_address}
    </p>
  </div>
);

export const PaymentHint = ({ method }: { method: PaymentMethod }) =>
  method === "COD" ? (
    <p className="flex items-center gap-1.5 text-xs text-muted">
      <HiOutlineCreditCard size={14} />
      Keep the exact amount ready for the rider.
    </p>
  ) : (
    <p className="flex items-center gap-1.5 text-xs text-muted">
      <HiOutlineBanknotes size={14} />
      Paid online with eSewa.
    </p>
  );
