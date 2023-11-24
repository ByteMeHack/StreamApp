import {
  Button,
  Checkbox,
  Input,
  InputGroup,
  InputRightElement,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  Stack,
  useDisclosure,
  useToast,
} from "@chakra-ui/react";
import { useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { userSelector } from "../store/selectors";
import { showErrorRoomName, showSuccessRoomCreate } from "../utils/Toasts";
import { addRoom } from "../store/roomSlice";

export default function ModalCreateRoom() {
  const user = useSelector(userSelector);
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [private_room, setPrivate] = useState(false);
  const [name, setName] = useState("");
  const [pass, setPass] = useState("");
  const [show, setShow] = useState(false);
  const handleClick = () => setShow(!show);
  const toast = useToast();
  const dispatch = useDispatch();

  return (
    <>
      <Button
        mb={3}
        onClick={onOpen}
        size="lg"
        colorScheme="green"
        color="#242424"
        bgGradient="linear(to-r, green.400, green.500, green.600)"
      >
        Create new room
      </Button>
      <Modal
        blockScrollOnMount={false}
        isOpen={isOpen}
        onClose={() => {
          onClose();
          setPrivate(false);
        }}
        className="blackBlock"
      >
        <ModalContent>
          <ModalHeader className="blackBlock">Create your room</ModalHeader>
          <ModalCloseButton bgGradient="linear(to-r, green.400, green.500, green.600)" />
          <ModalBody className="blackBlock">
            <form
              id="my-form"
              onSubmit={async (e) => {
                e.preventDefault();
                const room = {
                  name: name,
                  private: private_room,
                  password: pass,
                  owner_id: user.id,
                };
                onClose();
                setPrivate(false);
                dispatch(addRoom(room))
                  .unwrap()
                  .then(() => {
                    toast(showSuccessRoomCreate);
                  })
                  .catch(() => {
                    toast(showErrorRoomName);
                  });
              }}
            >
              <Stack gap={5} className="blackBlock">
                <Input
                  placeholder="Name of room"
                  type="text"
                  isRequired
                  onChange={(e) => setName(e.target.value)}
                  borderColor="#55bd79"
                />
                <Checkbox onChange={(e) => setPrivate(e.target.checked)}>
                  Private room
                </Checkbox>
                {private_room && (
                  <InputGroup>
                    <Input
                      placeholder="Room password"
                      type={show ? "text" : "password"}
                      isRequired={private_room}
                      onChange={(e) => setPass(e.target.value)}
                      borderColor="#55bd79"
                    />
                    <InputRightElement width="4.5rem">
                      <Button
                        size="sm"
                        colorScheme="green"
                        color="black"
                        onClick={handleClick}
                      >
                        {show ? "Hide" : "Show"}
                      </Button>
                    </InputRightElement>
                  </InputGroup>
                )}
              </Stack>
            </form>
          </ModalBody>
          <ModalFooter className="blackBlock">
            <Button
              onClick={() => {
                onClose();
              }}
              color="#242424"
              colorScheme="red"
              mr={3}
            >
              Close
            </Button>
            <Button
              form="my-form"
              type="submit"
              color="#242424"
              colorScheme="green"
              bgGradient="linear(to-r, green.400, green.500, green.600)"
            >
              Create
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
}
