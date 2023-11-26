import { Box, Button, Heading, Input, Stack } from "@chakra-ui/react";
import { useEffect, useRef, useState } from "react";
import Message from "./Message";

export default function Chat({ room_id }) {
  const [message, setMessage] = useState("");
  const [messages, setMessages] = useState([]);
  const [connection, setConnection] = useState(false);
  const [pending, setPending] = useState(false);
  const bottomRef = useRef(null);

  useEffect(() => {
    // 👇️ scroll to bottom every time messages change
    bottomRef.current?.scrollIntoView({ behavior: "smooth", block: "end" });
  }, [messages]);

  const socketRef = useRef(null);

  function sendMessage() {
    socketRef.current.send(
      JSON.stringify({
        message_type: 1,
        contents: message,
      })
    );
    setMessage("");
  }

  useEffect(() => {
    socketRef.current = new WebSocket(`ws://bytemehack.ru/api/room/${room_id}`);

    if (socketRef.current.readyState == WebSocket.CONNECTING) {
      setPending(true);
    }
    if (socketRef.current.readyState == WebSocket.OPEN) {
      setConnection(true);
    }

    const listenMessage = (event) => {
      setMessages((messages) => messages.concat(JSON.parse(event.data)));
    };
    socketRef.current.addEventListener("message", listenMessage);
    return function () {
      socketRef.current.removeEventListener("message", listenMessage);
      socketRef.current.close();
    };
  }, []);

  return (
    <Box className="blackBlock" p={3}>
      <Stack width="50%">
        <Heading size="sm" color="white" className="grayblock">
          Chat
        </Heading>
        <Stack overflow="scroll" maxHeight="30rem">
          {messages.map((message) => {
            return <Message key={message.contents} message={message} />;
          })}
          <div ref={bottomRef} />
        </Stack>
        {pending && (
          <Button
            onClick={() => {
              socketRef.current.send(
                JSON.stringify({
                  message_type: 0,
                  contents: "Entering room",
                })
              );
              setPending(false);
            }}
          >
            Connect
          </Button>
        )}
        {!pending && (
          <Stack direction="row" spacing={3}>
            <Input
              placeholder="Type message here..."
              onChange={(e) => setMessage(e.target.value)}
              value={message}
            />
            <Button isDisabled={connection} onClick={sendMessage}>
              Send
            </Button>
          </Stack>
        )}
      </Stack>
    </Box>
  );
}
