import { getTokenRemainingSeconds } from "@/helper/jwtExp";
import { AUTH_COOKIE_NAME, API_ENDPOINTS } from "@/lib/config";
import { cookies } from "next/headers";
import { NextResponse } from "next/server";

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const cookieStore = await cookies();
    const response = await fetch(API_ENDPOINTS.login, {
      method: "POST",
      body: JSON.stringify(body),
      headers: {
        "Content-Type": "application/json",
      },

      cache: "no-cache",
    });
    const data = await response.json();
    if (!response.ok) {
      return NextResponse.json(
        {
          message: data.message || "Authentication failed",
        },
        {
          status: response.status,
        },
      );
    }

    const accessToken: string = data.token;

    if (!accessToken) {
      return NextResponse.json(
        {
          message: "Backend response succeeded but missing access token.",
        },
        { status: 500 },
      );
    }
    const maxAge: number = getTokenRemainingSeconds(accessToken);
    cookieStore.set({
      name: AUTH_COOKIE_NAME,
      value: accessToken,
      httpOnly: true,
      path: "/",
      maxAge: maxAge,
      secure: false,
    });

    return NextResponse.json(data);
  } catch {
    return NextResponse.json(
      {
        message: "Internal Server Error",
      },
      {
        status: 500,
      },
    );
  }
}
