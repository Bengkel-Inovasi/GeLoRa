package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	_ = godotenv.Load(".env")

	dbPath := os.Getenv("SQLITE_DATABASE")
	if dbPath == "" {
		dbPath = "./gelora.db"
	}

	direction := "up"
	if len(os.Args) > 1 {
		direction = os.Args[1]
	}

	if err := run(dbPath, direction); err != nil {
		log.Fatal(err)
	}
}

func run(dbPath, direction string) error {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		f, err := os.Create(dbPath)
		if err != nil {
			return fmt.Errorf("create db file: %w", err)
		}
		f.Close()
		fmt.Printf("Created database file: %s\n", dbPath)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY)`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	files, err := filepath.Glob(fmt.Sprintf("database/migrations/*.%s.sql", direction))
	if err != nil {
		return err
	}
	sort.Strings(files)

	if direction == "down" {
		for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
			files[i], files[j] = files[j], files[i]
		}
	}

	for _, f := range files {
		base := filepath.Base(f)
		version := strings.TrimSuffix(strings.TrimSuffix(base, ".down.sql"), ".up.sql")

		if direction == "up" {
			var exists int
			_ = db.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = ?`, version).Scan(&exists)
			if exists > 0 {
				fmt.Printf("  skip (already applied): %s\n", base)
				continue
			}
		}

		content, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("read %s: %w", f, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("exec %s: %w", f, err)
		}

		if direction == "up" {
			if _, err := db.Exec(`INSERT OR IGNORE INTO schema_migrations (version) VALUES (?)`, version); err != nil {
				return fmt.Errorf("record migration %s: %w", version, err)
			}
			fmt.Printf("  applied: %s\n", base)
		} else {
			if _, err := db.Exec(`DELETE FROM schema_migrations WHERE version = ?`, version); err != nil {
				return fmt.Errorf("remove migration %s: %w", version, err)
			}
			fmt.Printf("  rolled back: %s\n", base)
		}
	}

	fmt.Println("Migration complete.")
	return nil
}
