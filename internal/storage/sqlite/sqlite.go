package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/theeeep/go-rest-api/internal/config"
)

type Sqlite struct {
	Db *sql.DB
}

// New initializes a new SQLite database connection and ensures the students table exists.
// It takes a configuration object `cfg` as input and returns a pointer to the Sqlite struct and an error.
func New(cfg *config.Config) (*Sqlite, error) {
	// Open a connection to the SQLite database using the storage path from the configuration.
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	// Execute a SQL statement to create the `students` table if it doesn't already exist.
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS students (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      name TEXT NOT NULL,
      email TEXT NOT NULL,
      age INTEGER NOT NULL
    );
  `)

	if err != nil {
		return nil, err
	}

	// Return a new Sqlite instance with the initialized database connection.
	return &Sqlite{
		Db: db,
	}, nil
}

// CreateStudent inserts a new student record into the `students` table.
// It takes the student's name, email, and age as input and returns the ID of the newly inserted record and an error.
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	// Prepare the SQL statement for inserting a new student.
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close() // Ensure the statement is closed after execution.

	// Execute the prepared statement with the provided student details.
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	// Retrieve the ID of the newly inserted student.
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Return the ID of the newly inserted student.
	return lastId, nil
}
