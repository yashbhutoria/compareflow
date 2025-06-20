package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ValidationStatus string

const (
	ValidationStatusPending   ValidationStatus = "pending"
	ValidationStatusRunning   ValidationStatus = "running"
	ValidationStatusCompleted ValidationStatus = "completed"
	ValidationStatusFailed    ValidationStatus = "failed"
)

type ValidationConfig map[string]interface{}

func (v ValidationConfig) Value() (driver.Value, error) {
	return json.Marshal(v)
}

func (v *ValidationConfig) Scan(value interface{}) error {
	if value == nil {
		*v = make(ValidationConfig)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("cannot scan non-byte value into ValidationConfig")
	}

	return json.Unmarshal(bytes, v)
}

type ValidationResults map[string]interface{}

func (v ValidationResults) Value() (driver.Value, error) {
	return json.Marshal(v)
}

func (v *ValidationResults) Scan(value interface{}) error {
	if value == nil {
		*v = make(ValidationResults)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("cannot scan non-byte value into ValidationResults")
	}

	return json.Unmarshal(bytes, v)
}

type Validation struct {
	ID                 uint              `json:"id"`
	Name               string            `gorm:"not null" json:"name"`
	SourceConnectionID uint              `json:"source_connection_id"`
	TargetConnectionID uint              `json:"target_connection_id"`
	SourceConnection   *Connection       `gorm:"foreignKey:SourceConnectionID" json:"source_connection,omitempty"`
	TargetConnection   *Connection       `gorm:"foreignKey:TargetConnectionID" json:"target_connection,omitempty"`
	Config             ValidationConfig  `gorm:"type:json" json:"config"`
	Status             ValidationStatus  `gorm:"default:'pending'" json:"status"`
	Results            ValidationResults `gorm:"type:json" json:"results,omitempty"`
	UserID             uint              `json:"user_id"`
	User               *User             `gorm:"foreignKey:UserID" json:"-"`
	CreatedAt          time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
}

func (v *Validation) BeforeCreate(tx *gorm.DB) error {
	if v.Status == "" {
		v.Status = ValidationStatusPending
	}
	return nil
}