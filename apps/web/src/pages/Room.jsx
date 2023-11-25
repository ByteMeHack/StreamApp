import {
  Card,
  CardBody,
  Box,
  Heading,
  Input,
  Stack,
  Button,
} from "@chakra-ui/react";
import { useParams } from "react-router-dom";
import Chat from "../components/Chat";
import { useEffect, useState } from "react";
import { getRoomById, registerToRoom } from "../api";

export default function Room() {
  const id = useParams().id;
  const [room, setRoom] = useState(null);
  const [needPass, setNeedPass] = useState(false);
  const [pass, setPass] = useState("");
  useEffect(() => {
    getRoomById(id)
      .then((res) => setRoom(res))
      .catch(() => {
        setNeedPass(true);
      });
  }, []);
  return (
    <Box className="grayBlock" display="flex" justifyContent="center">
      {needPass ? (
        <Stack direction="row" spacing={10} width={300} placeContent="center">
          <Input
            autocomplete="off"
            placeholder="Name of room"
            type="text"
            isRequired
            onChange={(e) => setPass(e.target.value)}
            borderColor="#e02525"
          />
          <Button
            onClick={async () => {
              await registerToRoom(id, pass).then((res) => {
                setRoom(res);
                setNeedPass(false);
              });
            }}
          >
            Enter
          </Button>
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
              <Input type="range" value={20} min={0} max={200} />
              <Chat />
            </CardBody>
          </Card>
        )
      )}
    </Box>
  );
}
Room.displayName = "Room";
