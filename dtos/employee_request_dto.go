package dtos

type EmployeeCreate struct {
	Name     string `json:"name" example:"John Doe"`
	Phone    string `json:"phone" example:"08123456789"`
	Email    string `json:"email" example:"jondoe@example.com"`
	Password string `json:"password" example:"password"`
}

type EmployeeUpdate struct {
	ID       uint   `json:"id" example:"1"`
	Name     string `json:"name" example:"John Doe"`
	Phone    string `json:"phone" example:"08123456789"`
	Email    string `json:"email" example:"jondoe@example.com"`
}

type EmployeeLogin struct {
    Email    string `json:"email" example:"`
    Password string `json:"password" example:"password"`
}
