# CompareFlow API Reference

## Overview

The CompareFlow API is a RESTful API that provides programmatic access to data validation and comparison functionality. All API endpoints return JSON responses and use standard HTTP response codes.

### Base URL
```
https://api.compareflow.com/api/v1
```

### Authentication
Most endpoints require authentication using JWT tokens. Include the token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

### Rate Limiting
- 1000 requests per hour for authenticated users
- 100 requests per hour for unauthenticated endpoints

### Response Format
```json
{
    "success": true,
    "data": {...},
    "error": null,
    "meta": {
        "page": 1,
        "per_page": 20,
        "total_count": 100
    }
}
```

## Authentication Endpoints

### Register User
Create a new user account.

**Endpoint:** `POST /auth/register`

**Request Body:**
```json
{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword123"
}
```

**Response:**
```json
{
    "success": true,
    "data": {
        "user": {
            "id": 1,
            "username": "johndoe",
            "email": "john@example.com",
            "created_at": "2024-01-15T10:00:00Z"
        },
        "access_token": "eyJhbGciOiJIUzI1NiIs..."
    }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid input data
- `409 Conflict` - Username or email already exists

---

### Login
Authenticate and receive access token.

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
    "username": "johndoe",
    "password": "securepassword123"
}
```

**Response:**
```json
{
    "success": true,
    "data": {
        "user": {
            "id": 1,
            "username": "johndoe",
            "email": "john@example.com"
        },
        "access_token": "eyJhbGciOiJIUzI1NiIs..."
    }
}
```

**Error Responses:**
- `401 Unauthorized` - Invalid credentials
- `429 Too Many Requests` - Rate limit exceeded

---

### Get Current User
Get authenticated user's profile.

