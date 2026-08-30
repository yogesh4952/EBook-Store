import BookCard, { Ibook } from './BookCard'
import { ListBooks } from '@/lib/book'




const HomeBook = async() => {

  const books:Ibook[] = await ListBooks();
  console.log(books)
  return (
      <div className='grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 lg:grid-cols-5 mt-4 gap-4'>
      {books.map((book,index)=>(
        <BookCard book={book}  key={index}/>
      ))}
      </div>
      
  )
}

export default HomeBook
