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
	InsertHistory(wtCmd string, cmds []string) error
	InsertFavourite(name string, wtCmd string, cmds []string) error
	ReadHistory() (Histories, error)
	ReadFavourite() (Favourites, error)
	DeleteFavourite(id int, name string) error
}

// Repository represents a repository for storing and retrieving history of executed commands
type Repository struct {
	db *sql.DB
}

// Histories represents a list of History returned from database
type Histories []struct {
	model.History
}

// Favourites represents a list of Favourite returned from database
type Favourites []struct {
	model.Favourite
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

// ReadFavourite reads all favourite entries from the database and returns them as a Favourites slice
func (r *Repository) ReadFavourite() (Favourites, error) {
	stmt := jetSqlite.SELECT(jetTable.Favourite.AllColumns).FROM(jetTable.Favourite)

	var f Favourites
	err := stmt.Query(r.db, &f)
	if err != nil {
		return nil, fmt.Errorf("failed to read FAVOURITE: %v", err)
	}

	return f, nil
}

// ReadHistory reads all history entries from the database and returns them as a Histories slice
func (r *Repository) ReadHistory() (Histories, error) {
	stmt := jetSqlite.SELECT(jetTable.History.AllColumns).FROM(jetTable.History).ORDER_BY(jetTable.History.ExecutedAt.DESC())

	var h Histories
	err := stmt.Query(r.db, &h)
	if err != nil {
		return nil, fmt.Errorf("failed to read HISTORY: %v", err)
	}

	return h, nil
}

// InsertFavourite insert a favourite entry into the database
func (r *Repository) InsertFavourite(name, wtCmd string, cmds []string) error {
	stmt := jetTable.Favourite.INSERT(
		jetTable.Favourite.Name,
		jetTable.Favourite.Wtcmd,
		jetTable.Favourite.Cmds).
		MODEL(model.Favourite{
			Name:  name,
			Wtcmd: wtCmd,
			Cmds:  strings.Join(cmds, ","),
		})

	_, err := stmt.Exec(r.db)
	if err != nil {
		return fmt.Errorf("failed to insert FAVOURITE: %v", err)
	}
	return nil
}

// InsertHistory insert a history entry into the database
func (r *Repository) InsertHistory(wtCmd string, cmds []string) error {
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

// DeleteFavourite deletes a favourite entry from the database by its id
func (r *Repository) DeleteFavourite(id int, name string) error {
	stmt := jetTable.Favourite.DELETE().WHERE(jetTable.Favourite.ID.IN(jetSqlite.Int(int64(id))))

	_, err := stmt.Exec(r.db)
	if err != nil {
		return fmt.Errorf("failed to delete %s: %v", name, err)
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
