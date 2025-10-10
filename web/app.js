document.getElementById("loginForm")?.addEventListener("submit", async (e) => {
    e.preventDefault();
    const res = await fetch("/auth/login", {
        method: "POST",
        headers: {"Content-Type":"application/json"},
        body: JSON.stringify({
            username: document.getElementById("username").value,
            password: document.getElementById("password").value
        })
    });
    const data = await res.json();
    localStorage.setItem("token", data.token);
    window.location.href = "dashboard.html";
});

async function loadDeliveries() {
    const token = localStorage.getItem("token");
    const res = await fetch("/deliveries?status=pending", {
        headers: { "Authorization": `Bearer ${token}` }
    });
    const deliveries = await res.json();
    const ul = document.getElementById("deliveries");
    if (ul) deliveries.forEach(d => {
        const li = document.createElement("li");
        li.innerText = `${d.pickupAddr} -> ${d.dropoffAddr} [${d.status}]`;
        ul.appendChild(li);
    });
}
loadDeliveries();