**Endpoint:** `GET /auth/me`

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
    "success": true,
    "data": {
        "id": 1,
        "username": "johndoe",
        "email": "john@example.com",
        "created_at": "2024-01-15T10:00:00Z",
        "updated_at": "2024-01-15T10:00:00Z"
    }
}
```

## Connection Endpoints

### List Connections
Get all connections for the authenticated user.

**Endpoint:** `GET /connections`

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `per_page` (integer, optional): Items per page (default: 20, max: 100)
- `type` (string, optional): Filter by connection type (`sqlserver`, `databricks`)

**Response:**
```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "name": "Production SQL Server",
            "type": "sqlserver",
            "created_at": "2024-01-15T10:00:00Z",
            "updated_at": "2024-01-15T10:00:00Z"
        },
        {
            "id": 2,
            "name": "Analytics Databricks",
            "type": "databricks",
            "created_at": "2024-01-16T10:00:00Z",
            "updated_at": "2024-01-16T10:00:00Z"
        }
    ],
    "meta": {
        "page": 1,
        "per_page": 20,
        "total_count": 2
    }
}
```

---

### Get Connection
Get a specific connection by ID.

**Endpoint:** `GET /connections/{id}`

**Response:**
```json
{
    "success": true,
    "data": {
        "id": 1,
        "name": "Production SQL Server",
        "type": "sqlserver",
        "config": {
            "server": "prod-sql.example.com",
            "port": 1433,
            "database": "sales_db",
            "username": "readonly_user",
            "encrypt": true,
            "trust_server_certificate": false
        },
        "created_at": "2024-01-15T10:00:00Z",
        "updated_at": "2024-01-15T10:00:00Z"
    }
}
```

**Note:** Password fields are never returned in responses.

---

### Create Connection
Create a new database connection.

**Endpoint:** `POST /connections`

**Request Body (SQL Server):**
```json
{
    "name": "Production SQL Server",
    "type": "sqlserver",
    "config": {
        "server": "prod-sql.example.com",
        "port": 1433,
        "database": "sales_db",
        "username": "readonly_user",
        "password": "securepassword",
        "encrypt": true,
        "trust_server_certificate": false
    }
}
```

**Request Body (Databricks):**
```json
{
    "name": "Analytics Databricks",
    "type": "databricks",
    "config": {
        "workspace": "https://myworkspace.cloud.databricks.com",
        "http_path": "/sql/1.0/endpoints/abc123def456",
        "access_token": "dapi1234567890abcdef"
    }
}
```

**Response:**
```json
{
    "success": true,
    "data": {
        "id": 3,
        "name": "Production SQL Server",
        "type": "sqlserver",
        "config": {
            "server": "prod-sql.example.com",
            "port": 1433,
            "database": "sales_db",
            "username": "readonly_user",
            "encrypt": true,
            "trust_server_certificate": false
        },
        "created_at": "2024-01-17T10:00:00Z",
        "updated_at": "2024-01-17T10:00:00Z"
    }
}
```

---

### Update Connection
Update an existing connection.

**Endpoint:** `PUT /connections/{id}`

**Request Body:**
```json
{
    "name": "Production SQL Server (Updated)",
    "config": {
        "server": "new-prod-sql.example.com",
        "port": 1433,
        "database": "sales_db",
        "username": "readonly_user",
        "password": "newpassword",
        "encrypt": true,
        "trust_server_certificate": false
    }
}
```

**Note:** Only include fields you want to update. Password is required if changing any config fields.

---

### Delete Connection
Delete a connection.

**Endpoint:** `DELETE /connections/{id}`

**Response:**
```json
{
    "success": true,
    "data": {
        "message": "Connection deleted successfully"
    }
}
```

**Error Responses:**
- `400 Bad Request` - Connection is used by active validations
- `404 Not Found` - Connection not found

---

### Test Connection
Test database connectivity.

**Endpoint:** `POST /connections/{id}/test`

**Response:**
```json
{
    "success": true,
    "data": {
        "success": true,
        "message": "Connection successful",
        "details": {
            "server_version": "Microsoft SQL Server 2019",
            "database": "sales_db",
            "response_time_ms": 45
        }
    }
}
```

**Error Response:**
```json
{
    "success": false,
    "error": {
        "code": "CONNECTION_FAILED",
        "message": "Failed to connect to database",
        "details": "Login failed for user 'readonly_user'"
    }
}
```

---

### Get Tables
Get list of tables from a connection.

**Endpoint:** `GET /connections/{id}/tables`

**Query Parameters:**
- `schema` (string, optional): Filter by schema name

**Response:**
```json
{
    "success": true,
    "data": [
        {
            "schema": "dbo",
            "name": "orders",
            "type": "TABLE",
            "row_count": 1500000
        },
        {
            "schema": "dbo",
            "name": "customers",
            "type": "TABLE",
            "row_count": 50000
        }
    ]
}
```

---

### Get Table Columns
Get column information for a specific table.

**Endpoint:** `GET /connections/{id}/tables/{table}/columns`

**Response:**
```json
{
    "success": true,
    "data": [
        {
            "name": "order_id",
            "data_type": "int",
            "is_nullable": false,
            "is_primary_key": true,
            "max_length": null
        },
        {
            "name": "customer_id",
            "data_type": "int",
            "is_nullable": false,
            "is_primary_key": false,
            "max_length": null
        },
        {
            "name": "order_date",
            "data_type": "datetime",
            "is_nullable": false,
            "is_primary_key": false,
            "max_length": null
        },
        {
            "name": "status",
            "data_type": "varchar",
            "is_nullable": true,
            "is_primary_key": false,
            "max_length": 50
        }
    ]
}
```

## Validation Endpoints

### List Validations
Get all validations for the authenticated user.

**Endpoint:** `GET /validations`

**Query Parameters:**
- `page` (integer, optional): Page number
- `per_page` (integer, optional): Items per page
- `status` (string, optional): Filter by status (`pending`, `running`, `completed`, `failed`)
- `source_connection_id` (integer, optional): Filter by source connection
- `target_connection_id` (integer, optional): Filter by target connection

**Response:**
```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "name": "Daily Sales Validation",
            "source_connection": {
                "id": 1,
                "name": "Production SQL Server"
            },
            "target_connection": {
                "id": 2,
                "name": "Data Warehouse"
            },
            "status": "completed",
            "last_run": "2024-01-17T02:00:00Z",
            "success_rate": 99.99,
            "created_at": "2024-01-10T10:00:00Z"
        }
    ],
    "meta": {
        "page": 1,
        "per_page": 20,
        "total_count": 15
    }
}
```

---

### Get Validation
Get a specific validation by ID.

**Endpoint:** `GET /validations/{id}`

**Response:**
```json
{
    "success": true,
    "data": {
        "id": 1,
        "name": "Daily Sales Validation",
        "source_connection_id": 1,
        "target_connection_id": 2,
        "config": {
            "comparison_type": "data_match",
            "source_query": "SELECT order_id, customer_id, amount, order_date FROM orders WHERE order_date = CAST(GETDATE() AS DATE)",
            "target_query": "SELECT order_id, customer_id, amount, order_date FROM fact_orders WHERE order_date = CURRENT_DATE",
            "key_columns": ["order_id"],
            "comparison_options": {
                "check_row_count": true,
                "check_column_count": true,
                "check_data_types": false,
                "check_nulls": true,
                "case_sensitive": false
            },
            "error_margin": {
                "type": "percentage",
                "value": 0.01
            }
        },
        "status": "completed",
        "results": {
            "execution_id": "550e8400-e29b-41d4-a716-446655440000",
            "start_time": "2024-01-17T02:00:00Z",
            "end_time": "2024-01-17T02:00:45Z",
            "duration_ms": 45000,
            "status": "failed",
            "summary": {
                "source_row_count": 15000,
                "target_row_count": 14999,
                "matched_rows": 14999,
                "mismatched_rows": 0,
                "missing_in_target": 1,
                "extra_in_target": 0,
                "success_rate": 99.99
            }
        },
        "created_at": "2024-01-10T10:00:00Z",
        "updated_at": "2024-01-17T02:00:45Z"
    }
}
```

---

### Create Validation
Create a new validation configuration.

**Endpoint:** `POST /validations`

**Request Body (Row Count):**
```json
{
    "name": "Order Count Validation",
    "source_connection_id": 1,
    "target_connection_id": 2,
    "config": {
        "comparison_type": "row_count",
        "source_query": "SELECT COUNT(*) FROM orders WHERE order_date >= '2024-01-01'",
        "target_query": "SELECT COUNT(*) FROM staging.orders WHERE order_date >= '2024-01-01'",
        "error_margin": {
            "type": "absolute",
            "value": 0
        }
    }
}
```

**Request Body (Data Match):**
```json
{
    "name": "Customer Data Validation",
    "source_connection_id": 1,
    "target_connection_id": 2,
    "config": {
        "comparison_type": "data_match",
        "source_query": "SELECT customer_id, name, email, phone FROM customers",
        "target_query": "SELECT cust_id as customer_id, full_name as name, email_address as email, phone_number as phone FROM dim_customers",
        "key_columns": ["customer_id"],
        "comparison_options": {
            "check_row_count": true,
            "check_nulls": true,
            "case_sensitive": false,
            "trim_strings": true,
            "decimal_precision": 2
        },
        "performance": {
            "batch_size": 10000,
            "timeout_seconds": 300,
            "max_differences": 1000
        }
    }
}
```

**Request Body (Schema Validation):**
```json
{
    "name": "Table Schema Validation",
    "source_connection_id": 1,
    "target_connection_id": 2,
    "config": {
        "comparison_type": "schema",
        "source_table": "orders",
        "target_table": "fact_orders",
        "check_options": {
            "column_names": true,
            "data_types": true,
            "nullable": true,
            "column_order": false,
            "constraints": false
        }
    }
}
```

---

### Update Validation
Update an existing validation.

**Endpoint:** `PUT /validations/{id}`

**Request Body:**
```json
{
    "name": "Updated Validation Name",
    "config": {
        "source_query": "SELECT * FROM orders WHERE order_date >= DATEADD(day, -7, GETDATE())"
    }
}
```

---

### Delete Validation
Delete a validation.

**Endpoint:** `DELETE /validations/{id}`

**Response:**
```json
{
    "success": true,
    "data": {
        "message": "Validation deleted successfully"
    }
}
```

---

### Run Validation
Execute a validation immediately.

**Endpoint:** `POST /validations/{id}/run`

**Request Body (optional):**
```json
{
    "parameters": {
        "start_date": "2024-01-01",
        "end_date": "2024-01-31"
    }
}
```

**Response:**
```json
{
    "success": true,
    "data": {
        "execution_id": "550e8400-e29b-41d4-a716-446655440001",
        "status": "running",
        "message": "Validation started successfully"
    }
}
```

---

### Get Validation Status
Get the current status of a validation execution.

**Endpoint:** `GET /validations/{id}/status`

**Query Parameters:**
- `execution_id` (string, optional): Specific execution ID

**Response:**
```json
{
    "success": true,
    "data": {
        "execution_id": "550e8400-e29b-41d4-a716-446655440001",
        "status": "running",
        "progress": {
            "stage": "comparing_data",
            "percentage": 45,
            "rows_processed": 450000,
            "total_rows": 1000000
        },
        "start_time": "2024-01-17T10:00:00Z",
        "estimated_completion": "2024-01-17T10:05:00Z"
    }
}
```

---

### Get Validation History
Get execution history for a validation.

**Endpoint:** `GET /validations/{id}/history`

**Query Parameters:**
- `page` (integer, optional): Page number
- `per_page` (integer, optional): Items per page
- `days` (integer, optional): Number of days to look back (default: 30)

**Response:**
```json
{
    "success": true,
    "data": [
        {
            "execution_id": "550e8400-e29b-41d4-a716-446655440001",
            "start_time": "2024-01-17T10:00:00Z",
            "end_time": "2024-01-17T10:02:30Z",
            "duration_seconds": 150,
            "status": "completed",
            "success_rate": 100,
            "source_row_count": 1000000,
            "target_row_count": 1000000
        }
    ],
    "meta": {
        "page": 1,
        "per_page": 20,
        "total_count": 90
    }
}
```

## System Endpoints

### Health Check
Check API health status.

**Endpoint:** `GET /health`

**Response:**
```json
{
    "service": "compareflow",
    "status": "healthy",
    "version": "1.0.0",
    "timestamp": "2024-01-17T10:00:00Z"
}
```

---

### Get System Statistics
Get system-wide statistics (requires admin role).

**Endpoint:** `GET /system/stats`

**Response:**
```json
{
    "success": true,
    "data": {
        "users": {
            "total": 150,
            "active_today": 45,
            "new_this_week": 5
        },
        "connections": {
            "total": 450,
            "by_type": {
                "sqlserver": 300,
                "databricks": 150
            }
        },
        "validations": {
            "total": 1500,
            "running": 5,
            "completed_today": 125,
            "failed_today": 3
        },
        "system": {
            "uptime_hours": 720,
            "database_size_mb": 2048,
            "api_calls_today": 15000
        }
    }
}
```

## Error Responses

All endpoints use standard HTTP status codes and return errors in a consistent format:

```json
{
    "success": false,
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "Validation failed",
        "details": {
            "field": "email",
            "reason": "Invalid email format"
        }
    }
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Missing or invalid authentication token |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `VALIDATION_ERROR` | 400 | Request validation failed |
| `DUPLICATE_RESOURCE` | 409 | Resource already exists |
| `RATE_LIMITED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Internal server error |
| `CONNECTION_FAILED` | 502 | Database connection failed |
| `TIMEOUT` | 504 | Request timeout |

## Pagination

List endpoints support pagination using query parameters:

- `page`: Page number (starting from 1)
- `per_page`: Number of items per page (max: 100)

Paginated responses include metadata:

```json
{
    "meta": {
        "page": 2,
        "per_page": 20,
        "total_count": 150,
        "total_pages": 8,
        "has_next": true,
        "has_prev": true
    }
}
```

## Filtering and Sorting

List endpoints support filtering and sorting:

**Filtering:**
```
GET /validations?status=failed&created_after=2024-01-01
```

**Sorting:**
```
GET /validations?sort=created_at&order=desc
```

Supported sort fields vary by endpoint. Common fields include:
- `created_at`
- `updated_at`
- `name`
- `status`

## WebSocket API (Coming Soon)

Real-time updates will be available via WebSocket:

```javascript
const ws = new WebSocket('wss://api.compareflow.com/ws');

