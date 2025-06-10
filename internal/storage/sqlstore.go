package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type sqlStore struct {
	db     *sql.DB
	driver string
}

func NewSQLite(path string) (Provider, error) {
	dsn := fmt.Sprintf("file:%s?_foreign_keys=on&_busy_timeout=5000", path)
	return newSQLStore("sqlite3", dsn)
}

func NewSQL(driver, dsn string) (Provider, error) {
	return newSQLStore(driver, dsn)
}

func newSQLStore(driver, dsn string) (Provider, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	store := &sqlStore{
		db:     db,
		driver: driver,
	}
	if err := store.initializeSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

func (s *sqlStore) initializeSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		client_id TEXT NOT NULL,
		agent_id TEXT NOT NULL,
		key BLOB NOT NULL,
		iv BLOB NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME
	);

	CREATE UNIQUE INDEX IF NOT EXISTS idx_client_agent ON sessions (client_id, agent_id);
	`

	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	return nil
}

func (s *sqlStore) SaveSession(ctx context.Context, session Session) error {
	if s.db == nil {
		return errors.New("database connection is not initialized")
	}

	query := `
	INSERT INTO sessions (client_id, agent_id, key, iv, created_at, expires_at)
	VALUES (?, ?, ?, ?, ?, ?)
	ON CONFLICT(client_id, agent_id) DO UPDATE SET
		key = excluded.key,
		iv = excluded.iv,
		created_at = excluded.created_at,
		expires_at = excluded.expires_at;
	`

	_, err := s.db.ExecContext(ctx, query,
		session.ClientID, session.AgentID, session.Key, session.IV,
		session.CreatedAt, session.ExpiresAt.Time)

	return err
}

func (s *sqlStore) GetSession(ctx context.Context, clientID, agentID string) (Session, error) {
	if s.db == nil {
		return Session{}, errors.New("database connection is not initialized")
	}

	query := `
	SELECT client_id, agent_id, key, iv, created_at, expires_at
	FROM sessions
	WHERE client_id = ? AND agent_id = ?;
	`

	row := s.db.QueryRowContext(ctx, query, clientID, agentID)

	var session Session
	var expiresAt sql.NullTime
	err := row.Scan(&session.ClientID, &session.AgentID, &session.Key, &session.IV,
		&session.CreatedAt, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Session{}, fmt.Errorf("session not found for client %s and agent %s", clientID, agentID)
		}
		return Session{}, fmt.Errorf("failed to get session: %w", err)
	}

	session.ExpiresAt = expiresAt
	return session, nil
}

func (s *sqlStore) DeleteSession(ctx context.Context, clientID, agentID string) error {
	if s.db == nil {
		return errors.New("database connection is not initialized")
	}

	query := `
	DELETE FROM sessions
	WHERE client_id = ? AND agent_id = ?;
	`

	result, err := s.db.ExecContext(ctx, query, clientID, agentID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no session found for client %s and agent %s", clientID, agentID)
	}

	return nil
}

func (s *sqlStore) Close() error {
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	s.db = nil
	return nil
}
