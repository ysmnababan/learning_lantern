package repository_test

import (
	"database/sql"
	"errors"
	"learning_lantern/helper"
	"learning_lantern/models"
	"learning_lantern/repository"
	"log"
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

func TestGetInfo(t *testing.T) {
	// Create a mock database connection
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	// Create a repository instance using the mock DB
	repo := &repository.Repo{DB: db}

	// Define expected data and mock responses
	userRows := sqlmock.NewRows([]string{"user_id", "username", "email", "deposit"}).
		AddRow(1, "john_doe", "john.doe@example.com", 100.0)

	userDetailRows := sqlmock.NewRows([]string{"user_id", "fname", "lname", "address", "age", "phone_number"}).
		AddRow(1, "John", "Doe", "123 Main St", 30, "1234567890")

	// Define the expected SQL queries and responses
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."user_id" = \$1 ORDER BY "users"."user_id" LIMIT \$2`).WithArgs(1, 1).WillReturnRows(userRows)
	mock.ExpectQuery(`SELECT \* FROM "user_details" WHERE user_id=\$1`).WithArgs(1, 1).WillReturnRows(userDetailRows)

	// Call the method under test
	result, err := repo.GetInfo(1)
	log.Println(err)
	// Assert the result
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, uint(1), result.UserID)
	assert.Equal(t, "john_doe", result.Username)
	assert.Equal(t, "john.doe@example.com", result.Email)
	assert.Equal(t, 100.0, result.Deposit)
	assert.Equal(t, "John", result.Fname)
	assert.Equal(t, "Doe", result.Lname)
	assert.Equal(t, "123 Main St", result.Address)
	assert.Equal(t, 30, result.Age)
	assert.Equal(t, "1234567890", result.PhoneNumber)

	// Ensure all expectations were met;
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestTopUp(t *testing.T) {
	// Create a mock database connection
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	// Create a repository instance using the mock DB
	repo := &repository.Repo{DB: db}

	// Define mock data for user

	// Expectations for querying user by user_id
	users := sqlmock.NewRows([]string{"user_id", "user_name", "email", "password", "role", "deposit", "last_login_date", "jwt_token"}).
		AddRow(12, "", "yoland", "eabyYY", "user", 100, time.Now(), "")

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."user_id" = \$1 ORDER BY "users"."user_id" LIMIT \$2`).WithArgs(12, 1).WillReturnRows(users)

	// Expectations for updating user deposit
	mock.ExpectExec(`UPDATE "users" SET "deposit"=\$1 WHERE "user_id"=\$2`).
		WithArgs(150.0, 12).
		WillReturnResult(sqlmock.NewResult(12, 1)) // Mock update success

	// Call the method under test
	updatedDeposit, err := repo.TopUp(12, 50.0)

	// Assert the result
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, 150.0, updatedDeposit, "expected updated deposit to match")

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
