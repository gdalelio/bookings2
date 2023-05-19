package dbrepo

import (
	"database/sql"

	"github.com/gdalelio/bookings/internal/config"
	"github.com/gdalelio/bookings/internal/repository"
)
/*
---------------------------------------------------------|
|    Create the database connections for the specific    |
|   database flavor(s) - initial is PostgreSQL.          |
|--------------------------------------------------------|
*/

// postgresDBRepo struct for the repository with app (pointer to config) and database connection pool
type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewPostgersRepo creates a new PostgreSQL database repository
// returns the App Config and DB connection pool -
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

/*--------------- initial stub for MySQL as the Repository ---------------*
// mySQLDBRepo
type mySQLDBRepo struct {
	App *config.AppConfig
	DB  string
}

func NewMysSQLRepo(conn, a *config.AppConfig) repository.DatabaseRepo {
	// conn needs to be a reference pointer to the MySQL driver package similar to *sql.DB
	return &mySQLDBRepo{
		App, a,
		DB, conn
	}
}
/--------------- initial stub for MySQL as the Repository ---------------*/
