import { API_ENDPOINTS } from "@/lib/config";
import { NextResponse } from "next/server";

export async function POST(request: Request) {
  try {
    const body = await request.json();

    const response = await fetch(API_ENDPOINTS.sendOtp, {
      method: "POST",
      body: JSON.stringify(body),
      headers: {
        "Content-type": "application/json",
      },
      cache: "no-cache",
    });

    const data = await response.json();

    if (!data.success) {
      return NextResponse.json(
        {
          message: data.message || "Invalid input",
        },
        {
          status: response.status,
        },
      );
    }

    return NextResponse.json({ success: true });
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
