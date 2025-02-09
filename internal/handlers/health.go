package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Health return the API server status and database status
func (app *App) Health(c *gin.Context) {
	dbHealth := "UP"

	dbErr := app.db.Ping()
	if dbErr != nil {
		dbHealth = "DOWN" + " ere:- " + dbErr.Error()
	}
	c.JSON(http.StatusOK, gin.H{
		"API service": "UP",
		"DB service":  dbHealth,
		"Timestamp":   time.Now().UTC(),
	})
}
