import React from "react";

export interface Book {
  _id: string;
  title: string;
  author: string;
  ISBN: string;
  price: string;
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
        gridTemplateColumns: "6fr 3fr 1fr",
        marginTop: 4,
      }}
    >
      <p>{[book.title, book.author].join(", ")}</p>
      <p>${book.price}</p>
      <button onClick={() => onAction(book._id)}>Add</button>
    </div>
  );
}
