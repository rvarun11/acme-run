package postgres

import (
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/domain"
	logger "github.com/CAS735-F23/macrun-teamvsl/challenge/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (r *Repository) CreateBadge(b *domain.Badge) (*domain.Badge, error) {
	pb := &postgresBadge{
		ID:          b.ID,
		PlayerID:    b.PlayerID,
		ChallengeID: b.Challenge.ID,
		CompletedOn: b.CompletedOn,
		Score:       b.Score,
	}
	logger.Info("HEY THE SCORE IN DB IS ", zap.Float64("score", b.Score))
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

func (r *Repository) ListBadgesByPlayerID(pid uuid.UUID) ([]*domain.Badge, error) {
	var badgesFromDB []postgresBadge
	// TODO: THIS MAY NOT WORK
	if err := r.db.Where("player_id = ?", pid).Find(&badgesFromDB).Error; err != nil {
		print("this failed as expected")
		return nil, err
	}
	var badges []*domain.Badge
	for _, pb := range badgesFromDB {
		ch, err := r.GetChallengeByID(pb.ChallengeID)
		if err != nil {
			continue
		}

		badge := pb.toAggregate(ch)
		badges = append(badges, badge)
	}

	return badges, nil
}
