import * as faker from "faker";
// placeholder fake data
const SIZE = 3;
const fakeBooks = [...Array(SIZE).keys()].map((_) => ({
  _id: faker.random.uuid(),
  title: faker.random.words(3),
  author: `${faker.name.firstName()} ${faker.name.lastName()}`,
  ISBN: faker.phone.phoneNumber(),
  price: faker.commerce.price(10, 300, 2),
}));

console.log(JSON.stringify(fakeBooks, null, 2));
