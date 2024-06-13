package repository_test

import (
	"database/sql"
	"errors"
	"learning_lantern/helper"
	"learning_lantern/models"
	"learning_lantern/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		t.Fatal(err)
	}
	return sqldb, gormdb, mock
}

func TestLogin_shouldNotFound(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	repo := &repository.Repo{DB: db}
	users := sqlmock.NewRows([]string{"user_id", "user_name", "email", "password", "role", "deposit", "last_login_date", "jwt_token"})

	expectedSQL := `SELECT \* FROM "users" WHERE email= \$1 ORDER BY "users"\."user_id" LIMIT \$2`
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)
	u := models.User{Email: "admins", Password: "admin"}
	_, res := repo.Login(u)
	assert.True(t, errors.Is(res, helper.ErrCredential))
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestLogin_shouldFound(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	// Clear any existing expectations
	defer mock.ExpectationsWereMet()

	repo := &repository.Repo{DB: db}

	// Define the rows that should be returned by the mock query
	users := sqlmock.NewRows([]string{"user_id", "user_name", "email", "password", "role", "deposit", "last_login_date", "jwt_token"}).
		AddRow(12, "", "yoland", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InlvbGFuZCIsImV4cCI6MTcxODQ2Mzk3OSwiaWQiOjEyLCJyb2xlIjoidXNlciJ9.nY6_jG_hh_37jGfcfG8pbyckPwNmhinYqLL7rGabyYY", "user", 0.0, time.Now(), "")

	// Define the expected SQL query with arguments
	expectedSQL := `SELECT \* FROM "users" WHERE email= \$1 ORDER BY "users"\."user_id" LIMIT \$2`
	mock.ExpectQuery(expectedSQL).WithArgs("yoland", 1).WillReturnRows(users)

	// Create a user instance for login
	u := models.User{Email: "yoland", Password: "admin"}

	// Perform the login operation
	_, res := repo.Login(u)

	// Check assertions
	assert.True(t, errors.Is(res, nil), "expected no error")

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRegister_shouldFound(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	// Clear any existing expectations
	defer mock.ExpectationsWereMet()

	repo := &repository.Repo{DB: db}

	// Define the rows that should be returned by the mock query
	users := sqlmock.NewRows([]string{"user_id", "user_name", "email", "password", "role", "deposit", "last_login_date", "jwt_token"}).
		AddRow(12, "", "yoland", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InlvbGFuZCIsImV4cCI6MTcxODQ2Mzk3OSwiaWQiOjEyLCJyb2xlIjoidXNlciJ9.nY6_jG_hh_37jGfcfG8pbyckPwNmhinYqLL7rGabyYY", "user", 0.0, time.Now(), "")

	// Define the expected SQL query with arguments
	expectedSQL := `SELECT \* FROM "users" WHERE email= \$1 ORDER BY "users"\."user_id" LIMIT \$2`
	mock.ExpectQuery(expectedSQL).WithArgs("yoland", 1).WillReturnRows(users)

	// Create a user instance for login
	u := models.User{Email: "yoland", Password: "admin"}

	// Perform the login operation
	_, res := repo.Register(u)

	// Check assertions
	assert.True(t, errors.Is(res, helper.ErrUserExists))

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
