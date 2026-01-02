package db

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // file:// source
)

func RunMigrations2(databaseURL string) error {
	// For local dev (repo root): file://migrations
	// For Docker runtime (we’ll copy to /migrations): file:///migrations
	// We’ll try /migrations first; if that fails, fall back to ./migrations.
	//sources := []string{"file:///migrations", "file://migrations"}
	sources := []string{"file:///Users/aarkay/github/go-service-template-sql/migrations"}

	var lastErr error
	for _, src := range sources {
		m, err := migrate.New(src, databaseURL)
		if err != nil {
			fmt.Printf("err:%v\n", err)
			lastErr = err
			continue
		}

		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				_, _ = m.Close()
				return nil
			}
			_, _ = m.Close()
			return fmt.Errorf("migrations up failed: %w", err)
		}

		_, _ = m.Close()
		return nil
	}

	return fmt.Errorf("failed to init migrate with any source: %w", lastErr)
}

func RunMigrations(databaseURL string) error {
	// 1) Try /migrations (Docker convention)
	candidates := []string{"/Users/aarkay/github/go-service-template-sql/migrations"}

	// 2) Try ./migrations relative to cwd
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(wd, "migrations"))
	}

	var lastErr error
	for _, p := range candidates {
		// ensure dir exists before calling migrate
		if st, err := os.Stat(p); err != nil || !st.IsDir() {
			lastErr = fmt.Errorf("migration dir not found: %s (%v)", p, err)
			continue
		}
		// optional: sanity check it contains files
		if entries, err := os.ReadDir(p); err != nil || len(entries) == 0 {
			lastErr = fmt.Errorf("migration dir empty/unreadable: %s (%v)", p, err)
			continue
		}

		// u := url.URL{Scheme: "file", Path: p} // properly forms file:///...
		src := "file://" + filepath.ToSlash(p) // IMPORTANT
		m, err := migrate.New(src, databaseURL)
		if err != nil {
			lastErr = err
			continue
		}

		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				_, _ = m.Close()
				return nil
			}
			_, _ = m.Close()
			return fmt.Errorf("migrations up failed: %w", err)
		}

		_, _ = m.Close()
		return nil
	}

	return fmt.Errorf("failed to init migrate: %w", lastErr)
}
