package driver

import (
	"database/sql"
	"log"
	"time"

	//anonymous import the database libraries to call run their init methods
	//to work with the database
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

/*------------ Set Up Database Constants ----------------*/
//initializing and creating an emptu connection
var dbConn = &DB{}

// maxOpenDbConn will limit number of open connections to the database
const maxOpenDbConn = 10

// maxIdleDbConn wll limit the number of idle connections in the pool
const maxIdleDbConn = 5

// maxDbLifetime will limit the time a connection can be connected (5 minutes)
const maxDbLifetime = 5 * time.Minute

/*--------------  Database Constants ----------------*/

// ConnectSQL creates database pool for database using the data source name (dsn) [connection string]
func ConnectSQL(dsn string) (*DB, error) {
	dbase, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	//connection has been successful - set the parameters from constants on dbase
	dbase.SetMaxOpenConns(maxOpenDbConn)
	dbase.SetConnMaxIdleTime(maxIdleDbConn)
	dbase.SetConnMaxLifetime(maxDbLifetime)

	//assign the dbConn to the dbase
	dbConn.SQL = dbase

	//test the database connectivity
	err = testDB(dbase)
	if err != nil {
		return dbConn, err
	}
	log.Println("connection test completed successfully")
	return dbConn, err
}

// testDB tries to ping the database
func testDB(dbase *sql.DB) error {
	err := dbase.Ping()
	if err != nil {
		return nil
	}
	return err
}

// NewDatabase creates a new database using data source name (dsn) [connection string]
// returns a pointer to sql.DB
func NewDatabase(dsn string) (*sql.DB, error) {
	database, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err

	}

	if err = database.Ping(); err != nil {
		return nil, err
	}

	//return the database
	return database, nil
}
