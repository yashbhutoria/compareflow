import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import {
  Container,
  Typography,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Chip,
  Box,
  Alert,
  Snackbar,
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  PlayArrow as TestIcon,
} from '@mui/icons-material';
import { 
  fetchConnections, 
  deleteConnection, 
  testConnection,
  clearTestResult 
} from '../store/slices/connectionSlice';
import { AppDispatch, RootState } from '../store';

export default function Connections() {
  const navigate = useNavigate();
  const dispatch = useDispatch<AppDispatch>();
  const { connections, testResult } = useSelector((state: RootState) => state.connections);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedConnection, setSelectedConnection] = useState<number | null>(null);
  const [snackbarOpen, setSnackbarOpen] = useState(false);

  useEffect(() => {
    dispatch(fetchConnections());
  }, [dispatch]);

  useEffect(() => {
    if (testResult) {
      setSnackbarOpen(true);
    }
  }, [testResult]);

  const handleDelete = async () => {
    if (selectedConnection) {
      await dispatch(deleteConnection(selectedConnection));
      setDeleteDialogOpen(false);
      setSelectedConnection(null);
    }
  };

  const handleTest = async (id: number) => {
    await dispatch(testConnection(id));
  };

  const handleSnackbarClose = () => {
    setSnackbarOpen(false);
    dispatch(clearTestResult());
  };

  return (
    <Container maxWidth="lg">
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4">Connections</Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => navigate('/connections/new')}
        >
          Add Connection
        </Button>
      </Box>

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Type</TableCell>
              <TableCell>Server/Workspace</TableCell>
              <TableCell>Database</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {connections.map((connection) => (
              <TableRow key={connection.id}>
                <TableCell>{connection.name}</TableCell>
                <TableCell>
                  <Chip
                    label={connection.type}
                    size="small"
                    color={connection.type === 'sqlserver' ? 'primary' : 'secondary'}
                  />
                </TableCell>
                <TableCell>
                  {connection.config.server || connection.config.workspace || '-'}
                </TableCell>
                <TableCell>{connection.config.database || '-'}</TableCell>
                <TableCell>
                  <IconButton
                    size="small"
                    onClick={() => handleTest(connection.id)}
                    title="Test Connection"
                  >
                    <TestIcon />
                  </IconButton>
                  <IconButton
                    size="small"
                    onClick={() => navigate(`/connections/${connection.id}/edit`)}
                    title="Edit"
                  >
                    <EditIcon />
                  </IconButton>
                  <IconButton
                    size="small"
                    onClick={() => {
                      setSelectedConnection(connection.id);
                      setDeleteDialogOpen(true);
                    }}
                    title="Delete"
                  >
                    <DeleteIcon />
                  </IconButton>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
        <DialogTitle>Delete Connection</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to delete this connection? This action cannot be undone.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleDelete} color="error">
            Delete
          </Button>
        </DialogActions>
      </Dialog>

      {/* Test Result Snackbar */}
      <Snackbar
        open={snackbarOpen}
        autoHideDuration={6000}
        onClose={handleSnackbarClose}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert
          onClose={handleSnackbarClose}
          severity={testResult?.success ? 'success' : 'error'}
          sx={{ width: '100%' }}
        >
          {testResult?.message}
        </Alert>
      </Snackbar>
    </Container>
  );
}