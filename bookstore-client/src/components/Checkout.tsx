import React from "react";
import { useForm } from "react-hook-form";
import { CheckoutFormValues } from "./interfaces";

interface Props {
  onSubmit: (data: CheckoutFormValues) => Promise<string>;
}
export function Checkout({ onSubmit }: Props) {
  const { register, handleSubmit, errors } = useForm();
  console.log(errors);

  return (
    <form
      style={{ display: "flex", flexDirection: "column", padding: 10 }}
      onSubmit={handleSubmit(onSubmit)}
    >
      <input
        type="text"
        placeholder="Email"
        name="email"
        ref={register({ required: true, pattern: /^\S+@\S+$/i })}
      />
      <input
        type="tel"
        placeholder="Mobile number"
        name="phone"
        ref={register({ required: true, maxLength: 12 })}
      />
      <textarea
        placeholder="Address"
        name="address"
        ref={register({ required: true })}
      />
      <input
        type="number"
        placeholder="Card Number"
        name="cardNum"
        ref={register}
      />
      <input
        type="text"
        placeholder="Secret Code"
        name="code"
        ref={register({ required: true })}
      />
      <input type="submit" />
    </form>
  );
}
