package request

type UserCreate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
