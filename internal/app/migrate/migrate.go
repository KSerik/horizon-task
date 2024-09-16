package migrate

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/pressly/goose/v3"
)

type App struct {
	conn *sql.DB
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func NewApp(conn *sql.DB) *App {
	return &App{
		conn: conn,
	}
}

func (a *App) Up() error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("clickhouse"); err != nil {
		return fmt.Errorf("error setting dialect: %w", err)
	}

	if err := goose.Up(a.conn, "migrations"); err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	return nil
}
