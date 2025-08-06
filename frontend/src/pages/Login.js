import React, { useState, useContext } from "react";
import api from "../utils/api";
import { AuthContext } from "../authContext";
import { useNavigate } from "react-router-dom";

function Login() {
  const { login } = useContext(AuthContext);

  const [form, setForm] = useState({ email: "", password: "" });
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(""); // Reset error message

    try {
      const res = await fetch("/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: form.email,
          password: form.password,
        }),
      });

      if (res.status === 200) {
        const data = await res.json();
        localStorage.setItem("token", data.token); // Store token in localStorage
        console.log("Token:", data.token);
        // navigate("/dashboard"); // Redirect to home page
      } else {
        const errText = await res.text(); // get server response
        setError(errText); // show error to user
      }
    } catch (err) {
      setError("Login request failed");
      console.error(err);
    }
  };

  return (
    // <div>
    //   <h2>Login</h2>
    //   <form onSubmit={handleSubmit}>
    //     <input
    //       name="email"
    //       type="email"
    //       placeholder="Email"
    //       onChange={handleChange}
    //     />
    //     <input
    //       name="password"
    //       type="password"
    //       placeholder="Password"
    //       onChange={handleChange}
    //     />
    //     <button type="submit">Login</button>
    //   </form>
    //   <p style={{ color: "red" }}>{error}</p>
    // </div>
    <div style={{ padding: "2rem" }}>
      <h2>Login</h2>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <form onSubmit={handleSubmit}>
        <div>
          <label>Email:</label>
          <br />
          <input
            type="email"
            name="email"
            value={form.email}
            onChange={handleChange}
            required
          />
        </div>
        <div style={{ marginTop: "1rem" }}>
          <label>Password:</label>
          <br />
          <input
            type="password"
            name="password"
            value={form.password}
            onChange={handleChange}
            required
          />
        </div>
        <button type="submit" style={{ marginTop: "1rem" }}>
          Login
        </button>
      </form>
    </div>
  );
}

export default Login;
