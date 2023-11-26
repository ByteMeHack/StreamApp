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
      borderRadius={10}
      p={3}
    >
      <Text color="black">User: {message.user_id}</Text>
      <Text color="black">Message: {message.contents}</Text>
      <Text color="black">Date:{message.timestamp}</Text>
    </Stack>
  );
}
