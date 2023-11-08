package postgres

import (
	"errors"

	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/core/domain"
	logger "github.com/CAS735-F23/macrun-teamvsl/challenge_manager/log"
	"github.com/google/uuid"
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
			// newCS := pcs.toAggregate(cs.Challenge)
			logger.Debug("challenge stats record not found, creating new one instead")
			return nil
		} else {
			logger.Fatal("error occured while finding challenge stats")
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
		EnemiesEscaped:  pcs.EnemiesEscaped + cs.EnemiesFought,
	}

	// Save or update ChallengeStats in the database
	if err := r.db.Save(&newPCS).Error; err != nil {
		logger.Fatal("unable to update challenge stats, returning old value instead")
		return err
	}
	// newCS := newPCS.toAggregate(cs.Challenge)
	return nil
}

func (r *Repository) ListChallengeStatsByPlayerID(pid uuid.UUID) ([]*domain.ChallengeStats, error) {
	var pcs []postgresChallengeStats
	res := r.db.First(&pcs, "player_id = ?", pid)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return []*domain.ChallengeStats{}, nil
		}
		return []*domain.ChallengeStats{}, res.Error
	}

	var csArr []*domain.ChallengeStats
	for _, pc := range pcs {
		ch, err := r.GetChallengeByID(pc.ChallengeID)
		if err != nil {
			return []*domain.ChallengeStats{}, err
		}
		cs := pc.toAggregate(ch)
		csArr = append(csArr, cs)
	}

	return csArr, nil
}

func (r *Repository) ListChallengeStatsByChallengeID(cid uuid.UUID) ([]*domain.ChallengeStats, error) {
	var pcs []postgresChallengeStats
	res := r.db.First(&pcs, "challenge_id = ?", cid)
	if res.Error != nil {
		return []*domain.ChallengeStats{}, res.Error
	}

	var csArr []*domain.ChallengeStats
	for _, pc := range pcs {
		ch, err := r.GetChallengeByID(pc.ChallengeID)
		if err != nil {
			return []*domain.ChallengeStats{}, err
		}
		cs := pc.toAggregate(ch)
		csArr = append(csArr, cs)
	}
	return csArr, nil
}

// func (r *Repository) ListEligibleChallengeStatsForChallenge(ch *domain.Challenge) ([]*domain.ChallengeStats, error) {
// 	var pcs []postgresChallengeStats
// 	res := r.db.First(&pcs, "challenge_id = ?", cid)
// 	if res.Error != nil {
// 		return []*domain.ChallengeStats{}, res.Error
// 	}

// 	var csArr []*domain.ChallengeStats
// 	for _, pc := range pcs {
// 		ch, err := r.GetChallengeByID(pc.ChallengeID)
// 		if err != nil {
// 			return []*domain.ChallengeStats{}, err
// 		}
// 		cs := pc.toAggregate(ch)
// 		csArr = append(csArr, cs)
// 	}
// 	return csArr, nil
// }

// TODO: Once the other parts work, add this
func (r *Repository) DeleteChallengeStats(pid uuid.UUID, cid uuid.UUID) error {
	return nil
}
