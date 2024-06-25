package migration_test

import (
	"context"
	"embed"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"github.com/anchorage-store/anchorage/migration"
)

//go:embed test_migrations/*
var migrations embed.FS

//go:embed test_migrations_bad/*
var badMigrations embed.FS

func TestMigrate_CreatesLogTable(t *testing.T) {
	dbx := sqlx.MustOpen("sqlite3", "file::memory:")
	t.Cleanup(func() {
		dbx.Close()
	})
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if _, err := migration.Migrate(ctx, dbx, migrations); err != nil {
		t.Fatalf("unexpected error running migrations: %s", err)
	}

	// Check for log table
	rows, err := dbx.QueryContext(ctx, "SELECT * FROM migration_log;")
	if err != nil {
		t.Fatalf("unexpected error querying the migration log: %s", err)
	}
	rows.Close()
}

func TestMigrate_InsertsLogs(t *testing.T) {
	dbx := sqlx.MustOpen("sqlite3", "file::memory:")
	t.Cleanup(func() {
		dbx.Close()
	})
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if _, err := migration.Migrate(ctx, dbx, migrations); err != nil {
		t.Fatalf("unexpected error running migrations: %s", err)
	}

	// Check that both migrations were logged
	got, err := migration.AppliedMigrations(ctx, dbx)
	if err != nil {
		t.Fatalf("unexpected error fetching applied migrations: %s", err)
	}

	want := []string{
		"test_migrations/2024_06_22_1_users.sql",
		"test_migrations/2024_06_22_2_logins.sql",
	}
	assert.Equal(t, got, want)
}

func TestMigrate_StopsAfterBadFile(t *testing.T) {
	dbx := sqlx.MustOpen("sqlite3", "file::memory:")
	t.Cleanup(func() {
		dbx.Close()
	})
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if _, err := migration.Migrate(ctx, dbx, badMigrations); err == nil {
		t.Fatal("did not get error after migrating")
	}

	// Check that both migrations were logged
	got, err := migration.AppliedMigrations(ctx, dbx)
	if err != nil {
		t.Fatalf("unexpected error fetching applied migrations: %s", err)
	}

	want := []string{
		"test_migrations_bad/2024_06_22_1_users.sql",
	}
	assert.Equal(t, got, want)
}
