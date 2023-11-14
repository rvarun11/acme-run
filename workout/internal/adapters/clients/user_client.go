package clients

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func calculateAge(dob string) (int, error) {
	birthdate, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return 0, err
	}

	now := time.Now()

	// Calculate the age
	age := now.Year() - birthdate.Year()

	// Adjust if this year's birthday has not occurred yet
	if now.Month() < birthdate.Month() || (now.Month() == birthdate.Month() && now.Day() < birthdate.Day()) {
		age--
	}

	return age, nil
}

type UserServiceClientImpl struct {
}

// Factory for creating a new WorkoutService
func NewUserServiceClient() *UserServiceClientImpl {
	return &UserServiceClientImpl{}
}

func (u *UserServiceClientImpl) GetWorkoutPreferenceOfUser(playerID uuid.UUID) (string, error) {

	url := "localhost:8000/api/v1/players/" + playerID.String()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Process the response body to extract preference
	var playerDTO playerDTO
	err = json.NewDecoder(resp.Body).Decode(&playerDTO)
	if err != nil {
		return "", err
	}
	return playerDTO.Preference, nil
}

func (u *UserServiceClientImpl) GetUserAge(playerID uuid.UUID) (uint8, error) {

	url := "localhost:8000/api/v1/players/" + playerID.String()

	// Create a new GET request to fetch the user's age.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Process the response body to extract preference
	var playerDTO playerDTO
	err = json.NewDecoder(resp.Body).Decode(&playerDTO)
	if err != nil {
		return 0, err
	}

	age, err := calculateAge(playerDTO.User.DateOfBirth)

	if err != nil {
		return 0, err
	}

	return uint8(age), nil
}
