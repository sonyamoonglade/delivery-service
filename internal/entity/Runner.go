package entity

type Runner struct {
	RunnerID    int64  `json:"runner_id,omitempty"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
}
