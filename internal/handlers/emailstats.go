package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LogResponse struct {
	Logs []Log
}

type Log struct {
	RecipientEmail string    `json:"recipient_email"`
	SenderEmail    string    `json:"sender_email"`
	Status         string    `json:"status"`
	Response       string    `json:"response"`
	CreatedAt      time.Time `json:"created_at"`
}

func (app *App) GetEmailStats(c *gin.Context) {
	messageID := c.Param("message_id")

	var logs []Log

	query := `SELECT recipient_email,sender_email, status, response, created_at 
              FROM delivery_logs 
              WHERE message_id = $1 
              ORDER BY created_at DESC`

	rows, err := app.db.Query(query, messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var log Log
		if err := rows.Scan(&log.RecipientEmail, &log.SenderEmail, &log.Status, &log.Response, &log.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning data"})
			return
		}
		logs = append(logs, log)
	}

	if len(logs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No records found for message_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message_id": messageID, "logs": LogResponse{Logs: logs}})
}
