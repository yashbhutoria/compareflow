package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/compareflow/compareflow/internal/models"
)

type ValidationHandler struct {
	db *gorm.DB
}

func NewValidationHandler(db *gorm.DB) *ValidationHandler {
	return &ValidationHandler{db: db}
}

type CreateValidationRequest struct {
	Name               string                   `json:"name" binding:"required"`
	SourceConnectionID uint                     `json:"source_connection_id" binding:"required"`
	TargetConnectionID uint                     `json:"target_connection_id" binding:"required"`
	Config             models.ValidationConfig  `json:"config"`
}

func (h *ValidationHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")

	var validations []models.Validation
	if err := h.db.Preload("SourceConnection").Preload("TargetConnection").
		Where("user_id = ?", userID).Find(&validations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch validations"})
		return
	}

	c.JSON(http.StatusOK, validations)
}

func (h *ValidationHandler) Get(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid validation ID"})
		return
	}

	var validation models.Validation
	if err := h.db.Preload("SourceConnection").Preload("TargetConnection").
		Where("id = ? AND user_id = ?", id, userID).First(&validation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Validation not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch validation"})
		}
		return
	}

	c.JSON(http.StatusOK, validation)
}

func (h *ValidationHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify that both connections belong to the user
	var sourceConn, targetConn models.Connection
	if err := h.db.Where("id = ? AND user_id = ?", req.SourceConnectionID, userID).First(&sourceConn).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source connection"})
		return
	}
	if err := h.db.Where("id = ? AND user_id = ?", req.TargetConnectionID, userID).First(&targetConn).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target connection"})
		return
	}

	validation := &models.Validation{
		Name:               req.Name,
		SourceConnectionID: req.SourceConnectionID,
		TargetConnectionID: req.TargetConnectionID,
		Config:             req.Config,
		UserID:             userID,
		Status:             models.ValidationStatusPending,
	}

	if err := h.db.Create(validation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create validation"})
		return
	}

	// Load associations
	h.db.Preload("SourceConnection").Preload("TargetConnection").First(validation, validation.ID)

	c.JSON(http.StatusCreated, validation)
}

func (h *ValidationHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid validation ID"})
		return
	}

	var validation models.Validation
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&validation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Validation not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch validation"})
		}
		return
	}

	var req CreateValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify that both connections belong to the user
	var sourceConn, targetConn models.Connection
	if err := h.db.Where("id = ? AND user_id = ?", req.SourceConnectionID, userID).First(&sourceConn).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source connection"})
		return
	}
	if err := h.db.Where("id = ? AND user_id = ?", req.TargetConnectionID, userID).First(&targetConn).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target connection"})
		return
	}

	validation.Name = req.Name
	validation.SourceConnectionID = req.SourceConnectionID
	validation.TargetConnectionID = req.TargetConnectionID
	validation.Config = req.Config

	if err := h.db.Save(&validation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update validation"})
		return
	}

	// Load associations
	h.db.Preload("SourceConnection").Preload("TargetConnection").First(&validation, validation.ID)

	c.JSON(http.StatusOK, validation)
}

func (h *ValidationHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid validation ID"})
		return
	}

	result := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Validation{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete validation"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Validation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Validation deleted successfully"})
}

func (h *ValidationHandler) Run(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid validation ID"})
		return
	}

	var validation models.Validation
	if err := h.db.Preload("SourceConnection").Preload("TargetConnection").
		Where("id = ? AND user_id = ?", id, userID).First(&validation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Validation not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch validation"})
		}
		return
	}

	// Update status to running
	validation.Status = models.ValidationStatusRunning
	h.db.Save(&validation)

	// TODO: Implement actual validation logic
	// For now, we'll just simulate a successful validation
	validation.Status = models.ValidationStatusCompleted
	validation.Results = models.ValidationResults{
		"summary": map[string]interface{}{
			"total_records": 1000,
			"matched":       950,
			"mismatched":    50,
			"success_rate":  95.0,
		},
	}
	h.db.Save(&validation)

	c.JSON(http.StatusOK, validation)
}

func (h *ValidationHandler) Status(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid validation ID"})
		return
	}

	var validation models.Validation
	if err := h.db.Select("id", "name", "status", "updated_at").
		Where("id = ? AND user_id = ?", id, userID).First(&validation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Validation not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch validation status"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         validation.ID,
		"name":       validation.Name,
		"status":     validation.Status,
		"updated_at": validation.UpdatedAt,
	})
}