package twilio

type Twilio struct {
	ID    string `json:"id" db:"id"`
	SID   string `json:"sid"`
	Token string `json:"token"`
}
