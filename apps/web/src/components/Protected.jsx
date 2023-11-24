import { useSelector } from "react-redux";
import { userSelector } from "../store/selectors";
import { Navigate } from "react-router-dom";

export default function Protected({ children }) {
  const user = useSelector(userSelector);
  if (user) return children;
  return <Navigate to="/" />;
}
