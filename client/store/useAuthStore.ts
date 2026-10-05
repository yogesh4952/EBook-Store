import { create } from "zustand";

export interface AuthUser {
  id?: number;
  email?: string;
  role?: string;
}

interface AuthState {
  user: AuthUser | null;
  loading: boolean;
  fetchUser: () => Promise<void>;
  logout: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  loading: true,

  fetchUser: async () => {
    set({ loading: true });
    try {
      const res = await fetch("/api/auth/me", { cache: "no-store" });

      if (!res.ok) {
        set({ user: null, loading: false });
        return;
      }

      const payload = await res.json();
      set({ user: payload.data ?? null, loading: false });
    } catch {
      set({ user: null, loading: false });
    }
  },

  logout: async () => {
    try {
      await fetch("/api/auth/logout", { method: "POST" });
    } finally {
      // Clear locally even if the request failed, so the UI never keeps
      // showing a signed-in state that the cookie no longer backs.
      set({ user: null });
      window.location.href = "/auth/login";
    }
  },
}));
