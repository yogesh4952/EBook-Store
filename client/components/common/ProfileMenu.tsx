"use client";

import Link from "next/link";
import { useEffect, useRef, useState } from "react";
import {
  HiOutlineArrowRightStartOnRectangle,
  HiOutlineChevronDown,
  HiOutlineClipboardDocumentList,
  HiOutlineMapPin,
} from "react-icons/hi2";

import { useAuthStore } from "@/store/useAuthStore";

const NAV_LINKS = [
  {
    label: "My orders",
    href: "/orders",
    Icon: HiOutlineClipboardDocumentList,
  },
  { label: "Delivery addresses", href: "/addresses", Icon: HiOutlineMapPin },
];

const ROLE_LABELS: Record<string, string> = {
  customer: "Customer",
  seller: "Seller",
  admin: "Admin",
};

/** First letter of the email local part, so the avatar is not a blank circle. */
const avatarInitial = (email?: string) =>
  email?.trim().charAt(0).toUpperCase() || "U";

const ProfileMenu = () => {
  const { user, fetchUser, logout } = useAuthStore();
  const [isOpen, setIsOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    fetchUser();
  }, [fetchUser]);

  useEffect(() => {
    if (!isOpen) return;

    const handlePointerDown = (event: MouseEvent) => {
      if (!containerRef.current?.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") setIsOpen(false);
    };

    document.addEventListener("mousedown", handlePointerDown);
    document.addEventListener("keydown", handleKeyDown);

    return () => {
      document.removeEventListener("mousedown", handlePointerDown);
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [isOpen]);

  const handleLogout = async () => {
    // logout() clears the cookie and hard-navigates away, so there is no state
    // left to update afterwards.
    await logout();
  };

  const handleNavigate = () => setIsOpen(false);

  return (
    <div ref={containerRef} className="relative">
      <button
        type="button"
        onClick={() => setIsOpen((prev) => !prev)}
        aria-haspopup="menu"
        aria-expanded={isOpen}
        aria-label="Account menu"
        className="flex cursor-pointer items-center gap-1.5 rounded-full transition-all hover:scale-105 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
      >
        <span className="flex h-9 w-9 items-center justify-center rounded-full bg-primary text-sm font-semibold text-white">
          {avatarInitial(user?.email)}
        </span>
        <HiOutlineChevronDown
          size={16}
          className={`text-muted transition-transform ${
            isOpen ? "rotate-180" : ""
          }`}
        />
      </button>

      {isOpen && (
        <div
          role="menu"
          className="absolute right-0 top-full z-50 mt-3 w-64 origin-top-right overflow-hidden rounded-2xl border border-border bg-surface shadow-xl"
        >
          {/* Identity */}
          <div className="border-b border-border bg-background px-4 py-4">
            <p className="truncate text-sm font-semibold text-primary">
              {user?.email ?? "Signed in"}
            </p>
            {user?.role && (
              <span className="mt-1.5 inline-block rounded-full bg-accent/30 px-2.5 py-0.5 text-[11px] font-medium text-primary">
                {ROLE_LABELS[user.role] ?? user.role}
              </span>
            )}
          </div>

          {/* Navigation */}
          <ul className="py-1.5">
            {NAV_LINKS.map(({ label, href, Icon }) => (
              <li key={href}>
                <Link
                  href={href}
                  role="menuitem"
                  onClick={handleNavigate}
                  className="flex items-center gap-3 px-4 py-2.5 text-sm text-text transition-colors hover:bg-accent/10"
                >
                  <Icon size={16} className="shrink-0 text-muted" />
                  {label}
                </Link>
              </li>
            ))}
          </ul>

          {/* Logout */}
          <div className="border-t border-border py-1.5">
            <button
              type="button"
              role="menuitem"
              onClick={handleLogout}
              className="flex w-full cursor-pointer items-center gap-3 px-4 py-2.5 text-sm text-error transition-colors hover:bg-error/10"
            >
              <HiOutlineArrowRightStartOnRectangle size={16} />
              Log out
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default ProfileMenu;
