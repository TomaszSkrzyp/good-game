import { Component, JSX, onMount } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { auth } from "../data/store";

interface Props {
  children: JSX.Element;
}

const ProtectedRoute: Component<Props> = (props) => {
  const navigate = useNavigate();

  onMount(() => {
    if (!auth.isLoggedIn) {
      navigate("/", { replace: true });
    }
  });

  return <>{auth.isLoggedIn ? props.children : null}</>;
};

export default ProtectedRoute;