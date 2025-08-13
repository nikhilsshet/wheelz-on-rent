import React, { useEffect, useState } from "react";
import { jwtDecode } from "jwt-decode";

export default function Dashboard() {
  const [email, setEmail] = useState("");
  const [role, setRole] = useState("");
  const [bookings, setBookings] = useState([]);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      const decoded = jwtDecode(token);
      setEmail(decoded.email);
      setRole(decoded.role);

      if (decoded.role === "customer") {
        fetchBookings(token);
      }
    }
  }, []);

  const fetchBookings = async (token) => {
    try {
      const res = await fetch("http://localhost:8080/api/mybookings", {
        method: "GET",
        credentials: "include", // Include cookies if needed
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      });

      if (!res.ok) {
        throw new Error("Failed to fetch bookings");
      }

      const data = await res.json();
      console.log("fetched bookings:", data);
      setBookings(data);
    } catch (error) {
      console.error("Error fetching bookings:", error);
    }
  };

  return (
    <div>
      <h1>Welcome to Dashboard</h1>
      <p>
        <strong>Email:</strong> {email}
      </p>
      <p>
        <strong>Role:</strong> {role}
      </p>

      {role === "customer" && (
        <>
          <h2>My Booking History</h2>
          {bookings.length === 0 ? (
            <p>No bookings found</p>
          ) : (
            <table border="1" cellPadding="5">
              <thead>
                <tr>
                  <th>Booking ID</th>
                  <th>Vehicle</th>
                  <th>Type</th>
                  <th>Model</th>
                  <th>Number Plate</th>
                  <th>Start Date</th>
                  <th>End Date</th>
                  <th>Total Price</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {bookings.map((b) => (
                  <tr key={b.booking_id}>
                    <td>{b.booking_id}</td>
                    <td>{b.vehicle_name}</td>
                    <td>{b.vehicle_type}</td>
                    <td>{b.model}</td>
                    <td>{b.number_plate}</td>
                    <td>{b.start_date}</td>
                    <td>{b.end_date}</td>
                    <td>{b.total_price}</td>
                    <td>{b.status}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </>
      )}
    </div>
  );
}
