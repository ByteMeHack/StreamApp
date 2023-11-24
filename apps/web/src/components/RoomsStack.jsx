import { Box, Button, Heading, Stack, Text } from "@chakra-ui/react";
import RoomCard from "./RoomCard";
import ModalCreateRoom from "./ModalCreateRoom";
import { useSelector } from "react-redux";
import { roomsSelector, userSelector } from "../store/selectors";
import { useState } from "react";

export default function RoomsStack() {
  const user = useSelector(userSelector);
  const allRooms = useSelector((state) => {
    return roomsSelector(state).slice((button - 1) * 6, button * 6);
  });
  const [button, setButton] = useState(1);
  let buttons = [],
    index = 1;
  for (let i = 0; i < allRooms.length; i += 6) {
    buttons.push(index);
    index++;
  }

  return (
    <Stack className="grayBlock" placeItems="center" spacing={15}>
      {user ? (
        <>
          <ModalCreateRoom />
          <Stack placeItems="center" mb={5}>
            <Stack
              direction="row"
              flexWrap="wrap"
              w="95%"
              gap={5}
              justifyContent="center"
            >
              {rooms.map((room) => {
                return <RoomCard key={room.id} id={room.id} name={room.name} />;
              })}
            </Stack>
          </Stack>
        </>
      ) : (
        <Box textAlign="center" className="grayBlock">
          <Heading
            display="inline-block"
            as="h2"
            size="2xl"
            bgGradient="linear(to-r, green.400, green.600)"
            backgroundClip="text"
          >
            Unauthorised
          </Heading>
          <Text fontSize="18px" color={"gray.500"} mt={4}>
            Login or Register to continue
          </Text>
        </Box>
      )}
      <Stack direction="row" spacing={5}>
        {buttons.map((button, index) => {
          return (
            <Button
              key={index}
              onClick={() => {
                setButton(index);
              }}
            >
              {index}
            </Button>
          );
        })}
      </Stack>
    </Stack>
  );
}
