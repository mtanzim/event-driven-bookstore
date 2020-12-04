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
  const isMax = item.qty === item.book.stock;
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
      <p style={{ color: isMax ? "red" : "black" }}>{item.qty}</p>
      <button onClick={() => onAction(item.book.id)}>Delete</button>
    </div>
  );
}
