package clients

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type PeripheralDeviceClientImpl struct {
}

// Factory for creating a NewPeripheralDeviceClient
func NewPeripheralDeviceClient() *PeripheralDeviceClientImpl {
	return &PeripheralDeviceClientImpl{}
}

func (p *PeripheralDeviceClientImpl) BindPeripheralData(trailID uuid.UUID, playerID uuid.UUID, workoutID uuid.UUID, hrmID uuid.UUID, hrmConnected bool, SendLiveLocationToTrailManager bool) error {
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

	url := "http://localhost:8004/api/v1/peripheral_bind"
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

func (p *PeripheralDeviceClientImpl) UnbindPeripheralData(workoutID uuid.UUID) error {
	// Prepare the data for the PUT request
	unbindData := UnbindPeripheralData{
		WorkoutID: workoutID,
	}

	// Marshal the data for the request body
	unbindPayload, err := json.Marshal(unbindData)
	if err != nil {
		return err
	}

	url := "http://localhost:8004/api/v1/peripheral_unbind"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(unbindPayload))
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

func (p *PeripheralDeviceClientImpl) GetAverageHeartRateOfUser(workoutID uuid.UUID) (uint8, error) {

	url := "http://localhost:8004/api/v1/peripheral/hrm/workout_id=" + workoutID.String() + "&type=avg"
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

	var averageHeartRate uint8
	err = json.NewDecoder(resp.Body).Decode(&averageHeartRate)
	if err != nil {
		return 0, err
	}

	return averageHeartRate, nil
}
