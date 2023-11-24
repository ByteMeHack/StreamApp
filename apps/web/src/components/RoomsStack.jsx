import { Box, Heading, Stack, Text } from "@chakra-ui/react";
import RoomCard from "./RoomCard";
import ModalCreateRoom from "./ModalCreateRoom";
import { useSelector } from "react-redux";
import { roomsSelector, userSelector } from "../store/selectors";

export default function RoomsStack() {
  const user = useSelector(userSelector);
  const rooms = useSelector(roomsSelector);
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
    </Stack>
  );
}
