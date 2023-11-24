import { useEffect } from "react";
import RoomsStack from "../components/RoomsStack";
import { useDispatch, useSelector } from "react-redux";
import { userSelector } from "../store/selectors";
import { setRooms } from "../store/roomSlice";

export default function Home() {
  const dispatch = useDispatch();
  const user = useSelector(userSelector);
  useEffect(() => {
    if (user) {
      dispatch(setRooms(user.id));
    }
  }, [user]);
  return <RoomsStack />;
}
