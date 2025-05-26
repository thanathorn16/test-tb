import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import "./common.css"; // Import the CSS file

function Register() {
  const navigate = useNavigate();
  const [form, setForm] = useState({
    username: "",
    password: "",
    confirmPassword: "",
  });
  const [error, setError] = useState("");

  const passwordsMatch = form.password === form.confirmPassword;
  const canSubmit = form.username && form.password && passwordsMatch;

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!passwordsMatch) {
      setError("Password และ Confirm Password ต้องตรงกัน");
      return;
    }
    setError("");

    try {
      const res = await fetch("http://localhost:8080/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username: form.username,
          password: form.password,
        }),
      });

      if (!res.ok) throw new Error("Register failed");

      navigate("/login");
    } catch (err) {
      alert("สมัครสมาชิกไม่สำเร็จ");
    }
  };

  return (
    <div className="login-container">
      <h2>สมัครสมาชิก</h2>
      <form onSubmit={handleSubmit}>
        <div className="input-row">
          <label htmlFor="username">User</label>
          <input
            id="username"
            type="text"
            placeholder="User"
            value={form.username}
            onChange={(e) => setForm({ ...form, username: e.target.value })}
            required
          />
        </div>

        <div className="input-row">
          <label htmlFor="password">Password</label>
          <input
            id="password"
            type="password"
            placeholder="Password"
            value={form.password}
            onChange={(e) => setForm({ ...form, password: e.target.value })}
            required
          />
        </div>

        <div className="input-row">
          <label htmlFor="confirmPassword">Confirm Password</label>
          <input
            id="confirmPassword"
            type="password"
            placeholder="Confirm Password"
            value={form.confirmPassword}
            onChange={(e) =>
              setForm({ ...form, confirmPassword: e.target.value })
            }
            required
          />
        </div>

        {error && <p style={{ color: "red", marginTop: "8px" }}>{error}</p>}

        <button type="submit" className="login-button" disabled={!canSubmit}>
          สมัครสมาชิก
        </button>
      </form>
    </div>
  );
}

export default Register;
