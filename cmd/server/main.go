package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DhruvikDonga/mock-ses/config"
	"github.com/DhruvikDonga/mock-ses/internal/database"
	"github.com/DhruvikDonga/mock-ses/internal/handlers"
	"github.com/DhruvikDonga/mock-ses/pkg/logger"
	"github.com/ardanlabs/conf"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Println("main api error :", err)
		os.Exit(1)
	}
}

func run() error {

	//setup config file
	cfg := config.Config{}
	var err error
	if err = conf.Parse(os.Args[1:], "CRUD", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("CRUD", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}

	log := logger.NewLogger(&cfg)
	log.Info("Mock SES service")
	log.Infof("main : Config :\n%v\n", out)

	//setup db
	db, err := database.Open(&cfg.DB)
	if err != nil {
		return errors.Wrap(err, "db connection failed")
	}

	//setup gin app
	gin.SetMode(cfg.Build)
	r := gin.Default()

	app := handlers.NewApp(r, db, log, &cfg)

	handlers.NewHandler(app)
	r.Run(cfg.Web.ServerURL)
	return nil
}
