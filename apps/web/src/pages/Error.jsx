"use client";

import { Box, Heading, Text, Button } from "@chakra-ui/react";
import { Link } from "react-router-dom";

export default function NotFound() {
  return (
    <Box textAlign="center" py={10} px={6} className="grayBlock">
      <Heading
        display="inline-block"
        as="h2"
        size="2xl"
        bgGradient="linear(to-r, green.400, green.600)"
        backgroundClip="text"
      >
        404
      </Heading>
      <Text fontSize="18px" mt={3} mb={2}>
        Page Not Found
      </Text>
      <Text color={"gray.500"} mb={6}>
        The page you&apos;re looking for does not seem to exist
      </Text>

      <Link to="/">
        <Button
          colorScheme="green"
          bgGradient="linear(to-r, green.400, green.500, green.600)"
          color="black"
          variant="solid"
        >
          Go to Home
        </Button>
      </Link>
    </Box>
  );
}
