package storage

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(path string) {
	var err error
	DB, err = sql.Open("sqlite", path)
	if err != nil {
		log.Fatal(err)
	}

	createTables()
}

func createTables() {
	folderTable := `
	CREATE TABLE IF NOT EXISTS folders (
		id TEXT PRIMARY KEY,
		name TEXT,
		parent_id TEXT,
		created_at INTEGER
	);`

	fileTable := `
	CREATE TABLE IF NOT EXISTS files (
		id TEXT PRIMARY KEY,
		name TEXT,
		size INTEGER,
		type TEXT,
		folder_id TEXT,
		status TEXT,
		created_at INTEGER
	);`

	jobTable := `
	CREATE TABLE IF NOT EXISTS jobs (
		id TEXT PRIMARY KEY,
		type TEXT,
		status TEXT,
		payload TEXT,
		retry_count INTEGER,
		created_at INTEGER
	);`

	_, err := DB.Exec(folderTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(fileTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(jobTable)
	if err != nil {
		log.Fatal(err)
	}
}
