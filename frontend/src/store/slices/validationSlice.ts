import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { validationService } from '../../services/validationService';
import { Validation } from '../../types';

interface ValidationState {
  validations: Validation[];
  currentValidation: Validation | null;
  loading: boolean;
  error: string | null;
}

const initialState: ValidationState = {
  validations: [],
  currentValidation: null,
  loading: false,
  error: null,
};

export const fetchValidations = createAsyncThunk('validations/fetchAll', async () => {
  const response = await validationService.getValidations();
  return response;
});

export const fetchValidation = createAsyncThunk(
  'validations/fetchOne',
  async (id: number) => {
    const response = await validationService.getValidation(id);
    return response;
  }
);

export const createValidation = createAsyncThunk(
  'validations/create',
  async (data: Omit<Validation, 'id'>) => {
    const response = await validationService.createValidation(data);
    return response;
  }
);

export const updateValidation = createAsyncThunk(
  'validations/update',
  async ({ id, data }: { id: number; data: Omit<Validation, 'id'> }) => {
    const response = await validationService.updateValidation(id, data);
    return response;
  }
);

export const deleteValidation = createAsyncThunk(
  'validations/delete',
  async (id: number) => {
    await validationService.deleteValidation(id);
    return id;
  }
);

export const runValidation = createAsyncThunk(
  'validations/run',
  async (id: number) => {
    const response = await validationService.runValidation(id);
    return response;
  }
);

const validationSlice = createSlice({
  name: 'validations',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch all
      .addCase(fetchValidations.pending, (state) => {
        state.loading = true;
      })
      .addCase(fetchValidations.fulfilled, (state, action) => {
        state.loading = false;
        state.validations = action.payload;
      })
      .addCase(fetchValidations.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch validations';
      })
      // Fetch one
      .addCase(fetchValidation.fulfilled, (state, action) => {
        state.currentValidation = action.payload;
      })
      // Create
      .addCase(createValidation.fulfilled, (state, action) => {
        state.validations.push(action.payload);
      })
      // Update
      .addCase(updateValidation.fulfilled, (state, action) => {
        const index = state.validations.findIndex((v) => v.id === action.payload.id);
        if (index !== -1) {
          state.validations[index] = action.payload;
        }
      })
      // Delete
      .addCase(deleteValidation.fulfilled, (state, action) => {
        state.validations = state.validations.filter((v) => v.id !== action.payload);
      })
      // Run
      .addCase(runValidation.fulfilled, (state, action) => {
        const index = state.validations.findIndex((v) => v.id === action.payload.id);
        if (index !== -1) {
          state.validations[index] = action.payload;
        }
      });
  },
});

export const { clearError } = validationSlice.actions;
export default validationSlice.reducer;