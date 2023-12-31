import { Button, Stack } from "@chakra-ui/react";
import RoomCard from "./RoomCard";
import { useState } from "react";

export default function RoomsStack({ rooms }) {
  const [currButton, setButton] = useState(1);
  let buttons = [],
    index = 1;
  for (let i = 0; i < rooms.length; i += 6) {
    buttons.push(index);
    index++;
  }

  return (
    <Stack className="grayBlock" placeItems="center" spacing={15}>
      <Stack placeItems="center" mb={5}>
        <Stack
          direction="row"
          flexWrap="wrap"
          w="95%"
          gap={5}
          justifyContent="center"
        >
          {rooms.slice((currButton - 1) * 6, currButton * 6).map((room) => {
            return <RoomCard key={room.id} id={room.id} name={room.name} />;
          })}
        </Stack>
      </Stack>
      <Stack direction="row" spacing={5}>
        {buttons.map((button, index) => {
          return (
            <Button
              isDisabled={currButton === button}
              colorScheme="red"
              bgGradient="linear(to-r, red.400, red.500, red.600)"
              key={button}
              onClick={() => {
                setButton(button);
              }}
            >
              {button}
            </Button>
          );
        })}
      </Stack>
    </Stack>
  );
}
