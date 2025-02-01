import React, { useState, useEffect, useRef } from "react";
import api from "../services/api";
import { Message } from "../types/message";
import "../chat.css"; // Ensure this file contains the updated styles

const Chat: React.FC = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [userName, setUserName] = useState<string>("");

  const messagesEndRef = useRef<HTMLDivElement | null>(null);
  const socket = useRef<WebSocket | null>(null);

  const generateRandomName = () => {
    const randomNumber = Math.floor(Math.random() * 1000) + 1;
    return `User${randomNumber}`;
  };

  const getUserName = () => {
    const storedName = localStorage.getItem("chatUserName");
    if (storedName) {
      return storedName;
    } else {
      const newUserName = generateRandomName();
      localStorage.setItem("chatUserName", newUserName); // Store in localStorage for the session
      return newUserName;
    }
  };

  const fetchMessages = async () => {
    try {
      const response = await api.get<Message[]>("/messages");
      console.log("Fetched messages from API:", response.data);
      setMessages(response.data);
    } catch (error) {
      console.error("Error fetching messages:", error);
    }
  };

  const sendMessage = async () => {
    if (!input) return;

    const sender = userName;
    const timestamp = Date.now(); // Get current timestamp
    console.log("Sending message:", input);

    try {
      const response = await api.post<Message>("/messages", {
        text: input,
        sender: sender,
        timestamp: timestamp, // Include timestamp in the sent message
      });
      console.log("Message sent:", response.data);
      setMessages((prev) => [...prev, response.data]); // Add to local state immediately
      setInput("");
    } catch (error) {
      console.error("Error sending message:", error);
    }

    // Send message over WebSocket as well
    if (socket.current) {
      const message = {
        text: input,
        sender: sender,
        timestamp: timestamp, // Include timestamp for consistency
      };
      socket.current.send(JSON.stringify(message));
    }
  };

  const handleKeyPress = (event: React.KeyboardEvent) => {
    if (event.key === "Enter") {
      sendMessage();
    }
  };

  useEffect(() => {
    const name = getUserName(); // Get or generate the user name
    setUserName(name); // Set the user name state
    console.log("User name set:", name);

    fetchMessages(); // Fetch initial messages

    // Initialize WebSocket connection
    socket.current = new WebSocket("ws://localhost:8083/ws");

    socket.current.onopen = () => {
      console.log("Connected to WebSocket server");
    };

    socket.current.onmessage = (event: MessageEvent) => {
      const messageData = JSON.parse(event.data); // Parse the incoming WebSocket message
      console.log("Received WebSocket message:", messageData);

      // Avoid adding duplicate messages using the 'id' as the unique identifier
      setMessages((prevMessages) => {
        const isMessageDuplicate = prevMessages.some((msg) => msg.id === messageData.id);
        if (isMessageDuplicate) {
          console.log("Duplicate message received, not adding to state.");
          return prevMessages; // Don't add the duplicate message
        }

        // If not a duplicate, add the message to the state
        return [
          ...prevMessages,
          {
            id: messageData.id, // Use the unique 'id' from server
            text: messageData.text,
            sender: messageData.sender,
            timestamp: Date.now(), // Add timestamp when receiving WebSocket message
          },
        ];
      });
    };

    socket.current.onerror = (error: Event) => {
      console.log("WebSocket error:", error);
    };

    socket.current.onclose = () => {
      console.log("WebSocket connection closed");
    };

    // Cleanup WebSocket when component unmounts
    return () => {
      if (socket.current) {
        socket.current.close();
      }
      console.log("WebSocket closed on unmount");
    };
  }, []); // Empty dependency array ensures it only runs once when the component mounts

  useEffect(() => {
    if (messagesEndRef.current) {
      console.log("Scrolling to the bottom of the chat");
      messagesEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages]);

  return (
    <div className="chat-container">
      <div className="chat-header">
        <h1>Navi</h1>
      </div>

      <div className="messages-container">
        <div className="messages">
          {messages.map((msg) => (
            <div
              key={msg.id}  // Use 'id' as a unique key
              className={`message ${msg.sender !== userName ? "incoming" : ""}`}
            >
              <strong>{msg.sender}: </strong>
              <span>{msg.text}</span>
            </div>
          ))}
        </div>
        <div ref={messagesEndRef} />
      </div>

      <div className="input-container">
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Type a message"
          className="message-input"
          onKeyPress={handleKeyPress} // Send message on Enter key press
        />
        <button onClick={sendMessage} className="send-button">
          Send
        </button>
      </div>
    </div>
  );
};

export default Chat;
