import { Heading, Spinner, Stack } from "@chakra-ui/react";
import { useParams } from "react-router-dom";
import { getUserByIdReq } from "../api";
import { useEffect } from "react";
import { useState } from "react";

export default function Profile() {
  const id = useParams().id;
  const [user, setUser] = useState();
  useEffect(() => {
    getUserByIdReq(id).then((res) => setUser(res));
  }, []);
  return (
    <Stack
      direction="column"
      gap={2}
      placeItems="center"
      className="grayBlock"
    >
      {user ? (
        <>
          <Heading className="grayBlock">{user.name}</Heading>
          <Heading className="grayBlock">{user.email}</Heading>
        </>
      ) : (
        <Spinner
          thickness="4px"
          placeContent="center"
          speed="0.65s"
          color="green.500"
          size="xl"
        />
      )}
    </Stack>
  );
}
