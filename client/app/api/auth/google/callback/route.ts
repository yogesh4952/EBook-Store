import { BACKEND_URL } from "@/lib/config";
import { NextResponse } from "next/server";

export async function GET(request: Request) {
  console.log("[BFF] CALLBACK HIT", new URL(request.url).toString());
  console.log("[BFF] BACKEND_URL:", process.env.BACKEND_URL || "NOT SET");
  const { searchParams } = new URL(request.url);
  const code = searchParams.get("code");
  const error = searchParams.get("error");

  if (error || !code) {
    return NextResponse.redirect(
      new URL("/login?error=OAuthCancelled", request.url),
    );
  }

  try {
    console.log("[BFF] sending to backend code length:", code?.length);
    const backendResponse = await fetch(`${BACKEND_URL}/api/auth/google`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ code }),
    });

    console.log(
      "[BFF] backend status:",
      backendResponse.status,
      "ok:",
      backendResponse.ok,
    );
    if (!backendResponse.ok) {
      const text = await backendResponse.text();
      console.log(
        "[BFF] backend error response:",
        text,
      );
      const redirectUrl = text.includes("duplicate key")
        ? new URL("/login?error=DBConflict", request.url)
        : new URL("/login?error=AuthenticationFailed", request.url);
      return NextResponse.redirect(redirectUrl);
    }
    const data = await backendResponse.json();

    const response = NextResponse.redirect(new URL("/", request.url));

    response.cookies.set("accessToken", data.token, {
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
