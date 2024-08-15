package domain

//Тут в будущем могли бы быть модели пользователей

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	// Refresh RefreshToken `json:"refresh_token"`
}
