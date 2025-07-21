package doubles

import (
	"fmt"
	"testing"
)

// --- Fake Notifier Implementation ---
type FakeUserNotifier struct {
	NotificationsSent []User // Stores users for whom notifications were "sent"
	ShouldFail        bool   // Control whether NotifyUserCreated should return an error
	FailureError      error  // The error to return if ShouldFail is true
}

func NewFakeUserNotifier() *FakeUserNotifier {
	return &FakeUserNotifier{NotificationsSent: []User{}}
}

func (f *FakeUserNotifier) NotifyUserCreated(user User) error {
	if f.ShouldFail {
		fmt.Printf("[FakeNotifier] Simulating failure for user: %s\n", user.ID)
		return f.FailureError
	}
	fmt.Printf("[FakeNotifier] 'Sent' notification for user: %+v\n", user)
	f.NotificationsSent = append(f.NotificationsSent, user)
	return nil
}

// Helper to check if a notification was sent for a specific user ID
func (f *FakeUserNotifier) WasNotificationSentFor(userID string) bool {
	for _, u := range f.NotificationsSent {
		if u.ID == userID {
			return true
		}
	}
	return false
}

func TestUserCreationService_WithFakeNotifier(t *testing.T) {
	t.Run("SuccessfulNotification", func(t *testing.T) {
		fakeNotifier := NewFakeUserNotifier()
		service := NewUserCreationService(fakeNotifier)

		user, err := service.CreateUserAndNotify("u123", "alice@example.com", "Alice")
		if err != nil {
			t.Fatalf("CreateUserAndNotify failed: %v", err)
		}

		if !fakeNotifier.WasNotificationSentFor(user.ID) {
			t.Errorf("Expected notification to be sent for user %s, but it wasn't.", user.ID)
		}
		if len(fakeNotifier.NotificationsSent) != 1 {
			t.Errorf("Expected 1 notification to be sent, got %d", len(fakeNotifier.NotificationsSent))
		}
	})

	t.Run("FailedNotification", func(t *testing.T) {
		fakeNotifier := NewFakeUserNotifier()
		fakeNotifier.ShouldFail = true
		fakeNotifier.FailureError = fmt.Errorf("simulated network error")

		service := NewUserCreationService(fakeNotifier)

		user, err := service.CreateUserAndNotify("u456", "bob@example.com", "Bob")
		if err != nil {
			// Assuming CreateUserAndNotify itself doesn't fail on notification error for this demo
			t.Fatalf("CreateUserAndNotify unexpectedly failed itself: %v", err)
		}

		// Check that notification was attempted but no successful notification recorded
		if fakeNotifier.WasNotificationSentFor(user.ID) {
			t.Errorf("Notification for user %s should have failed, but fake shows it as sent.", user.ID)
		}
		if len(fakeNotifier.NotificationsSent) != 0 {
			t.Errorf("Expected 0 successful notifications, got %d", len(fakeNotifier.NotificationsSent))
		}
		// In a real test, you might also check logs or other side effects of a failed notification attempt.
	})
}
