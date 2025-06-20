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
  CircularProgress,
  Alert,
} from '@mui/material';
import { createValidation, updateValidation, fetchValidation } from '../store/slices/validationSlice';
import { fetchConnections } from '../store/slices/connectionSlice';
import { AppDispatch, RootState } from '../store';

export default function ValidationForm() {
  const navigate = useNavigate();
  const { id } = useParams();
  const dispatch = useDispatch<AppDispatch>();
  const { currentValidation, loading } = useSelector((state: RootState) => state.validations);
  const { connections } = useSelector((state: RootState) => state.connections);
  
  const [formData, setFormData] = useState({
    name: '',
    source_connection_id: 0,
    target_connection_id: 0,
    config: {
      source_query: '',
      target_query: '',
      comparison_type: 'row_count' as 'row_count' | 'data_match' | 'schema',
      key_columns: [] as string[],
    },
  });

  const [error, setError] = useState('');

  useEffect(() => {
    dispatch(fetchConnections());
    if (id) {
      dispatch(fetchValidation(parseInt(id)));
    }
  }, [id, dispatch]);

  useEffect(() => {
    if (currentValidation && id) {
      setFormData({
        name: currentValidation.name,
        source_connection_id: currentValidation.source_connection_id,
        target_connection_id: currentValidation.target_connection_id,
        config: {
          source_query: currentValidation.config.source_query || '',
          target_query: currentValidation.config.target_query || '',
          comparison_type: currentValidation.config.comparison_type || 'row_count',
          key_columns: currentValidation.config.key_columns || [],
        },
      });
    }
  }, [currentValidation, id]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement> | any) => {
    const { name, value } = e.target;
    
    if (name === 'name' || name === 'source_connection_id' || name === 'target_connection_id') {
      setFormData({ 
        ...formData, 
        [name]: name.includes('connection_id') ? parseInt(value) : value 
      });
    } else if (name in formData.config) {
      setFormData({
        ...formData,
        config: {
          ...formData.config,
          [name]: value,
        },
      });
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    try {
      const data = {
        ...formData,
        status: 'pending' as const,
      };
      
      if (id) {
        await dispatch(updateValidation({ id: parseInt(id), data })).unwrap();
      } else {
        await dispatch(createValidation(data)).unwrap();
      }
      navigate('/validations');
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
          {id ? 'Edit Validation' : 'New Validation'}
        </Typography>
        
        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
        
        <Box component="form" onSubmit={handleSubmit}>
          <TextField
            fullWidth
            label="Validation Name"
            name="name"
            value={formData.name}
            onChange={handleChange}
            margin="normal"
            required
          />
          
          <FormControl fullWidth margin="normal" required>
            <InputLabel>Source Connection</InputLabel>
            <Select
              name="source_connection_id"
              value={formData.source_connection_id || ''}
              onChange={handleChange}
              label="Source Connection"
            >
              <MenuItem value="">Select a connection</MenuItem>
              {connections.map((conn) => (
                <MenuItem key={conn.id} value={conn.id}>
                  {conn.name} ({conn.type})
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <FormControl fullWidth margin="normal" required>
            <InputLabel>Target Connection</InputLabel>
            <Select
              name="target_connection_id"
              value={formData.target_connection_id || ''}
              onChange={handleChange}
              label="Target Connection"
            >
              <MenuItem value="">Select a connection</MenuItem>
              {connections.map((conn) => (
                <MenuItem key={conn.id} value={conn.id}>
                  {conn.name} ({conn.type})
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <FormControl fullWidth margin="normal">
            <InputLabel>Comparison Type</InputLabel>
            <Select
              name="comparison_type"
              value={formData.config.comparison_type}
              onChange={handleChange}
              label="Comparison Type"
            >
              <MenuItem value="row_count">Row Count</MenuItem>
              <MenuItem value="data_match">Data Match</MenuItem>
              <MenuItem value="schema">Schema Comparison</MenuItem>
            </Select>
          </FormControl>

          <TextField
            fullWidth
            label="Source Query"
            name="source_query"
            value={formData.config.source_query}
            onChange={handleChange}
            margin="normal"
            multiline
            rows={4}
            placeholder="SELECT * FROM table_name"
          />

          <TextField
            fullWidth
            label="Target Query"
            name="target_query"
            value={formData.config.target_query}
            onChange={handleChange}
            margin="normal"
            multiline
            rows={4}
            placeholder="SELECT * FROM table_name"
          />

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
              onClick={() => navigate('/validations')}
            >
              Cancel
            </Button>
          </Box>
        </Box>
      </Paper>
    </Container>
  );
}