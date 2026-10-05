import React from "react";
import Image from "next/image";
import Link from "next/link";
import { FaBookOpen } from "react-icons/fa";
import {
  FaCreditCard,
  FaEnvelope,
  FaFacebookF,
  FaHandHoldingDollar,
  FaInstagram,
  FaLocationDot,
  FaXTwitter,
  FaYoutube,
} from "react-icons/fa6";

const SHOP_LINKS = [
  { label: "Browse All Books", href: "/book" },
  { label: "Popular & Trending", href: "/book?stock=in_stock" },
  { label: "New Arrivals", href: "/book" },
  { label: "Under Rs. 1,000", href: "/book?max=1000" },
];

const CATEGORY_LINKS = [
  { label: "Fiction", href: "/book?category=Fiction" },
  { label: "Technology", href: "/book?category=Technology" },
  { label: "Classic", href: "/book?category=Classic" },
  { label: "Business", href: "/book?category=Business" },
  { label: "Young Adult", href: "/book?category=Young%20Adult" },
  { label: "Biography", href: "/book?category=Biography" },
];

const GENRE_LINKS = [
  { label: "Programming", href: "/book?genre=Programming" },
  { label: "Fantasy", href: "/book?genre=Fantasy" },
  { label: "Romance", href: "/book?genre=Romance" },
  { label: "Finance", href: "/book?genre=Finance" },
  { label: "Historical Fiction", href: "/book?genre=Historical%20Fiction" },
  { label: "Mystery", href: "/book?genre=Mystery" },
];

const HELP_LINKS = [
  { label: "How to Order", href: "/checkout" },
  { label: "Payment Options", href: "/checkout" },
  { label: "Delivery Info", href: "/checkout" },
  { label: "My Orders", href: "/orders" },
];

const SOCIALS = [
  { label: "Facebook", href: "https://facebook.com", Icon: FaFacebookF },
  { label: "Instagram", href: "https://instagram.com", Icon: FaInstagram },
  { label: "X", href: "https://x.com", Icon: FaXTwitter },
  { label: "YouTube", href: "https://youtube.com", Icon: FaYoutube },
];

const FOOTER_LINKS = [
  { label: "About Us", href: "/" },
  { label: "Privacy Policy", href: "/" },
  { label: "Terms of Service", href: "/" },
  { label: "Contact", href: "/book" },
];

