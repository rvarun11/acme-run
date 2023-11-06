package postgres

import (
	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository Functions

func (r *Repository) CreateChallenge(ch *domain.Challenge) (*domain.Challenge, error) {
	pc := &postgresChallenge{
		ID:          ch.ID,
		Name:        ch.Name,
		Description: ch.Description,
		BadgeURL:    ch.BadgeURL,
		Criteria:    string(ch.Criteria),
		Goal:        ch.Goal,
		CreatedAt:   ch.CreatedAt,
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&pc).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &domain.Challenge{}, err
	}

	p := pc.toAggregate()

	return p, nil
}

func (r *Repository) GetChallengeByID(cid uuid.UUID) (*domain.Challenge, error) {
	var pc postgresChallenge
	res := r.db.First(&pc, "id = ?", cid)
	if res.Error != nil {
		return &domain.Challenge{}, res.Error
	}

	ch := pc.toAggregate()
	return ch, nil
}

func (r *Repository) UpdateChallenge(req *domain.Challenge) (*domain.Challenge, error) {
	var pc *postgresChallenge
	if err := r.db.First(&pc, "id = ?", req.ID).Error; err != nil {
		return &domain.Challenge{}, err
	}

	pc.Name = req.Name
	pc.Description = req.Description
	pc.BadgeURL = req.BadgeURL
	pc.Criteria = string(req.Criteria)
	pc.Goal = req.Goal
	pc.Start = req.Start
	pc.End = req.End

	tx := r.db.Begin()
	if err := tx.Save(&pc).Error; err != nil {
		tx.Rollback()
		return &domain.Challenge{}, err
	}
	tx.Commit()

	ch := pc.toAggregate()

	return ch, nil
}

func (r *Repository) ListChallenges() ([]*domain.Challenge, error) {
	var challengesFromDB []postgresChallenge
	if err := r.db.Find(&challengesFromDB).Error; err != nil {
		return nil, err
	}

	var chs []*domain.Challenge
	for _, pc := range challengesFromDB {

		ch := pc.toAggregate()
		chs = append(chs, ch)
	}

	return chs, nil
}
