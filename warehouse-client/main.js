const BASE_URL = "http://localhost:8082";
async function fetchData() {
  const res = await fetch(`${BASE_URL}/api/shipment`);
  const shipments = await res.json();
  return shipments;
}

async function shipCart(cartId) {
  const res = await fetch(`${BASE_URL}/api/shipment`, {
    method: "POST",
    body: JSON.stringify({ cartId }),
  });
  if (res.status !== 200) {
    throw new Error("Something went wrong!");
  }
}

function generateBtnId(cartId) {
  return `btn-${cartId}`;
}

function generateBtnHandler(cartId) {
  return async () => {
    try {
      await shipCart(cartId);
      main();
    } catch (err) {
      console.log(err);
      alert(err);
    }
  };
}

function renderData(rows) {
  document.getElementById("root").innerHTML = `<div class="grid-container">
${rows
  .map((row) => {
    const {
      cart: { cartId, items, address, email, phone },
      paid,
    } = row;

    return `
      <div id=${cartId} class="grid-item">
        <h2><code>${cartId}</code>\t${paid ? "Paid" : "Pending"}</h2>
        <h3>Contact</h3>
        <p>${email}<p>
        <p>${phone}<p>
        <p>${address}<p>
        <h3>Books</h3>
        <ul>
        ${items
          .map(
            ({ book, qty }) =>
              `<li>${book.title} - ${book.author} - ${qty}</li>`
          )
          .join("")}
        </ul>
        <button ${paid ? "" : "disabled"} id="${generateBtnId(
      cartId
    )}">Ship it!</button>
      </div>
  `;
  })
  .join("")}
</div>`;
  return Promise.resolve();
}

function attachBtnHandlers(rows) {
  rows.forEach((row) => {
    const {
      cart: { cartId },
    } = row;
    const curBtnId = generateBtnId(cartId);
    const curBtnHandler = generateBtnHandler(cartId);
    document.getElementById(curBtnId).onclick = curBtnHandler;
  });
}

async function main() {
  try {
    document.getElementById("root").innerHTML = `<h4>Loading...</h4>`;
    const rows = await fetchData();
    console.log(JSON.stringify(rows));
    renderData(rows);
    attachBtnHandlers(rows);
  } catch (err) {
    console.log(err);
    document.getElementById("root").innerHTML = `<h4>Failed to load data!</h4>`;
    return;
  }
}

main();
