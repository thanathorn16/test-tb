import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import "./common.css"; // Import the CSS file

function Login() {
  const navigate = useNavigate();
  const [form, setForm] = useState({ username: "", password: "" });

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const res = await fetch("http://localhost:8080/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username: form.username,
          password: form.password,
        }),
      });

      if (!res.ok) throw new Error("Login failed");

      const raw = await res.json();
      const data = raw.data;
      console.log("Login successful:", data);
      console.log("Token:", data.token);
      localStorage.setItem("token", data.token);
      navigate("/welcome");
    } catch (err) {
      alert("เข้าสู่ระบบไม่สำเร็จ");
    }
  };

  return (
    <div className="login-container">
      <form onSubmit={handleSubmit}>
        <div className="input-row">
          <label htmlFor="username">User</label>
          <input
            id="username"
            type="text"
            value={form.username}
            onChange={(e) => setForm({ ...form, username: e.target.value })}
          />
        </div>

        <div className="input-row">
          <label htmlFor="password">Password</label>
          <input
            id="password"
            type="password"
            value={form.password}
            onChange={(e) => setForm({ ...form, password: e.target.value })}
          />
        </div>

        <button type="submit" className="login-button">
          ลงชื่อเข้าใช้งาน
        </button>
      </form>

      <div className="spacer-16" />

      <div onClick={() => navigate("/register")} className="register-link">
        สมัครสมาชิก
      </div>
    </div>
  );
}

export default Login;
