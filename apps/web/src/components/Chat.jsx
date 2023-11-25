import {
  Box,
  Button,
  Card,
  CardBody,
  Heading,
  Input,
  Stack,
  Text,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";

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
      console.log("Message from server ", event.data);
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
              return (
                <Box display="flex">
                  <Text>User {message.user_id}</Text>
                  <Text>{message.contents}</Text>
                  <Text>{new Date(+message.timestamp)}</Text>
                </Box>
              );
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
