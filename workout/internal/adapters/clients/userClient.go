package clients

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type UserServiceClientImpl struct {
}

// Factory for creating a new WorkoutService
func NewUserServiceClient() *UserServiceClientImpl {
	return &UserServiceClientImpl{}
}

func (u *UserServiceClientImpl) GetProfileOfUser(playerID uuid.UUID) (string, error) {
	// Make the GET request for user profile
	url := "YOUR_UNBIND_ENDPOINT/" // Modify the endpoint as required

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

	// Process the response body to extract profile
	var userProfile string
	err = json.NewDecoder(resp.Body).Decode(&userProfile)
	if err != nil {
		return "", err
	}

	return userProfile, nil
}

func (u *UserServiceClientImpl) GetHardcoreModeOfUser(playerID uuid.UUID) (bool, error) {
	// Make the GET request for user's hardcore mode
	url := "YOUR_UNBIND_ENDPOINT/" // Modify the endpoint as required

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Process the response body to extract hardcore mode
	var hardcoreMode string
	err = json.NewDecoder(resp.Body).Decode(&hardcoreMode)
	if err != nil {
		return false, err
	}

	return hardcoreMode == "true", nil
}

func (u *UserServiceClientImpl) GetUserAge(playerID uuid.UUID) (uint8, error) {

	url := "http://example.com/get_user_age/" + playerID.String()

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

	var age uint8
	err = json.NewDecoder(resp.Body).Decode(&age)
	if err != nil {
		return 0, err
	}

	return age, nil
}
