import { listBooks } from "./api";

export async function ListBooks() {
    const response = await listBooks();
    return response.data;
}