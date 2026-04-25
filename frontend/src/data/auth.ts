import { createStore } from "solid-js/store";
import { api } from "../utils/api";

// only the data structure
interface AuthState {
  token: string;
  userName: string;
  email: string;
  role: string; // added role here for your admin check
  isLoggedIn: boolean;
}

const getInitialState = (): AuthState => {
  const saved = localStorage.getItem("auth_data");
  return saved ? JSON.parse(saved) : { token: "", userName: "", email: "", role: "", isLoggedIn: false };
};

export const [auth, setAuth] = createStore<AuthState>(getInitialState());

export const loginUser = (data: any) => {
  const state = { ...data, isLoggedIn: true };
  setAuth(state);
  localStorage.setItem("auth_data", JSON.stringify(state));
};

export const logout = async (navigate?: (path: string) => void) => {
  try {
    await api.post('logout');
  } catch (err) {
    console.error("Logout request failed:", err);
  } finally {
    setAuth({ 
      token: "", 
      userName: "", 
      email: "", 
      role: "", 
      isLoggedIn: false 
    });

    localStorage.removeItem("auth_data");

    if (window.location.pathname.startsWith('/profile')) {
      window.location.href = '/';
    } else {
      window.location.reload();
    }
  }
};