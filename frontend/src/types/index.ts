export interface User {
  id: number;
  username: string;
  email: string;
}

export interface Connection {
  id: number;
  name: string;
  type: 'sqlserver' | 'databricks';
  config: {
    server?: string;
    port?: number;
    database?: string;
    username?: string;
    password?: string;
    encrypt?: boolean;
    trust_server_certificate?: boolean;
    workspace?: string;
    http_path?: string;
    access_token?: string;
  };
  created_at?: string;
  updated_at?: string;
}

export interface Validation {
  id: number;
  name: string;
  source_connection_id: number;
  target_connection_id: number;
  source_connection?: Connection;
  target_connection?: Connection;
  config: {
    source_query?: string;
    target_query?: string;
    comparison_type?: 'row_count' | 'data_match' | 'schema';
    key_columns?: string[];
  };
  status: 'pending' | 'running' | 'completed' | 'failed';
  results?: {
    execution_id?: string;
    start_time?: string;
    end_time?: string;
    duration_ms?: number;
    summary?: {
      source_row_count?: number;
      target_row_count?: number;
      matched_rows?: number;
      mismatched_rows?: number;
      missing_in_target?: number;
      extra_in_target?: number;
      success_rate?: number;
    };
    details?: {
      differences?: Array<{
        key: any;
        type: 'missing' | 'extra' | 'mismatch';
        source_data?: any;
        target_data?: any;
        columns?: string[];
      }>;
      column_stats?: Record<string, {
        source_min?: any;
        source_max?: any;
        source_avg?: number;
        target_min?: any;
        target_max?: any;
        target_avg?: number;
      }>;
    };
    errors?: Array<string | {
      timestamp?: string;
      message: string;
      details?: string;
    }>;
  };
  created_at?: string;
  updated_at?: string;
}

export interface TableInfo {
  name: string;
  columns?: ColumnInfo[];
}

export interface ColumnInfo {
  name: string;
  data_type: string;
  nullable: boolean;
}