ws.onopen = () => {
    ws.send(JSON.stringify({
        type: 'auth',
        token: 'your_jwt_token'
    }));
    
    ws.send(JSON.stringify({
        type: 'subscribe',
        channel: 'validation:123'
    }));
};

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('Update:', data);
};
```

## SDK Examples

### JavaScript/TypeScript
```typescript
import { CompareFlowClient } from '@compareflow/sdk';

const client = new CompareFlowClient({
    apiKey: 'your_jwt_token',
    baseUrl: 'https://api.compareflow.com'
});

// Create connection
const connection = await client.connections.create({
    name: 'My SQL Server',
    type: 'sqlserver',
    config: {
        server: 'localhost',
        port: 1433,
        database: 'mydb',
        username: 'user',
        password: 'pass'
    }
});

// Run validation
const result = await client.validations.run(validationId);
```

### Python
```python
from compareflow import Client

client = Client(
    api_key='your_jwt_token',
    base_url='https://api.compareflow.com'
)

# Create validation
validation = client.validations.create(
    name='Daily Check',
    source_connection_id=1,
    target_connection_id=2,
    config={
        'comparison_type': 'row_count',
        'source_query': 'SELECT COUNT(*) FROM orders',
        'target_query': 'SELECT COUNT(*) FROM orders'
    }
)

# Run and wait for result
result = client.validations.run_and_wait(validation.id)
print(f"Success rate: {result.summary.success_rate}%")
```

### Go
```go
import "github.com/compareflow/compareflow-go"

client := compareflow.NewClient("your_jwt_token")

// Test connection
result, err := client.Connections.Test(ctx, connectionID)
if err != nil {
    log.Fatal(err)
}

// List validations
validations, err := client.Validations.List(ctx, &compareflow.ListOptions{
    Status: "failed",
    Page: 1,
    PerPage: 20,
})
```

## Best Practices

1. **Authentication**
   - Store tokens securely
   - Implement token refresh logic
   - Never expose tokens in logs or URLs

2. **Error Handling**
   - Always check response status
   - Implement exponential backoff for retries
   - Log errors with context

3. **Performance**
   - Use pagination for large datasets
   - Implement caching where appropriate
   - Batch operations when possible

4. **Security**
   - Always use HTTPS
   - Validate SSL certificates
   - Sanitize user inputs
   - Follow principle of least privilege

## Support

For API support, please contact:
- Email: api-support@compareflow.com
- Documentation: https://docs.compareflow.com
- Status Page: https://status.compareflow.com