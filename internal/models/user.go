package models

type User struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   int    `json:"status"`
	Level    int    `json:"level"`
}
