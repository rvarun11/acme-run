// Package services holds all the services that connects repositories into a business flow
package service

import (
	player "github.com/rvarun11/macrun-teamvs/domain/player"
	"github.com/rvarun11/macrun-teamvs/domain/player/memory"
)

// PlayerConfiguration is an alias for a function that will take in a pointer to an PlayerService and modify it
type PlayerConfiguration func(os *PlayerService) error

// OrderService is a implementation of the OrderService
type PlayerService struct {
	players player.PlayerRepository
}

// NewPlayerService takes a variable amount of PlayerConfiguration functions and returns a new PlayerService
// Each PlayerConfiguration will be called in the order they are passed in
func NewOrderService(cfgs ...PlayerConfiguration) (*PlayerService, error) {
	// Create the PlayerService
	os := &PlayerService{}
	// Apply all Configurations passed in
	for _, cfg := range cfgs {
		// Pass the service into the configuration function
		err := cfg(os)
		if err != nil {
			return nil, err
		}
	}
	return os, nil
}

// WithPlayerRepository applies a given customer repository to the OrderService
func WithPlayerRepository(pr player.PlayerRepository) PlayerConfiguration {
	// return a function that matches the PlayerConfiguration alias,
	// You need to return this so that the parent function can take in all the needed parameters
	return func(ps *PlayerService) error {
		ps.players = pr
		return nil
	}
}

// WithMemoryCustomerRepository applies a memory customer repository to the PlayerService
func WithMemoryCustomerRepository() PlayerConfiguration {
	// Create the memory repo, if we needed parameters, such as connection strings they could be inputted here
	pr := memory.New()
	return WithPlayerRepository(pr)
}
