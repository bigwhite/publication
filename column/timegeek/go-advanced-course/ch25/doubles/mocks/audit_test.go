package doubles

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert" // For assertions
	"github.com/stretchr/testify/mock"   // For mocking
)

// --- Mock AuditService Implementation ---
type MockAuditService struct {
	mock.Mock // Embed mock.Mock
}

// RecordEvent is the mock implementation of the AuditService interface method.
func (m *MockAuditService) RecordEvent(event AuditEvent) error {
	// This method records that a call was made and with what arguments.
	// It will also return an error if one was specified with .Return().
	args := m.Called(event)
	return args.Error(0) // Return the first (and only, in this case) error argument
}

func TestImportantOperationService_WithMockAudit(t *testing.T) {
	t.Run("SuccessfulOperationAndAudit", func(t *testing.T) {
		// 1. Create an instance of the mock object.
		mockAudit := new(MockAuditService)

		// 2. Setup expectations on the mock.
		// We expect RecordEvent to be called once with a specific AuditEvent structure.
		// We use mock.MatchedBy to perform a custom match on the AuditEvent argument.
		expectedUserID := "user789"
		expectedData := "sensitive_data_processed"

		mockAudit.On("RecordEvent", mock.MatchedBy(func(event AuditEvent) bool {
			return event.Action == "IMPORTANT_ACTION_PERFORMED" &&
				event.UserID == expectedUserID &&
				event.Details["input_data"] == expectedData &&
				event.Details["status"] == "success"
		})).Return(nil).Once() // Expect it to be called once and return no error.

		// 3. Create the service instance, injecting the mock.
		service := NewImportantOperationService(mockAudit)

		// 4. Call the method on the service that should trigger the mock's method.
		err := service.PerformImportantAction(expectedUserID, expectedData)
		assert.NoError(t, err, "PerformImportantAction should not return an error")

		// 5. Assert that all expectations on the mock were met.
		mockAudit.AssertExpectations(t)
	})

	t.Run("OperationSucceedsButAuditFails", func(t *testing.T) {
		mockAudit := new(MockAuditService)
		simulatedAuditError := fmt.Errorf("audit system unavailable")

		// Expect RecordEvent to be called, but this time it will return an error.
		mockAudit.On("RecordEvent", mock.AnythingOfType("AuditEvent")).Return(simulatedAuditError).Once()

		service := NewImportantOperationService(mockAudit)
		err := service.PerformImportantAction("userFail", "dataFail")
		// In our SUT, PerformImportantAction logs the audit error but doesn't propagate it.
		assert.NoError(t, err, "PerformImportantAction should still succeed even if auditing fails (per demo SUT logic)")

		mockAudit.AssertExpectations(t)
	})
}
