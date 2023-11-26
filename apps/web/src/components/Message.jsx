import { Stack, Text } from "@chakra-ui/react";

export default function Message({ message }) {
  return (
    <Stack
      gap={5}
      display="flex"
      flexDirection="column"
      width={300}
      height={100}
    >
      <Text>{message.user_id}</Text>
      <Text>{message.contents}</Text>
      <Text>{message.timestamp}</Text>
    </Stack>
  );
}
