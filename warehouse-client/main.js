const BASE_URL = "http://localhost:8082/";
async function fetchData() {
  const res = await fetch(`${BASE_URL}/api/shipment`);
  const shipments = await res.json();
  return shipments;
}

async function main() {
  try {
    document.getElementById("root").innerHTML = `<h4>Loading...</h4>`;
    const rows = await fetchData();
    console.log(JSON.stringify(rows));
    document.getElementById("root").innerHTML = `<div class="grid-container"> 
  ${rows
    .map((row) => {
      const {
        cart: { cartId: _id, items, address, email, phone },
        shipped,
        paid,
      } = row;
      // TODO: do this from the backend
      if (shipped) {
        return "";
      }
      return `
        <div id=${_id} class="grid-item">
          <h2><code>${_id}</code>\t${paid ? "Paid" : "Pending"}</h2>
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
          <button ${paid ? "" : "disabled"} id="btn-${_id}">Ship it!</button>
        </div>
    `;
    })
    .join("")}
  </div>`;
  } catch (err) {
    console.log(err);
    document.getElementById("root").innerHTML = `<h4>Failed to load data!</h4>`;
  }
}

main();
