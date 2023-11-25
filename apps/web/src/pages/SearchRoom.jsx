import { Box, Button, Heading, Input, Stack, Text } from "@chakra-ui/react";
import { useState } from "react";
import { getRoomByName } from "../api";
import RoomCard from "../components/RoomCard";
import { Link } from "react-router-dom";

export default function SearchRoom() {
  const [room, setRoom] = useState(null);
  const [name, setName] = useState("");
  return (
    <Stack flexGrow={1} alignItems="center" width="90%" alignContent="center">
      <Stack direction="row" spacing={10} minWidth="100%">
        <Input
          placeholder="Type name of the room"
          onChange={() => setName(e.target.value)}
          isRequired
        />
        <Button
          onClick={() => {
            getRoomByName(name).then((res) => setRoom(res));
            console.log(room);
          }}
        >
          Search
        </Button>
      </Stack>

      {room && room.length > 0 ? (
        <Link to={`/rooms/${room.id}`}>
          <RoomCard id={room.id} name={room.name} />
        </Link>
      ) : (
        <Box textAlign="center" py={10} px={6} className="grayBlock">
          <Heading
            display="inline-block"
            as="h2"
            size="2xl"
            bgGradient="linear(to-r, red.400, red.600)"
            backgroundClip="text"
          >
            Hello
          </Heading>
          <Text fontSize="18px" mt={3} mb={2}>
            Here you can search for rooms you want to enter
          </Text>
          <Link to="/">
            <Button
              colorScheme="red"
              bgGradient="linear(to-r, red.400, red.500, red.600)"
              color="black"
              variant="solid"
            >
              Go to Home
            </Button>
          </Link>
        </Box>
      )}
    </Stack>
  );
}
