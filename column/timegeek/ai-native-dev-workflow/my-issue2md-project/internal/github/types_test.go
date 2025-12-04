package github

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  *GitHubClient
	}{
		{
			name:  "valid token",
			token: "test-token",
			want:  &GitHubClient{},
		},
		{
			name:  "empty token",
			token: "",
			want:  &GitHubClient{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewClient(tt.token)
			if got == nil {
				t.Errorf("NewClient() returned nil")
			}
			if got.Client == nil {
				t.Errorf("NewClient().Client is nil")
			}
		})
	}
}

func TestIssue(t *testing.T) {
	now := time.Now()
	issue := &Issue{
		Number:    123,
		Title:     "Test Issue",
		Body:      "This is a test issue",
		State:     "open",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if issue.Number != 123 {
		t.Errorf("Expected issue number 123, got %d", issue.Number)
	}

	if issue.Title != "Test Issue" {
		t.Errorf("Expected title 'Test Issue', got %s", issue.Title)
	}
}

func TestComment(t *testing.T) {
	now := time.Now()
	comment := &Comment{
		ID:        456,
		Body:      "This is a test comment",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if comment.ID != 456 {
		t.Errorf("Expected comment ID 456, got %d", comment.ID)
	}

	if comment.Body != "This is a test comment" {
		t.Errorf("Expected body 'This is a test comment', got %s", comment.Body)
	}
}

func TestUser(t *testing.T) {
	user := &User{
		Login:     "testuser",
		ID:        789,
		AvatarURL: "https://example.com/avatar.jpg",
		Type:      "User",
	}

	if user.Login != "testuser" {
		t.Errorf("Expected login 'testuser', got %s", user.Login)
	}

	if user.ID != 789 {
		t.Errorf("Expected ID 789, got %d", user.ID)
	}
}