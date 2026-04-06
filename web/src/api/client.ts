import axios from 'axios';
import { store } from '../store';
import { logout, setTokens } from '../store/slices/authSlice';

const client = axios.create({
  baseURL: '/api/v1',
  headers: { 'Content-Type': 'application/json' },
});

client.interceptors.request.use((config) => {
  const token = store.getState().auth.accessToken;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

client.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      const refreshToken = store.getState().auth.refreshToken;
      if (refreshToken) {
        try {
          const res = await axios.post('/api/v1/auth/refresh', { refresh_token: refreshToken });
          store.dispatch(setTokens(res.data));
          originalRequest.headers.Authorization = `Bearer ${res.data.access_token}`;
          return client(originalRequest);
        } catch {
          store.dispatch(logout());
        }
      }
    }
    return Promise.reject(error);
  }
);

export default client;
