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

export interface CheckoutFormValues {
  address: string;
  cardNum: string;
  code: string;
  email: string;
  phone: string;
}

export interface CheckoutDTO {
  items: CartItem[];
  cartUserInformation: CheckoutFormValues;
}

export async function postCart(body: CheckoutDTO) {
  return fetch(`${BASE_API}/api/checkout`, {
    method: "POST",
    body: JSON.stringify(body),
  });
}
