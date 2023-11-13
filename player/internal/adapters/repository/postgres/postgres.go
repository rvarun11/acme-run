package postgres

import (
	"fmt"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/player/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository() *Repository {
	host := conf.host
	port := conf.port
	user := conf.user
	password := conf.password
	dbname := conf.dbname
	encoding := conf.encoding

	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable client_encoding=%s",
		host, port, user, dbname, password, encoding,
	)

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&postgresUser{}, &postgresPlayer{})
	// db.Model(&postgresPlayer{}).Association("user_id")

	return &Repository{
		db: db,
	}
}

// Repository Types

type postgresUser struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`
	// Name of the user
	Name string
	// Email
	Email string `gorm:"unique;not null"`
	// Date of Birth of the player TODO: Fix type
	DateOfBirth string
}

type postgresPlayer struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key"`
	// User is the root entity of player
	UserID uuid.UUID `gorm:"type:uuid;unique not null"`
	// Weight of the player
	Weight float64 `gorm:"<-"`
	// Height of the player
	Height float64 `gorm:"<-"`
	// GeographicalZone is a group of trails in a region
	ZoneID uuid.UUID
	// CreatedAt is the time when the player registered
	CreatedAt time.Time
	// UpdatedAt is the time when the player last updated the profile
	UpdatedAt time.Time
}

// Repository Functions

func (r *Repository) Create(player *domain.Player) (*domain.Player, error) {
	pu := &postgresUser{
		ID:          player.User.ID,
		Name:        player.User.Name,
		Email:       player.User.Email,
		DateOfBirth: player.User.DateOfBirth,
	}

	pp := &postgresPlayer{
		ID:        player.ID,
		UserID:    player.User.ID,
		Weight:    player.Weight,
		Height:    player.Height,
		ZoneID:    player.ZoneID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
	println(pid.String())
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

// Helper for converting to/from domain model

func toAggregate(pu *postgresUser, pp *postgresPlayer) *domain.Player {
	return &domain.Player{
		ID: pp.ID,
		User: &domain.User{
			ID:          pu.ID,
			Email:       pu.Email,
			Name:        pu.Name,
			DateOfBirth: pu.DateOfBirth,
		},
		Weight:    pp.Weight,
		Height:    pp.Height,
		ZoneID:    pp.ZoneID,
		CreatedAt: pp.CreatedAt,
		UpdatedAt: pp.UpdatedAt,
	}
}

// May be needed later
// func fromAggregate(player *domain.Player) (*postgresUser, *postgresPlayer) {
// 	pu := &postgresUser{
// 		ID:          player.User.ID,
// 		Name:        player.User.Name,
// 		Email:       player.User.Email,
// 		DateOfBirth: player.User.DateOfBirth,
// 	}

// 	pp := &postgresPlayer{
// 		ID:        player.ID,
// 		UserID:    player.User.ID,
// 		Weight:    player.Weight,
// 		Height:    player.Height,
// 		ZoneID:    player.ZoneID,
// 		CreatedAt: player.CreatedAt,
// 		UpdatedAt: player.CreatedAt,
// 	}

// 	return pu, pp
// }
