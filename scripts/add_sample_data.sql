-- Sample data for CompareFlow
-- This creates two PostgreSQL connections and a validation between them

-- First, let's get the admin user ID (assuming username is 'admin')
DO $$
DECLARE
    admin_id INTEGER;
BEGIN
    SELECT id INTO admin_id FROM users WHERE username = 'admin';
    
    -- Create two PostgreSQL connections
    -- Source connection (simulating a production database)
    INSERT INTO connections (name, type, config, user_id, created_at, updated_at)
    VALUES (
        'Production PostgreSQL',
        'postgresql',
        '{
            "host": "localhost",
            "port": 5432,
            "database": "compareflow",
            "username": "compareflow",
            "password": "compareflow123",
            "sslmode": "disable"
        }'::json,
        admin_id,
        NOW(),
        NOW()
    );

    -- Target connection (simulating a data warehouse)
    INSERT INTO connections (name, type, config, user_id, created_at, updated_at)
    VALUES (
        'Data Warehouse PostgreSQL',
        'postgresql',
        '{
            "host": "localhost",
            "port": 5432,
            "database": "compareflow",
            "username": "compareflow",
            "password": "compareflow123",
            "sslmode": "disable"
        }'::json,
        admin_id,
        NOW(),
        NOW()
    );

    -- Create a sample validation
    INSERT INTO validations (name, source_connection_id, target_connection_id, config, status, user_id, created_at, updated_at)
    VALUES (
        'Daily User Count Validation',
        (SELECT id FROM connections WHERE name = 'Production PostgreSQL' AND user_id = admin_id),
        (SELECT id FROM connections WHERE name = 'Data Warehouse PostgreSQL' AND user_id = admin_id),
        '{
            "comparison_type": "row_count",
            "source_query": "SELECT COUNT(*) FROM users",
            "target_query": "SELECT COUNT(*) FROM users",
            "error_margin": {
                "type": "absolute",
                "value": 0
            }
        }'::json,
        'pending',
        admin_id,
        NOW(),
        NOW()
    );

    -- Create another validation for data matching
    INSERT INTO validations (name, source_connection_id, target_connection_id, config, status, user_id, created_at, updated_at)
    VALUES (
        'User Data Integrity Check',
        (SELECT id FROM connections WHERE name = 'Production PostgreSQL' AND user_id = admin_id),
        (SELECT id FROM connections WHERE name = 'Data Warehouse PostgreSQL' AND user_id = admin_id),
        '{
            "comparison_type": "data_match",
            "source_query": "SELECT id, username, email FROM users ORDER BY id",
            "target_query": "SELECT id, username, email FROM users ORDER BY id",
            "key_columns": ["id"],
            "comparison_options": {
                "check_row_count": true,
                "check_nulls": true,
                "case_sensitive": false
            }
        }'::json,
        'pending',
        admin_id,
        NOW(),
        NOW()
    );

    RAISE NOTICE 'Sample data created successfully!';
END $$;