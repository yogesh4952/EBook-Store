import HomeBook from "@/components/Book/HomeBook";
import { getAllBooks } from "@/lib/book-actions";
import Homepage from "@/components/Home/Homepage";

export default async function Home() {
  const books = await getAllBooks();

  return (
    <div>
      <Homepage />
      <HomeBook books={books} />
    </div>
  );
}