import React, { useEffect, useState } from "react";
import { Book } from "./Book";

export interface CartItem {
  book: Book;
  qty: number;
}

interface Props {
  item: CartItem;
  onAction: (id: string) => void;
}

export function Cart({ item, onAction }: Props) {
  return (
    <div
      style={{
        marginTop: 4,
        display: "grid",
        gridTemplateColumns: "6fr 2fr 1fr 1fr",
      }}
    >
      <p>{[item.book.title, item.book.author].join(", ")}</p>
      <p>${item.book.price}</p>
      <p>{item.qty}</p>
      <button onClick={() => onAction(item.book._id)}>Delete</button>
    </div>
  );
}
