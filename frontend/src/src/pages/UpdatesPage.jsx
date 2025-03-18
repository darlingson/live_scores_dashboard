import React, { useEffect, useState } from "react";

function UpdatesPage() {
  const [events, setEvents] = useState([]);

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/ws");

    ws.onmessage = (event) => {
      const eventData = JSON.parse(event.data);
      setEvents((prevEvents) => [...prevEvents, eventData]);
    };

    return () => {
      ws.close();
    };
  }, []);

  return (
    <div>
      <h1>Match Updates</h1>
      <ul>
        {events.map((event, index) => (
          <li key={index}>
            {event.type === "goal"
              ? `${event.scorer} scored at ${event.time}. Score: ${event.score}`
              : `${event.type.toUpperCase()} event at ${event.time}`}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default UpdatesPage;