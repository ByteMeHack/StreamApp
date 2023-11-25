import { Card, Heading, Image, Spinner } from "@chakra-ui/react";
import { Link } from "react-router-dom";

export default function RoomCard({ id, name }) {
  return (
    <Link to={`/rooms/${id}`}>
      <Card
        border="1px"
        borderColor="#e02525"
        overflow="hidden"
        borderRadius="lg"
        width={300}
        height={336}
        alignItems="center"
      >
        <Heading
          padding={3}
          className="blackBlock"
          color="#e02525"
          height={36}
          whiteSpace="nowrap"
          width={300}
          size="md"
        >
          {name}
        </Heading>
        <Image
          bgColor="#e02525"
          fallback={
            <Spinner
              bgColor="#2f3235"
              thickness="4px"
              placeContent="center"
              speed="0.65s"
              color="red.500"
              size="xl"
            />
          }
          src={`https://cataas.com/cat/says/${name}?fontSize=25&type=square&height=300&width=300&fontColor=white`}
          width={300}
          height={300}
          objectFit="cover"
          objectPosition="center"
        />
      </Card>
    </Link>
  );
}
