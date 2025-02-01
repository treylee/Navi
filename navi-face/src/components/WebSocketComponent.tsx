import React, { useState, useEffect } from 'react';

const WebSocketComponent: React.FC = () => {
  const [message, setMessage] = useState<string>(''); // State to hold the WebSocket message

  useEffect(() => {
    // Create the WebSocket connection
    const socket = new WebSocket("ws://localhost:8083/ws");

    // When WebSocket is connected
    socket.onopen = () => {
      console.log("Connected to WebSocket server");
    };

    // When a message is received
    socket.onmessage = (event: MessageEvent) => {
      console.log("Received message:", event.data);
      setMessage(event.data); // Update the message state with the received message
    };

    // Handle WebSocket errors
    socket.onerror = (error: Event) => {
      console.log("WebSocket error:", error);
    };

    // Handle WebSocket closure
    socket.onclose = () => {
      console.log("WebSocket connection closed");
    };

    // Cleanup when the component unmounts
    return () => {
      socket.close();
    };
  }, []); // Empty dependency array ensures it only runs once when the component mounts

  return (
    <div>
      <h1>WebSocket Messages</h1>
      <p>{message}</p>
    </div>
  );
}

export default WebSocketComponent;
