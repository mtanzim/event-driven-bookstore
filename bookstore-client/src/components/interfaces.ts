import { CartItem } from "../api";

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
