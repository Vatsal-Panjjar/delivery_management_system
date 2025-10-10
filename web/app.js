async function loadDeliveries() {
  const token = localStorage.getItem("token");
  if (!token) {
    window.location.href = "login.html";
    return;
  }

  const res = await fetch("/deliveries?status=pending", {
    headers: { Authorization: `Bearer ${token}` },
  });

  const deliveries = await res.json();
  const container = document.getElementById("deliveries");

  if (!deliveries || deliveries.length === 0) {
    container.innerHTML = "<p>No deliveries found.</p>";
    return;
  }

  container.innerHTML = deliveries
    .map(
      (d) => `
    <div class="delivery-card">
      <h4>Delivery ID: ${d.id}</h4>
      <p>Status: ${d.status}</p>
      <p>Pickup: ${d.pickup_address}</p>
      <p>Dropoff: ${d.dropoff_address}</p>
    </div>
  `
    )
    .join("");
}

function logout() {
  localStorage.removeItem("token");
  window.location.href = "login.html";
}

window.onload = loadDeliveries;
