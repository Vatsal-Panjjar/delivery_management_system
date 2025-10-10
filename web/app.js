// ------------------ LOGIN ------------------
if (document.getElementById("loginForm")) {
  const loginForm = document.getElementById("loginForm");
  loginForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    try {
      const res = await fetch("http://localhost:8080/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password })
      });
      const data = await res.json();
      if (res.ok) {
        localStorage.setItem("jwtToken", data.token);
        window.location.href = "dashboard.html";
      } else {
        document.getElementById("loginMsg").innerText = data.error || "Login failed";
      }
    } catch (err) {
      console.error(err);
      document.getElementById("loginMsg").innerText = "Network error";
    }
  });
}

// ------------------ SIGNUP ------------------
if (document.getElementById("signupForm")) {
  const signupForm = document.getElementById("signupForm");
  signupForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const username = document.getElementById("username").value;
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;

    try {
      const res = await fetch("http://localhost:8080/signup", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, email, password })
      });
      const data = await res.json();
      if (res.ok) {
        alert("Signup successful! Please login.");
        window.location.href = "index.html";
      } else {
        document.getElementById("signupMsg").innerText = data.error || "Signup failed";
      }
    } catch (err) {
      console.error(err);
      document.getElementById("signupMsg").innerText = "Network error";
    }
  });
}

// ------------------ DASHBOARD: FETCH D
