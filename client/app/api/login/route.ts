import { getTokenRemainingSeconds } from "@/helper/jwtExp";
import { cookies } from "next/headers";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const cookieStore = await cookies();
    const response = await fetch(`${process.env.BACKEND_URL}/api/Auth/login`, {
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
      name: "accessToken",
      value: accessToken,
      httpOnly: process.env.NODE_ENV == "production",
      path: "/",
      maxAge: maxAge,
      secure: process.env.NODE_ENV == "production",
    });

    return NextResponse.json(data);
  } catch (error) {
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
