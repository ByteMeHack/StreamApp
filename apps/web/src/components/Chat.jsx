import {
  Box,
  Button,
  Card,
  CardBody,
  Heading,
  Input,
  Stack,
} from "@chakra-ui/react";
import { useEffect, useRef, useState } from "react";
import Message from "./Message";

export default function Chat({ room_id }) {
  const [message, setMessage] = useState("");
  const [messages, setMessages] = useState([]);
  const socketRef = useRef(null);

  function sendMessage() {
    if (socketRef.current)
      socketRef.current.send(
        JSON.stringify({
          message_type: 1,
          contents: message,
        })
      );
  }

  useEffect(() => {
    let arrayMessages = [];
    socketRef.current = new WebSocket(`ws://bytemehack.ru/api/room/${room_id}`);
    socketRef.current.addEventListener("message", (event) => {
      arrayMessages.push(JSON.parse(event.data));
      setMessages(arrayMessages);
    });
  }, []);

  return (
    <Box className="blackBlock" p={3}>
      <Stack>
        <Heading size="sm" color="white" className="grayblock">
          Chat
        </Heading>
        <Stack overflow="scroll">
          {messages.map((message) => {
            return <Message key={message.id} message={message} />;
          })}
        </Stack>
        <Stack direction="row" spacing={3}>
          <Input
            placeholder="Type message here..."
            onChange={(e) => setMessage(e.target.value)}
          />
          <Button isDisabled={socketRef.current === null} onClick={sendMessage}>
            Send
          </Button>
        </Stack>
      </Stack>
    </Box>
  );
}
