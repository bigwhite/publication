package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/go-github/v56/github"
)

// mockIssueJSON 模拟GitHub Issue API响应
const mockIssueJSON = `{
	"id": 123456789,
	"number": 123,
	"title": "Test Issue Title",
	"body": "This is a test issue body.\n\n## Content\n\nSome markdown content here.",
	"state": "open",
	"user": {
		"login": "testuser",
		"id": 987654321,
		"avatar_url": "https://avatars.githubusercontent.com/u/987654321?v=4",
		"html_url": "https://github.com/testuser",
		"type": "User"
	},
	"labels": [
		{
			"id": 111111111,
			"name": "bug",
			"color": "d73a4a",
			"description": "Something isn't working"
		},
		{
			"id": 222222222,
			"name": "enhancement",
			"color": "a2eeef",
			"description": "New feature or request"
		}
	],
	"assignees": [
		{
			"login": "assignee1",
			"id": 555555555,
			"avatar_url": "https://avatars.githubusercontent.com/u/555555555?v=4",
			"html_url": "https://github.com/assignee1",
			"type": "User"
		}
	],
	"milestone": {
		"id": 333333333,
		"number": 5,
		"title": "v1.0.0",
		"description": "First major release",
		"state": "open",
		"created_at": "2023-01-01T00:00:00Z",
		"updated_at": "2023-01-15T00:00:00Z",
		"due_on": "2023-03-01T00:00:00Z"
	},
	"comments": 5,
	"created_at": "2023-01-10T10:00:00Z",
	"updated_at": "2023-01-20T15:30:00Z",
	"closed_at": null,
	"url": "https://api.github.com/repos/testowner/testrepo/issues/123",
	"html_url": "https://github.com/testowner/testrepo/issues/123"
}`

// mockCommentsJSON 模拟GitHub Issue Comments API响应
const mockCommentsJSON = `[
	{
		"id": 111111111,
		"body": "First comment on this issue.",
		"user": {
			"login": "commenter1",
			"id": 111111111,
			"avatar_url": "https://avatars.githubusercontent.com/u/111111111?v=4",
			"html_url": "https://github.com/commenter1",
			"type": "User"
		},
		"created_at": "2023-01-10T11:00:00Z",
		"updated_at": "2023-01-10T11:00:00Z",
		"url": "https://api.github.com/repos/testowner/testrepo/issues/comments/111111111",
		"html_url": "https://github.com/testowner/testrepo/issues/123#issuecomment-111111111"
	},
	{
		"id": 222222222,
		"body": "Second comment with **markdown** formatting.",
		"user": {
			"login": "commenter2",
			"id": 222222222,
			"avatar_url": "https://avatars.githubusercontent.com/u/222222222?v=4",
			"html_url": "https://github.com/commenter2",
			"type": "User"
		},
		"created_at": "2023-01-11T14:30:00Z",
		"updated_at": "2023-01-11T14:30:00Z",
		"url": "https://api.github.com/repos/testowner/testrepo/issues/comments/222222222",
		"html_url": "https://github.com/testowner/testrepo/issues/123#issuecomment-222222222"
	}
]`

