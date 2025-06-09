package migration

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type GooseMigration struct {
	log *slog.Logger
}

func NewGooseMigration(log *slog.Logger) *GooseMigration {
	return &GooseMigration{
		log: log,
	}
}

func (g *GooseMigration) Migrate(dbURL string, folderPath string) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		g.log.Error("Failed to open database connection", "error", err, "dbURL", dbURL)
		return err
	}
	defer db.Close()

	g.log.Info("Starting database migration", "folderPath", folderPath)
	goose.SetLogger(newGooseSlogLogger(g.log))
	if err := goose.Up(db, folderPath); err != nil {
		return err
	}

	g.log.Info("Database migration completed successfully", "folderPath", folderPath)
	return nil
}

func newGooseSlogLogger(l *slog.Logger) goose.Logger {
	return &gooseSlogLogger{
		log: l,
	}
}

type gooseSlogLogger struct {
	log *slog.Logger
}

func (l *gooseSlogLogger) Printf(format string, v ...any) {
	l.log.Info(fmt.Sprintf(format, v...))
}

func (l *gooseSlogLogger) Fatalf(format string, v ...any) {
	l.log.Error(fmt.Sprintf(format, v...))
}
