package postgres

import (
	"errors"

	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/domain"
	logger "github.com/CAS735-F23/macrun-teamvsl/challenge/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (r *Repository) CreateOrUpdateChallengeStats(cs *domain.ChallengeStats) error {
	var pcs postgresChallengeStats
	if err := r.db.First(&pcs, "player_id = ? AND challenge_id = ?", cs.PlayerID, cs.Challenge.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Record not found, create new ChallengeStats
			pcs = postgresChallengeStats{
				ID:              uuid.New(),
				PlayerID:        cs.PlayerID,
				ChallengeID:     cs.Challenge.ID,
				DistanceCovered: cs.DistanceCovered,
				EnemiesFought:   cs.EnemiesFought,
				EnemiesEscaped:  cs.EnemiesEscaped,
			}
			if err := r.db.Save(&pcs).Error; err != nil {
				return err
			}
			logger.Debug("challenge stat record not found, creating new one instead", zap.Any("challenge_stat", pcs))
			return nil
		} else {
			logger.Fatal("error occured while finding challenge stats", zap.Error(err))
			return err
		}
	}

	// Increment values in pcs as needed
	newPCS := &postgresChallengeStats{
		ID:              pcs.ID,
		PlayerID:        pcs.PlayerID,
		ChallengeID:     pcs.ChallengeID,
		DistanceCovered: pcs.DistanceCovered + cs.DistanceCovered,
		EnemiesFought:   pcs.EnemiesFought + cs.EnemiesFought,
		EnemiesEscaped:  pcs.EnemiesEscaped + cs.EnemiesEscaped,
	}
	logger.Debug("updating challenge stat record with new values", zap.Any("postgres ds", newPCS))
	if err := r.db.Save(&newPCS).Error; err != nil {
		logger.Fatal("unable to update challenge stats, returning old value instead")
		return err
	}

	return nil
}

func (r *Repository) ListChallengeStatsByChallengeID(cid uuid.UUID) ([]*domain.ChallengeStats, error) {
	var pcs []postgresChallengeStats
	res := r.db.Find(&pcs, "challenge_id = ?", cid)
	if res.Error != nil {
		logger.Error("unable to fetch challenge stats", zap.Any("error", res.Error))
		return []*domain.ChallengeStats{}, res.Error
	}

	ch, err := r.GetChallengeByID(cid)
	if err != nil {
		return []*domain.ChallengeStats{}, err
	}

	var csArr []*domain.ChallengeStats
	for _, pc := range pcs {
		cs := pc.toAggregate(ch)
		csArr = append(csArr, cs)
	}
	return csArr, nil
}

func (r *Repository) DeleteChallengeStats(pid uuid.UUID, cid uuid.UUID) error {
	var pcs postgresChallengeStats
	res := r.db.First(&pcs, "player_id = ? AND challenge_id = ?", pid, cid)
	if res.Error != nil {
		return res.Error
	}
	r.db.Delete(&pcs)
	return nil
}
