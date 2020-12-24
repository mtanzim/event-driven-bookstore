import * as faker from "faker";
import { MongoClient } from "mongodb";

export interface Book {
  title: string;
  author: string;
  price: string;
  stock: number;
  stagedQty: number;
}

const uri = process.env.MONGO_URI || "mongodb://localhost:27017";
const dbName = process.env.DB_NAME || "bookstore";
const collName = process.env.MONGO_COLL || "books";
const SIZE = 5;

const client = new MongoClient(uri, { useUnifiedTopology: true });

async function seed() {
  try {
    await client.connect();
    const db = client.db(dbName);
    const coll = db.collection(collName);
    try {
      await coll.drop();
    } catch (err) {
      console.warn(err);
    }
    // placeholder fake data
    const fakeBooks: Book[] = [...Array(SIZE).keys()].map((_) => ({
      title: faker.random.words(3),
      author: `${faker.name.firstName()} ${faker.name.lastName()}`,
      price: faker.commerce.price(10, 300, 2),
      stock: Math.floor(faker.random.number({ min: 2, max: 20 })),
      stagedQty: 0,
    }));
    const res = await coll.insertMany(fakeBooks, { ordered: true });
    console.log(`${res.insertedCount} documents were inserted.`);
  } catch (err) {
    console.error(err);
  } finally {
    await client.close();
  }
}
seed();
