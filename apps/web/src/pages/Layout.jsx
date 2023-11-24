import { Button, Heading, Stack, VStack } from "@chakra-ui/react";
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
          <Heading className="grayBlock">Stream audio</Heading>
        </Link>
        <Stack direction="row" gap={3}>
          {user ? (
            <>
              <Link to={`/users/${user.id}`}>
                <Button
                  colorScheme="green"
                  color="#242424"
                  bgGradient="linear(to-r, green.400, green.500, green.600)"
                  size="lg"
                >
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
