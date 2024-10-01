package repository

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"mpwt/internal/repository/.gen/model"
	jetTable "mpwt/internal/repository/.gen/table"

	jetSqlite "github.com/go-jet/jet/v2/sqlite"
)

//go:embed assets/*
var assets embed.FS

// IRepository is the interface for the repository
type IRepository interface {
	InsertHistory([]string, string) error
	ReadHistory() (Histories, error)
}

// Repository represents a repository for storing and retrieving history of executed commands
type Repository struct {
	db *sql.DB
}

// Histories represents a list of History returned from database
type Histories []struct {
	model.History
}

// NewDbConn creates a new connection to the SQLite database at the specified filepath
// If not exist, it will create a new database at the specified filepath
func NewDbConn(filepath string) (*Repository, error) {
	// Check if database file exists, if not create a new one
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		err := createDatabase(filepath)
		if err != nil {
			return nil, fmt.Errorf("failed to create new database: %v", err)
		}
	}

	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return &Repository{db: db}, nil
}

// Close closes the database connection
func (r *Repository) Close() {
	r.db.Close()
}

// ReadHistory reads all history entries from the database and returns them as a Histories slice
func (r *Repository) ReadHistory() (Histories, error) {
	stmt := jetSqlite.SELECT(jetTable.History.AllColumns).FROM(jetTable.History)

	var h Histories
	err := stmt.Query(r.db, &h)
	if err != nil {
		return nil, fmt.Errorf("failed to read HISTORY: %v", err)
	}

	return h, nil
}

// InsertHistory insert a history entry into the database
func (r *Repository) InsertHistory(cmds []string, wtCmd string) error {
	stmt := jetTable.History.INSERT(
		jetTable.History.ExecutedAt,
		jetTable.History.Cmds,
		jetTable.History.PaneCount,
		jetTable.History.Wtcmd).
		MODEL(model.History{
			ExecutedAt: time.Now(),
			Cmds:       strings.Join(cmds, ","),
			PaneCount:  int32(len(cmds)),
			Wtcmd:      wtCmd,
		})

	_, err := stmt.Exec(r.db)
	if err != nil {
		return fmt.Errorf("failed to insert HISTORY: %v", err)
	}
	return nil
}

// createDatabase creates a new database
func createDatabase(filepath string) error {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return fmt.Errorf("failed to open database file: %v", err)
	}
	defer db.Close()

	// Read content of init.sql file
	initSqlBytes, err := assets.ReadFile("assets/init.sql")
	if err != nil {
		return fmt.Errorf("failed to read init.sql file: %v", err)
	}

	// Execute the SQL statement to create the tables
	_, err = db.Exec(string(initSqlBytes))
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}
