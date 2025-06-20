import api from './api';
import { Validation } from '../types';

export const validationService = {
  async getValidations(): Promise<Validation[]> {
    const response = await api.get('/validations');
    return response.data;
  },

  async getValidation(id: number): Promise<Validation> {
    const response = await api.get(`/validations/${id}`);
    return response.data;
  },

  async createValidation(data: Omit<Validation, 'id'>): Promise<Validation> {
    const response = await api.post('/validations', data);
    return response.data;
  },

  async updateValidation(id: number, data: Omit<Validation, 'id'>): Promise<Validation> {
    const response = await api.put(`/validations/${id}`, data);
    return response.data;
  },

  async deleteValidation(id: number): Promise<void> {
    await api.delete(`/validations/${id}`);
  },

  async runValidation(id: number): Promise<Validation> {
    const response = await api.post(`/validations/${id}/run`);
    return response.data;
  },

  async getValidationStatus(id: number): Promise<{ id: number; status: string; updated_at: string }> {
    const response = await api.get(`/validations/${id}/status`);
    return response.data;
  },
};