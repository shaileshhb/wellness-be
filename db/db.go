package db

// Database will mimic a database
type Database struct{}

// NewDatabase will initialize a new database instance
func NewDatabase() *Database {
	db := &Database{}

	return db
}
