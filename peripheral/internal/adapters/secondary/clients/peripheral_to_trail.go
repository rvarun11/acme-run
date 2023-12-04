package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

type ZoneClient struct {
	clientURL string
}

// Factory for creating a new WorkoutService
func NewZoneServiceClient(cfg string) *ZoneClient {
	return &ZoneClient{
		clientURL: cfg,
	}
}

// GetTrailLocation makes a GET request to the GetTrailLocationInfo endpoint and saves the location data
func GetTrailLocation(serverURL, zoneID, trailID string) (*TrailLocationResponse, error) {
	// Construct the URL with the zone and trail IDs
	url := fmt.Sprintf("%s/zone/%s/trail/%s", serverURL, zoneID, trailID)

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-200 status code: %d", resp.StatusCode)
	}

	// Read the body of the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON response into a TrailLocationResponse struct
	var locationResponse TrailLocationResponse
	err = json.Unmarshal(body, &locationResponse)
	if err != nil {
		return nil, err
	}

	return &locationResponse, nil
}

func (z *ZoneClient) GetTrailLocation(trailID uuid.UUID) (float64, float64, float64, float64, error) {

	// Create a new GET request to fetch the user's age.
	serverURL := z.clientURL + "/api/v1"
	zoneID := uuid.New().String()

	locationInfo, err := GetTrailLocation(serverURL, zoneID, trailID.String())
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return locationInfo.StartLongitude, locationInfo.StartLatitude, locationInfo.EndLongitude, locationInfo.EndLatitude, nil
}
