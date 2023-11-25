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
import { useEffect, useRef, useState } from "react";

export default function Chat({ room_id }) {
  const [message, setMessage] = useState("");
  const socketRef = useRef(null);
  const messages = [];

  function sendMessage() {
    if (socketRef.current)
      socketRef.current.send({
        message_type: 1,
        contents: message,
      });
  }

  useEffect(() => {
    socketRef.current = new WebSocket(`ws://bytemehack.ru/api/room/${room_id}`);
    socketRef.current.onmessage = (event) => {
      messages.push(event.data);
    };
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
