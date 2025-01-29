// /src/ws-client/WebSocketClient.js

import React, { useEffect, useState } from "react";

const WebSocketClient = () => {
  const [messages, setMessages] = useState([]);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws"); // Your WebSocket server URL

    socket.onopen = () => {
      console.log("Connected to WebSocket server");
    };

    socket.onmessage = (event) => {
      const newMessage = event.data;
      setMessages((prevMessages) => [...prevMessages, newMessage]);
    };

    socket.onclose = () => {
      console.log("WebSocket connection closed");
    };

    // Clean up WebSocket connection when the component is unmounted
    return () => {
      socket.close();
    };
  }, []);

  return (
    <div>
      <h1>Kafka Messages</h1>
      <ul>
        {messages.map((msg, index) => (
          <li key={index}>{msg}</li>
        ))}
      </ul>
    </div>
  );
};

export default WebSocketClient;
