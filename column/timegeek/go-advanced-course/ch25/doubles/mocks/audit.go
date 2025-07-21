package doubles

import "fmt"

type AuditEvent struct {
	Action  string
	UserID  string
	Details map[string]interface{}
}

type AuditService interface {
	RecordEvent(event AuditEvent) error
}

type ImportantOperationService struct {
	audit AuditService
}

func NewImportantOperationService(audit AuditService) *ImportantOperationService {
	return &ImportantOperationService{audit: audit}
}

func (s *ImportantOperationService) PerformImportantAction(userID string, data string) error {
	fmt.Printf("[ImportantOperationService] Performing action for user %s with data: %s\n", userID, data)
	// ... perform the actual important operation ...

	event := AuditEvent{
		Action:  "IMPORTANT_ACTION_PERFORMED",
		UserID:  userID,
		Details: map[string]interface{}{"input_data": data, "status": "success"},
	}
	// We expect this audit event to be recorded.
	if err := s.audit.RecordEvent(event); err != nil {
		// Log failure to audit but operation itself might still be considered successful
		fmt.Printf("[ImportantOperationService] WARN: Failed to record audit event: %v\n", err)
		// Depending on requirements, this might or might not be a fatal error for the operation.
		// For this demo, we assume it's not fatal for PerformImportantAction itself.
	}
	return nil
}
