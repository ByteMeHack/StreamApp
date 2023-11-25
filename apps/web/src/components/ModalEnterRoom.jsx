import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { registerToRoom } from "../api";
import {
  Button,
  Input,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  Stack,
} from "@chakra-ui/react";

export default function ModalEnterRoom({ id, isOpen }) {
  const [pass, setPass] = useState("");
  const navigate = useNavigate();
  <Modal blockScrollOnMount={false} isOpen={true} className="blackBlock">
    <ModalContent>
      <ModalHeader className="blackBlock">
        Enter password to the room
      </ModalHeader>
      <Link to="/">
        <ModalCloseButton bgGradient="linear(to-r, red.400, red.500, red.600)" />
      </Link>
      <ModalBody className="blackBlock">
        <form
          id="my-form"
          onSubmit={async (e) => {
            e.preventDefault();
            try {
              const room = await registerToRoom(id, pass);
              navigate(`/rooms/${room.id}`);
            } catch (err) {
              console.log(err);
            }
          }}
        >
          <Stack gap={5} className="blackBlock">
            <Input
              autocomplete="off"
              placeholder="Name of room"
              type="text"
              isRequired
              onChange={(e) => setPass(e.target.value)}
              borderColor="#e02525"
            />
          </Stack>
        </form>
      </ModalBody>
      <ModalFooter className="blackBlock">
        <Link to="/">
          <Button
            color="#242424"
            colorScheme="red"
            bgGradient="linear(to-r, red.400, red.500, red.600)"
            mr={3}
          >
            Close
          </Button>
        </Link>
        <Button form="my-form" type="submit" color="#e02525" bgColor="white">
          Create
        </Button>
      </ModalFooter>
    </ModalContent>
  </Modal>;
}
