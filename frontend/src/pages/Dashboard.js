import React, { useEffect, useState } from "react";
import { jwtDecode } from "jwt-decode";
import {
  Box,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Paper,
  TableContainer,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
} from "@mui/material";

export default function Dashboard() {
  const [email, setEmail] = useState("");
  const [role, setRole] = useState("");
  const [bookings, setBookings] = useState([]);
  const [cancelDialogOpen, setCancelDialogOpen] = useState(false);
  const [selectedBooking, setSelectedBooking] = useState(null);

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
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      });

      if (!res.ok) {
        throw new Error("Failed to fetch bookings");
      }

      const data = await res.json();
      setBookings(data);
    } catch (error) {
      console.error("Error fetching bookings:", error);
    }
  };

  const handleCancelBooking = async () => {
    if (!selectedBooking) return;

    const token = localStorage.getItem("token");
    try {
      const res = await fetch(
        `http://localhost:8080/api/bookings/${selectedBooking.booking_id}/cancel`,
        {
          method: "PATCH",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!res.ok) {
        throw new Error("Failed to cancel booking");
      }

      // Close dialog and refresh bookings
      setCancelDialogOpen(false);
      setSelectedBooking(null);
      fetchBookings(token);
    } catch (error) {
      console.error("Error canceling booking:", error);
    }
  };

  return (
    <Box sx={{ p: 4 }}>
      <Typography variant="h4" gutterBottom>
        Welcome to Dashboard
      </Typography>
      <Typography variant="body1">
        <strong>Email:</strong> {email}
      </Typography>
      <Typography variant="body1" sx={{ mb: 3 }}>
        <strong>Role:</strong> {role}
      </Typography>

      {role === "customer" && (
        <>
          <Typography variant="h5" sx={{ mb: 2 }}>
            My Booking History
          </Typography>

          {bookings.length === 0 ? (
            <Typography>No bookings found</Typography>
          ) : (
            <TableContainer component={Paper}>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Booking ID</TableCell>
                    <TableCell>Vehicle</TableCell>
                    <TableCell>Type</TableCell>
                    <TableCell>Model</TableCell>
                    <TableCell>Number Plate</TableCell>
                    <TableCell>Start Date</TableCell>
                    <TableCell>End Date</TableCell>
                    <TableCell>Total Price</TableCell>
                    <TableCell>Status</TableCell>
                    <TableCell>Action</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {bookings.map((b) => (
                    <TableRow key={b.booking_id}>
                      <TableCell>{b.booking_id}</TableCell>
                      <TableCell>{b.vehicle_name}</TableCell>
                      <TableCell>{b.vehicle_type}</TableCell>
                      <TableCell>{b.model}</TableCell>
                      <TableCell>{b.number_plate}</TableCell>
                      <TableCell>
                        {new Date(b.start_date).toLocaleDateString()}
                      </TableCell>
                      <TableCell>
                        {new Date(b.end_date).toLocaleDateString()}
                      </TableCell>
                      <TableCell>{b.total_price}</TableCell>
                      <TableCell>{b.status}</TableCell>
                      <TableCell>
                        {b.status === "active" && (
                          <Button
                            variant="outlined"
                            color="error"
                            size="small"
                            onClick={() => {
                              setSelectedBooking(b);
                              setCancelDialogOpen(true);
                            }}
                          >
                            Cancel
                          </Button>
                        )}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          )}

          {/* Cancel Confirmation Dialog */}
          <Dialog
            open={cancelDialogOpen}
            onClose={() => setCancelDialogOpen(false)}
          >
            <DialogTitle>Cancel Booking</DialogTitle>
            <DialogContent>
              <DialogContentText>
                Are you sure you want to cancel booking ID{" "}
                {selectedBooking?.booking_id}? This action cannot be undone.
              </DialogContentText>
            </DialogContent>
            <DialogActions>
              <Button onClick={() => setCancelDialogOpen(false)}>No</Button>
              <Button color="error" onClick={handleCancelBooking} autoFocus>
                Yes, Cancel
              </Button>
            </DialogActions>
          </Dialog>
        </>
      )}
    </Box>
  );
}
