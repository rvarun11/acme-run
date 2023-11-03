package postgres

import (
	"fmt"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/core/domain"
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
	db.AutoMigrate(&postgresChallenge{})

	return &Repository{
		db: db,
	}
}

// TODO: unique constraint not working here
type postgresChallenge struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	Name        string    `gorm:"unique;not null"`
	Description string    `gorm:"not null"`
	BadgeURL    string
	Criteria    string `gorm:"not null"`
	Goal        float32
	Start       time.Time
	End         time.Time
	CreatedAt   time.Time
}

// Repository Functions

func (r *Repository) Create(ch *domain.Challenge) (*domain.Challenge, error) {
	pc := &postgresChallenge{
		ID:          ch.ID,
		Name:        ch.Name,
		Description: ch.Description,
		BadgeURL:    ch.BadgeURL,
		Criteria:    string(ch.Criteria),
		Goal:        ch.Goal,
		CreatedAt:   time.Now(),
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

	p := toAggregate(pc)

	return p, nil
}

func (r *Repository) GetByID(cid uuid.UUID) (*domain.Challenge, error) {
	var pc postgresChallenge
	res := r.db.First(&pc, "id = ?", cid)
	if res.Error != nil {
		return &domain.Challenge{}, res.Error
	}

	ch := toAggregate(&pc)
	return ch, nil
}

func (r *Repository) Update(req *domain.Challenge) (*domain.Challenge, error) {
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

	ch := toAggregate(pc)

	return ch, nil
}

// Helper for converting to/from domain model
func toAggregate(pc *postgresChallenge) *domain.Challenge {
	return &domain.Challenge{
		ID:          pc.ID,
		Name:        pc.Name,
		Description: pc.Description,
		Criteria:    domain.Criteria(pc.Criteria),
		Goal:        pc.Goal,
		Start:       pc.Start,
		End:         pc.End,
		BadgeURL:    pc.BadgeURL,
		CreatedAt:   pc.CreatedAt,
	}
}

func (r *Repository) List() ([]*domain.Challenge, error) {
	var challengesFromDB []postgresChallenge
	if err := r.db.Find(&challengesFromDB).Error; err != nil {
		return nil, err
	}

	var chs []*domain.Challenge
	for _, pc := range challengesFromDB {

		ch := toAggregate(&pc)
		chs = append(chs, ch)
	}

	return chs, nil
}
