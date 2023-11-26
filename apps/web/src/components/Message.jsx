import { Stack, Text } from "@chakra-ui/react";
import { useSelector } from "react-redux";
import { userSelector } from "../store/selectors";

export default function Message({ message }) {
  const user = useSelector(userSelector);
  const isCurrUser = message.user_id === user.id;
  return (
    <Stack
      gap={2}
      display="flex"
      flexDirection="column"
      width={300}
      height={100}
      bgColor={isCurrUser ? "red" : "gray"}
      alignSelf={isCurrUser ? "flex-end" : ""}
      borderRadius={10}
      p={3}
    >
      <Text color="black">User: {message.user_id}</Text>
      <Text color="black">Message: {message.contents}</Text>
      <Text color="black">Date: {message.timestamp}</Text>
    </Stack>
  );
}
