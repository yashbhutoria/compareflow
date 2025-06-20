package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type ConnectionType string

const (
	ConnectionTypeSQLServer  ConnectionType = "sqlserver"
	ConnectionTypeDatabricks ConnectionType = "databricks"
	ConnectionTypePostgreSQL ConnectionType = "postgresql"
)

type ConnectionConfig map[string]interface{}

func (c ConnectionConfig) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *ConnectionConfig) Scan(value interface{}) error {
	if value == nil {
		*c = make(ConnectionConfig)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("cannot scan non-byte value into ConnectionConfig")
	}

	return json.Unmarshal(bytes, c)
}

type Connection struct {
	ID        uint             `json:"id"`
	Name      string           `gorm:"not null" json:"name"`
	Type      ConnectionType   `gorm:"not null" json:"type"`
	Config    ConnectionConfig `gorm:"type:json" json:"config"`
	UserID    uint             `json:"user_id"`
	User      *User            `gorm:"foreignKey:UserID" json:"-"`
	CreatedAt time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
}

