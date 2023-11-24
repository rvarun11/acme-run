package clients

type TrailLocationResponse struct {
	Status         string  `json:"status"`
	Message        string  `json:"message"`
	StartLongitude float64 `json:"start_longitude"`
	StartLatitude  float64 `json:"start_latitude"`
	EndLongitude   float64 `json:"end_longitude"`
	EndLatitude    float64 `json:"end_latitude"`
}
