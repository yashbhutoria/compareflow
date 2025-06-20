import api from './api';
import { Connection, ColumnInfo } from '../types';

export const connectionService = {
  async getConnections(): Promise<Connection[]> {
    const response = await api.get('/connections');
    return response.data;
  },

  async getConnection(id: number): Promise<Connection> {
    const response = await api.get(`/connections/${id}`);
    return response.data;
  },

  async createConnection(data: Omit<Connection, 'id'>): Promise<Connection> {
    const response = await api.post('/connections', data);
    return response.data;
  },

  async updateConnection(id: number, data: Omit<Connection, 'id'>): Promise<Connection> {
    const response = await api.put(`/connections/${id}`, data);
    return response.data;
  },

  async deleteConnection(id: number): Promise<void> {
    await api.delete(`/connections/${id}`);
  },

  async testConnection(id: number): Promise<{ success: boolean; message: string }> {
    const response = await api.post(`/connections/${id}/test`);
    return response.data;
  },

  async getTables(id: number): Promise<string[]> {
    const response = await api.get(`/connections/${id}/tables`);
    return response.data.tables;
  },

  async getColumns(id: number, tableName: string): Promise<ColumnInfo[]> {
    const response = await api.get(`/connections/${id}/tables/${encodeURIComponent(tableName)}/columns`);
    return response.data.columns;
  },
};