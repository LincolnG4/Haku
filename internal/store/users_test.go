package store

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserStore_Create(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create a new UserStore with the mock database
	store := &UsersStore{db: db}

	userPassword := password{}
	userPassword.Set("MyPassword")

	// Create test data
	user := &User{
		Username: "Haku",
		Password: userPassword,
		Email:    "haku@ghibli.com",
	}

	// Set up expectations
	rows := sqlmock.NewRows([]string{"id", "created_at"}).
		AddRow(1, time.Now())

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Username,
			user.Password.hash,
			user.Email).
		WillReturnRows(rows)

	// Execute the test
	err = store.Create(context.Background(), user)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
	assert.NotEmpty(t, user.CreatedAt)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUserStore_GetByID(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create a new UserStore with the mock database
	store := &UsersStore{db: db}

	queryTest := "SELECT id, username, password, email,created_at,updated_at"
	// Test cases
	tests := []struct {
		name          string
		userID        int64
		mockSetup     func()
		expectedUser  User
		expectedError error
	}{
		{
			name:   "Success",
			userID: 1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "password", "email", "created_at", "updated_at"}).
					AddRow(1, "Haku", "h", "haku@ghibli.com", time.Now(), time.Now())

				mock.ExpectQuery(queryTest).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedUser: User{
				ID:       1,
				Username: "Haku",
				Email:    "haku@ghibli.com",
			},
			expectedError: nil,
		},
		{
			name:   "Not Found",
			userID: 999,
			mockSetup: func() {
				mock.ExpectQuery(queryTest).
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser:  User{},
			expectedError: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations
			tt.mockSetup()

			// Execute the test
			user, err := store.GetByID(context.Background(), tt.userID)

			// Assertions
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.NotEmpty(t, user.Password.hash)
				if tt.name == "Success" {
					userPassword := password{}
					userPassword.Set("MyPassword")
					assert.Equal(t, []byte("h"), user.Password.hash)
				}
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.NotEmpty(t, user.CreatedAt)
				assert.NotEmpty(t, user.UpdatedAt)
			}

			// Ensure all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}

}

func TestUserStore_GetByEmail(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	store := &UsersStore{db: db}

	queryTest := "SELECT id, username, email, password,created_at,updated_at"
	tests := []struct {
		name          string
		email         string
		mockSetup     func()
		expectedUser  User
		expectedError error
	}{
		{
			name:  "Success",
			email: "haku@ghibli.com",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
					AddRow(1, "Haku", "haku@ghibli.com", []byte("hash"), time.Now(), time.Now())
				mock.ExpectQuery(queryTest).
					WithArgs("haku@ghibli.com").
					WillReturnRows(rows)
			},
			expectedUser: User{
				ID:       1,
				Username: "Haku",
				Email:    "haku@ghibli.com",
			},
			expectedError: nil,
		},
		{
			name:  "Not Found",
			email: "notfound@ghibli.com",
			mockSetup: func() {
				mock.ExpectQuery(queryTest).
					WithArgs("notfound@ghibli.com").
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser:  User{},
			expectedError: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			user, err := store.GetByEmail(context.Background(), tt.email)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.NotEmpty(t, user.CreatedAt)
				assert.NotEmpty(t, user.UpdatedAt)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPassword_SetAndMatches(t *testing.T) {
	var p password
	plaintext := "supersecret"
	// Test Set
	err := p.Set(plaintext)
	assert.NoError(t, err)
	assert.NotEmpty(t, p.hash)
	assert.Equal(t, plaintext, p.plaintText)

	// Test Matches (correct password)
	assert.True(t, p.Matches(plaintext))
	// Test Matches (incorrect password)
	assert.False(t, p.Matches("wrongpassword"))
}

func TestUserStore_Create_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	store := &UsersStore{db: db}
	userPassword := password{}
	userPassword.Set("MyPassword")
	user := &User{
		Username: "Haku",
		Password: userPassword,
		Email:    "haku@ghibli.com",
	}
	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Username, user.Password.hash, user.Email).
		WillReturnError(sql.ErrConnDone)

	err = store.Create(context.Background(), user)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
