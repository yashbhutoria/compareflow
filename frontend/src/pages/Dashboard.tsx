import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import {
  Container,
  Grid,
  Paper,
  Typography,
  Box,
  Card,
  CardContent,
  CardActions,
  Button,
} from '@mui/material';
import {
  Storage as StorageIcon,
  CheckCircle as CheckCircleIcon,
  PlayArrow as PlayArrowIcon,
  Add as AddIcon,
} from '@mui/icons-material';
import { fetchConnections } from '../store/slices/connectionSlice';
import { fetchValidations } from '../store/slices/validationSlice';
import { AppDispatch, RootState } from '../store';

export default function Dashboard() {
  const navigate = useNavigate();
  const dispatch = useDispatch<AppDispatch>();
  const { connections } = useSelector((state: RootState) => state.connections);
  const { validations } = useSelector((state: RootState) => state.validations);

  useEffect(() => {
    dispatch(fetchConnections());
    dispatch(fetchValidations());
  }, [dispatch]);

  const stats = [
    {
      title: 'Total Connections',
      value: connections.length,
      icon: <StorageIcon fontSize="large" />,
      color: '#1976d2',
    },
    {
      title: 'Total Validations',
      value: validations.length,
      icon: <CheckCircleIcon fontSize="large" />,
      color: '#388e3c',
    },
    {
      title: 'Running Validations',
      value: validations.filter(v => v.status === 'running').length,
      icon: <PlayArrowIcon fontSize="large" />,
      color: '#f57c00',
    },
  ];

  return (
    <Container maxWidth="lg">
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>
      
      <Grid container spacing={3} sx={{ mb: 4 }}>
        {stats.map((stat, index) => (
          <Grid item xs={12} sm={6} md={4} key={index}>
            <Paper
              sx={{
                p: 3,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
              }}
            >
              <Box>
                <Typography color="textSecondary" gutterBottom>
                  {stat.title}
                </Typography>
                <Typography variant="h3">
                  {stat.value}
                </Typography>
              </Box>
              <Box sx={{ color: stat.color }}>
                {stat.icon}
              </Box>
            </Paper>
          </Grid>
        ))}
      </Grid>

      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Quick Actions
              </Typography>
              <Box sx={{ mt: 2 }}>
                <Button
                  variant="outlined"
                  startIcon={<AddIcon />}
                  fullWidth
                  sx={{ mb: 2 }}
                  onClick={() => navigate('/connections/new')}
                >
                  Add New Connection
                </Button>
                <Button
                  variant="outlined"
                  startIcon={<AddIcon />}
                  fullWidth
                  onClick={() => navigate('/validations/new')}
                >
                  Create New Validation
                </Button>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Recent Validations
              </Typography>
              {validations.slice(0, 5).map((validation) => (
                <Box
                  key={validation.id}
                  sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    py: 1,
                    borderBottom: '1px solid #e0e0e0',
                    '&:last-child': { borderBottom: 'none' },
                  }}
                >
                  <Typography variant="body2">{validation.name}</Typography>
                  <Typography
                    variant="caption"
                    sx={{
                      px: 1,
                      py: 0.5,
                      borderRadius: 1,
                      bgcolor: validation.status === 'completed' ? '#e8f5e9' : '#fff3e0',
                      color: validation.status === 'completed' ? '#2e7d32' : '#e65100',
                    }}
                  >
                    {validation.status}
                  </Typography>
                </Box>
              ))}
            </CardContent>
            <CardActions>
              <Button size="small" onClick={() => navigate('/validations')}>
                View All
              </Button>
            </CardActions>
          </Card>
        </Grid>
      </Grid>
    </Container>
  );
}