const Footer = () => {
  return (
    <footer className="mt-auto border-t border-border bg-surface">
      {/* Newsletter */}
      <div className="border-b border-border bg-primary">
        <div className="mx-auto flex max-w-7xl flex-col items-center gap-6 px-6 py-10 lg:flex-row lg:justify-between">
          <div className="max-w-md text-center lg:text-left">
            <h2 className="text-xl font-semibold text-white">
              Get new arrivals in your inbox
            </h2>
            <p className="mt-1 text-sm text-accent">
              Occasional reading lists and seller drops. No spam, unsubscribe
              anytime.
            </p>
          </div>

          <div className="w-full max-w-md">
            <fieldset disabled className="flex flex-col gap-2 sm:flex-row">
              <legend className="sr-only">Newsletter signup</legend>
              <input
                type="email"
                placeholder="you@example.com"
                className="h-11 flex-1 rounded-lg border border-transparent bg-surface px-4 text-sm text-text outline-none placeholder:text-muted disabled:cursor-not-allowed disabled:bg-surface/60 disabled:text-muted"
              />
              <button
                type="submit"
                className="h-11 shrink-0 cursor-not-allowed rounded-lg bg-accent px-6 text-sm font-semibold text-primary opacity-50"
              >
                Subscribe
              </button>
            </fieldset>
            <p className="mt-2 text-xs text-accent/80">
              Newsletter signup is coming soon.
            </p>
          </div>
        </div>
      </div>

      {/* Link columns */}
      <div className="mx-auto grid max-w-7xl gap-10 px-6 py-12 sm:grid-cols-2 lg:grid-cols-6 lg:items-start">
        {/* Brand */}
        <div className="lg:col-span-2">
          <div className="flex items-center gap-3">
            <FaBookOpen size={26} className="text-primary" />
            <span className="text-lg font-bold text-primary">
              Book-Store Nepal
            </span>
          </div>
          <p className="mt-3 max-w-sm text-sm leading-relaxed text-muted">
            Nepal&apos;s marketplace for books — thousands of titles from
            trusted sellers, delivered nationwide with cash on delivery and
            eSewa.
          </p>

          <ul className="mt-6 space-y-2.5 text-sm text-muted">
            <li className="flex items-start gap-2.5">
              <FaLocationDot
                size={16}
                className="mt-0.5 shrink-0 text-primary-light"
              />
              <span>New Road, Kathmandu 44600, Nepal</span>
            </li>
            <li className="flex items-start gap-2.5">
              <FaEnvelope
                size={16}
                className="mt-0.5 shrink-0 text-primary-light"
              />
              <a
                href="mailto:hello@bookstorenepal.com"
                className="transition-colors hover:text-primary"
              >
                hello@bookstorenepal.com
              </a>
            </li>
            <li className="flex items-start gap-2.5">
              <FaCreditCard
                size={16}
                className="mt-0.5 shrink-0 text-primary-light"
              />
              <span>
                Cash on delivery &amp; eSewa · Free delivery over Rs. 2,000
              </span>
            </li>
          </ul>

          <div className="mt-6 flex items-center gap-2">
            {SOCIALS.map(({ label, href, Icon }) => (
              <a
                key={label}
                href={href}
                target="_blank"
                rel="noopener noreferrer"
                aria-label={label}
                className="flex size-9 items-center justify-center rounded-full border border-border bg-background text-primary transition-all hover:-translate-y-0.5 hover:border-accent hover:bg-primary hover:text-white"
              >
                <Icon size={16} />
              </a>
            ))}
          </div>
        </div>

        {/* Link lists */}
        <nav aria-labelledby="footer-shop">
          <h3
            id="footer-shop"
            className="text-sm font-semibold tracking-wide text-primary uppercase"
          >
            Shop
          </h3>
          <ul className="mt-4 space-y-2.5 text-sm">
            {SHOP_LINKS.map((link) => (
              <li key={link.label}>
                <Link
                  href={link.href}
                  className="text-muted transition-colors hover:text-primary"
                >
                  {link.label}
                </Link>
              </li>
            ))}
          </ul>
        </nav>

        <nav aria-labelledby="footer-categories">
          <h3
            id="footer-categories"
            className="text-sm font-semibold tracking-wide text-primary uppercase"
          >
            Categories
          </h3>
          <ul className="mt-4 space-y-2.5 text-sm">
            {CATEGORY_LINKS.map((link) => (
              <li key={link.label}>
                <Link
                  href={link.href}
                  className="text-muted transition-colors hover:text-primary"
                >
                  {link.label}
                </Link>
              </li>
            ))}
          </ul>
        </nav>

        <nav aria-labelledby="footer-help">
          <h3
            id="footer-help"
            className="text-sm font-semibold tracking-wide text-primary uppercase"
          >
            Help
          </h3>
          <ul className="mt-4 space-y-2.5 text-sm">
            {HELP_LINKS.map((link) => (
              <li key={link.label}>
                <Link
                  href={link.href}
                  className="text-muted transition-colors hover:text-primary"
                >
                  {link.label}
                </Link>
              </li>
            ))}
          </ul>
        </nav>

        <nav aria-labelledby="footer-genres">
          <h3
            id="footer-genres"
            className="text-sm font-semibold tracking-wide text-primary uppercase"
          >
            Popular Genres
          </h3>
          <ul className="mt-4 space-y-2.5 text-sm">
            {GENRE_LINKS.map((link) => (
              <li key={link.label}>
                <Link
                  href={link.href}
                  className="text-muted transition-colors hover:text-primary"
                >
                  {link.label}
                </Link>
              </li>
            ))}
          </ul>
        </nav>
      </div>

      {/* Bottom bar */}
      <div className="border-t border-border bg-background">
        <div className="mx-auto flex max-w-7xl flex-col items-center gap-4 px-6 py-6 lg:flex-row lg:justify-between">
          <p className="text-xs text-muted">
            © {new Date().getFullYear()} Book-Store Nepal. All rights
            reserved.
          </p>

          <ul className="flex flex-wrap items-center justify-center gap-x-5 gap-y-2 text-xs">
            {FOOTER_LINKS.map((link) => (
              <li key={link.label}>
                <Link
                  href={link.href}
                  className="text-muted transition-colors hover:text-primary"
                >
                  {link.label}
                </Link>
              </li>
            ))}
          </ul>

          <div className="flex flex-wrap items-center justify-center gap-2">
            <span className="text-xs text-muted">We accept</span>
            <span className="flex items-center gap-1.5 rounded-md border border-border bg-surface px-2 py-1 text-[11px] font-medium text-muted">
              <Image
                src="/esewa_logo.jpg"
                alt=""
                width={14}
                height={14}
                className="rounded"
              />
              eSewa
            </span>
            <span className="flex items-center gap-1.5 rounded-md border border-border bg-surface px-2 py-1 text-[11px] font-medium text-muted">
              <FaHandHoldingDollar size={13} />
              Cash on Delivery
            </span>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;