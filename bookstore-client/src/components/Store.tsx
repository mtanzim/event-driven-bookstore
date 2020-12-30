import React, { useEffect, useState } from "react";
import { BookComponent } from "./Book";
import { Cart } from "./CartItem";
import { Checkout } from "./Checkout";
import {
  fetchBooks,
  Book,
  CartItem,
  postCart,
  CheckoutFormValues,
  CheckoutDTO,
} from "../api";

export function Store() {
  const [books, setBooks] = useState(new Map<string, Book>());
  const [cart, setCart] = useState(new Map<string, CartItem>());
  const [cartTotal, setCartTotal] = useState(0.0);

  async function getBooks() {
    const booksRaw = await fetchBooks();
    const bookMap = new Map<string, Book>();
    booksRaw.forEach((book) => {
      bookMap.set(book.id, book);
    });
    setBooks(bookMap);
  }

  function checkoutCart(cart: CartItem[]) {
    return async function (values: CheckoutFormValues) {
      const body: CheckoutDTO = {
        items: cart,
        cartUserInformation: values,
      };
      const res = await postCart(body);
      if (res.status == 200) {
        alert("Order requested!");
        getBooks();
        clearCart();
      }
    };
  }

  useEffect(() => {
    getBooks();
  }, []);
  useEffect(calculateTotal, [cart]);

  function calculateTotal() {
    const total = [...cart.values()].reduce(
      (acc, cur) => acc + Number(cur.book.price) * cur.qty,
      0
    );
    setCartTotal(total);
  }

  function addToCart(id: string) {
    if (!books.has(id)) {
      throw new Error("Book not found!");
    }
    const book = books.get(id)!;
    if (cart.has(id)) {
      const bookInCart = cart.get(id)!;
      if (bookInCart.qty === book.stock) {
        return;
      }
      const updatedBookInCart: CartItem = {
        book: bookInCart.book,
        qty: bookInCart.qty + 1,
      };
      cart.set(id, updatedBookInCart);
      setCart(new Map(cart));
    } else {
      cart.set(id, {
        book,
        qty: 1,
      });
      setCart(new Map(cart));
    }
  }

  function removeOneFromCart(id: string) {
    if (!cart.has(id)) {
      throw new Error("Book not found!");
    }
    const cartItem = cart.get(id)!;
    if (cartItem.qty === 1) {
      cart.delete(id);
      setCart(new Map(cart));
    } else {
      const updatedBookInCart: CartItem = {
        book: cartItem.book,
        qty: cartItem.qty - 1,
      };
      cart.set(id, updatedBookInCart);
      setCart(new Map(cart));
    }
  }

  function clearCart() {
    setCart(new Map());
  }

  return (
    <div>
      <h1>Welcome to the Bookstore</h1>
      <div style={{ display: "grid", gridTemplateColumns: "800px 800px" }}>
        <div>
          <h2>Books</h2>
          <ul style={{ listStyle: "none" }}>
            {[...books.values()].map((item) => (
              <BookComponent key={item.id} book={item} onAction={addToCart} />
            ))}
          </ul>
        </div>
        <div>
          <h2>Cart</h2>
          <ul style={{ listStyle: "none" }}>
            {[...cart.values()].map((item) => (
              <Cart
                item={item}
                key={item.book.id}
                onAction={removeOneFromCart}
              />
            ))}
          </ul>
          <button onClick={clearCart}>Clear Cart</button>
        </div>
      </div>
      <div>
        <div style={{ maxWidth: 600, margin: "auto" }}>
          <h2>Checkout</h2>
          <p>Total Cost: ${cartTotal.toFixed(2)}</p>
          <Checkout onSubmit={checkoutCart([...cart.values()])} />
        </div>
      </div>
    </div>
  );
}
