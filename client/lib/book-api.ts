export interface GenreResponse {
  success: boolean;
  message: string;
  data: string[];
}

export interface CategoryResponse {
  success: boolean;
  message: string;
  data: string[];
}

const BACKEND_URL = process.env.BACKEND_URL || "http://localhost:8080";

export async function getGenresFromApi(): Promise<string[]> {
  let res: Response;
  try {
    res = await fetch(`${BACKEND_URL}/api/book/genre`, {
      next: { revalidate: 3600 },
    });
  } catch (error) {
    console.warn("Failed to connect to book API while loading genres", error);
    return [];
  }
  if (!res.ok) return [];
  const payload = await res.json();
  return payload.data ?? [];
}

export async function getCategoriesFromApi(): Promise<string[]> {
  let res: Response;
  try {
    res = await fetch(`${BACKEND_URL}/api/book/category`, {
      next: { revalidate: 3600 },
    });
  } catch (error) {
    console.warn("Failed to connect to book API while loading categories", error);
    return [];
  }
  if (!res.ok) return [];
  const payload = await res.json();
  return payload.data ?? [];
}
