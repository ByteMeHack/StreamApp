import { Box, Button, CardBody, Heading, Input, Stack } from "@chakra-ui/react";
import { useState } from "react";

export default function Chat({ room_id }) {
  const [message, setMessage] = useState("");
  const socket = new WebSocket(`ws://bytemehack.ru/api/room/${room_id}`);
  socket.onmessage = (event) => {
    console.log(event.data);
  };
  return (
    <Box className="blackBlock" p={3}>
      <Card>
        <Heading size="sm" color="white" className="grayblock">
          Chat
        </Heading>
        <CardBody bgColor="gray">
          <Stack>Elements</Stack>
          <Stack direction="row" spacing={3}>
            <Input
              placeholder="Type message here..."
              onChange={(e) => setMessage(e.target.value)}
            />
            <Button
              onClick={() => {
                socket.send(
                  JSON.stringify({
                    user_id: 1,
                    message_type: 1,
                    contents: message,
                  })
                );
              }}
            >
              Send
            </Button>
          </Stack>
        </CardBody>
      </Card>
    </Box>
  );
}
