import { API_ENDPOINTS } from "@/lib/config";
import { NextResponse } from "next/server";

export async function POST(request: Request) {
  try {
    const body = await request.json();

    const response = await fetch(API_ENDPOINTS.register, {
      method: "POST",
      body: JSON.stringify(body),
      headers: { "Content-Type": "application/json" },
      cache: "no-cache",
    });

    const data = await response.json();
    if (!response.ok) {
      return NextResponse.json(
        { message: data.message || "Registration failed" },
        { status: response.status },
      );
    }

    return NextResponse.json({ success: true, message: data.message });
  } catch {
    return NextResponse.json(
      { message: "Internal Server Error" },
      { status: 500 },
    );
  }
}
