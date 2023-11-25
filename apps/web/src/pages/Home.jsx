import { useEffect } from "react";
import RoomsStack from "../components/RoomsStack";
import { useDispatch, useSelector } from "react-redux";
import { roomsSelector, userSelector } from "../store/selectors";
import { setRooms } from "../store/roomSlice";
import { Stack } from "@chakra-ui/react";
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
      <ModalCreateRoom />
      <RoomsStack rooms={allRooms} />
    </Stack>
  );
}
