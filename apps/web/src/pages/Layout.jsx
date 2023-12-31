import { Button, Heading, Image, Stack } from "@chakra-ui/react";
import DrawerRegister from "../components/DrawerRegister";
import { Link, Outlet } from "react-router-dom";
import DrawerLogin from "../components/DrawLogin";
import { useDispatch, useSelector } from "react-redux";
import { userSelector } from "../store/selectors";
import { setUser } from "../store/userSlice";
import { useEffect } from "react";
import { getUserLocal, setUserLocal } from "../utils/localStorage";
import { getCurrUserReq, logoutReq } from "../api";

export default function Layout() {
  let user = useSelector(userSelector);
  const dispatch = useDispatch();
  useEffect(() => {
    if (!user) user = getUserLocal();
    getCurrUserReq().then((res) => {
      setUserLocal(res);
      dispatch(setUser(res));
    });
  }, []);
  return (
    <Stack className="grayBlock" flexGrow={1}>
      <Stack direction="row" p={5} justifyContent="space-between">
        <Link to="/">
          <Stack direction="row" spacing={5}>
            <Image src="/vite.svg" width={50} height={50} />
            <Heading className="grayBlock">Stream audio</Heading>
          </Stack>
        </Link>
        <Stack direction="row" spacing={5} alignItems="center">
          {user ? (
            <>
              <Link to="/rooms/search">
                <Button colorScheme="red" variant="link" size="lg">
                  Search room
                </Button>
              </Link>
              <Link to={`/users/${user.id}`}>
                <Button bgColor="white" color="#e02525" size="lg">
                  {user.name} profile
                </Button>
              </Link>
              <Button
                colorScheme="red"
                color="#242424"
                onClick={() => {
                  logoutReq();
                  setUserLocal(null);
                  dispatch(setUser(null));
                }}
                size="lg"
              >
                Logout
              </Button>
            </>
          ) : (
            <>
              <DrawerLogin />
              <DrawerRegister />
            </>
          )}
        </Stack>
      </Stack>
      <Stack placeContent="center" flexGrow={1}>
        <Outlet />
      </Stack>
    </Stack>
  );
}
