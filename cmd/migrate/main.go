package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Migration struct {
	Version   string
	Name      string
	Filename  string
	Checksum  string
	AppliedAt *time.Time
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	db, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := ensureMigrationTable(db); err != nil {
		log.Fatalf("Failed to create migration table: %v", err)
	}

	// Determine command
	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "up":
		if err := runMigrations(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	case "status":
		if err := printStatus(db); err != nil {
			log.Fatalf("Failed to get status: %v", err)
		}
	default:
		fmt.Println("Usage: migrate [up|status]")
		fmt.Println("  up     - Run all pending migrations (default)")
		fmt.Println("  status - Show migration status")
		os.Exit(1)
	}
}

func connectDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		sslmode := os.Getenv("DB_SSLMODE")

		if host == "" {
			host = "localhost"
		}
		if port == "" {
			port = "5432"
		}
		if sslmode == "" {
			sslmode = "disable"
		}

		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode,
		)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot reach database: %w", err)
	}
	return db, nil
}

// ensureMigrationTable creates the tracking table if it doesn't exist
func ensureMigrationTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS navisha_english_migration (
			id          SERIAL PRIMARY KEY,
			version     VARCHAR(20)  NOT NULL UNIQUE,
			name        VARCHAR(255) NOT NULL,
			filename    VARCHAR(255) NOT NULL,
			checksum    VARCHAR(64)  NOT NULL,
			applied_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

// runMigrations finds and applies all pending migration files
func runMigrations(db *sql.DB) error {
	files, err := collectMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to collect migration files: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No migration files found in migrations/schema or migrations/seeds")
		return nil
	}

	applied, err := appliedMigrations(db)
	if err != nil {
		return err
	}

	pending := 0
	for _, f := range files {
		if _, ok := applied[f.Version]; ok {
			continue
		}

		content, err := os.ReadFile(f.Filename)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", f.Filename, err)
		}

		checksum := fmt.Sprintf("%x", sha256.Sum256(content))

		fmt.Printf("Applying %s_%s ... ", f.Version, f.Name)

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to apply %s: %w", f.Filename, err)
		}

		_, err = tx.Exec(`
			INSERT INTO navisha_english_migration (version, name, filename, checksum)
			VALUES ($1, $2, $3, $4)`,
			f.Version, f.Name, filepath.Base(f.Filename), checksum,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", f.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", f.Version, err)
		}

		fmt.Println("OK")
		pending++
	}

	if pending == 0 {
		fmt.Println("Nothing to migrate — all migrations already applied.")
	} else {
		fmt.Printf("\n%d migration(s) applied successfully.\n", pending)
	}

	return nil
}

// printStatus shows which migrations have been applied and which are pending
func printStatus(db *sql.DB) error {
	files, err := collectMigrationFiles()
	if err != nil {
		return err
	}

	applied, err := appliedMigrations(db)
	if err != nil {
		return err
	}

	fmt.Printf("\n%-6s %-20s %-40s %-10s %s\n", "ID", "Version", "Name", "Status", "Applied At")
	fmt.Println(strings.Repeat("-", 100))

	for i, f := range files {
		m, ok := applied[f.Version]
		status := "pending"
		appliedAt := ""
		if ok {
			status = "applied"
			if m.AppliedAt != nil {
				appliedAt = m.AppliedAt.Format("2006-01-02 15:04:05")
			}
		}
		fmt.Printf("%-6d %-20s %-40s %-10s %s\n", i+1, f.Version, f.Name, status, appliedAt)
	}

	total := len(files)
	appliedCount := len(applied)
	fmt.Printf("\n%d/%d migrations applied.\n", appliedCount, total)
	return nil
}

// collectMigrationFiles scans migrations/schema and migrations/seeds, returns sorted list
func collectMigrationFiles() ([]Migration, error) {
	dirs := []string{
		filepath.Join("migrations", "schema"),
		filepath.Join("migrations", "seeds"),
	}

	var migrations []Migration
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}

		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
				continue
			}

			// Expected format: 001_name_of_migration.sql
			parts := strings.SplitN(strings.TrimSuffix(e.Name(), ".sql"), "_", 2)
			if len(parts) != 2 {
				log.Printf("Skipping %s — filename must follow pattern NNN_name.sql", e.Name())
				continue
			}

			migrations = append(migrations, Migration{
				Version:  parts[0],
				Name:     parts[1],
				Filename: filepath.Join(dir, e.Name()),
			})
		}
	}

	// Sort by version ascending
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// appliedMigrations returns a map of version -> Migration for already-applied migrations
func appliedMigrations(db *sql.DB) (map[string]Migration, error) {
	rows, err := db.Query(`
		SELECT version, name, filename, checksum, applied_at
		FROM navisha_english_migration ORDER BY version`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]Migration)
	for rows.Next() {
		var m Migration
		if err := rows.Scan(&m.Version, &m.Name, &m.Filename, &m.Checksum, &m.AppliedAt); err != nil {
			continue
		}
		result[m.Version] = m
	}
	return result, nil
}
