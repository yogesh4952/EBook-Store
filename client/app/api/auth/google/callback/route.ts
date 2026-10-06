import { BACKEND_URL } from "@/lib/config";
import { NextResponse } from "next/server";

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const code = searchParams.get("code");
  const error = searchParams.get("error");

  if (error || !code) {
    return NextResponse.redirect(
      new URL("/login?error=OAuthCancelled", request.url),
    );
  }

  try {
    const backendResponse = await fetch(`${BACKEND_URL}/api/auth/google`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ code }),
    });

    if (!backendResponse.ok) {
      return NextResponse.redirect(
        new URL("/login?error=AuthenticationFailed", request.url),
      );
    }
    const data = await backendResponse.json();

    const response = NextResponse.redirect(new URL("/", request.url));

    response.cookies.set("app_session", data.token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      path: "/",
    });

    return response;
  } catch (err) {
    console.error("BFF OAuth Error:", err);
    return NextResponse.redirect(
      new URL("/login?error=ServerError", request.url),
    );
  }
}
