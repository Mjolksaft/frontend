package session

type UserSession struct {
	UserID    string
	Username  string
	AuthToken string // e.g., for JWT or cookies
}

func NewUserSession(userID, username, authToken string) *UserSession {
	return &UserSession{
		UserID:    userID,
		Username:  username,
		AuthToken: authToken,
	}
}
