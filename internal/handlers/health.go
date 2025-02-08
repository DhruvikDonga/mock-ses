package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) Health(c *gin.Context) {
	dbHealth := "UP"

	dbErr := app.db.Ping()
	if dbErr != nil {
		dbHealth = "DOWN" + " ere:- " + dbErr.Error()
	}
	c.JSON(http.StatusOK, gin.H{
		"API service": "UP",
		"DB Health":   dbHealth,
	})
}
