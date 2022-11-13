// File: forum/cmd/api/main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"forum.kevin.net/internal/data"
	_ "github.com/lib/pq"
)

// configuration settings
type config struct {
	port int
	env  string
	db   struct { // development, staging, production, etc.
		dsn          string
		maxOpenConns int
		maxIdleConns int
		MaxIdleTime  string
	}
}

// The application version number
const version = "1.0.0"

// Dependency Injections
type application struct {
	config config
	logger *log.Logger
	models data.Models
}

// main
func main() {
	//initialize the config
	var cfg config

	// the settings that are needed to populate our config
	cfg.port = 4000
	cfg.env = "development"
	cfg.db.dsn = os.Getenv("TESTFORUM_DB_DSN")
	cfg.db.maxOpenConns = 25
	cfg.db.maxIdleConns = 25
	cfg.db.MaxIdleTime = "15m"

	//creating logger to log issues or state changes
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//creating connection
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	//ensuring that the connection to the database is closed
	defer db.Close()

	logger.Println("database connection pool established")

	//initializing the app struct
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	//initializing http server dependencies
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	//staring the web server
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

// OpenDB() function returns a *sql.DB connection pool
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	//Creating a context with a 5-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
