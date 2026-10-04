import BookCard, { type Ibook } from "./BookCard";
import { ListBooks } from "@/lib/book";

const HomeBook = async () => {
  const books: Ibook[] = await ListBooks();

  if (!books?.length) {
    return (
      <p className="mt-6 rounded-xl border border-border bg-surface px-6 py-10 text-center text-sm text-muted">
        No books published yet.
      </p>
    );
  }

  return (
    <section className="mt-5">
      <div className="mb-3 flex items-baseline justify-between">
        <h2 className="text-lg font-semibold text-primary">Browse books</h2>
        <span className="text-xs text-muted">{books.length} titles</span>
      </div>

      <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-5 xl:grid-cols-6">
        {books.map((book) => (
          <BookCard book={book} key={book.id} />
        ))}
      </div>
    </section>
  );
};

export default HomeBook;