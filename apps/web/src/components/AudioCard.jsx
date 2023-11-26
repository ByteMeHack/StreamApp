import { Card, Heading } from "@chakra-ui/react";

export default function AudioCard() {
  return (
    <Card>
      <Heading>Playing: Audio Name</Heading>

      <CardBody>Search for Audio</CardBody>
      <Select placeholder="Select option">
        <option value="option1">Option 1</option>
        <option value="option2">Option 2</option>
        <option value="option3">Option 3</option>
      </Select>
    </Card>
  );
}
