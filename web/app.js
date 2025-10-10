// Detect which page we are on
if (document.getElementById("loginForm")) {
  // Login page
  const loginForm = document.getElementById("loginForm");
  loginForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    const res = await fetch("http://localhost:8080/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password })
    });

    const data = await res.json();
    if (res.ok) {
      localStorage.setItem("token", data.token);
      window.location.href = "dashboard.html";
    } else {
      document.getElementById("loginMsg").innerText = data.error || "Login failed";
    }
  });
}

if (document.getElementById("deliveriesList")) {
  // Dashboard page
  const token = localStorage.getItem("token");
  if (!token) window.location.href = "index.html";

  async function loadDeliveries() {
    const res = await fetch("http://localhost:8080/deliveries?status=pending", {
      headers: { "Authorization": "Bearer " + token }
    });
    const deliveries = await res.json();
    const list = document.getElementById("deliveriesList");
    list.innerHTML = "";
    deliveries.forEach(d => {
      const li = document.createElement("li");
      li.innerText = `Pickup: ${d.pickup_address}, Dropoff: ${d.dropoff_address}, Status: ${d.status}`;
      list.appendChild(li);
    });
  }

  loadDeliveries();

  document.getElementById("logoutBtn").addEventListener("click", () => {
    localStorage.removeItem("token");
    window.location.href = "index.html";
  });
}
