package database

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithTx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockDB := &DB{DB: db}
	mock.ExpectBegin()
	mock.ExpectCommit()

	err = mockDB.WithTx(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		return nil
	})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWithTx_ErrorRollback(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockDB := &DB{DB: db}
	mock.ExpectBegin()
	mock.ExpectRollback()

	err = mockDB.WithTx(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		return errors.New("transaction error")
	})
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWithTx_CommitFailureRetry(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockDB := &DB{DB: db}
	mock.ExpectBegin()
	mock.ExpectCommit().WillReturnError(errors.New("pq: deadlock detected"))
	mock.ExpectBegin()
	mock.ExpectCommit()

	err = mockDB.WithTx(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		return nil
	})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWithTx_DeadlockRetryExhausted(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockDB := &DB{DB: db}
	for i := 0; i < 3; i++ {
		mock.ExpectBegin()
		mock.ExpectCommit().WillReturnError(errors.New("pq: deadlock detected"))
	}

	err = mockDB.WithTx(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		return nil
	})
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWithTx_BeginError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockDB := &DB{DB: db}
	mock.ExpectBegin().WillReturnError(errors.New("begin error"))

	err = mockDB.WithTx(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		return nil
	})
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIsDeadlock(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"pq deadlock", errors.New("pq: deadlock detected"), true},
		{"ERROR #40P01", errors.New("ERROR #40P01 deadlock detected"), true},
		{"other error", errors.New("some other error"), false},
		{"contains deadlock", errors.New("deadlock detected in transaction"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isDeadlock(tt.err))
		})
	}
}
