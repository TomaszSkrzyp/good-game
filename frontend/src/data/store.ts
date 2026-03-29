import { createStore } from "solid-js/store";
import { jwtDecode } from 'jwt-decode'; 

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

export const getToken = () => auth.token;

export const setToken = (token: string) => {
  setAuth('token', token);
  // Update localStorage
  const current = getInitialState();
  current.token = token;
  localStorage.setItem("auth_data", JSON.stringify(current));
};

export const logout = () => {
  clearSession();
  window.location.href = '/login';
};

export const fetchWithAuth = async (url: string, options: RequestInit = {}): Promise<Response> => {
  let token = getToken();
  if (!token) {
    logout();
    throw new Error('No token available');
  }

  // Check if token expires in less than 1 minute
  const decoded: any = jwtDecode(token);
  if (decoded.exp * 1000 < Date.now() + 60000) {
    try {
      const refreshResponse = await fetch('/api/refresh', {
        method: 'POST',
        credentials: 'include',
      });
      if (refreshResponse.ok) {
        const data = await refreshResponse.json();
        setToken(data.token);
        token = data.token;
      } else {
        logout();
        throw new Error('Session expired');
      }
    } catch (error) {
      logout();
      throw new Error('Session expired');
    }
  }

  const response = await fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      Authorization: `Bearer ${token}`,
    },
  });

  if (response.status === 401) {
    // Fallback: if still 401, logout
    logout();
    throw new Error('Session expired');
  }

  return response;
};