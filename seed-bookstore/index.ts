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
const SIZE = 3;
const STOCK_RANGE = [10, 50];
const STAGED_RANGE = [0, 0];

const client = new MongoClient(uri, { useUnifiedTopology: true });

async function seed() {
  try {
    await client.connect();
    const db = client.db(dbName);
    const coll = db.collection(collName);
    try {
      await coll.drop();
      console.log(`Dropped collection ${collName}`);
    } catch (err) {
      console.warn(err);
    }
    // placeholder fake data
    const [stockMin, stockMax] = STOCK_RANGE;
    const [stagedMin, stagedMax] = STAGED_RANGE;
    const fakeBooks: Book[] = [...Array(SIZE).keys()].map((_) => ({
      title: faker.random.words(3),
      author: `${faker.name.firstName()} ${faker.name.lastName()}`,
      price: faker.commerce.price(10, 300, 2),
      stock: Math.floor(faker.random.number({ min: stockMin, max: stockMax })),
      stagedQty: Math.floor(
        faker.random.number({ min: stagedMin, max: stagedMax })
      ),
    }));
    const res = await coll.insertMany(fakeBooks, { ordered: true });
    console.log(JSON.stringify(fakeBooks, null, 2));
    console.log(`${res.insertedCount} documents were inserted.`);
  } catch (err) {
    console.error(err);
  } finally {
    await client.close();
  }
}
seed();
