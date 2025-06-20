import { configureStore } from '@reduxjs/toolkit';
import authReducer from './slices/authSlice';
import connectionReducer from './slices/connectionSlice';
import validationReducer from './slices/validationSlice';

export const store = configureStore({
  reducer: {
    auth: authReducer,
    connections: connectionReducer,
    validations: validationReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;