import React from "react";
import { useForm } from "react-hook-form";
import { CheckoutFormValues } from "../api";

interface Props {
  onSubmit: (data: CheckoutFormValues) => Promise<void>;
}

const isEmptyError = (error: object) =>
  Object.keys(error).length === 0 && error.constructor === Object;

export function Checkout({ onSubmit }: Props) {
  const { register, handleSubmit, errors, reset } = useForm();
  console.log(errors);

  async function submitAndClear(data: CheckoutFormValues) {
    try {
      await onSubmit(data);
      reset();
    } catch (error) {
      console.log(error);
    }
  }

  return (
    <form
      style={{ display: "flex", flexDirection: "column", padding: 10 }}
      onSubmit={handleSubmit(submitAndClear)}
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
        ref={register({ required: true })}
      />
      <input
        type="text"
        placeholder="Secret Code"
        name="code"
        ref={register({ required: true })}
      />
      <input type="submit" />
      <p>{isEmptyError(errors) ? "" : "Error on form"}</p>
    </form>
  );
}
