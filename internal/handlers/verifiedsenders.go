package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) AddVerifiedSender(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	_, err := app.db.Exec("INSERT INTO verified_senders (email) VALUES ($1) ON CONFLICT (email) DO NOTHING", req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add verified sender"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sender added successfully"})
}

func (app *App) DeleteVerifiedSender(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	res, err := app.db.Exec("DELETE FROM verified_senders WHERE email = $1", req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete verified sender"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sender removed successfully"})
}
