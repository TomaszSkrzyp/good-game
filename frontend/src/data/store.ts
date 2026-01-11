import { createStore } from "solid-js/store";

const getInitialState = () => {
  const saved = localStorage.getItem("auth_data");
  return saved ? JSON.parse(saved) : { token: "", userName: "", email: "", isLoggedIn: false };
};

export const [auth, setAuth] = createStore(getInitialState());

export const loginUser = (data: any) => {
  const state = { ...data, isLoggedIn: true };
  setAuth(state);
  localStorage.setItem("auth_data", JSON.stringify(state));
};

export const clearSession = () => {
  setAuth({ 
    token: "", 
    userName: "", 
    email: "", 
    isLoggedIn: false 
  });
  
  localStorage.removeItem("auth_data");
};