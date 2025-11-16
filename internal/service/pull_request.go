package service

import (
	"avito_intership_2025/internal/domain"
	"avito_intership_2025/internal/repo"
	"avito_intership_2025/pkg/postgres/txmanager"
	"context"
	"math/rand"
	"slices"
)

type prService struct {
	tx   *txmanager.TxManager
	pr   repo.PullRequest
	team repo.Team
	user repo.User
}

func newPRService(tx *txmanager.TxManager, pr repo.PullRequest, team repo.Team, user repo.User) *prService {
	return &prService{
		tx:   tx,
		pr:   pr,
		team: team,
		user: user,
	}
}

func (s *prService) Create(ctx context.Context, prID, name, authorID string) (domain.PullRequest, error) {
	var pr domain.PullRequest
	err := s.tx.ExecInTx(ctx, func(ctx context.Context) error {
		author, err := s.user.Find(ctx, authorID)
		if err != nil {
			return err
		}

		members, err := s.team.GetMembers(ctx, author.TeamName)
		if err != nil {
			return err
		}

		rand.Shuffle(len(members), func(i, j int) {
			members[i], members[j] = members[j], members[i]
		})

		active := make([]string, 0, min(2, len(members)))
		for _, m := range members {
			if m.ID != authorID && m.IsActive {
				active = append(active, m.ID)
			}
			if len(active) == 2 {
				break
			}
		}

		pr, err = s.pr.Create(ctx, domain.PullRequest{
			ID:       prID,
			Name:     name,
			AuthorID: authorID,
		})
		if err != nil {
			return err
		}

		if err = s.pr.AddReviewers(ctx, prID, active); err != nil {
			return err
		}

		pr.AssignedReviewers = active
		return nil
	})
	if err != nil {
		return domain.PullRequest{}, err
	}
	return pr, nil
}

func (s *prService) Merge(ctx context.Context, prID string) (domain.PullRequest, error) {
	if err := s.pr.Merge(ctx, prID); err != nil {
		return domain.PullRequest{}, err
	}
	return s.pr.Find(ctx, prID)
}

func (s *prService) Reassign(ctx context.Context, prID, oldReviewer string) (string, domain.PullRequest, error) {
	var (
		pr          domain.PullRequest
		newReviewer string
	)

	err := s.tx.ExecInTx(ctx, func(ctx context.Context) error {
		var err error
		pr, err = s.pr.Find(ctx, prID)
		if err != nil {
			return err
		}
		if pr.Status == domain.PRStatusMerged {
			return domain.ErrPRMerged
		}

		if !slices.Contains(pr.AssignedReviewers, oldReviewer) {
			return domain.ErrReviewerNotAssigned
		}

		members, err := s.team.GetUserMembers(ctx, oldReviewer)
		if err != nil {
			return err
		}

		var candidates []string

		for _, m := range members {
			if m.ID != oldReviewer && m.IsActive && m.ID != pr.AuthorID {
				if !slices.Contains(pr.AssignedReviewers, m.ID) {
					candidates = append(candidates, m.ID)
				}
			}
		}
		if len(candidates) == 0 {
			return domain.ErrNoCandidate
		}

		newReviewer = candidates[rand.Intn(len(candidates))]

		if err = s.pr.ReplaceReviewer(ctx, prID, oldReviewer, newReviewer); err != nil {
			return err
		}

		for i, r := range pr.AssignedReviewers {
			if r == oldReviewer {
				pr.AssignedReviewers[i] = newReviewer
				break
			}
		}
		return nil
	})
	return newReviewer, pr, err
}

func (s *prService) GetUserReview(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	return s.pr.FindUserReview(ctx, userID)
}
