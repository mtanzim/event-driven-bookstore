import { items as rows } from "./mock.js";

function main() {
  console.log(JSON.stringify(rows));
  document.getElementById("root").innerHTML = `<div class="grid-container"> 
  ${rows
    .map((row) => {
      const {
        cart: { _id, items, address, email, phone },
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
}

main();
