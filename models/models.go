package models

type User struct {
	ID       string `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"email"`
}

type UserLogin struct {
	Username string `json:"username" validate:"required,min=4,max=32,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type UserRegistration struct {
	UserLogin
	Email string `json:"email" validate:"required,max=255,email"`
}

type UserList struct {
	Users []User `json:"users"`
}

type UserSession struct {
	ApiKey string `json:"api_key" gorm:"primaryKey"`
	UserId string `json:"-"`
}

type ImageName struct {
	Name string `json:"name" validate:"required,max=255"`
}

type ImageData struct {
	ImageName
	Data string `json:"data,omitempty" validate:"required,base64"`
}

type Image struct {
	ID     string `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId string `json:"-"`
	ImageData
}
