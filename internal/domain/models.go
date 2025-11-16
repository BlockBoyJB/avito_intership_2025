package domain

import "time"

type User struct {
	ID       string `json:"user_id" db:"id" validate:"required"`
	Username string `json:"username" db:"username" validate:"required"`
	TeamName string `json:"team_name,omitempty" db:"team_name"`
	IsActive bool   `json:"is_active" db:"is_active"`
}

type Team struct {
	Name    string `json:"team_name" validate:"required"`
	Members []User `json:"members" validate:"required,dive"`
}

type PullRequest struct {
	ID                string     `json:"pull_request_id" db:"id"`
	Name              string     `json:"pull_request_name" db:"name"`
	AuthorID          string     `json:"author_id" db:"author_id"`
	Status            string     `json:"status" db:"status"`
	CreatedAt         time.Time  `json:"createdAt" db:"created_at"`
	MergedAt          *time.Time `json:"mergedAt" db:"merged_at"`
	AssignedReviewers []string   `json:"assigned_reviewers,omitempty" db:"reviewers"`
}

const (
	PRStatusOpen   = "OPEN"
	PRStatusMerged = "MERGED"
)

type UserStats struct {
	ReviewerID  string `json:"reviewer_id" db:"reviewer_id"`
	Assignments int    `json:"assignments" db:"assignments"`
}

type AuthorStats struct {
	AuthorID string `json:"author_id" db:"author_id"`
	PRCount  int    `json:"pr_count" db:"pr_count"`
}

type TeamStats struct {
	TeamName string `json:"team_name" db:"team_name"`
	Members  int    `json:"members" db:"members"`
}
