package github

import (
	"context"
	"net/http"
	"time"

	"github.com/google/go-github/v56/github"
)

// Client 定义GitHub客户端接口
type Client interface {
	GetIssue(ctx context.Context, owner, repo string, issueNumber int) (*Issue, error)
	GetIssueComments(ctx context.Context, owner, repo string, issueNumber int) ([]*Comment, error)
}

// GitHubClient GitHub客户端实现
type GitHubClient struct {
	Client *github.Client
}

// NewClient 创建新的GitHub客户端
func NewClient(token string) *GitHubClient {
	client := github.NewClient(nil).WithAuthToken(token)
	return &GitHubClient{
		Client: client,
	}
}

// NewClientWithHTTPClient 使用自定义HTTP客户端创建GitHub客户端
func NewClientWithHTTPClient(httpClient *http.Client, token string) *GitHubClient {
	client := github.NewClient(httpClient).WithAuthToken(token)
	return &GitHubClient{
		Client: client,
	}
}

// Issue 表示一个GitHub Issue
type Issue struct {
	Number      int       `json:"number"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	State       string    `json:"state"`
	User        User      `json:"user"`
	Labels      []Label   `json:"labels"`
	Assignees   []User    `json:"assignees"`
	Milestone   *Milestone `json:"milestone,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`
	URL         string    `json:"url"`
	HTMLURL     string    `json:"html_url"`
}

// Comment 表示Issue评论
type Comment struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"url"`
	HTMLURL   string    `json:"html_url"`
}

// User 表示GitHub用户
type User struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
}

// Label 表示Issue标签
type Label struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description,omitempty"`
}

// Milestone 表示Issue里程碑
type Milestone struct {
	Title       string     `json:"title"`
	Number      int        `json:"number"`
	State       string     `json:"state"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`
}

// Repository 表示GitHub仓库
type Repository struct {
	Owner    string `json:"owner"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	URL      string `json:"url"`
	HTMLURL  string `json:"html_url"`
}