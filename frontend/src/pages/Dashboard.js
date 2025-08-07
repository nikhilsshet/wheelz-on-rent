import React, { useEffect, useState } from "react";
import { jwtDecode } from "jwt-decode";

function Dashboard() {
  const [user, setUser] = useState({ role: "", email: "" });

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      return;
    }

    try {
      const decoded = jwtDecode(token);
      setUser({ role: decoded.role, email: decoded.email });
    } catch (err) {
      console.error("Invalid token");
    }
  }, []);

  return (
    <div style={{ padding: "2rem" }}>
      <h2>Welcome to Dashboard</h2>
      <p>
        <strong>Email:</strong> {user.email}
      </p>
      <p>
        <strong>Role:</strong> {user.role}
      </p>

      {user.role === "customer" && <p>Show customer dashboard content</p>}
      {user.role === "staff" && <p>Show staff dashboard content</p>}
      {user.role === "admin" && <p>Show admin dashboard content</p>}
    </div>
  );
}

export default Dashboard;
