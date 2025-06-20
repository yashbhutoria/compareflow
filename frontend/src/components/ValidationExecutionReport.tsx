import React from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  Box,
  Chip,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Alert,
  CircularProgress,
  Grid,
  Card,
  CardContent,
  Divider,
  IconButton,
} from '@mui/material';
import {
  Close as CloseIcon,
  CheckCircle as CheckCircleIcon,
  Error as ErrorIcon,
  Warning as WarningIcon,
  Info as InfoIcon,
} from '@mui/icons-material';
import { Validation } from '../types';

interface ValidationExecutionReportProps {
  open: boolean;
  onClose: () => void;
  validation: Validation | null;
  isRunning: boolean;
}

export default function ValidationExecutionReport({
  open,
  onClose,
  validation,
  isRunning,
}: ValidationExecutionReportProps) {
  if (!validation) return null;

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed':
        return <CheckCircleIcon color="success" />;
      case 'failed':
        return <ErrorIcon color="error" />;
      case 'running':
        return <CircularProgress size={20} />;
      default:
        return <InfoIcon color="info" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'success';
      case 'failed':
        return 'error';
      case 'running':
        return 'info';
      default:
        return 'default';
    }
  };

  const formatDateTime = (dateString?: string) => {
    if (!dateString) return 'N/A';
    return new Date(dateString).toLocaleString();
  };

  const results = validation.results || {};
  const summary = results.summary || {};
  const details = results.details || {};
  const errors = results.errors || [];

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="lg"
      fullWidth
      PaperProps={{
        sx: { minHeight: '70vh' }
      }}
    >
      <DialogTitle>
        <Box display="flex" alignItems="center" justifyContent="space-between">
          <Box display="flex" alignItems="center" gap={2}>
            {getStatusIcon(validation.status)}
            <Typography variant="h6">
              Validation Execution Report: {validation.name}
            </Typography>
          </Box>
          <IconButton onClick={onClose}>
            <CloseIcon />
          </IconButton>
        </Box>
      </DialogTitle>
      
      <DialogContent dividers>
        {isRunning ? (
          <Box display="flex" flexDirection="column" alignItems="center" justifyContent="center" py={5}>
            <CircularProgress size={60} />
            <Typography variant="h6" mt={3}>
              Running validation...
            </Typography>
            <Typography variant="body2" color="text.secondary" mt={1}>
              This may take a few moments depending on the data size
            </Typography>
          </Box>
        ) : (
          <>
            {/* Status and Basic Info */}
            <Box mb={3}>
              <Grid container spacing={2}>
                <Grid item xs={12} md={6}>
                  <Card>
                    <CardContent>
                      <Typography variant="subtitle2" color="text.secondary">
                        Status
                      </Typography>
                      <Box display="flex" alignItems="center" gap={1} mt={1}>
                        <Chip
                          label={validation.status.toUpperCase()}
                          color={getStatusColor(validation.status)}
                          size="small"
                        />
                        <Typography variant="body2">
                          Last run: {formatDateTime(validation.updated_at)}
                        </Typography>
                      </Box>
                    </CardContent>
                  </Card>
                </Grid>
                <Grid item xs={12} md={6}>
                  <Card>
                    <CardContent>
                      <Typography variant="subtitle2" color="text.secondary">
                        Configuration
                      </Typography>
                      <Typography variant="body2" mt={1}>
                        Type: <strong>{validation.config.comparison_type || 'N/A'}</strong>
                      </Typography>
                      <Typography variant="body2">
                        Source: <strong>{validation.source_connection?.name || 'N/A'}</strong>
                      </Typography>
                      <Typography variant="body2">
                        Target: <strong>{validation.target_connection?.name || 'N/A'}</strong>
                      </Typography>
                    </CardContent>
                  </Card>
                </Grid>
              </Grid>
            </Box>

            {/* Summary Statistics */}
            {validation.status !== 'pending' && (
              <>
                <Typography variant="h6" gutterBottom>
                  Summary Statistics
                </Typography>
                <Grid container spacing={2} mb={3}>
                  <Grid item xs={6} md={3}>
                    <Card>
                      <CardContent>
                        <Typography variant="subtitle2" color="text.secondary">
                          Source Rows
                        </Typography>
                        <Typography variant="h4">
                          {summary.source_row_count || 0}
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>
                  <Grid item xs={6} md={3}>
                    <Card>
                      <CardContent>
                        <Typography variant="subtitle2" color="text.secondary">
                          Target Rows
                        </Typography>
                        <Typography variant="h4">
                          {summary.target_row_count || 0}
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>
                  <Grid item xs={6} md={3}>
                    <Card>
                      <CardContent>
                        <Typography variant="subtitle2" color="text.secondary">
                          Matched Rows
                        </Typography>
                        <Typography variant="h4" color="success.main">
                          {summary.matched_rows || 0}
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>
                  <Grid item xs={6} md={3}>
                    <Card>
                      <CardContent>
                        <Typography variant="subtitle2" color="text.secondary">
                          Success Rate
                        </Typography>
                        <Typography variant="h4" color={summary.success_rate === 100 ? 'success.main' : 'warning.main'}>
                          {summary.success_rate || 0}%
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>
                </Grid>

                {/* Discrepancies */}
                {(summary.mismatched_rows > 0 || summary.missing_in_target > 0 || summary.extra_in_target > 0) && (
                  <>
                    <Typography variant="h6" gutterBottom>
                      Discrepancies Found
                    </Typography>
                    <Grid container spacing={2} mb={3}>
                      {summary.mismatched_rows > 0 && (
                        <Grid item xs={12} md={4}>
                          <Alert severity="warning">
                            <Typography variant="subtitle2">
                              Mismatched Rows: <strong>{summary.mismatched_rows}</strong>
                            </Typography>
                          </Alert>
                        </Grid>
                      )}
                      {summary.missing_in_target > 0 && (
                        <Grid item xs={12} md={4}>
                          <Alert severity="error">
                            <Typography variant="subtitle2">
                              Missing in Target: <strong>{summary.missing_in_target}</strong>
                            </Typography>
                          </Alert>
                        </Grid>
                      )}
                      {summary.extra_in_target > 0 && (
                        <Grid item xs={12} md={4}>
                          <Alert severity="info">
                            <Typography variant="subtitle2">
                              Extra in Target: <strong>{summary.extra_in_target}</strong>
                            </Typography>
                          </Alert>
                        </Grid>
                      )}
                    </Grid>
                  </>
                )}

                {/* Queries Used */}
                <Typography variant="h6" gutterBottom>
                  Queries Executed
                </Typography>
                <Grid container spacing={2} mb={3}>
                  <Grid item xs={12} md={6}>
                    <Card>
                      <CardContent>
                        <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                          Source Query
                        </Typography>
                        <Paper variant="outlined" sx={{ p: 2, bgcolor: 'grey.50' }}>
                          <Typography variant="body2" component="pre" sx={{ fontFamily: 'monospace', fontSize: '0.85rem' }}>
                            {validation.config.source_query || 'N/A'}
                          </Typography>
                        </Paper>
                      </CardContent>
                    </Card>
                  </Grid>
                  <Grid item xs={12} md={6}>
                    <Card>
                      <CardContent>
                        <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                          Target Query
                        </Typography>
                        <Paper variant="outlined" sx={{ p: 2, bgcolor: 'grey.50' }}>
                          <Typography variant="body2" component="pre" sx={{ fontFamily: 'monospace', fontSize: '0.85rem' }}>
                            {validation.config.target_query || 'N/A'}
                          </Typography>
                        </Paper>
                      </CardContent>
                    </Card>
                  </Grid>
                </Grid>

                {/* Errors */}
                {errors.length > 0 && (
                  <>
                    <Typography variant="h6" gutterBottom color="error">
                      Errors
                    </Typography>
                    <Box mb={3}>
                      {errors.map((error, index) => (
                        <Alert severity="error" key={index} sx={{ mb: 1 }}>
                          {typeof error === 'string' ? error : error.message || 'Unknown error'}
                        </Alert>
                      ))}
                    </Box>
                  </>
                )}

                {/* Detailed Differences (if available) */}
                {details.differences && details.differences.length > 0 && (
                  <>
                    <Typography variant="h6" gutterBottom>
                      Sample Differences (First 10)
                    </Typography>
                    <TableContainer component={Paper} sx={{ mb: 3 }}>
                      <Table size="small">
                        <TableHead>
                          <TableRow>
                            <TableCell>Type</TableCell>
                            <TableCell>Key</TableCell>
                            <TableCell>Source Value</TableCell>
                            <TableCell>Target Value</TableCell>
                            <TableCell>Columns</TableCell>
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          {details.differences.slice(0, 10).map((diff, index) => (
                            <TableRow key={index}>
                              <TableCell>
                                <Chip
                                  label={diff.type}
                                  size="small"
                                  color={diff.type === 'missing' ? 'error' : diff.type === 'extra' ? 'info' : 'warning'}
                                />
                              </TableCell>
                              <TableCell>{JSON.stringify(diff.key)}</TableCell>
                              <TableCell>{JSON.stringify(diff.source_data)}</TableCell>
                              <TableCell>{JSON.stringify(diff.target_data)}</TableCell>
                              <TableCell>{diff.columns?.join(', ') || 'N/A'}</TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </TableContainer>
                  </>
                )}
              </>
            )}
          </>
        )}
      </DialogContent>
      
      <DialogActions>
        <Button onClick={onClose}>Close</Button>
      </DialogActions>
    </Dialog>
  );
}