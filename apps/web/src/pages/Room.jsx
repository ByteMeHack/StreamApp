import { Card, CardBody, Box, Heading, Input } from "@chakra-ui/react";
import { useParams } from "react-router-dom";
import Chat from "../components/Chat";
import { useEffect, useState } from "react";
import { getRoomById } from "../api";
import ModalEnterRoom from "../components/ModalEnterRoom";

export default function Room() {
  const id = useParams().id;
  const [room, setRoom] = useState(null);
  const [needPass, setNeedPass] = useState(false);
  useEffect(() => {
    getRoomById(id)
      .then((res) => setRoom(res))
      .catch(() => {
        setNeedPass(true);
        console.log(needPass);
      });
  }, []);
  return (
    <Box className="grayBlock" display="flex" justifyContent="center">
      <ModalEnterRoom id={id} isOpen={needPass} />
      vvvvvv
      {room && (
        <Card
          w="90%"
          className="blackBlock"
          overflow="hidden"
          borderRadius="lg"
          border="4px"
          borderColor="#2f3235"
        >
          <Heading
            className="blackBlock"
            color="gray"
            size="lg"
            pt={3}
            pl={3}
            pr={3}
          >
            {room.name}
          </Heading>
          <CardBody className="blackBlock">
            <Input type="range" value={20} min={0} max={200} />
            <Chat />
          </CardBody>
        </Card>
      )}
    </Box>
  );
}
