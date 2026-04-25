import { jwtDecode } from 'jwt-decode';
import { auth } from './auth';

export const getToken = () => auth.token;

export const isAdmin = () => {
  if (!auth.token) return false;
  try {
    const decoded: any = jwtDecode(auth.token);
    return decoded.role === 'Admin';
  } catch {
    return false;
  }
};

export const isTokenExpired = (token: string) => {
  try {
    const decoded: any = jwtDecode(token);
    const currentTime = Math.floor(Date.now() / 1000);
    // if token expires in less than 30 seconds, consider it expired
    return decoded.exp < currentTime + 30;
  } catch {
    return true;
  }
};