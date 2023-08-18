package domain

type User struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	IsDelayed bool   `json:"isDelayed"`
	HasError  bool   `json:"hasError"`
}
