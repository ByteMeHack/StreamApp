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
      <Button mb={3} onClick={onOpen} size="lg" bgColor="white" color="#e02525">
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
          <ModalCloseButton bgGradient="linear(to-r, red.400, red.500, red.600)" />
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
                  .catch((err) => {
                    if (err.status !== 500) toast(showErrorRoomName);
                  });
              }}
            >
              <Stack gap={5} className="blackBlock">
                <Input
                  autocomplete="off"
                  placeholder="Name of room"
                  type="text"
                  isRequired
                  onChange={(e) => setName(e.target.value)}
                  borderColor="#e02525"
                />
                <Checkbox onChange={(e) => setPrivate(e.target.checked)}>
                  Private room
                </Checkbox>
                {private_room && (
                  <InputGroup>
                    <Input
                      autocomplete="off"
                      placeholder="Room password"
                      type={show ? "text" : "password"}
                      isRequired={private_room}
                      onChange={(e) => setPass(e.target.value)}
                      borderColor="#e02525"
                    />
                    <InputRightElement width="4.5rem">
                      <Button
                        size="sm"
                        color="#e02525"
                        bgColor="white"
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
              bgGradient="linear(to-r, red.400, red.500, red.600)"
              mr={3}
            >
              Close
            </Button>
            <Button
              form="my-form"
              type="submit"
              color="#e02525"
              bgColor="white"
            >
              Create
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
}
