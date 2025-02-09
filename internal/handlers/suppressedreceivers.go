package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) AddToSuppressionList(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	_, err := app.db.Exec("INSERT INTO suppression_list (email) VALUES ($1) ON CONFLICT (email) DO NOTHING", req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to suppression list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email added to suppression list"})
}

func (app *App) DeleteFromSuppressionList(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	res, err := app.db.Exec("DELETE FROM suppression_list WHERE email = $1", req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete from suppression list"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email removed from suppression list"})
}
