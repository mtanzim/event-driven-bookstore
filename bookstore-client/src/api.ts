const BASE_API = "http://localhost:8080";

export interface Book {
  id: string;
  title: string;
  author: string;
  price: string;
  stock: number;
}

export interface CartItem {
  book: Book;
  qty: number;
}

export async function fetchBooks(): Promise<Book[]> {
  const res = await fetch(`${BASE_API}/api/books`);
  const books = await res.json();
  return books || [];
}
