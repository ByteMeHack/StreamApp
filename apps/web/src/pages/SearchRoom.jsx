import { Box, Button, Heading, Input, Stack, Text } from "@chakra-ui/react";
import { useState } from "react";
import { getRoomByName } from "../api";
import { Link } from "react-router-dom";
import RoomsStack from "../components/RoomsStack";

export default function SearchRoom() {
  const [rooms, setRooms] = useState([]);
  const [name, setName] = useState("");
  return (
    <Stack
      flexGrow={1}
      width="90%"
      placeContent="center"
      alignSelf="center"
      alignItems="center"
      spacing={5}
    >
      <Stack direction="row" spacing={10} width="50%" placeContent="center">
        <Input
          placeholder="Type name of the room"
          onChange={(e) => setName(e.target.value)}
          isRequired
        />
        <Button
          onClick={() => {
            getRoomByName(name).then((res) => {
              setRooms(res);
            });
          }}
        >
          Search
        </Button>
      </Stack>

      {rooms.length > 0 ? (
        <RoomsStack rooms={rooms} />
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
          <Text fontSize="18px" mt={3} mb={2} color="gray.500">
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
