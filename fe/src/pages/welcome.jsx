import React, { useEffect, useState } from "react";

function Welcome() {
  const [username, setUsername] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchProfile = async () => {
      const token = localStorage.getItem("token");
      if (!token) {
        alert("Token missing, please login again.");
        return;
      }

      try {
        const res = await fetch("http://localhost:8080/profile", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!res.ok) throw new Error("Failed to fetch profile");

        const raw = await res.json();
        const data = raw.data;
        console.log("Profile data:", data);
        setUsername(data.userName);
        setLoading(false);
      } catch (err) {
        alert("Failed to load profile");
        setLoading(false);
      }
    };

    fetchProfile();
  }, []);

  if (loading) return <div style={styles.container}>Loading...</div>;

  return (
    <div style={styles.container}>
      <h1>Welcome user: {username}</h1>
    </div>
  );
}

const styles = {
  container: {
    height: "100vh",
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    fontFamily: "sans-serif",
  },
};

export default Welcome;
