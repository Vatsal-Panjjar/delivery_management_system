// Base API URL
const API_BASE = "http://localhost:8080";

// ===== LOGIN =====
const loginForm = document.getElementById("loginForm");
if (loginForm) {
    loginForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const username = document.getElementById("username").value;
        const password = document.getElementById("password").value;

        const res = await fetch(`${API_BASE}/login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, password }),
        });

        if (res.ok) {
            const data = await res.json();
            localStorage.setItem("token", data.token);
            localStorage.setItem("role", data.role);
            window.location.href = "dashboard.html";
        } else {
            document.getElementById("loginError").innerText = "Invalid credentials";
        }
    });
}

// ===== DASHBOARD ACTIONS =====
const token = localStorage.getItem("token");

if (document.getElementById("createOrderBtn")) {
    document.getElementById("createOrderBtn").addEventListener("click", async () => {
        const customer = document.getElementById("orderCustomer").value;
        const res = await fetch(`${API_BASE}/orders`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify({ customer })
        });
        alert(res.ok ? "Order Created!" : "Failed to create order");
    });
}

if (document.getElementById("cancelOrderBtn")) {
    document.getElementById("cancelOrderBtn").addEventListener("click", async () => {
        const orderId = document.getElementById("orderIdCancel").value;
        const res = await fetch(`${API_BASE}/orders/${orderId}/cancel`, {
            method: "PUT",
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });
        alert(res.ok ? "Order Cancelled!" : "Failed to cancel order");
    });
}

if (document.getElementById("trackOrderBtn")) {
    document.getElementById("trackOrderBtn").addEventListener("click", async () => {
        const orderId = document.getElementById("orderIdTrack").value;
        const res = await fetch(`${API_BASE}/orders/${orderId}`, {
            method: "GET",
            headers: { "Authorization": `Bearer ${token}` }
        });
        if (res.ok) {
            const data = await res.json();
            document.getElementById("orderStatus").innerText = `Status: ${data.status}`;
        } else {
            document.getElementById("orderStatus").innerText = "Failed to fetch order";
        }
    });
}
