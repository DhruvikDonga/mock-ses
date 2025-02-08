package config

type (
	Config struct {
		Build string `conf:"default:release,env:BUILD"`
		Web   Web
		DB    DB
	}
	Web struct {
		ServerURL string `conf:"default:0.0.0.0:8080,env:API_HOST"`
	}
	DB struct {
		User     string `conf:"default:postgres,env:DB_USER"`
		Password string `conf:"default:xxxxxxxxx,noprint,env:DB_PASSWORD"`
		Host     string `conf:"default:xxxxxxxxx:5432,env:DB_HOST"`
		Name     string `conf:"default:xxxxxxxxx,env:DB_NAME"`
		Port     string `conf:"default:5432,env:DB_PORT"`
	}
)
