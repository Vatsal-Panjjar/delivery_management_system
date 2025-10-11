document.addEventListener('DOMContentLoaded', () => {

    // Handle login
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const username = loginForm.username.value;
            const password = loginForm.password.value;
            const res = await fetch('/api/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password })
            });
            if (res.ok) {
                window.location.href = '/dashboard';
            } else {
                alert('Login failed');
            }
        });
    }

    // Handle registration
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const username = registerForm.username.value;
            const password = registerForm.password.value;
            const res = await fetch('/api/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password })
            });
            if (res.ok) {
                window.location.href = '/dashboard';
            } else {
                alert('Registration failed');
            }
        });
    }

    // Handle order creation and listing
    const createOrderForm = document.getElementById('createOrderForm');
    const ordersTable = document.getElementById('ordersTable');
    if (createOrderForm && ordersTable) {
        const fetchOrders = async () => {
            const res = await fetch('/api/orders');
            const orders = await res.json();
            const tbody = ordersTable.querySelector('tbody');
            tbody.innerHTML = '';
            orders.forEach(o => {
                const tr = document.createElement('tr');
                tr.innerHTML = `
                    <td>${o.id}</td>
                    <td>${o.source}</td>
                    <td>${o.destination}</td>
                    <td>${o.status}</td>
                    <td>
                        ${o.status !== 'cancelled' ? `<button onclick="cancelOrder('${o.id}')">Cancel</button>` : ''}
                    </td>
                `;
                tbody.appendChild(tr);
            });
        };

        createOrderForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const source = createOrderForm.source.value;
            const destination = createOrderForm.destination.value;
            const res = await fetch('/api/orders', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ source, destination })
            });
            if (res.ok) {
                createOrderForm.reset();
                fetchOrders();
            } else {
                alert('Failed to create order');
            }
        });

        window.cancelOrder = async (id) => {
            const res = await fetch(`/api/orders/${id}/cancel`, { method: 'POST' });
            if (res.ok) fetchOrders();
        };

        fetchOrders();
    }

    // Handle admin dashboard
    const adminOrdersTable = document.getElementById('adminOrdersTable');
    if (adminOrdersTable) {
        const fetchAdminOrders = async () => {
            const res = await fetch('/admin');
            const orders = await res.json();
            const tbody = adminOrdersTable.querySelector('tbody');
            tbody.innerHTML = '';
            orders.forEach(o => {
                const tr = document.createElement('tr');
                tr.innerHTML = `
                    <td>${o.id}</td>
                    <td>${o.user_id}</td>
                    <td>${o.source}</td>
                    <td>${o.destination}</td>
                    <td>${o.status}</td>
                    <td>
                        <select onchange="updateStatus('${o.id}', this.value)">
                            <option value="pending" ${o.status==='pending'?'s
