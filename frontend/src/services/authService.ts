import api from './api';
import { User } from '../types';

interface LoginResponse {
  access_token: string;
  user: User;
}

interface RegisterResponse {
  access_token: string;
  user: User;
}

export const authService = {
  async login(username: string, password: string): Promise<LoginResponse> {
    const response = await api.post('/auth/login', { username, password });
    return response.data;
  },

  async register(username: string, email: string, password: string): Promise<RegisterResponse> {
    const response = await api.post('/auth/register', { username, email, password });
    return response.data;
  },

  async getCurrentUser(): Promise<User> {
    const response = await api.get('/auth/me');
    return response.data;
  },

  logout() {
    localStorage.removeItem('token');
  },
};