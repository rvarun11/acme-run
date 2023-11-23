package postgres

import (
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/user/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository Functions

func (r *Repository) Create(player *domain.Player) (*domain.Player, error) {
	pu := &postgresUser{
		ID:          player.User.ID,
		Name:        player.User.Name,
		Email:       player.User.Email,
		DateOfBirth: player.User.DateOfBirth,
	}

	pp := &postgresPlayer{
		ID:         player.ID,
		UserID:     player.User.ID,
		Weight:     player.Weight,
		Height:     player.Height,
		Preference: string(player.Preference),
		ZoneID:     player.ZoneID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	// pu, pp := fromAggregate(player)
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&pu).Error; err != nil {
			return err
		}
		if err := tx.Create(&pp).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &domain.Player{}, err
	}

	p := toAggregate(pu, pp)

	return p, nil
}

func (r *Repository) Get(pid uuid.UUID) (*domain.Player, error) {
	var pp postgresPlayer
	var pu postgresUser
	res := r.db.First(&pp, "id = ?", pid)
	if res.Error != nil {
		return &domain.Player{}, res.Error
	}

	res = r.db.First(&pu, "id = ?", pp.UserID)
	if res.Error != nil {
		return &domain.Player{}, res.Error
	}

	player := toAggregate(&pu, &pp)
	return player, nil
}

func (r *Repository) GetByEmail(email string) (*domain.Player, error) {
	var pp postgresPlayer
	var pu postgresUser

	res := r.db.First(&pu, "email = ?", email)
	if res.Error != nil {
		return &domain.Player{}, res.Error
	}

	res = r.db.First(&pp, pu.ID)
	if res.Error != nil {
		return &domain.Player{}, res.Error
	}

	player := toAggregate(&pu, &pp)

	return player, nil
}

func (r *Repository) Update(player *domain.Player) (*domain.Player, error) {
	var pu postgresUser
	var pp postgresPlayer
	if err := r.db.First(&pu, "id = ?", player.User.ID).Error; err != nil {
		return &domain.Player{}, err
	}

	if err := r.db.First(&pp, "user_id = ?", player.User.ID).Error; err != nil {
		return &domain.Player{}, err
	}

	pu.Name = player.User.Name
	pu.Email = player.User.Email
	pu.DateOfBirth = player.User.DateOfBirth
	pp.Weight = player.Weight
	pp.Height = player.Height
	pp.Preference = string(player.Preference)
	pp.ZoneID = player.ZoneID
	pp.UpdatedAt = time.Now()

	tx := r.db.Begin()
	if err := tx.Save(&pu).Error; err != nil {
		tx.Rollback()
		return &domain.Player{}, err
	}

	if err := tx.Save(&pp).Error; err != nil {
		tx.Rollback()
		return &domain.Player{}, err
	}

	tx.Commit()
	player = toAggregate(&pu, &pp)

	return player, nil
}

func (r *Repository) List() ([]*domain.Player, error) {
	var playersFromDB []postgresPlayer
	if err := r.db.Find(&playersFromDB).Error; err != nil {
		return nil, err
	}

	var players []*domain.Player
	for _, pp := range playersFromDB {
		var pu postgresUser
		if err := r.db.First(&pu, "id = ?", pp.UserID).Error; err != nil {
			return nil, err
		}

		player := toAggregate(&pu, &pp)
		players = append(players, player)
	}

	return players, nil
}
