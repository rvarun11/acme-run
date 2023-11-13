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

func (p *PeripheralDeviceClientImpl) BindPeripheralData(playerID uuid.UUID, workoutID uuid.UUID, hrmID uuid.UUID, SendLiveLocationToTrailManager bool) error {
	// Prepare the data for the POST request
	bindData := BindPeripheralData{
		PlayerID:                       playerID,
		WorkoutID:                      workoutID,
		HRMId:                          hrmID,
		SendLiveLocationToTrailManager: SendLiveLocationToTrailManager,
	}

	bindPayload, err := json.Marshal(bindData)
	if err != nil {
		return err
	}

	// Make the POST request
	// TODO
	url := "YOUR_PERIPHERAL_SERVICE_ENDPOINT" // Replace with your actual endpoint
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

	// Make the PUT request
	// TODO: Replace with your actual endpoint
	url := "YOUR_UNBIND_ENDPOINT/" + workoutID.String() // Modify the endpoint as required
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

func (p *PeripheralDeviceClientImpl) GetAverageHeartRateOfUser(workoutID uuid.UUID) (uint8, error) {
	// Make the GET request
	// TODO: Replace with your actual endpoint
	url := "YOUR_AVERAGE_HEARTRATE_ENDPOINT/" + workoutID.String() // Modify the endpoint as required
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

	// Parse the response to retrieve average heart rate
	// Assuming response body contains the average heart rate as uint8
	var averageHeartRate uint8
	err = json.NewDecoder(resp.Body).Decode(&averageHeartRate)
	if err != nil {
		return 0, err
	}

	return averageHeartRate, nil
}
