import { Component } from "solid-js";
import Header from "./Header";

const Layout: Component<{ children?: any }> = (props) => {
  return (
    <>
      <Header />
      <main class="p-4">{props.children}</main>
    </>
  );
};

export default Layout;