import {
  Button,
  Drawer,
  DrawerBody,
  DrawerCloseButton,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerOverlay,
  useDisclosure,
  Input,
  Stack,
  InputRightElement,
  InputGroup,
  useToast,
} from "@chakra-ui/react";
import { useRef, useState } from "react";
import { loginReq } from "../api";
import { useDispatch } from "react-redux";
import { setUser } from "../store/userSlice";
import { showErrorLogOpts, showSuccessLogOpts } from "../utils/Toasts";
import { setUserLocal } from "../utils/localStorage";

async function loginUser(email, pass) {
  return await loginReq(email, pass);
}

export default function DrawerLogin() {
  const dispatch = useDispatch();
  const [show, setShow] = useState(false);
  const handleClick = () => setShow(!show);
  const { isOpen, onOpen, onClose } = useDisclosure();
  const btnRef = useRef();
  const [email, setEmail] = useState("");
  const [pass, setPass] = useState("");
  const toast = useToast();

  return (
    <>
      <Button
        ref={btnRef}
        color="#c23838"
        bgColor="white"
        onClick={onOpen}
        size="lg"
      >
        Login
      </Button>
      <Drawer
        isOpen={isOpen}
        placement="right"
        onClose={onClose}
        finalFocusRef={btnRef}
      >
        <DrawerOverlay />
        <DrawerContent>
          <DrawerCloseButton textColor="white" />
          <DrawerHeader className="blackBlock" fontSize="x-large">
            Login to your account
          </DrawerHeader>
          <DrawerBody className="grayBlock">
            <form
              id="my-form"
              onSubmit={async (e) => {
                e.preventDefault();
                let user;
                try {
                  user = await loginUser(email, pass);
                  dispatch(setUser(user));
                  setUserLocal(user);
                  toast(showSuccessLogOpts);
                } catch {
                  toast(showErrorLogOpts);
                }
              }}
            >
              <Stack gap={5}>
                <Input
                  placeholder="Type your email"
                  type="email"
                  required
                  onChange={(e) => setEmail(e.target.value)}
                  borderColor="#c23838"
                />
                <InputGroup>
                  <Input
                    placeholder="Type your password"
                    type={show ? "text" : "password"}
                    required
                    onChange={(e) => setPass(e.target.value)}
                    borderColor="#c23838"
                  />
                  <InputRightElement width="4.5rem">
                    <Button
                      size="sm"
                      color="#c23838"
                      bgColor="white"
                      onClick={handleClick}
                    >
                      {show ? "Hide" : "Show"}
                    </Button>
                  </InputRightElement>
                </InputGroup>
              </Stack>
            </form>
          </DrawerBody>

          <DrawerFooter className="blackBlock">
            <Button
              colorScheme="red"
              textColor="black"
              mr={3}
              onClick={onClose}
            >
              Cancel
            </Button>
            <Button
              form="my-form"
              type="submit"
              bgColor="#c23838"
              textColor="black"
            >
              Login
            </Button>
          </DrawerFooter>
        </DrawerContent>
      </Drawer>
    </>
  );
}
