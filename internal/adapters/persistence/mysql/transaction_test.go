package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MatheusHenrique129/bemax-api/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestTransactionMySQL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	loggerMock := mocks.NewMockLogger()
	baseRepository := &BaseRepository{dbClient: db, logger: loggerMock}

	testCases := []struct {
		name                string
		expectError         bool
		expectRollbackError bool
		errorPanic          bool
		errorType           error
		sqlResult           driver.Result
	}{
		{name: "success", expectError: false, errorType: nil, sqlResult: sqlmock.NewResult(1, 1)},
		{name: "error - rollback", expectError: true, errorType: ErrDBRollback},
		{name: "error - rollback failed", expectError: true, expectRollbackError: true, errorType: ErrDBRollback},
		{name: "error - panic recover with error", expectError: true, errorPanic: true, errorType: fmt.Errorf("panic attack")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectBegin()

			if tc.expectError {
				switch tc.expectRollbackError {
				case true:
					mock.ExpectRollback().WillReturnError(tc.errorType)
				default:
					mock.ExpectRollback()
				}
			} else {
				mock.ExpectExec("SELECT field FROM table").WillReturnResult(tc.sqlResult)
				mock.ExpectCommit()
			}

			err = baseRepository.WithTransaction(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
				if tc.errorPanic {
					panic(tc.errorType)
				}

				if tc.errorType != nil {
					return tc.errorType
				}

				query := "SELECT field FROM table"
				_, e := tx.ExecContext(ctx, query)
				if e != nil {
					return e
				}
				return nil
			})

			if tc.expectError {
				assert.ErrorIs(t, err, tc.errorType)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
