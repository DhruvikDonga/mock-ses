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
}
