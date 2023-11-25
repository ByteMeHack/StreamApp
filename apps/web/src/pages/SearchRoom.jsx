import { Button, Input, Stack } from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { getRoomByName } from "../api";
import RoomCard from "../components/RoomCard";
import { Link } from "react-router-dom";

export default function SearchRoom() {
  const [room, setRoom] = useState(null);
  const [name, setName] = useState("");
  return (
    <Stack flexGrow={1} justifyContent="space-evenly">
      <Input
        placeholder="Type name of the room"
        onChange={() => setName(e.target.value)}
        isRequired
      />
      <Button
        onClick={() => {
          getRoomByName(name).then((res) => setRoom(res));
        }}
      >
        Search
      </Button>
      {room && (
        <Link to={`/rooms/${room.id}`}>
          <RoomCard id={room.id} name={room.name} />
        </Link>
      )}
    </Stack>
  );
}
