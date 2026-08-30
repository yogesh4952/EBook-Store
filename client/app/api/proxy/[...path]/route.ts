import { cookies } from "next/headers";
import { NextResponse } from "next/server";

const HOP_BY_HOP_HEADERS = new Set([
  "connection",
  "keep-alive",
  "proxy-authenticate",
  "proxy-authorization",
  "te",
  "trailer",
  "transfer-encoding",
  "upgrade",
  "proxy-connection",
  "host",
  "content-length",
]);

const SAFE_RESPONSE_HEADERS = [
  "content-type",
  "cache-control",
  "etag",
  "last-modified",
  "x-request-id",
  "x-correlation-id",
  "x-ratelimit-limit",
  "x-ratelimit-remaining",
  "x-ratelimit-reset",
  "retry-after",
  "access-control-allow-origin",
  "access-control-allow-methods",
  "access-control-allow-headers",
  "vary",
  "location",
];

export function getBackendUrl(): string {
  const url = process.env.BACKEND_URL;
  if (!url) throw new Error("BACKEND_URL environment variable is not defined");
  return url;
}

async function buildForwardHeaders(request: Request): Promise<Headers> {
  const forwarded = new Headers();

  for (const [key, value] of request.headers.entries()) {
    if (!HOP_BY_HOP_HEADERS.has(key.toLowerCase())) {
      forwarded.set(key, value);
    }
  }

  if (!forwarded.has("authorization")) {
    const cookieStore = await cookies();
    const accessToken = cookieStore.get("accessToken")?.value;
    if (accessToken) {
      forwarded.set("Authorization", `Bearer ${accessToken}`);
    }
  }

  return forwarded;
}

function buildResponseHeaders(backendResponse: Response): Headers {
  const headers = new Headers();

  for (const key of SAFE_RESPONSE_HEADERS) {
    const value = backendResponse.headers.get(key);
    if (value) headers.set(key, value);
  }

  const setCookies = backendResponse.headers.getSetCookie?.() ?? [];
  for (const cookie of setCookies) {
    headers.append("set-cookie", cookie);
  }

  return headers;
}

async function proxy(
  request: Request,
  { params }: { params: Promise<{ path: string[] }> },
) {
  let targetUrl = "";

  try {
    const backendUrl = getBackendUrl();
    const { path } = await params;
    const url = new URL(request.url);

    const joinedPath = path?.join("/") ?? "";
    targetUrl = `${backendUrl}/api/${joinedPath}${url.search}`;

    const isBodyAllowed = request.method !== "GET" && request.method !== "HEAD";

    const response = await fetch(targetUrl, {
      method: request.method,
      headers: await buildForwardHeaders(request),
      body: isBodyAllowed ? request.body : undefined,
      // @ts-expect-error duplex is required for streaming request bodies in Node.js
      duplex: isBodyAllowed ? "half" : undefined,
      cache: "no-store",
      signal: AbortSignal.timeout(30_000),
    });

    if (response.status === 204 || response.status === 304) {
      return new NextResponse(null, {
        status: response.status,
        headers: buildResponseHeaders(response),
      });
    }

    return new NextResponse(response.body, {
      status: response.status,
      headers: buildResponseHeaders(response),
    });
  } catch (error) {
    console.error("Proxy error:", error);

    const isTimeout =
      error instanceof Error &&
      (error.name === "TimeoutError" || error.name === "AbortError");

    return NextResponse.json(
      {
        error: isTimeout
          ? "Backend request timed out"
          : "Failed to connect to backend",
        targetUrl,
        ...(process.env.NODE_ENV !== "production" && {
          details: error instanceof Error ? error.message : String(error),
        }),
      },
      { status: isTimeout ? 504 : 502 },
    );
  }
}

export const GET = proxy;
export const POST = proxy;
export const PUT = proxy;
export const PATCH = proxy;
export const DELETE = proxy;
export const HEAD = proxy;
export const OPTIONS = proxy;