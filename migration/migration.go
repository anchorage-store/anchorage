// Package migrations provides a set of functions to perform db migrations
// and track their runs against the database.
package migration

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"slices"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// Migrate applies the mgirations in the given directory using the connection.
//
// It ensures that a `migration_log` table is created to keep track of and checkpoint
// the state of the database in case of a bad/failing migration.
//
// Returns the name of the migrations that were applied.
func Migrate(ctx context.Context, db *sqlx.DB, migrationDir fs.FS) ([]string, error) {
	// Determine the files in the directory as the
	// super set of all migrations that need to be run.
	var migrations []string
	if err := fs.WalkDir(migrationDir, ".", func(path string, entry fs.DirEntry, err error) error {
		// Skip anything that isn't a file
		if entry.IsDir() {
			return nil
		}
		// Skip anything that isn't a `.sql` file
		if !strings.HasSuffix(entry.Name(), ".sql") {
			return nil
		}

		migrations = append(migrations, path)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("error walking migration directory: %s", err)
	}

	// Migrations are applied in alphabetical order of their names
	slices.Sort(migrations)

	// The migration log table needs to be created to checkpoint migrations as they are applied
	if err := ensureMigrationLog(ctx, db); err != nil {
		return nil, fmt.Errorf("error ensuring migration log: %s", err)
	}

	// For each migration in the directory, check if it has already been applied.
	// If not, read it, apply it, and then checkpoint it in the log.
	//
	// Otherwise go to the next one.
	existing, err := AppliedMigrations(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("error getting existing migrations: %s", err)
	}
	var applied []string
	for _, m := range migrations {
		if slices.Contains(existing, m) {
			continue
		}

		f, err := migrationDir.Open(m)
		if err != nil {
			return nil, fmt.Errorf("error opening migration '%s': %s", m, err)
		}
		defer f.Close()
		byts, err := io.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("error reading migration '%s': %s", m, err)
		}

		if _, err := db.ExecContext(ctx, string(byts)); err != nil {
			return nil, fmt.Errorf("error appying migration '%s': %s", m, err)
		}

		if err := checkpointMigration(ctx, db, m); err != nil {
			return nil, fmt.Errorf("error appying migration '%s': %s", m, err)
		}

		applied = append(applied, m)
	}

	return applied, nil
}

type migration struct {
	ID        uint      `db:"id"`
	Path      string    `db:"path"`
	CreatedAt time.Time `db:"created_at"`
}

// Runs sql to make sure the `migration_log` table is created in the database.
func ensureMigrationLog(ctx context.Context, dbx *sqlx.DB) error {
	q := `
    CREATE TABLE IF NOT EXISTS migration_log (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        path TEXT NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    `

	if _, err := dbx.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("error executing migration_log statement: %s", err)
	}

	return nil
}

// AppliedMigrations returns the migrations currently in the log.
func AppliedMigrations(ctx context.Context, dbx *sqlx.DB) ([]string, error) {
	q := "SELECT * FROM migration_log ORDER BY id ASC;"

	var migrations []migration
	if err := dbx.SelectContext(ctx, &migrations, q); err != nil {
		return nil, fmt.Errorf("error selecting migiration logs: %s", err)
	}

	paths := make([]string, 0, len(migrations))
	for _, m := range migrations {
		paths = append(paths, m.Path)
	}
	return paths, nil
}

// Inserts a log into the migration table, signifying that it has been run.
func checkpointMigration(ctx context.Context, db *sqlx.DB, name string) error {
	q := "INSERT INTO migration_log (path) VALUES (?);"

	if _, err := db.ExecContext(ctx, q, name); err != nil {
		return fmt.Errorf("error inserting path into migration_log: %s", err)
	}

	return nil
}