// createMockServer 创建模拟GitHub API服务器
func createMockServer() *httptest.Server {
	mux := http.NewServeMux()

	// Mock Issue API endpoint
	mux.HandleFunc("/repos/testowner/testrepo/issues/123", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockIssueJSON))
	})

	// Mock Comments API endpoint
	mux.HandleFunc("/repos/testowner/testrepo/issues/123/comments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockCommentsJSON))
	})

	// Mock error endpoint
	mux.HandleFunc("/repos/errorowner/errorrepo/issues/999", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{
			"message": "Not Found",
			"documentation_url": "https://docs.github.com/rest/reference/issues#get-an-issue"
		}`))
	})

	return httptest.NewServer(mux)
}

func TestGetIssue(t *testing.T) {
	// 创建模拟服务器
	server := createMockServer()
	defer server.Close()

	// 创建测试用例
	tests := []struct {
		name          string
		owner         string
		repo          string
		issueNumber   int
		wantErr       bool
		errorContains string
		setupClient   func() *GitHubClient
	}{
		{
			name:        "Successful Issue Retrieval",
			owner:       "testowner",
			repo:        "testrepo",
			issueNumber: 123,
			wantErr:     false,
			setupClient: func() *GitHubClient {
				// 创建一个指向Mock Server的HTTP客户端
				httpClient := server.Client()
				// 创建GitHub客户端，并设置基础URL指向Mock Server
				githubClient := NewClientWithHTTPClient(httpClient, "test-token")
				// 覆盖基础URL
				githubClient.Client.BaseURL, _ = url.Parse(server.URL + "/")
				return githubClient
			},
		},
		{
			name:        "Issue Not Found",
			owner:       "errorowner",
			repo:        "errorrepo",
			issueNumber: 999,
			wantErr:     true,
			setupClient: func() *GitHubClient {
				githubClient := NewClientWithHTTPClient(server.Client(), "test-token")
				githubClient.Client.BaseURL, _ = url.Parse(server.URL + "/")
				return githubClient
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := tt.setupClient()
			ctx := context.Background()

			// 调用 GetIssue 方法
			got, err := client.GetIssue(ctx, tt.owner, tt.repo, tt.issueNumber)

			// 检查错误情况
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetIssue() expected error, but got nil")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("GetIssue() error = %v, expected to contain '%s'", err.Error(), tt.errorContains)
					return
				}
				// 如果期望错误，不再检查返回值
				return
			}

			// 不期望错误但发生了错误
			if err != nil {
				t.Errorf("GetIssue() unexpected error = %v", err)
				return
			}

			// 验证返回的Issue数据
			if got == nil {
				t.Error("GetIssue() returned nil, expected Issue")
				return
			}

			// 验证关键字段
			if got.Number != tt.issueNumber {
				t.Errorf("GetIssue().Number = %v, want %v", got.Number, tt.issueNumber)
			}
			if got.Title != "Test Issue Title" {
				t.Errorf("GetIssue().Title = %v, want %v", got.Title, "Test Issue Title")
			}
			if got.State != "open" {
				t.Errorf("GetIssue().State = %v, want %v", got.State, "open")
			}
			if got.User.Login != "testuser" {
				t.Errorf("GetIssue().User.Login = %v, want %v", got.User.Login, "testuser")
			}

			// 验证Labels
			if len(got.Labels) != 2 {
				t.Errorf("GetIssue().Labels length = %v, want %v", len(got.Labels), 2)
			} else {
				if got.Labels[0].Name != "bug" {
					t.Errorf("GetIssue().Labels[0].Name = %v, want %v", got.Labels[0].Name, "bug")
				}
				if got.Labels[1].Name != "enhancement" {
					t.Errorf("GetIssue().Labels[1].Name = %v, want %v", got.Labels[1].Name, "enhancement")
				}
			}

			// 验证Assignees
			if len(got.Assignees) != 1 {
				t.Errorf("GetIssue().Assignees length = %v, want %v", len(got.Assignees), 1)
			} else {
				if got.Assignees[0].Login != "assignee1" {
					t.Errorf("GetIssue().Assignees[0].Login = %v, want %v", got.Assignees[0].Login, "assignee1")
				}
			}

			// 验证Milestone
			if got.Milestone == nil {
				t.Error("GetIssue().Milestone = nil, expected Milestone")
			} else {
				if got.Milestone.Title != "v1.0.0" {
					t.Errorf("GetIssue().Milestone.Title = %v, want %v", got.Milestone.Title, "v1.0.0")
				}
				if got.Milestone.Number != 5 {
					t.Errorf("GetIssue().Milestone.Number = %v, want %v", got.Milestone.Number, 5)
				}
			}
		})
	}
}

func TestGetIssueComments(t *testing.T) {
	// 创建模拟服务器
	server := createMockServer()
	defer server.Close()

	// 创建测试用例
	tests := []struct {
		name        string
		owner       string
		repo        string
		issueNumber int
		wantErr     bool
		setupClient func() *GitHubClient
	}{
		{
			name:        "Successful Comments Retrieval",
			owner:       "testowner",
			repo:        "testrepo",
			issueNumber: 123,
			wantErr:     false,
			setupClient: func() *GitHubClient {
				githubClient := NewClientWithHTTPClient(server.Client(), "test-token")
				githubClient.Client.BaseURL, _ = url.Parse(server.URL + "/")
				return githubClient
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := tt.setupClient()
			ctx := context.Background()

			// 调用 GetIssueComments 方法
			got, err := client.GetIssueComments(ctx, tt.owner, tt.repo, tt.issueNumber)

			// 检查错误情况
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetIssueComments() expected error, but got nil")
					return
				}
				// 如果期望错误，不再检查返回值
				return
			}

			// 不期望错误但发生了错误
			if err != nil {
				t.Errorf("GetIssueComments() unexpected error = %v", err)
				return
			}

			// 验证返回的Comments数据
			if got == nil {
				t.Error("GetIssueComments() returned nil, expected Comments")
				return
			}

			// 验证评论数量
			if len(got) != 2 {
				t.Errorf("GetIssueComments() length = %v, want %v", len(got), 2)
				return
			}

			// 验证第一条评论
			if got[0].ID != 111111111 {
				t.Errorf("GetIssueComments()[0].ID = %v, want %v", got[0].ID, 111111111)
			}
			if got[0].User.Login != "commenter1" {
				t.Errorf("GetIssueComments()[0].User.Login = %v, want %v", got[0].User.Login, "commenter1")
			}

			// 验证第二条评论
			if got[1].ID != 222222222 {
				t.Errorf("GetIssueComments()[1].ID = %v, want %v", got[1].ID, 222222222)
			}
			if got[1].User.Login != "commenter2" {
				t.Errorf("GetIssueComments()[1].User.Login = %v, want %v", got[1].User.Login, "commenter2")
			}
		})
	}
}

func TestConvertGitHubIssue(t *testing.T) {
	// 测试转换函数的基本功能
	tests := []struct {
		name     string
		input    *github.Issue
		expected *Issue
	}{
		{
			name:     "Nil Input",
			input:    nil,
			expected: nil,
		},
		{
			name: "Valid Issue",
			input: &github.Issue{
				ID:    github.Int64(123456789),
				Number: github.Int(123),
				Title: github.String("Test Issue Title"),
				Body:  github.String("Test body"),
				State: github.String("open"),
				User: &github.User{
					Login:     github.String("testuser"),
					ID:        github.Int64(987654321),
					AvatarURL: github.String("https://avatars.githubusercontent.com/u/987654321?v=4"),
					HTMLURL:   github.String("https://github.com/testuser"),
					Type:      github.String("User"),
				},
				CreatedAt: &github.Timestamp{Time: time.Date(2023, 1, 10, 10, 0, 0, 0, time.UTC)},
				UpdatedAt: &github.Timestamp{Time: time.Date(2023, 1, 20, 15, 30, 0, 0, time.UTC)},
			},
			expected: &Issue{
				Number:    123,
				Title:     "Test Issue Title",
				Body:      "Test body",
				State:     "open",
				User: User{
					Login:     "testuser",
					ID:        987654321,
					AvatarURL: "https://avatars.githubusercontent.com/u/987654321?v=4",
					HTMLURL:   "https://github.com/testuser",
					Type:      "User",
				},
				CreatedAt: time.Date(2023, 1, 10, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 1, 20, 15, 30, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertGitHubIssue(tt.input)

			if tt.expected == nil {
				if got != nil {
					t.Errorf("convertGitHubIssue() expected nil, got %v", got)
				}
				return
			}

			if got == nil {
				t.Errorf("convertGitHubIssue() got nil, expected %v", tt.expected)
				return
			}

			if got.Number != tt.expected.Number {
				t.Errorf("convertGitHubIssue().Number = %v, want %v", got.Number, tt.expected.Number)
			}
			if got.Title != tt.expected.Title {
				t.Errorf("convertGitHubIssue().Title = %v, want %v", got.Title, tt.expected.Title)
			}
			if got.User.Login != tt.expected.User.Login {
				t.Errorf("convertGitHubIssue().User.Login = %v, want %v", got.User.Login, tt.expected.User.Login)
			}
		})
	}
}

func TestConvertGitHubComment(t *testing.T) {
	// 测试转换函数的基本功能
	tests := []struct {
		name     string
		input    *github.IssueComment
		expected *Comment
	}{
		{
			name:     "Nil Input",
			input:    nil,
			expected: nil,
		},
		{
			name: "Valid Comment",
			input: &github.IssueComment{
				ID:   github.Int64(111111111),
				Body: github.String("First comment on this issue."),
				User: &github.User{
					Login:     github.String("commenter1"),
					ID:        github.Int64(111111111),
					AvatarURL: github.String("https://avatars.githubusercontent.com/u/111111111?v=4"),
					HTMLURL:   github.String("https://github.com/commenter1"),
					Type:      github.String("User"),
				},
				CreatedAt: &github.Timestamp{Time: time.Date(2023, 1, 10, 11, 0, 0, 0, time.UTC)},
				UpdatedAt: &github.Timestamp{Time: time.Date(2023, 1, 10, 11, 0, 0, 0, time.UTC)},
			},
			expected: &Comment{
				ID:   111111111,
				Body: "First comment on this issue.",
				User: User{
					Login:     "commenter1",
					ID:        111111111,
					AvatarURL: "https://avatars.githubusercontent.com/u/111111111?v=4",
					HTMLURL:   "https://github.com/commenter1",
					Type:      "User",
				},
				CreatedAt: time.Date(2023, 1, 10, 11, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 1, 10, 11, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertGitHubComment(tt.input)

			if tt.expected == nil {
				if got != nil {
					t.Errorf("convertGitHubComment() expected nil, got %v", got)
				}
				return
			}

			if got == nil {
				t.Errorf("convertGitHubComment() got nil, expected %v", tt.expected)
				return
			}

			if got.ID != tt.expected.ID {
				t.Errorf("convertGitHubComment().ID = %v, want %v", got.ID, tt.expected.ID)
			}
			if got.Body != tt.expected.Body {
				t.Errorf("convertGitHubComment().Body = %v, want %v", got.Body, tt.expected.Body)
			}
			if got.User.Login != tt.expected.User.Login {
				t.Errorf("convertGitHubComment().User.Login = %v, want %v", got.User.Login, tt.expected.User.Login)
			}
		})
	}
}