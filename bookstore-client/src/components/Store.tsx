import React, { useEffect, useState } from "react";
import * as faker from "faker";

interface Book {
  _id: string;
  title: string;
  author: string;
  ISBN: string;
  price: string;
}

interface CartItem {
  book: Book;
  qty: number;
}

// placeholder fake data
const fakeBooks = [...Array(10).keys()].map((_) => ({
  _id: faker.random.uuid(),
  title: faker.random.words(4),
  author: `${faker.name.firstName()} ${faker.name.lastName()}`,
  ISBN: faker.phone.phoneNumber(),
  price: faker.commerce.price(10, 300, 2),
}));

interface BookProps {
  book: Book;
  onAction: (id: string) => void;
}

function BookComponent({ book, onAction }: BookProps) {
  return (
    <div style={{ display: "flex" }}>
      <li>{[book.title, book.author, `$${book.price}`].join(", ")}</li>
      <button onClick={() => onAction(book._id)} style={{ marginLeft: 8 }}>
        Add to Cart
      </button>
    </div>
  );
}

interface CartProps {
  item: CartItem;
  onAction: (id: string) => void;
}

function Cart({ item, onAction }: CartProps) {
  return (
    <div style={{ display: "flex" }}>
      <li>
        {[
          item.book.title,
          item.book.author,
          `$${item.book.price}`,
          item.qty,
        ].join(", ")}
      </li>
      <button onClick={() => onAction(item.book._id)} style={{ marginLeft: 8 }}>
        Remove
      </button>
    </div>
  );
}

export function Store() {
  const [books, setBooks] = useState(new Map<string, Book>());
  const [cart, setCart] = useState(new Map<string, CartItem>());
  const [cartTotal, setCartTotal] = useState(0.0);
  useEffect(() => {
    const bookMap = new Map<string, Book>();
    fakeBooks.forEach((book) => bookMap.set(book._id, book));
    setBooks(bookMap);
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
      <h2>Books</h2>
      <ul style={{ listStyle: "none" }}>
        {[...books.values()].map((item) => (
          <BookComponent book={item} onAction={addToCart} />
        ))}
      </ul>
      <h2>Cart</h2>
      <ul style={{ listStyle: "none" }}>
        {[...cart.values()].map((item) => (
          <Cart item={item} onAction={removeOneFromCart} />
        ))}
      </ul>
      <button onClick={clearCart}>Clear Cart</button>
      <h2>Checkout</h2>
      <p>Total Cost: ${cartTotal.toFixed(2)}</p>
    </div>
  );
}
