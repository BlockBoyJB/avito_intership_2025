package domain

import "errors"

var (
	ErrTeamExists          = errors.New("team exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrPRExists            = errors.New("PR id already exists")
	ErrPRNotFound          = errors.New("PR not found")
	ErrPRMerged            = errors.New("cannot reassign on merged PR")
	ErrReviewerNotAssigned = errors.New("reviewer not assigned to this PR")
	ErrNoCandidate         = errors.New("no active replacement candidate in team")
)

var (
	ErrNotFound = errors.Join(ErrUserNotFound, ErrPRNotFound)
)
