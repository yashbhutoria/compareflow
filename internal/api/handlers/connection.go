package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/compareflow/compareflow/internal/models"
	"github.com/compareflow/compareflow/internal/services"
)

type ConnectionHandler struct {
	db      *gorm.DB
	service *services.ConnectionService
}

func NewConnectionHandler(db *gorm.DB) *ConnectionHandler {
	return &ConnectionHandler{
		db:      db,
		service: services.NewConnectionService(),
	}
}

type CreateConnectionRequest struct {
	Name   string                   `json:"name" binding:"required"`
	Type   models.ConnectionType    `json:"type" binding:"required"`
	Config models.ConnectionConfig  `json:"config" binding:"required"`
}

func (h *ConnectionHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")

	var connections []models.Connection
	if err := h.db.Where("user_id = ?", userID).Find(&connections).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch connections"})
		return
	}

	c.JSON(http.StatusOK, connections)
}

func (h *ConnectionHandler) Get(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
		return
	}

	var connection models.Connection
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&connection).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Connection not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch connection"})
		}
		return
	}

	c.JSON(http.StatusOK, connection)
}

func (h *ConnectionHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	connection := &models.Connection{
		Name:   req.Name,
		Type:   req.Type,
		Config: req.Config,
		UserID: userID,
	}

	if err := h.db.Create(connection).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create connection"})
		return
	}

	c.JSON(http.StatusCreated, connection)
}

func (h *ConnectionHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
		return
	}

	var connection models.Connection
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&connection).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Connection not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch connection"})
		}
		return
	}

	var req CreateConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	connection.Name = req.Name
	connection.Type = req.Type
	connection.Config = req.Config

	if err := h.db.Save(&connection).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update connection"})
		return
	}

	c.JSON(http.StatusOK, connection)
}

func (h *ConnectionHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
		return
	}

	result := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Connection{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete connection"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Connection not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connection deleted successfully"})
}

func (h *ConnectionHandler) Test(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
		return
	}

	var connection models.Connection
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&connection).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Connection not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch connection"})
		}
		return
	}

	// Test the connection
	if err := h.service.TestConnection(&connection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Connection test successful",
	})
}

func (h *ConnectionHandler) GetTables(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
		return
	}

	var connection models.Connection
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&connection).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Connection not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch connection"})
		}
		return
	}

	// Get tables
	tables, err := h.service.GetTables(&connection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tables": tables,
	})
}

func (h *ConnectionHandler) GetColumns(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
		return
	}

	tableName := c.Param("table")
	if tableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table name is required"})
		return
	}

	var connection models.Connection
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&connection).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Connection not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch connection"})
		}
		return
	}

	// Get columns
	columns, err := h.service.GetColumns(&connection, tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"columns": columns,
	})
}