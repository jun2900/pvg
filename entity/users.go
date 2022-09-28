package entity

type UserRegisterInput struct {
	Username    string `json:"username"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Birthday    string `json:"birthday"`
}

type UserUpdateInput struct {
	Username        string `json:"username"`
	FirstName       string `json:"firstname"`
	LastName        string `json:"lastname"`
	UpdatePassword  string `json:"update_password"`
	PhoneNumber     string `json:"phone_number"`
	Email           string `json:"email"`
	Birthday        string `json:"birthday"`
	CurrentPassword string `json:"current_password"`
}

type CheckUserExist struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type OtpInput struct {
	Email    string `json:"email"`
	Passcode string `json:"passcode"`
}

type BasicResponse struct {
	Message string `json:"message"`
}
