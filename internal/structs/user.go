package structs

// form нужен для парсинга query параметров с помощью gin.Bind
type User struct {
	Uid string `json:"uid"`
	Username string `json:"username" binding:"required"`
	Name string `json:"name" binding:"required"`
	Surname string `json:"surname" binding:"required"`
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}