# CompareFlow Development Guide

## Table of Contents
1. [Development Environment Setup](#development-environment-setup)
2. [Project Structure](#project-structure)
3. [Coding Standards](#coding-standards)
4. [Development Workflow](#development-workflow)
5. [Adding Features](#adding-features)
6. [Testing](#testing)
7. [Debugging](#debugging)
8. [Performance Optimization](#performance-optimization)
9. [Contributing](#contributing)

## 1. Development Environment Setup

### 1.1 Required Tools

```bash
# Install Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install Node.js (via nvm)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 18
nvm use 18

# Install development tools
go install github.com/cosmtrek/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
npm install -g prettier eslint
```

### 1.2 IDE Setup

#### VS Code Extensions
```json
{
  "recommendations": [
    "golang.go",
    "dbaeumer.vscode-eslint",
    "esbenp.prettier-vscode",
    "bradlc.vscode-tailwindcss",
    "formulahendry.auto-rename-tag",
    "yzhang.markdown-all-in-one"
  ]
}
```

#### VS Code Settings
```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package",
  "editor.formatOnSave": true,
  "[go]": {
    "editor.defaultFormatter": "golang.go"
  },
  "[typescript]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[typescriptreact]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  }
}
```

### 1.3 Local Development Setup

```bash
# Clone repository
git clone https://github.com/compareflow/compareflow.git
cd compareflow

# Install dependencies
go mod download
cd frontend && npm install && cd ..

# Copy environment template
cp .env.example .env

# Start local PostgreSQL
./start-local.sh

# Run with hot reload
air
```

## 2. Project Structure

### 2.1 Backend Structure

```
internal/
├── api/                    # HTTP layer
│   ├── handlers/          # Request handlers
│   │   ├── auth.go       # Authentication endpoints
│   │   ├── connection.go # Connection CRUD
│   │   └── validation.go # Validation operations
│   ├── middleware/        # HTTP middleware
│   │   ├── auth.go       # JWT validation
│   │   ├── cors.go       # CORS handling
│   │   └── logger.go     # Request logging
│   └── routes.go         # Route definitions
├── config/               # Configuration
│   └── config.go        # Config struct and loader
├── database/            # Database setup
│   └── database.go      # GORM initialization
├── models/              # Domain models
│   ├── user.go         # User model
│   ├── connection.go   # Connection model
│   └── validation.go   # Validation model
├── services/           # Business logic
│   ├── auth_service.go
│   ├── connection_service.go
│   └── validation_service.go
└── validation/         # Validation engine
    ├── engine.go      # Core engine
    ├── validators/    # Validation types
    └── processors/    # Result processors
```

### 2.2 Frontend Structure

```
frontend/src/
├── components/          # Reusable components
│   ├── common/         # Generic components
│   ├── forms/          # Form components
│   └── layouts/        # Layout components
├── pages/              # Page components
│   ├── Dashboard.tsx
│   ├── Connections.tsx
│   └── Validations.tsx
├── services/           # API services
│   ├── api.ts         # Axios instance
│   └── authService.ts # Auth API calls
├── store/             # Redux store
│   ├── index.ts
│   └── slices/       # Redux slices
├── types/            # TypeScript types
├── utils/            # Helper functions
└── hooks/            # Custom React hooks
```

## 3. Coding Standards

### 3.1 Go Coding Standards

#### Naming Conventions
```go
// Package names: lowercase, single word
package auth

// Interface names: end with 'er'
type Validator interface {
    Validate() error
}

// Struct names: PascalCase
type UserService struct {
    db *gorm.DB
}

// Function names: PascalCase for exported, camelCase for private
func GetUserByID(id uint) (*User, error) {}
func validateEmail(email string) error {}

// Constants: PascalCase
const MaxRetries = 3
const defaultTimeout = 30 * time.Second
```

#### Error Handling
```go
// Always handle errors explicitly
user, err := GetUserByID(id)
if err != nil {
    return fmt.Errorf("failed to get user %d: %w", id, err)
}

// Custom error types
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

// Error checking
if errors.Is(err, ErrNotFound) {
    return c.JSON(404, gin.H{"error": "User not found"})
}
```

#### Comments and Documentation
```go
// Package auth provides authentication and authorization functionality.
package auth

// UserService handles user-related business logic.
// It provides methods for user creation, authentication, and profile management.
type UserService struct {
    db    *gorm.DB
    cache *Cache
}

// CreateUser creates a new user with the provided details.
// It validates the input, hashes the password, and stores the user in the database.
// Returns the created user or an error if validation fails or the user already exists.
func (s *UserService) CreateUser(username, email, password string) (*User, error) {
    // Implementation
}
```

### 3.2 TypeScript/React Coding Standards

#### Component Structure
```tsx
// Use functional components with TypeScript
interface DashboardProps {
  user: User;
  onRefresh: () => void;
}

export const Dashboard: React.FC<DashboardProps> = ({ user, onRefresh }) => {
  // Hooks at the top
  const [loading, setLoading] = useState(false);
  const dispatch = useAppDispatch();
  
  // Event handlers
  const handleRefresh = useCallback(() => {
    setLoading(true);
    onRefresh();
  }, [onRefresh]);
  
  // Effects
  useEffect(() => {
    // Effect logic
  }, []);
  
  // Render
  return (
    <div className="dashboard">
      {/* Component JSX */}
    </div>
  );
};
```

#### State Management
```typescript
// Redux Toolkit slice
import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface AuthState {
  user: User | null;
  token: string | null;
  loading: boolean;
  error: string | null;
}

const authSlice = createSlice({
  name: 'auth',
  initialState: {
    user: null,
    token: null,
    loading: false,
    error: null,
  } as AuthState,
  reducers: {
    loginStart: (state) => {
      state.loading = true;
      state.error = null;
    },
    loginSuccess: (state, action: PayloadAction<{ user: User; token: string }>) => {
      state.user = action.payload.user;
      state.token = action.payload.token;
      state.loading = false;
    },
    loginFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
  },
});
```

### 3.3 SQL and Database Standards

```sql
-- Table naming: plural, snake_case
CREATE TABLE users (
    -- Column naming: snake_case
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index naming: idx_table_columns
CREATE INDEX idx_users_username ON users(username);

-- Foreign key naming: fk_table_referenced
ALTER TABLE validations 
ADD CONSTRAINT fk_validations_user 
FOREIGN KEY (user_id) REFERENCES users(id);
```

## 4. Development Workflow

### 4.1 Git Workflow

```bash
# Feature branch workflow
git checkout -b feature/add-validation-export

# Make changes
git add .
git commit -m "feat: add CSV export for validation results"

# Push and create PR
git push origin feature/add-validation-export
```

#### Commit Message Convention
```
type(scope): subject

body

footer
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test additions/changes
- `chore`: Build process or auxiliary tool changes

Example:
```
feat(validation): add support for Snowflake connections

- Implement Snowflake driver
- Add connection configuration UI
- Update documentation

Closes #123
```

### 4.2 Development Commands

```bash
# Backend development
air                          # Run with hot reload
go test ./...               # Run all tests
go test -race ./...         # Run tests with race detection
golangci-lint run           # Run linters
go mod tidy                 # Clean up dependencies

# Frontend development
npm run dev                 # Start dev server
npm run build              # Build for production
npm run test               # Run tests
npm run lint               # Run ESLint
npm run format             # Format with Prettier

# Database
make migrate-up            # Run migrations
make migrate-down          # Rollback migration
make migrate-create name=x # Create new migration

# Full stack
make dev                   # Run backend and frontend
make test                  # Run all tests
make build                 # Build production binary
```

## 5. Adding Features

### 5.1 Adding a New Database Connector

Example: Adding Databricks Support

1. **Add the driver dependency** (`go.mod`):
```bash
go get github.com/databricks/databricks-sql-go
```

2. **Import the driver** (`internal/services/connection_service.go`):
```go
import (
    _ "github.com/databricks/databricks-sql-go"
)
```

3. **Implement connection testing**:
```go
func (s *ConnectionService) testDatabricksConnection(conn *models.Connection) error {
    // Extract configuration
    workspace, _ := conn.Config["workspace"].(string)
    httpPath, _ := conn.Config["http_path"].(string)
    accessToken, _ := conn.Config["access_token"].(string)
    
    // Build connection string
    host := strings.TrimPrefix(workspace, "https://")
    connString := fmt.Sprintf("databricks://token:%s@%s:443%s", accessToken, host, httpPath)
    
    // Test connection
    db, err := sql.Open("databricks", connString)
    if err != nil {
        return fmt.Errorf("failed to open connection: %w", err)
    }
    defer db.Close()
    
    // Verify with a simple query
    var result int
    err = db.QueryRow("SELECT 1").Scan(&result)
    if err != nil {
        return fmt.Errorf("failed to execute test query: %w", err)
    }
    
    return nil
}
```

4. **Implement table listing**:
```go
func (s *ConnectionService) getDatabricksTables(conn *models.Connection) ([]string, error) {
    // Connect to Databricks
    db, err := s.connectToDatabricks(conn)
    if err != nil {
        return nil, err
    }
    defer db.Close()
    
    // Query information schema
    query := `
        SELECT table_schema || '.' || table_name AS full_table_name
        FROM information_schema.tables
        WHERE table_type = 'TABLE'
        ORDER BY table_schema, table_name
    `
    
    rows, err := db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to query tables: %w", err)
    }
    defer rows.Close()
    
    var tables []string
    for rows.Next() {
        var table string
        if err := rows.Scan(&table); err != nil {
            return nil, err
        }
        tables = append(tables, table)
    }
    
    return tables, nil
}
```

### 5.2 Adding a New API Endpoint

1. **Define the handler** (`internal/api/handlers/feature.go`):
```go
type FeatureHandler struct {
    service *services.FeatureService
}

func NewFeatureHandler(service *services.FeatureService) *FeatureHandler {
    return &FeatureHandler{service: service}
}

func (h *FeatureHandler) Create(c *gin.Context) {
    var req CreateFeatureRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    userID := c.GetUint("user_id")
    feature, err := h.service.Create(userID, req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, feature)
}
```

2. **Add the route** (`internal/api/routes.go`):
```go
featureHandler := handlers.NewFeatureHandler(featureService)
protected.POST("/features", featureHandler.Create)
protected.GET("/features", featureHandler.List)
protected.GET("/features/:id", featureHandler.Get)
```

3. **Implement the service** (`internal/services/feature_service.go`):
```go
type FeatureService struct {
    db *gorm.DB
}

func (s *FeatureService) Create(userID uint, req CreateFeatureRequest) (*models.Feature, error) {
    feature := &models.Feature{
        UserID: userID,
        Name:   req.Name,
        Config: req.Config,
    }
    
    if err := s.db.Create(feature).Error; err != nil {
        return nil, fmt.Errorf("failed to create feature: %w", err)
    }
    
    return feature, nil
}
```

### 5.2 Adding a New Frontend Page

1. **Create the page component** (`frontend/src/pages/Feature.tsx`):
```tsx
import React, { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { fetchFeatures } from '../store/slices/featureSlice';

export const Feature: React.FC = () => {
    const dispatch = useAppDispatch();
    const { features, loading } = useAppSelector(state => state.features);
    
    useEffect(() => {
        dispatch(fetchFeatures());
    }, [dispatch]);
    
    if (loading) return <LoadingSpinner />;
    
    return (
        <div>
            <h1>Features</h1>
            {features.map(feature => (
                <FeatureCard key={feature.id} feature={feature} />
            ))}
        </div>
    );
};
```

2. **Add Redux slice** (`frontend/src/store/slices/featureSlice.ts`):
```typescript
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { featureService } from '../../services/featureService';

export const fetchFeatures = createAsyncThunk(
    'features/fetchAll',
    async () => {
        return await featureService.getAll();
    }
);

const featureSlice = createSlice({
    name: 'features',
    initialState: {
        features: [],
        loading: false,
        error: null,
    },
    reducers: {},
    extraReducers: (builder) => {
        builder
            .addCase(fetchFeatures.pending, (state) => {
                state.loading = true;
            })
            .addCase(fetchFeatures.fulfilled, (state, action) => {
                state.features = action.payload;
                state.loading = false;
            })
            .addCase(fetchFeatures.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message;
            });
    },
});
```

3. **Add route** (`frontend/src/App.tsx`):
```tsx
<Route path="/features" element={<Feature />} />
```

### 5.3 Adding a New Database Table

1. **Create migration**:
```bash
make migrate-create name=add_features_table
```

2. **Write migration** (`migrations/xxx_add_features_table.sql`):
```sql
-- +migrate Up
CREATE TABLE features (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    config JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_features_user_id ON features(user_id);

-- +migrate Down
DROP TABLE IF EXISTS features;
```

3. **Create model** (`internal/models/feature.go`):
```go
type Feature struct {
    ID        uint           `gorm:"primarykey"`
    UserID    uint           `gorm:"not null"`
    User      User           `gorm:"foreignKey:UserID"`
    Name      string         `gorm:"not null"`
    Config    datatypes.JSON `gorm:"type:jsonb;default:'{}'"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

## 6. Testing

### 6.1 Backend Testing

#### Unit Tests
```go
// internal/services/auth_service_test.go
func TestCreateUser(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(db)
    
    service := NewAuthService(db)
    
    tests := []struct {
        name     string
        username string
        email    string
        password string
        wantErr  bool
    }{
        {
            name:     "valid user",
            username: "testuser",
            email:    "test@example.com",
            password: "password123",
            wantErr:  false,
        },
        {
            name:     "duplicate username",
            username: "testuser",
            email:    "test2@example.com",
            password: "password123",
            wantErr:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            user, err := service.CreateUser(tt.username, tt.email, tt.password)
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
            }
            if !tt.wantErr && user == nil {
                t.Error("CreateUser() returned nil user")
            }
        })
    }
}
```

#### Integration Tests
```go
// tests/integration/api_test.go
func TestAPIFlow(t *testing.T) {
    app := setupTestApp(t)
    
    // Register user
    resp := httptest.NewRecorder()
    body := `{"username":"test","email":"test@example.com","password":"password123"}`
    req, _ := http.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    app.ServeHTTP(resp, req)
    
    assert.Equal(t, 201, resp.Code)
    
    var result map[string]interface{}
    json.Unmarshal(resp.Body.Bytes(), &result)
    token := result["token"].(string)
    
    // Use token to access protected endpoint
    resp = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/api/v1/connections", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    app.ServeHTTP(resp, req)
    
    assert.Equal(t, 200, resp.Code)
}
```

### 6.2 Frontend Testing

#### Component Tests
```tsx
// frontend/src/pages/__tests__/Dashboard.test.tsx
import { render, screen, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { Dashboard } from '../Dashboard';
import { store } from '../../store';

describe('Dashboard', () => {
    it('renders loading state initially', () => {
        render(
            <Provider store={store}>
                <Dashboard />
            </Provider>
        );
        
        expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
    });
    
    it('renders data after loading', async () => {
        render(
            <Provider store={store}>
                <Dashboard />
            </Provider>
        );
        
        await waitFor(() => {
            expect(screen.getByText('Total Connections')).toBeInTheDocument();
        });
    });
});
```

#### Service Tests
```typescript
// frontend/src/services/__tests__/authService.test.ts
import { authService } from '../authService';
import api from '../api';

jest.mock('../api');

describe('AuthService', () => {
    it('logs in successfully', async () => {
        const mockResponse = {
            data: {
                user: { id: 1, username: 'test' },
                token: 'fake-token'
            }
        };
        
        (api.post as jest.Mock).mockResolvedValueOnce(mockResponse);
        
        const result = await authService.login('test', 'password');
        
        expect(api.post).toHaveBeenCalledWith('/auth/login', {
            username: 'test',
            password: 'password'
        });
        expect(result).toEqual(mockResponse.data);
    });
});
```

### 6.3 End-to-End Testing

```typescript
// e2e/login.spec.ts
import { test, expect } from '@playwright/test';

test('user can log in', async ({ page }) => {
    await page.goto('http://localhost:8080');
    
    // Navigate to login
    await page.click('text=Login');
    
    // Fill form
    await page.fill('input[name="username"]', 'testuser');
    await page.fill('input[name="password"]', 'password123');
    
    // Submit
    await page.click('button[type="submit"]');
    
    // Verify redirect to dashboard
    await expect(page).toHaveURL('http://localhost:8080/dashboard');
    await expect(page.locator('h1')).toContainText('Dashboard');
});
```

## 7. Debugging

### 7.1 Backend Debugging

#### Using Delve
```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug the application
dlv debug cmd/compareflow/main.go

# Set breakpoint
(dlv) break main.go:20

# Continue execution
(dlv) continue

# Print variable
(dlv) print variableName

# Step through code
(dlv) next
(dlv) step
```

#### VS Code Debugging
```json
// .vscode/launch.json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug Backend",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/compareflow",
            "env": {
                "DATABASE_URL": "postgresql://compareflow:compareflow123@localhost:5432/compareflow?sslmode=disable",
                "JWT_SECRET": "debug-secret",
                "GIN_MODE": "debug"
            }
        }
    ]
}
```

#### Logging
```go
// Structured logging with context
log := logger.WithFields(logrus.Fields{
    "user_id": userID,
    "action":  "create_validation",
})

log.Info("Starting validation creation")

// Log errors with stack trace
if err != nil {
    log.WithError(err).Error("Failed to create validation")
}

// Debug logging
if config.Debug {
    log.WithField("query", query).Debug("Executing SQL query")
}
```

### 7.2 Frontend Debugging

#### React Developer Tools
```typescript
// Add debug info to components
export const Dashboard: React.FC = () => {
    // Debug render count
    const renderCount = useRef(0);
    renderCount.current++;
    console.log(`Dashboard rendered ${renderCount.current} times`);
    
    // Debug props
    console.log('Dashboard props:', props);
    
    // Debug state
    const state = useAppSelector(state => state);
    console.log('Redux state:', state);
};
```

#### Network Debugging
```typescript
// Add request/response interceptors
api.interceptors.request.use(request => {
    console.log('Starting Request:', request);
    return request;
});

api.interceptors.response.use(
    response => {
        console.log('Response:', response);
        return response;
    },
    error => {
        console.error('Response Error:', error.response);
        return Promise.reject(error);
    }
);
```

### 7.3 Database Debugging

```sql
-- Enable query logging
ALTER SYSTEM SET log_statement = 'all';
ALTER SYSTEM SET log_duration = on;
SELECT pg_reload_conf();

-- Check slow queries
SELECT 
    query,
    calls,
    total_time,
    mean_time,
    max_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;

-- Check locks
SELECT 
    pid,
    usename,
    query,
    state,
    wait_event_type,
    wait_event
FROM pg_stat_activity
WHERE state != 'idle';
```

## 8. Performance Optimization

### 8.1 Backend Optimization

#### Query Optimization
```go
// Use preloading to avoid N+1 queries
var validations []models.Validation
db.Preload("SourceConnection").
   Preload("TargetConnection").
   Where("user_id = ?", userID).
   Find(&validations)

// Use select to load only needed fields
db.Select("id", "name", "status").Find(&validations)

// Use pagination
db.Offset((page - 1) * perPage).Limit(perPage).Find(&validations)
```

#### Caching
```go
// In-memory cache with TTL
type Cache struct {
    data sync.Map
    ttl  time.Duration
}

func (c *Cache) Get(key string) (interface{}, bool) {
    if val, ok := c.data.Load(key); ok {
        item := val.(*cacheItem)
        if time.Now().Before(item.expiry) {
            return item.value, true
        }
        c.data.Delete(key)
    }
    return nil, false
}

func (c *Cache) Set(key string, value interface{}) {
    c.data.Store(key, &cacheItem{
        value:  value,
        expiry: time.Now().Add(c.ttl),
    })
}
```

#### Concurrent Processing
```go
// Process validations concurrently
func (e *Engine) RunValidations(validations []Validation) []Result {
    results := make([]Result, len(validations))
    var wg sync.WaitGroup
    semaphore := make(chan struct{}, 5) // Limit to 5 concurrent
    
    for i, validation := range validations {
        wg.Add(1)
        go func(idx int, val Validation) {
            defer wg.Done()
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            results[idx] = e.runSingle(val)
        }(i, validation)
    }
    
    wg.Wait()
    return results
}
```

### 8.2 Frontend Optimization

#### React Optimization
```tsx
// Memoize expensive computations
const expensiveValue = useMemo(() => {
    return computeExpensiveValue(data);
}, [data]);

// Memoize callbacks
const handleClick = useCallback(() => {
    doSomething(id);
}, [id]);

// Use React.memo for pure components
export const ExpensiveComponent = React.memo(({ data }) => {
    return <div>{/* Render logic */}</div>;
}, (prevProps, nextProps) => {
    // Custom comparison
    return prevProps.data.id === nextProps.data.id;
});

// Lazy load components
const HeavyComponent = lazy(() => import('./HeavyComponent'));

// Use virtualization for long lists
import { FixedSizeList } from 'react-window';

const VirtualizedList = ({ items }) => (
    <FixedSizeList
        height={600}
        itemCount={items.length}
        itemSize={50}
        width="100%"
    >
        {({ index, style }) => (
            <div style={style}>
                {items[index].name}
            </div>
        )}
    </FixedSizeList>
);
```

#### Bundle Optimization
```javascript
// vite.config.ts
export default defineConfig({
    build: {
        rollupOptions: {
            output: {
                manualChunks: {
                    'react-vendor': ['react', 'react-dom', 'react-router-dom'],
                    'ui-vendor': ['@mui/material', '@emotion/react'],
                    'redux-vendor': ['@reduxjs/toolkit', 'react-redux'],
                },
            },
        },
        // Enable compression
        minify: 'terser',
        terserOptions: {
            compress: {
                drop_console: true,
                drop_debugger: true,
            },
        },
    },
});
```

### 8.3 Database Optimization

```sql
-- Add indexes for common queries
CREATE INDEX idx_validations_user_status_created 
ON validations(user_id, status, created_at DESC);

-- Partial indexes for specific conditions
CREATE INDEX idx_validations_pending 
ON validations(status) 
WHERE status = 'pending';

-- Use materialized views for complex reports
CREATE MATERIALIZED VIEW validation_summary AS
SELECT 
    user_id,
    DATE(created_at) as date,
    COUNT(*) as total,
    COUNT(*) FILTER (WHERE status = 'success') as successful,
    AVG(EXTRACT(EPOCH FROM (updated_at - created_at))) as avg_duration
FROM validations
GROUP BY user_id, DATE(created_at);

-- Refresh periodically
REFRESH MATERIALIZED VIEW CONCURRENTLY validation_summary;
```

## 9. Contributing

### 9.1 Pull Request Process

1. **Fork and clone**
```bash
git clone https://github.com/yourusername/compareflow.git
cd compareflow
git remote add upstream https://github.com/compareflow/compareflow.git
```

2. **Create feature branch**
```bash
git checkout -b feature/your-feature-name
```

3. **Make changes and test**
```bash
# Make your changes
# Run tests
make test

# Run linters
make lint

# Build to ensure it compiles
make build
```

4. **Commit with conventional commits**
```bash
git add .
git commit -m "feat: add amazing feature"
```

5. **Push and create PR**
```bash
git push origin feature/your-feature-name
```

### 9.2 Code Review Checklist

- [ ] Tests pass
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] No sensitive data exposed
- [ ] Performance impact considered
- [ ] Database migrations included (if needed)
- [ ] UI responsive on mobile (if frontend)
- [ ] Accessibility considered (if frontend)

### 9.3 Release Process

```bash
# 1. Update version
VERSION=v1.2.0

# 2. Update CHANGELOG.md
# Add release notes

# 3. Commit version bump
git add .
git commit -m "chore: bump version to $VERSION"

# 4. Tag release
git tag -a $VERSION -m "Release $VERSION"

# 5. Push
git push origin main
git push origin $VERSION

# 6. Build release artifacts
make release

# 7. Create GitHub release
# Upload artifacts
```

## Conclusion

This development guide provides the foundation for contributing to CompareFlow. Always prioritize code quality, test coverage, and user experience. Happy coding!