import Home from "./pages/Home";
import { createBrowserRouter } from "react-router-dom";
import Room from "./pages/Room";
import Layout from "./pages/Layout";
import Profile from "./pages/Profile";
import Protected from "./components/Protected";
import NotFound from "./pages/Error";
export const router = createBrowserRouter([
  {
    element: <Layout />,
    path: "/",
    children: [
      {
        element: <Home />,
        index: true,
      },
      {
        element: (
          <Protected>
            <Room />
          </Protected>
        ),
        path: "/rooms/:id",
      },
      {
        element: (
          <Protected>
            <Profile />
          </Protected>
        ),
        path: "/users/:id",
      },
      {
        element: <NotFound />,
        path: "*",
      },
    ],
  },
]);
