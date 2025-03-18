import React, { useState } from "react";

function SendEventPage() {
  const [eventType, setEventType] = useState("");
  const [scorer, setScorer] = useState("");
  const [time, setTime] = useState("");
  const [score, setScore] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();

    const event = {
      type: eventType,
      scorer,
      time,
      score,
    };

    try {
      await fetch("http://localhost:8080/event", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(event),
      });
      alert("Event sent successfully!");
    } catch (error) {
      console.error("Error sending event:", error);
      alert("Failed to send event.");
    }
  };

  return (
    <div>
      <h1>Send Match Event</h1>
      <form onSubmit={handleSubmit}>
        <label>
          Event Type:
          <input
            type="text"
            value={eventType}
            onChange={(e) => setEventType(e.target.value)}
            required
          />
        </label>
        <br />
        <label>
          Scorer (for goals):
          <input
            type="text"
            value={scorer}
            onChange={(e) => setScorer(e.target.value)}
          />
        </label>
        <br />
        <label>
          Time:
          <input
            type="text"
            value={time}
            onChange={(e) => setTime(e.target.value)}
            required
          />
        </label>
        <br />
        <label>
          Score:
          <input
            type="text"
            value={score}
            onChange={(e) => setScore(e.target.value)}
          />
        </label>
        <br />
        <button type="submit">Send Event</button>
      </form>
    </div>
  );
}

export default SendEventPage;