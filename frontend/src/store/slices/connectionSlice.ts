import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { connectionService } from '../../services/connectionService';
import { Connection } from '../../types';

interface ConnectionState {
  connections: Connection[];
  currentConnection: Connection | null;
  loading: boolean;
  error: string | null;
  testResult: { success: boolean; message: string } | null;
}

const initialState: ConnectionState = {
  connections: [],
  currentConnection: null,
  loading: false,
  error: null,
  testResult: null,
};

export const fetchConnections = createAsyncThunk('connections/fetchAll', async () => {
  const response = await connectionService.getConnections();
  return response;
});

export const fetchConnection = createAsyncThunk(
  'connections/fetchOne',
  async (id: number) => {
    const response = await connectionService.getConnection(id);
    return response;
  }
);

export const createConnection = createAsyncThunk(
  'connections/create',
  async (data: Omit<Connection, 'id'>) => {
    const response = await connectionService.createConnection(data);
    return response;
  }
);

export const updateConnection = createAsyncThunk(
  'connections/update',
  async ({ id, data }: { id: number; data: Omit<Connection, 'id'> }) => {
    const response = await connectionService.updateConnection(id, data);
    return response;
  }
);

export const deleteConnection = createAsyncThunk(
  'connections/delete',
  async (id: number) => {
    await connectionService.deleteConnection(id);
    return id;
  }
);

export const testConnection = createAsyncThunk(
  'connections/test',
  async (id: number) => {
    const response = await connectionService.testConnection(id);
    return response;
  }
);

const connectionSlice = createSlice({
  name: 'connections',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    clearTestResult: (state) => {
      state.testResult = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch all
      .addCase(fetchConnections.pending, (state) => {
        state.loading = true;
      })
      .addCase(fetchConnections.fulfilled, (state, action) => {
        state.loading = false;
        state.connections = action.payload;
      })
      .addCase(fetchConnections.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch connections';
      })
      // Fetch one
      .addCase(fetchConnection.fulfilled, (state, action) => {
        state.currentConnection = action.payload;
      })
      // Create
      .addCase(createConnection.fulfilled, (state, action) => {
        state.connections.push(action.payload);
      })
      // Update
      .addCase(updateConnection.fulfilled, (state, action) => {
        const index = state.connections.findIndex((c) => c.id === action.payload.id);
        if (index !== -1) {
          state.connections[index] = action.payload;
        }
      })
      // Delete
      .addCase(deleteConnection.fulfilled, (state, action) => {
        state.connections = state.connections.filter((c) => c.id !== action.payload);
      })
      // Test
      .addCase(testConnection.fulfilled, (state, action) => {
        state.testResult = action.payload;
      });
  },
});

export const { clearError, clearTestResult } = connectionSlice.actions;
export default connectionSlice.reducer;