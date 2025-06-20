-- Add more sample validations using the sample tables

DO $$
DECLARE
    admin_id INTEGER;
    source_conn_id INTEGER;
    target_conn_id INTEGER;
BEGIN
    -- Get admin user ID
    SELECT id INTO admin_id FROM users WHERE username = 'admin';
    
    -- Get connection IDs
    SELECT id INTO source_conn_id FROM connections WHERE name = 'Production PostgreSQL' AND user_id = admin_id;
    SELECT id INTO target_conn_id FROM connections WHERE name = 'Data Warehouse PostgreSQL' AND user_id = admin_id;

    -- Validation 1: Products row count
    INSERT INTO validations (name, source_connection_id, target_connection_id, config, status, user_id, created_at, updated_at)
    VALUES (
        'Products Table Row Count',
        source_conn_id,
        target_conn_id,
        '{
            "comparison_type": "row_count",
            "source_query": "SELECT COUNT(*) FROM products",
            "target_query": "SELECT COUNT(*) FROM products_dw",
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

    -- Validation 2: Orders row count (this will fail - missing 1 order)
    INSERT INTO validations (name, source_connection_id, target_connection_id, config, status, user_id, created_at, updated_at)
    VALUES (
        'Orders Table Row Count',
        source_conn_id,
        target_conn_id,
        '{
            "comparison_type": "row_count",
            "source_query": "SELECT COUNT(*) FROM orders",
            "target_query": "SELECT COUNT(*) FROM orders_dw",
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

    -- Validation 3: Product data match
    INSERT INTO validations (name, source_connection_id, target_connection_id, config, status, user_id, created_at, updated_at)
    VALUES (
        'Product Data Integrity',
        source_conn_id,
        target_conn_id,
        '{
            "comparison_type": "data_match",
            "source_query": "SELECT id, name, price, category FROM products ORDER BY id",
            "target_query": "SELECT id, name, price, category FROM products_dw ORDER BY id",
            "key_columns": ["id"],
            "comparison_options": {
                "check_row_count": true,
                "check_nulls": true,
                "case_sensitive": false,
                "decimal_precision": 2
            }
        }'::json,
        'pending',
        admin_id,
        NOW(),
        NOW()
    );

    -- Validation 4: Daily orders summary
    INSERT INTO validations (name, source_connection_id, target_connection_id, config, status, user_id, created_at, updated_at)
    VALUES (
        'Daily Orders Summary Check',
        source_conn_id,
        target_conn_id,
        '{
            "comparison_type": "data_match",
            "source_query": "SELECT order_date, COUNT(*) as order_count, SUM(total_amount) as total_revenue FROM orders GROUP BY order_date ORDER BY order_date",
            "target_query": "SELECT order_date, COUNT(*) as order_count, SUM(total_amount) as total_revenue FROM orders_dw GROUP BY order_date ORDER BY order_date",
            "key_columns": ["order_date"],
            "comparison_options": {
                "check_row_count": true,
                "decimal_precision": 2
            }
        }'::json,
        'pending',
        admin_id,
        NOW(),
        NOW()
    );

    -- Validation 5: Customer email uniqueness check
    INSERT INTO validations (name, source_connection_id, target_connection_id, config, status, user_id, created_at, updated_at)
    VALUES (
        'Customer Email Consistency',
        source_conn_id,
        target_conn_id,
        '{
            "comparison_type": "data_match",
            "source_query": "SELECT email, name, city FROM customers ORDER BY email",
            "target_query": "SELECT email, name, city FROM customers_dw ORDER BY email",
            "key_columns": ["email"],
            "comparison_options": {
                "check_row_count": true,
                "case_sensitive": false,
                "trim_strings": true
            }
        }'::json,
        'pending',
        admin_id,
        NOW(),
        NOW()
    );

    -- Validation 6: High value orders check
    INSERT INTO validations (name, source_connection_id, target_connection_id, config, status, user_id, created_at, updated_at)
    VALUES (
        'High Value Orders (>$500)',
        source_conn_id,
        target_conn_id,
        '{
            "comparison_type": "row_count",
            "source_query": "SELECT COUNT(*) FROM orders WHERE total_amount > 500",
            "target_query": "SELECT COUNT(*) FROM orders_dw WHERE total_amount > 500",
            "error_margin": {
                "type": "percentage",
                "value": 0
            }
        }'::json,
        'pending',
        admin_id,
        NOW(),
        NOW()
    );

    RAISE NOTICE 'Additional validations created successfully!';
END $$;