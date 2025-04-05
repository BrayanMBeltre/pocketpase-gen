package commands

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func GetCollections() ([]core.Collection, error) {
	db, err := sql.Open("sqlite3", FlagDBPath)

	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite3 database: %w", err)
	}

	open := dbx.NewFromDB(db, "sqlite3")

	var collections []core.Collection

	err = open.NewQuery("SELECT * FROM _collections").All(&collections)

	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	return collections, nil
}
