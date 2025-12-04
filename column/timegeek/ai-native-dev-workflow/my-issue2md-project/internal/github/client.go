package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v56/github"
)

// GetIssue 获取Issue信息
func (c *GitHubClient) GetIssue(ctx context.Context, owner, repo string, issueNumber int) (*Issue, error) {
	// 调用 GitHub API 获取 Issue
	gitHubIssue, _, err := c.Client.Issues.Get(ctx, owner, repo, issueNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue %d from %s/%s: %w", issueNumber, owner, repo, err)
	}

	// 转换为内部结构
	issue := convertGitHubIssue(gitHubIssue)
	return issue, nil
}

// GetIssueComments 获取Issue评论
func (c *GitHubClient) GetIssueComments(ctx context.Context, owner, repo string, issueNumber int) ([]*Comment, error) {
	// 调用 GitHub API 获取 Issue 评论列表
	gitHubComments, _, err := c.Client.Issues.ListComments(ctx, owner, repo, issueNumber, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments for issue %d from %s/%s: %w", issueNumber, owner, repo, err)
	}

	// 转换为内部结构
	var comments []*Comment
	if gitHubComments != nil {
		for _, gitHubComment := range gitHubComments {
			if gitHubComment != nil {
				comment := convertGitHubComment(gitHubComment)
				comments = append(comments, comment)
			}
		}
	}

	return comments, nil
}

// convertGitHubIssue 将GitHub API的Issue转换为内部Issue结构
func convertGitHubIssue(gitHubIssue *github.Issue) *Issue {
	if gitHubIssue == nil {
		return nil
	}

	var labels []Label
	if gitHubIssue.Labels != nil {
		for _, label := range gitHubIssue.Labels {
			if label != nil {
				l := Label{
					Name:        label.GetName(),
					Color:       label.GetColor(),
					Description: label.GetDescription(),
				}
				labels = append(labels, l)
			}
		}
	}

	var assignees []User
	if gitHubIssue.Assignees != nil {
		for _, assignee := range gitHubIssue.Assignees {
			if assignee != nil {
				a := User{
					Login:     assignee.GetLogin(),
					ID:        assignee.GetID(),
					AvatarURL: assignee.GetAvatarURL(),
					HTMLURL:   assignee.GetHTMLURL(),
					Type:      assignee.GetType(),
				}
				assignees = append(assignees, a)
			}
		}
	}

	var milestone *Milestone
	if gitHubIssue.Milestone != nil {
		var dueDate, closedAt *time.Time
		if gitHubIssue.Milestone.DueOn != nil {
			dueDate = &gitHubIssue.Milestone.DueOn.Time
		}
		if gitHubIssue.Milestone.ClosedAt != nil {
			closedAt = &gitHubIssue.Milestone.ClosedAt.Time
		}

		milestone = &Milestone{
			Title:       gitHubIssue.Milestone.GetTitle(),
			Number:      gitHubIssue.Milestone.GetNumber(),
			State:       gitHubIssue.Milestone.GetState(),
			Description: gitHubIssue.Milestone.GetDescription(),
			CreatedAt:   gitHubIssue.Milestone.GetCreatedAt().Time,
			UpdatedAt:   gitHubIssue.Milestone.GetUpdatedAt().Time,
			DueDate:     dueDate,
			ClosedAt:    closedAt,
		}
	}

	var user User
	if gitHubIssue.User != nil {
		user = User{
			Login:     gitHubIssue.User.GetLogin(),
			ID:        gitHubIssue.User.GetID(),
			AvatarURL: gitHubIssue.User.GetAvatarURL(),
			HTMLURL:   gitHubIssue.User.GetHTMLURL(),
			Type:      gitHubIssue.User.GetType(),
		}
	}

	var closedAt *time.Time
	if gitHubIssue.ClosedAt != nil {
		closedAt = &gitHubIssue.ClosedAt.Time
	}

	return &Issue{
		Number:    gitHubIssue.GetNumber(),
		Title:     gitHubIssue.GetTitle(),
		Body:      gitHubIssue.GetBody(),
		State:     gitHubIssue.GetState(),
		User:      user,
		Labels:    labels,
		Assignees: assignees,
		Milestone: milestone,
		CreatedAt: gitHubIssue.GetCreatedAt().Time,
		UpdatedAt: gitHubIssue.GetUpdatedAt().Time,
		ClosedAt:  closedAt,
		URL:       gitHubIssue.GetURL(),
		HTMLURL:   gitHubIssue.GetHTMLURL(),
	}
}

// convertGitHubComment 将GitHub API的Comment转换为内部Comment结构
func convertGitHubComment(gitHubComment *github.IssueComment) *Comment {
	if gitHubComment == nil {
		return nil
	}

	var user User
	if gitHubComment.User != nil {
		user = User{
			Login:     gitHubComment.User.GetLogin(),
			ID:        gitHubComment.User.GetID(),
			AvatarURL: gitHubComment.User.GetAvatarURL(),
			HTMLURL:   gitHubComment.User.GetHTMLURL(),
			Type:      gitHubComment.User.GetType(),
		}
	}

	return &Comment{
		ID:        gitHubComment.GetID(),
		Body:      gitHubComment.GetBody(),
		User:      user,
		CreatedAt: gitHubComment.GetCreatedAt().Time,
		UpdatedAt: gitHubComment.GetUpdatedAt().Time,
		URL:       gitHubComment.GetURL(),
		HTMLURL:   gitHubComment.GetHTMLURL(),
	}
}