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

export default function ModalEnterRoom({ id }) {
  const [pass, setPass] = useState("");
  const navigate = useNavigate();
  <Stack direction="row" spacing={10} width="50%" placeContent="center">
    <Input
      autocomplete="off"
      placeholder="Name of room"
      type="text"
      isRequired
      onChange={(e) => setPass(e.target.value)}
      borderColor="#e02525"
    />
    <Button
      onClick={async () => {
        await registerToRoom(id, pass).then((res) => {
          navigate(`/rooms/${res.id}`);
        });
      }}
    >
      Enter
    </Button>
  </Stack>;
}
