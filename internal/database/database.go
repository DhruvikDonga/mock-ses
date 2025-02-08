package database

import (
	"context"
	"net/url"

	"github.com/DhruvikDonga/mock-ses/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
)

// Open knows how to open a database connection based on the configuration.
func Open(cfg *config.DB) (*sqlx.DB, error) {

	// Query parameters.
	q := make(url.Values)
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	// Construct url.
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sqlx.DB) error {

	// Run a simple query to determine connectivity. The db has a "Ping" method
	// but it can false-positive when it was previously able to talk to the
	// database but the database has since gone away. Running this query forces a
	// round trip to the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}
