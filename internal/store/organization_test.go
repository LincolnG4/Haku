package store

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestOrganizationStore_Create(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create a new OrganizationStore with the mock database
	store := &OrganizationStore{db: db}

	// Create test data
	org := &Organization{
		Name:        "Test Org",
		Description: "Test Description",
	}

	// Set up expectations
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(1, time.Now(), time.Now())

	mock.ExpectQuery("INSERT INTO organizations").
		WithArgs(org.Name, org.Description).
		WillReturnRows(rows)

	// Execute the test
	err = store.Create(context.Background(), org)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, int64(1), org.ID)
	assert.NotEmpty(t, org.CreatedAt)
	assert.NotEmpty(t, org.UpdatedAt)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestOrganizationStore_GetByID(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create a new OrganizationStore with the mock database
	store := &OrganizationStore{db: db}

	// Test cases
	tests := []struct {
		name          string
		orgID         int64
		mockSetup     func()
		expectedOrg   *Organization
		expectedError error
	}{
		{
			name:  "Success",
			orgID: 1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
					AddRow(1, "Test Org", "Test Description", time.Now(), time.Now())

				mock.ExpectQuery("SELECT id, name, description, created_at, updated_at").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedOrg: &Organization{
				ID:          1,
				Name:        "Test Org",
				Description: "Test Description",
			},
			expectedError: nil,
		},
		{
			name:  "Not Found",
			orgID: 999,
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, name, description, created_at, updated_at").
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			expectedOrg:   nil,
			expectedError: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations
			tt.mockSetup()

			// Execute the test
			org, err := store.GetByID(context.Background(), tt.orgID)

			// Assertions
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, org)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, org)
				assert.Equal(t, tt.expectedOrg.ID, org.ID)
				assert.Equal(t, tt.expectedOrg.Name, org.Name)
				assert.Equal(t, tt.expectedOrg.Description, org.Description)
				assert.NotEmpty(t, org.CreatedAt)
				assert.NotEmpty(t, org.UpdatedAt)
			}

			// Ensure all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}
