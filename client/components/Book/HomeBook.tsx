import BookCard, { Ibook } from "./BookCard";
import { ListBooks } from "@/lib/book";
import FilterBar from "@/components/common/FilterBar";

const HomeBook = async () => {
  const books: Ibook[] = await ListBooks();
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4 mt-4">
      {books.map((book, index) => (
        <BookCard book={book} key={index} />
      ))}
    </div>
  );
};

export default HomeBook;
