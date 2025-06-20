import { Routes, Route, Navigate } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from './store';
import MainLayout from './components/Layout/MainLayout';
import PrivateRoute from './components/PrivateRoute';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import Connections from './pages/Connections';
import ConnectionForm from './pages/ConnectionForm';
import Validations from './pages/Validations';
import ValidationForm from './pages/ValidationForm';

function App() {
  const isAuthenticated = useSelector((state: RootState) => state.auth.isAuthenticated);

  return (
    <Routes>
      <Route path="/login" element={!isAuthenticated ? <Login /> : <Navigate to="/" />} />
      <Route path="/register" element={!isAuthenticated ? <Register /> : <Navigate to="/" />} />
      <Route
        path="/"
        element={
          <PrivateRoute>
            <MainLayout />
          </PrivateRoute>
        }
      >
        <Route index element={<Dashboard />} />
        <Route path="connections" element={<Connections />} />
        <Route path="connections/new" element={<ConnectionForm />} />
        <Route path="connections/:id/edit" element={<ConnectionForm />} />
        <Route path="validations" element={<Validations />} />
        <Route path="validations/new" element={<ValidationForm />} />
        <Route path="validations/:id/edit" element={<ValidationForm />} />
      </Route>
    </Routes>
  );
}

export default App;