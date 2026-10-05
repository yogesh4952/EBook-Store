import { decodeJwtClaims } from "@/helper/jwtExp";
import { AUTH_COOKIE_NAME } from "@/lib/config";
import { cookies } from "next/headers";
import { NextResponse } from "next/server";

/**
 * Returns the caller's identity from the access token cookie.
 *
 * Claims are read straight out of the token instead of calling the backend,
 * because the Go API has no /auth/me endpoint to call. That is acceptable
 * here: the token's signature was already verified when it was issued, and
 * every other request re-verifies it server-side. This is a convenience
 * endpoint for the UI, not an authorization gate.
 */
export async function GET() {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get(AUTH_COOKIE_NAME)?.value;

    if (!token) {
      return NextResponse.json(
        { message: "Unauthorized" },
        { status: 401 },
      );
    }

    const claims = decodeJwtClaims(token);

    if (!claims) {
      return NextResponse.json(
        { message: "Invalid access token" },
        { status: 401 },
      );
    }

    if (claims.exp && claims.exp * 1000 <= Date.now()) {
      return NextResponse.json(
        { message: "Access token expired" },
        { status: 401 },
      );
    }

    return NextResponse.json({
      success: true,
      data: {
        id: claims.user_id,
        email: claims.email,
        role: claims.role,
      },
    });
  } catch {
    return NextResponse.json(
      { message: "Internal Server Error" },
      { status: 500 },
    );
  }
}