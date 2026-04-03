// Package app - application package
package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aditya-sutar-45/finance-dashboard-api/internal/database"
	_ "github.com/lib/pq"
)

type App struct {
	port   string
	router http.Handler
	DB     *database.Queries
}

func New(port string, dbURL string) (*App, error) {
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	db := database.New(conn)

	app := &App{
		port:   port,
		router: loadRoutes(db),
		DB:     db,
	}

	return app, nil
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    a.port,
		Handler: a.router,
	}

	log.Println("starting server on PORT", a.port)

	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}
}
