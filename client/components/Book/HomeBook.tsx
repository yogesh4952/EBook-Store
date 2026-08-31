import BookCard, { Ibook } from "./BookCard";
import { ListBooks } from "@/lib/book";
import FilterBar from "@/components/common/FilterBar";

const HomeBook = async () => {
  const books: Ibook[] = await ListBooks();
  console.log(books);
  return (
    <div className="grid grid-cols-1 md:grid-cols-12 gap-6">
      {/* Sidebar takes 3 out of 12 columns (~25%) */}
      <div className="md:col-span-3 ">
        <FilterBar />
      </div>

      {/* Main content takes 9 out of 12 columns (~75%) */}
      <div className="md:col-span-9 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mt-4">
        {books.map((book, index) => (
          <BookCard book={book} key={index} />
        ))}
      </div>
    </div>
  );
};

export default HomeBook;
