package main

// Admin struct represents the structure of an admin
type Admin struct {
	GoogleID      string  `json:"googleID"`
	Name          string  `json:"name"`
	ContactNumber string  `json:"contactNumber"`
	Email         string  `json:"email"`
	BasePay       float64 `json:"basePay"`
	Incentive     int     `json:"incentive"`
	Role          string  `json:"role"`
}
