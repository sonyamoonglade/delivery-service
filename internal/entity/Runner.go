package entity

type Runner struct {
	RunnerID    int64  `json:"runner_id,omitempty" db:"runner_id"`
	Username    string `json:"username" db:"username"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
}
