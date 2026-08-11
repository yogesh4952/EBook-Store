import { NextRequest } from "next/server";

export async function GET(request: NextRequest) {
  const body = request.json();
  console.log(body);
}
