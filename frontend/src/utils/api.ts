import ky from 'ky';
import { auth, logout, setAuth } from '../data/auth';
import { isTokenExpired } from '../data/jwt';

export const api = ky.create({
  prefix: '/api',
  hooks: {
    beforeRequest: [
      async ({ request }: any) => { 
        let token = auth.token;

        // refresh logic
        if (token && isTokenExpired(token)) {
          try {
            const refreshResponse = await window.fetch('/api/refresh', {
              method: 'POST',
              credentials: 'include',
            });

            if (refreshResponse.ok) {
              const data = await refreshResponse.json();
              setAuth('token', data.token);
              token = data.token;
            } else {
              logout();
              return;
            }
          } catch {
            logout();
            return;
          }
        }

        if (token) {
          request.headers.set('Authorization', `Bearer ${token}`);
        }
      }
    ],
    afterResponse: [
      async ({ response }) => {
        if (response.status === 401 && !window.location.pathname.includes('/login')) {
      await logout(); 
    }
      }
    ]
  }
});