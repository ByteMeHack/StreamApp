import { useEffect } from "react";
import RoomsStack from "../components/RoomsStack";
import { useDispatch, useSelector } from "react-redux";
import { roomsSelector, userSelector } from "../store/selectors";
import { setRooms } from "../store/roomSlice";
import { Box, Heading, Stack, Text } from "@chakra-ui/react";
import ModalCreateRoom from "../components/ModalCreateRoom";

export default function Home() {
  const dispatch = useDispatch();
  const user = useSelector(userSelector);
  const allRooms = useSelector(roomsSelector);
  useEffect(() => {
    if (user) {
      dispatch(setRooms(user.id));
    }
  }, [user]);
  return (
    <Stack placeItems="center" spacing={15}>
      {user ? (
        <>
          <ModalCreateRoom />
          <RoomsStack rooms={allRooms} />
        </>
      ) : (
        <Box textAlign="center" className="grayBlock">
          <Heading
            display="inline-block"
            as="h2"
            size="2xl"
            bgGradient="linear(to-r, red.400, red.600)"
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
