import ky from 'ky';
import { auth } from '../data/store';
export const api = ky.create({
  prefix: '/api',
  hooks: {
    beforeRequest: [
      ({ request }) => { 
        request.headers.set('Authorization', `Bearer ${auth.token}`);
      }
  ]
  }
});