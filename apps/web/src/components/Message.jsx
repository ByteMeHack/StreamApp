import { Stack, Text } from "@chakra-ui/react";

export default function Message({ message }) {
  return (
    <Stack
      gap={2}
      display="flex"
      flexDirection="column"
      width={300}
      height={100}
      bgColor="gray"
    >
      <Text>User: {message.user_id}</Text>
      <Text>Message: {message.contents}</Text>
      <Text>Date:{message.timestamp}</Text>
    </Stack>
  );
}
