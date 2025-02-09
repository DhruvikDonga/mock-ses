package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Set email sending limits for a sender email
func (app *App) SetEmailLimits(c *gin.Context) {
	var req struct {
		SenderEmail string `json:"sender_email"`
		DailyQuota  int    `json:"daily_quota"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := app.db.Exec(
		"INSERT INTO email_limits (sender_email, daily_quota) VALUES ($1, $2) ON CONFLICT (sender_email) DO UPDATE SET daily_quota = EXCLUDED.daily_quota",
		req.SenderEmail, req.DailyQuota,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set email limits"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email limits updated successfully"})
}

// GetDailyQuota gets the dailyquota also checks if reset is neededs
func (app *App) GetDailyQuota(senderEmail string) (int, error) {
	var dailyQuota, sentToday int
	var lastReset time.Time

	err := app.db.QueryRow("SELECT daily_quota, sent_today, last_reset FROM email_limits WHERE sender_email = $1", senderEmail).
		Scan(&dailyQuota, &sentToday, &lastReset)
	if err != nil {
		return 0, err
	}

	// Check if the last reset was more than 24 hours ago
	if time.Since(lastReset) >= 24*time.Hour {
		sentToday = 0 // Reset count
		_, err = app.db.Exec("UPDATE email_limits SET sent_today = 0, last_reset = NOW() WHERE sender_email = $1", senderEmail)
		if err != nil {
			return 0, err
		}
	}

	return dailyQuota - sentToday, nil
}
