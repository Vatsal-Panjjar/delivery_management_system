document.getElementById('createOrderForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const form = e.target;
    const data = {
        pickup_address: form.pickup.value,
        delivery_address: form.delivery.value
    };

    const res = await fetch('/orders', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });

    const result = await res.text();
    document.getElementById('ordersList').innerText = result;
});
