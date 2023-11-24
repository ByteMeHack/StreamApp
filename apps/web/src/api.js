import axios from "axios";

const instance = axios.create({
  baseURL: import.meta.env.VITE_BACKEND_URL,
  withCredentials: true,
});

export async function registerReq(name, email, password) {
  return instance
    .post("/signup", {
      name: name,
      email: email,
      password: password,
    })
    .then(function (response) {
      return response.data;
    });
}

export async function loginReq(email, password) {
  return instance
    .post("/login", {
      email: email,
      password: password,
    })
    .then(function (response) {
      return response.data;
    });
}

export async function getCurrUserReq() {
  return instance.get("/me").then(function (response) {
    return response.data;
  });
}

export async function getUserByIdReq(id) {
  return instance.get(`/users/${id}`).then(function (response) {
    return response.data;
  });
}

export async function getUserRooms(id) {
  return instance.get(`/users/${id}/rooms`).then(function (response) {
    return response.data;
  });
}

export async function createRoom(name, private_room, pass, owner_id) {
  return instance
    .post("/rooms", {
      name: name,
      private: private_room,
      password: pass,
      owner_id: owner_id,
    })
    .then(function (response) {
      return response.data;
    });
}

export async function logoutReq() {
  await instance.post(`/logout`);
}

export async function getRoomById(id) {
  return instance.get(`/rooms/${id}`).then(function (response) {
    return response.data;
  });
}
