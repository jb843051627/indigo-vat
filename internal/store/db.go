package store

import (
	"context"
	"database/sql"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"sync"
)

type DB struct {
	sql  *sql.DB
	mu   sync.RWMutex
	path string
}

func Open(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db := &DB{sql: conn, path: path}
	if _, err := conn.Exec("PRAGMA busy_timeout=5000; PRAGMA journal_mode=WAL;"); err != nil {
		conn.Close()
		return nil, err
	}
	if err := db.InitSchema(context.Background()); err != nil {
		conn.Close()
		return nil, err
	}
	return db, nil
}
func (d *DB) InitSchema(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, err := d.sql.ExecContext(ctx, Schema)
	return err
}
func (d *DB) Close() error { d.mu.Lock(); defer d.mu.Unlock(); return d.sql.Close() }
func (d *DB) Ping(ctx context.Context) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.sql.PingContext(ctx)
}
func (d *DB) Path() string { return d.path }
func (d *DB) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := context.Background().Err(); err != nil {
		return err
	}
	tx, err := d.sql.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
func (d *DB) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.sql.ExecContext(ctx, query, args...)
}
func (d *DB) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.sql.QueryContext(ctx, query, args...)
}
func (d *DB) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return d.sql.QueryRowContext(ctx, query, args...)
}
