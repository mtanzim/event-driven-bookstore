import React, { useEffect, useState } from "react";
import { Book, BookComponent } from "./Book";
import { CartItem, Cart } from "./CartItem";
import { Checkout } from "./Checkout";
import * as faker from "faker";

async function fetchBooks() {
  const res = await fetch("http://localhost:8080/api/books");
  const books = await res.json();
  return books;
}

export function Store() {
  const [books, setBooks] = useState(new Map<string, Book>());
  const [cart, setCart] = useState(new Map<string, CartItem>());
  const [cartTotal, setCartTotal] = useState(0.0);

  useEffect(() => {
    (async () => {
      const booksRaw: Book[] = await fetchBooks();
      const bookMap = new Map<string, Book>();
      booksRaw.forEach((book) => {
        bookMap.set(book.id, book);
      });
      setBooks(bookMap);
    })();
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

  function onSubmit(data: any) {
    console.log(data);
    console.log(JSON.stringify([...cart.values()]));
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
          <Checkout onSubmit={onSubmit} />
        </div>
      </div>
    </div>
  );
}
