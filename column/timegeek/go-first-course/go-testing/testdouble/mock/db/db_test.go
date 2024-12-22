package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Define a mock struct that implements the `Database` interface
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Save(data string) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockDatabase) Get(id int) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

func TestSaveData(t *testing.T) {
	// Create a new mock database
	db := new(MockDatabase)

	// Expect the `Save` method to be called with "test data"
	db.On("Save", "test data").Return(nil)

	// Call the code that uses the database
	err := saveData(db, "test data")

	// Assert that the `Save` method was called with the correct argument
	db.AssertCalled(t, "Save", "test data")

	// Assert that no errors were returned
	assert.NoError(t, err)
}

func TestGetData(t *testing.T) {
	// Create a new mock database
	db := new(MockDatabase)

	// Expect the `Get` method to be called with ID 123 and return "test data"
	db.On("Get", 123).Return("test data", nil)

	// Call the code that uses the database
	data, err := getData(db, 123)

	// Assert that the `Get` method was called with the correct argument
	db.AssertCalled(t, "Get", 123)

	// Assert that the correct data was returned
	assert.Equal(t, "test data", data)

	// Assert that no errors were returned
	assert.NoError(t, err)
}
