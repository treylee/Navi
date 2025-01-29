import React, { useState, useEffect } from "react";
import api from "../services/api";
import { Message } from "../types/message";

const Chat: React.FC = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");

  const fetchMessages = async () => {
    try {
      const response = await api.get<Message[]>("/messages");
      setMessages(response.data);
    } catch (error) {
      console.error("Error fetching messages:", error);
    }
  };

  const sendMessage = async () => {
    if (!input) return;
    try {
      const response = await api.post<Message>("/messages", {
        text: input,
        sender: "User1", // Hardcoded sender for now
      });
      setMessages((prev) => [...prev, response.data]);
      setInput("");
    } catch (error) {
      console.error("Error sending message:", error);
    }
  };

  useEffect(() => {
    fetchMessages();
  }, []);

  return (
    <div className="chat-container">
      <h1>Messenger</h1>
      <div className="messages">
        {messages.map((msg) => (
          <div key={msg.id}>
            <strong>{msg.sender}: </strong>
            {msg.text}
          </div>
        ))}
      </div>
      <input
        type="text"
        value={input}
        onChange={(e) => setInput(e.target.value)}
        placeholder="Type a message"
      />
      <button onClick={sendMessage}>Send</button>
    </div>
  );
};

export default Chat;
