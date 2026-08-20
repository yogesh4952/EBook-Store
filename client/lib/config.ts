/**
 * Static frontend configuration.
 *
 * All values below are hardcoded so the app works without any
 * environment variables / dynamic configuration. Replace the
 * placeholder values with your real values when you are ready.
 */

/** Base URL of the Go REST API backend. */
export const BACKEND_URL = "http://localhost:8080";

/** Name of the cookie used to store the JWT access token. */
export const AUTH_COOKIE_NAME = "accessToken";

/** Backend endpoints (kept in sync with the Go API route table). */
export const API_ENDPOINTS = {
  login: `${BACKEND_URL}/api/auth/login`,
  register: `${BACKEND_URL}/api/auth/register`,
  me: `${BACKEND_URL}/api/auth/me`,
  sendOtp: `${BACKEND_URL}/api/auth/send-otp`,
} as const;
