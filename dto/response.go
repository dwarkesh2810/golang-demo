package dto

type UserDetailsResponse struct {
	Id          int    `json:"user_id"`
	FirstName   string `json:"first_name" `
	LastName    string `json:"last_name" `
	Email       string `json:"email"`
	Country     int    `json:"country_id"`
	CreatedDate string `json:"CreatedDate"`
}
