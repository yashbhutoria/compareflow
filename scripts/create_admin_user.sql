-- Create admin user for CompareFlow
-- Password: admin123 (bcrypt hash)

INSERT INTO users (username, email, password_hash, created_at, updated_at)
VALUES (
    'admin',
    'admin@compareflow.com',
    '$2a$10$K.0HwpsoPDGaB/atFBmmXOGTw4ceeg33.WrxJx/FeC9.gCyYvIbs6',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT (username) DO NOTHING;

-- You can generate a new bcrypt hash for a different password using:
-- go run -e 'import "golang.org/x/crypto/bcrypt"; import "fmt"; hash, _ := bcrypt.GenerateFromPassword([]byte("your-password"), bcrypt.DefaultCost); fmt.Println(string(hash))'