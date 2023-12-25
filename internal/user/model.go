package user

type User struct {
	// Уникальный идентификатор пользователя
	UserId int32 `json:"user_id"`
	// Имя пользователя
	Username string `json:"username,omitempty"`
	// Адрес электронной почты пользователя
	Email string `json:"email"`
}
