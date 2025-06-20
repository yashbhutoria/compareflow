package sqlserver

import (
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: Config{
				Server:   "localhost",
				Port:     1433,
				Database: "testdb",
				Username: "sa",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name: "missing server",
			config: Config{
				Database: "testdb",
				Username: "sa",
				Password: "password",
			},
			wantErr: true,
			errMsg:  "server is required",
		},
		{
			name: "missing database",
			config: Config{
				Server:   "localhost",
				Username: "sa",
				Password: "password",
			},
			wantErr: true,
			errMsg:  "database is required",
		},
		{
			name: "default port",
			config: Config{
				Server:   "localhost",
				Database: "testdb",
				Username: "sa",
				Password: "password",
			},
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("Config.Validate() error = %v, want %v", err.Error(), tt.errMsg)
			}
			// Check default port
			if !tt.wantErr && tt.config.Port == 0 {
				t.Error("Config.Validate() should set default port")
			}
		})
	}
}

func TestConnector_ParseConfig(t *testing.T) {
	connector := New()
	
	tests := []struct {
		name      string
		configMap map[string]interface{}
		wantErr   bool
	}{
		{
			name: "valid config map",
			configMap: map[string]interface{}{
				"server":   "localhost",
				"port":     1433,
				"database": "testdb",
				"username": "sa",
				"password": "password",
			},
			wantErr: false,
		},
		{
			name: "invalid port type",
			configMap: map[string]interface{}{
				"server":   "localhost",
				"port":     "not-a-number",
				"database": "testdb",
				"username": "sa",
				"password": "password",
			},
			wantErr: true,
		},
		{
			name: "missing required field",
			configMap: map[string]interface{}{
				"server": "localhost",
				"port":   1433,
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := connector.ParseConfig(tt.configMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connector.ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && config == nil {
				t.Error("Connector.ParseConfig() returned nil config")
			}
		})
	}
}

func TestConnector_Type(t *testing.T) {
	connector := New()
	if got := connector.Type(); got != "sqlserver" {
		t.Errorf("Connector.Type() = %v, want %v", got, "sqlserver")
	}
}

func TestBuildConnectionString(t *testing.T) {
	connector := New()
	
	tests := []struct {
		name   string
		config Config
		want   string
	}{
		{
			name: "basic config",
			config: Config{
				Server:   "localhost",
				Port:     1433,
				Database: "testdb",
				Username: "sa",
				Password: "password",
			},
			want: "server=localhost;port=1433;database=testdb;user id=sa;password=password;encrypt=false;TrustServerCertificate=true",
		},
		{
			name: "with encryption",
			config: Config{
				Server:   "localhost",
				Port:     1433,
				Database: "testdb",
				Username: "sa",
				Password: "password",
				Encrypt:  true,
			},
			want: "server=localhost;port=1433;database=testdb;user id=sa;password=password;encrypt=true;TrustServerCertificate=true",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := connector.buildConnectionString(&tt.config); got != tt.want {
				t.Errorf("buildConnectionString() = %v, want %v", got, tt.want)
			}
		})
	}
}