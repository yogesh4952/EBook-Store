import { NextResponse } from "next/server";

export async function POST(request: Request) {
  try {
    const body = request.json();
    console.log(body);
    const response = await fetch(`${process.env.BACKEND_URL}/api/Auth/login`, {
      method: "POST",
      body: JSON.stringify(body),
    });
    const data = await response.json();

    console.log(data.user);
    return NextResponse.json({
      message: "success",
    });
  } catch (error) {
    return NextResponse.json(error);
  }
}
