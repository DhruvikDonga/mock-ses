package handlers

import (
	"github.com/DhruvikDonga/mock-ses/config"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type App struct {
	r      *gin.Engine
	config *config.Config
	db     *sqlx.DB
	log    *zap.SugaredLogger
}

func NewApp(r *gin.Engine, db *sqlx.DB, zapLog *zap.SugaredLogger, config *config.Config) *App {
	return &App{
		r:      r,
		db:     db,
		log:    zapLog,
		config: config,
	}
}

func NewHandler(app *App) {

	//Health Check API
	app.r.GET("/ping", app.Health)

	//SendMail mock API
	app.r.POST("/v2/email/outbound-emails", app.SendEmail)

	//SetEmailLimits
	app.r.POST("/email-limits", app.SetEmailLimits)

	//Add/Remove Verified Emails
	app.r.POST("/verified-senders", app.AddVerifiedSender)
	app.r.DELETE("/verified-senders", app.DeleteVerifiedSender)

	//Add/Remove Suppressed Emails
	app.r.POST("/suppression-list", app.AddToSuppressionList)
	app.r.DELETE("/suppression-list", app.DeleteFromSuppressionList)

}
