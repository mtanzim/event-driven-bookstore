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
