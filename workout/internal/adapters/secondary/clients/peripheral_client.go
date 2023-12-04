package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

type PeripheralClientImpl struct {
	clientURL string
}

// Factory for creating a NewPeripheralClient
func NewPeripheralClient(cfg string) *PeripheralClientImpl {
	return &PeripheralClientImpl{
		clientURL: cfg,
	}
}

func (p *PeripheralClientImpl) BindPeripheralData(trailID uuid.UUID, playerID uuid.UUID, workoutID uuid.UUID, hrmID uuid.UUID, hrmConnected bool, SendLiveLocationToTrailManager bool) error {
	// Prepare the data for the POST request
	bindData := BindPeripheralData{
		TrailID:                        trailID,
		PlayerID:                       playerID,
		WorkoutID:                      workoutID,
		HRMId:                          hrmID,
		HRMConnected:                   hrmConnected,
		SendLiveLocationToTrailManager: SendLiveLocationToTrailManager,
	}

	bindPayload, err := json.Marshal(bindData)
	if err != nil {
		return err
	}

	url := p.clientURL + "/api/v1/peripheral"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bindPayload))

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}

func (p *PeripheralClientImpl) UnbindPeripheralData(workoutID uuid.UUID) error {
	// Prepare the data for the PUT request
	unbindData := UnbindPeripheralData{
		WorkoutID: workoutID,
	}

	// Marshal the data for the request body
	unbindPayload, err := json.Marshal(unbindData)
	if err != nil {
		return err
	}

	url := p.clientURL + "/api/v1/peripheral"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(unbindPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}

func (p *PeripheralClientImpl) GetAverageHeartRateOfUser(workoutID uuid.UUID) (uint8, error) {
	// Ensure workoutID is valid
	if workoutID == uuid.Nil {
		return 0, errors.New("invalid workout ID")
	}

	url := p.clientURL + "/api/v1/peripheral/hrm?workout_id=" + workoutID.String() + "&type=avg"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	if resp == nil || resp.Body == nil {
		return 0, errors.New("received nil response or nil body")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var averageHeartRate AverageHeartRate
	err = json.Unmarshal(body, &averageHeartRate)
	if err != nil {
		return 0, err
	}

	return averageHeartRate.AverageHeartRate, nil
}
