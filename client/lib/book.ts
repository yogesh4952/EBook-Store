export async function ListBooks() {
    const res = await fetch("http://localhost:8080/api/book/list-books")

    const data = await res.json()
    return data.data
    
}