package handlers

import (
	"math/rand"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// EmailStatus represents possible statuses of an email
const (
	StatusSend          = "Send"
	StatusRenderingFail = "RenderingFailure"
	StatusReject        = "Reject"
	StatusDelivery      = "Delivery"
	StatusBounce        = "Bounce"
	StatusComplaint     = "Complaint"
	StatusDeliveryDelay = "DeliveryDelay"
	StatusSubscription  = "Subscription"
	StatusOpen          = "Open"
	StatusClick         = "Click"
)

// EmailRequest represents the structure of the SES SendEmail request payload
type EmailRequest struct {
	FromEmailAddress string `json:"FromEmailAddress" binding:"required"`
	Destination      struct {
		ToAddresses  []string `json:"ToAddresses" binding:"required"`
		CcAddresses  []string `json:"CcAddresses,omitempty"`
		BccAddresses []string `json:"BccAddresses,omitempty"`
	} `json:"Destination" binding:"required"`
	Content struct {
		Simple struct {
			Subject struct {
				Data    string `json:"Data" binding:"required"`
				Charset string `json:"Charset,omitempty"`
			} `json:"Subject" binding:"required"`
			Body struct {
				Text struct {
					Data    string `json:"Data" binding:"required"`
					Charset string `json:"Charset,omitempty"`
				} `json:"Text,omitempty"`
				Html struct {
					Data    string `json:"Data" binding:"required"`
					Charset string `json:"Charset,omitempty"`
				} `json:"Html,omitempty"`
			} `json:"Body" binding:"required"`
		} `json:"Simple" binding:"required"`
	} `json:"Content" binding:"required"`
}

// SendEmail will process the email request and won't send email but will return a message id and mock the result
func (app *App) SendEmail(c *gin.Context) {

	// Validate AWS Authorization headers
	if !validateAWSHeaders(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid AWS authorization headers"})
		return
	}
	var req EmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sender := req.FromEmailAddress
	messageID := uuid.New().String()
	tx, err := app.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database transaction failed"})
		return
	}

	// Validate domain verification
	if !app.isVerifiedSender(sender) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Sender email not verified."})
		tx.Rollback()
		return
	}

	// Check Daily Limit
	quota, err := app.GetDailyQuota(sender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve sender limit"})
		tx.Rollback()
		return
	}
	if quota <= 0 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Daily email sending limit reached."})
		tx.Rollback()
		return
	}

	// Insert Email
	var emailID string
	err = tx.QueryRow(
		"INSERT INTO emails (message_id, subject, body, from_email) VALUES ($1, $2, $3, $4) RETURNING email_id",
		messageID, req.Content.Simple.Subject.Data, req.Content.Simple.Body.Html.Data, req.FromEmailAddress,
	).Scan(&emailID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert email"})
		return
	}

	recipients := append(req.Destination.ToAddresses, req.Destination.CcAddresses...)
	recipients = append(recipients, req.Destination.BccAddresses...)

	statuses := []string{
		StatusSend, StatusRenderingFail, StatusReject, StatusDelivery,
		StatusBounce, StatusComplaint, StatusDeliveryDelay, StatusSubscription,
		StatusOpen, StatusClick,
	}

	var suppressedRecipients []string

	for _, recipient := range recipients {

		// Check suppression list
		if app.isSuppressed(recipient) {
			_, err = tx.Exec(
				"INSERT INTO delivery_logs (email_id, message_id, recipient_email, status, response) VALUES ($1, $2, $3, $4, $5)",
				emailID, messageID, recipient, "Suppressed", "Recipient is suppressed and email was not sent.",
			)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert suppression log"})
				return
			}
			suppressedRecipients = append(suppressedRecipients, recipient)
			continue
		}

		//Email must be send from here we are skipping that thing actually and use the email status random

		//Update delivery_logs
		randomStatus := statuses[rand.Intn(len(statuses))]
		_, err = tx.Exec(
			"INSERT INTO delivery_logs (email_id,message_id,recipient_email, status, response) VALUES ($1, $2, $3,$4,$5)",
			emailID, messageID, recipient, randomStatus, getMockResponse(randomStatus),
		)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert delivery log"})
			return
		}
	}
	// Update daily email count
	_, err = tx.Exec("UPDATE email_limits SET sent_today = sent_today + 1 WHERE sender_email = $1", sender)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update email limit"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"MessageId": messageID})
}

// isVerifiedSender checks if the sender's email is verified in the system
func (app *App) isVerifiedSender(sender string) bool {
	var count int
	err := app.db.QueryRow("SELECT COUNT(*) FROM verified_senders WHERE email = $1", sender).Scan(&count)
	if err != nil || count == 0 {
		return false
	}
	return true
}

// isSuppressed checks if the recipient's email is in the suppression list
func (app *App) isSuppressed(recipient string) bool {
	var count int
	err := app.db.QueryRow("SELECT COUNT(*) FROM suppression_list WHERE email = $1", recipient).Scan(&count)
	if err != nil || count > 0 {
		return true
	}
	return false
}

func validateAWSHeaders(c *gin.Context) bool {
	authHeader := c.GetHeader("Authorization")
	dateHeader := c.GetHeader("X-Amz-Date")

	if authHeader == "" || dateHeader == "" {
		return false
	}

	if !strings.HasPrefix(authHeader, "AWS4-HMAC-SHA256") {
		return false
	}

	return true
}

func getMockResponse(status string) string {
	switch status {
	case StatusSend:
		return "Email has been sent successfully."
	case StatusRenderingFail:
		return "Failed to render email content."
	case StatusReject:
		return "Email was rejected due to policy."
	case StatusDelivery:
		return "Email was successfully delivered."
	case StatusBounce:
		return "Email bounced back. Invalid recipient address."
	case StatusComplaint:
		return "Recipient marked email as spam."
	case StatusDeliveryDelay:
		return "Email is delayed due to network issues."
	case StatusSubscription:
		return "Recipient has unsubscribed."
	case StatusOpen:
		return "Email was opened by recipient."
	case StatusClick:
		return "Recipient clicked on a link in the email."
	default:
		return "Unknown email status."
	}
}
