package dto

type LoginResponse struct {
	User        User         `json:"user"`
	Restaurants []Restaurant `json:"restaurants"`
}
