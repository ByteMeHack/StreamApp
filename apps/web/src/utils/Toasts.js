const basicToast = {
  duration: 3000,
  isClosable: true,
  position: "bottom-left",
};

export const showErrorLogOpts = {
  title: "Account wasn't found.",
  description: "Check for errors in your info",
  status: "error",
  ...basicToast,
};

export const showSuccessLogOpts = {
  title: "Success",
  description: "You are logged in",
  status: "success",
  ...basicToast,
};

export function showErrorRegOpts(desc) {
  return {
    title: "Error creating account",
    description: desc,
    status: "error",
    ...basicToast,
  };
}

export const showErrorRoomName = {
  title: "Error creating room",
  description: "Room with this name already exists",
  status: "error",
  ...basicToast,
};

export const showSuccessRoomCreate = {
  title: "Success",
  description: "New room was created!",
  status: "success",
  ...basicToast,
};
