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
  Tooltip,
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  PlayArrow as RunIcon,
  Visibility as ViewIcon,
} from '@mui/icons-material';
import { fetchValidations, deleteValidation, runValidation, fetchValidation } from '../store/slices/validationSlice';
import { AppDispatch, RootState } from '../store';
import ValidationExecutionReport from '../components/ValidationExecutionReport';

export default function Validations() {
  const navigate = useNavigate();
  const dispatch = useDispatch<AppDispatch>();
  const { validations, currentValidation } = useSelector((state: RootState) => state.validations);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedValidation, setSelectedValidation] = useState<number | null>(null);
  const [reportOpen, setReportOpen] = useState(false);
  const [runningValidation, setRunningValidation] = useState<number | null>(null);

  useEffect(() => {
    dispatch(fetchValidations());
  }, [dispatch]);

  const handleDelete = async () => {
    if (selectedValidation) {
      await dispatch(deleteValidation(selectedValidation));
      setDeleteDialogOpen(false);
      setSelectedValidation(null);
    }
  };

  const handleRun = async (id: number) => {
    setRunningValidation(id);
    setSelectedValidation(id);
    setReportOpen(true);
    
    // Fetch the validation details first
    await dispatch(fetchValidation(id));
    
    // Run the validation
    const result = await dispatch(runValidation(id));
    
    // Refresh validations list
    await dispatch(fetchValidations());
    
    // Update current validation with results
    if (result.payload) {
      await dispatch(fetchValidation(id));
    }
    
    setRunningValidation(null);
  };
  
  const handleViewReport = async (id: number) => {
    setSelectedValidation(id);
    await dispatch(fetchValidation(id));
    setReportOpen(true);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'success';
      case 'running':
        return 'info';
      case 'failed':
        return 'error';
      default:
        return 'default';
    }
  };

  return (
    <Container maxWidth="lg">
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4">Validations</Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => navigate('/validations/new')}
        >
          New Validation
        </Button>
      </Box>

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Source Connection</TableCell>
              <TableCell>Target Connection</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Success Rate</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {validations.map((validation) => (
              <TableRow key={validation.id}>
                <TableCell>{validation.name}</TableCell>
                <TableCell>{validation.source_connection?.name || 'N/A'}</TableCell>
                <TableCell>{validation.target_connection?.name || 'N/A'}</TableCell>
                <TableCell>
                  <Chip
                    label={validation.status}
                    size="small"
                    color={getStatusColor(validation.status)}
                  />
                </TableCell>
                <TableCell>
                  {validation.results?.summary?.success_rate
                    ? `${validation.results.summary.success_rate}%`
                    : '-'}
                </TableCell>
                <TableCell>
                  <Tooltip title="Run Validation">
                    <IconButton
                      size="small"
                      onClick={() => handleRun(validation.id)}
                      disabled={validation.status === 'running' || runningValidation === validation.id}
                      color="primary"
                    >
                      <RunIcon />
                    </IconButton>
                  </Tooltip>
                  {validation.status !== 'pending' && (
                    <Tooltip title="View Report">
                      <IconButton
                        size="small"
                        onClick={() => handleViewReport(validation.id)}
                        color="info"
                      >
                        <ViewIcon />
                      </IconButton>
                    </Tooltip>
                  )}
                  <Tooltip title="Edit">
                    <IconButton
                      size="small"
                      onClick={() => navigate(`/validations/${validation.id}/edit`)}
                    >
                      <EditIcon />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="Delete">
                    <IconButton
                      size="small"
                      onClick={() => {
                        setSelectedValidation(validation.id);
                        setDeleteDialogOpen(true);
                      }}
                      color="error"
                    >
                      <DeleteIcon />
                    </IconButton>
                  </Tooltip>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
        <DialogTitle>Delete Validation</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to delete this validation? This action cannot be undone.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleDelete} color="error">
            Delete
          </Button>
        </DialogActions>
      </Dialog>
      
      {/* Validation Execution Report Modal */}
      <ValidationExecutionReport
        open={reportOpen}
        onClose={() => {
          setReportOpen(false);
          setSelectedValidation(null);
          setRunningValidation(null);
        }}
        validation={currentValidation}
        isRunning={runningValidation !== null}
      />
    </Container>
  );
}