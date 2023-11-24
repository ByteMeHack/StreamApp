export function setUserLocal(user) {
  localStorage.setItem("user", JSON.stringify(user));
}

export function getUserLocal() {
  const user = localStorage.getItem("user");
  if (user) return JSON.parse(user);
  return null;
}
