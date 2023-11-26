import { Box, Text } from "@chakra-ui/react";

export default function Message({ message }) {
  return (
    <Box gap={5} display="flex" flexDirection="column">
      <Text>{message.user_id}</Text>
      <Text>{message.contents}</Text>
      <Text>{message.timestamp}</Text>
    </Box>
  );
}
