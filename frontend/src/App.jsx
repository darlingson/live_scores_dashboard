import React from "react";
import { Routes, Route, Link } from "react-router-dom";
import SendEventPage from "./pages/SendEventPage";
import UpdatesPage from "./pages/UpdatesPage";

function App() {
  return (
    <div>
      <nav>
        <Link to="/send-event">Send Event</Link> |{" "}
        <Link to="/updates">View Updates</Link>
      </nav>
      <Routes>
        <Route path="/send-event" element={<SendEventPage />} />
        <Route path="/updates" element={<UpdatesPage />} />
      </Routes>
    </div>
  );
}

export default App;