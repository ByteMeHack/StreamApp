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
    socketRef.current = new WebSocket(`ws://bytemehack.ru/api/room/${room_id}`);
    socketRef.current.addEventListener("message", (event) => {
      setMessages([...messages, event.data]);
    });
  }, []);

  return (
    <Box className="blackBlock" p={3}>
      <Card>
        <Heading size="sm" color="white" className="grayblock">
          Chat
        </Heading>
        <CardBody bgColor="gray">
          <Stack>
            {messages.map((message) => {
              return <Message message={message} />;
            })}
          </Stack>
          <Stack direction="row" spacing={3}>
            <Input
              placeholder="Type message here..."
              onChange={(e) => setMessage(e.target.value)}
            />
            <Button
              isDisabled={socketRef.current === null}
              onClick={sendMessage}
            >
              Send
            </Button>
          </Stack>
        </CardBody>
      </Card>
    </Box>
  );
}
