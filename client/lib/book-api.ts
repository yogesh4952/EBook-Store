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

export async function getGenresFromApi(): Promise<string[]> {
  const res = await fetch("http://localhost:8080/api/book/genre", { next: { revalidate: 3600 } });
  if (!res.ok) return [];
  const payload = await res.json();
  return payload.data ?? [];
}

export async function getCategoriesFromApi(): Promise<string[]> {
  const res = await fetch("http://localhost:8080/api/book/category", { next: { revalidate: 3600 } });
  if (!res.ok) return [];
  const payload = await res.json();
  return payload.data ?? [];
}
