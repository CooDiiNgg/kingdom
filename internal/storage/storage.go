package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Session struct {
	ID        int64
	ClientID  string
	AgentID   string
	Key       []byte
	IV        []byte
	CreatedAt time.Time
	ExpiresAt sql.NullTime
}

type Provider interface {
	SaveSession(ctx context.Context, session Session) error
	GetSession(ctx context.Context, clientID, agentID string) (Session, error)
	DeleteSession(ctx context.Context, clientID, agentID string) error
	Close() error
}

var defaultProvider Provider

func SetDefaultProvider(provider Provider) {
	defaultProvider = provider
}

func NewSession(clientID, agentID string, key, iv []byte) Session {
	return Session{
		ClientID:  clientID,
		AgentID:   agentID,
		Key:       key,
		IV:        iv,
		CreatedAt: time.Now(),
		ExpiresAt: sql.NullTime{Time: time.Now().Add(24 * time.Hour), Valid: true},
	}
}

func (s *Session) IsExpired() bool {
	if s.ExpiresAt.Valid {
		return time.Now().After(s.ExpiresAt.Time)
	}
	return false
}

func SaveSession(session Session) error {
	if defaultProvider == nil {
		return errors.New("default storage provider is not set")
	}
	ctx := context.Background()
	return defaultProvider.SaveSession(ctx, session)
}

func GetSession(clientID, agentID string) (Session, bool) {
	if defaultProvider == nil {
		return Session{}, false
	}
	ctx := context.Background()
	session, err := defaultProvider.GetSession(ctx, clientID, agentID)
	if err != nil {
		return Session{}, false
	}
	return session, true
}

func DeleteSession(clientID, agentID string) error {
	if defaultProvider == nil {
		return errors.New("default storage provider is not set")
	}
	ctx := context.Background()
	return defaultProvider.DeleteSession(ctx, clientID, agentID)
}

func Close() error {
	if defaultProvider == nil {
		return errors.New("default storage provider is not set")
	}
	return defaultProvider.Close()
}
