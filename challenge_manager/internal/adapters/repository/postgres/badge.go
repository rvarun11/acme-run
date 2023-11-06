package postgres

import (
	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Repository) CreateBadge(b *domain.Badge) (*domain.Badge, error) {
	pb := &postgresBadge{
		ID:          b.ID,
		PlayerID:    b.PlayerID,
		ChallengeID: b.Challenge.ID,
		CompletedOn: b.CompletedOn,
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&pb).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &domain.Badge{}, err
	}

	badge := pb.toAggregate(b.Challenge)

	return badge, nil
}

func (r *Repository) ListBadgesByPlayerID(pid uuid.UUID) ([]*domain.Challenge, error) {
	var challengesFromDB []postgresChallenge
	// TODO: THIS MAY NOT WORK
	if err := r.db.Where("player_id = ?", pid).Find(&challengesFromDB).Error; err != nil {
		print("this failed as expected")
		return nil, err
	}

	var chs []*domain.Challenge
	for _, pc := range challengesFromDB {

		ch := pc.toAggregate()
		chs = append(chs, ch)
	}

	return chs, nil
}
