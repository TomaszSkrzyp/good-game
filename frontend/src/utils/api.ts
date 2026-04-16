import ky from 'ky';
import { auth, logout } from '../data/store';

export const api = ky.create({
  prefix: '/api',
  hooks: {
    beforeRequest: [
      ({ request }) => { 
        if (auth.token) {
          request.headers.set('Authorization', `Bearer ${auth.token}`);
        }
      }
    ],
    afterResponse: [
      async ({ response }) => {
        if (response.status === 401 || response.status === 502) {
          logout(); 
        }
      }
    ]
  }
});