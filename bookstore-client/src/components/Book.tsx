import React from "react";

export interface Book {
  id: string;
  title: string;
  author: string;
  price: string;
  stock: number;
}

interface Props {
  book: Book;
  onAction: (id: string) => void;
}

export function BookComponent({ book, onAction }: Props) {
  return (
    <div
      style={{
        display: "grid",
        gridTemplateColumns: "6fr 2fr 1fr 1fr",
        marginTop: 4,
      }}
    >
      <p>{[book.title, book.author].join(", ")}</p>
      <p>${book.price}</p>
      <p>{book.stock}</p>
      <button onClick={() => onAction(book.id)}>Add</button>
    </div>
  );
}
