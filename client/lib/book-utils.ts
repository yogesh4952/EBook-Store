export interface Ibook {
  id: number;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
  title: string;
  author_name: string;
  genre: string;
  category: string;
  pages: number;
  publication: string;
  price: number;
  units: number;
  status: "in_stock" | "out_of_stock";
  cover_page_url: string;
  seller_id: number | null;
  seller?: {
    id: number;
    created_at: string;
    updated_at: string;
    deleted_at: string | null;
    user_id: number;
    seller_number: number;
    user?: {
      id: number;
      created_at: string;
      updated_at: string;
      deleted_at: string | null;
      first_name: string;
      last_name: string;
      email: string;
      role: "seller";
      phone_number: string;
    };
  };
}

interface BrowseByCategory {
  category: string;
  count: number;
}

export const getCategories = (books: Ibook[]): BrowseByCategory[] => {
  const categoryCounts = books.reduce((acc, book) => {
    acc[book.category] = (acc[book.category] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  return Object.entries(categoryCounts)
    .map(([category, count]) => ({ category, count }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 10);
};

type GenreOrder = {
  Programming: number;
  Fantasy: number;
  Romance: number;
  "Self Help": number;
  Psychology: number;
  Finance: number;
  "Historical Fiction": number;
  Mystery: number;
  "Psychological Thriller": number;
  Drama: number;
};

export const getGenres = (books: Ibook[]): string[] => {
  const genreSet = new Set<string>();
  books.forEach((book) => {
    if (book.genre) {
      genreSet.add(book.genre);
    }
  });

  const order: GenreOrder = {
    Programming: 1,
    Fantasy: 2,
    Romance: 3,
    "Self Help": 4,
    Psychology: 5,
    Finance: 6,
    "Historical Fiction": 7,
    Mystery: 8,
    "Psychological Thriller": 9,
    Drama: 10,
  };

  return Array.from(genreSet)
    .filter(Boolean)
    .sort((a, b) => {
      const rankA = order[a as keyof GenreOrder] || 99;
      const rankB = order[b as keyof GenreOrder] || 99;
      return rankA - rankB;
    });
};