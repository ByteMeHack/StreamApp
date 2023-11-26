import {
  Card,
  CardBody,
  Box,
  Heading,
  Input,
  Stack,
  Button,
  InputGroup,
  InputRightElement,
  useToast,
} from "@chakra-ui/react";
import { useParams } from "react-router-dom";
import Chat from "../components/Chat";
import { useEffect, useState } from "react";
import { getRoomById, registerToRoom } from "../api";
import { showError } from "../utils/Toasts";

export default function Room() {
  const id = useParams().id;
  const toast = useToast();
  const [room, setRoom] = useState(null);
  const [needPass, setNeedPass] = useState(false);
  const [pass, setPass] = useState("");
  const [show, setShow] = useState(false);
  const handleClick = () => setShow(!show);
  useEffect(() => {
    getRoomById(id)
      .then((res) => setRoom(res))
      .catch(() => {
        setNeedPass(true);
      });
  }, []);
  return (
    <Box
      className="grayBlock"
      display="flex"
      justifyContent="center"
      alignItems="center"
    >
      {needPass ? (
        <Stack spacing={5}>
          <Heading color="red" size="lg">
            This room is private. Write a password to enter
          </Heading>
          <Stack
            direction="row"
            spacing={10}
            width={300}
            placeContent="center"
            placeItems="center"
          >
            <InputGroup>
              <Input
                placeholder="Type your password"
                type={show ? "text" : "password"}
                required
                onChange={(e) => setPass(e.target.value)}
                borderColor="#e02525"
              />
              <InputRightElement width="4.5rem">
                <Button
                  size="sm"
                  color="#e02525"
                  bgColor="white"
                  onClick={handleClick}
                >
                  {show ? "Hide" : "Show"}
                </Button>
              </InputRightElement>
            </InputGroup>
            <Button
              onClick={async () => {
                await registerToRoom(id, pass)
                  .then((res) => {
                    setRoom(res);
                    setNeedPass(false);
                  })
                  .catch((err) =>
                    toast(showError("Can't enter the room. Check the password"))
                  );
              }}
            >
              Enter
            </Button>
          </Stack>
        </Stack>
      ) : (
        room && (
          <Card
            w="90%"
            className="blackBlock"
            overflow="hidden"
            borderRadius="lg"
            border="4px"
            borderColor="#2f3235"
          >
            <Heading
              className="blackBlock"
              color="gray"
              size="lg"
              pt={3}
              pl={3}
              pr={3}
            >
              {room.name}
            </Heading>
            <CardBody className="blackBlock">
              <Chat room_id={room.id} />
            </CardBody>
          </Card>
        )
      )}
    </Box>
  );
}
Room.displayName = "Room";
