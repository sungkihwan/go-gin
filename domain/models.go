package domain

type User struct {
	ID        *uint  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	IsDelayed *bool  `json:"isDelayed,omitempty"`
	HasError  *bool  `json:"hasError,omitempty"`
}

type XMLRoot struct {
	Maps []XMLMap `xml:"map"`
}

type XMLMap struct {
	ID  string `xml:"id,attr"`
	Key string `xml:"key"`
}
