import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import {
  Container,
  Paper,
  TextField,
  Button,
  Typography,
  Box,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  FormControlLabel,
  Checkbox,
  CircularProgress,
  Alert,
} from '@mui/material';
import { createConnection, updateConnection, fetchConnection } from '../store/slices/connectionSlice';
import { AppDispatch, RootState } from '../store';

export default function ConnectionForm() {
  const navigate = useNavigate();
  const { id } = useParams();
  const dispatch = useDispatch<AppDispatch>();
  const { currentConnection, loading } = useSelector((state: RootState) => state.connections);
  
  const [formData, setFormData] = useState({
    name: '',
    type: 'sqlserver' as 'sqlserver' | 'databricks',
    config: {
      server: '',
      port: 1433,
      database: '',
      username: '',
      password: '',
      encrypt: false,
      trust_server_certificate: true,
      workspace: '',
      http_path: '',
      access_token: '',
    },
  });

  const [error, setError] = useState('');

  useEffect(() => {
    if (id) {
      dispatch(fetchConnection(parseInt(id)));
    }
  }, [id, dispatch]);

  useEffect(() => {
    if (currentConnection && id) {
      setFormData({
        name: currentConnection.name,
        type: currentConnection.type,
        config: {
          ...formData.config,
          ...currentConnection.config,
        },
      });
    }
  }, [currentConnection, id]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement> | any) => {
    const { name, value, checked } = e.target;
    
    if (name === 'name' || name === 'type') {
      setFormData({ ...formData, [name]: value });
    } else if (name in formData.config) {
      setFormData({
        ...formData,
        config: {
          ...formData.config,
          [name]: name === 'port' ? parseInt(value) || 0 : 
                   (name === 'encrypt' || name === 'trust_server_certificate') ? checked : value,
        },
      });
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    try {
      if (id) {
        await dispatch(updateConnection({ id: parseInt(id), data: formData })).unwrap();
      } else {
        await dispatch(createConnection(formData)).unwrap();
      }
      navigate('/connections');
    } catch (err: any) {
      setError(err.message || 'An error occurred');
    }
  };

  if (loading) {
    return (
      <Container maxWidth="md">
        <Box display="flex" justifyContent="center" mt={4}>
          <CircularProgress />
        </Box>
      </Container>
    );
  }

  return (
    <Container maxWidth="md">
      <Paper sx={{ p: 4 }}>
        <Typography variant="h5" gutterBottom>
          {id ? 'Edit Connection' : 'New Connection'}
        </Typography>
        
        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
        
        <Box component="form" onSubmit={handleSubmit}>
          <TextField
            fullWidth
            label="Connection Name"
            name="name"
            value={formData.name}
            onChange={handleChange}
            margin="normal"
            required
          />
          
          <FormControl fullWidth margin="normal">
            <InputLabel>Connection Type</InputLabel>
            <Select
              name="type"
              value={formData.type}
              onChange={handleChange}
              label="Connection Type"
              required
            >
              <MenuItem value="sqlserver">SQL Server</MenuItem>
              <MenuItem value="databricks">Databricks</MenuItem>
            </Select>
          </FormControl>

          {formData.type === 'sqlserver' && (
            <>
              <TextField
                fullWidth
                label="Server"
                name="server"
                value={formData.config.server}
                onChange={handleChange}
                margin="normal"
                required
              />
              <TextField
                fullWidth
                label="Port"
                name="port"
                type="number"
                value={formData.config.port}
                onChange={handleChange}
                margin="normal"
                required
              />
              <TextField
                fullWidth
                label="Database"
                name="database"
                value={formData.config.database}
                onChange={handleChange}
                margin="normal"
                required
              />
              <TextField
                fullWidth
                label="Username"
                name="username"
                value={formData.config.username}
                onChange={handleChange}
                margin="normal"
                required
              />
              <TextField
                fullWidth
                label="Password"
                name="password"
                type="password"
                value={formData.config.password}
                onChange={handleChange}
                margin="normal"
                required
              />
              <FormControlLabel
                control={
                  <Checkbox
                    name="encrypt"
                    checked={formData.config.encrypt}
                    onChange={handleChange}
                  />
                }
                label="Encrypt Connection"
              />
              <FormControlLabel
                control={
                  <Checkbox
                    name="trust_server_certificate"
                    checked={formData.config.trust_server_certificate}
                    onChange={handleChange}
                  />
                }
                label="Trust Server Certificate"
              />
            </>
          )}

          {formData.type === 'databricks' && (
            <>
              <TextField
                fullWidth
                label="Workspace URL"
                name="workspace"
                value={formData.config.workspace}
                onChange={handleChange}
                margin="normal"
                required
                placeholder="https://your-workspace.databricks.com"
              />
              <TextField
                fullWidth
                label="HTTP Path"
                name="http_path"
                value={formData.config.http_path}
                onChange={handleChange}
                margin="normal"
                required
                placeholder="/sql/1.0/endpoints/your-endpoint"
              />
              <TextField
                fullWidth
                label="Access Token"
                name="access_token"
                type="password"
                value={formData.config.access_token}
                onChange={handleChange}
                margin="normal"
                required
              />
            </>
          )}

          <Box sx={{ mt: 3, display: 'flex', gap: 2 }}>
            <Button
              variant="contained"
              type="submit"
              disabled={loading}
            >
              {id ? 'Update' : 'Create'}
            </Button>
            <Button
              variant="outlined"
              onClick={() => navigate('/connections')}
            >
              Cancel
            </Button>
          </Box>
        </Box>
      </Paper>
    </Container>
  );
}